package service

import (
	"context"
	"fmt"
	"gopher/internal/clients"
	"gopher/internal/models"
	"io"
	"strconv"
)

type TransactionManager interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}

type MetaDataBase interface {
	SaveRequest(context.Context) (int64, error)
	SaveFile(context.Context, int64, string) (int64, error)
	HandleRequestJsonb(ctx context.Context, msg clients.OutputMessage) error
	GetRequestJsonb(ctx context.Context, requestID int) (*clients.OutputMessage, error)
}

type FileStorage interface {
	SaveFile(context.Context, string, io.Reader, int64) error
}

type KafkaService interface {
	SendJSON(ctx context.Context, topic string, key string, data interface{}) error
	StartConsume(context.Context, string, string, func(context.Context, clients.OutputMessage))
}

type SaverService struct {
	MetaDataBase MetaDataBase
	FileStorage  FileStorage
	KafkaService KafkaService
	trm          TransactionManager
}

func NewSaverService(metaDataBase MetaDataBase, fileStorage FileStorage, kafkaService KafkaService, trm TransactionManager) *SaverService {
	return &SaverService{
		MetaDataBase: metaDataBase,
		FileStorage:  fileStorage,
		KafkaService: kafkaService,
		trm:          trm,
	}
}

func (s *SaverService) Save(ctx context.Context, file models.File, r io.Reader, size int64) (int, error) {
	const op = "service.Save"
	var id int
	err := s.trm.Run(ctx, func(ctx context.Context) error {
		// 1. Сохраняем нужную информацию
		// 1.1 Сохраняем request
		requestID, err := s.MetaDataBase.SaveRequest(ctx)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		id = int(requestID)

		// 1.2 Сохраняем файл в S3 хранилище
		file.FileName = strconv.FormatInt(requestID, 10) + "_" + file.FileName
		err = s.FileStorage.SaveFile(ctx, file.FileName, r, size)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		// 1.3 Сохраняем мета информацию файла
		_, err = s.MetaDataBase.SaveFile(ctx, requestID, file.FileName)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		msg := clients.InputMessage{
			RequestID: int(requestID),
			FilePath:  file.FileName,
		}

		// 1.4 Отправка event на обработку
		err = s.KafkaService.SendJSON(ctx, "input_topic", strconv.Itoa(int(requestID)), msg)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *SaverService) Consume(ctx context.Context, msg clients.OutputMessage) {
	s.MetaDataBase.HandleRequestJsonb(ctx, msg)
}

func (s *SaverService) GetRequestJson(ctx context.Context, requestID int) (*clients.OutputMessage, error) {
	return s.MetaDataBase.GetRequestJsonb(ctx, requestID)
}
