-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS menu_item_prices(
    id BIGSERIAL PRIMARY KEY,
    menu_item_id BIGINT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    price_cents INT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'USD',
    valid_from TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_to TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_menu_item_price_currency CHECK (currency ~ '^[A-Z]{3}$')
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_menu_item_active_price ON menu_item_prices(menu_item_id) WHERE valid_to IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_menu_item_price_menu_item_id;
DROP TABLE IF EXISTS menu_item_prices;
-- +goose StatementEnd
