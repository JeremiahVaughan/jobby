package controllers

import (
    "context"
    "github.com/JeremiahVaughan/jobby/clients" 
)

type Controller interface {
    Start(context.Context) error
}

type Controllers []Controller

func New(clients *clients.Clients) Controllers {
    var result []Controller
    dbBackup := NewDatabaseBackupController(clients)
    result = append(result, dbBackup)
    return result
}
