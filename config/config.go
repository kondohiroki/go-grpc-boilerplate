package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

var config *Config
var m sync.Mutex

type Config struct {
	Env        string     `yaml:"env"`
	App        App        `yaml:"app"`
	GRPCServer GRPCServer `yaml:"grpcServer"`
	Log        Log        `yaml:"log"`
	Scheduler  Scheduler  `yaml:"scheduler"`
	Schedules  []Schedule `yaml:"schedules"`
	Postgres   Postgres   `yaml:"postgres"`
	Redis      []Redis    `yaml:"redis"`
	Services   Services   `yaml:"services"`
}

type GRPCServer struct {
	Port           int    `yaml:"port"`
	Reflection     bool   `yaml:"reflection"`
	MaxSendMsgSize int    `yaml:"maxSendMsgSize"` // in MB
	MaxRecvMsgSize int    `yaml:"maxRecvMsgSize"` // in MB
	UseTls         bool   `yaml:"useTls"`
	TlsCertFile    string `yaml:"tlsCertFile"`
	TlsKeyFile     string `yaml:"tlsKeyFile"`
}

type Log struct {
	Level           string `yaml:"level"`
	StacktraceLevel string `yaml:"stacktraceLevel"`
	FileEnabled     bool   `yaml:"fileEnabled"`
	FileSize        int    `yaml:"fileSize"`
	FilePath        string `yaml:"filePath"`
	FileCompress    bool   `yaml:"fileCompress"`
	MaxAge          int    `yaml:"maxAge"`
	MaxBackups      int    `yaml:"maxBackups"`
}

type Label struct {
	En string `json:"en"`
	Th string `json:"th"`
}

type App struct {
	Name     string `yaml:"name"`
	NameSlug string `yaml:"nameSlug"`
}

type Postgres struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Schema          string `yaml:"schema"`
	MaxConnections  int32  `yaml:"maxConnections"`
	MaxConnIdleTime int32  `yaml:"maxConnIdleTime"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type Scheduler struct {
	Timezone string `yaml:"timezone"`
}

type Schedule struct {
	Job       string `yaml:"job"`
	Cron      string `yaml:"cron"`
	IsEnabled bool   `yaml:"isEnabled"`
}

type Services struct {
	Email   EmailConfig `yaml:"email"`
	S3      S3Config    `yaml:"s3"`
	Example Example     `yaml:"example"`
}

type Example struct {
	Endpoint           string         `yaml:"endpoint"`
	Authentication     Authentication `yaml:"authentication"`
	TransactionOwnerID string         `yaml:"transactionOwnerID"`
	Channel            string         `yaml:"channel"`
}

type S3Config struct {
	AwsRegion          string `yaml:"awsRegion"`
	AwsAccessKeyID     string `yaml:"awsAccessKeyID"`
	AwsSecretAccessKey string `yaml:"awsSecretAccessKey"`
	Bucket             string `yaml:"bucket"`
	Path               string `yaml:"path"`
}

type EmailConfig struct {
	SmtpHost string `yaml:"smtpHost"`
	SmtpPort int    `yaml:"smtpPort"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type Authentication struct {
	Endpoint string `yaml:"endpoint"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func GetConfig() *Config {
	return config
}

func SetConfig(configFile string) {
	m.Lock()
	defer m.Unlock()

	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error getting config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct, ", err)
	}
}
