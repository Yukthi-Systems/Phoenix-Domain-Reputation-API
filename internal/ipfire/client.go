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

package ipfire

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// maxResponseBytes caps how much of a category list we will read, guarding
// against an unexpectedly large or unbounded response.
const maxResponseBytes = 32 << 20 // 32MB

// Downloader fetches and parses the domains for a single IPFire category.
// It is the only interface the updater depends on, which makes it trivial
// to fake in tests.
type Downloader interface {
	DownloadDomains(ctx context.Context, category Category) ([]string, error)
}

// Client downloads IPFire DBL category lists over HTTP. A single Client
// (and its underlying http.Client) is reused across all categories and
// update cycles.
type Client struct {
	httpClient *http.Client
}

// NewClient builds a Client with the given per-request timeout. The
// returned Client is safe for concurrent use.
func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// DownloadDomains fetches and parses the domain list for category. It
// returns an error for any non-2xx response, transport failure, or a
// response containing no valid domains.
func (c *Client) DownloadDomains(ctx context.Context, category Category) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, category.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request for %s list: %w", category.Name, err)
	}
	req.Header.Set("User-Agent", "domain-reputation-api/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download %s list: %w", category.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("download %s list: unexpected status %s", category.Name, resp.Status)
	}

	limited := io.LimitReader(resp.Body, maxResponseBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("read %s list: %w", category.Name, err)
	}
	if len(body) > maxResponseBytes {
		return nil, fmt.Errorf("download %s list: response exceeds %d byte limit", category.Name, maxResponseBytes)
	}

	domains := ParseDomains(bytes.NewReader(body))
	if len(domains) == 0 {
		return nil, fmt.Errorf("download %s list: no valid domains parsed", category.Name)
	}
	return domains, nil
}
