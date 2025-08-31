package configs

var Conf = new(Config)

type Config struct {
	App   ApplicationConfig `mapstructure:"app"`
	Log   LoggerConfig      `mapstructure:"log"`
	DB    DBConfig          `mapstructure:"database"`
	Redis RedisConfig       `mapstructure:"redis"`
	QQ    QQConfig          `mapstructure:"qq"`
	AI    []AIConfig        `mapstructure:"ai"`
}

type ApplicationConfig struct {
	Reload      bool   `mapstructure:"reload"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Env         string `mapstructure:"env"`
	LogfilePath string `mapstructure:"logfilePath"`
	ImagesPath  string `mapstructure:"imagesPath"`
}
type LoggerConfig struct {
	Reload   bool   `mapstructure:"reload"`
	Level    int8   `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Director string `mapstructure:"director"`
	ShowLine bool   `mapstructure:"show-line"`
}

type DBConfig struct {
	Reload      bool   `mapstructure:"reload"`
	Driver      string `mapstructure:"driver"`
	AutoMigrate bool   `mapstructure:"migrate"`
	Dsn         string `mapstructure:"dsn"`
}
type RedisConfig struct {
	Reload   bool   `mapstructure:"reload"`
	Enable   bool   `mapstructure:"enable"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type KafkaConfig struct {
	Reload bool   `mapstructure:"reload"`
	host   string `mapstructure:"host"`
	port   int    `mapstructure:"port"`
}

type QQConfig struct {
	AppID       string `mapstructure:"appID"`
	AppKey      string `mapstructure:"appKey"`
	RedirectUrl string `mapstructure:"redirectUrl"`
}

type AIConfig struct {
	Model  string `mapstructure:"model"`
	ApiKey string `mapstructure:"apiKey"`
	ApiUrl string `mapstructure:"apiUrl"`
}
