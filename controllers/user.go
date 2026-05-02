package controllers

import (
	"parking/models"
	"parking/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	username := c.Query("username")
	status := c.Query("status")

	users, total, err := models.GetUserList(page, pageSize, username, status)
	if err != nil {
		c.JSON(200, utils.Error("获取用户列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(users, total))
}

func GetUserDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := models.GetUserByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("用户不存在"))
		return
	}

	c.JSON(200, utils.Success(user))
}

func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		RoleID   uint   `json:"role_id" binding:"required"`
		Status   int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误: "+err.Error()))
		return
	}

	exist, _ := models.GetUserByUsername(req.Username)
	if exist != nil {
		c.JSON(200, utils.Error("用户名已存在"))
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: utils.MD5(req.Password),
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Email:    req.Email,
		RoleID:   req.RoleID,
		Status:   req.Status,
	}

	if err := models.CreateUser(user); err != nil {
		c.JSON(200, utils.Error("创建用户失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("创建成功"))
}

func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := models.GetUserByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("用户不存在"))
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		RoleID   uint   `json:"role_id"`
		Status   int    `json:"status"`
		Password string `json:"password"`
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
	if req.RoleID > 0 {
		user.RoleID = req.RoleID
	}
	if req.Status != 0 {
		user.Status = req.Status
	}
	if req.Password != "" {
		user.Password = utils.MD5(req.Password)
	}

	if err := models.UpdateUser(user); err != nil {
		c.JSON(200, utils.Error("更新用户失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("更新成功"))
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if uint(id) == 1 {
		c.JSON(200, utils.Error("不能删除超级管理员"))
		return
	}

	if err := models.DeleteUser(uint(id)); err != nil {
		c.JSON(200, utils.Error("删除失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("删除成功"))
}
