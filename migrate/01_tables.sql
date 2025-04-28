CREATE TABLE cert_challenge (
    id INTEGER,
    user_key_generated INTEGER,
    user_registered INTEGER, 
    cert_expires_at INTEGER NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO cert_challenge (user_key_generated, user_registered, cert_expires_at)
         VALUES (0, 0, 0);

CREATE TABLE database_backup (
    id INTEGER,
    backups_completed_at INTEGER NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO database_backup (backups_completed_at)
VALUES (0);



