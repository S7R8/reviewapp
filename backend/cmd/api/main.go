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

	"github.com/s7r8/reviewapp/internal/di"
	"github.com/s7r8/reviewapp/internal/infrastructure/config"
	"github.com/s7r8/reviewapp/internal/infrastructure/persistence/postgres"
)

func main() {
	// 1. 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("🚀 Starting ReviewApp API (env: %s)\n", cfg.Env)
	
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

	// 3. Wire で依存関係を自動解決
	knowledgeHandler, err := di.InitializeKnowledgeHandler(db.DB)
	if err != nil {
		log.Fatalf("Failed to initialize knowledge handler: %v", err)
	}

	reviewHandler, err := di.InitializeReviewHandler(db.DB, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize review handler: %v", err)
	}

	// 4. Echoサーバー初期化
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ヘルスチェック
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

	// ルート
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "ReviewApp API",
			"version": "0.1.0",
			"env":     cfg.Env,
		})
	})

	// API v1 ルーティング
	api := e.Group("/api/v1")

	// ナレッジエンドポイント
	api.POST("/knowledge", knowledgeHandler.CreateKnowledge) // KN-001: ナレッジ作成
	api.GET("/knowledge", knowledgeHandler.ListKnowledge)    // KN-002: ナレッジ一覧取得
	api.POST("/review", reviewHandler.ReviewCode)            // RV-001: コードレビュー実行

	// 5. サーバー起動（グレースフルシャットダウン対応）
	go func() {
		addr := "0.0.0.0:" + cfg.Server.Port
		fmt.Printf("📍 Server listening on %s\n", addr)
		fmt.Printf("💡 API Endpoint: http://localhost:%s/api/v1\n", cfg.Server.Port)
		fmt.Printf("🏥 Health Check: http://localhost:%s/health\n", cfg.Server.Port)

		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server: ", err)
		}
	}()

	// 6. グレースフルシャットダウン
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
