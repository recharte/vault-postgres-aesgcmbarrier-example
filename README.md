# vault-postgres-aesgcmbarrier-example

This is an example of how to use Hashicorp Vault's Barrier interface to store AES-GCM encrypted secrets in PostgreSQL.

## Running

```sh
docker compose up -d

PGPASSWORD=pwd psql -h localhost -U admin -d postgres < schema/vault_kv_store_up.sql

go run main.go
```
