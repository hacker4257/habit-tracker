package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"habit-tracker/internal/model"
	"habit-tracker/internal/service"
	"habit-tracker/pkg/logger"
)

type RecordHandler struct {
	service service.RecordService
}

func NewRecordHandler(svc service.RecordService) *RecordHandler {
	return &RecordHandler{service: svc}
}

func (h *RecordHandler) HandleRecords(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *RecordHandler) HandleRecord(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/records/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *RecordHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stats, err := h.service.GetStats()
	if err != nil {
		logger.Error("Failed to get stats: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get stats")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func (h *RecordHandler) getAll(w http.ResponseWriter, r *http.Request) {
	records, err := h.service.GetAll()
	if err != nil {
		logger.Error("Failed to get records: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get records")
		return
	}

	respondJSON(w, http.StatusOK, records)
}

func (h *RecordHandler) getByID(w http.ResponseWriter, id int64) {
	record, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "record not found")
			return
		}
		logger.Error("Failed to get record: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get record")
		return
	}

	respondJSON(w, http.StatusOK, record)
}

func (h *RecordHandler) create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	record, err := h.service.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			respondError(w, http.StatusBadRequest, "invalid input")
			return
		}
		logger.Error("Failed to create record: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create record")
		return
	}

	respondJSON(w, http.StatusCreated, record)
}

func (h *RecordHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var req model.UpdateRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	record, err := h.service.Update(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "record not found")
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			respondError(w, http.StatusBadRequest, "invalid input")
			return
		}
		logger.Error("Failed to update record: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update record")
		return
	}

	respondJSON(w, http.StatusOK, record)
}

func (h *RecordHandler) delete(w http.ResponseWriter, id int64) {
	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "record not found")
			return
		}
		logger.Error("Failed to delete record: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete record")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.APIResponse{
		Success: false,
		Error:   message,
	})
}
