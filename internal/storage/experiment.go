package storage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arvaliullin/wapa/internal/domain"
)

type ExperimentStorage struct {
	DataDirStorage string
}

func NewExperimentStorage(baseDir string) (*ExperimentStorage, error) {
	if baseDir == "" {
		baseDir = os.TempDir() + "/runner_files"
	}
	err := os.MkdirAll(baseDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать директорию: %v", err)
	}
	log.Printf("Директория успешно создана: %s", baseDir)
	return &ExperimentStorage{DataDirStorage: baseDir}, nil
}

func (s *ExperimentStorage) ExperimentDir(id string) (string, error) {
	dir := filepath.Join(s.DataDirStorage, id)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("не удалось создать директорию эксперимента: %v", err)
	}
	return dir, nil
}

// DownloadFile скачивает файл и сохраняет его в локальное хранилище в папке эксперимента
func (s *ExperimentStorage) DownloadFile(design domain.DesignPayload, filetype, apiUrl string) (string, error) {
	id := design.ID

	if id == "" || filetype == "" || apiUrl == "" {
		return "", fmt.Errorf("все параметры (id, type, apiUrl) обязательны для скачивания файла")
	}

	url := fmt.Sprintf("%s/api/design/%s/files/%s", apiUrl, id, filetype)

	expDir, err := s.ExperimentDir(id)
	if err != nil {
		log.Printf("Ошибка при создании каталога эксперимента: %v", err)
		return "", err
	}

	filePath := filepath.Join(expDir, filetype)

	if filetype == "js" {
		filePath = filepath.Join(expDir, design.JS)
	}

	if filetype == "wasm" {
		filePath = filepath.Join(expDir, design.Wasm)
	}

	log.Printf("Начало загрузки файла: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка скачивания файла [%s] для ID %s: %v", filetype, id, err)
		return "", fmt.Errorf("не удалось скачать файл (%s). Проверьте соединение с API", filetype)
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("Файл успешно скачан и сохранён: %s", filePath)
	return filePath, nil
}

// DeleteFile удаляет файл по пути
func (s *ExperimentStorage) DeleteFile(filePath string) error {
	if filePath == "" {
		log.Printf("Ошибка: Путь к файлу обязателен для удаления")
		return nil
	}
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Ошибка при удалении файла %s: %v", filePath, err)
		return err
	}
	log.Printf("Файл успешно удалён: %s", filePath)
	return nil
}

// CleanUp удаляет все файлы и папки из хранилища (рекурсивно)
func (s *ExperimentStorage) CleanUp() error {
	files, err := os.ReadDir(s.DataDirStorage)
	if err != nil {
		log.Printf("Ошибка при очистке временного хранилища: %v", err)
		return fmt.Errorf("не удалось очистить временное хранилище")
	}
	log.Printf("Найдено %d каталог(ов) для удаления", len(files))
	for _, file := range files {
		filePath := filepath.Join(s.DataDirStorage, file.Name())
		if file.IsDir() {
			err := os.RemoveAll(filePath)
			if err != nil {
				log.Printf("Ошибка при удалении каталога %s: %v", filePath, err)
			} else {
				log.Printf("Каталог успешно удалён: %s", filePath)
			}
		} else {
			_ = s.DeleteFile(filePath)
		}
	}
	log.Printf("Очистка хранилища завершена")
	return nil
}
