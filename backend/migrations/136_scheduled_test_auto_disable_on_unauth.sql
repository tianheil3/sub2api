-- 136: Add auto_disable_on_unauth column to scheduled_test_plans
-- When enabled, automatically pauses (sets schedulable=false) the account
-- if a scheduled test fails with a 401/Unauthorized response.

ALTER TABLE scheduled_test_plans ADD COLUMN IF NOT EXISTS auto_disable_on_unauth BOOLEAN NOT NULL DEFAULT false;
