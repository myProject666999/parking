package models

import (
	"time"

	"gorm.io/gorm"
)

type Finance struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FinanceNo   string         `json:"finance_no" gorm:"type:varchar(50);uniqueIndex;not null"`
	RecordID    uint           `json:"record_id" gorm:"not null"`
	VehicleNo   string         `json:"vehicle_no" gorm:"type:varchar(20);not null;index"`
	Type        int            `json:"type" gorm:"type:tinyint;not null"`
	Amount      float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentType int            `json:"payment_type" gorm:"type:tinyint;default:0"`
	PayTime     time.Time      `json:"pay_time"`
	OperatorID  uint           `json:"operator_id" gorm:"not null"`
	Remark      string         `json:"remark" gorm:"type:varchar(255)"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Finance) TableName() string {
	return "finance"
}

type OperationLog struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Username    string         `json:"username" gorm:"type:varchar(50);not null"`
	Module      string         `json:"module" gorm:"type:varchar(50);not null"`
	Action      string         `json:"action" gorm:"type:varchar(50);not null"`
	Method      string         `json:"method" gorm:"type:varchar(10)"`
	URL         string         `json:"url" gorm:"type:varchar(255)"`
	IP          string         `json:"ip" gorm:"type:varchar(50)"`
	UserAgent   string         `json:"user_agent" gorm:"type:varchar(255)"`
	RequestData string         `json:"request_data" gorm:"type:text"`
	ResponseData string        `json:"response_data" gorm:"type:text"`
	Status      int            `json:"status" gorm:"type:tinyint;default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (OperationLog) TableName() string {
	return "operation_log"
}

type Partner struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null"`
	Code        string         `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Type        int            `json:"type" gorm:"type:tinyint;default:1"`
	Contact     string         `json:"contact" gorm:"type:varchar(50)"`
	Phone       string         `json:"phone" gorm:"type:varchar(20)"`
	Email       string         `json:"email" gorm:"type:varchar(100)"`
	Address     string         `json:"address" gorm:"type:varchar(255)"`
	DiscountRate float64       `json:"discount_rate" gorm:"type:decimal(5,2);default:0"`
	Remark      string         `json:"remark" gorm:"type:varchar(255)"`
	Status      int            `json:"status" gorm:"type:tinyint;default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Partner) TableName() string {
	return "partner"
}

func GetFinanceList(page, pageSize int, vehicleNo, financeType, startTime, endTime string) ([]Finance, int64, error) {
	var finances []Finance
	var total int64

	query := DB.Model(&Finance{})
	if vehicleNo != "" {
		query = query.Where("vehicle_no LIKE ?", "%"+vehicleNo+"%")
	}
	if financeType != "" {
		query = query.Where("type = ?", financeType)
	}
	if startTime != "" {
		query = query.Where("pay_time >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("pay_time <= ?", endTime+" 23:59:59")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&finances).Error

	return finances, total, err
}

func GetFinanceByID(id uint) (*Finance, error) {
	var finance Finance
	err := DB.First(&finance, id).Error
	if err != nil {
		return nil, err
	}
	return &finance, nil
}

func CreateFinance(finance *Finance) error {
	return DB.Create(finance).Error
}

func GetLogList(page, pageSize int, username, module, action string) ([]OperationLog, int64, error) {
	var logs []OperationLog
	var total int64

	query := DB.Model(&OperationLog{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error

	return logs, total, err
}

func CreateLog(log *OperationLog) error {
	return DB.Create(log).Error
}

func GetPartnerList(page, pageSize int, name, code, status string) ([]Partner, int64, error) {
	var partners []Partner
	var total int64

	query := DB.Model(&Partner{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&partners).Error

	return partners, total, err
}

func GetAllPartners() ([]Partner, error) {
	var partners []Partner
	err := DB.Where("status = 1").Order("id ASC").Find(&partners).Error
	return partners, err
}

func GetPartnerByID(id uint) (*Partner, error) {
	var partner Partner
	err := DB.First(&partner, id).Error
	if err != nil {
		return nil, err
	}
	return &partner, nil
}

func CreatePartner(partner *Partner) error {
	return DB.Create(partner).Error
}

func UpdatePartner(partner *Partner) error {
	return DB.Save(partner).Error
}

func DeletePartner(id uint) error {
	return DB.Delete(&Partner{}, id).Error
}

func GetTotalRevenue(startTime, endTime string) (float64, error) {
	var total float64
	query := DB.Model(&Finance{}).Where("type = 1")
	if startTime != "" {
		query = query.Where("pay_time >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("pay_time <= ?", endTime+" 23:59:59")
	}
	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}
