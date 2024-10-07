CREATE TABLE IF NOT EXISTS kv (
  id INTEGER PRIMARY KEY,
  key TEXT NOT NULL,
  value TEXT,
  expires_in TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS history (
  seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  seeded BOOLEAN DEFAULT FALSE
);

INSERT INTO kv (key, value) VALUES ('name', 'Test User');
INSERT INTO kv (key, value) VALUES ('age', '24');
INSERT INTO kv (key, value) VALUES ('email', 'test@example.com');

INSERT INTO history (seeded) VALUES (TRUE);