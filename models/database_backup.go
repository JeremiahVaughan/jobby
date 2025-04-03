package models

import (
    "log"

    "github.com/JeremiahVaughan/jobby/clients/database" 
    "github.com/JeremiahVaughan/jobby/clients/bucket" 
    "github.com/JeremiahVaughan/jobby/config" 
)

type DatabaseBackupModel struct {
    db *database.Client
    bucket *bucket.Client
}

func NewDatabaseBackupModel(config config.Config) *DatabaseBackupModel {
    dbClient := database.New(config.Clients.Database)
    bucketClient := bucket.New(config.Clients.Bucket)
    return &DatabaseBackupModel{ 
        db: dbClient, 
        bucket: bucketClient, 
    }
}

func (m *DatabaseBackupModel) BackupDatabase() error {
    // todo implement
    log.Printf("database backup starting")
    // do work

    log.Printf("database backup completed")
    return nil
}
