-- 1) データベース
CREATE DATABASE IF NOT EXISTS minkan
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_0900_ai_ci;
USE minkan;

-- 2) users: OIDCとユーザー属性
CREATE TABLE users (
  user_id        BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  oidc_iss       VARCHAR(255) NOT NULL,
  oidc_sub       VARCHAR(255) NOT NULL,
  display_name   VARCHAR(255) NOT NULL,
  email          VARCHAR(320) NULL,
  email_verified TINYINT(1) NOT NULL DEFAULT 0,
  created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  last_login_at  TIMESTAMP NULL,
  UNIQUE KEY uk_users_oidc (oidc_iss, oidc_sub),
  UNIQUE KEY uk_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3) minkan_states: ユーザーごとのマインド＋カンバン状態（JSON）
-- 一旦はJSONで保存するが、今後はオジェクトストレージ保存も検討（パスをDBに保存）
CREATE TABLE minkan_states (
  user_id        BIGINT NOT NULL PRIMARY KEY,
  state_json     JSON   NOT NULL,          -- { currentPjID, projects:[...], kanbanIndex, kanbanColumns }
  schema_version SMALLINT NOT NULL DEFAULT 1,  -- JSONスキーマのバージョン
  version        BIGINT  NOT NULL DEFAULT 1,   -- 楽観ロック用（更新ごとに+1, オーバーフロー懸念ゼロ）
  updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_states_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
