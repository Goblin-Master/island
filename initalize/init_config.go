package initalize

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/goccy/go-json"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
	"sync"
	"tgwp/configs"
	"tgwp/global"
	"tgwp/log/zlog"
	"time"
)

var (
	configLock = &sync.RWMutex{}
	// 需要动态监听的环境变量
	criticalEnvKeys = []string{
		"LOG_LEVEL",
		"CONFIG_PATH", // 环境变量配置路径
	}
)

func InitConfig() {
	initTimeZone()
	initConfigSources()
	loadConfig()
	startConfigWatcher()
}

func initTimeZone() {
	// 兼容容器化环境时区配置
	if tz := os.Getenv("TZ"); tz != "" {
		loc, err := time.LoadLocation(tz)
		if err == nil {
			time.Local = loc
			return
		}
	}
	time.Local = time.FixedZone("CST", 8*3600)
}

func initConfigSources() {
	// 初始化配置源优先级：环境变量 > 命令行 > 配置文件
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ISLAND")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var configPath string
	flag.StringVar(&configPath, "c", getDefaultConfigPath(), "配置文件路径")
	flag.Parse()

	// 环境变量覆盖命令行参数
	if envPath := viper.GetString("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	viper.SetConfigFile(configPath)
}

func loadConfig() {
	configLock.Lock()
	defer configLock.Unlock()

	if err := viper.ReadInConfig(); err != nil {
		zlog.Panicf("配置加载失败: %v", err)
	}

	if err := viper.Unmarshal(&configs.Conf); err != nil {
		zlog.Panicf("配置解析失败: %v", err)
	}

	global.Config = configs.Conf
	zlog.Infof("成功加载配置: %+v", global.Config)
}

func startConfigWatcher() {
	// 文件配置监听
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		zlog.Warnf("配置文件变更: %s", e.Name)
		reloadConfig()
	})

	// 环境变量监听
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			checkEnvChanges()
		}
	}()
}

func reloadConfig() {
	configLock.Lock()
	defer configLock.Unlock()

	oldConfig := deepCopyConfig(global.Config)

	newConfig := new(configs.Config)
	if err := viper.Unmarshal(newConfig); err != nil {
		zlog.Errorf("配置重载失败: %v", err)
		return
	}

	global.Config = newConfig
	configs.Conf = newConfig
	zlog.Infof("配置更新完成: %+v", global.Config)

	// 触发关联初始化
	if !reflect.DeepEqual(oldConfig, newConfig) {
		Eve()
		ReInit(newConfig)
	}
}

func checkEnvChanges() {
	for _, key := range criticalEnvKeys {
		envVal := viper.GetString(key)
		if viper.Get(key) != envVal {
			zlog.Warnf("环境变量变更: %s=%s", key, envVal)
			viper.Set(key, envVal)
			reloadConfig()
		}
	}
}

func getDefaultConfigPath() string {
	// 多环境默认配置路径
	if env := os.Getenv("APP_ENV"); env != "" {
		return global.Path + "/config-" + env + ".yaml"
	}
	return global.Path + global.DEFAULT_CONFIG_FILE_PATH
}

// 深拷贝辅助函数
func deepCopyConfig(src *configs.Config) *configs.Config {
	if src == nil {
		return nil
	}
	dst := new(configs.Config)
	b, _ := json.Marshal(src)
	_ = json.Unmarshal(b, dst)
	return dst
}
