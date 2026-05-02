package models

import (
	"time"

	"gorm.io/gorm"
)

type ParkingRecord struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	RecordNo      string         `json:"record_no" gorm:"type:varchar(50);uniqueIndex;not null"`
	LotID         uint           `json:"lot_id" gorm:"not null"`
	SpaceID       uint           `json:"space_id" gorm:"not null"`
	VehicleNo     string         `json:"vehicle_no" gorm:"type:varchar(20);not null;index"`
	EntryTime     time.Time      `json:"entry_time"`
	ExitTime      *time.Time     `json:"exit_time"`
	ParkingHours  float64        `json:"parking_hours" gorm:"type:decimal(10,2);default:0"`
	Amount        float64        `json:"amount" gorm:"type:decimal(10,2);default:0"`
	Discount      float64        `json:"discount" gorm:"type:decimal(10,2);default:0"`
	PaidAmount    float64        `json:"paid_amount" gorm:"type:decimal(10,2);default:0"`
	PaymentType   int            `json:"payment_type" gorm:"type:tinyint;default:0"`
	PaymentTime   *time.Time     `json:"payment_time"`
	EntryImage    string         `json:"entry_image" gorm:"type:varchar(255)"`
	ExitImage     string         `json:"exit_image" gorm:"type:varchar(255)"`
	Status        int            `json:"status" gorm:"type:tinyint;default:0"`
	Remark        string         `json:"remark" gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ParkingRecord) TableName() string {
	return "parking_record"
}

func GetParkingRecordList(page, pageSize int, vehicleNo, status, startTime, endTime string) ([]ParkingRecord, int64, error) {
	var records []ParkingRecord
	var total int64

	query := DB.Model(&ParkingRecord{})
	if vehicleNo != "" {
		query = query.Where("vehicle_no LIKE ?", "%"+vehicleNo+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startTime != "" {
		query = query.Where("entry_time >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("entry_time <= ?", endTime+" 23:59:59")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error

	return records, total, err
}

func GetParkingRecordByID(id uint) (*ParkingRecord, error) {
	var record ParkingRecord
	err := DB.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func GetActiveRecordByVehicle(vehicleNo string) (*ParkingRecord, error) {
	var record ParkingRecord
	err := DB.Where("vehicle_no = ? AND status = 0", vehicleNo).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func CreateParkingRecord(record *ParkingRecord) error {
	return DB.Create(record).Error
}

func UpdateParkingRecord(record *ParkingRecord) error {
	return DB.Save(record).Error
}

func DeleteParkingRecord(id uint) error {
	return DB.Delete(&ParkingRecord{}, id).Error
}

func GetTodayCount() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := DB.Model(&ParkingRecord{}).
		Where("DATE(entry_time) = ?", today).
		Count(&count).Error
	return count, err
}

func GetTodayRevenue() (float64, error) {
	var total float64
	today := time.Now().Format("2006-01-02")
	err := DB.Model(&ParkingRecord{}).
		Where("DATE(entry_time) = ? AND status = 1", today).
		Select("COALESCE(SUM(paid_amount), 0)").
		Scan(&total).Error
	return total, err
}

func GetActiveCount() (int64, error) {
	var count int64
	err := DB.Model(&ParkingRecord{}).
		Where("status = 0").
		Count(&count).Error
	return count, err
}
