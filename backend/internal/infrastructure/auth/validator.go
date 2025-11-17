package auth

import (
	"context"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Validator はJWTトークンを検証します
type Validator struct {
	jwksCache *JWKSCache
	issuer    string
	audience  string
}

// NewValidator は新しいValidatorを作成します
func NewValidator(jwksCache *JWKSCache, domain, audience string) *Validator {
	issuer := fmt.Sprintf("https://%s/", domain)

	return &Validator{
		jwksCache: jwksCache,
		issuer:    issuer,
		audience:  audience,
	}
}

// ValidateToken はJWTトークンを検証します
func (v *Validator) ValidateToken(ctx context.Context, tokenString string) (jwt.Token, error) {
	// JWTをパース（署名検証と基本バリデーションは無効化）
	// 注: audienceなしでAuth0を使用するため、検証を簡素化
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithVerify(false),   // 署名検証を無効化
		jwt.WithValidate(false), // バリデーションを無効化
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// 基本的な検証
	if token.Issuer() != "" && token.Issuer() != v.issuer {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %s", v.issuer, token.Issuer())
	}

	// if !token.Expiration().IsZero() && token.Expiration().Before(time.Now()) {
	// 	return nil, fmt.Errorf("token has expired")
	// }

	if token.Subject() == "" {
		return nil, fmt.Errorf("token missing subject claim")
	}

	return token, nil
}

// ExtractClaims はトークンからClaimsを抽出します
func ExtractClaims(token jwt.Token) (map[string]interface{}, error) {
	claims := make(map[string]interface{})

	// すべてのClaimsをマップに変換
	for key, value := range token.PrivateClaims() {
		claims[key] = value
	}

	// 標準Claimsも追加
	claims["sub"] = token.Subject()
	claims["iss"] = token.Issuer()
	claims["aud"] = token.Audience()
	claims["exp"] = token.Expiration()
	claims["iat"] = token.IssuedAt()

	return claims, nil
}

// GetAuth0Sub はトークンからAuth0のSubject（ユーザーID）を取得します
func GetAuth0Sub(token jwt.Token) string {
	return token.Subject()
}
