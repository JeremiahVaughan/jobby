package controllers


import (
    "context"
    "time"
    "fmt"

    "github.com/JeremiahVaughan/jobby/models" 
    "github.com/JeremiahVaughan/jobby/config" 
)

type DatabaseBackupController struct {
    model *models.DatabaseBackupModel
}

func NewDatabaseBackupController(config config.Config) *DatabaseBackupController {
    dbModel := models.NewDatabaseBackupModel(config)
    return &DatabaseBackupController{ 
        model: dbModel, 
    }
}

func (c *DatabaseBackupController) Start(ctx context.Context) error {
    ticker := time.Tick(time.Hour) 
    select {
    case t := <- ticker:
        if t.Hour() == 5 {
            err := c.model.BackupDatabase()
            if err != nil {
                return fmt.Errorf("error, when model.BackupDatabase() for DatabaseBackupController.Start(). Error: %v", err) 
            }
        }
    case <- ctx.Done():
    }
    return nil
}
