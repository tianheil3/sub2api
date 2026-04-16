-- Store full request-body version chains for admin troubleshooting by (request_id, api_key_id).
-- Intentionally separate from usage_logs to keep hot-path query rows slim.

CREATE TABLE IF NOT EXISTS usage_log_conversations (
    id            BIGSERIAL   PRIMARY KEY,
    request_id    TEXT        NOT NULL,
    user_id       BIGINT      NOT NULL,
    api_key_id    BIGINT      NOT NULL,
    sequence      INTEGER     NOT NULL,
    stage         VARCHAR(64) NOT NULL,
    snapshot_kind VARCHAR(64) NOT NULL DEFAULT '',
    platform      VARCHAR(64) NOT NULL DEFAULT '',
    account_id    BIGINT,
    account_name  TEXT        NOT NULL DEFAULT '',
    upstream_url  TEXT        NOT NULL DEFAULT '',
    payload       TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (request_id, api_key_id, sequence)
);

CREATE INDEX IF NOT EXISTS idx_usage_log_conversations_created_at
    ON usage_log_conversations (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_usage_log_conversations_user_created_at
    ON usage_log_conversations (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_usage_log_conversations_api_key_created_at
    ON usage_log_conversations (api_key_id, created_at DESC);

COMMENT ON TABLE usage_log_conversations IS
    'Full request-body version chains keyed by usage request identity for admin troubleshooting.';
COMMENT ON COLUMN usage_log_conversations.request_id IS
    'Usage request identifier (same semantic as usage_logs.request_id).';
COMMENT ON COLUMN usage_log_conversations.user_id IS
    'Owning user identifier for grouping and admin lookups.';
COMMENT ON COLUMN usage_log_conversations.api_key_id IS
    'API key identifier used for grouping and uniqueness with request_id.';
COMMENT ON COLUMN usage_log_conversations.sequence IS
    'Monotonic snapshot order inside one request conversation.';
COMMENT ON COLUMN usage_log_conversations.stage IS
    'Snapshot stage such as inbound, gateway_rewritten, upstream_retry, upstream_final.';
COMMENT ON COLUMN usage_log_conversations.snapshot_kind IS
    'Additional stage detail such as retry_exhausted/http_error/success.';
COMMENT ON COLUMN usage_log_conversations.payload IS
    'Full request payload text as seen at this stage (unredacted and untruncated).';
