package api

import (
	"encoding/json"
	"net/http"

	"github.com/iam43x/interview-help-api/internal/gpt"
	"github.com/iam43x/interview-help-api/internal/util"
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
		util.SendErrorResponse(w, http.StatusBadRequest, "Ошибка при декодировании запроса")
		return
	}

	resp, err := q.gptClient.AskGpt3Dot5Turbo16K(r.Context(), req.Text)
	if err != nil {
		util.SendErrorResponse(w, http.StatusInternalServerError, "Ошибка при запросе к OpenAI")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseQuery{Response: resp})
}
