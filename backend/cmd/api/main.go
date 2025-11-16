package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/s7r8/reviewapp/internal/application/usecase/user"
	"github.com/s7r8/reviewapp/internal/di"
	"github.com/s7r8/reviewapp/internal/infrastructure/auth"
	"github.com/s7r8/reviewapp/internal/infrastructure/config"
	"github.com/s7r8/reviewapp/internal/infrastructure/persistence/postgres"
	"github.com/s7r8/reviewapp/internal/interfaces/http/handler"
	httpmiddleware "github.com/s7r8/reviewapp/internal/interfaces/http/middleware"
)

func main() {
	// 1. 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("🚀 Starting ReviewApp API (env: %s)\n", cfg.Env)

	// Auth0設定の確認
	if cfg.Auth.Domain == "" || cfg.Auth.Audience == "" {
		log.Fatal("❌ AUTH0_DOMAIN and AUTH0_AUDIENCE must be set!")
	}
	fmt.Printf("✅ Auth0 configured (domain: %s)\n", cfg.Auth.Domain)

	// デバッグ: APIキーの確認
	if cfg.LLM.ClaudeAPIKey == "" {
		log.Println("⚠️  WARNING: CLAUDE_API_KEY is not set!")
	} else {
		log.Printf("✅ Claude API Key loaded (length: %d)\n", len(cfg.LLM.ClaudeAPIKey))
	}

	// 2. データベース接続
	db, err := postgres.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 3. Auth0認証の初期化
	ctx := context.Background()

	// JWKSキャッシュの作成と起動
	jwksCache := auth.NewJWKSCache(cfg.Auth.Domain, 15*time.Minute)
	if err := jwksCache.Start(ctx); err != nil {
		log.Fatalf("Failed to start JWKS cache: %v", err)
	}
	fmt.Println("✅ JWKS cache initialized")

	// JWT Validatorの作成
	validator := auth.NewValidator(jwksCache, cfg.Auth.Domain, cfg.Auth.Audience)
	fmt.Println("✅ JWT validator initialized")

	// 4. ユーザー関連の初期化
	userRepo := postgres.NewUserRepository(db.DB)
	userUC := user.NewUseCase(userRepo)
	authHandler := handler.NewAuthHandler(userUC)
	fmt.Println("✅ User usecase initialized")

	// 認証ミドルウェアの作成（UserRepositoryを渡す）
	authMiddleware := httpmiddleware.NewAuthMiddleware(validator, userRepo)
	fmt.Println("✅ Auth middleware initialized")

	// 5. Wire で依存関係を自動解決
	knowledgeHandler, err := di.InitializeKnowledgeHandler(db.DB)
	if err != nil {
		log.Fatalf("Failed to initialize knowledge handler: %v", err)
	}

	reviewHandler, err := di.InitializeReviewHandler(db.DB, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize review handler: %v", err)
	}

	// 6. Echoサーバー初期化
	e := echo.New()

	// グローバルミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS設定（開発環境用）
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// ヘルスチェック（認証不要）
	e.GET("/health", func(c echo.Context) error {
		// DB接続確認
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := db.HealthCheck(ctx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status":   "error",
				"database": "disconnected",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status":   "ok",
			"service":  "ReviewApp API",
			"database": "connected",
		})
	})

	// ルート（認証不要）
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "ReviewApp API",
			"version": "0.1.0",
			"env":     cfg.Env,
		})
	})

	// 7. API v1 ルーティング
	api := e.Group("/api/v1")

	// 公開エンドポイント（認証不要）
	public := api.Group("/public")
	public.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// 保護されたエンドポイント（認証必須）
	protected := api.Group("")
	// OPTIONSリクエスト（CORS Preflight）は認証をスキップ
	protected.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Printf("🔵 Request: %s %s\n", c.Request().Method, c.Request().URL.Path)
			if c.Request().Method == http.MethodOptions {
				fmt.Println("✅ Skipping auth for OPTIONS")
				return next(c)
			}
			fmt.Println("🔐 Authenticating request...")
			return authMiddleware.Authenticate(next)(c)
		}
	})

	// ユーザー同期エンドポイント（初回ログイン時に呼ばれる）
	protected.POST("/auth/sync", authHandler.SyncUser)

	// ナレッジエンドポイント（認証必須）
	protected.POST("/knowledge", knowledgeHandler.CreateKnowledge) // KN-001: ナレッジ作成
	protected.GET("/knowledge", knowledgeHandler.ListKnowledge)    // KN-002: ナレッジ一覧取得

	// レビューエンドポイント（認証必須）
	protected.POST("/reviews", reviewHandler.ReviewCode)                 // RV-001: コードレビュー実行
	protected.GET("/reviews", reviewHandler.ListReviews)                 // RV-002: レビュー履歴一覧取得
	protected.GET("/reviews/:id", reviewHandler.GetReviewByID)          // RV-003: レビュー詳細取得 ★ 追加
	protected.PUT("/reviews/:id/feedback", reviewHandler.UpdateFeedback) // RV-004: フィードバック更新

	// 8. サーバー起動（グレースフルシャットダウン対応）
	go func() {
		addr := "127.0.0.1:" + cfg.Server.Port
		fmt.Printf("📍 Server listening on %s\n", addr)
		fmt.Printf("💡 API Endpoint: http://localhost:%s/api/v1\n", cfg.Server.Port)
		fmt.Printf("🔒 Protected endpoints require Authorization: Bearer <token>\n")
		fmt.Printf("🏥 Health Check: http://localhost:%s/health\n", cfg.Server.Port)
		fmt.Printf("🌐 CORS: Allowing localhost:5173, localhost:3000\n")

		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server: ", err)
		}
	}()

	// 9. グレースフルシャットダウン
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("✅ Server stopped gracefully")
}
