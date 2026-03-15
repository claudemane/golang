-- Users table
CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    name       TEXT      NOT NULL,
    email      TEXT      NOT NULL UNIQUE,
    gender     TEXT      NOT NULL,
    birth_date DATE      NOT NULL
);

-- Many-to-many friends table
CREATE TABLE IF NOT EXISTS user_friends (
    user_id   INTEGER REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    CHECK (user_id <> friend_id)   -- cannot befriend yourself
);

-- Seed 20 users
INSERT INTO users (name, email, gender, birth_date) VALUES
  ('Alice Johnson',   'alice@example.com',   'female', '1990-03-15'),
  ('Bob Smith',       'bob@example.com',     'male',   '1988-07-22'),
  ('Carol White',     'carol@example.com',   'female', '1995-11-05'),
  ('David Brown',     'david@example.com',   'male',   '1992-01-30'),
  ('Eva Martinez',    'eva@example.com',     'female', '1997-06-18'),
  ('Frank Lee',       'frank@example.com',   'male',   '1985-09-10'),
  ('Grace Kim',       'grace@example.com',   'female', '1993-04-25'),
  ('Henry Wilson',    'henry@example.com',   'male',   '1991-12-03'),
  ('Isla Davis',      'isla@example.com',    'female', '1996-08-14'),
  ('Jack Taylor',     'jack@example.com',    'male',   '1989-02-27'),
  ('Karen Anderson',  'karen@example.com',   'female', '1994-05-09'),
  ('Liam Thomas',     'liam@example.com',    'male',   '1998-10-21'),
  ('Mia Jackson',     'mia@example.com',     'female', '1990-07-16'),
  ('Noah Harris',     'noah@example.com',    'male',   '1987-03-08'),
  ('Olivia Martin',   'olivia@example.com',  'female', '1999-01-11'),
  ('Paul Garcia',     'paul@example.com',    'male',   '1986-11-29'),
  ('Quinn Robinson',  'quinn@example.com',   'female', '1993-09-04'),
  ('Ryan Clark',      'ryan@example.com',    'male',   '1995-06-17'),
  ('Sophia Lewis',    'sophia@example.com',  'female', '1991-04-22'),
  ('Tom Walker',      'tom@example.com',     'male',   '1988-08-31')
ON CONFLICT DO NOTHING;

-- Friendships (bidirectional: insert both directions)
-- Alice(1) and Bob(2) are both friends with Carol(3), David(4), Eva(5) — 3 common friends
INSERT INTO user_friends (user_id, friend_id) VALUES
  (1, 3), (3, 1),
  (1, 4), (4, 1),
  (1, 5), (5, 1),
  (2, 3), (3, 2),
  (2, 4), (4, 2),
  (2, 5), (5, 2),
  -- extra friendships
  (1, 6), (6, 1),
  (2, 7), (7, 2),
  (3, 8), (8, 3),
  (4, 9), (9, 4),
  (5, 10),(10, 5),
  (6, 11),(11, 6),
  (7, 12),(12, 7),
  (8, 13),(13, 8),
  (9, 14),(14, 9),
  (10,15),(15,10),
  (11,16),(16,11),
  (12,17),(17,12),
  (13,18),(18,13),
  (14,19),(19,14),
  (15,20),(20,15)
ON CONFLICT DO NOTHING;
