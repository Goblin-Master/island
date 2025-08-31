package manager

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"tgwp/middleware"
)

//主要管理路由组和中间件的注册

// PathHandler 是一个用于注册路由组的函数类型
type PathHandler func(rg *gin.RouterGroup)

// Middleware 是一个用于生成中间件的函数类型
type Middleware func() gin.HandlerFunc

// RouteManager 管理不同的路由组，按业务功能分组
type RouteManager struct {
	LoginRouter    *gin.RouterGroup // 登录相关的路由组
	CommonRouter   *gin.RouterGroup //通用功能相关的路由组
	IslandRouter   *gin.RouterGroup //岛屿相关的路由组
	AiRouter       *gin.RouterGroup //ai相关的路由组
	QuestionRouter *gin.RouterGroup //刷题功能相关的路由组
	ArticleRouter  *gin.RouterGroup //文章相关的路由组
	ChatRouter     *gin.RouterGroup //聊天相关的路由组
	PKRouter       *gin.RouterGroup //pk相关的路由组
	UserRouter     *gin.RouterGroup //用户相关的路由组
}

// NewRouteManager 创建一个新的 RouteManager 实例，包含各业务功能的路由组
func NewRouteManager(router *gin.Engine) *RouteManager {
	return &RouteManager{
		LoginRouter:    router.Group("/api/login"),    // 初始化登录路由组
		CommonRouter:   router.Group("/api/common"),   //通用功能相关的路由组
		IslandRouter:   router.Group("/api/island"),   //岛屿相关的路由组
		AiRouter:       router.Group("/api/ai"),       //ai相关的路由组
		QuestionRouter: router.Group("/api/question"), //刷题功能相关的路由组
		ArticleRouter:  router.Group("/api/article"),  //文章相关的路由组
		ChatRouter:     router.Group("/api/chat"),     //聊天相关的路由组
		PKRouter:       router.Group("/api/pk"),       //pk相关的路由组
		UserRouter:     router.Group("/api/user"),     //用户相关的路由组
	}
}

// RegisterLoginRoutes 注册登录相关的路由处理函数
func (rm *RouteManager) RegisterLoginRoutes(handler PathHandler) {
	handler(rm.LoginRouter)
}

// RegisterCommonRoutes 通用功能相关的路由组
func (rm *RouteManager) RegisterCommonRoutes(handler PathHandler) {
	handler(rm.CommonRouter)
}

// RegisterIslandRoutes 注册岛屿相关的路由处理函数
func (rm *RouteManager) RegisterIslandRoutes(handler PathHandler) {
	handler(rm.IslandRouter)
}

// RegisterAiRoutes 注册ai相关的路由处理函数
func (rm *RouteManager) RegisterAiRoutes(handler PathHandler) {
	handler(rm.AiRouter)
}

// RegisterQuestionRoutes 注册刷题功能相关的路由处理函数
func (rm *RouteManager) RegisterQuestionRoutes(handler PathHandler) {
	handler(rm.QuestionRouter)
}

// RegisterArticleRoutes 注册文章相关的路由处理函数
func (rm *RouteManager) RegisterArticleRoutes(handler PathHandler) {
	handler(rm.ArticleRouter)
}

// RegisterChatRoutes 注册聊天相关的路由处理函数
func (rm *RouteManager) RegisterChatRoutes(handler PathHandler) {
	handler(rm.ChatRouter)
}

// RegisterPKRoutes 注册PK相关的路由处理函数
func (rm *RouteManager) RegisterPKRoutes(handler PathHandler) {
	handler(rm.PKRouter)
}

// RegisterUserRoutes 用户相关的路由处理函数
func (rm *RouteManager) RegisterUserRoutes(handler PathHandler) {
	handler(rm.UserRouter)
}

// RegisterMiddleware 根据组名为对应的路由组注册中间件
// group 参数为 "login"、"profile"、"team"或"Common"，分别对应不同的路由组
func (rm *RouteManager) RegisterMiddleware(group string, middleware Middleware) {
	switch group {
	case "login":
		rm.LoginRouter.Use(middleware())
	case "common":
		rm.CommonRouter.Use(middleware())
	case "island":
		rm.IslandRouter.Use(middleware())
	case "ai":
		rm.AiRouter.Use(middleware())
	case "article":
		rm.ArticleRouter.Use(middleware())
	case "question":
		rm.QuestionRouter.Use(middleware())
	case "chat":
		rm.ChatRouter.Use(middleware())
	case "pk":
		rm.CommonRouter.Use(middleware())
	case "user":
		rm.UserRouter.Use(middleware())
	}

}

// RequestGlobalMiddleware 注册全局中间件，应用于所有路由
func RequestGlobalMiddleware(r *gin.Engine) {
	r.Use(requestid.New())
	r.Use(middleware.AddTraceId())
	r.Use(middleware.Cors())
}
