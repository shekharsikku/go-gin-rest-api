CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY,
    owner INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,   
    location TEXT NOT NULL,   
    date DATETIME NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (id) ON DELETE CASCADE
);
