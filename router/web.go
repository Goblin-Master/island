package routerg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"tgwp/configs"
	"tgwp/internal/api"
	"tgwp/log/zlog"
	"tgwp/manager"
	"tgwp/middleware"
)

// RunServer 启动服务器 路由层
func RunServer() {
	r, err := listen()
	if err != nil {
		zlog.Errorf("Listen error: %v", err)
		panic(err.Error())
	}
	err = r.Run(fmt.Sprintf("%s:%d", configs.Conf.App.Host, configs.Conf.App.Port)) // 启动 Gin 服务器
	if err != nil {
		zlog.Errorf("Run error: %v", err)
		panic(err.Error())
	}
}

// listen 配置 Gin 服务器
func listen() (*gin.Engine, error) {
	r := gin.Default() // 创建默认的 Gin 引擎
	// 注册全局中间件（例如获取 Trace ID）
	manager.RequestGlobalMiddleware(r)
	//配置静态路由，用于访问上传的文件
	r.Static("/uploads", "uploads")
	// 创建 RouteManager 实例
	routeManager := manager.NewRouteManager(r)
	// 注册各业务路由组的具体路由
	registerRoutes(routeManager)
	return r, nil
}

// registerRoutes 注册各业务路由的具体处理函数
func registerRoutes(routeManager *manager.RouteManager) {

	routeManager.RegisterCommonRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/token", api.RefreshToken)
		rg.POST("/images", middleware.Authentication, api.UploadImages)
		rg.DELETE("/images", middleware.Authentication, api.DeleteImages)
		rg.GET("/images", api.GetImages)
	})

	routeManager.RegisterLoginRoutes(func(rg *gin.RouterGroup) {
		rg.GET("/qq", api.QQLogin)
		rg.GET("/auto", middleware.Authentication, api.AutoLogin)
	})

	routeManager.RegisterIslandRoutes(func(rg *gin.RouterGroup) {
		rg.GET("/", api.GetIsland)                                  //获取岛屿列表
		rg.POST("/", middleware.Authentication, api.CreateIsland)   //创建岛屿
		rg.PUT("/", middleware.Authentication, api.ModifyIsland)    //修改岛屿
		rg.DELETE("/", middleware.Authentication, api.DeleteIsland) //删除岛屿
		rg.GET("/detail", api.IslandDetail)                         //岛屿详情
	})

	routeManager.RegisterAiRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/analysis", middleware.Authentication, api.AiGenerateAbstract)
		rg.GET("/chat", middleware.Authentication, api.AiChatStream)
		rg.DELETE("/", middleware.Authentication, api.ClearHistory)
	})

	routeManager.RegisterQuestionRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/run-code", api.RunCode)

		rg.POST("/create-question", api.CreateQuestion)
		rg.POST("/create-question-bank", api.CreateQuestionBank)
		rg.POST("/add-question", api.AddQuestion)
		rg.POST("/add-question-bank", api.AddQuestionBank)

		rg.GET("/get-question", api.GetQuestion)
		rg.GET("/get-question-bank", api.GetQuestionBank)
		rg.GET("/get-question-list", api.GetQuestionList)
		rg.GET("/get-question-bank-list", api.GetQuestionBankList)

	})

	routeManager.RegisterArticleRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/", middleware.Authentication, api.ArticleCreate)
		rg.GET("/", api.ArticleList)
		rg.DELETE("/", middleware.Authentication, api.ArticleDelete)
		rg.GET("/digg", middleware.Authentication, api.ArticleDigg)
		rg.GET("/collect", middleware.Authentication, api.ArticleCollect)
		rg.GET("/owner", middleware.Authentication, api.ArticleListByUerID)
		rg.GET("/collect/owner", middleware.Authentication, api.ArticleCollectList)
		rg.GET("/digg/owner", middleware.Authentication, api.ArticleDiggList)
	})

	routeManager.RegisterChatRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/send", api.ChatSendMessage)
		rg.GET("/get", api.ChatGetMessages)
	})

	routeManager.RegisterPKRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/matching", middleware.Authentication, api.PKMatching)
		rg.GET("/room-info", middleware.Authentication, api.GetRoomInfo)
		rg.POST("/submit", middleware.Authentication, api.SubmitQuestion)
		rg.POST("/set-rule", middleware.Authentication, api.PKSetRule)
		rg.GET("/get-rule", middleware.Authentication, api.PKGetRule)
	})

	routeManager.RegisterUserRoutes(func(rg *gin.RouterGroup) {
		rg.GET("/detail", middleware.Authentication, api.UserDetail)
		rg.GET("/focus", middleware.Authentication, api.FocusUser)
		rg.GET("/focus/list", middleware.Authentication, api.FocusList)
	})

	// 两个可以用来测试的用户 Token
	//fmt.Println(jwtUtils.ForTest(1898984453741481984, time.Hour*24*365*100))
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOjE4OTg5ODQ0NTM3NDE0ODE5ODQsImNsYXNzIjoiYXRva2VuIiwiaXNzIjoiaXNsYW5kIiwiZXhwIjo0ODk1MjU2Nzc4LCJuYmYiOjE3NDE2NTY3Nzh9.ptTgeLknT9lR_Tbj9GIVISLRFImA_x4S-oLMdKupfqs
	//fmt.Println(jwtUtils.ForTest(1898900897681903616, time.Hour*24*365*100))
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOjE4OTg5MDA4OTc2ODE5MDM2MTYsImNsYXNzIjoiYXRva2VuIiwiaXNzIjoiaXNsYW5kIiwiZXhwIjo0ODk1Mjc1OTEzLCJuYmYiOjE3NDE2NzU5MTN9.FMbMnACeKAZf5bGi02MVk6mpMAbQIChl65bAapP1o30
}
