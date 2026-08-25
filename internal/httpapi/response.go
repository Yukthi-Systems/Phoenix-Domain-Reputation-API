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
	"encoding/json"
	"log/slog"
	"net/http"
)

// writeJSON writes v to w as a JSON body with the given status code. An
// encoding failure is logged rather than returned, since the status line
// and headers have already been written by the time encoding runs.
func writeJSON(w http.ResponseWriter, logger *slog.Logger, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

// errorResponse is the JSON body returned for non-2xx API responses.
type errorResponse struct {
	Error string `json:"error"`
}

// writeError writes a JSON error body with the given status code and a
// client-safe message. Internal error details are never included, so
// handlers must pass a message appropriate for API consumers.
func writeError(w http.ResponseWriter, logger *slog.Logger, status int, message string) {
	writeJSON(w, logger, status, errorResponse{Error: message})
}
