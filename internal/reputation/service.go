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

package reputation

import "github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/config"

// BuildSnapshot aggregates per-category domain lists into a single lookup
// map. A domain listed under several categories accumulates every
// category's score; a domain repeated within a category is only counted
// once, since domainsByCategory is expected to already be deduplicated
// (see ipfire.ParseDomains).
//
// Categories are applied in the order given, which fixes the order of the
// resulting Categories slice and keeps output deterministic.
func BuildSnapshot(domainsByCategory map[string][]string, categories []config.CategoryConfig) map[string]DomainReputation {
	acc := make(map[string]*DomainReputation)

	for _, cat := range categories {
		for _, domain := range domainsByCategory[cat.Name] {
			rep, ok := acc[domain]
			if !ok {
				rep = &DomainReputation{}
				acc[domain] = rep
			}
			rep.Score += cat.Score
			rep.Categories = append(rep.Categories, CategoryScore{Category: cat.Name, Score: cat.Score})
		}
	}

	snapshot := make(map[string]DomainReputation, len(acc))
	for domain, rep := range acc {
		snapshot[domain] = *rep
	}
	return snapshot
}
