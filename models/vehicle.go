package models

import (
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	VehicleNo   string         `json:"vehicle_no" gorm:"type:varchar(20);uniqueIndex;not null"`
	OwnerName   string         `json:"owner_name" gorm:"type:varchar(50)"`
	OwnerPhone  string         `json:"owner_phone" gorm:"type:varchar(20)"`
	VehicleType int            `json:"vehicle_type" gorm:"type:tinyint;default:1"`
	PlateColor  int            `json:"plate_color" gorm:"type:tinyint;default:1"`
	Brand       string         `json:"brand" gorm:"type:varchar(50)"`
	Model       string         `json:"model" gorm:"type:varchar(50)"`
	Color       string         `json:"color" gorm:"type:varchar(20)"`
	IsVIP       int            `json:"is_vip" gorm:"type:tinyint;default:0"`
	PartnerID   uint           `json:"partner_id" gorm:"default:0"`
	Remark      string         `json:"remark" gorm:"type:varchar(255)"`
	Status      int            `json:"status" gorm:"type:tinyint;default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Vehicle) TableName() string {
	return "vehicle"
}

func GetVehicleList(page, pageSize int, vehicleNo, ownerName, status string) ([]Vehicle, int64, error) {
	var vehicles []Vehicle
	var total int64

	query := DB.Model(&Vehicle{})
	if vehicleNo != "" {
		query = query.Where("vehicle_no LIKE ?", "%"+vehicleNo+"%")
	}
	if ownerName != "" {
		query = query.Where("owner_name LIKE ?", "%"+ownerName+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&vehicles).Error

	return vehicles, total, err
}

func GetVehicleByID(id uint) (*Vehicle, error) {
	var vehicle Vehicle
	err := DB.First(&vehicle, id).Error
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func GetVehicleByNo(vehicleNo string) (*Vehicle, error) {
	var vehicle Vehicle
	err := DB.Where("vehicle_no = ?", vehicleNo).First(&vehicle).Error
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func CreateVehicle(vehicle *Vehicle) error {
	return DB.Create(vehicle).Error
}

func UpdateVehicle(vehicle *Vehicle) error {
	return DB.Save(vehicle).Error
}

func DeleteVehicle(id uint) error {
	return DB.Delete(&Vehicle{}, id).Error
}
