package encode

import (
	"fmt"

	"github.com/iam43x/interview-help-4u/internal/config"
	"github.com/iam43x/interview-help-4u/internal/util"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type Encoder struct {
	SampleRate int
	BitDepth   int
	Channels   int
}

func NewEncoder(conf *config.Config) *Encoder {
	return &Encoder{
		SampleRate: conf.SampleRate,
		BitDepth:   conf.BitDepth,
		Channels:   conf.Channels,
	}
}

func (e *Encoder) RawToWav(audioData []int) (*util.WriteSeeker, error) {
	buf := &util.WriteSeeker{Filename: "sample.wav"}
	encoder := wav.NewEncoder(buf, e.SampleRate, e.BitDepth, e.Channels, 1)
	defer encoder.Close()

	// Создание audio.IntBuffer
	intDataBuf := &audio.IntBuffer{
		Data: audioData,
		Format: &audio.Format{
			SampleRate:  e.SampleRate,
			NumChannels: e.Channels,
		},
		SourceBitDepth: e.BitDepth,
	}
	// Записываем аудио-данные в WAV
	if err := encoder.Write(intDataBuf); err != nil {
		return nil, fmt.Errorf("ошибка записи аудио-данных: %w", err)
	}
	return buf, nil
}
