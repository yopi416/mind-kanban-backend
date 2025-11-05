package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/yopi416/mind-kanban-backend/api"
	"github.com/yopi416/mind-kanban-backend/internal/middleware"
)

func (s *Server) GetUsersMe(w http.ResponseWriter, r *http.Request) {
	lg := slog.Default().With("handler", "GetUsersMe")

	// 念のための nil ガード
	if s.UserRepository == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		lg.Error("missing dependency",
			"hasUserRepository", s.UserRepository != nil,
		)
		return
	}

	// ContextからUserIDを取得
	userID, ok := middleware.GetUserIDFromContext(r.Context())

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		lg.Warn("userID not found in context")
		return
	}

	// UserIDからユーザー情報を取得
	userData, err := s.UserRepository.FindUserByUserID(r.Context(), userID)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		lg.Error("find user error", "err", err)
		return
	}

	// ユーザーデータがnilの場合（未登録）、404を返す
	if userData == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		lg.Warn("userData not found", "userID", userID)
		return
	}

	// openapi.ymlで指定したスキーマを基にレスポンスするJSONを作成
	response := api.User{
		DisplayName: &userData.DisplayName,
		Email:       &userData.Email,
		// UserId:
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		lg.Error("failed to encode userData", "err", err)
	}
}

func (s *Server) DeleteUsersMe(w http.ResponseWriter, r *http.Request) {
	lg := slog.Default().With("handler", "DeleteUsersMe")

	// 念のための nil ガード
	if s.UserRepository == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		lg.Error("missing dependency",
			"hasUserRepository", s.UserRepository != nil,
		)
		return
	}

	// ContextからUserIDを取得
	userID, ok := middleware.GetUserIDFromContext(r.Context())

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		lg.Warn("userID not found in context")
		return
	}

	// userIDに該当するuserデータをDBから削除
	err := s.UserRepository.DeleteUser(r.Context(), userID)

	// 削除失敗処理
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		lg.Error("failed to delete user", "err", err)
		return
	}

	// 成功だが返すデータなし(204レスポンス)
	w.WriteHeader(http.StatusNoContent)
}
