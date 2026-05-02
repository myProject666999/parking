package models

import (
	"time"

	"gorm.io/gorm"
)

type ParkingLot struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null"`
	Code        string         `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Address     string         `json:"address" gorm:"type:varchar(255)"`
	TotalSpaces int            `json:"total_spaces" gorm:"type:int;not null;default:0"`
	UsedSpaces  int            `json:"used_spaces" gorm:"type:int;default:0"`
	FreeSpaces  int            `json:"free_spaces" gorm:"type:int;default:0"`
	Longitude   float64        `json:"longitude" gorm:"type:decimal(10,7)"`
	Latitude    float64        `json:"latitude" gorm:"type:decimal(10,7)"`
	Contact     string         `json:"contact" gorm:"type:varchar(50)"`
	Phone       string         `json:"phone" gorm:"type:varchar(20)"`
	Remark      string         `json:"remark" gorm:"type:varchar(255)"`
	Status      int            `json:"status" gorm:"type:tinyint;default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ParkingLot) TableName() string {
	return "parking_lot"
}

type ParkingSpace struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	LotID        uint           `json:"lot_id" gorm:"not null;index"`
	SpaceNo      string         `json:"space_no" gorm:"type:varchar(50);not null"`
	Area         string         `json:"area" gorm:"type:varchar(50)"`
	Type         int            `json:"type" gorm:"type:tinyint;default:1"`
	Status       int            `json:"status" gorm:"type:tinyint;default:0"`
	VehicleNo    string         `json:"vehicle_no" gorm:"type:varchar(20)"`
	StartTime    *time.Time     `json:"start_time"`
	Remark       string         `json:"remark" gorm:"type:varchar(255)"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ParkingSpace) TableName() string {
	return "parking_space"
}

func GetParkingLotList(page, pageSize int, name, code, status string) ([]ParkingLot, int64, error) {
	var lots []ParkingLot
	var total int64

	query := DB.Model(&ParkingLot{})
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
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&lots).Error

	return lots, total, err
}

func GetAllParkingLots() ([]ParkingLot, error) {
	var lots []ParkingLot
	err := DB.Where("status = 1").Order("id ASC").Find(&lots).Error
	return lots, err
}

func GetParkingLotByID(id uint) (*ParkingLot, error) {
	var lot ParkingLot
	err := DB.First(&lot, id).Error
	if err != nil {
		return nil, err
	}
	return &lot, nil
}

func CreateParkingLot(lot *ParkingLot) error {
	return DB.Create(lot).Error
}

func UpdateParkingLot(lot *ParkingLot) error {
	return DB.Save(lot).Error
}

func DeleteParkingLot(id uint) error {
	return DB.Delete(&ParkingLot{}, id).Error
}

func GetParkingSpaceList(page, pageSize int, lotID uint, spaceNo, area, status string) ([]ParkingSpace, int64, error) {
	var spaces []ParkingSpace
	var total int64

	query := DB.Model(&ParkingSpace{})
	if lotID > 0 {
		query = query.Where("lot_id = ?", lotID)
	}
	if spaceNo != "" {
		query = query.Where("space_no LIKE ?", "%"+spaceNo+"%")
	}
	if area != "" {
		query = query.Where("area LIKE ?", "%"+area+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&spaces).Error

	return spaces, total, err
}

func GetParkingSpaceByID(id uint) (*ParkingSpace, error) {
	var space ParkingSpace
	err := DB.First(&space, id).Error
	if err != nil {
		return nil, err
	}
	return &space, nil
}

func CreateParkingSpace(space *ParkingSpace) error {
	return DB.Create(space).Error
}

func UpdateParkingSpace(space *ParkingSpace) error {
	return DB.Save(space).Error
}

func DeleteParkingSpace(id uint) error {
	return DB.Delete(&ParkingSpace{}, id).Error
}

func UpdateSpaceCount(lotID uint) error {
	var total, used int64
	DB.Model(&ParkingSpace{}).Where("lot_id = ?", lotID).Count(&total)
	DB.Model(&ParkingSpace{}).Where("lot_id = ? AND status = 1", lotID).Count(&used)

	return DB.Model(&ParkingLot{}).Where("id = ?", lotID).
		Updates(map[string]interface{}{
			"total_spaces": total,
			"used_spaces":  used,
			"free_spaces":  total - used,
		}).Error
}
