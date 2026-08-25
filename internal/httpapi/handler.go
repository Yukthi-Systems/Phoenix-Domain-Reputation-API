/*
Copyright (C) 2026 Yukthi Systems Private Limited

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License version 3
as published by the Free Software Foundation.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
version 3 along with this program. If not, see
<https://www.gnu.org/licenses/>.
*/

// Package httpapi exposes the reputation store over HTTP. Handlers only
// translate between HTTP and the reputation package; they contain no
// business logic.
package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/ipfire"
	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/reputation"
)

// Batch request limits, chosen to keep one call cheap and bounded even if
// a client sends a very large or malicious payload.
const (
	maxBatchDomains  = 1000    // maximum domains accepted per batch request
	maxBatchBodySize = 1 << 20 // maximum request body size (1MB)
)

// Handler serves the reputation HTTP API. It holds no state of its own
// beyond references to the store and logger.
type Handler struct {
	store  *reputation.Store
	logger *slog.Logger
}

// NewHandler builds a Handler backed by store.
func NewHandler(store *reputation.Store, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

// domainResponse is the JSON body returned by the single-domain and batch
// reputation endpoints.
type domainResponse struct {
	Domain     string                     `json:"domain"`
	Score      int                        `json:"score"`
	Categories []reputation.CategoryScore `json:"categories"`
}

// toDomainResponse builds the JSON representation of a domain's
// reputation, normalizing a nil Categories slice to an empty one so
// unknown domains serialize as "categories": [] rather than null.
func toDomainResponse(domain string, rep reputation.DomainReputation) domainResponse {
	categories := rep.Categories
	if categories == nil {
		categories = []reputation.CategoryScore{}
	}
	return domainResponse{Domain: domain, Score: rep.Score, Categories: categories}
}

// Health reports whether the process is up.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.logger, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready reports whether the reputation store has ever published a
// snapshot. Consumers should treat "not ready" as "do not send traffic
// yet".
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if !h.store.Ready() {
		writeJSON(w, h.logger, http.StatusServiceUnavailable, map[string]string{"status": "not ready"})
		return
	}
	writeJSON(w, h.logger, http.StatusOK, map[string]string{"status": "ready"})
}

// DomainReputation handles GET /v1/reputation/domain/{domain}.
func (h *Handler) DomainReputation(w http.ResponseWriter, r *http.Request) {
	domain, ok := ipfire.NormalizeDomain(r.PathValue("domain"))
	if !ok {
		writeError(w, h.logger, http.StatusBadRequest, "invalid domain")
		return
	}
	rep := h.store.Get(domain)
	writeJSON(w, h.logger, http.StatusOK, toDomainResponse(domain, rep))
}

// scoreResponse is the minimal JSON body returned by the score-only
// endpoint, for consumers that only need the numeric score.
type scoreResponse struct {
	Domain string `json:"domain"`
	Score  int    `json:"score"`
}

// DomainScore handles GET /v1/reputation/domain/{domain}/score.
func (h *Handler) DomainScore(w http.ResponseWriter, r *http.Request) {
	domain, ok := ipfire.NormalizeDomain(r.PathValue("domain"))
	if !ok {
		writeError(w, h.logger, http.StatusBadRequest, "invalid domain")
		return
	}
	rep := h.store.Get(domain)
	writeJSON(w, h.logger, http.StatusOK, scoreResponse{Domain: domain, Score: rep.Score})
}

// batchRequest is the JSON body accepted by the batch reputation endpoint.
type batchRequest struct {
	Domains []string `json:"domains"`
}

// batchResponse is the JSON body returned by the batch reputation
// endpoint.
type batchResponse struct {
	Results []domainResponse `json:"results"`
}

// BatchDomainReputation handles POST /v1/reputation/domains.
func (h *Handler) BatchDomainReputation(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBatchBodySize)

	var req batchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, h.logger, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Domains) == 0 {
		writeError(w, h.logger, http.StatusBadRequest, "domains must not be empty")
		return
	}
	if len(req.Domains) > maxBatchDomains {
		writeError(w, h.logger, http.StatusBadRequest, "too many domains in one request")
		return
	}

	results := make([]domainResponse, 0, len(req.Domains))
	for _, raw := range req.Domains {
		domain, ok := ipfire.NormalizeDomain(raw)
		if !ok {
			continue
		}
		results = append(results, toDomainResponse(domain, h.store.Get(domain)))
	}

	writeJSON(w, h.logger, http.StatusOK, batchResponse{Results: results})
}
