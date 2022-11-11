package config

import (
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"github.com/tuanp/go-gin-boilerplate/pkg/constants"
)

type LogType int32

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	Zap    LogType = 0
	Logrus LogType = 1
)

type (
	Config struct {
		Mysql    MysqlConfig
		HTTP     HTTPConfig
		Redis    RedisConfig
		Logger   LoggerConfig
		Server   ServerConfig
		CacheTTL time.Duration `mapstructure:"ttl"`
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	// MysqlConfig is settings of a MySQL server. It contains almost same fields as mysql.Config,
	// but with some different field names and tags.
	MysqlConfig struct {
		Username  string            `yaml:"username" mapstructure:"username"`
		Password  string            `yaml:"password" mapstructure:"password"`
		Protocol  string            `yaml:"protocol" mapstructure:"protocol"`
		Address   string            `yaml:"address" mapstructure:"address"`
		Database  string            `yaml:"database" mapstructure:"database"`
		Params    map[string]string `yaml:"params" mapstructure:"params"`
		Collation string            `yaml:"collation" mapstructure:"collation"`
		Loc       *time.Location    `yaml:"location" mapstructure:"loc"`
		TLSConfig string            `yaml:"tlsConfig" mapstructure:"tlsConfig"`

		Timeout      time.Duration `yaml:"timeout" mapstructure:"timeout"`
		ReadTimeout  time.Duration `yaml:"readTimeout" mapstructure:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout" mapstructure:"writeTimeout"`

		AllowAllFiles           bool   `yaml:"allowAllFiles" mapstructure:"allowAllFiles"`
		AllowCleartextPasswords bool   `yaml:"allowCleartextPasswords" mapstructure:"allowCleartextPasswords"`
		AllowOldPasswords       bool   `yaml:"allowOldPasswords" mapstructure:"allowOldPasswords"`
		ClientFoundRows         bool   `yaml:"clientFoundRows" mapstructure:"clientFoundRows"`
		ColumnsWithAlias        bool   `yaml:"columnsWithAlias" mapstructure:"columnsWithAlias"`
		InterpolateParams       bool   `yaml:"interpolateParams" mapstructure:"interpolateParams"`
		MultiStatements         bool   `yaml:"multiStatements" mapstructure:"multiStatements"`
		ParseTime               bool   `yaml:"parseTime" mapstructure:"parseTime"`
		GoogleAuthFile          string `yaml:"googleAuthFile" mapstructure:"googleAuthFile"`
	}

	RedisConfig struct {
		Addr         string `mapstructure:"addr"`
		Password     string `mapstructure:"password"`
		DB           int    `mapstructure:"db"`
		PoolSize     int    `mapstructure:"poolSize"`
		MinIdleConns int    `mapstructure:"minIdleConns"`
		PoolTimeout  int    `mapstructure:"poolTimeout"`
	}

	LoggerConfig struct {
		Development       bool    `yaml:"development" mapstructure:"development"`
		DisableCaller     bool    `yaml:"disableCaller" mapstructure:"disableCaller"`
		DisableStacktrace bool    `yaml:"disableStacktrace" mapstructure:"disableStacktrace"`
		Encoding          string  `yaml:"encoding" mapstructure:"encoding"`
		LogLevel          string  `yaml:"level" mapstructure:"level"`
		LogType           LogType `yaml:"logType" mapstructure:"logType"`
	}

	ServerConfig struct {
		AppVersion string `yaml:"appVersion" mapstructure:"appVersion"`
		Mode       string `yaml:"mode" mapstructure:"mode"`
		Debug      bool   `yaml:"debug" mapstructure:"debug"`
	}
)

// FormatDSN returns MySQL DSN from settings.
func (m *MysqlConfig) FormatDSN() string {
	um := &mysql.Config{
		User:                    m.Username,
		Passwd:                  m.Password,
		Net:                     m.Protocol,
		Addr:                    m.Address,
		DBName:                  m.Database,
		Params:                  m.Params,
		Collation:               m.Collation,
		Loc:                     m.Loc,
		TLSConfig:               m.TLSConfig,
		Timeout:                 m.Timeout,
		ReadTimeout:             m.ReadTimeout,
		WriteTimeout:            m.WriteTimeout,
		AllowAllFiles:           m.AllowAllFiles,
		AllowCleartextPasswords: m.AllowCleartextPasswords,
		AllowOldPasswords:       m.AllowOldPasswords,
		ClientFoundRows:         m.ClientFoundRows,
		ColumnsWithAlias:        m.ColumnsWithAlias,
		InterpolateParams:       m.InterpolateParams,
		MultiStatements:         m.MultiStatements,
		ParseTime:               m.ParseTime,
		AllowNativePasswords:    true,
	}
	return um.FormatDSN()
}

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init(configsDir string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Override configuration from environment settings
	if err := setFromEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("cache.ttl", &cfg.CacheTTL); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("mysql", &cfg.Mysql); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("logger", &cfg.Logger); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("server", &cfg.Server); err != nil {
		return err
	}

	return viper.UnmarshalKey("http", &cfg.HTTP)
}

func setFromEnv(cfg *Config) error {

	if err := envconfig.Process("mysql", &cfg.Mysql); err != nil {
		return err
	}

	if err := envconfig.Process("redis", &cfg.Redis); err != nil {
		return err
	}

	if err := envconfig.Process("http", &cfg.HTTP); err != nil {
		return err
	}

	return nil
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("base")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == constants.EnvDev {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
}
