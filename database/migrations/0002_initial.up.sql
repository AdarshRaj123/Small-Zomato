CREATE TABLE restaurants (
                            name VARCHAR,
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            latitude TEXT,
                            longitude TEXT
)
