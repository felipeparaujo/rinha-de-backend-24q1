CREATE TABLE u (
    id INTEGER PRIMARY KEY,
    l INTEGER, -- limit
    b INTEGER  -- balance
);

CREATE TABLE t (
    id INTEGER PRIMARY KEY,
    t TIMESTAMP DEFAULT CURRENT_TIMESTAMP, --time
    a INTEGER, -- amount, may be negative.
    d VARCHAR(10) -- description
);

