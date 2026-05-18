CREATE TABLE IF NOT EXISTS customers (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document     VARCHAR(50)  NOT NULL UNIQUE,
    name         VARCHAR(255) NOT NULL,
    score        INT          NOT NULL CHECK (score >= 0 AND score <= 1000),
    risk_level   VARCHAR(10)  NOT NULL CHECK (risk_level IN ('LOW', 'MEDIUM', 'HIGH')),
    income_range VARCHAR(50)  NOT NULL,
    status       VARCHAR(20)  NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'UNDER_REVIEW')),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()    
);

CREATE INDEX IF NOT EXISTS idx_customers_document
    ON customers(document);

CREATE INDEX IF NOT EXISTS idx_customers_status
    ON customers(status);

CREATE INDEX IF NOT EXISTS idx_customers_risk_level
    ON customers(risk_level);
