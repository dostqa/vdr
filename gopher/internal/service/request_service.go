package service

import (
	"context"
	"fmt"
	"gopher/internal/infrastructure/kafka/message"
	"io"
	"strconv"
)

type RequestRepository interface {
	Save(ctx context.Context) (int, error)
	IsReady(ctx context.Context, requestID int) (bool, error)
	SaveOutputMessage(ctx context.Context, msg message.OutputMessage) error
	OutputMessage(ctx context.Context, requestID int) (message.OutputMessage, error)
}

type Minio interface {
	SaveAudioFile(ctx context.Context, objectName string, r io.Reader, size int64) error
	AudioFile(ctx context.Context, objectName string) (io.ReadCloser, error)
}

type Kafka interface {
	SendJSON(ctx context.Context, topic string, key string, data message.InputMessage) error
}

type TransactionManager interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}

type RequestService struct {
	requestRepo RequestRepository
	minio       Minio
	kafka       Kafka
	trm         TransactionManager
}

func NewRequestService(requestRepo RequestRepository, minio Minio, kafka Kafka, trm TransactionManager) *RequestService {
	return &RequestService{
		requestRepo: requestRepo,
		kafka:       kafka,
		minio:       minio,
		trm:         trm,
	}
}

func (s *RequestService) Save(ctx context.Context, filename string, r io.Reader, size int64) (int, error) {
	const op = "service.Save"

	var id int
	err := s.trm.Run(ctx, func(ctx context.Context) error {
		// 1. Сохраняем нужную информацию
		// 1.1 Сохраняем request
		requestID, err := s.requestRepo.Save(ctx)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		id = int(requestID)

		// 1.2 Сохраняем файл в S3 хранилище
		filename = strconv.Itoa(requestID) + "_" + filename
		err = s.minio.SaveAudioFile(ctx, filename, r, size)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		msg := message.InputMessage{
			RequestID: requestID,
			FilePath:  filename,
		}

		// 1.3 Отправка event на обработку
		err = s.kafka.SendJSON(ctx, "input_topic", strconv.Itoa(requestID), msg)
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

func (s *RequestService) OutputMessage(ctx context.Context, requestID int) (message.OutputMessage, error) {
	const op = "requestService.OutputMessage"

	msg, err := s.requestRepo.OutputMessage(ctx, requestID)
	if err != nil {
		return message.OutputMessage{}, fmt.Errorf("%s: %w", op, err)
	}
	return msg, nil
}

func (s *RequestService) IsReady(ctx context.Context, requestID int) (bool, error) {
	const op = "requestService.IsReadt"

	status, err := s.requestRepo.IsReady(ctx, requestID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return status, nil
}

func (s *RequestService) SaveOutputMessage(ctx context.Context, msg message.OutputMessage) error {
	const op = "requestService.SaveOutputMessage"

	if err := s.requestRepo.SaveOutputMessage(ctx, msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *RequestService) File(ctx context.Context, filename string) (io.ReadCloser, error) {
	const op = "requestService.File"

	file, err := s.minio.AudioFile(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return file, nil
}
