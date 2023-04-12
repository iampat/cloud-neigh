package embeddings

import (
	"bytes"
	"encoding/json"
	"github/iampat/cloudy-neigh/vector"
	"log"
	"net/http"
)

type OpenAIClient struct {
	key      string
	endpoint string
	model    string
}

func NewOpenAIClient(key string) *OpenAIClient {
	return &OpenAIClient{
		key:      key,
		endpoint: "https://api.openai.com/v1/embeddings",
		model:    "text-embedding-ada-002",
	}
}

func (e *OpenAIClient) Get(input []string) ([]*vector.Vector32, error) {
	embd, _, err := e.GetWithUsage(input)
	return embd, err
}
func (e *OpenAIClient) GetWithUsage(input []string) ([]*vector.Vector32, int, error) {

	// Set the request body parameters
	reqBody := struct {
		Input []string `json:"input"`
		Model string   `json:"model"`
	}{
		Input: input,
		Model: e.model,
	}
	j, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest(http.MethodPost, e.endpoint, bytes.NewBuffer(j))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", "Bearer "+e.key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	if res.StatusCode > 299 {
		log.Fatalln("oops!", buf.String())
	}

	resBody := struct {
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
		}
	}{}
	json.NewDecoder(buf).Decode(&resBody)
	if len(resBody.Data) != len(input) {
		log.Fatalf("oops! len(input)=%d  len(response)=%d\n", len(input), len(resBody.Data))
	}
	embeddings := make([]*vector.Vector32, len(resBody.Data))
	for idx, d := range resBody.Data {
		if idx != d.Index {
			log.Fatalf("oops! %d != %d", idx, d.Index)
		}
		embeddings[idx] = &vector.Vector32{
			Values: d.Embedding,
		}

	}
	return embeddings, resBody.Usage.TotalTokens, nil
}
