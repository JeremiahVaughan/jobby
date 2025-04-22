package clients

import (
    "context"
    "fmt"

    "github.com/JeremiahVaughan/jobby/clients/database" 
    "github.com/JeremiahVaughan/jobby/clients/bucket" 
    "github.com/JeremiahVaughan/jobby/clients/sqlite" 
    "github.com/JeremiahVaughan/jobby/clients/lego" 
    "github.com/JeremiahVaughan/jobby/config" 
)

type Clients struct {
    Databases []*database.Client
    Bucket *bucket.Client
    Sqlite *sqlite.Client
    Lego *lego.Client
}

func New(ctx context.Context, config config.Clients) (*Clients, error) {
    theClients := Clients{}
    for _, db := range config.Databases {
        c, err := database.New(ctx, db)
        if err != nil {
            return nil, fmt.Errorf("error, when creating new DB client for clients.New(). Error: %v", err)
        }
        theClients.Databases = append(theClients.Databases, c)
    }
    var err error
    theClients.Bucket, err = bucket.New(ctx, config.Bucket)
    if err != nil {
        return nil, fmt.Errorf("error, when creating new bucket client for clients.New(). Error: %v", err)
    }
    theClients.Sqlite, err = sqlite.New(config.Sqlite)
    if err != nil {
        return nil, fmt.Errorf("error, when creating new sqlite client for clients.New(). Error: %v", err)
    }
    // todo implement and uncomment
    // theClients.Lego, err = lego.New(config.Lego)
    // if err != nil {
    //     return nil, fmt.Errorf("error, when creating new lego client for clients.New(). Error: %v", err)
    // }
    return &theClients, nil
}

func (c *Clients) Destroy() {
    for _, db := range c.Databases {
        db.Destroy()
    }
}
