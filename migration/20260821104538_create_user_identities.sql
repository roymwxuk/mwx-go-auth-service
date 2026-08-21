-- +goose Up
CREATE TABLE user_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    provider TEXT NOT NULL,
    provider_user_id TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_provider_identity
        UNIQUE (provider, provider_user_id)
);

CREATE INDEX idx_user_identities_user_id
ON user_identities(user_id);


-- +goose Down
DROP TABLE IF EXISTS user_identities;
