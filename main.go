package main

import (
	"strings"
	"time"

	"example.com/basic-gin/webook/internal/repository"
	"example.com/basic-gin/webook/internal/repository/cache"
	"example.com/basic-gin/webook/internal/repository/dao"
	"example.com/basic-gin/webook/internal/service"
	"example.com/basic-gin/webook/internal/service/sms"
	"example.com/basic-gin/webook/internal/web"
	"example.com/basic-gin/webook/internal/web/middleware"
	"example.com/basic-gin/webook/pkg/ginx/middleware/ratelimit"
	"example.com/basic-gin/webook/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	goRedis "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	rdb := initRedisClient()

	dao.InitTables(db) //自动建表
	server := initWebServer(rdb)
	//(分散式)注册(初始化)「用户」路由
	initRegisterUserHandler(db, server, rdb)

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

func initRedisClient() *goRedis.Client {
	return goRedis.NewClient(&goRedis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func initWebServer(rdb *goRedis.Client) *gin.Engine {
	server := gin.Default()

	//中间件的可以去仓库：https://github.com/gin-gonic/contrib 中去查看
	//middleware：解决跨域问题的preflight（预检请求）
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"content-Type", "Authorization"},
		//这个是允许前端访问你的后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token"},
		//AllowMethods: []string{"POST"}, //跨域允许的请求方法
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "contentreview.com")
		}, //允许的请求源判断方法
		MaxAge: 12 * time.Hour,
	}))

	//限流中间件
	/* 因为要要使用压测，暂时先限流一下*/
	rswLimiter := limiter.NewRedisSlidingWindowLimiter(rdb, time.Second, 1)
	server.Use(ratelimit.NewBuilder(rswLimiter).Build())

	//useSession(server) //用户校验：session
	useJWT(server) //用户校验：JWT

	return server
}

func initRegisterUserHandler(db *gorm.DB, server *gin.Engine, rdb *goRedis.Client) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)

	cc := cache.NewCodeRedisCache(rdb)
	cr := repository.NewCodeRepository(cc)
	ss := sms.NewService("local")
	cs := service.NewCodeService(cr, ss)
	hdl := web.NewUserHandler(us, cs)

	hdl.RegisterRoutes(server) //注册“用户”路由
}

func useJWT(server *gin.Engine) {
	loginWJWT := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(loginWJWT.CheckLogin())
}

func useSession(server *gin.Engine) {
	//登陆校验中间件
	loginW := &middleware.LoginMiddlewareBuilder{}

	//方式一：基于cookie实现的session：不安全，
	//store := cookie.NewStore([]byte("secret"))

	//方式二：基于内存实现的session：适用于单机部署
	//store := memstore.NewStore([]byte("yfmgmDb7VGQZh0fRQnqCzA2V51fGJVUY"), []byte("9bErGWxgl7P7mddPM3fTifhD3hWbGF7e"))

	//方式三：基于第三方存储实现session：可以用于多实例，下面是以redis为例，还可以使用memchache的，甚至sql
	//参数一：最大空闲连接数
	//参数二：传输层协议：
	//----tcp（面向连接协议，通过三次握手建立可靠连接；优：可靠，有序；缺：开销大，速度慢），大部分的选择
	//----udp（无连接协议，数据以数据报的形式独立发送；优：开销小，速度快；缺：不可靠，顺序性差），基本不会使用
	//参数五：authentication key,身份认证
	//参数六：encryption key，数据加密
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("yfmgmDb7VGQZh0fRQnqCzA2V51fGJVUY"), []byte("9bErGWxgl7P7mddPM3fTifhD3hWbGF7e"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), loginW.CheckLogin())
}
