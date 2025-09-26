package upload

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// UploadFile: Belirtilen dosyayı izin verilen mime tipleri ve max 5MB kontrolüyle unique isimle kaydeder, ismini döner
func UploadFile(file *multipart.FileHeader, destFolder string, allowedMimeTypes []string) (string, error) {
	if file == nil {
		return "", fmt.Errorf("Dosya bulunamadı")
	}
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("Dosya boyutu 5MB'dan büyük olamaz")
	}
	mimeType := file.Header.Get("Content-Type")
	log.Println("file header ---> ", file.Header)
	log.Println("file mime type---> ", mimeType)
	allowed := false
	for _, t := range allowedMimeTypes {
		if strings.ToLower(t) == strings.ToLower(mimeType) {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", fmt.Errorf("Sadece belirtilen mime tipleri yüklenebilir")
	}

	ext := filepath.Ext(file.Filename)
	uniqueName := uuid.New().String() + ext
	dstPath := filepath.Join(destFolder, uniqueName)

	// Hedef klasörü oluştur (yoksa)
	if err := os.MkdirAll(destFolder, 0755); err != nil {
		return "", fmt.Errorf("Klasör oluşturulamadı: %v", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("Dosya açılamadı: %v", err)
	}
	defer src.Close()

	out, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("Dosya kaydedilemedi: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", fmt.Errorf("Dosya kopyalanamadı: %v", err)
	}

	return uniqueName, nil
}
