package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/yopi416/mind-kanban-backend/internal/session"
)

type ctxKey string

const ctxKeyUserID ctxKey = "userID"

// ミドルウェア内では使用しないが、他関数から呼び出すユーティリティ関数
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	v := ctx.Value(ctxKeyUserID)
	id, ok := v.(int64)
	if !ok || id == 0 {
		return 0, false
	}
	return id, true
}

type RequireLoginOptions struct {
	SessionManager   *session.SessionManager                      // 既存の SessionManager を直接利用
	SkipPaths        []string                                     // ログイン検証を行わないパス
	RequireCSRFToken bool                                         // CSRF検証を行うかどうか
	OnUnauthorized   func(w http.ResponseWriter, r *http.Request) // ログイン検証失敗時の処理
}

func RequireLogin(next http.Handler, opt RequireLoginOptions) http.Handler {
	if opt.OnUnauthorized == nil {
		opt.OnUnauthorized = func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		}
	}

	requireLoginHandler := func(w http.ResponseWriter, r *http.Request) {
		lg := slog.Default().With("middleware", "RequireLogin", "path", r.URL.Path)

		// OPTIONSは素通し
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// ログイン検証除外パスも素通し
		for _, skipPath := range opt.SkipPaths {
			// lg.Info("path debug", "r.URL.Path", r.URL.Path, "check ok", strings.HasPrefix(r.URL.Path, skipPath))

			if strings.HasPrefix(r.URL.Path, skipPath) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Cookieからsession_idを取得
		sessCookie, err := r.Cookie("session_id")

		if err != nil || sessCookie == nil || sessCookie.Value == "" {
			lg.Warn("no session cookie")
			opt.OnUnauthorized(w, r)
			return
		}

		sessID := sessCookie.Value

		// GetSessionにて検証 & UserIDを取得
		userID, ok := opt.SessionManager.GetSession(sessID)
		if !ok || userID == 0 {
			lg.Warn("invalid or expired session")
			opt.OnUnauthorized(w, r)
			return
		}

		// CSRFトークンの検証
		// - Get, Head, OPTIONSは検証しない
		// - ダブルサブミット方式で検証
		if opt.RequireCSRFToken && r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
			csrfHeader := r.Header.Get("X-CSRF-Token")

			csrfCookie, err := r.Cookie("csrf_token")
			csrfToken := ""
			if err == nil {
				csrfToken = csrfCookie.Value
			}

			if csrfHeader == "" || csrfToken == "" || csrfHeader != csrfToken {
				http.Error(w, "csrf invalid", http.StatusForbidden)
				lg.Warn("csrf check failed")
				return
			}

		}

		// contextに取得したUserIDを保存し、次の処理を実行
		ctx := context.WithValue(r.Context(), ctxKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	}

	return http.HandlerFunc(requireLoginHandler)

}
