package controllers

import (
    "fmt"
    "time"
    "context"

    "github.com/JeremiahVaughan/jobby/models" 
    "github.com/JeremiahVaughan/jobby/clients/healthy" 
)

type AcmeChallengerController struct {
    model *models.AcmeChallengerModel
    healthy *models.HealthyModel
}

func NewAcmeChallengerController(models *models.Models) *AcmeChallengerController {
    return &AcmeChallengerController{
        model: models.AcmeChallenger,
        healthy: models.Healthy,
    }
}

func (c *AcmeChallengerController) Start(ctx context.Context) {

    // initial check so we don't have to wait for job to run if jobby should be stood up for the first time
    err := c.model.RefreshCertificate(ctx)
    if err != nil {
        err = fmt.Errorf("error, when RefreshCertificate() for NewAcmeChallengerModel(). Error: %v", err)
        c.healthy.ReportUnexpectedError(nil, err)
    }

    ticker := time.Tick(time.Hour) 
    for {
        select {
        case t := <- ticker:
            if t.Hour() == 5 {
                err := c.model.RefreshCertificate(ctx)
                if err != nil {
                    err = fmt.Errorf("error, when model.RefreshAllCertificates() for AcmeChallengerController.Start(). Error: %v", err) 
                    c.healthy.ReportUnexpectedError(nil, err)
                }
            }
        case <- ctx.Done():
            return 
        }
    }
}

func (c *AcmeChallengerController) StartChallengeProvider(port string) {
    err := c.model.StartChallengeProvider(port)
    if err != nil {
        err = fmt.Errorf("error, when StartChallengeProvider(). Error: %v", err)
        c.healthy.ReportUnexpectedError(nil, err)
        return
    }
}

func (c *AcmeChallengerController) GetHealthStatus() (*healthy.HealthStatus, error) {
    isHealthy, msg, err := c.model.HealthyCheck()
    if err != nil {
        return nil, fmt.Errorf("error, when HealthyCheck() for HealthyCheck(). Error: %v", err)
    }
    return &healthy.HealthStatus{
        Healthy: isHealthy,
        Message: msg,
    }, nil
}

func (c *AcmeChallengerController) GetStatusKey() string {
    return "certificate_renewal"
}
