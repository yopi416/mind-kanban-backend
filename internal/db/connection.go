package db

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/yopi416/mind-kanban-backend/configs"
)

// データベース接続（初期化・維持・クローズ）

func InitDB(cfg *configs.ConfigList) (*sql.DB, error) {
	lg := slog.Default().With("db", "InitDB")

	// 接続用環境変数読み込み
	mysqlConfig := mysql.NewConfig()

	mysqlConfig.User = cfg.DBUser
	mysqlConfig.Passwd = cfg.DBPassword
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = cfg.DBHost + ":" + cfg.DBPort
	mysqlConfig.DBName = cfg.DBName
	mysqlConfig.ParseTime = true // TIMESTAMP型をtime.Timeに変換
	mysqlConfig.Loc = time.Local // time.Timeへの変換時にローカルのタイムゾーンで解釈
	mysqlConfig.Params = map[string]string{
		"charset": "utf8mb4",
	}

	conn, err := mysql.NewConnector(mysqlConfig)

	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(conn)

	// コネクションプール設定
	// - まずは低めの数値から初めて、チューニング
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(7)
	db.SetConnMaxLifetime(20 * time.Minute)

	// 接続確認（疎通チェック）
	if err := db.Ping(); err != nil {
		return nil, err
	}

	lg.Info("db connected", "dbAddr", mysqlConfig.Addr, "dbName", mysqlConfig.DBName)
	return db, nil

}
