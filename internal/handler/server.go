package handler

import (
	"database/sql"

	"github.com/yopi416/mind-kanban-backend/configs"
	"github.com/yopi416/mind-kanban-backend/internal/auth"
	"github.com/yopi416/mind-kanban-backend/internal/repository"
	"github.com/yopi416/mind-kanban-backend/internal/session"
)

// Server は api.ServerInterface を実装する
type Server struct {
	OIDC                   *auth.OIDC
	SessionManager         *session.SessionManager
	RedirectURLAfterLogin  string
	RedirectURLAfterLogout string
	UserRepository         *repository.UserRepository
}

func NewServer(cfg *configs.ConfigList, db *sql.DB) (*Server, error) {
	oidc, err := auth.NewOIDCFromEnv(cfg)
	if err != nil {
		return nil, err
	}

	sm := session.NewSessionManager(cfg.SessionTTL)
	userRepo := repository.NewUserRepository(db)

	return &Server{
		OIDC:                   oidc,
		SessionManager:         sm,
		RedirectURLAfterLogin:  cfg.RedirectURLAfterLogin,
		RedirectURLAfterLogout: cfg.RedirectURLAfterLogout,
		UserRepository:         userRepo,
	}, nil
}
