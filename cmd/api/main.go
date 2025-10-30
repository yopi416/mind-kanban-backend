package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yopi416/mind-kanban-backend/api"
	"github.com/yopi416/mind-kanban-backend/configs"
	"github.com/yopi416/mind-kanban-backend/internal/handler"
	"github.com/yopi416/mind-kanban-backend/internal/middleware"
)

func newLogger(cfg *configs.ConfigList) {
	isDev := cfg.IsDevelopment()

	var h slog.Handler
	if isDev {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {

		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	lg := slog.New(h).With("service", "minkan-api") // サービス名は固定で付与
	slog.SetDefault(lg)                             // デフォルトロガーをカスタムロガーに設定
}

func main() {
	err := realMain()
	if err != nil {
		slog.Error("main exit with error", "err", err)
		os.Exit(1)
	}
}

func realMain() error {

	// 環境変数の取得
	cfg, err := configs.LoadEnv()
	if err != nil {
		return err
	}

	port := cfg.APIPort

	// dbPath := os.Getenv("DB_PATH")
	// if dbPath == "" {
	// 	dbPath = defaultDBPath
	// }

	// ログ設定
	newLogger(cfg)

	// 確認用ログ出力
	slog.Debug("Environment loaded", "env", cfg.Env)
	slog.Debug("DB config",
		"user", cfg.DBUser,
		"host", cfg.DBHost,
		"port", cfg.DBPort,
		"name", cfg.DBName,
	)

	// set time zone
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	// todoDB, err := db.NewDB(dbPath)
	// if err != nil {
	// 	return err
	// }
	// defer todoDB.Close()

	// api.ServerInterface を取得(http.Serverではないので注意)
	s, err := handler.NewServer(cfg)
	if err != nil {
		return err
	}

	mux := api.HandlerWithOptions(s, api.StdHTTPServerOptions{
		BaseURL: "/v1",
	})

	// ミドルウェア適用
	handlerWithMW := middleware.RequireLogin(mux, middleware.RequireLoginOptions{
		SessionManager:   s.SessionManager,
		SkipPaths:        []string{"/v1/healthz", "/v1/auth/"},
		RequireCSRFToken: true,
		OnUnauthorized:   nil, // デフォルトを利用
	})

	handlerWithMW = middleware.AccessLog(handlerWithMW)

	// handlerWithMW := middleware.AccessLog(
	// 	middleware.CORS(
	// 		middleware.OIDC(mux),
	// 	),
	// )

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handlerWithMW,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// graceful shutdown処理
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	serverErrCh := make(chan error, 1)

	go func() {
		slog.Info("HTTP server starting on", "addr", server.Addr)
		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
			return
		}

		serverErrCh <- nil
	}()

	// どちらが先でも拾えるようにする（←ここがポイント）
	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = server.Shutdown(shutdownCtx) // 受け付けたリクエストが終わるのを5秒間だけ待ってShutdown
		if err != nil {
			return err
		}

		// go routineリーク防止
		err = <-serverErrCh
		if err != nil {
			return err
		}

		slog.Info("server shut down cleanly")
		return nil

	case err := <-serverErrCh:
		if err != nil {
			return err
		}
		slog.Info("server closed")
		return nil
	}

}
