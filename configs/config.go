package configs

import (
	"os"
	"strconv"
	"time"
)

func GetEnvDefault(key, defVal string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		return defVal // 環境変数設定がなければdefault値
	}

	return val // 環境変数設定があれば設定値
}

type ConfigList struct {
	// バックエンド
	Env              string
	APIPort          string // HTTPサーバのポート
	CorsAllowOrigins string // CORSで許諾するURL（フロントエンド）

	// Open Id Connect
	OIDCGoogleIssuer       string
	OIDCGoogleClientID     string
	OIDCGoogleClientSecret string
	OIDCGoogleEnablePKCE   bool
	OIDCCookieKey          string
	OIDCRedirectURL        string

	// ログイン
	RedirectURLAfterLogin  string
	RedirectURLAfterLogout string
	SessionTTL             time.Duration

	// DB
	DBHost     string
	DBPort     string
	DBDriver   string
	DBName     string
	DBUser     string
	DBPassword string
}

func (c *ConfigList) IsDevelopment() bool {
	return c.Env == "development"
}

func LoadEnv() (*ConfigList, error) {

	// string ⇒ intに変換
	// DBPort, err := strconv.Atoi(GetEnvDefault("MYSQL_PORT", "3306"))
	// if err != nil {
	// 	return nil, err
	// }

	// string ⇒ boolに変換
	oidcGoogleEnablePKCE, err := strconv.ParseBool(GetEnvDefault("OIDC_GOOGLE_ENABLE_PKCE", "true"))
	if err != nil {
		return nil, err
	}

	// string ⇒ time.Durationに変換
	sessionTTL, err := time.ParseDuration(GetEnvDefault("SESSION_TTL", "12h"))
	if err != nil {
		return nil, err
	}

	cfg := &ConfigList{
		// バックエンド
		Env:              GetEnvDefault("APP_ENV", "development"),
		APIPort:          GetEnvDefault("APP_PORT", "8080"),
		CorsAllowOrigins: GetEnvDefault("CORS_ALLOW_ORIGINS", "http://localhost:5173"),

		// Open ID Connect
		OIDCGoogleIssuer:       GetEnvDefault("OIDC_GOOGLE_ISSUER", "https://accounts.google.com"),
		OIDCGoogleClientID:     GetEnvDefault("OIDC_GOOGLE_CLIENT_ID", ""),
		OIDCGoogleClientSecret: GetEnvDefault("OIDC_GOOGLE_CLIENT_SECRET", ""),
		OIDCGoogleEnablePKCE:   oidcGoogleEnablePKCE,
		OIDCCookieKey:          GetEnvDefault("OIDC_COOKIE_KEY", "dummykey"),
		OIDCRedirectURL:        GetEnvDefault("OIDC_REDIRECT_URL", "http://localhost:8080/v1/auth/callback"),

		// login
		RedirectURLAfterLogin:  GetEnvDefault("REDIRECT_URL_AFTER_LOGIN", "http://localhost:5173/app/mindmap"),
		RedirectURLAfterLogout: GetEnvDefault("REDIRECT_URL_AFTER_LOGOUT", "http://localhost:5173/login"),
		SessionTTL:             sessionTTL,

		// DB
		DBDriver:   GetEnvDefault("DB_DRIVER", "mysql"),
		DBHost:     GetEnvDefault("DB_HOST", "0.0.0.0"),
		DBPort:     GetEnvDefault("MYSQL_PORT", "3306"),
		DBName:     GetEnvDefault("DB_NAME", "api_database"),
		DBUser:     GetEnvDefault("DB_USER", "app"),
		DBPassword: GetEnvDefault("DB_PASSWORD", "password"),
	}

	return cfg, nil
}
