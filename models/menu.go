package models

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ParentID  uint           `json:"parent_id" gorm:"not null;default:0"`
	Name      string         `json:"name" gorm:"type:varchar(50);not null"`
	Code      string         `json:"code" gorm:"type:varchar(50)"`
	Icon      string         `json:"icon" gorm:"type:varchar(100)"`
	URL       string         `json:"url" gorm:"type:varchar(255)"`
	Type      int            `json:"type" gorm:"type:tinyint;default:1"`
	Sort      int            `json:"sort" gorm:"type:int;default:0"`
	Status    int            `json:"status" gorm:"type:tinyint;default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Children  []Menu         `json:"children,omitempty" gorm:"-"`
}

func (Menu) TableName() string {
	return "sys_menu"
}

func GetAllMenus() ([]Menu, error) {
	var menus []Menu
	err := DB.Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func GetMenusByRoleID(roleID uint) ([]Menu, error) {
	var menus []Menu
	err := DB.Joins("INNER JOIN sys_role_menu ON sys_menu.id = sys_role_menu.menu_id").
		Where("sys_role_menu.role_id = ?", roleID).
		Where("sys_menu.status = 1").
		Order("sys_menu.sort ASC, sys_menu.id ASC").
		Find(&menus).Error
	return menus, err
}

func BuildMenuTree(menus []Menu, parentID uint) []Menu {
	var tree []Menu
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = BuildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}

func GetMenuByID(id uint) (*Menu, error) {
	var menu Menu
	err := DB.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func CreateMenu(menu *Menu) error {
	return DB.Create(menu).Error
}

func UpdateMenu(menu *Menu) error {
	return DB.Save(menu).Error
}

func DeleteMenu(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id = ?", id).Delete(&RoleMenu{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Menu{}, id).Error
	})
}

func HasChildren(parentID uint) bool {
	var count int64
	DB.Model(&Menu{}).Where("parent_id = ?", parentID).Count(&count)
	return count > 0
}
