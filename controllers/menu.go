package controllers

import (
	"parking/models"
	"parking/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMenuList(c *gin.Context) {
	menus, err := models.GetAllMenus()
	if err != nil {
		c.JSON(200, utils.Error("获取菜单列表失败"))
		return
	}

	menuTree := models.BuildMenuTree(menus, 0)
	c.JSON(200, utils.Success(menuTree))
}

func GetMenuAll(c *gin.Context) {
	menus, err := models.GetAllMenus()
	if err != nil {
		c.JSON(200, utils.Error("获取菜单列表失败"))
		return
	}

	c.JSON(200, utils.Success(menus))
}

func GetMenuDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	menu, err := models.GetMenuByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("菜单不存在"))
		return
	}

	c.JSON(200, utils.Success(menu))
}

func CreateMenu(c *gin.Context) {
	var req struct {
		ParentID uint   `json:"parent_id"`
		Name     string `json:"name" binding:"required"`
		Code     string `json:"code"`
		Icon     string `json:"icon"`
		URL      string `json:"url"`
		Type     int    `json:"type"`
		Sort     int    `json:"sort"`
		Status   int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	menu := &models.Menu{
		ParentID: req.ParentID,
		Name:     req.Name,
		Code:     req.Code,
		Icon:     req.Icon,
		URL:      req.URL,
		Type:     req.Type,
		Sort:     req.Sort,
		Status:   req.Status,
	}

	if err := models.CreateMenu(menu); err != nil {
		c.JSON(200, utils.Error("创建菜单失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("创建成功"))
}

func UpdateMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	menu, err := models.GetMenuByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("菜单不存在"))
		return
	}

	var req struct {
		ParentID uint   `json:"parent_id"`
		Name     string `json:"name"`
		Code     string `json:"code"`
		Icon     string `json:"icon"`
		URL      string `json:"url"`
		Type     int    `json:"type"`
		Sort     int    `json:"sort"`
		Status   int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	if req.Name != "" {
		menu.Name = req.Name
	}
	menu.ParentID = req.ParentID
	menu.Code = req.Code
	menu.Icon = req.Icon
	menu.URL = req.URL
	menu.Type = req.Type
	menu.Sort = req.Sort
	if req.Status != 0 {
		menu.Status = req.Status
	}

	if err := models.UpdateMenu(menu); err != nil {
		c.JSON(200, utils.Error("更新菜单失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("更新成功"))
}

func DeleteMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if models.HasChildren(uint(id)) {
		c.JSON(200, utils.Error("存在子菜单，不能删除"))
		return
	}

	if err := models.DeleteMenu(uint(id)); err != nil {
		c.JSON(200, utils.Error("删除失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("删除成功"))
}
