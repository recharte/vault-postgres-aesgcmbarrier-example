package main

import (
	"context"
	"crypto/rand"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	physPostgreSQL "github.com/hashicorp/vault/physical/postgresql"
	vault "github.com/hashicorp/vault/vault"
)

func main() {
	logger := logging.NewVaultLogger(log.Debug)

	connURL := "postgres://admin:pwd@localhost:5432/postgres?sslmode=disable"
	table := "vault_kv_store"
	haenabled := "false"

	physical, err := physPostgreSQL.NewPostgreSQLBackend(map[string]string{
		"connection_url": connURL,
		"table":          table,
		"ha_enabled":     haenabled,
	}, logger)
	if err != nil {
		logger.Error("failed to create PostgreSQL backend", "error", err)
		return
	}

	logger.Info("connected to PostgreSQL backend", "connection_url", connURL)

	barrier, err := vault.NewAESGCMBarrier(physical)
	if err != nil {
		logger.Error("failed to create barrier", "error", err)
		return
	}

	logger.Info("created barrier")

	ctx := context.Background()

	// Initialize the barrier
	err = barrier.Initialize(ctx, []byte("12345678901234567890123456789012"), nil, rand.Reader)
	if err != nil && err != vault.ErrBarrierAlreadyInit {
		logger.Error("failed to initialize barrier", "error", err)
		return
	}

	logger.Info("initialized barrier")

	// Unseal the barrier
	if err := barrier.Unseal(ctx, []byte("12345678901234567890123456789012")); err != nil {
		logger.Error("failed to unseal barrier", "error", err)
		return
	}

	logger.Info("unsealed barrier")

	// List entries in root
	list, err := barrier.List(ctx, "")
	if err != nil {
		logger.Error("Error listing entries", "error", err)
		return
	}

	logger.Info("listed entries", "entries", list)

	// Put an entry in nested/path
	entry := &logical.StorageEntry{
		Key:   "nested/path/foo",
		Value: []byte("bar"),
	}
	err = barrier.Put(ctx, entry)
	if err != nil {
		logger.Error("Error putting entry", "error", err)
		return
	}

	logger.Info("put entry", "entry", entry)

	// List entries in nested/path
	list, err = barrier.List(ctx, "nested/path/")
	if err != nil {
		logger.Error("Error listing entries", "error", err)
		return
	}

	logger.Info("listed entries", "entries", list)

	// Get the foo entry from nested/path
	storedEntry, err := barrier.Get(ctx, "nested/path/foo")
	if err != nil {
		logger.Error("Error getting entry", "error", err)
		return
	}

	logger.Info("got entry", "entry", storedEntry)
}
