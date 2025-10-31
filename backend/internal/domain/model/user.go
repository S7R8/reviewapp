package model

import "time"

// User - ユーザーエンティティ
type User struct {
	ID           string     `json:"id"`
	Auth0UserID  string     `json:"auth0_user_id"`  // Auth0のユーザーID
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	AvatarURL    *string    `json:"avatar_url"`     // Auth0のpicture
	Preferences  string     `json:"preferences"`    // JSONB
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// NewUser - ユーザーを生成（ファクトリーメソッド）
func NewUser(auth0UserID, email, name string) *User {
	now := time.Now()
	return &User{
		Auth0UserID: auth0UserID,
		Email:       email,
		Name:        name,
		Preferences: "{}",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate - ユーザーのバリデーション
func (u *User) Validate() error {
	// ドメインルールのチェック
	// TODO: 実装
	return nil
}
