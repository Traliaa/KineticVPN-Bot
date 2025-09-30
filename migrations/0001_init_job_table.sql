-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_settings (
                          id bigserial PRIMARY KEY,
                          args jsonb NOT NULL DEFAULT '{}',
                          attempt smallint NOT NULL DEFAULT 0,
                          attempted_at timestamptz,
                          attempted_by text[],
                          created_at timestamptz NOT NULL DEFAULT NOW(),
                          errors jsonb[],
                          finalized_at timestamptz,
                          work_name text NOT NULL,
                          max_attempts smallint NOT NULL,
                          metadata jsonb NOT NULL DEFAULT '{}',
                          priority smallint NOT NULL DEFAULT 1,
                          queue text NOT NULL DEFAULT 'default',
                          state cronak_job_state NOT NULL DEFAULT 'available',
                          scheduled_at timestamptz NOT NULL DEFAULT NOW(),
                          tags text[] NOT NULL DEFAULT '{}',
                          unique_key bytea,
                          unique_states bit(8),
                          CONSTRAINT finalized_or_finalized_at_null CHECK (
                              (finalized_at IS NULL AND state NOT IN ('cancelled', 'completed', 'discarded')) OR
                              (finalized_at IS NOT NULL AND state IN ('cancelled', 'completed', 'discarded'))
                              ),
                          CONSTRAINT priority_in_range CHECK (priority >= 1 AND priority <= 4),
                          CONSTRAINT queue_length CHECK (char_length(queue) > 0 AND char_length(queue) < 128),
                          CONSTRAINT work_name_length CHECK (char_length(work_name) > 0 AND char_length(work_name) < 128)
);
