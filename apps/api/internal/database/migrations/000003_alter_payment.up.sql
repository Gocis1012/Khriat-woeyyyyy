CREATE TABLE IF NOT EXISTS credit_ledger (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    amount      NUMERIC(10,2) NOT NULL,  -- บวก = เติม, ลบ = ใช้
    type        VARCHAR(50) NOT NULL,    -- 'topup', 'usage', 'daily_bonus', 'refund'
    ref_id      UUID,                   -- payment_id หรือ trans_word_id
    ref_type    VARCHAR(50),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE payments 
ADD COLUMN provider_tx_id VARCHAR(255) UNIQUE,  -- Omise charge_id / Stripe payment_intent_id
ADD COLUMN webhook_processed_at TIMESTAMPTZ;