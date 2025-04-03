package controllers

import (
    "context"
    "github.com/JeremiahVaughan/jobby/config" 
)

type Controller interface {
    Start(context.Context) error
}

type Controllers []Controller

func New(config config.Config) Controllers {
    var result []Controller
    dbBackup := NewDatabaseBackupController(config)
    result = append(result, dbBackup)
    return result
}
