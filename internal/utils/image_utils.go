package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// SaveUploadedImages speichert mehrere multipart.FileHeader Bilder und gibt die relativen Pfade zurück.
func SaveUploadedImages(files []*multipart.FileHeader) ([]string, error) {
	var imagePaths []string

	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(currentFile), "../../")
	imageDir := filepath.Join(projectRoot, "static", "images")

	err := os.MkdirAll(imageDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Erstellen des Bildverzeichnisses: %w", err)
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("Fehler beim Öffnen der Datei %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
		filePath := filepath.Join(imageDir, fileName)

		outFile, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("Fehler beim Speichern der Datei: %w", err)
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, file); err != nil {
			return nil, fmt.Errorf("Fehler beim Kopieren der Datei: %w", err)
		}

		imagePaths = append(imagePaths, fileName)
	}

	return imagePaths, nil
}
