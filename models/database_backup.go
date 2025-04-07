package models

import (
    "log"
    "os"
    "os/exec"
    "fmt"
    "context"
    "strings"

    "github.com/JeremiahVaughan/jobby/clients" 
    "github.com/JeremiahVaughan/jobby/clients/database" 
    "github.com/JeremiahVaughan/jobby/clients/bucket" 
)

type DatabaseBackupModel struct {
    databases []*database.Client
    bucket *bucket.Client
}

func NewDatabaseBackupModel(clients *clients.Clients) *DatabaseBackupModel {
    return &DatabaseBackupModel{ 
        databases: clients.Databases, 
        bucket: clients.Bucket, 
    }
}

func (m *DatabaseBackupModel) BackupDatabases(ctx context.Context) error {
    log.Printf("database backup starting", )
    pgPassContent := strings.Builder{}
    for _, db := range m.databases {
        line := fmt.Sprintf(
            "%s:%s:%s:%s:%s\n",
            db.Config.Host,
            db.Config.Port,
            db.Config.Name,
            db.Config.Username,
            db.Config.Password,
        )
        pgPassContent.WriteString(line)
    }
    file, err := os.OpenFile("/root/.pgpass", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
    if err != nil {
        return fmt.Errorf("error, when creating .pgpass file. Error: %s", err)
    }
    _, err = file.WriteString(pgPassContent.String())
    if err != nil {
        return fmt.Errorf("error, when attempting to write .pgpass file. Error: %v", err)
    }
    file.Close()
    for _, db := range m.databases {
        objectName := fmt.Sprintf("%s_dump.sql", db.Config.Name)
        dumpFileName := fmt.Sprintf("/tmp/%s", objectName)
        cmd := exec.Command(
            "pg_dump",
            "-h", 
            db.Config.Host,
            "-p",
            db.Config.Port,
            "-U",
            db.Config.Username,
            "-d", 
            db.Config.Name,
            "-f",
            dumpFileName,
        )
        output, err := cmd.CombinedOutput()
        if err != nil {
            return fmt.Errorf("error, when executing pg_dump command. stderr: %s. Error: %v", output, err)
        }
        err = m.bucket.BackupFile(ctx, dumpFileName, objectName)
        if err != nil {
            return fmt.Errorf("error, when bucket.BackupFile() for models.DatabaseBackupModel.BackupDatabases(). Error: %v", err)
        }
    }
    log.Printf("database backup completed")
    return nil
}
