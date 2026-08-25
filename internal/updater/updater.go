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

// Package updater periodically synchronizes the reputation store from
// IPFire DBL category lists.
package updater

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/config"
	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/ipfire"
	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/reputation"
)

// Updater periodically downloads all configured IPFire categories, builds a
// new reputation snapshot, and publishes it. A failed cycle leaves the
// previously published snapshot untouched.
type Updater struct {
	downloader ipfire.Downloader
	store      *reputation.Store
	categories []config.CategoryConfig
	interval   time.Duration
	logger     *slog.Logger
}

// New builds an Updater. downloader is typically an *ipfire.Client but any
// ipfire.Downloader works, which keeps this package testable without
// network access.
func New(downloader ipfire.Downloader, store *reputation.Store, categories []config.CategoryConfig, interval time.Duration, logger *slog.Logger) *Updater {
	return &Updater{
		downloader: downloader,
		store:      store,
		categories: categories,
		interval:   interval,
		logger:     logger,
	}
}

// UpdateOnce runs a single synchronization cycle: download every category,
// and only if all succeed, build and publish a new snapshot. It is
// all-or-nothing — a single category failure aborts the cycle and leaves
// the store serving whatever it served before.
func (u *Updater) UpdateOnce(ctx context.Context) error {
	domainsByCategory := make(map[string][]string, len(u.categories))

	for _, cat := range u.categories {
		domains, err := u.downloader.DownloadDomains(ctx, ipfire.Category{Name: cat.Name, URL: cat.URL})
		if err != nil {
			return fmt.Errorf("update category %s: %w", cat.Name, err)
		}
		u.logger.Info("ipfire category downloaded", "category", cat.Name, "domains", len(domains))
		domainsByCategory[cat.Name] = domains
	}

	domains := reputation.BuildSnapshot(domainsByCategory, u.categories)
	u.store.Replace(&reputation.Snapshot{Domains: domains, UpdatedAt: time.Now()})
	u.logger.Info("reputation snapshot published", "domains", len(domains))
	return nil
}

// Run blocks, executing UpdateOnce every interval until ctx is cancelled.
// It does not perform an initial update on entry; callers should run
// UpdateOnce once, synchronously, before starting Run so the service has a
// snapshot as soon as possible after startup.
func (u *Updater) Run(ctx context.Context) {
	ticker := time.NewTicker(u.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			u.logger.Info("stopping ipfire updater")
			return
		case <-ticker.C:
			if err := u.UpdateOnce(ctx); err != nil {
				u.logger.Error("ipfire update failed, keeping previous snapshot", "error", err)
			}
		}
	}
}
