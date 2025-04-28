package controllers

import (
    "context"

    "github.com/JeremiahVaughan/jobby/models" 
    "github.com/JeremiahVaughan/jobby/clients/healthy" 
)

type Controller interface {
    Start(context.Context)
    GetStatusKey() string
    healthy.Controller
}

// Controllers key is the HealthyStatus.StatusKey for the controller
type Controllers map[string]Controller

type HttpControllers struct {
    AcmeChallenger *AcmeChallengerController
    Health *HealthController
}

func New(models *models.Models) (Controllers, *HttpControllers) {
    con := make(map[string]Controller)

    dbBackup := NewDatabaseBackupController(models)
    con[dbBackup.GetStatusKey()] = dbBackup

    acmeChallenger := NewAcmeChallengerController(models) 
    con[acmeChallenger.GetStatusKey()] = acmeChallenger

    httpCon := &HttpControllers{
        AcmeChallenger: NewAcmeChallengerController(models),
        Health: NewHealthController(),
    }
    return con, httpCon
}
