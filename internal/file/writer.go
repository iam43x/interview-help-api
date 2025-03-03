package file

import (
	"fmt"
	"os"
	"time"

	"github.com/iam43x/interview-help-4u/internal/util"
)

func SaveAudio(r *util.WriteSeeker) error {
	filepath := fmt.Sprintf("out/voice-sample-%v.wav", time.Now().Unix())
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(r.Bytes()); err != nil {
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}
	return nil
}
