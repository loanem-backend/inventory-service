CREATE TYPE loan_status AS ENUM (
    'UPCOMING',
    'ONGOING',
    'DONE',
    'CANCELLED',
    'EXPIRED'
);

CREATE TABLE IF NOT EXISTS loans (
    id CHAR(26) PRIMARY KEY,
    toolkit_id SMALLINT NULL REFERENCES toolkits(id) ON DELETE SET NULL,
    team_id CHAR(36) NULL,
    submitter_id CHAR(36) NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    session_number SMALLINT,
    status loan_status,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);