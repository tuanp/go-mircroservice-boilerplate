package config

import (
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	EnvLocal = "local"
	Prod     = "prod"
)

type (
	Config struct {
		Environment string
		Mysql       MysqlConfig
		HTTP        HTTPConfig
		Redis       RedisConfig
		Logger      LoggerConfig
		Server      ServerConfig
		CacheTTL    time.Duration `mapstructure:"ttl"`
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
		Development       bool   `yaml:"development" mapstructure:"development"`
		DisableCaller     bool   `yaml:"disableCaller" mapstructure:"disableCaller"`
		DisableStacktrace bool   `yaml:"disableStacktrace" mapstructure:"disableStacktrace"`
		Encoding          string `yaml:"encoding" mapstructure:"encoding"`
		Level             string `yaml:"level" mapstructure:"level"`
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
	// TODO use envconfig https://github.com/kelseyhightower/envconfig

	cfg.Mysql.Address = os.Getenv("MYSQL_ADDRESS")
	cfg.Mysql.Username = os.Getenv("MYSQL_USER")
	cfg.Mysql.Password = os.Getenv("MYSQL_PASS")
	cfg.Mysql.Database = os.Getenv("MYSQL_DATABASE")
	cfg.Mysql.Protocol = os.Getenv("MYSQL_PROTOCOL")

	cfg.Redis.Addr = os.Getenv("REDIS_ADDRESS")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		return err
	}
	cfg.Redis.DB = redisDB
	cfg.Redis.Password = os.Getenv("REDIS_PASS")
	redisPoolSize, err1 := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
	if err1 != nil {
		return err1
	}
	cfg.Redis.PoolSize = redisPoolSize

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")

	cfg.Environment = os.Getenv("APP_ENV")

	return nil
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
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
