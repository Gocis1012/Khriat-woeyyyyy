ALTER TABLE payments
DROP COLUMN IF EXISTS webhook_processed_at,
DROP COLUMN IF EXISTS provider_tx_id;

DROP TABLE IF EXISTS credit_ledger;