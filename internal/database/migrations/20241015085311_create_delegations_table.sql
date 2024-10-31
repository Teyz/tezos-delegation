-- +goose Up
-- +goose StatementBegin
CREATE TABLE delegations (
    id                          BIGSERIAL       PRIMARY KEY NOT NULL,
    timestamp                   TIMESTAMP(6)    NOT NULL,
    amount                      BIGINT          NOT NULL,
    delegator                   VARCHAR(36)     NOT NULL,
    level                       BIGINT          NOT NULL
);

CREATE INDEX idx_delegations_timestamp ON delegations(timestamp);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_delegations_timestamp;
DROP TABLE delegations;
-- +goose StatementEnd
