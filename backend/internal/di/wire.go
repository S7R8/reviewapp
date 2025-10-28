//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/s7r8/reviewapp/internal/application/usecase/knowledge"
	"github.com/s7r8/reviewapp/internal/application/usecase/review"
	"github.com/s7r8/reviewapp/internal/domain/repository"
	"github.com/s7r8/reviewapp/internal/domain/service"
	"github.com/s7r8/reviewapp/internal/infrastructure/config"
	"github.com/s7r8/reviewapp/internal/infrastructure/external"
	"github.com/s7r8/reviewapp/internal/infrastructure/persistence/postgres"
	"github.com/s7r8/reviewapp/internal/interfaces/http/handler"
)

// InitializeKnowledgeHandler - KnowledgeHandlerを初期化（Wireが自動生成）
func InitializeKnowledgeHandler(db *sql.DB) (*handler.KnowledgeHandler, error) {
	wire.Build(
		// Repository
		postgres.NewKnowledgeRepository,
		wire.Bind(new(repository.KnowledgeRepository), new(*postgres.KnowledgeRepository)),

		// UseCase
		knowledge.NewCreateKnowledgeUseCase,
		knowledge.NewListKnowledgeUseCase,

		// Handler
		handler.NewKnowledgeHandler,
	)
	return nil, nil
}

// InitializeReviewHandler - ReviewHandlerを初期化（Wireが自動生成）
func InitializeReviewHandler(db *sql.DB, cfg *config.Config) (*handler.ReviewHandler, error) {
	wire.Build(
		// Repository
		postgres.NewKnowledgeRepository,
		wire.Bind(new(repository.KnowledgeRepository), new(*postgres.KnowledgeRepository)),
		postgres.NewReviewRepository,
		wire.Bind(new(repository.ReviewRepository), new(*postgres.ReviewRepository)),

		// Service
		service.NewReviewService,

		// External
		ProvideClaudeClient,
		wire.Bind(new(external.ClaudeClientInterface), new(*external.ClaudeClient)),

		// UseCase
		review.NewReviewCodeUseCase,
		review.NewUpdateFeedbackUseCase,

		// Handler
		handler.NewReviewHandler,
	)
	return nil, nil
}

// ProvideClaudeClient - ClaudeClientのプロバイダ
func ProvideClaudeClient(cfg *config.Config) *external.ClaudeClient {
	return external.NewClaudeClient(
		cfg.LLM.ClaudeAPIKey,
		cfg.LLM.ClaudeModel,
		cfg.LLM.ClaudeMaxTokens,
	)
}
