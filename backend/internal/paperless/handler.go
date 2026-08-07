package paperless

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/reusing-code/kontor/backend/internal/cryptoutil"
	"github.com/reusing-code/kontor/backend/internal/httputil"
	"github.com/reusing-code/kontor/backend/internal/middleware"
	"github.com/reusing-code/kontor/backend/internal/storage"
)

type Handler struct {
	store         *Store
	logger        *slog.Logger
	encryptionKey []byte
}

func NewHandler(store *Store, logger *slog.Logger, encryptionKeyRaw string) *Handler {
	key, err := cryptoutil.NormalizeEncryptionKey(encryptionKeyRaw)
	if err != nil {
		logger.Warn("paperless integration requires EMAIL_ENCRYPTION_KEY", "error", err)
		key = nil
	}
	return &Handler{store: store, logger: logger, encryptionKey: key}
}

type configResponse struct {
	Configured bool   `json:"configured"`
	BaseURL    string `json:"baseUrl,omitempty"`
}

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.store.GetConfig(r.Context(), middleware.GetUserID(r.Context()))
	if errors.Is(err, storage.ErrNotFound) {
		httputil.WriteJSON(h.logger, w, http.StatusOK, configResponse{Configured: false})
		return
	}
	if err != nil {
		httputil.StoreError(h.logger, w, err)
		return
	}
	httputil.WriteJSON(h.logger, w, http.StatusOK, configResponse{Configured: true, BaseURL: cfg.BaseURL})
}

func (h *Handler) PutConfig(w http.ResponseWriter, r *http.Request) {
	if len(h.encryptionKey) == 0 {
		httputil.Error(h.logger, w, http.StatusInternalServerError, "EMAIL_ENCRYPTION_KEY is not configured")
		return
	}
	var input ConfigInput
	if err := httputil.ReadJSON(r, &input); err != nil {
		httputil.Error(h.logger, w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := input.Validate(); err != nil {
		httputil.Error(h.logger, w, http.StatusBadRequest, err.Error())
		return
	}
	userID := middleware.GetUserID(r.Context())
	cfg, err := h.store.GetConfig(r.Context(), userID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		httputil.StoreError(h.logger, w, err)
		return
	}
	existed := err == nil
	if !existed && input.Token == "" {
		httputil.Error(h.logger, w, http.StatusBadRequest, "token is required")
		return
	}
	if cfg.BaseURL != input.BaseURL {
		// The cached custom field ID belongs to the old instance.
		cfg.CustomFieldID = 0
	}
	cfg.BaseURL = input.BaseURL
	if input.Token != "" {
		encrypted, err := cryptoutil.EncryptString(input.Token, h.encryptionKey)
		if err != nil {
			httputil.Error(h.logger, w, http.StatusInternalServerError, err.Error())
			return
		}
		cfg.EncryptedToken = encrypted
	}
	cfg.UpdatedAt = time.Now().UTC()
	if err := h.store.PutConfig(r.Context(), userID, cfg); err != nil {
		httputil.StoreError(h.logger, w, err)
		return
	}
	httputil.WriteJSON(h.logger, w, http.StatusOK, configResponse{Configured: true, BaseURL: cfg.BaseURL})
}

func (h *Handler) DeleteConfig(w http.ResponseWriter, r *http.Request) {
	if err := h.store.DeleteConfig(r.Context(), middleware.GetUserID(r.Context())); err != nil {
		httputil.StoreError(h.logger, w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// clientForUser builds a client from the stored config. It writes the error
// response itself and returns nil when the integration is unusable.
func (h *Handler) clientForUser(w http.ResponseWriter, r *http.Request) (*Client, Config, bool) {
	userID := middleware.GetUserID(r.Context())
	cfg, err := h.store.GetConfig(r.Context(), userID)
	if errors.Is(err, storage.ErrNotFound) {
		httputil.Error(h.logger, w, http.StatusConflict, "paperless is not configured")
		return nil, Config{}, false
	}
	if err != nil {
		httputil.StoreError(h.logger, w, err)
		return nil, Config{}, false
	}
	if len(h.encryptionKey) == 0 {
		httputil.Error(h.logger, w, http.StatusInternalServerError, "EMAIL_ENCRYPTION_KEY is not configured")
		return nil, Config{}, false
	}
	token, err := cryptoutil.DecryptString(cfg.EncryptedToken, h.encryptionKey)
	if err != nil {
		httputil.Error(h.logger, w, http.StatusInternalServerError, "could not decrypt stored paperless token")
		return nil, Config{}, false
	}
	return NewClient(cfg.BaseURL, token), cfg, true
}

func (h *Handler) TestConfig(w http.ResponseWriter, r *http.Request) {
	client, _, ok := h.clientForUser(w, r)
	if !ok {
		return
	}
	if err := client.Ping(r.Context()); err != nil {
		httputil.Error(h.logger, w, http.StatusBadGateway, err.Error())
		return
	}
	httputil.WriteJSON(h.logger, w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	client, _, ok := h.clientForUser(w, r)
	if !ok {
		return
	}
	page := 1
	if raw := r.URL.Query().Get("page"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			httputil.Error(h.logger, w, http.StatusBadRequest, "invalid page")
			return
		}
		page = parsed
	}
	result, err := client.SearchDocuments(r.Context(), r.URL.Query().Get("query"), page)
	if err != nil {
		httputil.Error(h.logger, w, http.StatusBadGateway, err.Error())
		return
	}
	httputil.WriteJSON(h.logger, w, http.StatusOK, result)
}

func (h *Handler) Thumbnail(w http.ResponseWriter, r *http.Request) {
	documentID, err := strconv.Atoi(r.PathValue("documentId"))
	if err != nil || documentID <= 0 {
		httputil.Error(h.logger, w, http.StatusBadRequest, "invalid documentId")
		return
	}
	client, _, ok := h.clientForUser(w, r)
	if !ok {
		return
	}
	body, contentType, err := client.Thumbnail(r.Context(), documentID)
	if err != nil {
		httputil.Error(h.logger, w, http.StatusBadGateway, err.Error())
		return
	}
	defer body.Close()
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Header().Set("Cache-Control", "private, max-age=3600")
	if _, err := io.Copy(w, body); err != nil {
		h.logger.Warn("streaming paperless thumbnail", "error", err)
	}
}

func (h *Handler) entityParams(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	entityType := r.PathValue("entityType")
	if !ValidEntityType(entityType) {
		httputil.Error(h.logger, w, http.StatusBadRequest, "invalid entityType")
		return "", "", false
	}
	entityID, err := uuid.Parse(r.PathValue("entityId"))
	if err != nil {
		httputil.Error(h.logger, w, http.StatusBadRequest, "invalid entityId")
		return "", "", false
	}
	return entityType, entityID.String(), true
}

func (h *Handler) ListLinks(w http.ResponseWriter, r *http.Request) {
	entityType, entityID, ok := h.entityParams(w, r)
	if !ok {
		return
	}
	links, err := h.store.ListLinks(r.Context(), middleware.GetUserID(r.Context()), entityType, entityID)
	if err != nil {
		httputil.StoreError(h.logger, w, err)
		return
	}
	httputil.WriteJSON(h.logger, w, http.StatusOK, links)
}

type attachResponse struct {
	Links    []Link   `json:"links"`
	Warnings []string `json:"warnings"`
}

func (h *Handler) AttachLinks(w http.ResponseWriter, r *http.Request) {
	entityType, entityID, ok := h.entityParams(w, r)
	if !ok {
		return
	}
	var input AttachInput
	if err := httputil.ReadJSON(r, &input); err != nil {
		httputil.Error(h.logger, w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := input.Validate(); err != nil {
		httputil.Error(h.logger, w, http.StatusBadRequest, err.Error())
		return
	}
	client, cfg, ok := h.clientForUser(w, r)
	if !ok {
		return
	}
	userID := middleware.GetUserID(r.Context())
	now := time.Now().UTC()
	links := make([]Link, 0, len(input.Documents))
	for _, doc := range input.Documents {
		l := Link{
			EntityType: entityType,
			EntityID:   entityID,
			DocumentID: doc.ID,
			Title:      doc.Title,
			EntityURL:  input.EntityURL,
			CreatedAt:  now,
		}
		if err := h.store.PutLink(r.Context(), userID, l); err != nil {
			httputil.StoreError(h.logger, w, err)
			return
		}
		links = append(links, l)
	}

	// Back-links are best effort: the links above are already stored, so
	// paperless being unreachable only produces warnings.
	warnings := []string{}
	fieldID, err := h.ensureCustomField(r, client, cfg)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("could not prepare the %s custom field in paperless: %v", CustomFieldName, err))
	} else {
		for _, l := range links {
			if err := client.SetCustomField(r.Context(), l.DocumentID, fieldID, l.EntityURL); err != nil {
				h.logger.Warn("setting paperless back-link", "documentId", l.DocumentID, "error", err)
				warnings = append(warnings, fmt.Sprintf("could not set the back-link on document %d: %v", l.DocumentID, err))
			}
		}
	}
	httputil.WriteJSON(h.logger, w, http.StatusOK, attachResponse{Links: links, Warnings: warnings})
}

func (h *Handler) ensureCustomField(r *http.Request, client *Client, cfg Config) (int, error) {
	if cfg.CustomFieldID > 0 {
		return cfg.CustomFieldID, nil
	}
	fieldID, err := client.EnsureCustomField(r.Context())
	if err != nil {
		return 0, err
	}
	cfg.CustomFieldID = fieldID
	if err := h.store.PutConfig(r.Context(), middleware.GetUserID(r.Context()), cfg); err != nil {
		h.logger.Warn("caching paperless custom field id", "error", err)
	}
	return fieldID, nil
}

func (h *Handler) DetachLink(w http.ResponseWriter, r *http.Request) {
	entityType, entityID, ok := h.entityParams(w, r)
	if !ok {
		return
	}
	documentID, err := strconv.Atoi(strings.TrimSpace(r.PathValue("documentId")))
	if err != nil || documentID <= 0 {
		httputil.Error(h.logger, w, http.StatusBadRequest, "invalid documentId")
		return
	}
	userID := middleware.GetUserID(r.Context())
	deleted, err := h.store.DeleteLink(r.Context(), userID, entityType, entityID, documentID)
	if err != nil {
		httputil.StoreError(h.logger, w, err)
		return
	}
	// Best-effort back-link cleanup; only clears the field when it still
	// points at this entity.
	cfg, err := h.store.GetConfig(r.Context(), userID)
	if err == nil && cfg.CustomFieldID > 0 && len(h.encryptionKey) > 0 {
		if token, err := cryptoutil.DecryptString(cfg.EncryptedToken, h.encryptionKey); err == nil {
			client := NewClient(cfg.BaseURL, token)
			if err := client.ClearCustomFieldIfMatches(r.Context(), documentID, cfg.CustomFieldID, deleted.EntityURL); err != nil {
				h.logger.Warn("clearing paperless back-link", "documentId", documentID, "error", err)
			}
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
