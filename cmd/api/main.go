package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yopi416/mind-kanban-backend/internal/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort = ":8080"
		// defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// dbPath := os.Getenv("DB_PATH")
	// if dbPath == "" {
	// 	dbPath = defaultDBPath
	// }

	// set time zone
	var err error
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

	mux := router.NewRouter()

	// graceful shutdown処理
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrCh := make(chan error, 1)

	go func() {
		log.Println("HTTP server starting on", server.Addr)
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
		log.Println("shutdown signal received")
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

		log.Println("server shut down cleanly")
		return nil

	case err := <-serverErrCh:
		if err != nil {
			return err
		}
		log.Println("server closed")
		return nil
	}

}
