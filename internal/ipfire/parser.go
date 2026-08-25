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
	"bufio"
	"io"
	"regexp"
	"strings"
)

// domainPattern is a pragmatic hostname validator: labels of alphanumerics
// and hyphens (not leading/trailing with a hyphen), separated by dots.
var domainPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)

// NormalizeDomain lowercases, trims whitespace and a trailing dot, and
// validates the result as a plausible hostname. It reports false for empty
// or malformed input.
func NormalizeDomain(raw string) (string, bool) {
	d := strings.ToLower(strings.TrimSpace(raw))
	d = strings.TrimSuffix(d, ".")
	if d == "" || len(d) > 253 || !domainPattern.MatchString(d) {
		return "", false
	}
	return d, true
}

// ParseDomains reads an IPFire DBL list, one domain per line. Blank lines
// and lines starting with '#' or ';' are treated as comments. Malformed
// lines are skipped and duplicate domains are collapsed.
func ParseDomains(r io.Reader) []string {
	seen := make(map[string]struct{})
	var domains []string

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Some lines may carry trailing whitespace-separated comments; the
		// domain is always the first field.
		field, _, _ := strings.Cut(line, " ")
		domain, ok := NormalizeDomain(field)
		if !ok {
			continue
		}
		if _, dup := seen[domain]; dup {
			continue
		}
		seen[domain] = struct{}{}
		domains = append(domains, domain)
	}
	if err := scanner.Err(); err != nil {
		return domains
	}

	return domains
}
