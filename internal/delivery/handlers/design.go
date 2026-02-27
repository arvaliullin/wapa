package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arvaliullin/wapa/internal/delivery"
	"github.com/arvaliullin/wapa/internal/domain"
	"github.com/arvaliullin/wapa/internal/persistence"
	"github.com/arvaliullin/wapa/internal/storage"
	"github.com/labstack/echo/v4"
)

type DesignHandler struct {
	Repo    persistence.DesignRepositoryContract
	Storage storage.DesignStorageContract
}

func RegisterDesignHandler(httpService delivery.HttpService,
	repo persistence.DesignRepositoryContract,
	storage storage.DesignStorageContract) {
	handler := &DesignHandler{Repo: repo, Storage: storage}
	e := httpService.(*delivery.EchoHttpService).Echo
	e.GET("/api/designs", handler.GetAll)
	e.POST("/api/design", handler.Create)
	e.DELETE("/api/design/:id", handler.Delete)
	e.GET("/api/design/:id/files/:type", handler.DownloadFile)
}

// GetAll получает список планов экспериментов
//
//	@Summary		получает список планов экспериментов
//	@Description	Возвращает список планов экспериментов
//	@Tags			Design
//	@Produce		json
//	@Success		200	{array}		domain.Design
//	@Failure		500	{object}	string	"Ошибка сервера"
//	@Router			/api/designs [get]
func (h *DesignHandler) GetAll(c echo.Context) error {
	designs, err := h.Repo.GetAll()

	if err != nil {
		msg := fmt.Sprintf("Ошибка получения планов экспериментов: %v", err)
		return c.JSON(http.StatusInternalServerError, msg)
	}

	return c.JSON(http.StatusOK, designs)
}

// Create создает новый эксперимент
//
//	@Summary		Создать новый эксперимент
//	@Description	Создает новый эксперимент и загружает связанные файлы (JS и/или Wasm)
//	@Tags			Design
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name	formData	string	true	"Имя эксперимента"
//	@Param			lang	formData	string	true	"Язык реализации модуля"
//	@Param			js		formData	file	false	"Файл JavaScript"
//	@Param			wasm	formData	file	false	"Файл WebAssembly"
//	@Param			functions	formData	string	true	"JSON-строка с функциями"
//	@Success		201		{string}	string	"ID созданного эксперимента"
//	@Failure		400		{string}	string	"Неверный запрос"
//	@Failure		500		{string}	string	"Ошибка сервера"
//	@Router			/api/design [post]
func (h *DesignHandler) Create(c echo.Context) error {
	name := c.FormValue("name")
	lang := c.FormValue("lang")
	functions := c.FormValue("functions")

	if name == "" || lang == "" || functions == "" {
		return c.JSON(http.StatusBadRequest,
			"Необходимо указать имя эксперимента, язык и функции")
	}

	var funcs []domain.Function
	if err := json.Unmarshal([]byte(functions), &funcs); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Ошибка парсинга функций: %v", err))
	}

	jsFile, jsHdr, _ := c.Request().FormFile("js")
	wasmFile, wasmHdr, _ := c.Request().FormFile("wasm")

	var jsFileName, wasmFileName string
	if jsHdr != nil {
		jsFileName = jsHdr.Filename
	}
	if wasmHdr != nil {
		wasmFileName = wasmHdr.Filename
	}

	design := domain.Design{
		Name:      name,
		Lang:      lang,
		Functions: funcs,
		JS:        jsFileName,
		Wasm:      wasmFileName,
	}

	designID, err := h.Repo.Create(design)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			fmt.Sprintf("Ошибка создания эксперимента: %v", err))
	}

	design.ID = designID

	if err := h.Storage.SaveDesignFiles(&design, jsFile, wasmFile); err != nil {
		return c.JSON(http.StatusInternalServerError,
			fmt.Sprintf("Ошибка сохранения файлов эксперимента: %v", err))
	}

	return c.JSON(http.StatusCreated, designID)
}

// Delete удаляет эксперимент
//
//	@Summary		Удалить эксперимент
//	@Description	Удаляет эксперимент и связанные файлы по ID
//	@Tags			Design
//	@Param			id	path	string	true	"ID эксперимента"
//	@Success		204
//	@Failure		404	{string}	string	"Эксперимент не найден"
//	@Failure		500	{string}	string	"Ошибка сервера"
//	@Router			/api/design/{id} [delete]
func (h *DesignHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	design, err := h.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, "Эксперимент не найден")
		}
		return c.JSON(http.StatusInternalServerError,
			fmt.Sprintf("Ошибка поиска эксперимента: %v", err))
	}

	if err := h.Repo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError,
			fmt.Sprintf("Ошибка удаления эксперимента: %v", err))
	}

	if err := h.Storage.DeleteDesignFiles(&design); err != nil {
		return c.JSON(http.StatusInternalServerError,
			fmt.Sprintf("Ошибка удаления файлов эксперимента: %v", err))
	}

	return c.NoContent(http.StatusNoContent)
}

// DownloadFile скачивает файл JavaScript или WebAssembly
//
//	@Summary		Скачать файл эксперимента
//	@Description	Скачивает JavaScript или WebAssembly файл эксперимента по ID
//	@Tags			Design
//	@Produce		octet-stream
//	@Param			id	path	string	true	"ID эксперимента"
//	@Param			type	path	string	true	"Тип файла (js или wasm)"
//	@Success		200
//	@Failure		404	{string}	string	"Файл или эксперимент не найден"
//	@Failure		500	{string}	string	"Ошибка сервера"
//	@Router			/api/design/{id}/files/{type} [get]
func (h *DesignHandler) DownloadFile(c echo.Context) error {
	id := c.Param("id")
	fileType := c.Param("type")

	design, err := h.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, "Эксперимент не найден")
		}
		return c.JSON(http.StatusInternalServerError,
			fmt.Sprintf("Ошибка поиска эксперимента: %v", err))
	}

	var filePath string

	switch fileType {
	case "js":
		filePath, err = h.Storage.GetDesignJSFilePath(&design)
	case "wasm":
		filePath, err = h.Storage.GetDesignWasmFilePath(&design)
	default:
		return c.JSON(http.StatusBadRequest,
			"Неверный тип файла (должен быть js или wasm)")
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("Файл не найден: %v", err))
	}

	return c.File(filePath)
}
