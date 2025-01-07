package initalize

import (
	"tgwp/configs"
	"tgwp/global"
	"tgwp/pkg/database"
	"tgwp/pkg/mysqlx"
	"tgwp/pkg/redisx"

	"tgwp/log/zlog"
)

func InitDataBase(config configs.Config) {
	switch config.DB.Driver {
	case "mysql":
		database.InitDataBases(mysqlx.NewMySql(), config)
	default:

		zlog.Fatalf("不支持的数据库驱动：%s", config.DB.Driver)
	}
	if config.App.Env != "pro" {
		err := global.DB.AutoMigrate()

		//迁移数据库所有的表
		migrateTables()

		if err != nil {
			zlog.Fatalf("数据库迁移失败！")
		}
	}
	zlog.Infof("数据库初始化成功！")
}
func InitRedis(config configs.Config) {
	if config.Redis.Enable {
		var err error
		global.Rdb, err = redisx.GetRedisClient(config)
		if err != nil {
			zlog.Errorf("无法初始化Redis : %v", err)
		}
	} else {
		zlog.Warnf("不使用Redis")
	}

}

func migrateTables() {
	//自动迁移某一个表，确保表结构存在
	//global.DB.AutoMigrate(&model.xxx)
}