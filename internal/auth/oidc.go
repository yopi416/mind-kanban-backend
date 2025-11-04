package auth

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net/http"
	"time"

	"github.com/yopi416/mind-kanban-backend/configs"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type OIDC struct {
	RP rp.RelyingParty
}

// Issuer / ClientID / Secret / RedirectURL / CookieKeyなど、OIDCに関わる設定を反映した
// rp.RelyingParty インスタンス を作成し、それをラップした *OIDC 構造体を返す
func NewOIDCFromEnv(cfg *configs.ConfigList) (*OIDC, error) {
	issuer := cfg.OIDCGoogleIssuer
	clientID := cfg.OIDCGoogleClientID
	clientSecret := cfg.OIDCGoogleClientSecret
	enablePKCE := cfg.OIDCGoogleEnablePKCE
	redirectURL := cfg.OIDCRedirectURL
	cookieKeyB64 := cfg.OIDCCookieKey
	cookieKey, err := base64.StdEncoding.DecodeString(cookieKeyB64)

	if err != nil {
		return nil, err
	}

	// state 値や PKCE の code_verifier, セッション情報を暗号化・署名付き Cookieとして保管
	// 一旦 HTTPで開発するためunsecureにするが後程変更
	cookieHandler := httphelper.NewCookieHandler(cookieKey, cookieKey, httphelper.WithUnsecure())

	// HTTPクライアント準備
	httpClient := &http.Client{
		Timeout: time.Minute, // OIDCメタデータ/トークン/ユーザー情報などのHTTP呼び出しに使用
	}

	// ロガーの準備
	logger := slog.Default().With("module", "oidc")

	// RPオプション
	options := []rp.Option{
		rp.WithCookieHandler(cookieHandler),                         // state/nonce/PKCEの保存
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)), // iat許容オフセット
		rp.WithHTTPClient(httpClient),                               // 上で作成したHTTPクライアント
		rp.WithLogger(logger),                                       // ロガーを組み込み
		rp.WithSigningAlgsFromDiscovery(),                           // OPのdiscoveryから署名アルゴリズム取得
		// rp.WithPKCE(cookieHandler),                                  // PKCE対応
	}

	// // clientSecretが空ならpublic client扱い → PKCE必須（セキュリティ向上）
	// if clientSecret == "" {
	// 	options = append(options, rp.WithPKCE(cookieHandler))
	// }

	// PKCEを有効化
	if enablePKCE {
		options = append(options, rp.WithPKCE(cookieHandler))
	}

	// RP(NewRelyingPartyOIDC)を作成
	rpClient, err := rp.NewRelyingPartyOIDC(
		context.Background(),
		issuer,
		clientID,
		clientSecret,
		redirectURL,
		[]string{oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeProfile},
		options...,
	)
	if err != nil {
		return nil, err
	}
	return &OIDC{RP: rpClient}, nil
}
