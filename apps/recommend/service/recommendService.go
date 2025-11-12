package service

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/GoFurry/gofurry-game-backend/apps/recommend/dao"
	"github.com/GoFurry/gofurry-game-backend/apps/recommend/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/log"
	cs "github.com/GoFurry/gofurry-game-backend/common/service"
	"github.com/GoFurry/gofurry-game-backend/common/util"

	gd "github.com/GoFurry/gofurry-game-backend/apps/game/dao"
	gm "github.com/GoFurry/gofurry-game-backend/apps/game/models"
)

type recommendService struct{}

var recommendSingleton = new(recommendService)

func GetRecommendService() *recommendService { return recommendSingleton }

// Redis 定义
const (
	redisTagMappingKey = "recommend:tag-mapping"
	redisTagIDsKey     = "recommend:tag-ids"
	cacheExpireTime    = 1 * time.Hour // 缓存过期时间
)

var once sync.Once

// 随机获取一个 GameID
func (s recommendService) GetRandomGameID() (string, common.GFError) {
	count, err := gd.GetGameDao().Count(gm.GfgGame{})
	if err != nil {
		return "", common.NewServiceError("统计数量出错")
	}
	intCnt, _ := util.String2Int(util.Int642String(count))

	// Go 1.20 + 无需初始化
	once.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	// 生成一个 0 到 intCnt-1 之间的随机整数
	randomInt := rand.Intn(intCnt)
	gameRecord, gfError := gd.GetGameDao().GetByNum(randomInt)
	if gfError != nil {
		return "", gfError
	}

	return util.Int642String(gameRecord.ID), nil
}

// Content-based Filter 返回物品A的余弦相似度最高的物品
func (s recommendService) RecommendByCBF(id string, lang string) (gameListVo []gm.GameRespVo, err common.GFError) {
	intID, parseErr := util.String2Int64(id)
	if parseErr != nil {
		return nil, common.NewServiceError(parseErr.Error())
	}
	if gameListVo, err = getGameCBF(intID, lang); err != nil {
		return nil, err
	}
	return
}

// CBF 获取一组推荐的游戏记录
func getGameCBF(id int64, lang string) (recommendContent []gm.GameRespVo, err common.GFError) {
	// 执行 CBF
	start := time.Now()
	similarities, err := processContentBasedFilter(id)
	if err != nil {
		return nil, err
	}
	fmt.Println(time.Since(start))

	// 从相似度结果生成推荐视图 前12随机选8
	const topN = 8
	const candidateN = 12

	filtered := make([]models.ContentSimilarities, 0, candidateN)
	for _, sim := range similarities {
		if sim.ID == id || sim.Similarity <= 0 {
			continue
		}
		filtered = append(filtered, sim)
		if len(filtered) >= candidateN {
			break
		}
	}

	// 打乱候选列表 增加多样性
	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})
	if len(filtered) > topN {
		filtered = filtered[:topN]
	}

	// 转换为GameRespVo
	// TODO:暂不修改

	return recommendContent, nil
}

// CBF 算法
func processContentBasedFilter(gameID int64) ([]models.ContentSimilarities, common.GFError) {
	// 获取标签映射和标签ID列表
	tagMappingMap, tagIDs, err := getTagToMap()
	if err != nil {
		return nil, err
	}

	// 初始化标签ID到维度索引的映射
	tagIDToIndex := buildTagIndexMap(tagIDs)

	// 特征提取 - 独热编码
	targetContent, contentFeatures := execFeature(tagMappingMap, tagIDToIndex, gameID)

	// 校验目标游戏是否存在有效特征
	if len(targetContent.Tag) == 0 {
		// 游戏ID不在映射中
		if _, exists := tagMappingMap[gameID]; !exists {
			return nil, common.NewServiceError("目标游戏不存在或未关联标签")
		}
		// 游戏存在但无标签
		log.Warn("游戏ID=", gameID, "未关联任何标签，无法生成推荐")
		return []models.ContentSimilarities{}, nil
	}

	// 计算相似度
	similarities := execSimilarity(targetContent, contentFeatures)
	return similarities, nil
}

// 构建标签 ID到维度索引的映射 map[tagID]index
func buildTagIndexMap(tagIDs []int64) map[int64]int {
	tagIDToIndex := make(map[int64]int, len(tagIDs))
	for idx, tagID := range tagIDs {
		tagIDToIndex[tagID] = idx
	}
	return tagIDToIndex
}

// 特征提取 - 独热编码
func execFeature(tagMapping map[int64][]int64, tagIDToIndex map[int64]int, targetGameID int64) (models.ContentSimilarities, []models.ContentSimilarities) {
	var targetContent models.ContentSimilarities
	var contentFeatures []models.ContentSimilarities

	for gameID, tagIDs := range tagMapping {
		// 构建独热特征
		feature := make([]float64, 0, len(tagIDs))
		seen := make(map[int]struct{}) // 去重标签

		for _, tagID := range tagIDs {
			idx, ok := tagIDToIndex[tagID]
			if !ok {
				continue // 忽略未注册的标签
			}
			if _, exists := seen[idx]; exists {
				continue // 跳过重复标签
			}
			seen[idx] = struct{}{}
			feature = append(feature, float64(idx)) // 存储维度索引
		}

		// 区分目标游戏和其他游戏
		if gameID == targetGameID {
			targetContent = models.ContentSimilarities{
				ID:  gameID,
				Tag: feature,
			}
		} else {
			contentFeatures = append(contentFeatures, models.ContentSimilarities{
				ID:  gameID,
				Tag: feature,
			})
		}
	}

	return targetContent, contentFeatures
}

// 计算相似度
func execSimilarity(target models.ContentSimilarities, others []models.ContentSimilarities) []models.ContentSimilarities {
	similarities := make([]models.ContentSimilarities, 0, len(others))

	// 将目标特征转换为字典
	targetSet := make(map[float64]struct{}, len(target.Tag))
	for _, idx := range target.Tag {
		targetSet[idx] = struct{}{}
	}

	// 计算每个游戏与目标的相似度
	for _, other := range others {
		// 跳过没有特征的项
		if len(other.Tag) == 0 {
			continue
		}

		// 计算共同标签数量
		commonCount := 0
		for _, idx := range other.Tag {
			if _, exists := targetSet[idx]; exists {
				commonCount++
			}
		}

		// 计算余弦相似度 commonCount / (sqrt(len(target)) * sqrt(len(other)))
		magTarget := math.Sqrt(float64(len(target.Tag)))
		magOther := math.Sqrt(float64(len(other.Tag)))
		sim := float64(commonCount) / (magTarget * magOther)

		if sim > 0 {
			similarities = append(similarities, models.ContentSimilarities{
				ID:         other.ID,
				Similarity: sim,
			})
		}
	}

	// 相似度排序
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	return similarities
}

// 获取标签映射
func getTagToMap() (tagMapping map[int64][]int64, tagIDs []int64, err common.GFError) {
	// Redis 读缓存
	tagMapping, tagIDs, err = loadFromRedis()
	if err == nil && tagMapping != nil && len(tagIDs) > 0 {
		// 缓存命中，直接返回
		return tagMapping, tagIDs, nil
	}

	// 缓存未命中
	mappingRecords, err := dao.GetRecommendDao().GetTagMappingList()
	if err != nil {
		return nil, nil, common.NewServiceError("获取标签映射记录失败: " + err.GetMsg())
	}
	tagMapping = make(map[int64][]int64)
	for _, rec := range mappingRecords {
		gameID := rec.GameID
		tagID := rec.TagID

		tags := tagMapping[gameID]
		exists := false
		for _, t := range tags {
			if t == tagID {
				exists = true
				break
			}
		}
		if !exists {
			tagMapping[gameID] = append(tags, tagID)
		}
	}

	tagRecords, err := dao.GetRecommendDao().GetTagList()
	if err != nil {
		return nil, nil, common.NewServiceError("获取标签记录失败: " + err.GetMsg())
	}
	tagIDs = make([]int64, 0, len(tagRecords))
	for idx := range tagRecords {
		tagIDs = append(tagIDs, tagRecords[idx].ID)
	}

	// 异步写入Redis缓存
	go saveToRedis(tagMapping, tagIDs)

	return tagMapping, tagIDs, nil
}

// 从Redis加载缓存
func loadFromRedis() (tagMapping map[int64][]int64, tagIDs []int64, err common.GFError) {
	// 读取tagMapping
	mappingStr, err := cs.GetString(redisTagMappingKey)
	if err != nil || mappingStr == "" {
		// 缓存不存在或读取失败
		return nil, nil, nil
	}
	// 反序列化map[int64][]int64
	tagMapping = make(map[int64][]int64)
	if err := json.Unmarshal([]byte(mappingStr), &tagMapping); err != nil {
		log.Error("tagMapping反序列化失败: " + err.Error())
		return nil, nil, common.NewServiceError("缓存数据格式错误")
	}

	// 读取tagIDs
	idsStr, err := cs.GetString(redisTagIDsKey)
	if err != nil || idsStr == "" {
		return nil, nil, nil
	}
	// 反序列化[]int64
	if err := json.Unmarshal([]byte(idsStr), &tagIDs); err != nil {
		log.Error("tagIDs反序列化失败: " + err.Error())
		return nil, nil, common.NewServiceError("缓存数据格式错误")
	}

	return tagMapping, tagIDs, nil
}

// 保存数据到Redis
func saveToRedis(tagMapping map[int64][]int64, tagIDs []int64) {
	// 序列化tagMapping并保存
	mappingBytes, err := json.Marshal(tagMapping)
	if err != nil {
		log.Error("tagMapping序列化失败: " + err.Error())
		return
	}
	if err := cs.SetExpire(redisTagMappingKey, mappingBytes, cacheExpireTime); err != nil {
		log.Error("tagMapping缓存写入失败: " + err.GetMsg())
	}

	// 序列化tagIDs并保存
	idsBytes, err := json.Marshal(tagIDs)
	if err != nil {
		log.Error("tagIDs序列化失败: " + err.Error())
		return
	}
	if err := cs.SetExpire(redisTagIDsKey, idsBytes, cacheExpireTime); err != nil {
		log.Error("tagIDs缓存写入失败: " + err.GetMsg())
	}
}
