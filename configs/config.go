package configs

import (
	"os"
	"strconv"
)

func GetEnvDefault(key, defVal string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		return defVal // 環境変数設定がなければdefault値
	}

	return val // 環境変数設定があれば設定値
}

type ConfigList struct {
	Env                 string
	Port                string // HTTPサーバのポート
	DBHost              string
	DBPort              int
	DBDriver            string
	DBName              string
	DBUser              string
	DBPassword          string
	APICorsAllowOrigins []string
}

func (c *ConfigList) IsDevelopment() bool {
	return c.Env == "development"
}

func LoadEnv() (*ConfigList, error) {
	DBPort, err := strconv.Atoi(GetEnvDefault("MYSQL_PORT", "3306"))
	if err != nil {
		return nil, err
	}

	cfg := &ConfigList{
		Env:                 GetEnvDefault("APP_ENV", "development"),
		Port:                GetEnvDefault("APP_PORT", "8080"),
		DBDriver:            GetEnvDefault("DB_DRIVER", "mysql"),
		DBHost:              GetEnvDefault("DB_HOST", "0.0.0.0"),
		DBPort:              DBPort,
		DBUser:              GetEnvDefault("DB_USER", "app"),
		DBPassword:          GetEnvDefault("DB_PASSWORD", "password"),
		DBName:              GetEnvDefault("DB_NAME", "api_database"),
		APICorsAllowOrigins: []string{"http://0.0.0.0:8001"},
	}

	return cfg, nil
}
