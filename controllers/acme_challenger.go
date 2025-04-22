package controllers

import (
    "fmt"
    "time"
    "context"
    "net/http"

    "github.com/JeremiahVaughan/jobby/models" 
)

type AcmeChallengerController struct {
    model *models.AcmeChallengerModel
}

func NewAcmeChallengerController(models *models.Models) *AcmeChallengerController {
    return &AcmeChallengerController{
        model: models.AcmeChallenger,
    }
}

func (c *AcmeChallengerController) Start(ctx context.Context) error {
    ticker := time.Tick(time.Hour) 
    for {
        select {
        case t := <- ticker:
            if t.Hour() == 5 {
                err := c.model.RefreshAllCertificates(ctx)
                if err != nil {
                    return fmt.Errorf("error, when model.RefreshAllCertificates() for AcmeChallengerController.Start(). Error: %v", err) 
                }
            }
        case <- ctx.Done():
            return nil
        }
    }
}

func (c *AcmeChallengerController) Challenge(w http.ResponseWriter, r *http.Request) {

}
