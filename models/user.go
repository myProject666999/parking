package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"`
	Nickname  string         `json:"nickname" gorm:"type:varchar(50)"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255)"`
	Phone     string         `json:"phone" gorm:"type:varchar(20)"`
	Email     string         `json:"email" gorm:"type:varchar(100)"`
	RoleID    uint           `json:"role_id" gorm:"not null"`
	Status    int            `json:"status" gorm:"type:tinyint;default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "sys_user"
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint) (*User, error) {
	var user User
	err := DB.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserList(page, pageSize int, username, status string) ([]User, int64, error) {
	var users []User
	var total int64

	query := DB.Model(&User{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error

	return users, total, err
}

func CreateUser(user *User) error {
	return DB.Create(user).Error
}

func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

func DeleteUser(id uint) error {
	return DB.Delete(&User{}, id).Error
}
