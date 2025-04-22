package controllers

import (
    "context"
    "github.com/JeremiahVaughan/jobby/models" 
)

type Controller interface {
    Start(context.Context) error
}

type Controllers []Controller

type HttpControllers struct {
    AcmeChallenger *AcmeChallengerController
}

func New(models *models.Models) Controllers {
    var result []Controller
    dbBackup := NewDatabaseBackupController(models.DatabaseBackup)
    result = append(result, dbBackup)
    return result
}

func NewHttpControllers(models *models.Models) *HttpControllers {
    return &HttpControllers{
        AcmeChallenger: NewAcmeChallengerController(models),
    }
}
