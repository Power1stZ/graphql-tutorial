CREATE TABLE contact (
    contact_id INTEGER primary key autoincrement,
    name TEXT,
    first_name TEXT,
    last_name TEXT,
    gender_id INTEGER,
    dob DATE,
    email TEXT,
    phone TEXT,
    address TEXT,
    photo_path TEXT,
    created_date DATETIME,
    created_by TEXT
);

CREATE TABLE account (
    account_id INTEGER primary key autoincrement,
    name TEXT UNIQUE,
    status TEXT,
    balance INTEGER
);

CREATE TABLE deposit_transaction (
    trans_id INTEGER primary key autoincrement,
    amount INTEGER,
    status TEXT,
    created_at DATETIME,
    finished_at DATETIME,
    account_id INTEGER
);

CREATE TABLE user (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL
);
