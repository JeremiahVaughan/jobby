package sqlite


import (
    "database/sql"
    "os"
    "fmt"
    "errors"

    "github.com/JeremiahVaughan/jobby/config" 
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Client struct {
    conn *sql.DB
    migrationDir string
}


func New(config config.Sqlite) (*Client, error) {
    var err error
    _, err = os.Stat(config.DataDirectory)
    if os.IsNotExist(err) {
        err = os.MkdirAll(config.DataDirectory, 0700)
        if err != nil {
            return nil, fmt.Errorf("error, when creating database data directory. Error: %v", err)
        }
    }
    c := Client{
        migrationDir: config.MigrationDirectory,
    }
    dbFile := fmt.Sprintf("%s/data", config.DataDirectory)
    c.conn, err = sql.Open("sqlite3", dbFile)
    if err != nil {
        return nil, fmt.Errorf("error, when entablishing database connection. Error: %v", err)
    }
    err = c.migrate()
    if err != nil {
        return nil, fmt.Errorf("error, when migrating database files. Error: %v", err)
    }
    return &c, nil
}

func (c *Client) IsUserKeyGenerated() (bool, error) {
    var result bool
    err := c.conn.QueryRow(
        `SELECT user_key_generated
        FROM cert_challenge
        WHERE user_registered = 1`,
    ).Scan(
        &result,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return false, nil
        } else {
            return false, fmt.Errorf("error, when attempting to execute sql statement. Error: %v", err)
        }
    }
    return result, nil
}


func (c *Client) IsUserRegistered() (bool, error) {
    var result bool
    err := c.conn.QueryRow(
        `SELECT user_registered
        FROM cert_challenge
        WHERE user_registered = 1`,
    ).Scan(
        &result,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return false, nil
        } else {
            return false, fmt.Errorf("error, when attempting to execute sql statement. Error: %v", err)
        }
    }
    return result, nil
}

func (c *Client) UpdateCertExpiration(newExpiration int64) error {
    _, err := c.conn.Exec(
        `UPDATE cert_challenge
         SET cert_expires_at = ?
         WHERE user_registered = 1`,
        newExpiration,
    )
    if err != nil {
        return fmt.Errorf("error, when updating certificate expiration. Error: %v", err)
    }
    return nil
}

func (c *Client) MarkUserAsKeyGenerated() error {
    _, err := c.conn.Exec(
        `UPDATE cert_challenge
        SET user_key_generated = 1
        WHERE id = 1`,
    )
    if err != nil {
        return fmt.Errorf("error when inserting new certificate. Error: %v", err)
    }
    return nil
}


func (c *Client) MarkUserAsRegistered() error {
    _, err := c.conn.Exec(
        `UPDATE cert_challenge
        SET user_registered = 1
        WHERE id = 1`,
    )
    if err != nil {
        return fmt.Errorf("error when inserting new certificate. Error: %v", err)
    }
    return nil
}

func (c *Client) FetchCurrentCertExpiration() (int64, error) {
    var result int64
    err := c.conn.QueryRow(
        `SELECT cert_expires_at
        FROM cert_challenge
        WHERE id = 1`,
    ).Scan(
        &result,
    )
    if err != nil {
        return 0, fmt.Errorf("error, when attempting to execute sql statement. Error: %v", err)
    }
    return result, nil
}

func (c *Client) UpdateDatabaseBackupLastUpdated(currentTime int64) error {
    _, err := c.conn.Exec(
        `UPDATE database_backup
        SET backups_completed_at = ?
        WHERE id = 1`,
        currentTime,
    )
    if err != nil {
        return fmt.Errorf("error, when executing query. Error: %v", err)
    }
    return nil
}

func (c *Client) FetchDatabaseBackupLastUpdated() (int64, error) {
    var lastUpdated int64
    err := c.conn.QueryRow(
        `SELECT backups_completed_at
        FROM database_backup
        WHERE id = 1`,
    ).Scan(
        &lastUpdated,
    )
    if err != nil {
        return 0, fmt.Errorf("error, when executing sql statement. Error: %v", err)
    }
    return lastUpdated, nil
}
