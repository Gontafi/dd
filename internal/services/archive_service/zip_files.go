package archive_service

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"errors"
	"io"
	"mime/multipart"
	"sync"
)

func (z *ZipService) CreateZip(files []multipart.File, filenames []string) (io.Reader, error) {
	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)
	defer zipWriter.Close()

	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	type fileData struct {
		FileName string
		Content  []byte
		Err      error
	}

	const maxFileSize = 50 * 1024 * 1024 // 50MB example limit

	fileChan := make(chan fileData, len(files))
	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once

	setError := func(err error) {
		errOnce.Do(func() {
			firstErr = err
		})
	}

	for i, file := range files {
		wg.Add(1)
		go func(file multipart.File, index int) {
			defer wg.Done()
			defer file.Close()

			fileName := filenames[index]

			content, err := io.ReadAll(io.LimitReader(file, maxFileSize))
			if err != nil {
				setError(err)
				return
			}

			select {
			case fileChan <- fileData{FileName: fileName, Content: content}:
			default:
				setError(errors.New("channel full"))
			}
		}(file, i)
	}

	go func() {
		wg.Wait()
		close(fileChan)
	}()

	for file := range fileChan {
		if firstErr != nil {
			return nil, firstErr
		}

		zipFile, err := zipWriter.Create(file.FileName)
		if err != nil {
			return nil, err
		}

		if _, err = zipFile.Write(file.Content); err != nil {
			return nil, err
		}
	}

	if firstErr != nil {
		return nil, firstErr
	}

	return buffer, nil
}
