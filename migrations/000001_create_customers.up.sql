CREATE TABLE IF NOT EXISTS customers (
    id           UUID PRIMARY KEY,
    document     VARCHAR(50)  NOT NULL,
    name         VARCHAR(255) NOT NULL,
    score        INT          NOT NULL,
    risk_level   VARCHAR(10)  NOT NULL,
    income_range VARCHAR(50)  NOT NULL,
    status       VARCHAR(20)  NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL,
    updated_at   TIMESTAMPTZ  NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS customers_document_unique_idx
    ON customers (document);
