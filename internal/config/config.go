package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var c = Config{
	Server: ServerConfig{
		Port:      "0",
		TcpPort:   "0",
		Id:        "0",
		Mode:      "debug",
		DeployEnv: "",
	},
	LoggerConfig: LoggerConfig{
		Level:      "info",
		Filename:   "",
		MaxSize:    0,
		MaxBackups: 0,
		MaxAge:     0,
		Compress:   false,
	},
	MysqlConfig: MysqlConfig{
		DataSource: "",
		MaxIdle:    0,
		MaxOpen:    0,
	},
	RedisConfig: RedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	},
}

func GetConfig() *Config {
	return &c
}

type Config struct {
	Server           ServerConfig     `mapstructure:"server"`
	LoggerConfig     LoggerConfig     `mapstructure:"logger"`
	MysqlConfig      MysqlConfig      `mapstructure:"mysql"`
	RedisConfig      RedisConfig      `mapstructure:"redis"`
	OssConfig        OssConfig        `mapstructure:"oss"`
	OpenSearchConfig OpenSearchConfig `mapstructure:"opensearch"`
	PabConfig        PabConfig        `mapstructure:"pab"`
	AlipayConfig     AlipayConfig     `mapstructure:"alipay"`
	WxpayConfig      WxpayConfig      `mapstructure:"wxpay"`
}

type ServerConfig struct {
	Port      string `mapstructure:"port"`
	TcpPort   string `mapstructure:"tcp_port"`
	WsPort    string `mapstructure:"ws_port"`
	Id        string `mapstructure:"id"`
	Mode      string `mapstructure:"mode"`
	DeployEnv string `mapstructure:"deploy_env"`
}

type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type MysqlConfig struct {
	DataSource string `mapstructure:"data_source"`
	MaxIdle    int    `mapstructure:"max_idle"`
	MaxOpen    int    `mapstructure:"max_open"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type OssConfig struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	RoleArn         string `mapstructure:"role_arn"`
	Endpoint        string `mapstructure:"endpoint"`
	BucketName      string `mapstructure:"bucket_name"`
	CdnDomain       string `mapstructure:"cdn_domain"`
}

type OpenSearchConfig struct {
	AppName string `mapstructure:"app_name"`
}

type PabConfig struct {
	Url               string `mapstructure:"url"`
	AppId             string `mapstructure:"app_id"`
	CertPath          string `mapstructure:"cert_path"`
	PfxPath           string `mapstructure:"pfx_path"`
	Dn                string `mapstructure:"dn"`
	UserShortNo       string `mapstructure:"user_short_no"`
	FundSummaryAcctNo string `mapstructure:"fund_summary_acct_no"`
	FileUrl           string `mapstructure:"file_url"`
	Passwd            string `mapstructure:"passwd"`
}

type AlipayConfig struct {
	AppId               string `mapstructure:"app_id"`
	AppSecret           string `mapstructure:"app_secret"`
	AppCertPublicKey    string `mapstructure:"app_cert_public_key"`
	AlipayCertPublicKey string `mapstructure:"alipay_cert_public_key"`
	AlipayRootCert      string `mapstructure:"alipay_root_cert"`
	NotifyUrl           string `mapstructure:"notify_url"`
}

type WxpayConfig struct {
	MchId     string `mapstructure:"mch_id"`
	AppId     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	ApiKey    string `mapstructure:"api_key"`
	Apiv3Key  string `mapstructure:"apiv3_key"`
	P12Cert   string `mapstructure:"p12_cert"`
	NotifyUrl string `mapstructure:"notify_url"`
}

func LoadConfig(cfgFile string) (err error) {
	if cfgFile == "" {
		err = fmt.Errorf("config file is empty")
		return
	}

	viper.SetConfigFile(cfgFile)
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&c); err != nil {
		return
	}

	return nil
}
