package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nguyentk31/url-shortening/internal/database"
	"github.com/nguyentk31/url-shortening/internal/utils"
)

type Service struct {
	queries *database.Queries
}

func NewService(queries *database.Queries) *Service {
	return &Service{queries: queries}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

func (s *Service) CreateUrl(w http.ResponseWriter, r *http.Request) {
	// parse url
	var req ShortenRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if err := utils.FormatUrl(&req.URL); err != nil {
		log.Println("Invalid URL format:", err)
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	shortCode := utils.ConvertBase10ToBase62(time.Now().UnixMilli())

	// Save to database
	urlEntity, err := s.queries.CreateUrl(
		r.Context(),
		database.CreateUrlParams{
			Url:       req.URL,
			ShortCode: shortCode,
		},
	)

	if err != nil {
		log.Println("Failed to create short URL:", err)
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(urlEntity)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (s *Service) RetrieveUrl(w http.ResponseWriter, r *http.Request) {
	shortenID := chi.URLParam(r, "shortenID")
	if shortenID == "" {
		http.Error(w, "shortenID is required", http.StatusBadRequest)
		return
	}

	urlEntity, err := s.queries.GetUrl(r.Context(), shortenID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "shortenID not found", http.StatusNotFound)
			return
		}
		log.Println("Failed to retrieve URL:", err)
		http.Error(w, "Failed to retrieve URL", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(urlEntity)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (s *Service) UpdateUrl(w http.ResponseWriter, r *http.Request) {
	// parse url
	var req ShortenRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if err := utils.FormatUrl(&req.URL); err != nil {
		log.Println("Invalid URL format:", err)
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	shortenID := chi.URLParam(r, "shortenID")
	if shortenID == "" {
		http.Error(w, "shortenID is required", http.StatusBadRequest)
		return
	}

	urlEntity, err := s.queries.UpdateUrl(
		r.Context(),
		database.UpdateUrlParams{
			Url:       req.URL,
			ShortCode: shortenID,
		},
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "shortenID not found", http.StatusNotFound)
			return
		}
		log.Println("Failed to update URL:", err)
		http.Error(w, "Failed to update URL", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(urlEntity)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (s *Service) DeleteUrl(w http.ResponseWriter, r *http.Request) {
	shortenID := chi.URLParam(r, "shortenID")
	if shortenID == "" {
		http.Error(w, "shortenID is required", http.StatusBadRequest)
		return
	}

	_, err := s.queries.DeleteUrl(r.Context(), shortenID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "shortenID not found", http.StatusNotFound)
			return
		}
		log.Println("Failed to delete URL:", err)
		http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) StatsUrls(w http.ResponseWriter, r *http.Request) {
	shortenID := chi.URLParam(r, "shortenID")
	if shortenID == "" {
		http.Error(w, "shortenID is required", http.StatusBadRequest)
		return
	}

	urlStats, err := s.queries.StatUrls(r.Context(), shortenID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "shortenID not found", http.StatusNotFound)
			return
		}
		log.Println("Failed to retrieve URL stats:", err)
		http.Error(w, "Failed to retrieve URL stats", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(urlStats)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (s *Service) IncrementAccessCount(w http.ResponseWriter, r *http.Request) {
	shortenID := chi.URLParam(r, "shortenID")
	if shortenID == "" {
		http.Error(w, "shortenID is required", http.StatusBadRequest)
		return
	}

	access_count, err := s.queries.IncrementAccessCount(r.Context(), shortenID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "shortenID not found", http.StatusNotFound)
			return
		}
		log.Println("Failed to increment access count:", err)
		http.Error(w, "Failed to increment access count", http.StatusInternalServerError)
		return
	}

	respMap := map[string]int64{"access_count": access_count}
	resp, err := json.Marshal(respMap)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
