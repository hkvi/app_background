package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gateway/cache"
	"gateway/config"
	"gateway/utils"
	"gateway/utils/hkvilog"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	cfg *config.Config
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone,omitempty"`
}

// SMSRequest 短信请求结构
type SMSRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// SMSLoginRequest 短信登录请求结构
type SMSLoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// RefreshTokenRequest 刷新令牌请求结构
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// BusinessResponse 业务服务响应结构
type BusinessResponse struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	User    interface{} `json:"user,omitempty"`
}

// Login 用户登录处理器
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 转发请求到业务服务
	resp, err := h.forwardToBusinessService("POST", "/api/auth/login", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "登录服务暂不可用",
		})
		return
	}

	// 如果业务服务返回成功，生成双token
	if resp.StatusCode == http.StatusOK {
		var businessResp BusinessResponse
		if err := json.NewDecoder(resp.Body).Decode(&businessResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "响应解析失败",
			})
			return
		}

		// 从业务服务响应中提取用户信息
		userData, ok := businessResp.Data.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "用户信息格式错误",
			})
			return
		}

		user, ok := userData["user"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "用户信息不存在",
			})
			return
		}

		userID := int(user["id"].(float64))
		username := user["username"].(string)

		// 生成双token
		accessToken, refreshToken, err := utils.GenerateTokenPair(
			userID,
			username,
			h.cfg.JWT.AccessSecretKey,
			h.cfg.JWT.RefreshSecretKey,
			h.cfg.JWT.AccessExpire,
			h.cfg.JWT.RefreshExpire,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "令牌生成失败",
			})
			return
		}

		// 存储刷新令牌到Redis
		if err := cache.StoreRefreshToken(userID, refreshToken, time.Duration(h.cfg.JWT.RefreshExpire)*time.Second); err != nil {
			hkvilog.Errorf("存储刷新令牌失败: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
			"data": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"expires_in":    h.cfg.JWT.AccessExpire,
				"user":          user,
			},
		})
		return
	}

	// 转发业务服务的错误响应
	h.forwardResponse(c, resp)
}

// Register 用户注册处理器
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 转发请求到业务服务
	resp, err := h.forwardToBusinessService("POST", "/api/auth/register", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注册服务暂不可用",
		})
		return
	}

	// 转发业务服务的响应
	h.forwardResponse(c, resp)
}

// SendSMS 发送短信验证码处理器
func (h *AuthHandler) SendSMS(c *gin.Context) {
	var req SMSRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 转发请求到业务服务
	resp, err := h.forwardToBusinessService("POST", "/api/sms/send", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "短信服务暂不可用",
		})
		return
	}

	// 转发业务服务的响应
	h.forwardResponse(c, resp)
}

// SMSLogin 短信验证码登录处理器
func (h *AuthHandler) SMSLogin(c *gin.Context) {
	var req SMSLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 转发请求到业务服务
	resp, err := h.forwardToBusinessService("POST", "/api/sms/login", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "短信登录服务暂不可用",
		})
		return
	}

	// 如果业务服务返回成功，生成双token
	if resp.StatusCode == http.StatusOK {
		var businessResp BusinessResponse
		if err := json.NewDecoder(resp.Body).Decode(&businessResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "响应解析失败",
			})
			return
		}

		// 从业务服务响应中提取用户信息
		userData, ok := businessResp.Data.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "用户信息格式错误",
			})
			return
		}

		user, ok := userData["user"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "用户信息不存在",
			})
			return
		}

		userID := int(user["id"].(float64))
		username := ""
		if user["username"] != nil {
			username = user["username"].(string)
		}

		// 生成双token
		accessToken, refreshToken, err := utils.GenerateTokenPair(
			userID,
			username,
			h.cfg.JWT.AccessSecretKey,
			h.cfg.JWT.RefreshSecretKey,
			h.cfg.JWT.AccessExpire,
			h.cfg.JWT.RefreshExpire,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "令牌生成失败",
			})
			return
		}

		// 存储刷新令牌到Redis
		if err := cache.StoreRefreshToken(userID, refreshToken, time.Duration(h.cfg.JWT.RefreshExpire)*time.Second); err != nil {
			hkvilog.Errorf("存储刷新令牌失败: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
			"data": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"expires_in":    h.cfg.JWT.AccessExpire,
				"user":          user,
			},
		})
		return
	}

	// 转发业务服务的错误响应
	h.forwardResponse(c, resp)
}

// RefreshToken 刷新令牌处理器
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证刷新令牌
	claims, err := utils.ValidateRefreshToken(req.RefreshToken, h.cfg.JWT.RefreshSecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无效的刷新令牌",
		})
		return
	}

	// 检查Redis中的刷新令牌
	storedToken, err := cache.GetRefreshToken(claims.UserID)
	if err != nil || storedToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "刷新令牌已失效",
		})
		return
	}

	// 生成新的双token
	newAccessToken, newRefreshToken, err := utils.GenerateTokenPair(
		claims.UserID,
		claims.Username,
		h.cfg.JWT.AccessSecretKey,
		h.cfg.JWT.RefreshSecretKey,
		h.cfg.JWT.AccessExpire,
		h.cfg.JWT.RefreshExpire,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "令牌生成失败",
		})
		return
	}

	// 更新Redis中的刷新令牌
	if err := cache.StoreRefreshToken(claims.UserID, newRefreshToken, time.Duration(h.cfg.JWT.RefreshExpire)*time.Second); err != nil {
		hkvilog.Errorf("更新刷新令牌失败: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "令牌刷新成功",
		"data": gin.H{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
			"expires_in":    h.cfg.JWT.AccessExpire,
		},
	})
}

// Logout 用户退出处理器
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户信息不存在",
		})
		return
	}

	// tokenID, exists := c.Get("token_id")
	// if !exists {
	//	c.JSON(http.StatusUnauthorized, gin.H{
	//		"error": "令牌信息不存在",
	//	})
	//	return
	// }

	// 删除Redis中的刷新令牌
	if err := cache.DeleteRefreshToken(userID.(int)); err != nil {
		hkvilog.Errorf("删除刷新令牌失败: %v", err)
	}

	// 将访问令牌加入黑名单 (暂时注释)
	// if err := cache.BlacklistToken(tokenID.(string), time.Duration(h.cfg.JWT.AccessExpire)*time.Second); err != nil {
	//	hkvilog.Errorf("添加令牌黑名单失败: %v", err)
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "退出成功",
	})
}

// forwardToBusinessService 转发请求到业务服务
func (h *AuthHandler) forwardToBusinessService(method, path string, data interface{}) (*http.Response, error) {
	// 序列化请求数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 构建请求URL
	url := h.cfg.BusinessAPI.BaseURL + path

	// 创建HTTP请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{
		Timeout: time.Duration(h.cfg.BusinessAPI.Timeout) * time.Second,
	}

	return client.Do(req)
}

// forwardResponse 转发响应
func (h *AuthHandler) forwardResponse(c *gin.Context, resp *http.Response) {
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "读取响应失败",
		})
		return
	}

	// 设置状态码并返回响应体
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
