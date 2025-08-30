package services

import (
	"fmt"
	"login/cache"
	"login/config"
	"login/utils/hkvilog"
	"math/rand"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// SMSService 短信服务
type SMSService struct {
	client *dysmsapi.Client
	config *config.SMSConfig
}

// NewSMSService 创建短信服务实例
func NewSMSService(cfg *config.SMSConfig) (*SMSService, error) {
	client, err := dysmsapi.NewClientWithAccessKey(
		cfg.RegionID,
		cfg.AccessKeyID,
		cfg.AccessKeySecret,
	)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云短信客户端失败: %v", err)
	}

	return &SMSService{
		client: client,
		config: cfg,
	}, nil
}

// GenerateSMSCode 生成6位数字验证码
func (s *SMSService) GenerateSMSCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000 // 生成100000-999999之间的数字
	return fmt.Sprintf("%06d", code)
}

// SendSMSCode 发送短信验证码
func (s *SMSService) SendSMSCode(phone string) (string, error) {
	// 生成验证码，暂时写死为201707
	code := "201707"
	//code := s.GenerateSMSCode()
	//
	// 创建短信请求
	//request := dysmsapi.CreateSendSmsRequest()
	//request.Scheme = "https"
	//request.PhoneNumbers = phone
	//request.SignName = s.config.SignName
	//request.TemplateCode = s.config.TemplateCode
	//request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)
	//
	//// 发送短信
	//response, err := s.client.SendSms(request)
	//if err != nil {
	//	return "", fmt.Errorf("发送短信失败: %v", err)
	//}
	//
	//// 检查发送结果
	//if response.Code != "OK" {
	//	return "", fmt.Errorf("短信发送失败: %s", response.Message)
	//}

	// 将验证码存储到Redis，5分钟过期
	err := cache.SetSMSCode(phone, code, 5*time.Minute)
	if err != nil {
		hkvilog.Errorf("存储验证码失败: %v", err)
		return "", fmt.Errorf("存储验证码失败")
	}

	hkvilog.Infof("短信验证码已发送到 %s", phone)
	return code, nil
}

// VerifySMSCode 验证短信验证码
func (s *SMSService) VerifySMSCode(phone, code string) (bool, error) {
	// 从Redis获取存储的验证码
	storedCode, err := cache.GetSMSCode(phone)
	if err != nil {
		return false, fmt.Errorf("验证码不存在或已过期")
	}

	// 比较验证码
	if storedCode != code {
		return false, fmt.Errorf("验证码错误")
	}

	// 验证成功后删除验证码
	err = cache.DeleteSMSCode(phone)
	if err != nil {
		hkvilog.Errorf("删除验证码失败: %v", err)
	}

	return true, nil
}

// CheckRateLimit 检查限流
func (s *SMSService) CheckRateLimit(phone, clientIP string) error {
	// 检查手机号限流（1分钟内只能发送1次）
	phoneKey := fmt.Sprintf("sms_rate_limit:phone:%s", phone)
	count, err := cache.IncrementRateLimit(phoneKey, time.Minute)
	if err != nil {
		return fmt.Errorf("检查手机号限流失败")
	}
	if count > 1 {
		return fmt.Errorf("该手机号发送过于频繁，请稍后再试")
	}

	// 检查IP限流（1分钟内最多发送10次）
	ipKey := fmt.Sprintf("sms_rate_limit:ip:%s", clientIP)
	count, err = cache.IncrementRateLimit(ipKey, time.Minute)
	if err != nil {
		return fmt.Errorf("检查IP限流失败")
	}
	if count > 10 {
		return fmt.Errorf("发送过于频繁，请稍后再试")
	}

	return nil
}

// ValidatePhoneNumber 验证手机号格式
func (s *SMSService) ValidatePhoneNumber(phone string) error {
	// 简单的手机号格式验证（中国大陆手机号）
	if len(phone) != 11 || !strings.HasPrefix(phone, "1") {
		return fmt.Errorf("手机号格式不正确")
	}
	return nil
}
