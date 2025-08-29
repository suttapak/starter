package service

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type (
	ImageFileService interface {
		IsImageFromFileHeader(fh *multipart.FileHeader) (bool, error)
		GetUuidFileNameFromFileHeader(fh *multipart.FileHeader) string
		GetImageStatsFromFileHeader(fh *multipart.FileHeader) (*imageStats, error)
		GetMimeTypeFromFileHeader(fh *multipart.FileHeader) (string, error)
		SaveFileFromFileHeader(fh *multipart.FileHeader, filePath string) error
		DeleteFile(path string) error
	}
	imageFileService struct{}
	imageStats       struct {
		width, height int
		size          float64
		mimeType      string
	}
)

// DeleteImage implements ImageFileService.
func (s *imageFileService) DeleteFile(path string) error {
	// ensure dir
	if err := s.ensureFileDir(filepath.Join(deletedProductImagePath, filepath.Base(path))); err != nil {
		return err
	}
	// not remove ,but rename it to /deleted
	return os.Rename(path, filepath.Join(deletedProductImagePath, filepath.Base(path)))
}

// IsImageFromFileHeader implements ImageFileService.
func (s *imageFileService) IsImageFromFileHeader(fh *multipart.FileHeader) (bool, error) {
	file, err := fh.Open()
	if err != nil {
		return false, err
	}
	defer file.Close()
	// Try to decode the file as an image
	if _, _, err = image.Decode(file); err != nil {
		return false, err
	}

	if _, seekErr := file.Seek(0, 0); seekErr != nil {
		return false, seekErr
	}
	// Decoding succeeded: file is an image
	return true, nil
}

func (s *imageFileService) GetUuidFileNameFromFileHeader(fh *multipart.FileHeader) string {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	return uuid.New().String() + ext
}

func (s *imageFileService) GetImageStatsFromFileHeader(fh *multipart.FileHeader) (*imageStats, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	imgBound := img.Bounds()
	width := imgBound.Dx()
	height := imgBound.Dy()
	size := fh.Size
	mimeType, err := s.GetMimeTypeFromFileHeader(fh)
	if err != nil {
		return nil, err
	}
	return &imageStats{
		width:    width,
		height:   height,
		size:     float64(size),
		mimeType: mimeType,
	}, nil
}

func (s *imageFileService) GetMimeTypeFromFileHeader(fh *multipart.FileHeader) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	if _, err := f.Read(buf); err != nil {
		return "", err
	}
	return http.DetectContentType(buf), nil
}

func (s *imageFileService) ensureFileDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func (s *imageFileService) SaveFileFromFileHeader(fh *multipart.FileHeader, filePath string) error {
	f, err := fh.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	if err := s.ensureFileDir(filePath); err != nil {
		return err
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, f)
	return err
}

func NewImageFileService() ImageFileService {
	return &imageFileService{}
}
