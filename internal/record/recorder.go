package record

import (
	"fmt"

	"github.com/iam43x/interview-help-api/internal/config"
	"github.com/iam43x/interview-help-api/internal/util"

	"github.com/gordonklaus/portaudio"
)

type Recorder struct {
	SampleRate      float64
	BitDepth        int
	Channels        int
	OutChannels     int
	FramesPerBuffer int
}

func NewRecorder(conf *config.Config) *Recorder {
	return &Recorder{
		SampleRate:      16_000,
		BitDepth:        16,
		Channels:        1,
		OutChannels:     0,
		FramesPerBuffer: 1024,
	}
}

func (r *Recorder) RecordAudio() ([]int, error) {
	var audioData []int
	fmt.Print("Нажмите Enter для начала записи...")
	fmt.Scanln() // Ожидание нажатия клавиши для начала записи
	stream, err := portaudio.OpenDefaultStream(r.Channels, r.OutChannels, r.SampleRate, r.FramesPerBuffer, func(in []int16) {
		for _, v := range in {
			audioData = append(audioData, int(v))
		}
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия аудиопотока: %v", err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return nil, fmt.Errorf("ошибка запуска потока: %v", err)
	}
	defer stream.Stop()

	fmt.Printf("%s Нажмите Enter для завершения записи...", util.Red("◉ REC"))
	fmt.Scanln() // Ожидание нажатия клавиши для завершения записи
	return audioData, nil
}
