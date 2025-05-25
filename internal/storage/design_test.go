package storage

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/arvaliullin/wapa/internal/domain"
)

func TestDesignStorage_SaveDesignFiles(t *testing.T) {
	tempDir := t.TempDir()

	storage := &DesignStorage{DataPath: tempDir}

	design := &domain.Design{
		ID:   "testDesign",
		JS:   "test.js",
		Wasm: "test.wasm",
	}

	jsContent := bytes.NewReader([]byte("console.log('Hello, JS!');"))
	wasmContent := bytes.NewReader([]byte("\x00asm..."))

	err := storage.SaveDesignFiles(design, jsContent, wasmContent)
	if err != nil {
		t.Fatalf("ошибка при сохранении файлов: %v", err)
	}

	jsFilePath := filepath.Join(tempDir, design.ID, design.JS)
	wasmFilePath := filepath.Join(tempDir, design.ID, design.Wasm)

	if _, err := os.Stat(jsFilePath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("JS файл %s не был создан", jsFilePath)
	}

	if _, err := os.Stat(wasmFilePath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("Wasm файл %s не был создан", wasmFilePath)
	}
}

func TestDesignStorage_GetDesignJSFilePath(t *testing.T) {
	tempDir := t.TempDir()
	storage := &DesignStorage{DataPath: tempDir}

	design := &domain.Design{
		ID: "testDesign",
		JS: "test.js",
	}

	expectedPath := filepath.Join(tempDir, design.ID, design.JS)
	err := os.MkdirAll(filepath.Dir(expectedPath), os.ModePerm)
	if err != nil {
		t.Fatalf("ошибка при создании директории: %v", err)
	}

	file, err := os.Create(expectedPath)
	if err != nil {
		t.Fatalf("ошибка при создании JS файла: %v", err)
	}
	file.Close()

	actualPath, err := storage.GetDesignJSFilePath(design)
	if err != nil {
		t.Errorf("ошибка при получении пути к JS файлу: %v", err)
	}

	if actualPath != expectedPath {
		t.Errorf("путь к JS файлу неверный: ожидается %s, получено %s",
			expectedPath,
			actualPath)
	}
}

func TestDesignStorage_GetDesignWasmFilePath(t *testing.T) {
	tempDir := t.TempDir()
	storage := &DesignStorage{DataPath: tempDir}

	design := &domain.Design{
		ID:   "testDesign",
		Wasm: "test.wasm",
	}

	expectedPath := filepath.Join(tempDir, design.ID, design.Wasm)
	err := os.MkdirAll(filepath.Dir(expectedPath), os.ModePerm)
	if err != nil {
		t.Fatalf("ошибка при создании директории: %v", err)
	}

	file, err := os.Create(expectedPath)
	if err != nil {
		t.Fatalf("ошибка при создании Wasm файла: %v", err)
	}
	file.Close()

	actualPath, err := storage.GetDesignWasmFilePath(design)
	if err != nil {
		t.Errorf("ошибка при получении пути к Wasm файлу: %v", err)
	}

	if actualPath != expectedPath {
		t.Errorf("путь к Wasm файлу неверный: ожидается %s, получено %s",
			expectedPath,
			actualPath)
	}
}

func TestDesignStorage_DeleteDesignFiles(t *testing.T) {
	tempDir := t.TempDir()
	storage := &DesignStorage{DataPath: tempDir}

	design := &domain.Design{
		ID:   "testDeleteDesign",
		JS:   "test.js",
		Wasm: "test.wasm",
	}

	dirPath := filepath.Join(tempDir, design.ID)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		t.Fatalf("ошибка при создании директории: %v", err)
	}

	jsFilePath := filepath.Join(dirPath, design.JS)
	wasmFilePath := filepath.Join(dirPath, design.Wasm)
	_, err = os.Create(jsFilePath)
	if err != nil {
		t.Fatalf("ошибка при создании JS файла: %v", err)
	}
	_, err = os.Create(wasmFilePath)
	if err != nil {
		t.Fatalf("ошибка при создании Wasm файла: %v", err)
	}

	err = storage.DeleteDesignFiles(design)
	if err != nil {
		t.Errorf("ошибка при удалении файлов: %v", err)
	}

	if _, err := os.Stat(dirPath); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("директория %s не была удалена", dirPath)
	}
}
