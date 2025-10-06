package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// レスポンスのステータスコードを記録するためのラッパー
type statusWriter struct {
	http.ResponseWriter     // 既存の ResponseWriter を埋め込み
	status              int // WriteHeader で設定されたステータスコードを保持
}

// WriteHeader が呼ばれたときにステータスを記録する
func (w *statusWriter) WriteHeader(code int) {
	w.status = code                    // ステータスコードを保存
	w.ResponseWriter.WriteHeader(code) // 元の WriteHeader をそのまま呼ぶ
}

// アクセスログ + panic リカバリ（ベーシック）
func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // 遅延計測用

		// statusの初期値は200だが、今後の処理でWriteHeaderが呼ばれれば上書きする
		// - http.Errorや、next.ServeHTTP内で呼ばれる
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		// 最低限のリカバリ（panic→500を返す＆ログ）
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "err", rec)
				http.Error(sw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		// 次のハンドラへ
		next.ServeHTTP(sw, r)

		// 最低限のアクセスログ
		slog.Info("access",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"latency_ms", time.Since(start).Milliseconds(),
		)
	})
}
