package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gateway/config"
	"gateway/utils/hkvilog"

	"github.com/gin-gonic/gin"
)

// ProxyToBusiness 代理请求到业务服务
func ProxyToBusiness(c *gin.Context) {
	cfg := config.LoadConfig()

	// 解析业务服务URL
	target, err := url.Parse(cfg.BusinessAPI.BaseURL)
	if err != nil {
		hkvilog.Errorf("解析业务服务URL失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "业务服务配置错误",
		})
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(target)

	// 修改请求
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// 修改请求路径，去掉/api/business前缀
		req.URL.Path = strings.Replace(req.URL.Path, "/api/business", "/api", 1)

		// 添加用户信息到请求头
		if userID, exists := c.Get("user_id"); exists {
			req.Header.Set("X-User-ID", fmt.Sprintf("%d", userID.(int)))
		}
		if username, exists := c.Get("username"); exists {
			req.Header.Set("X-Username", username.(string))
		}

		// 设置目标主机
		req.Host = target.Host

		hkvilog.Infof("代理请求到业务服务: %s %s", req.Method, req.URL.String())
	}

	// 修改响应
	proxy.ModifyResponse = func(resp *http.Response) error {
		// 可以在这里修改响应
		return nil
	}

	// 错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		hkvilog.Errorf("代理请求失败: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"error":"业务服务暂不可用"}`)
	}

	// 执行代理
	proxy.ServeHTTP(c.Writer, c.Request)
}
