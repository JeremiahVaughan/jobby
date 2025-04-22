package config

import (
    "fmt"
    "errors"
    "context"
    "encoding/json"
)

type Config struct {
    Clients Clients `json:"clients"`
}

type Clients struct {
    Databases []Database `json:"databases"`
    Bucket Bucket `json:"bucket"`
    Sqlite Sqlite `json:"sqlite"`
    Lego Lego `json:"lego"`
}

type Lego struct {
    CADirURL string `json"caDirURL"`
    Email string `json"email"`
}

type Sqlite struct {
    DataDirectory string `json:"dataDirectory"`
    MigrationDirectory string `json:"migrationDirectory"`
}

type Database struct {
    Host string `json:"host"`
    Name string `json:"name"`
    Username string `json:"username"`
    Password string `json:"password"`
    Port string `json:"port"`
}

type Bucket struct {
    Name string `json:"name"`
}

func New(ctx context.Context) (Config, error) {
    c, err := fetchConfigFromS3(ctx, "jobby")
    if err != nil {
        return Config{}, fmt.Errorf("error, when fetching config file. Error: %v", err)
    }

    result := Config{}
    err = json.Unmarshal(c, &result)
    if err != nil {
        return Config{}, fmt.Errorf("error, when decoding config file. Error: %v", err)
    }

    err = result.isValid() 
    if err != nil {
        return Config{}, fmt.Errorf("error, configuration validation failed. Error: %v", err)
    }

    return result, nil
}

func (cfg *Config) isValid() error {                               
   clients := cfg.Clients                                            
                                                                     
   // Validate Bucket                                                
   if clients.Bucket.Name == "" {                                    
       return errors.New("bucket name must not be empty")            
   }                                                                 
                                                                     
   // Validate Sqlite                                                
   if clients.Sqlite.DataDirectory == "" {                           
       return errors.New("sqlite data directory must not be empty")  
   }                                                                 
   if clients.Sqlite.MigrationDirectory == "" {                      
       return errors.New("sqlite migration directory must not be empty")                                                               
   }                                                                 
                                                                     
   // Validate Databases                                             
   if len(clients.Databases) == 0 {                                  
       return errors.New("at least one database must be specified")  
   }                                                                 

   for _, db := range clients.Databases {                            
       if db.Host == "" {                                            
           return errors.New("database host must not be empty")      
       }                                                             
       if db.Name == "" {                                            
           return errors.New("database name must not be empty")      
       }                                                             
       if db.Username == "" {                                        
           return errors.New("database username must not be empty")  
       }                                                             
       if db.Password == "" {                                        
           return errors.New("database password must not be empty")  
       }                                                             
       if db.Port == "" {                                            
           return errors.New("database port must not be empty")      
       }                                                             
   }                                                                 

   if cfg.Clients.Lego.Email == "" {
       return errors.New("Lego email is required")
   }
                                                                     
   return nil                                                        
}                                                                     


