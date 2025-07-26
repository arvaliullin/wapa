package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/wapa/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Mock repo that implements DesignRepositoryContract.
type MockDesignRepo struct {
	designs map[string]domain.Design
}

func (repo *MockDesignRepo) Create(design domain.Design) (string, error) {
	design.ID = "test-id"
	repo.designs[design.ID] = design
	return design.ID, nil
}

func (repo *MockDesignRepo) Update(design domain.Design) error {
	if _, exists := repo.designs[design.ID]; !exists {
		return errors.New("design not found")
	}
	repo.designs[design.ID] = design
	return nil
}

func (repo *MockDesignRepo) Delete(id string) error {
	if _, exists := repo.designs[id]; !exists {
		return errors.New("design not found")
	}
	delete(repo.designs, id)
	return nil
}

func (repo *MockDesignRepo) GetAll() ([]domain.Design, error) {
	var designs []domain.Design
	for _, design := range repo.designs {
		designs = append(designs, design)
	}
	return designs, nil
}

func (repo *MockDesignRepo) GetByID(id string) (domain.Design, error) {
	design, exists := repo.designs[id]
	if !exists {
		return domain.Design{}, errors.New("design not found")
	}
	return design, nil
}

// Mock storage that implements DesignStorageContract.
type MockDesignStorage struct {
	files map[string]map[string][]byte
}

func (s *MockDesignStorage) SaveDesignFiles(design *domain.Design, jsContent io.Reader, wasmContent io.Reader) error {
	if design.ID == "" {
		return errors.New("design ID is required")
	}
	if design.JS != "" && jsContent != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(jsContent)
		s.saveFile(design.ID, design.JS, buf.Bytes())
	}
	if design.Wasm != "" && wasmContent != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(wasmContent)
		s.saveFile(design.ID, design.Wasm, buf.Bytes())
	}
	return nil
}

func (s *MockDesignStorage) GetDesignJSFilePath(design *domain.Design) (string, error) {
	if design.JS == "" {
		return "", errors.New("no JS file specified")
	}
	return "/mock/" + design.JS, nil
}

func (s *MockDesignStorage) GetDesignWasmFilePath(design *domain.Design) (string, error) {
	if design.Wasm == "" {
		return "", errors.New("no Wasm file specified")
	}
	return "/mock/" + design.Wasm, nil
}

func (s *MockDesignStorage) DeleteDesignFiles(design *domain.Design) error {
	delete(s.files, design.ID)
	return nil
}

func (s *MockDesignStorage) saveFile(designID string, fileName string, content []byte) {
	if _, exists := s.files[designID]; !exists {
		s.files[designID] = make(map[string][]byte)
	}
	s.files[designID][fileName] = content
}

func TestDesignHandler_GetAll(t *testing.T) {
	e := echo.New()
	repo := &MockDesignRepo{
		designs: map[string]domain.Design{
			"test-id": {
				ID:   "test-id",
				Name: "Test Design",
				Lang: "go",
			},
		},
	}
	storage := &MockDesignStorage{}

	handler := &DesignHandler{Repo: repo, Storage: storage}

	req := httptest.NewRequest(http.MethodGet, "/api/designs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var designs []domain.Design
		json.Unmarshal(rec.Body.Bytes(), &designs)

		assert.Len(t, designs, 1)
		assert.Equal(t, "test-id", designs[0].ID)
	}
}

func TestDesignHandler_Delete(t *testing.T) {
	e := echo.New()
	repo := &MockDesignRepo{
		designs: map[string]domain.Design{
			"test-id": {
				ID:   "test-id",
				Name: "Test Design",
				Lang: "go",
			},
		},
	}
	storage := &MockDesignStorage{files: map[string]map[string][]byte{
		"test-id": {
			"test.js":   []byte("console.log('test');"),
			"test.wasm": []byte("\x00asm..."),
		},
	}}

	handler := &DesignHandler{Repo: repo, Storage: storage}

	req := httptest.NewRequest(http.MethodDelete, "/api/design/test-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test-id")

	if assert.NoError(t, handler.Delete(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.NotContains(t, repo.designs, "test-id")
		assert.NotContains(t, storage.files, "test-id")
	}
}
