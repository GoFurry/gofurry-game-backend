package service

/*
 * @Desc: email服务
 * @author: 福狼
 * @version: v1.0.0
 */

import (
	"fmt"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/util"
	"github.com/GoFurry/gofurry-game-backend/roof/env"
	"gopkg.in/gomail.v2"
	"regexp"
)

// 发送邮箱验证码
func EmailSendCode(email string) (code string, gfsError common.GFError) {
	// 生成6位随机验证码
	code = util.GenerateRandomCode(common.EMAIL_CODE_LENGTH)
	m := gomail.NewMessage()
	m.SetHeader("From", env.GetServerConfig().Email.EmailUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "验证码")
	msg := fmt.Sprintf("您的注册验证码为: [%s], 3分钟内有效.", code)
	m.SetBody("text/html", msg)
	d := gomail.NewDialer(env.GetServerConfig().Email.EmailHost, env.GetServerConfig().Email.EmailPort, env.GetServerConfig().Email.EmailUser, env.GetServerConfig().Email.EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		gfsError = common.NewServiceError("邮件发送失败..." + err.Error())
	}
	return code, gfsError
}

// 校验邮箱是否合法
func IsEmailValid(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
