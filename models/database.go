package models

import (
	"fmt"
	"parking/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.AppConfig.Database.Username,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Database,
		config.AppConfig.Database.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = DB.AutoMigrate(
		&User{},
		&Role{},
		&Menu{},
		&RoleMenu{},
		&Config{},
		&ParkingLot{},
		&ParkingSpace{},
		&Vehicle{},
		&ParkingRecord{},
		&Finance{},
		&OperationLog{},
		&Partner{},
	)

	if err != nil {
		return err
	}

	initData()

	return nil
}

func initData() {
	var count int64
	DB.Model(&Role{}).Count(&count)
	if count == 0 {
		adminRole := Role{
			Name:     "超级管理员",
			Code:     "admin",
			Sort:     1,
			Status:   1,
		}
		DB.Create(&adminRole)

		adminUser := User{
			Username: "admin",
			Password: "e10adc3949ba59abbe56e057f20f883e",
			Nickname: "管理员",
			RoleID:   adminRole.ID,
			Status:   1,
		}
		DB.Create(&adminUser)

		menus := []Menu{
			{ParentID: 0, Name: "首页", Code: "dashboard", Icon: "layui-icon-home", Type: 1, Sort: 1, Status: 1},
			{ParentID: 0, Name: "系统管理", Code: "system", Icon: "layui-icon-set", Type: 1, Sort: 2, Status: 1},
			{ParentID: 2, Name: "用户管理", Code: "user", Type: 2, Sort: 1, Status: 1},
			{ParentID: 2, Name: "角色管理", Code: "role", Type: 2, Sort: 2, Status: 1},
			{ParentID: 2, Name: "菜单管理", Code: "menu", Type: 2, Sort: 3, Status: 1},
			{ParentID: 2, Name: "全局配置", Code: "config", Type: 2, Sort: 4, Status: 1},
			{ParentID: 2, Name: "系统日志", Code: "log", Type: 2, Sort: 5, Status: 1},
			{ParentID: 0, Name: "停车场管理", Code: "parking", Icon: "layui-icon-car", Type: 1, Sort: 3, Status: 1},
			{ParentID: 8, Name: "停车场列表", Code: "parking_lot", Type: 2, Sort: 1, Status: 1},
			{ParentID: 8, Name: "车位管理", Code: "parking_space", Type: 2, Sort: 2, Status: 1},
			{ParentID: 0, Name: "车牌识别", Code: "recognize", Icon: "layui-icon-face-camera", Type: 1, Sort: 4, Status: 1},
			{ParentID: 0, Name: "车辆管理", Code: "vehicle", Icon: "layui-icon-diamond", Type: 1, Sort: 5, Status: 1},
			{ParentID: 0, Name: "停车记录", Code: "record", Icon: "layui-icon-list", Type: 1, Sort: 6, Status: 1},
			{ParentID: 0, Name: "财务管理", Code: "finance", Icon: "layui-icon-rmb", Type: 1, Sort: 7, Status: 1},
			{ParentID: 0, Name: "合作单位", Code: "partner", Icon: "layui-icon-group", Type: 1, Sort: 8, Status: 1},
		}
		DB.Create(&menus)

		var menuIDs []uint
		for _, m := range menus {
			menuIDs = append(menuIDs, m.ID)
		}
		for _, mid := range menuIDs {
			DB.Create(&RoleMenu{RoleID: adminRole.ID, MenuID: mid})
		}
	}
}
