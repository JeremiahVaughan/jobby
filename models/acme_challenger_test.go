package models

import (
    "time"
    "testing"
)

func Test_doDomainsMatch(t *testing.T) {

    m := AcmeChallengerModel{}
    t.Run("they do not match as they contain different items", func(t *testing.T) {
        oldDomains := []string{
            "b",
            "a",
            "d",
        }
        newDomains := []string{
            "a",
            "b",
            "c",
        }
        got := m.doDomainsMatch(oldDomains, newDomains)
        if got {
            t.Errorf("error, expected match but did not get")
        }
    })

    t.Run("they do not match as they are different lengths", func(t *testing.T) {
        oldDomains := []string{
            "b",
            "a",
        }
        newDomains := []string{
            "a",
            "b",
            "c",
        }
        got := m.doDomainsMatch(oldDomains, newDomains)
        if got {
            t.Errorf("error, expected match but did not get")
        }
    })

    t.Run("they match even in different orders", func(t *testing.T) {
        oldDomains := []string{
            "b",
            "a",
            "c",
        }
        newDomains := []string{
            "a",
            "b",
            "c",
        }
        got := m.doDomainsMatch(oldDomains, newDomains)
        if !got {
            t.Errorf("error, expected match but did not get")
        }
    })

    t.Run("they match", func(t *testing.T) {
        oldDomains := []string{
            "a",
            "b",
            "c",
        }
        newDomains := []string{
            "a",
            "b",
            "c",
        }
        got := m.doDomainsMatch(oldDomains, newDomains)
        if !got {
            t.Errorf("error, expected match but did not get")
        }
    })

}

func Test_isEligableForRenewal(t *testing.T) {

    expiresTime := int64(1745409458)
    secondsInDay := 86400
    m := AcmeChallengerModel{
        certRenewalWindowDurationInDays: 30,
    }

    t.Run("is eligable equals", func(t *testing.T) {
        currentTime := expiresTime - int64(secondsInDay * m.certRenewalWindowDurationInDays) 
        isEligable, err := m.isWithinRenewalWindow(expiresTime, currentTime)
        if err != nil {
            t.Errorf("error, was expired but did not expect")
        }
        if !isEligable {
            t.Errorf("error, expected eligable for renewal but it was not")
        }
    })

    t.Run("is eligable in window", func(t *testing.T) {
        currentTime := expiresTime - int64(secondsInDay * 28) 
        isEligable, err := m.isWithinRenewalWindow(expiresTime, currentTime)
        if err != nil {
            t.Errorf("error, was expired but did not expect")
        }
        if !isEligable {
            t.Errorf("error, expected eligable for renewal but it was not")
        }
    })

    t.Run("is not in window", func(t *testing.T) {
        currentTime := expiresTime - int64(secondsInDay * 31) 
        isEligable, err := m.isWithinRenewalWindow(expiresTime, currentTime)
        if err != nil {
            t.Errorf("error, was expired but did not expect")
        }
        if isEligable {
            t.Errorf("error, expected not eligable for renewal but it was")
        }
    })

    t.Run("is expired even with an hour buffer to give time for the cert renewal activity duration itself", func(t *testing.T) {
        currentTime := expiresTime - int64(3600)
        isEligable, err := m.isWithinRenewalWindow(expiresTime, currentTime)
        if err == nil {
            t.Errorf("error, expected expired but it wasn't")
        }
        if isEligable {
            t.Errorf("error, expected not eligable for renewal but it was")
        }
    })

    t.Run("is expired even with an hour buffer to give time for the cert renewal activity duration itself 2", func(t *testing.T) {
        currentTime := expiresTime - int64(3599)
        isEligable, err := m.isWithinRenewalWindow(expiresTime, currentTime)
        if err == nil {
            t.Errorf("error, expected expired but it wasn't")
        }
        if isEligable {
            t.Errorf("error, expected not eligable for renewal but it was")
        }
    })

    t.Run("is expired equals", func(t *testing.T) {
        currentTime := expiresTime 
        isEligable, err := m.isWithinRenewalWindow(expiresTime, currentTime)
        if err == nil {
            t.Errorf("error, expected expired but it wasn't")
        }
        if isEligable {
            t.Errorf("error, expected not eligable for renewal but it was")
        }
    })

    t.Run("is well expired", func(t *testing.T) {
        currentTime := expiresTime + int64(secondsInDay + 3) 
        isEligable, err := m.isWithinRenewalWindow(expiresTime, currentTime)
        if err == nil {
            t.Errorf("error, expected expired but it wasn't")
        }
        if isEligable {
            t.Errorf("error, expected not eligable for renewal but it was")
        }
    })


}

func Test_isTooCloseToExpiration(t *testing.T) {
    currentTime := time.Now()

    t.Run("is not too close", func(t *testing.T) {
        expiresAt := currentTime.AddDate(0, 0, 28)
        tooClose := isTooCloseToExpiration(currentTime, expiresAt)
        if tooClose {
            t.Errorf("error, did not want too close but got")
        }
    })

    t.Run("is too close", func(t *testing.T) {
        expiresAt := currentTime.AddDate(0, 0, 15)
        tooClose := isTooCloseToExpiration(currentTime, expiresAt)
        if !tooClose {
            t.Errorf("error, wanted too close but did not get")
        }
    })

}
