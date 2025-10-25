package model

import "time"

// User - ユーザーエンティティ
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser - ユーザーを生成（ファクトリーメソッド）
func NewUser(email, name string) *User {
	return &User{
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate - ユーザーのバリデーション
func (u *User) Validate() error {
	// ドメインルールのチェック
	// TODO: 実装
	return nil
}
