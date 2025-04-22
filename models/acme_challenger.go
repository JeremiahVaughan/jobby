package models

import (
    "context"

    "github.com/JeremiahVaughan/jobby/clients"
    "github.com/JeremiahVaughan/jobby/clients/sqlite"
)

type AcmeChallengerModel struct {
    sqlite *sqlite.Client
}

func NewAcmeChallengerModel(clients *clients.Clients) *AcmeChallengerModel {
    return &AcmeChallengerModel{
        sqlite: clients.Sqlite,
    }
}

func (c *AcmeChallengerModel) RefreshAllCertificates(ctx context.Context) error {
    return nil
}

func (c *AcmeChallengerModel) register() {

}

func (c *AcmeChallengerModel) renew() {

}

func (c *AcmeChallengerModel) getExpiring() {
}

