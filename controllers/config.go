package controllers

import (
	"parking/models"
	"parking/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetConfigList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.Query("name")
	code := c.Query("code")

	configs, total, err := models.GetConfigList(page, pageSize, name, code)
	if err != nil {
		c.JSON(200, utils.Error("获取配置列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(configs, total))
}

func GetConfigDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	config, err := models.GetConfigByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("配置不存在"))
		return
	}

	c.JSON(200, utils.Success(config))
}

func GetConfigByCode(c *gin.Context) {
	code := c.Param("code")

	config, err := models.GetConfigByCode(code)
	if err != nil {
		c.JSON(200, utils.Error("配置不存在"))
		return
	}

	c.JSON(200, utils.Success(config))
}

func CreateConfig(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Code   string `json:"code" binding:"required"`
		Value  string `json:"value"`
		Type   int    `json:"type"`
		Remark string `json:"remark"`
		Sort   int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	exist, _ := models.GetConfigByCode(req.Code)
	if exist != nil {
		c.JSON(200, utils.Error("配置编码已存在"))
		return
	}

	config := &models.Config{
		Name:   req.Name,
		Code:   req.Code,
		Value:  req.Value,
		Type:   req.Type,
		Remark: req.Remark,
		Sort:   req.Sort,
	}

	if err := models.CreateConfig(config); err != nil {
		c.JSON(200, utils.Error("创建配置失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("创建成功"))
}

func UpdateConfig(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	config, err := models.GetConfigByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("配置不存在"))
		return
	}

	var req struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Type   int    `json:"type"`
		Remark string `json:"remark"`
		Sort   int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	if req.Name != "" {
		config.Name = req.Name
	}
	config.Value = req.Value
	config.Type = req.Type
	config.Remark = req.Remark
	config.Sort = req.Sort

	if err := models.UpdateConfig(config); err != nil {
		c.JSON(200, utils.Error("更新配置失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("更新成功"))
}

func DeleteConfig(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := models.DeleteConfig(uint(id)); err != nil {
		c.JSON(200, utils.Error("删除失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("删除成功"))
}
