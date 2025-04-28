package models

import (
    "testing"
)

func Test_isDatabaseBackupHealthy(t *testing.T) {

    t.Run("is not healthy", func(t *testing.T) {
        currentTime := int64(200000)
        updatedAt := int64(100)
        healthy := isDatabaseBackupHealthy(updatedAt, currentTime)
        if healthy {
            t.Errorf("error, expected not healthy but it was")
        }
    })

    t.Run("is healthy", func(t *testing.T) {
        currentTime := int64(100000)
        updatedAt := int64(100)
        healthy := isDatabaseBackupHealthy(updatedAt, currentTime)
        if !healthy {
            t.Errorf("error, expected healthy but it was not")
        }
    })

}
