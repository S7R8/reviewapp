package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const openAIEmbeddingURL = "https://api.openai.com/v1/embeddings"

// OpenAIClient - OpenAI API クライアント
type OpenAIClient struct {
	apiKey     string
	model      string
	timeout    time.Duration
	httpClient *http.Client
}

// NewOpenAIClient - コンストラクタ
func NewOpenAIClient(apiKey, model string, timeout time.Duration) *OpenAIClient {
	return &OpenAIClient{
		apiKey:  apiKey,
		model:   model,
		timeout: timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// embeddingRequest - OpenAI Embedding APIのリクエスト
type embeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// embeddingBatchRequest - バッチ用リクエスト
type embeddingBatchRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

// embeddingResponse - OpenAI Embedding APIのレスポンス
type embeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// errorResponse - OpenAI APIのエラーレスポンス
type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// GenerateEmbedding - テキストからEmbeddingベクトルを生成
func (c *OpenAIClient) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// リクエストボディ作成
	reqBody := embeddingRequest{
		Input: text,
		Model: c.model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// HTTPリクエスト作成
	req, err := http.NewRequestWithContext(ctx, "POST", openAIEmbeddingURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// APIリクエスト実行
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	// レスポンス読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// エラーチェック
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("OpenAI API error: %s (type: %s)", errResp.Error.Message, errResp.Error.Type)
	}

	// レスポンスパース
	var embResp embeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// データ検証
	if len(embResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	return embResp.Data[0].Embedding, nil
}

// GenerateEmbeddings - 複数テキストから一括でEmbedding生成
func (c *OpenAIClient) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	// リクエストボディ作成
	reqBody := embeddingBatchRequest{
		Input: texts,
		Model: c.model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// HTTPリクエスト作成
	req, err := http.NewRequestWithContext(ctx, "POST", openAIEmbeddingURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// APIリクエスト実行
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	// レスポンス読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// エラーチェック
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("OpenAI API error: %s (type: %s)", errResp.Error.Message, errResp.Error.Type)
	}

	// レスポンスパース
	var embResp embeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// データ検証
	if len(embResp.Data) != len(texts) {
		return nil, fmt.Errorf("expected %d embeddings, got %d", len(texts), len(embResp.Data))
	}

	// 結果を配列に格納（インデックス順にソート）
	embeddings := make([][]float32, len(texts))
	for _, data := range embResp.Data {
		if data.Index >= len(embeddings) {
			return nil, fmt.Errorf("invalid index in response: %d", data.Index)
		}
		embeddings[data.Index] = data.Embedding
	}

	return embeddings, nil
}
