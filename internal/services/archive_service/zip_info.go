package archive_service

import (
	"archive/zip"
	"doodocs-archive/internal/models"
	"doodocs-archive/pkg/utils"
	"io"
	"sync"
)

func (z *ZipService) GetZipInfo(r io.ReaderAt, size int64) (*models.ZipInfo, error) {
	zipReader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	var totalSize float64
	var files []models.FileMeta
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(zipReader.File))

	for _, file := range zipReader.File {
	    if file.FileInfo().IsDir() {
	        // Skip directories
	        continue
	    }
	
	    wg.Add(1)
	
	    go func(f *zip.File) {
	        defer wg.Done()
	
	        fileSize := float64(f.UncompressedSize64)
	        mimeType := utils.DetectMimeTypeZip(f)
	
	        mu.Lock()
	        files = append(files, models.FileMeta{
	            FilePath: f.Name,
	            Size:     fileSize,
	            MimeType: mimeType,
	        })
	        totalSize += fileSize
	        mu.Unlock()
	    }(file)
	}


	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return &models.ZipInfo{
		TotalFiles:  float64(len(files)),
		TotalSize:   totalSize,
		ArchiveSize: float64(size),
		Files:       files,
	}, nil
}
