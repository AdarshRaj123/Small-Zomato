CREATE TABLE user_profile(
                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             name TEXT NOT NULL,
                             email TEXT NOT NULL,
                             password TEXT NOT NULL,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                             archived_at TIMESTAMP WITH TIME ZONE

);
CREATE TYPE role_type AS ENUM (
    'admin',
    'subadmin',
    'user'
    );
CREATE TABLE IF NOT EXISTS user_roles (
                                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                          user_id UUID REFERENCES user_profile(id) NOT NULL,
                                          role role_type NOT NULL,
                                          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                          archived_at TIMESTAMP WITH TIME ZONE
);
CREATE TABLE IF NOT EXISTS user_address(
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                           user_id UUID REFERENCES user_profile(id) NOT NULL,
                                           latitude TEXT,
                                           longitude TEXT
);
CREATE TABLE IF NOT EXISTS user_session (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                            user_id UUID REFERENCES user_profile(id) NOT NULL,
                                            session_token TEXT NOT NULL,
                                            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
