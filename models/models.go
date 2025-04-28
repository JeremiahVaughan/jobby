package models

import (
    "fmt"
    "context"

    "github.com/JeremiahVaughan/jobby/clients" 
    "github.com/JeremiahVaughan/jobby/config" 
)

type Models struct {
    AcmeChallenger *AcmeChallengerModel
    DatabaseBackup *DatabaseBackupModel
    Healthy *HealthyModel
}

func New(ctx context.Context, clients *clients.Clients, config config.Config) (*Models, error) {
    acmeModel, err := NewAcmeChallengerModel(ctx, clients, config)
    if err != nil {
        return nil, fmt.Errorf("error, when NewAcmeChallengerModel() for New(). Error: %v", err)
    }
    return &Models{
        DatabaseBackup: NewDatabaseBackupModel(clients),
        AcmeChallenger: acmeModel,
        Healthy: NewHealthyModel(clients),
    }, nil
}
