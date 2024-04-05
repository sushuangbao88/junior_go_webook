package main

import (
	"strings"
	"time"

	"example.com/basic-gin/webook/internal/repository"
	"example.com/basic-gin/webook/internal/repository/dao"
	"example.com/basic-gin/webook/internal/service"
	"example.com/basic-gin/webook/internal/web"
	"example.com/basic-gin/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	dao.InitTables(db) //自动建表
	server := initWebServer()
	//(分散式)注册(初始化)「用户」路由
	initRegisterUserHandler(db, server)

	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	//中间件的可以去仓库：https://github.com/gin-gonic/contrib 中去查看
	//middleware：解决跨域问题的preflight（预检请求）
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"content-Type"},
		//AllowMethods: []string{"POST"}, //跨域允许的请求方法
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "contentreview.com")
		}, //允许的请求源判断方法
		MaxAge: 12 * time.Hour,
	}))

	//登陆校验中间件
	loginW := &middleware.LoginMiddlewareBuilder{}
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store), loginW.CheckLogin())

	return server
}

func initRegisterUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)

	hdl.RegisterRoutes(server) //注册“用户”路由
}
