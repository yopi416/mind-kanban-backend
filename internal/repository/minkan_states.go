package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// MinkanState は minkan_states テーブル1行を表す構造体
type MinkanState struct {
	UserID        int64           `json:"user_id"`
	StateJSON     json.RawMessage `json:"state_json"`
	SchemaVersion int             `json:"schema_version"`
	Version       int64           `json:"version"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type MinkanStatesRepository struct {
	DB *sql.DB
}

func NewMinkanStatesRepository(DB *sql.DB) *MinkanStatesRepository {
	return &MinkanStatesRepository{DB: DB}
}

const (
	// デフォルトノードタイプ（フロントと共通）
	defaultNodeType = "custom"

	// ルートノードID（固定）
	rootNodeID = "root"
)

// 初回ユーザー登録時にデフォルトの state を挿入
// 同トランザクションにて、user テーブルの初期化を行うのでtxを引数に
func (msr *MinkanStatesRepository) InitState(ctx context.Context, tx *sql.Tx, userID int64) error {
	// 空のマインドマップ＋カンバン構造など、初期JSONを定義
	now := time.Now().UTC()
	pjID, err := gonanoid.New()
	if err != nil {
		return err
	}

	rootNode := Node{
		Id:   rootNodeID,
		Type: defaultNodeType,
		Data: NodeData{
			Label:    "input",
			ParentId: nil,
			IsDone:   false,
			Comments: []NodeComment{},
		},
		Position: struct {
			X float32 `json:"x"`
			Y float32 `json:"y"`
		}{X: 0, Y: 0},
	}

	project := Project{
		Id:        pjID,
		Name:      "New Project",
		Nodes:     []Node{rootNode},
		Edges:     []Edge{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	defaultState := Minkan{
		CurrentPjId: pjID,
		Projects:    Projects{pjID: project},
		KanbanIndex: KanbanIndex{pjID: {}},
		KanbanColumns: KanbanColumns{
			Backlog: []KanbanCardRef{},
			Todo:    []KanbanCardRef{},
			Doing:   []KanbanCardRef{},
			Done:    []KanbanCardRef{},
		},
	}

	// 挿入するjsonデータを[]byte化
	stateBytes, err := json.Marshal(defaultState)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO minkan_states (user_id, state_json, schema_version, version)
		VALUES (?, ?, 1, 1)
	`

	_, err = tx.ExecContext(ctx, query, userID, stateBytes)
	return err
}

// userIDからminkan_stateを探す
// 見つからない場合、return, nil, nil
func (msr *MinkanStatesRepository) FindStateByUserID(ctx context.Context, userID int64) (*MinkanState, error) {

	query := `
		SELECT user_id, state_json, schema_version, version, updated_at
		FROM minkan_states
		WHERE user_id = ?
	`

	row := msr.DB.QueryRowContext(ctx, query, userID)
	state := &MinkanState{}
	err := row.Scan(
		&state.UserID,
		&state.StateJSON,
		&state.SchemaVersion,
		&state.Version,
		&state.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // 該当stateなしの場合nilを返す
	}

	if err != nil {
		return nil, err
	}
	return state, nil
}

// jsonデータを受け取り、userIDに対応するminkan_statesを更新
func (msr *MinkanStatesRepository) UpdateStateByUserID(ctx context.Context, newStateJSON []byte, userID int64) error {

	query := `
		UPDATE minkan_states
		SET state_json = ?, version = version + 1
		WHERE user_id = ?
	`

	res, err := msr.DB.ExecContext(ctx, query, newStateJSON, userID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	// 更新行が0行の時
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil

}
