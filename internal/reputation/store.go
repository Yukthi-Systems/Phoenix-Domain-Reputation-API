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

import (
	"sync/atomic"
	"time"
)

// Snapshot is one immutable, fully-built view of domain reputations. A new
// Snapshot is only ever published in full; it is never mutated in place.
type Snapshot struct {
	Domains   map[string]DomainReputation
	UpdatedAt time.Time
}

// Store holds the currently-serving Snapshot behind an atomic pointer so
// that reads on the HTTP hot path never block on, or are blocked by, the
// background updater publishing a new one.
type Store struct {
	v atomic.Pointer[Snapshot]
}

// NewStore returns an empty, not-yet-ready Store.
func NewStore() *Store {
	return &Store{}
}

// Replace atomically publishes snap as the current snapshot.
func (s *Store) Replace(snap *Snapshot) {
	s.v.Store(snap)
}

// Ready reports whether at least one snapshot has been published.
func (s *Store) Ready() bool {
	return s.v.Load() != nil
}

// Snapshot returns the currently-serving snapshot, or nil if none has been
// published yet.
func (s *Store) Snapshot() *Snapshot {
	return s.v.Load()
}

// Get returns the reputation for a normalized domain. It returns the zero
// value (score 0, no categories) if the domain is unknown or no snapshot
// has been published yet.
func (s *Store) Get(domain string) DomainReputation {
	snap := s.v.Load()
	if snap == nil {
		return DomainReputation{}
	}
	return snap.Domains[domain]
}
