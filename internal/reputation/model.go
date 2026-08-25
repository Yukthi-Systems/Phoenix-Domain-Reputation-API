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

// Package reputation builds and serves the in-memory domain reputation
// snapshot consumed by the HTTP API.
package reputation

// CategoryScore is the contribution of a single IPFire category to a
// domain's total score.
type CategoryScore struct {
	Category string `json:"category"`
	Score    int    `json:"score"`
}

// DomainReputation is the aggregated reputation of a single domain across
// all configured categories.
type DomainReputation struct {
	Score      int             `json:"score"`
	Categories []CategoryScore `json:"categories"`
}
