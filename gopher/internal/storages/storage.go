package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopher/internal/models"
	"io"
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

const (
	filePerm = 0644
	dirPerm  = 0755
)

var ErrDuplicateAudioID = errors.New("duplicate audio id")

type Storage struct {
	StoragePath string                      `json:"-"`
	AudioFiles  map[string]models.AudioFile `json:"audioFiles"`
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.NewStorage"
	storage := &Storage{}

	data, err := os.ReadFile(storagePath)
	if errors.Is(err, os.ErrNotExist) {
		storage = &Storage{
			StoragePath: storagePath,
			AudioFiles:  make(map[string]models.AudioFile),
		}
		err := storage.save()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return storage, nil

	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = json.Unmarshal(data, storage)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if storage.AudioFiles == nil {
		storage.AudioFiles = make(map[string]models.AudioFile)
	}

	storage.StoragePath = storagePath

	return storage, nil
}

func (s *Storage) save() error {
	const op = "storage.save"

	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = os.WriteFile(s.StoragePath, data, filePerm)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveNewAudioFile(id string, filename string, r io.Reader) error {
	const op = "storage.SaveNewAudioFile"

	_, exists := s.AudioFiles[id]
	if exists {
		return fmt.Errorf("%s: %w", op, ErrDuplicateAudioID)
	}

	// 1. Создаём директорию
	dir := filepath.Join(filepath.Dir(s.StoragePath), id)
	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// 2. Создаём файл
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	// 3. Стримим файл
	if _, err := io.Copy(file, r); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.AudioFiles[id] = models.NewAudioFile(id, filename)
	if err := s.save(); err != nil {
		delete(s.AudioFiles, id)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
