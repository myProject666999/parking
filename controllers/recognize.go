package controllers

import (
	"encoding/base64"
	"io/ioutil"
	"parking/models"
	"parking/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RecognizePlate(c *gin.Context) {
	var req struct {
		ImageBase64 string `json:"image_base64"`
		IsMock      bool   `json:"is_mock"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		file, _, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(200, utils.Error("请上传图片或提供图片base64"))
			return
		}
		defer file.Close()

		imgData, err := ioutil.ReadAll(file)
		if err != nil {
			c.JSON(200, utils.Error("读取图片失败"))
			return
		}

		req.ImageBase64 = base64.StdEncoding.EncodeToString(imgData)
	}

	var result *utils.PlateResult
	var err error

	if req.IsMock || req.ImageBase64 == "" {
		result, err = utils.RecognizePlateWithMock()
	} else {
		result, err = utils.RecognizePlate(req.ImageBase64)
	}

	if err != nil {
		c.JSON(200, utils.Error("识别失败: "+err.Error()))
		return
	}

	c.JSON(200, utils.Success(result))
}

func VehicleEntry(c *gin.Context) {
	var req struct {
		VehicleNo string `json:"vehicle_no" binding:"required"`
		LotID     uint   `json:"lot_id" binding:"required"`
		SpaceID   uint   `json:"space_id" binding:"required"`
		EntryImage string `json:"entry_image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	space, err := models.GetParkingSpaceByID(req.SpaceID)
	if err != nil {
		c.JSON(200, utils.Error("车位不存在"))
		return
	}

	if space.Status == 1 {
		c.JSON(200, utils.Error("该车位已被占用"))
		return
	}

	activeRecord, _ := models.GetActiveRecordByVehicle(req.VehicleNo)
	if activeRecord != nil {
		c.JSON(200, utils.Error("该车辆已在停车场内"))
		return
	}

	now := time.Now()
	record := &models.ParkingRecord{
		RecordNo:  utils.GenerateRecordNo(),
		LotID:     req.LotID,
		SpaceID:   req.SpaceID,
		VehicleNo: req.VehicleNo,
		EntryTime: now,
		EntryImage: req.EntryImage,
		Status:    0,
	}

	if err := models.CreateParkingRecord(record); err != nil {
		c.JSON(200, utils.Error("创建停车记录失败"))
		return
	}

	space.Status = 1
	space.VehicleNo = req.VehicleNo
	space.StartTime = &now
	if err := models.UpdateParkingSpace(space); err != nil {
		c.JSON(200, utils.Error("更新车位状态失败"))
		return
	}

	models.UpdateSpaceCount(req.LotID)

	c.JSON(200, utils.Success(gin.H{
		"record_no": record.RecordNo,
		"entry_time": now,
	}))
}

func VehicleExit(c *gin.Context) {
	var req struct {
		RecordNo  string `json:"record_no" binding:"required"`
		ExitImage string `json:"exit_image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	var record *models.ParkingRecord

	if req.RecordNo != "" {
		activeRecord, _ := models.GetActiveRecordByVehicle(req.RecordNo)
		if activeRecord != nil {
			record = activeRecord
		}
	}

	if record == nil {
		c.JSON(200, utils.Error("未找到停车记录"))
		return
	}

	if record.Status == 1 {
		c.JSON(200, utils.Error("该车辆已出场"))
		return
	}

	now := time.Now()
	exitTime := &now

	duration := exitTime.Sub(record.EntryTime)
	hours := duration.Hours()

	hourlyRate := 5.0
	amount := float64(int(hours+0.5)) * hourlyRate
	if amount < 5 {
		amount = 5
	}

	record.ExitTime = exitTime
	record.ParkingHours = hours
	record.Amount = amount
	record.ExitImage = req.ExitImage
	record.Status = 1

	if err := models.UpdateParkingRecord(record); err != nil {
		c.JSON(200, utils.Error("更新停车记录失败"))
		return
	}

	space, _ := models.GetParkingSpaceByID(record.SpaceID)
	if space != nil {
		space.Status = 0
		space.VehicleNo = ""
		space.StartTime = nil
		models.UpdateParkingSpace(space)
		models.UpdateSpaceCount(space.LotID)
	}

	c.JSON(200, utils.Success(gin.H{
		"record_no":     record.RecordNo,
		"entry_time":    record.EntryTime,
		"exit_time":     exitTime,
		"parking_hours": hours,
		"amount":        amount,
	}))
}

func GetParkingRecordList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	vehicleNo := c.Query("vehicle_no")
	status := c.Query("status")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	records, total, err := models.GetParkingRecordList(page, pageSize, vehicleNo, status, startTime, endTime)
	if err != nil {
		c.JSON(200, utils.Error("获取停车记录列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(records, total))
}

func GetParkingRecordDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	record, err := models.GetParkingRecordByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("停车记录不存在"))
		return
	}

	c.JSON(200, utils.Success(record))
}
