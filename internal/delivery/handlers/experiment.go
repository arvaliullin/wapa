package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arvaliullin/wapa/internal/broker"
	"github.com/arvaliullin/wapa/internal/delivery"
	"github.com/arvaliullin/wapa/internal/domain"
	"github.com/arvaliullin/wapa/internal/persistence"
	"github.com/labstack/echo/v4"
)

type ExperimentHandler struct {
	DesignRepo      persistence.DesignRepositoryContract
	DesignPublisher *broker.DesignPublisher
}

func RegisterExperimentHandler(httpService delivery.HttpService,
	repo persistence.DesignRepositoryContract, publisher *broker.DesignPublisher) {

	handler := &ExperimentHandler{
		DesignRepo:      repo,
		DesignPublisher: publisher,
	}

	e := httpService.(*delivery.EchoHttpService).Echo

	e.POST("/api/experiment/:uid/start", handler.StartExperiment)
}

// StartExperiment обрабатывает запуск нового эксперимента.
//
// @Summary      Запуск эксперимента
// @Description  Запуск нового эксперимента по UID и параметрам из тела запроса
// @Tags         Experiment
// @Accept       json
// @Produce      json
// @Param        uid      path      string      true  "Уникальный идентификатор дизайна (Design ID)"
// @Param        payload  body      object      true  "Параметры запуска эксперимента"  (example: {"repeats": 1, "warmup": false})
// @Success      200      {object}  object      "Эксперимент успешно запущен"
// @Failure      400      {object}  object      "Ошибка в запросе из-за отсутствия UID или некорректных данных"
// @Failure      404      {object}  object      "Дизайн не найден по переданному UID"
// @Failure      500      {object}  object      "Внутренняя ошибка сервера"
// @Router       /api/experiment/{uid}/start [post]
func (h *ExperimentHandler) StartExperiment(c echo.Context) error {
	uid := c.Param("uid")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{Error: "uid is required"})
	}

	type RequestPayload struct {
		Repeats int  `json:"repeats"`
		Warmup  bool `json:"warmup"`
	}

	var reqPayload RequestPayload
	if err := c.Bind(&reqPayload); err != nil {
		return c.JSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{Error: "invalid request payload"})
	}

	design, err := h.DesignRepo.GetByID(uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, struct {
				Error string `json:"error"`
			}{Error: "design not found"})
		}
		return c.JSON(http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
	}

	consumerPayload := domain.DesignPayload{
		ID:        design.ID,
		Name:      design.Name,
		Lang:      design.Lang,
		JS:        design.JS,
		Wasm:      design.Wasm,
		Repeats:   reqPayload.Repeats,
		Warmup:    reqPayload.Warmup,
		Functions: design.Functions,
	}

	message, err := json.Marshal(consumerPayload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: "failed to marshal payload"})
	}

	err = h.DesignPublisher.Publish(message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: "failed to publish message to NATS"})
	}

	return c.JSON(http.StatusOK, struct {
		Status string `json:"status"`
		UID    string `json:"uid"`
	}{
		Status: "experiment started",
		UID:    uid,
	})
}
