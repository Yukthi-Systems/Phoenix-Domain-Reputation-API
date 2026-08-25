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

package httpapi

import (
	"log/slog"
	"net/http"
)

// NewRouter builds the HTTP routing table for the reputation API.
//
// /health and /ready are always open, so orchestrator liveness/readiness
// probes never need a key. Every /v1/ route requires a valid X-API-Key
// header if apiKey is non-empty; if apiKey is empty, the /v1/ routes are
// left unauthenticated (intended for local development only — callers
// should log a warning when running this way, see cmd/server).
func NewRouter(h *Handler, apiKey string, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /ready", h.Ready)

	v1 := http.NewServeMux()
	v1.HandleFunc("GET /v1/reputation/domain/{domain}/score", h.DomainScore)
	v1.HandleFunc("GET /v1/reputation/domain/{domain}", h.DomainReputation)
	v1.HandleFunc("POST /v1/reputation/domains", h.BatchDomainReputation)

	var v1Handler http.Handler = v1
	if apiKey != "" {
		v1Handler = RequireAPIKey(apiKey, logger)(v1)
	}
	mux.Handle("/v1/", v1Handler)

	return mux
}
