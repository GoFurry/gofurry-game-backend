package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/roof/env"
	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/experimental"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/gofiber/fiber/v2"
)

/*
 * @Desc: Coraza-WAF中间件
 * @author: 福狼
 * @version: v1.0.1
 */

// 全局 WAF 单例与初始化锁
var (
	globalWAF  coraza.WAF
	wafOnce    sync.Once
	wafInitErr error
)

// initGlobalWAF 初始化全局 WAF 单例
func initGlobalWAF() {
	wafOnce.Do(func() {
		globalWAF, wafInitErr = createWAF()
	})
}

// CorazaMiddleware 中间件
func CorazaMiddleware() fiber.Handler {
	initGlobalWAF()

	return func(context *fiber.Ctx) (err error) {

		// 直接使用全局 WAF 实例

		if wafInitErr != nil {
			fmt.Println("WAF全局实例初始化失败: ", wafInitErr)
			return common.NewResponse(context).ErrorWithCode("WAF 初始化失败", http.StatusInternalServerError)
		}
		if globalWAF == nil {
			return common.NewResponse(context).ErrorWithCode("WAF 实例未初始化", http.StatusInternalServerError)
		}

		// 事件句柄匿名函数
		newTX := func(*http.Request) types.Transaction {
			return globalWAF.NewTransaction()
		}
		// 事件句柄匿名函数
		if ctxwaf, ok := globalWAF.(experimental.WAFWithOptions); ok {
			newTX = func(r *http.Request) types.Transaction {
				return ctxwaf.NewTransactionWithOptions(experimental.Options{
					Context: r.Context(),
				})
			}
		}

		stdReq, err := convertFasthttpToStdRequest(context)
		if err != nil {
			return common.NewResponse(context).ErrorWithCode("请求转换失败", http.StatusInternalServerError)
		}

		// 开启事件
		tx := newTX(stdReq)
		defer func() {
			// 捕获 panic
			if r := recover(); r != nil {
				slog.Error(fmt.Sprintf("WAF transaction panicked: %v", r))
			}
			// 打印日志
			tx.ProcessLogging()
			// 关闭事件
			if err = tx.Close(); err != nil {
				tx.DebugLogger().Error().Err(err).Msg("Failed to close the transaction")
			}
		}()

		// 没开规则就返回
		if tx.IsRuleEngineOff() {
			return context.Next()
		}

		// 处理请求
		if it, err := processRequest(tx, stdReq); err != nil {
			tx.DebugLogger().Error().Err(err).Msg("Failed to process request")
			return common.NewResponse(context).ErrorWithCode("WAF处理请求失败", http.StatusInternalServerError)
		} else if it != nil {
			status := obtainStatusCodeFromInterruptionOrDefault(it, http.StatusOK)
			context.Status(status)
			context.Set("X-WAF-Blocked", "true")
			return common.NewResponse(context).ErrorWithCode("WAF拦截", status)
		}

		return context.Next()
	}
}

// 创建WAF
func createWAF() (coraza.WAF, error) {
	directivesFile := env.GetServerConfig().Waf.ConfPath
	if s := os.Getenv("DIRECTIVES_FILE"); s != "" {
		directivesFile = s
	}

	waf, err := coraza.NewWAF(
		coraza.NewWAFConfig().
			WithErrorCallback(logError).
			WithDirectivesFromFile(directivesFile),
	)

	return waf, err
}

// WAF 错误日志
func logError(error types.MatchedRule) {
	slog.Warn("WAF rule matched",
		slog.String("severity", string(error.Rule().Severity())),
		slog.String("error_log", error.ErrorLog()),
		slog.Int("rule_id", error.Rule().ID()),
	)
}

// 处理请求
func processRequest(tx types.Transaction, req *http.Request) (*types.Interruption, error) {
	var (
		client string
		cport  int
	)
	// IMPORTANT: Some http.Request.RemoteAddr implementations will not contain port or contain IPV6: [2001:db8::1]:8080
	idx := strings.LastIndexByte(req.RemoteAddr, ':')
	if idx != -1 {
		client = req.RemoteAddr[:idx]
		cport, _ = strconv.Atoi(req.RemoteAddr[idx+1:])
	}

	var in *types.Interruption
	// There is no socket access in the request object, so we neither know the server client nor port.
	tx.ProcessConnection(client, cport, "", 0)
	tx.ProcessURI(req.URL.String(), req.Method, req.Proto)
	for k, vr := range req.Header {
		for _, v := range vr {
			tx.AddRequestHeader(k, v)
		}
	}

	// Host will always be removed from req.Headers() and promoted to the
	// Request.Host field, so we manually add it
	if req.Host != "" {
		tx.AddRequestHeader("Host", req.Host)
		// This connector relies on the host header (now host field) to populate ServerName
		tx.SetServerName(req.Host)
	}

	// Transfer-Encoding header is removed by go/http
	// We manually add it to make rules relying on it work (E.g. CRS rule 920171)
	if req.TransferEncoding != nil {
		tx.AddRequestHeader("Transfer-Encoding", req.TransferEncoding[0])
	}

	in = tx.ProcessRequestHeaders()
	if in != nil {
		return in, nil
	}

	if tx.IsRequestBodyAccessible() {
		// We only do body buffering if the transaction requires request
		// body inspection, otherwise we just let the request follow its
		// regular flow.
		if req.Body != nil && req.Body != http.NoBody {
			it, _, err := tx.ReadRequestBodyFrom(req.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to append request body: %w", err)
			}

			if it != nil {
				return it, nil
			}

			rbr, err := tx.RequestBodyReader()
			if err != nil {
				return nil, fmt.Errorf("failed to get the request body: %w", err)
			}

			// Adds all remaining bytes beyond the coraza limit to its buffer
			// It happens when the partial body has been processed and it did not trigger an interruption
			bodyReader := io.MultiReader(rbr, req.Body)
			// req.Body is transparently reinizialied with a new io.ReadCloser.
			// The http handler will be able to read it.
			req.Body = io.NopCloser(bodyReader)
		}
	}

	return tx.ProcessRequestBody()
}

// "deny" Action 拒绝策略
func obtainStatusCodeFromInterruptionOrDefault(it *types.Interruption, defaultStatusCode int) int {
	if it.Action == "deny" {
		statusCode := it.Status
		if statusCode == 0 {
			statusCode = 403
		}

		return statusCode
	}
	return defaultStatusCode
}

func convertFasthttpToStdRequest(c *fiber.Ctx) (*http.Request, error) {
	// 复制请求体
	bodyCopy := make([]byte, len(c.Body()))
	copy(bodyCopy, c.Body())

	uri := c.OriginalURL()
	fullURL := fmt.Sprintf("%s://%s%s", c.Protocol(), c.Hostname(), uri)

	// 创建标准库请求
	stdReq, err := http.NewRequest(
		c.Method(),
		fullURL,
		bytes.NewReader(bodyCopy),
	)
	if err != nil {
		return nil, err
	}

	// 复制请求头
	c.Request().Header.VisitAll(func(key, value []byte) {
		stdReq.Header.Add(string(key), string(value))
	})

	// 获取真实客户端远程地址
	stdReq.RemoteAddr = net.JoinHostPort(c.IP(), c.Port())

	// 手动设置Host
	stdReq.Host = c.Hostname()

	return stdReq, nil
}
