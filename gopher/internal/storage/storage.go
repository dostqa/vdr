package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopher/internal/models"
	"os"
	"path/filepath"
)

// Надо хранить:
// Оригинальный аудиофайл
// Оригинальная транскрибация
// Модифицированный аудиофайл
// Модифицированная транскрибация

// Сервисные данные:
// Корень хранилища
// Мапа сохраненных данных

type Storage struct {
	StoragePath string                      `json:"-"`
	AudioFiles  map[string]models.AudioFile `json:"audioFiles"`
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.NewStorage"
	var storage = &Storage{}

	data, err := os.ReadFile(storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = json.Unmarshal(data, storage)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage.StoragePath = storagePath

	return storage, nil
}

func (s *Storage) Close() error {
	const op = "storage.Close"

	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = os.WriteFile(s.StoragePath, data, 0644)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveNewAudioFile(audioFile models.AudioFile) error {
	const op = "storage.SaveNewAudioFile"

	_, exists := s.AudioFiles[audioFile.ID]
	if exists {
		return fmt.Errorf("%s: %w", op, errors.New("duplicate audio id"))
	}

	// 1. Создаём директорию
	dir := filepath.Dir(s.StoragePath) + "/" + audioFile.ID

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// 2. Создаём файл
	err = os.WriteFile(dir+audioFile.Filename, audioFile.Data, 0644)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.AudioFiles[audioFile.ID] = audioFile

	return nil
}
