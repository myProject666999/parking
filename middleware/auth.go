package middleware

import (
	"net/http"
	"parking/models"
	"parking/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			authHeader = c.Query("token")
		}

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.Error("未登录或登录已过期"))
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Error("登录已过期，请重新登录"))
			c.Abort()
			return
		}

		user, err := models.GetUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Error("用户不存在"))
			c.Abort()
			return
		}

		if user.Status != 1 {
			c.JSON(http.StatusForbidden, utils.Error("账号已被禁用"))
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("role_id", user.RoleID)

		c.Next()
	}
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		userIDVal, _ := c.Get("user_id")
		username, _ := c.Get("username")
		
		var userID uint
		if userIDVal != nil {
			switch v := userIDVal.(type) {
			case uint:
				userID = v
			case int:
				userID = uint(v)
			case int64:
				userID = uint(v)
			case float64:
				userID = uint(v)
			default:
				userID = 0
			}
		} else {
			userID = 0
		}

		if username == nil {
			username = "guest"
		}

		log := &models.OperationLog{
			UserID:      userID,
			Username:    username.(string),
			Module:      getModule(path),
			Action:      getAction(method, path),
			Method:      method,
			URL:         path,
			IP:          utils.GetClientIP(c.Request),
			UserAgent:   c.Request.UserAgent(),
			Status:      1,
		}

		models.CreateLog(log)

		latency := time.Since(start)
		c.Set("latency", latency)
	}
}

func getModule(path string) string {
	if strings.Contains(path, "/api/login") || strings.Contains(path, "/api/logout") {
		return "登录"
	} else if strings.Contains(path, "/api/user") {
		return "用户管理"
	} else if strings.Contains(path, "/api/role") {
		return "角色管理"
	} else if strings.Contains(path, "/api/menu") {
		return "菜单管理"
	} else if strings.Contains(path, "/api/config") {
		return "系统配置"
	} else if strings.Contains(path, "/api/lot") || strings.Contains(path, "/api/space") {
		return "停车场管理"
	} else if strings.Contains(path, "/api/vehicle") {
		return "车辆管理"
	} else if strings.Contains(path, "/api/record") {
		return "停车记录"
	} else if strings.Contains(path, "/api/finance") {
		return "财务管理"
	} else if strings.Contains(path, "/api/log") {
		return "系统日志"
	} else if strings.Contains(path, "/api/partner") {
		return "合作单位"
	} else if strings.Contains(path, "/api/recognize") {
		return "车牌识别"
	} else if strings.Contains(path, "/api/dashboard") {
		return "控制台"
	}
	return "其他"
}

func getAction(method, path string) string {
	if method == "GET" {
		return "查询"
	} else if method == "POST" {
		if strings.Contains(path, "login") {
			return "登录"
		}
		return "新增"
	} else if method == "PUT" {
		return "修改"
	} else if method == "DELETE" {
		return "删除"
	}
	return method
}
