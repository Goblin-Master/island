package initalize

import (
	"tgwp/cmd/flags"
	"tgwp/global"
	"tgwp/initalize/corn"
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
func InitPath() {
	global.Path = utils.GetRootPath("")
}
