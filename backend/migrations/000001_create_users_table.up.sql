CREATE TABLE IF NOT EXISTS enterprises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    enterprise_id INTEGER REFERENCES enterprises(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS participants (
    id SERIAL PRIMARY KEY,
    event_id INTEGER REFERENCES events(id),
    name VARCHAR(255) NOT NULL
);

CREATE TABLE posts(
    id SERIAL PRIMARY KEY,
    event_id INTEGER REFERENCES events(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments(
     id SERIAL PRIMARY KEY,
     post_id INTEGER REFERENCES posts(id),
     participant_id INTEGER REFERENCES participants(id),
     content TEXT NOT NULL,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)