package models

import (
    "net/http"

    "github.com/JeremiahVaughan/jobby/clients/healthy"
    "github.com/JeremiahVaughan/jobby/clients"
)

type HealthyModel struct {
    healthy *healthy.Client
}

func NewHealthyModel(clients *clients.Clients) *HealthyModel {
    return &HealthyModel{
        healthy: clients.Healthy,
    }
}

func (m *HealthyModel) ReportUnexpectedError(w http.ResponseWriter, err error) {
    m.healthy.ReportUnexpectedError(w, err)
}

