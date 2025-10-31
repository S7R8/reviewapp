package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

// JWKSCache はJWKSをキャッシュするための構造体
type JWKSCache struct {
	cache    jwk.Set
	jwksURL  string
	interval time.Duration
}

// NewJWKSCache は新しいJWKSCacheを作成します
func NewJWKSCache(domain string, refreshInterval time.Duration) *JWKSCache {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)

	return &JWKSCache{
		jwksURL:  jwksURL,
		interval: refreshInterval,
	}
}

// Start はバックグラウンドでJWKSを定期的に更新します
func (j *JWKSCache) Start(ctx context.Context) error {
	// 初回取得
	if err := j.refresh(ctx); err != nil {
		return fmt.Errorf("initial JWKS fetch failed: %w", err)
	}

	// 定期更新
	go func() {
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := j.refresh(ctx); err != nil {
					// エラーログ（本番ではloggerを使用）
					fmt.Printf("JWKS refresh failed: %v\n", err)
				}
			}
		}
	}()

	return nil
}

// refresh はJWKSを取得して更新します
func (j *JWKSCache) refresh(ctx context.Context) error {
	set, err := jwk.Fetch(ctx, j.jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	j.cache = set
	return nil
}

// Get はキャッシュされているJWKSを返します
func (j *JWKSCache) Get() jwk.Set {
	return j.cache
}
