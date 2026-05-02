package controllers

import (
	"parking/models"
	"parking/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetVehicleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	vehicleNo := c.Query("vehicle_no")
	ownerName := c.Query("owner_name")
	status := c.Query("status")

	vehicles, total, err := models.GetVehicleList(page, pageSize, vehicleNo, ownerName, status)
	if err != nil {
		c.JSON(200, utils.Error("获取车辆列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(vehicles, total))
}

func GetVehicleDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	vehicle, err := models.GetVehicleByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("车辆不存在"))
		return
	}

	c.JSON(200, utils.Success(vehicle))
}

func CreateVehicle(c *gin.Context) {
	var req struct {
		VehicleNo   string `json:"vehicle_no" binding:"required"`
		OwnerName   string `json:"owner_name"`
		OwnerPhone  string `json:"owner_phone"`
		VehicleType int    `json:"vehicle_type"`
		PlateColor  int    `json:"plate_color"`
		Brand       string `json:"brand"`
		Model       string `json:"model"`
		Color       string `json:"color"`
		IsVIP       int    `json:"is_vip"`
		PartnerID   uint   `json:"partner_id"`
		Remark      string `json:"remark"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	exist, _ := models.GetVehicleByNo(req.VehicleNo)
	if exist != nil {
		c.JSON(200, utils.Error("车牌号已存在"))
		return
	}

	vehicle := &models.Vehicle{
		VehicleNo:   req.VehicleNo,
		OwnerName:   req.OwnerName,
		OwnerPhone:  req.OwnerPhone,
		VehicleType: req.VehicleType,
		PlateColor:  req.PlateColor,
		Brand:       req.Brand,
		Model:       req.Model,
		Color:       req.Color,
		IsVIP:       req.IsVIP,
		PartnerID:   req.PartnerID,
		Remark:      req.Remark,
		Status:      req.Status,
	}

	if err := models.CreateVehicle(vehicle); err != nil {
		c.JSON(200, utils.Error("创建车辆失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("创建成功"))
}

func UpdateVehicle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	vehicle, err := models.GetVehicleByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("车辆不存在"))
		return
	}

	var req struct {
		OwnerName   string `json:"owner_name"`
		OwnerPhone  string `json:"owner_phone"`
		VehicleType int    `json:"vehicle_type"`
		PlateColor  int    `json:"plate_color"`
		Brand       string `json:"brand"`
		Model       string `json:"model"`
		Color       string `json:"color"`
		IsVIP       int    `json:"is_vip"`
		PartnerID   uint   `json:"partner_id"`
		Remark      string `json:"remark"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	vehicle.OwnerName = req.OwnerName
	vehicle.OwnerPhone = req.OwnerPhone
	vehicle.VehicleType = req.VehicleType
	vehicle.PlateColor = req.PlateColor
	vehicle.Brand = req.Brand
	vehicle.Model = req.Model
	vehicle.Color = req.Color
	vehicle.IsVIP = req.IsVIP
	vehicle.PartnerID = req.PartnerID
	vehicle.Remark = req.Remark
	if req.Status != 0 {
		vehicle.Status = req.Status
	}

	if err := models.UpdateVehicle(vehicle); err != nil {
		c.JSON(200, utils.Error("更新车辆失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("更新成功"))
}

func DeleteVehicle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := models.DeleteVehicle(uint(id)); err != nil {
		c.JSON(200, utils.Error("删除失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("删除成功"))
}
