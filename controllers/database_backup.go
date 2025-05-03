package controllers


import (
    "context"
    "time"
    "fmt"

    "github.com/JeremiahVaughan/jobby/models" 
    "github.com/JeremiahVaughan/jobby/clients/healthy" 
)

type DatabaseBackupController struct {
    model *models.DatabaseBackupModel
    healthy *models.HealthyModel
}

func NewDatabaseBackupController(models *models.Models) *DatabaseBackupController {
    return &DatabaseBackupController{ 
        model: models.DatabaseBackup,
        healthy: models.Healthy,
    }
}

func (c *DatabaseBackupController) Start(ctx context.Context) {
    ticker := time.Tick(time.Hour) 
    for {
        select {
        case t := <- ticker:
            if t.Hour() == 5 {
                err := c.model.BackupDatabases(ctx)
                if err != nil {
                    err = fmt.Errorf("error, when model.BackupDatabase() for DatabaseBackupController.Start(). Error: %v", err) 
                    c.healthy.ReportUnexpectedError(nil, err)
                }
            }
        case <- ctx.Done():
            return
        }
    }
}


func (c *DatabaseBackupController) GetHealthStatus() (*healthy.HealthStatus, error) {
    isHealthy, err := c.model.HealthyCheck()
    if err != nil {
        return nil, fmt.Errorf("error, when HealthyCheck() for GetHealthStatus(). Error: %v", err)
    }
    var msg string
    if isHealthy {
        msg = "database backups are working nicely"
    } else {
        msg = "database backups are failing"
    }
    return &healthy.HealthStatus{
        Healthy: isHealthy,
        Message: msg,
    }, nil
}

func (c *DatabaseBackupController) GetStatusKey() string {
    return "database_backups"
}
