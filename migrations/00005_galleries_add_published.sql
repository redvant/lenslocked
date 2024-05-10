-- +goose Up
-- +goose StatementBegin
ALTER TABLE galleries
ADD published boolean NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE galleries
DROP COLUMN published;
-- +goose StatementEnd
