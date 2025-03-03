package api

import (
	"encoding/json"
	"net/http"

	"github.com/iam43x/interview-help-4u/internal/gpt"
)

type QueryAPI struct {
	gptClient *gpt.ChatGptClient
}

type RequestQuery struct {
	Text string `json:"text"`
}

type ResponseQuery struct {
	Response string `json:"response"`
}

func NewQueryAPI(g *gpt.ChatGptClient) *QueryAPI {
	return &QueryAPI{
		gptClient: g,
	}
}

func (q *QueryAPI) QueryHttpHandler(w http.ResponseWriter, r *http.Request) {
	var req RequestQuery
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка при декодировании запроса", http.StatusBadRequest)
		return
	}

	resp, err := q.gptClient.AskGpt3Dot5Turbo16K(r.Context(), req.Text)
	if err != nil {
		http.Error(w, "Ошибка при запросе к OpenAI", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseQuery{Response: resp})
}