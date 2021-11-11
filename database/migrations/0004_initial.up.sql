CREATE TABLE dishes(
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       res_id UUID REFERENCES restaurants(id) NOT NULL,
                       name VARCHAR
)