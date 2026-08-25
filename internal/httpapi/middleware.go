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
	"crypto/subtle"
	"log/slog"
	"net/http"
)

// apiKeyHeader is the header clients must send a valid API key in.
const apiKeyHeader = "X-API-Key"

// RequireAPIKey returns middleware that rejects any request whose
// X-API-Key header does not match apiKey with 401 Unauthorized. apiKey
// must be non-empty; callers should only wrap routes with this middleware
// once a key has been configured (see NewRouter).
//
// The comparison uses crypto/subtle.ConstantTimeCompare so that response
// timing does not leak how much of a guessed key was correct.
func RequireAPIKey(apiKey string, logger *slog.Logger) func(http.Handler) http.Handler {
	want := []byte(apiKey)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := []byte(r.Header.Get(apiKeyHeader))

			// ConstantTimeCompare requires equal-length inputs; a length
			// mismatch is itself a safe, non-timing-sensitive signal that
			// the key is wrong.
			valid := len(got) == len(want) && subtle.ConstantTimeCompare(got, want) == 1
			if !valid {
				logger.Warn("rejected request with invalid API key",
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
				)
				writeError(w, logger, http.StatusUnauthorized, "missing or invalid API key")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
