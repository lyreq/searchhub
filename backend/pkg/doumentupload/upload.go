package documentupload

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const MaxFileSize = 5 * 1024 * 1024 // 50MB

var allowedExtensions = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
}

var allowedMimeTypes = map[string]bool{
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
}

// SaveOfficeFile uploads and validates office-type files (PDF, Excel, Word)
func SaveOfficeFile(subfolder string, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", nil
	}

	// Dosya boyutu kontrolü
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("Dosya boyutu 5MB'den büyük olamaz")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("Uzantı '%s' desteklenmiyor", ext)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// İlk 512 byte oku → mime kontrolü
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	if !allowedMimeTypes[contentType] {
		return "", fmt.Errorf("MIME türü '%s' desteklenmiyor", contentType)
	}
	_, err = src.Seek(0, 0)
	if err != nil {
		return "", err
	}

	if contentType == "application/pdf" && !strings.HasPrefix(string(buffer), "%PDF") {
		return "", fmt.Errorf("PDF dosyası geçersiz, içeriği bozulmuş olabilir")
	}

	filename := uuid.New().String() + ext
	uploadPath := filepath.Join("uploads", subfolder)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}

	dstPath := filepath.Join(uploadPath, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", err
	}

	return filename, nil
}
