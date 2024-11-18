package archive_service

import (
	"doodocs-archive/internal/models"
	"io"
	"mime/multipart"
)

type ZipServiceInterface interface {
	GetZipInfo(filename string, r io.ReaderAt, size int64) (*models.ZipInfo, error)
	CreateZip(files []multipart.File, filenames []string) (io.Reader, error)
}

type ZipService struct{}

func ZipInit() ZipServiceInterface {
	return &ZipService{}
}
