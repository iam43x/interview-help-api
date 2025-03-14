package main

import (
	"context"
	"fmt"
	"log"

	"github.com/iam43x/interview-help-api/internal/config"
	"github.com/iam43x/interview-help-api/internal/encode"
	"github.com/iam43x/interview-help-api/internal/file"
	"github.com/iam43x/interview-help-api/internal/gpt"
	"github.com/iam43x/interview-help-api/internal/record"
	"github.com/iam43x/interview-help-api/internal/util"

	"github.com/gordonklaus/portaudio"
)

var cnf *config.Config

func init() {
	cnf = config.LoadConfig()
}

func main() {
	ctx, cancel := util.GracefulShutdown()
	defer cancel()

	// Инициализация PortAudio
	if err := portaudio.Initialize(); err != nil {
		log.Fatalf("Ошибка инициализации PortAudio: %v", err)
	}
	defer portaudio.Terminate()

	recorder := record.NewRecorder(cnf)
	encoder := encode.NewEncoder(cnf)
	gptClient := gpt.NewChatGptClient(cnf.ApiKey)

	for {
		if err := run(ctx, recorder, encoder, gptClient); err != nil {
			fmt.Println(util.Red(err.Error()))
		}
	}

}

func run(ctx context.Context, recorder *record.Recorder, encoder *encode.Encoder, gptClient *gpt.ChatGptClient) error {
	// Запись аудио
	audioData, err := recorder.RecordAudio()
	if err != nil {
		return fmt.Errorf("ошибка записи аудио: %v", err)
	}

	audioReader, err := encoder.RawToWav(audioData)
	if err != nil {
		return fmt.Errorf("ошибка создания аудио-ридера: %v", err)
	}

	// background saving sample to file
	go file.SaveAudio(audioReader)

	// Отправка аудио в OpenAI API
	transcription, err := gptClient.TranscribeAudio(ctx, audioReader)
	if err != nil {
		return fmt.Errorf("ошибка преобразования речи в текст: %v", err)
	}

	// Вывод результата
	fmt.Println(util.Green(transcription))

	answer, err := gptClient.AskGpt3Dot5Turbo16K(ctx, transcription)
	if err != nil {
		return fmt.Errorf("ошибка запроса к GPT: %v", err)
	}
	fmt.Println(util.Blue(answer))
	return nil
}
