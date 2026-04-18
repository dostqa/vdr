package service

import (
	"context"
	"fmt"
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
}

type FileStorage interface {
	SaveFile(context.Context, string, io.Reader, int64) error
}

type SaverService struct {
	MetaDataBase MetaDataBase
	FileStorage  FileStorage
	trm          TransactionManager
}

func NewSaverService(metaDataBase MetaDataBase, fileStorage FileStorage, trm TransactionManager) *SaverService {
	return &SaverService{
		MetaDataBase: metaDataBase,
		FileStorage:  fileStorage,
		trm:          trm,
	}
}

func (s *SaverService) Save(ctx context.Context, file models.File, r io.Reader, size int64) error {
	const op = "service.Save"

	err := s.trm.Run(ctx, func(ctx context.Context) error {
		// 1. Сохраняем нужную информацию
		// 1.1 Сохраняем request
		requestID, err := s.MetaDataBase.SaveRequest(ctx)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

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
		return nil
	})

	return err
}
