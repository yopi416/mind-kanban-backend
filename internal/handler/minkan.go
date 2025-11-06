package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/yopi416/mind-kanban-backend/api"
	"github.com/yopi416/mind-kanban-backend/internal/middleware"
)

// あるユーザーのminkan_statesのjsonとversion(楽観ロック用）を取得しレスポンス
func (s *Server) GetMinkan(w http.ResponseWriter, r *http.Request) {
	lg := slog.Default().With("handler", "GetMinkan")

	// 念のための nil ガード
	if s.MinkanStatesRepository == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		lg.Error("missing dependency",
			"hasMinkanRepository", s.MinkanStatesRepository != nil,
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

	// UserIDからminkan_state情報を取得
	minkanState, err := s.MinkanStatesRepository.FindStateByUserID(r.Context(), userID)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		lg.Error("find minkan_state error", "err", err)
		return
	}

	// DBデータ用いてをレスポンス用Go構造体を作成
	response := api.MinkanGetRes{
		Minkan:  minkanState.StateJSON, // json.RawMessage
		Version: minkanState.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Go構造体をJSONエンコードして書き込み
	if err := json.NewEncoder(w).Encode(response); err != nil {
		lg.Error("failed to encode MinkanGetRes", "err", err)
	}
}

// func (s *Server) PostMinkan(w http.ResponseWriter, r *http.Request) {
// 	lg := slog.Default().With("handler", "GetMinkan")

// 	http.Error(w, "not implemented", http.StatusNotImplemented)
// }

// あるユーザーのminkan_statesを置き換え
func (s *Server) PutMinkan(w http.ResponseWriter, r *http.Request) {
	lg := slog.Default().With("handler", "PutMinkan")

	// 念のための nil ガード
	if s.MinkanStatesRepository == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		lg.Error("missing dependency",
			"hasMinkanRepository", s.MinkanStatesRepository != nil,
		)
		return
	}

	// ContextからUserIDを取得
	// userID, ok := middleware.GetUserIDFromContext(r.Context())

	// if !ok {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	lg.Warn("userID not found in context")
	// 	return
	// }

	// 楽観ロック判定

	// 置き換えた後のversionは返さないとかも(楽観ロック用)

	http.Error(w, "not implemented", http.StatusNotImplemented)
}
