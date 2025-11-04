package handler

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

// 認可リクエスト URL を生成してユーザーを IdP(Google等)にリダイレクト
func (s *Server) GetAuthLogin(w http.ResponseWriter, r *http.Request) {
	lg := slog.Default().With("handler", "GetAuthLogin")

	// 念のための nil ガード
	if s.OIDC == nil || s.OIDC.RP == nil {
		http.Error(w, "OIDC not initialized", http.StatusInternalServerError)
		lg.Error("oidc not initialized")
		return
	}

	// oidcのgenState生成関数を定義
	genState := func() string {
		return uuid.New().String()
	}

	// http.handlerFuncを返すので、それを実行
	// stateの生成や、cookieへの保存、認可リクエストURLの生成、http.Redirect(w, r, authURL, 302) を実行
	rp.AuthURLHandler(
		genState,
		s.OIDC.RP,
		// urlOptions...,
	)(w, r)

	// 正常時は Info ログを出力
	lg.Info("login redirect ok", "remote", r.RemoteAddr)

}

func (s *Server) GetAuthCallback(w http.ResponseWriter, r *http.Request) {
	lg := slog.Default().With("handler", "GetAuthCallback")

	// 念のための nil ガード
	if s.OIDC == nil || s.OIDC.RP == nil || s.SessionManager == nil {
		http.Error(w, "server not initialized", http.StatusInternalServerError)
		lg.Error("missing dependency",
			"hasOIDC", s.OIDC != nil,
			"hasRP", s.OIDC != nil && s.OIDC.RP != nil,
			"hasSession", s.SessionManager != nil,
		)
		return
	}

	// 後に、トークン取得後に実行するコールバック関数
	// インメモリへのセッション登録や、sessionIDをset-Cookie, 302リダイレクトを実施
	callback := func(
		w http.ResponseWriter,
		r *http.Request,
		tokens *oidc.Tokens[*oidc.IDTokenClaims],
		state string,
		_ rp.RelyingParty,
	) {
		claims := tokens.IDTokenClaims
		iss := claims.Issuer
		sub := claims.Subject

		// iss,subからユーザーIDを取得
		// ⇒今は暫定でuserIDを0としておく
		var userID int64 = 10000 // TODO: 実装後に DB から実IDを取得

		// セッション発行、登録
		sessionID := uuid.New().String()
		sessionTTL := s.SessionManager.GetTTL()
		s.SessionManager.CreateSession(sessionID, userID)

		// クッキー付与（本番は Secure: true）
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,  // jsから読み取れないようにする
			Secure:   false, // ← HTTPS 運用時は true に
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(sessionTTL.Seconds()), // ブラウザ閉後もCookieをキープ
		})

		// csrfトークンの発行してCookieに格納
		csrfToken := uuid.New().String()

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Path:     "/",
			HttpOnly: false,                // (JS)が読めるように
			Secure:   false,                // 本番は true（HTTPS）
			SameSite: http.SameSiteLaxMode, // 運用に応じて
			MaxAge:   int(sessionTTL.Seconds()),
		})

		lg.Info("login success", "iss", iss, "sub", sub)

		// フロントエンドへ 302
		redirectURL := s.RedirectURLAfterLogin
		http.Redirect(w, r, redirectURL, http.StatusFound)

	}

	// CodeExchangeHandlerは、受信した認可レスポンスからトークンリクエストを作成
	// 取得したトークンを検証し?, その後トークンなどを引数として、コールバック関数を実行
	// この際にuserInfoEndポイントを叩くこともできるが、今回は実行しない
	rp.CodeExchangeHandler(callback, s.OIDC.RP)(w, r)

}

func (s *Server) PostAuthLogout(w http.ResponseWriter, r *http.Request) {
	lg := slog.Default().With("handler", "PostAuthLogout")

	// CookieからセッションIDを取得
	sessCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "no session", http.StatusUnauthorized)
		lg.Warn("logout request without session cookie")
		return
	}

	sessionID := sessCookie.Value

	// CSRF: ヘッダとCookieの一致を確認（ダブルサブミット）
	// - /v1/authはCSRFトークン検証省略対象のため個別で記載
	csrfCookie, _ := r.Cookie("csrf_token")
	csrfHeader := r.Header.Get("X-CSRF-Token")
	if csrfHeader == "" || csrfCookie == nil || csrfHeader != csrfCookie.Value {
		http.Error(w, "csrf invalid", http.StatusForbidden)
		lg.Warn("csrf check failed")
		return
	}

	// セッションを削除
	if s.SessionManager != nil {
		s.SessionManager.DeleteSession(sessionID)
	}

	// Cookieを失効
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // 本番は true
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		Secure:   false, // 本番は true
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	lg.Info("logout success")
	w.WriteHeader(http.StatusOK)
	// lg.Info("logout success", "session_id", sessionID)

	// リダイレクト
	// redirectURL := s.RedirectURLAfterLogout
	// http.Redirect(w, r, redirectURL, http.StatusFound)

}
