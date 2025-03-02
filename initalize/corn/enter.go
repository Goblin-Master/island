package corn

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"time"
)

func Cron() {
	zone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logrus.Warn("加载时区错误:%v", err)
	}
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(zone))
	// 每天2点同步文章数据
	_, err = crontab.AddFunc("* * 2 * * *", SyncArticle)
	if err != nil {
		logrus.Warn("添加文章内容同步定时任务错误:%v", err)
	}
	crontab.Start()
}
