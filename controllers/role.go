package controllers

import (
	"parking/models"
	"parking/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRoleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.Query("name")
	status := c.Query("status")

	roles, total, err := models.GetRoleList(page, pageSize, name, status)
	if err != nil {
		c.JSON(200, utils.Error("获取角色列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(roles, total))
}

func GetRoleAll(c *gin.Context) {
	roles, err := models.GetAllRoles()
	if err != nil {
		c.JSON(200, utils.Error("获取角色列表失败"))
		return
	}

	c.JSON(200, utils.Success(roles))
}

func GetRoleDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	role, err := models.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("角色不存在"))
		return
	}

	menuIDs := make([]uint, 0)
	for _, menu := range role.Menus {
		menuIDs = append(menuIDs, menu.ID)
	}

	c.JSON(200, utils.Success(gin.H{
		"id":       role.ID,
		"name":     role.Name,
		"code":     role.Code,
		"remark":   role.Remark,
		"sort":     role.Sort,
		"status":   role.Status,
		"menu_ids": menuIDs,
	}))
}

func CreateRole(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Code    string `json:"code" binding:"required"`
		Remark  string `json:"remark"`
		Sort    int    `json:"sort"`
		Status  int    `json:"status"`
		MenuIDs []uint `json:"menu_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	role := &models.Role{
		Name:   req.Name,
		Code:   req.Code,
		Remark: req.Remark,
		Sort:   req.Sort,
		Status: req.Status,
	}

	if err := models.CreateRole(role, req.MenuIDs); err != nil {
		c.JSON(200, utils.Error("创建角色失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("创建成功"))
}

func UpdateRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	role, err := models.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("角色不存在"))
		return
	}

	var req struct {
		Name    string `json:"name"`
		Code    string `json:"code"`
		Remark  string `json:"remark"`
		Sort    int    `json:"sort"`
		Status  int    `json:"status"`
		MenuIDs []uint `json:"menu_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Code != "" {
		role.Code = req.Code
	}
	role.Remark = req.Remark
	role.Sort = req.Sort
	if req.Status != 0 {
		role.Status = req.Status
	}

	if err := models.UpdateRole(role, req.MenuIDs); err != nil {
		c.JSON(200, utils.Error("更新角色失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("更新成功"))
}

func DeleteRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if uint(id) == 1 {
		c.JSON(200, utils.Error("不能删除超级管理员角色"))
		return
	}

	if err := models.DeleteRole(uint(id)); err != nil {
		c.JSON(200, utils.Error("删除失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("删除成功"))
}
