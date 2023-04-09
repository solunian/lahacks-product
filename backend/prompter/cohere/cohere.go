package cohere

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type GenerateInput struct {
	Model             string   `json:"model"`
	K                 int      `json:"k"`
	Prompt            string   `json:"prompt"`
	StopSequences     []string `json:"stop_sequences"`
	ReturnLikelihoods string   `json:"return_likelihoods"`
	MaxTokens         int      `json:"max_tokens"`
	Temperature       float32  `json:"temperature"`
	Stream            bool     `json:"stream"`
}

type Result struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func Generate(token string, input *GenerateInput, results chan Result) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://production.api.cohere.ai/v1/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	for dec.More() {
		var r Result
		if err := dec.Decode(&r); err != nil {
			return err
		}
		results <- r
	}

	close(results)
	return nil
}

type EmbeddingRequest struct {
	Model    string   `json:"model"`
	Texts    []string `json:"texts"`
	Truncate string   `json:"truncate"`
}

type EmbeddingResponse struct {
	ID         string      `json:"id"`
	Texts      []string    `json:"texts"`
	Embeddings [][]float32 `json:"embeddings"`
	Meta       struct {
		APIVersion struct {
			Version string `json:"version"`
		} `json:"api_version"`
	} `json:"meta"`
}

func GetEmbeddings(token string, texts []string) ([][]float32, error) {
	reqBody, err := json.Marshal(EmbeddingRequest{
		Model:    "large",
		Texts:    texts,
		Truncate: "NONE",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://production.api.cohere.ai/v1/embed", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var embeddingResp EmbeddingResponse
	err = json.NewDecoder(resp.Body).Decode(&embeddingResp)
	if err != nil {
		return nil, err
	}

	return embeddingResp.Embeddings, nil
}
