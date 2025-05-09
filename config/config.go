package config

import (
    "fmt"
    "errors"
    "context"
    "encoding/json"
)

type Config struct {
    Clients Clients `json:"clients"`
    CertRenewalWindowDurationInDays int `json:"certRenewalWindowDurationInDays"`
    CertValidDurationInDays int `json:"certValidDurationInDays"`
    CertEmplacementLocation string `json:"certEmplacementLocation"`
}

type Clients struct {
    Databases []Database `json:"databases"`
    Bucket Bucket `json:"bucket"`
    Sqlite Sqlite `json:"sqlite"`
    Lego Lego `json:"lego"`
    Nats Nats `json:"nats"`
}

type Nats struct {
    Host string `json:"host"`
    Port int `json:"port"`
}

type Lego struct {
    Email string `json:"email"`
    Domains []string `json:"domains"`
    CaDirUrl string `json:"caDirUrl"`
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
    BitBunkerBucketName string `json:"bitBunkerBucketName"`
    ConfigBunkerBucketName string `json:"configBunkerBucketName"`
    LegoRegistrationFileName string `json:"legoRegistrationFileName"`
    LegoRegistrationKeyFileName string `json:"legoRegistrationKeyFileName"`
    CertName string `json:"certName"`
    CertKeyName string `json:"certKeyName"`
    DomainsFileName string `json:"domainsFileName"`
    UseTestDir bool `json:"useTestDir"`
}

func New(ctx context.Context, serviceName string) (Config, error) {
    bytes, err := fetchConfigFromS3(ctx, serviceName)
    if err != nil {
        return Config{}, fmt.Errorf("error, when fetching config file. Error: %v", err)
    }

    result := Config{}
    err = json.Unmarshal(bytes, &result)
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
   if clients.Bucket.BitBunkerBucketName == "" {                                    
       return errors.New("bitBunkerBucketName must not be empty")            
   }                                                                 
   if clients.Bucket.ConfigBunkerBucketName == "" {
       return errors.New("configBunkerBucketName must not be empty")
   }
   if clients.Bucket.LegoRegistrationFileName == "" {
       return errors.New("legoRegistrationFileName must not be empty")
   }
   if clients.Bucket.LegoRegistrationKeyFileName == "" {
       return errors.New("legoRegistrationKeyFileName must not be empty")
   }
   if clients.Bucket.CertName == "" {
       return errors.New("certName must not be empty")
   }
   if clients.Bucket.CertKeyName == "" {
       return errors.New("certKeyName must not be empty")
   }
   if clients.Bucket.DomainsFileName == "" {
       return errors.New("domainsFileName must not be empty")
   }
   if clients.Sqlite.DataDirectory == "" {                           
       return errors.New("sqlite data directory must not be empty")  
   }                                                                 
   if clients.Sqlite.MigrationDirectory == "" {                      
       return errors.New("sqlite migration directory must not be empty")                                                               
   }                                                                 
                                                                     
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

   if len(cfg.Clients.Lego.Domains) < 1 {
       return errors.New("must provide at least one domain")
   }

   if hasDuplicates(cfg.Clients.Lego.Domains) {
       return errors.New("duplicate domains detected")
   }

   if cfg.CertRenewalWindowDurationInDays == 0 {
       return errors.New("must provide a cert renewal duration")
   }
                                                                     
   if cfg.Clients.Nats.Host == "" {
       return errors.New("must provide nats host")
   }

   if cfg.Clients.Nats.Port == 0 {
       return errors.New("must provide nats port")
   }

   return nil                                                        
}                                                                     

func hasDuplicates(domains []string) bool {                         
   domainMap := make(map[string]bool)                                
   for _, domain := range domains {                                  
       if _, exists := domainMap[domain]; exists {                   
           return true                                              
       }                                                             
       domainMap[domain] = true                                      
   }                                                                 
   return false                                                       
}
