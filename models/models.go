package models

import (
    "github.com/JeremiahVaughan/jobby/clients" 
)

type Models struct {
    AcmeChallenger *AcmeChallengerModel
    DatabaseBackup *DatabaseBackupModel
}

func New(clients *clients.Clients) *Models {
    return &Models{
        DatabaseBackup: NewDatabaseBackupModel(clients),
        AcmeChallenger: NewAcmeChallengerModel(clients),
    }
}
