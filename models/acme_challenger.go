package models

import (
    "os"
    "fmt"
    "sort"
    "os/exec"
    "time"
    "context"

    "github.com/JeremiahVaughan/jobby/config"
    "github.com/JeremiahVaughan/jobby/clients"
    "github.com/JeremiahVaughan/jobby/clients/sqlite"
    "github.com/JeremiahVaughan/jobby/clients/bucket"
    "github.com/JeremiahVaughan/jobby/clients/lego"
	"github.com/go-acme/lego/v4/certificate"
)

type AcmeChallengerModel struct {
    sqlite *sqlite.Client
    lego *lego.Client
    bucket *bucket.Client
    certRenewalWindowDurationInDays int
    certValidDurationInDays int
    certEmplacementLocation string
    domains []string
}

func NewAcmeChallengerModel(
    ctx context.Context,
    clients *clients.Clients,
    config config.Config,
) (*AcmeChallengerModel, error) {
    model := &AcmeChallengerModel{
        sqlite: clients.Sqlite,
        lego: clients.Lego,
        bucket: clients.Bucket,
        certRenewalWindowDurationInDays: config.CertRenewalWindowDurationInDays,
        certValidDurationInDays: config.CertValidDurationInDays,
        domains: config.Clients.Lego.Domains,
    }

    return model, nil
}


func (c *AcmeChallengerModel) StartChallengeProvider(port string) error {
    err := c.lego.StartChallengeProvider(port)
    if err != nil {
        return err
    }
    return nil
}


func (m *AcmeChallengerModel) RefreshCertificate(ctx context.Context) error {
    expiresAt, err := m.sqlite.FetchCurrentCertExpiration()
    if err != nil {
        return fmt.Errorf("error, when FetchCurrentCertExpiration() for RefreshCertificate(). Error: %v", err) 
    }

    withinWindow, err := m.isWithinRenewalWindow(
        expiresAt,
        time.Now().Unix(),
    ) 
    if err != nil {
        var cert *certificate.Resource
        cert, err = m.lego.Obtain(m.domains)
        if err != nil {
            return fmt.Errorf("error, when Obtain() for RefreshCertificate(). Error: %v", err)
        }
        err = m.recordCertificates(ctx, cert, true)
        if err != nil {
            return fmt.Errorf("error, when recordCertificates() for RefreshCertificate(). Error: %v", err)
        }
        return nil
    }

    var cert *certificate.Resource
    if withinWindow {
        oldDomains, err := m.bucket.DownloadDomainsFromBucket(ctx)
        if err != nil {
            return fmt.Errorf("error, when DownloadLegoRegistrationKeyFromBucket() for RefreshCertificate(). Error: %v", err)
        }
        domainsMatch := m.doDomainsMatch(oldDomains, m.domains)
        if domainsMatch {
            oldCert, err := m.bucket.DownloadCertFromConfigBunker(ctx)
            if err != nil {
                return fmt.Errorf("error, when DownloadCertFromConfigBunker() for RefreshCertificate(). Error: %v", err)
            }
            cert, err = m.lego.Renew(oldCert)
            if err != nil {
                return fmt.Errorf("error, when Renew() for RefreshCertificate(). Error: %v", err)
            }
        } else {
            cert, err = m.lego.Obtain(m.domains)
            if err != nil {
                return fmt.Errorf("error, when Obtain() for RefreshCertificate(). Error: %v", err)
            }
        }
        err = m.recordCertificates(ctx, cert, !domainsMatch)
        if err != nil {
            return fmt.Errorf("error, when recordCertificates() while withinWindow for RefreshCertificate(). Error: %v", err)
        }
        return nil
    } 

    _, err = os.Stat(m.certEmplacementLocation)
    if os.IsNotExist(err) {
        privateKey, err := m.bucket.DownloadCertKeyFromConfigBunker(ctx)
        if err != nil {
            return fmt.Errorf("error, when DownloadCertKeyFromConfigBunker() for RefreshCertificate(). Error: %v", err)
        }
        cert, err := m.bucket.DownloadCertFromConfigBunker(ctx)
        if err != nil {
            return fmt.Errorf("error, when DownloadCertFromConfigBunker() for RefreshCertificate(). Error: %v", err)
        }
        err = m.emplaceCertificates(privateKey, cert)
        if err != nil {
            return fmt.Errorf("error, when emplaceCertificates() for RefreshCertificate(). Error: %v", err)
        }
        return nil
    } 
    return nil
}

func (m *AcmeChallengerModel) recordCertificates(ctx context.Context, cert *certificate.Resource, domainChange bool) error {
    var err error

    err = m.bucket.UploadCertToConfigBunker(ctx, cert.PrivateKey)
    if err != nil {
        return fmt.Errorf("error, when UploadToConfigBunker() for RefreshCertificate() for cert key. Error: %v", err)
    }

    err = m.bucket.UploadCertKeyToConfigBunker(ctx, cert.Certificate)
    if err != nil {
        return fmt.Errorf("error, when UploadToConfigBunker() for RefreshCertificate() for cert. Error: %v", err)
    }

    err = m.emplaceCertificates(cert.PrivateKey, cert.Certificate)
    if err != nil {
        return fmt.Errorf("error, when emplaceCertificates() for recordCertificates(). Error: %v", err)
    }

    err = m.reloadHaProxy()
    if err != nil {
        return fmt.Errorf("error, when reloadHaProxy() for recordCertificates(). Error: %v", err)
    }

    // record success to db, but then store new domains in the bucket after recording success to db,
    // this means that should the upload to bucket fail then at least it would just execute an obtain
    // instead of a renewal on the next renewal cycle rather than having incorrect domains in the cert.
    newExpiration := m.getNewExpiration()
    err = m.sqlite.UpdateCertExpiration(newExpiration)
    if err != nil {
        return fmt.Errorf("error, when UpdateCertExpiration() for recordCertificates(). Error: %v", err)
    }

    if domainChange {
        err = m.bucket.UploadDomainsToBucket(ctx, m.domains)
        if err != nil {
            return fmt.Errorf("error, when UploadDomainsToBucket() while withinWindow for recordCertificates(). Error: %v", err)
        }
    }
    return nil
}

func (m *AcmeChallengerModel) reloadHaProxy() error {
    cmd := exec.Command("systemctl", "reload", "haproxy.service")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("error, when attempting to reload HaProxy. Ouput: %s. Error: %v", output, err)
    }
    return nil
}

func (m *AcmeChallengerModel) getNewExpiration() int64 {
    return time.Now().AddDate(0, 0, m.certValidDurationInDays).Unix()
}

func (m *AcmeChallengerModel) emplaceCertificates(privateKey, cert []byte) error {
    file, err := os.Open(m.certEmplacementLocation)
    if err != nil {
        return fmt.Errorf("error, when opening file for emplaceCertificates(). Error: %v", err)
    }
    certs := m.stackCert(privateKey, cert)
    _, err = file.Write(certs)
    if err != nil {
        return fmt.Errorf("error, when writing to file for emplaceCertificates(). Error: %v", err)
    }
    return nil
}




func (m *AcmeChallengerModel) doDomainsMatch(oldDomains, newDomains []string) bool {
    if len(oldDomains) != len(newDomains) {
        return false
    }
    oldMap := make(map[string]bool)
    for _, o := range oldDomains {
        oldMap[o] = true
    }
    for _, n := range newDomains {
        _, ok := oldMap[n]
        if !ok {
            return false
        }
    }
    return true
}

func (m *AcmeChallengerModel) sortDomains(domains []string) {
    sort.Slice(domains, func (i, j int) bool { return domains[i] > domains[i] })
}

func (m *AcmeChallengerModel) isWithinRenewalWindow(
    expiresAt int64,
    currentTime int64,
) (bool, error) {
    oneHourBufferForRenewalActivity := int64(3600)
    expiresAt = expiresAt - oneHourBufferForRenewalActivity

    secondsInDay := int64(86400)
	renewalWindowStart := expiresAt - (secondsInDay * int64(m.certRenewalWindowDurationInDays))

	if currentTime >= expiresAt {
		return false, fmt.Errorf("cannot renew: already expired at %d", expiresAt)
	}

	return currentTime >= renewalWindowStart, nil
}

func (m *AcmeChallengerModel) stackCert(privateKey, certificate []byte) []byte {
    return append(privateKey, certificate...)
}

func (m *AcmeChallengerModel) HealthyCheck() (bool, string, error) {
    expiresAt, err := m.sqlite.FetchCurrentCertExpiration()
    if err != nil {
        return false, "", fmt.Errorf("error, when FetchCurrentCertExpiration() for HealthyCheck(). Error: %v", err) 
    }

    withinWindow, err := m.isWithinRenewalWindow(
        expiresAt,
        time.Now().Unix(),
    ) 
    if err != nil {
        return false, "certificate has expired", nil
    }
    if !withinWindow {
        return true, "certificate is not within renewal window yet", nil
    }
    currentTime := time.Now()
    expireTime := time.Unix(expiresAt, 0)
    tooClose := isTooCloseToExpiration(currentTime, expireTime)
    var msg string
    if tooClose {
        msg = "certificate is not getting renewed for some reason"
    } else {
        msg = "certificate is within window, and should be renewed soon"
    }
    return tooClose, msg, nil
}

func isTooCloseToExpiration(currentTime, expireTime time.Time) bool {
    threshold := currentTime.AddDate(0, 0, 20)
    return expireTime.Before(threshold)
}

