package clients

import (
    "context"
    "fmt"

    "github.com/JeremiahVaughan/jobby/clients/database" 
    "github.com/JeremiahVaughan/jobby/clients/bucket" 
    "github.com/JeremiahVaughan/jobby/clients/sqlite" 
    "github.com/JeremiahVaughan/jobby/clients/healthy" 
    "github.com/JeremiahVaughan/jobby/clients/lego" 
    "github.com/JeremiahVaughan/jobby/config" 
)

type Clients struct {
    Databases []*database.Client
    Bucket *bucket.Client
    Sqlite *sqlite.Client
    Lego *lego.Client
    Healthy *healthy.Client
}

func New(
    ctx context.Context,
    config config.Clients,
    serviceName string,
) (*Clients, error) {
    theClients := Clients{}
    theClients.Databases = make([]*database.Client, len(config.Databases))
    for i, db := range config.Databases {
        c, err := database.New(ctx, db)
        if err != nil {
            return nil, fmt.Errorf("error, when creating new DB client for clients.New(). Error: %v", err)
        }
        theClients.Databases[i] = c
    }
    var err error
    theClients.Bucket, err = bucket.New(ctx, config.Bucket, serviceName)
    if err != nil {
        return nil, fmt.Errorf("error, when creating new bucket client for clients.New(). Error: %v", err)
    }
    theClients.Sqlite, err = sqlite.New(config.Sqlite)
    if err != nil {
        return nil, fmt.Errorf("error, when creating new sqlite client for clients.New(). Error: %v", err)
    }
    theClients.Healthy, err = healthy.New(config.Nats, serviceName)
    if err != nil {
        return nil, fmt.Errorf("error, when creating new healthy client for clients.New(). Error: %v", err)
    }
    theClients.Lego, err = lego.New(ctx, config.Lego, theClients.Bucket, theClients.Sqlite, theClients.Healthy)
    if err != nil {
        return nil, fmt.Errorf("error, when creating new lego client for clients.New(). Error: %v", err)
    }
    return &theClients, nil
}

func (c *Clients) Destroy() {
    for _, db := range c.Databases {
        db.Destroy()
    }
    c.Healthy.Close()
}
