package models

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	Code      string         `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Value     string         `json:"value" gorm:"type:text"`
	Type      int            `json:"type" gorm:"type:tinyint;default:1"`
	Remark    string         `json:"remark" gorm:"type:varchar(255)"`
	Sort      int            `json:"sort" gorm:"type:int;default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Config) TableName() string {
	return "sys_config"
}

func GetConfigList(page, pageSize int, name, code string) ([]Config, int64, error) {
	var configs []Config
	var total int64

	query := DB.Model(&Config{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&configs).Error

	return configs, total, err
}

func GetConfigByCode(code string) (*Config, error) {
	var config Config
	err := DB.Where("code = ?", code).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetConfigByID(id uint) (*Config, error) {
	var config Config
	err := DB.First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func CreateConfig(config *Config) error {
	return DB.Create(config).Error
}

func UpdateConfig(config *Config) error {
	return DB.Save(config).Error
}

func DeleteConfig(id uint) error {
	return DB.Delete(&Config{}, id).Error
}

func GetConfigValue(code string, defaultValue string) string {
	config, err := GetConfigByCode(code)
	if err != nil {
		return defaultValue
	}
	return config.Value
}
