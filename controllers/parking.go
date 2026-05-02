package controllers

import (
	"parking/models"
	"parking/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParkingLotList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.Query("name")
	code := c.Query("code")
	status := c.Query("status")

	lots, total, err := models.GetParkingLotList(page, pageSize, name, code, status)
	if err != nil {
		c.JSON(200, utils.Error("获取停车场列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(lots, total))
}

func GetParkingLotAll(c *gin.Context) {
	lots, err := models.GetAllParkingLots()
	if err != nil {
		c.JSON(200, utils.Error("获取停车场列表失败"))
		return
	}

	c.JSON(200, utils.Success(lots))
}

func GetParkingLotDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	lot, err := models.GetParkingLotByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("停车场不存在"))
		return
	}

	c.JSON(200, utils.Success(lot))
}

func CreateParkingLot(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Code        string  `json:"code" binding:"required"`
		Address     string  `json:"address"`
		TotalSpaces int     `json:"total_spaces"`
		Longitude   float64 `json:"longitude"`
		Latitude    float64 `json:"latitude"`
		Contact     string  `json:"contact"`
		Phone       string  `json:"phone"`
		Remark      string  `json:"remark"`
		Status      int     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	lot := &models.ParkingLot{
		Name:        req.Name,
		Code:        req.Code,
		Address:     req.Address,
		TotalSpaces: req.TotalSpaces,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		Contact:     req.Contact,
		Phone:       req.Phone,
		Remark:      req.Remark,
		Status:      req.Status,
	}

	if err := models.CreateParkingLot(lot); err != nil {
		c.JSON(200, utils.Error("创建停车场失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("创建成功"))
}

func UpdateParkingLot(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	lot, err := models.GetParkingLotByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("停车场不存在"))
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Code        string  `json:"code"`
		Address     string  `json:"address"`
		TotalSpaces int     `json:"total_spaces"`
		Longitude   float64 `json:"longitude"`
		Latitude    float64 `json:"latitude"`
		Contact     string  `json:"contact"`
		Phone       string  `json:"phone"`
		Remark      string  `json:"remark"`
		Status      int     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	if req.Name != "" {
		lot.Name = req.Name
	}
	if req.Code != "" {
		lot.Code = req.Code
	}
	lot.Address = req.Address
	lot.TotalSpaces = req.TotalSpaces
	lot.Longitude = req.Longitude
	lot.Latitude = req.Latitude
	lot.Contact = req.Contact
	lot.Phone = req.Phone
	lot.Remark = req.Remark
	if req.Status != 0 {
		lot.Status = req.Status
	}

	if err := models.UpdateParkingLot(lot); err != nil {
		c.JSON(200, utils.Error("更新停车场失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("更新成功"))
}

func DeleteParkingLot(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := models.DeleteParkingLot(uint(id)); err != nil {
		c.JSON(200, utils.Error("删除失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("删除成功"))
}

func GetParkingSpaceList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	lotID, _ := strconv.Atoi(c.Query("lot_id"))
	spaceNo := c.Query("space_no")
	area := c.Query("area")
	status := c.Query("status")

	spaces, total, err := models.GetParkingSpaceList(page, pageSize, uint(lotID), spaceNo, area, status)
	if err != nil {
		c.JSON(200, utils.Error("获取车位列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(spaces, total))
}

func GetParkingSpaceDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	space, err := models.GetParkingSpaceByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("车位不存在"))
		return
	}

	c.JSON(200, utils.Success(space))
}

func CreateParkingSpace(c *gin.Context) {
	var req struct {
		LotID   uint   `json:"lot_id" binding:"required"`
		SpaceNo string `json:"space_no" binding:"required"`
		Area    string `json:"area"`
		Type    int    `json:"type"`
		Status  int    `json:"status"`
		Remark  string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	space := &models.ParkingSpace{
		LotID:   req.LotID,
		SpaceNo: req.SpaceNo,
		Area:    req.Area,
		Type:    req.Type,
		Status:  req.Status,
		Remark:  req.Remark,
	}

	if err := models.CreateParkingSpace(space); err != nil {
		c.JSON(200, utils.Error("创建车位失败"))
		return
	}

	models.UpdateSpaceCount(req.LotID)

	c.JSON(200, utils.SuccessMsg("创建成功"))
}

func UpdateParkingSpace(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	space, err := models.GetParkingSpaceByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("车位不存在"))
		return
	}

	var req struct {
		SpaceNo string `json:"space_no"`
		Area    string `json:"area"`
		Type    int    `json:"type"`
		Status  int    `json:"status"`
		Remark  string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	if req.SpaceNo != "" {
		space.SpaceNo = req.SpaceNo
	}
	space.Area = req.Area
	space.Type = req.Type
	space.Status = req.Status
	space.Remark = req.Remark

	if err := models.UpdateParkingSpace(space); err != nil {
		c.JSON(200, utils.Error("更新车位失败"))
		return
	}

	models.UpdateSpaceCount(space.LotID)

	c.JSON(200, utils.SuccessMsg("更新成功"))
}

func DeleteParkingSpace(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	space, err := models.GetParkingSpaceByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("车位不存在"))
		return
	}

	if err := models.DeleteParkingSpace(uint(id)); err != nil {
		c.JSON(200, utils.Error("删除失败"))
		return
	}

	models.UpdateSpaceCount(space.LotID)

	c.JSON(200, utils.SuccessMsg("删除成功"))
}
