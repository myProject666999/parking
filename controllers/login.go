package controllers

import (
	"parking/models"
	"parking/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("用户名或密码不能为空"))
		return
	}

	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(200, utils.Error("用户名或密码错误"))
		return
	}

	if !utils.VerifyPassword(req.Password, user.Password) {
		c.JSON(200, utils.Error("用户名或密码错误"))
		return
	}

	if user.Status != 1 {
		c.JSON(200, utils.Error("账号已被禁用，请联系管理员"))
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleID)
	if err != nil {
		c.JSON(200, utils.Error("生成Token失败"))
		return
	}

	c.JSON(200, utils.Success(gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
	}))
}

func Logout(c *gin.Context) {
	c.JSON(200, utils.SuccessMsg("退出登录成功"))
}

func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := models.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(200, utils.Error("获取用户信息失败"))
		return
	}

	menus, err := models.GetMenusByRoleID(user.RoleID)
	if err != nil {
		c.JSON(200, utils.Error("获取菜单失败"))
		return
	}

	menuTree := models.BuildMenuTree(menus, 0)

	c.JSON(200, utils.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"phone":    user.Phone,
		"email":    user.Email,
		"role_id":  user.RoleID,
		"menus":    menuTree,
	}))
}

func ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	user, err := models.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(200, utils.Error("用户不存在"))
		return
	}

	if !utils.VerifyPassword(req.OldPassword, user.Password) {
		c.JSON(200, utils.Error("原密码错误"))
		return
	}

	user.Password = utils.MD5(req.NewPassword)
	if err := models.UpdateUser(user); err != nil {
		c.JSON(200, utils.Error("修改密码失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("密码修改成功"))
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := models.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(200, utils.Error("用户不存在"))
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := models.UpdateUser(user); err != nil {
		c.JSON(200, utils.Error("更新失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("更新成功"))
}
