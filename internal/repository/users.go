package repository

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	UserID        int64
	OIDCIss       string
	OIDCSub       string
	DisplayName   string
	Email         string
	EmailVerified bool
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) *UserRepository {
	return &UserRepository{DB: DB}
}

// 新規ユーザーを登録し、生成された user_id を返す
// 同トランザクションにて、minkan_state テーブルの初期化を行うのでtxを引数に
func (ur *UserRepository) CreateUser(ctx context.Context, tx *sql.Tx, user *User) (int64, error) {
	// user_idはAuto Incrementなので未登録でOK
	query := `
		INSERT INTO users (oidc_iss, oidc_sub, display_name, email, email_verified)
		VALUES (?, ?, ?, ?, ?)
	`

	res, err := tx.ExecContext(ctx, query,
		user.OIDCIss, user.OIDCSub, user.DisplayName, user.Email, user.EmailVerified,
	)
	if err != nil {
		return 0, err
	}

	// 主キー（userIDを取得）を取得
	userID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// oidcのiss, subからuserDataを探す
// 見つからない場合、return, nil, nil
func (ur *UserRepository) FindUserByOIDC(ctx context.Context, oidcIss, oidcSub string) (*User, error) {

	query := `
		SELECT user_id, oidc_iss, oidc_sub, display_name, email, email_verified
		FROM users
		WHERE oidc_iss = ? AND oidc_sub = ?
	`

	row := ur.DB.QueryRowContext(ctx, query, oidcIss, oidcSub)
	user := &User{}
	err := row.Scan(
		&user.UserID,
		&user.OIDCIss,
		&user.OIDCSub,
		&user.DisplayName,
		&user.Email,
		&user.EmailVerified,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // 該当ユーザーなし
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

// userIDからuserを探す
// userIDが見つからない場合、return, nil, nil
func (ur *UserRepository) FindUserByUserID(ctx context.Context, userID int64) (*User, error) {

	query := `
		SELECT user_id, oidc_iss, oidc_sub, display_name, email, email_verified
		FROM users
		WHERE user_id = ?
	`

	row := ur.DB.QueryRowContext(ctx, query, userID)
	user := &User{}
	err := row.Scan(
		&user.UserID,
		&user.OIDCIss,
		&user.OIDCSub,
		&user.DisplayName,
		&user.Email,
		&user.EmailVerified,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // 該当ユーザーなし
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

// 最終ログイン日時を現在時刻で更新
func (ur *UserRepository) UpdateLastLoginAt(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET last_login_at = NOW()
		WHERE user_id = ?
	`

	_, err := ur.DB.ExecContext(ctx, query, userID)

	return err
}

// 指定されたユーザーを削除する（ON DELETE CASCADE により関連 minkan_states も削除される）
func (ur *UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM users
		WHERE user_id = ?
	`

	_, err := ur.DB.ExecContext(ctx, query, userID)

	return err
}
