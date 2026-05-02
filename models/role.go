package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(50);not null"`
	Code      string         `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Remark    string         `json:"remark" gorm:"type:varchar(255)"`
	Sort      int            `json:"sort" gorm:"type:int;default:0"`
	Status    int            `json:"status" gorm:"type:tinyint;default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Menus     []Menu         `json:"menus,omitempty" gorm:"many2many:sys_role_menu;joinForeignKey:role_id;joinReferences:menu_id"`
}

func (Role) TableName() string {
	return "sys_role"
}

type RoleMenu struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	RoleID uint `json:"role_id" gorm:"not null"`
	MenuID uint `json:"menu_id" gorm:"not null"`
}

func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

func GetRoleList(page, pageSize int, name, status string) ([]Role, int64, error) {
	var roles []Role
	var total int64

	query := DB.Model(&Role{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&roles).Error

	return roles, total, err
}

func GetAllRoles() ([]Role, error) {
	var roles []Role
	err := DB.Where("status = 1").Order("sort ASC").Find(&roles).Error
	return roles, err
}

func GetRoleByID(id uint) (*Role, error) {
	var role Role
	err := DB.Preload("Menus").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func CreateRole(role *Role, menuIDs []uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		if len(menuIDs) > 0 {
			for _, mid := range menuIDs {
				rm := RoleMenu{RoleID: role.ID, MenuID: mid}
				if err := tx.Create(&rm).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func UpdateRole(role *Role, menuIDs []uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}

		if err := tx.Where("role_id = ?", role.ID).Delete(&RoleMenu{}).Error; err != nil {
			return err
		}

		if len(menuIDs) > 0 {
			for _, mid := range menuIDs {
				rm := RoleMenu{RoleID: role.ID, MenuID: mid}
				if err := tx.Create(&rm).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func DeleteRole(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&RoleMenu{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Role{}, id).Error
	})
}
