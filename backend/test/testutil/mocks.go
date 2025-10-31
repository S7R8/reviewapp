package testutil

import (
	"context"
	"errors"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/infrastructure/external"
)

// MockKnowledgeRepository - ナレッジリポジトリのモック
type MockKnowledgeRepository struct {
	knowledges []*model.Knowledge
	err        error
}

func NewMockKnowledgeRepository() *MockKnowledgeRepository {
	return &MockKnowledgeRepository{
		knowledges: make([]*model.Knowledge, 0),
	}
}

func (m *MockKnowledgeRepository) SetError(err error) {
	m.err = err
}

func (m *MockKnowledgeRepository) SetKnowledges(knowledges []*model.Knowledge) {
	m.knowledges = knowledges
}

func (m *MockKnowledgeRepository) Create(ctx context.Context, knowledge *model.Knowledge) error {
	if m.err != nil {
		return m.err
	}
	m.knowledges = append(m.knowledges, knowledge)
	return nil
}

func (m *MockKnowledgeRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Knowledge, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.knowledges, nil
}

func (m *MockKnowledgeRepository) Update(ctx context.Context, knowledge *model.Knowledge) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *MockKnowledgeRepository) FindByID(ctx context.Context, id string) (*model.Knowledge, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, k := range m.knowledges {
		if k.ID == id {
			return k, nil
		}
	}
	return nil, errors.New("knowledge not found")
}

func (m *MockKnowledgeRepository) FindByCategory(ctx context.Context, userID, category string) ([]*model.Knowledge, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*model.Knowledge
	for _, k := range m.knowledges {
		if k.UserID == userID && k.Category == category {
			result = append(result, k)
		}
	}
	return result, nil
}

func (m *MockKnowledgeRepository) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	for i, k := range m.knowledges {
		if k.ID == id {
			m.knowledges = append(m.knowledges[:i], m.knowledges[i+1:]...)
			return nil
		}
	}
	return errors.New("knowledge not found")
}

// MockReviewRepository - レビューリポジトリのモック
type MockReviewRepository struct {
	reviews []*model.Review
	err     error
}

func NewMockReviewRepository() *MockReviewRepository {
	return &MockReviewRepository{
		reviews: make([]*model.Review, 0),
	}
}

func (m *MockReviewRepository) SetError(err error) {
	m.err = err
}

func (m *MockReviewRepository) Create(ctx context.Context, review *model.Review) error {
	if m.err != nil {
		return m.err
	}
	m.reviews = append(m.reviews, review)
	return nil
}

func (m *MockReviewRepository) FindByID(ctx context.Context, id string) (*model.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, r := range m.reviews {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, errors.New("review not found")
}

func (m *MockReviewRepository) FindByUserID(ctx context.Context, userID string, limit int) ([]*model.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*model.Review
	for _, r := range m.reviews {
		if r.UserID == userID {
			result = append(result, r)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *MockReviewRepository) Update(ctx context.Context, review *model.Review) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *MockReviewRepository) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	for i, r := range m.reviews {
		if r.ID == id {
			m.reviews = append(m.reviews[:i], m.reviews[i+1:]...)
			return nil
		}
	}
	return errors.New("review not found")
}

func (m *MockReviewRepository) FindRecentByUserID(ctx context.Context, userID string, limit int) ([]*model.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	// 簡単な実装：最新のlimit件を返す
	return m.FindByUserID(ctx, userID, limit)
}

func (m *MockReviewRepository) UpdateFeedback(ctx context.Context, reviewID string, score int, comment string) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

// MockUserRepository - ユーザーリポジトリのモック
type MockUserRepository struct {
	users map[string]*model.User // key: auth0_user_id
	err   error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*model.User),
	}
}

func (m *MockUserRepository) SetError(err error) {
	m.err = err
}

func (m *MockUserRepository) SetUser(user *model.User) {
	m.users[user.Auth0UserID] = user
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[user.Auth0UserID] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindByAuth0UserID(ctx context.Context, auth0UserID string) (*model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[auth0UserID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[user.Auth0UserID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	for auth0ID, user := range m.users {
		if user.ID == id {
			delete(m.users, auth0ID)
			return nil
		}
	}
	return errors.New("user not found")
}

// ClaudeClientInterface - Claude Clientのインターフェース
type ClaudeClientInterface interface {
	ReviewCode(ctx context.Context, input external.ReviewCodeInput) (*external.ReviewCodeOutput, error)
}

// MockClaudeClient - Claude APIクライアントのモック
type MockClaudeClient struct {
	response *external.ReviewCodeOutput
	err      error
}

func NewMockClaudeClient() *MockClaudeClient {
	return &MockClaudeClient{
		response: &external.ReviewCodeOutput{
			ReviewResult: "Mock review result",
			TokensUsed:   100,
		},
	}
}

func (m *MockClaudeClient) SetResponse(response *external.ReviewCodeOutput) {
	m.response = response
}

func (m *MockClaudeClient) SetError(err error) {
	m.err = err
}

func (m *MockClaudeClient) ReviewCode(ctx context.Context, input external.ReviewCodeInput) (*external.ReviewCodeOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}
