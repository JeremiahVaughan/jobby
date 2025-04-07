package database

import (
    "context"
    "fmt"

    "github.com/JeremiahVaughan/jobby/config" 
    "github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
    Pool *pgxpool.Pool
    Config config.Database
}

func New(ctx context.Context, config config.Database) (*Client, error) {
    connectionString := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        config.Username,
        config.Password,
        config.Host,
        config.Port,
        config.Name,
    )

    pool, err := pgxpool.New(ctx, connectionString)
    if err != nil {
        return nil, fmt.Errorf("error, when conecting to database. Error: %v", err) 
    }

    return &Client{
        Pool: pool,
        Config: config,
    }, nil
}

func (c *Client) Destroy() {
    c.Pool.Close()
}


