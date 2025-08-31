package initalize

import (
	"tgwp/cmd/flags"
	"tgwp/configs"
	"tgwp/global"
	"tgwp/initalize/corn"
	"tgwp/log/zlog"
	"tgwp/utils"
)

func Init() {
	flags.Parse()
	introduce()
	InitLog(global.Config)
	InitPath()
	InitConfig()
	InitLog(global.Config)
	InitDataBase(*global.Config)
	InitRedis(*global.Config)
	flags.Migrate()
	corn.Cron() // 开启定时任务
	flags.Run() // 会通过特殊手段执行数据库表的迁移
}

func ReInit(config *configs.Config) {
	flags.Parse()
	introduce()
	if config.Log.Reload {
		zlog.Warnf("日志配置已更新，重新加载日志配置")
		InitLog(config)
	}
	if config.Redis.Reload {
		zlog.Warnf("Redis配置已更新，重新加载Redis配置")
		InitRedis(*config)
	}
	if config.DB.Reload {
		zlog.Warnf("数据库配置已更新，重新加载数据库配置")
		InitDataBase(*config)
	}
}

func InitPath() {
	global.Path = utils.GetRootPath("")
}
