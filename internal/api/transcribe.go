package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/iam43x/interview-help-api/internal/encode"
	"github.com/iam43x/interview-help-api/internal/gpt"
	"github.com/iam43x/interview-help-api/internal/util"
)

type Transcriptor struct {
	encoder   *encode.Encoder
	gptClient *gpt.ChatGptClient
}

type ResponseRecognize struct {
	Text string `json:"text"`
}

func NewTranscriptor(e *encode.Encoder, g *gpt.ChatGptClient) *Transcriptor {
	return &Transcriptor{
		encoder:   e,
		gptClient: g,
	}
}

func (t *Transcriptor) TranscribeHttpHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем файл из формы
	file, handler, err := r.FormFile("file") // "file" - имя поля в форме
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, "Ошибка при получении файла")
		return
	}
	defer file.Close()

	split := strings.Split(handler.Filename, ".")
	ft := split[len(split)-1]

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		util.SendErrorResponse(w, http.StatusInternalServerError, "Ошибка при чтении файла")
		return
	}
	var audio *util.WriteSeeker
	switch ft {
	case "wav":
		var audioData []int
		for _, v := range fileBytes {
			audioData = append(audioData, int(v))
		}
		audio, err = t.encoder.RawToWav(audioData)
		if err != nil {
			util.SendErrorResponse(w, http.StatusInternalServerError, "Ошибка при конвертации файла")
			return
		}
	case "webm":
		audio = &util.WriteSeeker{Filename: "sample.webm"}
		audio.Write(fileBytes)
	default:
		util.SendErrorResponse(w, http.StatusBadRequest, "Неверный формат файла")
	}

	if strings.HasSuffix(handler.Filename, ".wav") {
		util.SendErrorResponse(w, http.StatusBadRequest, "Неверный формат файла. Ожидается .wav")
		return
	}

	text, err := t.gptClient.TranscribeAudio(r.Context(), audio)
	if err != nil {
		util.SendErrorResponse(w, http.StatusInternalServerError, "Ошибка при транскрибации аудио")
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(ResponseRecognize{Text: text})
}
