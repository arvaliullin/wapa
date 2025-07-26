package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arvaliullin/wapa/internal/domain"
)

type DesignStorageContract interface {
	SaveDesignFiles(design *domain.Design,
		jsContent io.Reader, wasmContent io.Reader) error
	GetDesignJSFilePath(design *domain.Design) (string, error)
	GetDesignWasmFilePath(design *domain.Design) (string, error)
	DeleteDesignFiles(design *domain.Design) error
}

// DesignStorage отвечает за взаимодействие с файловым хранилищем JavaScript и WebAssembly модулей.
type DesignStorage struct {
	DataPath string
}

// SaveDesignFiles сохраняет файлы JavaScript и WebAssembly.
func (ds *DesignStorage) SaveDesignFiles(design *domain.Design,
	jsContent io.Reader, wasmContent io.Reader) error {
	if design.JS != "" && jsContent != nil {
		_, err := ds.saveFile(design.ID, design.JS, jsContent)
		if err != nil {
			return fmt.Errorf("ошибка при сохранении JS файла: %w", err)
		}
	}

	if design.Wasm != "" && wasmContent != nil {
		_, err := ds.saveFile(design.ID, design.Wasm, wasmContent)
		if err != nil {
			return fmt.Errorf("ошибка при сохранении Wasm файла: %w", err)
		}
	}

	return nil
}

// GetDesignJSFilePath возвращает полный путь к JavaScript файлу.
func (ds *DesignStorage) GetDesignJSFilePath(design *domain.Design) (string, error) {
	if design.JS == "" {
		return "", errors.New("указанный объект Design не содержит имени JS файла")
	}
	return ds.getFilePath(design.ID, design.JS)
}

// GetDesignWasmFilePath возвращает полный путь к WebAssembly файлу.
func (ds *DesignStorage) GetDesignWasmFilePath(design *domain.Design) (string, error) {
	if design.Wasm == "" {
		return "", errors.New("указанный объект Design не содержит имени Wasm файла")
	}
	return ds.getFilePath(design.ID, design.Wasm)
}

// DeleteDesignFiles удаляет все файлы и директорию, связанную с объектом Design.
func (ds *DesignStorage) DeleteDesignFiles(design *domain.Design) error {
	return ds.DeleteFiles(design.ID)
}

// saveFile сохраняет файл в рамках директории Design.
func (ds *DesignStorage) saveFile(designID string,
	fileName string, fileContent io.Reader) (string, error) {
	dirPath := filepath.Join(ds.DataPath, designID)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("ошибка при создании директории %s: %w", dirPath, err)
	}

	filePath := filepath.Join(dirPath, fileName)
	destFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании файла %s: %w", filePath, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, fileContent)
	if err != nil {
		return "", fmt.Errorf("ошибка при записи в файл %s: %w", filePath, err)
	}

	return filePath, nil
}

// getFilePath возвращает полный путь к файлу.
func (ds *DesignStorage) getFilePath(designID string, fileName string) (string, error) {
	filePath := filepath.Join(ds.DataPath, designID, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("файл %s не существует", filePath)
	}
	return filePath, nil
}

// DeleteFiles удаляет все файлы и папку для указанного DesignID.
func (ds *DesignStorage) DeleteFiles(designID string) error {
	dirPath := filepath.Join(ds.DataPath, designID)
	return os.RemoveAll(dirPath)
}
