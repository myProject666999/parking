package router

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"parking/controllers"
	"parking/middleware"
	"strings"

	"github.com/gin-gonic/gin"
)

func loadAllTemplates(baseDir string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Loading templates from: %s\n", absBaseDir)

	var root *template.Template

	htmlFiles := []string{}
	filepath.Walk(absBaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".html") {
			htmlFiles = append(htmlFiles, path)
		}
		return nil
	})

	fmt.Printf("Found %d HTML files\n", len(htmlFiles))

	for _, path := range htmlFiles {
		relPath := strings.TrimPrefix(strings.ReplaceAll(path, "\\", "/"), strings.ReplaceAll(absBaseDir, "\\", "/")+"/")
		fmt.Printf("Parsing template: %s\n", relPath)

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var tmpl *template.Template
		if root == nil {
			root = template.New(relPath).Funcs(funcMap).Delims("<<", ">>")
			tmpl = root
		} else {
			tmpl = root.New(relPath).Funcs(funcMap).Delims("<<", ">>")
		}

		_, err = tmpl.Parse(string(content))
		if err != nil {
			return nil, err
		}
	}

	if root != nil {
		fmt.Printf("Total templates loaded: %d\n", len(root.Templates()))
		for _, t := range root.Templates() {
			fmt.Printf("  - %s\n", t.Name())
		}
	}

	return root, nil
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./static")

	templates, err := loadAllTemplates("templates")
	if err != nil {
		fmt.Printf("Error loading templates: %v\n", err)
	}

	if templates != nil {
		r.SetHTMLTemplate(templates)
	}

	r.Use(middleware.LogMiddleware())

	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/logout", controllers.Logout)

		auth := api.Group("")
		auth.Use(middleware.Auth())
		{
			auth.GET("/user/info", controllers.GetUserInfo)
			auth.POST("/user/password", controllers.ChangePassword)
			auth.PUT("/user/profile", controllers.UpdateProfile)

			auth.GET("/dashboard/stats", controllers.GetDashboardStats)

			auth.GET("/user/list", controllers.GetUserList)
			auth.GET("/user/:id", controllers.GetUserDetail)
			auth.POST("/user", controllers.CreateUser)
			auth.PUT("/user/:id", controllers.UpdateUser)
			auth.DELETE("/user/:id", controllers.DeleteUser)

			auth.GET("/role/list", controllers.GetRoleList)
			auth.GET("/role/all", controllers.GetRoleAll)
			auth.GET("/role/:id", controllers.GetRoleDetail)
			auth.POST("/role", controllers.CreateRole)
			auth.PUT("/role/:id", controllers.UpdateRole)
			auth.DELETE("/role/:id", controllers.DeleteRole)

			auth.GET("/menu/list", controllers.GetMenuList)
			auth.GET("/menu/all", controllers.GetMenuAll)
			auth.GET("/menu/:id", controllers.GetMenuDetail)
			auth.POST("/menu", controllers.CreateMenu)
			auth.PUT("/menu/:id", controllers.UpdateMenu)
			auth.DELETE("/menu/:id", controllers.DeleteMenu)

			auth.GET("/config/list", controllers.GetConfigList)
			auth.GET("/config/code/:code", controllers.GetConfigByCode)
			auth.GET("/config/:id", controllers.GetConfigDetail)
			auth.POST("/config", controllers.CreateConfig)
			auth.PUT("/config/:id", controllers.UpdateConfig)
			auth.DELETE("/config/:id", controllers.DeleteConfig)

			auth.GET("/lot/list", controllers.GetParkingLotList)
			auth.GET("/lot/all", controllers.GetParkingLotAll)
			auth.GET("/lot/:id", controllers.GetParkingLotDetail)
			auth.POST("/lot", controllers.CreateParkingLot)
			auth.PUT("/lot/:id", controllers.UpdateParkingLot)
			auth.DELETE("/lot/:id", controllers.DeleteParkingLot)

			auth.GET("/space/list", controllers.GetParkingSpaceList)
			auth.GET("/space/:id", controllers.GetParkingSpaceDetail)
			auth.POST("/space", controllers.CreateParkingSpace)
			auth.PUT("/space/:id", controllers.UpdateParkingSpace)
			auth.DELETE("/space/:id", controllers.DeleteParkingSpace)

			auth.GET("/vehicle/list", controllers.GetVehicleList)
			auth.GET("/vehicle/:id", controllers.GetVehicleDetail)
			auth.POST("/vehicle", controllers.CreateVehicle)
			auth.PUT("/vehicle/:id", controllers.UpdateVehicle)
			auth.DELETE("/vehicle/:id", controllers.DeleteVehicle)

			auth.POST("/recognize", controllers.RecognizePlate)
			auth.POST("/vehicle/entry", controllers.VehicleEntry)
			auth.POST("/vehicle/exit", controllers.VehicleExit)

			auth.GET("/record/list", controllers.GetParkingRecordList)
			auth.GET("/record/:id", controllers.GetParkingRecordDetail)

			auth.GET("/finance/list", controllers.GetFinanceList)
			auth.GET("/finance/:id", controllers.GetFinanceDetail)
			auth.POST("/finance/payment", controllers.Payment)

			auth.GET("/partner/list", controllers.GetPartnerList)
			auth.GET("/partner/all", controllers.GetPartnerAll)
			auth.GET("/partner/:id", controllers.GetPartnerDetail)
			auth.POST("/partner", controllers.CreatePartner)
			auth.PUT("/partner/:id", controllers.UpdatePartner)
			auth.DELETE("/partner/:id", controllers.DeletePartner)

			auth.GET("/log/list", controllers.GetLogList)
		}
	}

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/login")
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login/index.html", gin.H{})
	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "layout/index.html", gin.H{})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard/index.html", gin.H{})
	})

	r.GET("/user/list", func(c *gin.Context) {
		c.HTML(200, "user/list.html", gin.H{})
	})

	r.GET("/user/form", func(c *gin.Context) {
		c.HTML(200, "user/form.html", gin.H{})
	})

	r.GET("/user/profile", func(c *gin.Context) {
		c.HTML(200, "user/profile.html", gin.H{})
	})

	r.GET("/user/password", func(c *gin.Context) {
		c.HTML(200, "user/password.html", gin.H{})
	})

	r.GET("/role/list", func(c *gin.Context) {
		c.HTML(200, "role/list.html", gin.H{})
	})

	r.GET("/role/form", func(c *gin.Context) {
		c.HTML(200, "role/form.html", gin.H{})
	})

	r.GET("/menu/list", func(c *gin.Context) {
		c.HTML(200, "menu/list.html", gin.H{})
	})

	r.GET("/config/list", func(c *gin.Context) {
		c.HTML(200, "config/list.html", gin.H{})
	})

	r.GET("/lot/list", func(c *gin.Context) {
		c.HTML(200, "lot/list.html", gin.H{})
	})

	r.GET("/space/list", func(c *gin.Context) {
		c.HTML(200, "space/list.html", gin.H{})
	})

	r.GET("/recognize/index", func(c *gin.Context) {
		c.HTML(200, "recognize/index.html", gin.H{})
	})

	r.GET("/recognize/entry", func(c *gin.Context) {
		c.HTML(200, "recognize/entry.html", gin.H{})
	})

	r.GET("/recognize/exit", func(c *gin.Context) {
		c.HTML(200, "recognize/exit.html", gin.H{})
	})

	r.GET("/vehicle/list", func(c *gin.Context) {
		c.HTML(200, "vehicle/list.html", gin.H{})
	})

	r.GET("/record/list", func(c *gin.Context) {
		c.HTML(200, "record/list.html", gin.H{})
	})

	r.GET("/finance/list", func(c *gin.Context) {
		c.HTML(200, "finance/list.html", gin.H{})
	})

	r.GET("/partner/list", func(c *gin.Context) {
		c.HTML(200, "partner/list.html", gin.H{})
	})

	r.GET("/log/list", func(c *gin.Context) {
		c.HTML(200, "log/list.html", gin.H{})
	})

	return r
}
