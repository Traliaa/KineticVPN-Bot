-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_settings (
                          id bigserial PRIMARY KEY,
                          chatID bigint NOT NULL,
                          name text NOT NULL,
                          router_url text,
                          username text,
                          password text,
                          SelectedApps text[],
                          auth_code text,
                          step text

);
