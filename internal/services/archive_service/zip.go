package archive_service

import (
	"doodocs-archive/internal/models"
	"io"
	"mime/multipart"
)

type ZipServiceInterface interface {
	GetZipInfo(r io.ReaderAt, size int64) (*models.ZipInfo, error)
	CreateZip(files []multipart.File, filenames []string) (io.Reader, error)
}

type ZipService struct{}

func ZipInit() ZipServiceInterface {
	return &ZipService{}
}
