package controllers

import (
	"parking/models"
	"parking/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetFinanceList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	vehicleNo := c.Query("vehicle_no")
	financeType := c.Query("type")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	finances, total, err := models.GetFinanceList(page, pageSize, vehicleNo, financeType, startTime, endTime)
	if err != nil {
		c.JSON(200, utils.Error("获取财务列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(finances, total))
}

func GetFinanceDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	finance, err := models.GetFinanceByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("财务记录不存在"))
		return
	}

	c.JSON(200, utils.Success(finance))
}

func Payment(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		RecordID    uint   `json:"record_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		PaymentType int    `json:"payment_type"`
		Remark      string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	record, err := models.GetParkingRecordByID(req.RecordID)
	if err != nil {
		c.JSON(200, utils.Error("停车记录不存在"))
		return
	}

	if record.Status == 1 {
		c.JSON(200, utils.Error("该订单已支付"))
		return
	}

	now := time.Now()
	finance := &models.Finance{
		FinanceNo:   utils.GenerateFinanceNo(),
		RecordID:    req.RecordID,
		VehicleNo:   record.VehicleNo,
		Type:        1,
		Amount:      req.Amount,
		PaymentType: req.PaymentType,
		PayTime:     now,
		OperatorID:  userID.(uint),
		Remark:      req.Remark,
	}

	if err := models.CreateFinance(finance); err != nil {
		c.JSON(200, utils.Error("创建财务记录失败"))
		return
	}

	record.PaidAmount = req.Amount
	record.PaymentType = req.PaymentType
	record.PaymentTime = &now
	record.Status = 1

	if err := models.UpdateParkingRecord(record); err != nil {
		c.JSON(200, utils.Error("更新停车记录失败"))
		return
	}

	c.JSON(200, utils.Success(gin.H{
		"finance_no": finance.FinanceNo,
		"pay_time":   now,
		"amount":     req.Amount,
	}))
}

func GetDashboardStats(c *gin.Context) {
	todayCount, err := models.GetTodayCount()
	if err != nil {
		todayCount = 0
	}

	todayRevenue, err := models.GetTodayRevenue()
	if err != nil {
		todayRevenue = 0
	}

	activeCount, err := models.GetActiveCount()
	if err != nil {
		activeCount = 0
	}

	var totalSpaces, usedSpaces int64
	models.DB.Model(&models.ParkingSpace{}).Count(&totalSpaces)
	models.DB.Model(&models.ParkingSpace{}).Where("status = 1").Count(&usedSpaces)

	c.JSON(200, utils.Success(gin.H{
		"today_count":    todayCount,
		"today_revenue":  todayRevenue,
		"active_count":   activeCount,
		"total_spaces":   totalSpaces,
		"used_spaces":    usedSpaces,
		"free_spaces":    totalSpaces - usedSpaces,
	}))
}

func GetPartnerList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.Query("name")
	code := c.Query("code")
	status := c.Query("status")

	partners, total, err := models.GetPartnerList(page, pageSize, name, code, status)
	if err != nil {
		c.JSON(200, utils.Error("获取合作单位列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(partners, total))
}

func GetPartnerAll(c *gin.Context) {
	partners, err := models.GetAllPartners()
	if err != nil {
		c.JSON(200, utils.Error("获取合作单位列表失败"))
		return
	}

	c.JSON(200, utils.Success(partners))
}

func GetPartnerDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	partner, err := models.GetPartnerByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("合作单位不存在"))
		return
	}

	c.JSON(200, utils.Success(partner))
}

func CreatePartner(c *gin.Context) {
	var req struct {
		Name         string  `json:"name" binding:"required"`
		Code         string  `json:"code" binding:"required"`
		Type         int     `json:"type"`
		Contact      string  `json:"contact"`
		Phone        string  `json:"phone"`
		Email        string  `json:"email"`
		Address      string  `json:"address"`
		DiscountRate float64 `json:"discount_rate"`
		Remark       string  `json:"remark"`
		Status       int     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	partner := &models.Partner{
		Name:         req.Name,
		Code:         req.Code,
		Type:         req.Type,
		Contact:      req.Contact,
		Phone:        req.Phone,
		Email:        req.Email,
		Address:      req.Address,
		DiscountRate: req.DiscountRate,
		Remark:       req.Remark,
		Status:       req.Status,
	}

	if err := models.CreatePartner(partner); err != nil {
		c.JSON(200, utils.Error("创建合作单位失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("创建成功"))
}

func UpdatePartner(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	partner, err := models.GetPartnerByID(uint(id))
	if err != nil {
		c.JSON(200, utils.Error("合作单位不存在"))
		return
	}

	var req struct {
		Name         string  `json:"name"`
		Type         int     `json:"type"`
		Contact      string  `json:"contact"`
		Phone        string  `json:"phone"`
		Email        string  `json:"email"`
		Address      string  `json:"address"`
		DiscountRate float64 `json:"discount_rate"`
		Remark       string  `json:"remark"`
		Status       int     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Error("参数错误"))
		return
	}

	if req.Name != "" {
		partner.Name = req.Name
	}
	partner.Type = req.Type
	partner.Contact = req.Contact
	partner.Phone = req.Phone
	partner.Email = req.Email
	partner.Address = req.Address
	partner.DiscountRate = req.DiscountRate
	partner.Remark = req.Remark
	if req.Status != 0 {
		partner.Status = req.Status
	}

	if err := models.UpdatePartner(partner); err != nil {
		c.JSON(200, utils.Error("更新合作单位失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("更新成功"))
}

func DeletePartner(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := models.DeletePartner(uint(id)); err != nil {
		c.JSON(200, utils.Error("删除失败"))
		return
	}

	c.JSON(200, utils.SuccessMsg("删除成功"))
}

func GetLogList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	username := c.Query("username")
	module := c.Query("module")
	action := c.Query("action")

	logs, total, err := models.GetLogList(page, pageSize, username, module, action)
	if err != nil {
		c.JSON(200, utils.Error("获取日志列表失败"))
		return
	}

	c.JSON(200, utils.SuccessWithCount(logs, total))
}
