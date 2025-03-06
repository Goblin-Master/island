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
		rg.POST("/qq", api.QQLogin)
	})

	routeManager.RegisterIslandRoutes(func(rg *gin.RouterGroup) {
		rg.GET("/", api.GetIsland)                                  //获取岛屿列表
		rg.POST("/", middleware.Authentication, api.CreateIsland)   //创建岛屿
		rg.PUT("/", middleware.Authentication, api.ModifyIsland)    //修改岛屿
		rg.DELETE("/", middleware.Authentication, api.DeleteIsland) //删除岛屿
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
		rg.POST("/", api.ArticleCreate)
		rg.GET("/", api.ArticleList)
		rg.DELETE("/", api.ArticleDelete)
		rg.GET("/digg", api.ArticleDigg)
		rg.GET("/collect", api.ArticleCollect)
		rg.GET("/owner", api.ArticleListByUerID)
		rg.GET("/collect/owner", api.ArticleCollectList)
		rg.GET("/digg/owner", api.ArticleDiggList)
	})

	routeManager.RegisterChatRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/send", api.ChatSendMessage)
		rg.GET("/get", api.ChatGetMessages)
	})
}
