package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/yopi416/mind-kanban-backend/api"
	"github.com/yopi416/mind-kanban-backend/internal/middleware"
	"github.com/yopi416/mind-kanban-backend/internal/repository"
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Go構造体をJSONエンコードして書き込み
	if err := json.NewEncoder(w).Encode(response); err != nil {
		lg.Error("failed to encode MinkanGetRes", "err", err)
	}
}

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
	userID, ok := middleware.GetUserIDFromContext(r.Context())

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		lg.Warn("userID not found in context")
		return
	}

	// リクエストボディから、minkanデータとversionを取得
	defer func() {
		if err := r.Body.Close(); err != nil {
			lg.Error("failed to close request body", "err", err)
		}
	}()

	var reqBody api.MinkanPutReq
	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		lg.Warn("decode error", "err", err)
		return
	}

	// minkanデータとversion + 1をDBに登録
	err = s.MinkanStatesRepository.UpdateStateByUserID(r.Context(), reqBody.Minkan, userID, reqBody.Version)

	if errors.Is(err, repository.ErrOptimisticLock) {
		http.Error(w, "version conflict", http.StatusConflict)
		lg.Warn("optimistic lock error", "err", err)
		return
	}

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		lg.Error("update state error", "err", err)
		return
	}

	// 置き換えた後のversionを返す
	resBody := api.MinkanPutRes{
		Version: reqBody.Version + 1,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resBody); err != nil {
		lg.Error("failed to encode MinkanPutRes", "err", err)
	}

}
