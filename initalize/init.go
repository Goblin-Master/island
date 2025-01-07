package initalize

import (
	"tgwp/global"
	"tgwp/util"
)

func Init() {
	introduce()
	InitLog(global.Config)
	InitPath()
	InitConfig()
	InitLog(global.Config)
	InitDataBase(*global.Config)
	InitRedis(*global.Config)
}
func InitPath() {
	global.Path = util.GetRootPath("")
}
