package controllers


import (
    "context"
    "time"
    "fmt"

    "github.com/JeremiahVaughan/jobby/models" 
)

type DatabaseBackupController struct {
    model *models.DatabaseBackupModel
}

func NewDatabaseBackupController(model *models.DatabaseBackupModel) *DatabaseBackupController {
    return &DatabaseBackupController{ 
        model: model,
    }
}

func (c *DatabaseBackupController) Start(ctx context.Context) error {
    ticker := time.Tick(time.Hour) 
    for {
        select {
        case t := <- ticker:
            if t.Hour() == 5 {
                err := c.model.BackupDatabases(ctx)
                if err != nil {
                    return fmt.Errorf("error, when model.BackupDatabase() for DatabaseBackupController.Start(). Error: %v", err) 
                }
            }
        case <- ctx.Done():
            return nil
        }
    }
}
