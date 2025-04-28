package lego

import (
    "fmt"
    "log"
    "crypto"
    "crypto/rand"
    "crypto/ecdsa"
    "crypto/elliptic"
    "context"
    "github.com/JeremiahVaughan/jobby/config" 
    "github.com/JeremiahVaughan/jobby/clients/sqlite" 
    "github.com/JeremiahVaughan/jobby/clients/bucket" 
    "github.com/JeremiahVaughan/jobby/clients/healthy" 
    "github.com/go-acme/lego/v4/registration"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/certcrypto"
    "github.com/go-acme/lego/v4/lego"
    "github.com/go-acme/lego/v4/challenge/tlsalpn01"
)

// Client this client depends on other clients due to a race condition issue seen below
type Client struct {
    config config.Lego
    legoClient *lego.Client
    sqlite *sqlite.Client
    bucket *bucket.Client
    healthy *healthy.Client
}

func (c *Client) GetEmail() string {
    log.Printf("email being used: %v", c.config.Email)
	return c.config.Email
}

// GetPrivateKey can't fail fast because the lego client doesn't allow error returns. 
// This shouldn't be a huge deal though as letsencrypt is just going to reject the request with a 400.
// we are inline grabbing the registration because otherwise we keep hitting weird race conditions
func (c *Client) GetRegistration() *registration.Resource {
    ctx := context.Background()
    userRegistered, err := c.sqlite.IsUserRegistered()
    if err != nil {
        err = fmt.Errorf("error, when IsUserRegistered() for NewAcmeChallengerModel(). Error: %v", err)
        c.healthy.ReportUnexpectedError(nil, err)
        return nil
    }

    var reg *registration.Resource
    if !userRegistered {
        reg, err = c.register(ctx, c.GetPrivateKey())
        if err != nil {
            err = fmt.Errorf("error, when register() for NewAcmeChallengerModel(). Error: %v", err)
            c.healthy.ReportUnexpectedError(nil, err)
            return nil
        }
        log.Printf("registration response: %v", reg)
    } else {
        reg, err = c.bucket.DownloadLegoRegistrationFromBucket(ctx)
        if err != nil {
            err = fmt.Errorf("error, when DownloadLegoRegistrationFromBucket() for NewAcmeChallengerModel(). Error: %v", err)
            c.healthy.ReportUnexpectedError(nil, err)
            return nil
        }
        log.Printf("registration download: %v", reg)
    }
	return reg
}

// GetPrivateKey can't fail fast because the lego client doesn't allow error returns
// This shouldn't be a huge deal though as letsencrypt is just going to reject the request with a 400.
// we are inline grabbing the key because otherwise we keep hitting weird race conditions
func (c *Client) GetPrivateKey() crypto.PrivateKey {
    ctx := context.Background()
    userKeyGenerated, err := c.sqlite.IsUserKeyGenerated()
    if err != nil {
        err = fmt.Errorf("error, when IsUserKeyGenerated() for NewAcmeChallengerModel(). Error: %v", err)
        c.healthy.ReportUnexpectedError(nil, err)
        return nil
    }
    var privateKey *ecdsa.PrivateKey
    if !userKeyGenerated {
        privateKey, err = c.generateKey(ctx)
        if err != nil {
            err = fmt.Errorf("error, when generateKey() for NewAcmeChallengerModel(). Error: %v", err)
            c.healthy.ReportUnexpectedError(nil, err)
            return nil
        }
    } else {
        privateKey, err = c.bucket.DownloadLegoRegistrationKeyFromBucket(ctx)
        if err != nil {
            err = fmt.Errorf("error, when DownloadLegoRegistrationKeyFromBucket() for NewAcmeChallengerModel(). Error: %v", err)
            c.healthy.ReportUnexpectedError(nil, err)
            return nil
        }
    }
	return privateKey
}

func New(config config.Lego, bucket *bucket.Client, sqlite *sqlite.Client, healthy *healthy.Client) (*Client, error) {
    c := &Client{
		config: config,
        bucket: bucket,
        sqlite: sqlite,
        healthy: healthy,
    }

	legoConfig := lego.NewConfig(c)

    // uncomment for testing
    // legoConfig.CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"

    // todo look into providing the newer smaller certificates as well should clients support it
	legoConfig.Certificate.KeyType = certcrypto.RSA2048

    var err error
	c.legoClient, err = lego.NewClient(legoConfig)
	if err != nil {
        return nil, fmt.Errorf("error, when creating new lego client. Error: %v", err)
	}

    return c, nil
}

func (c *Client) register(ctx context.Context, key crypto.PrivateKey) (*registration.Resource, error) {
    reg, err := c.legoClient.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
    if err != nil {
        return nil, fmt.Errorf("error, when registering new lego user. Error: %v", err)
    }

    err = c.bucket.UploadLegoRegistrationToBucket(ctx, reg)
    if err != nil {
        return nil, fmt.Errorf("error, when UploadLegoRegistrationToBucket() for register(). Error: %v", err)
    }

    err = c.sqlite.MarkUserAsRegistered()
    if err != nil {
        return nil, fmt.Errorf("error, when MarkUserAsRegistered() for register(). Error: %v", err)
    }

    return reg, nil
}

func (c *Client) Obtain(domains []string) (*certificate.Resource, error) {
    request := certificate.ObtainRequest{
		Domains: c.config.Domains,
		Bundle:  true,
	}
	resource, err := c.legoClient.Certificate.Obtain(request)
	if err != nil {
        return nil, fmt.Errorf("error, when requesting new certificate. Error: %v", err)
	}
    return resource, nil
}

func (c *Client) Renew(oldCert []byte) (resource *certificate.Resource, err error) {
    resource, err = c.legoClient.Certificate.RenewWithOptions(
        certificate.Resource{
            Certificate: []byte(oldCert),
        },
        &certificate.RenewOptions{
            Bundle:  true,
        },
    )
    if err != nil {
        return nil, fmt.Errorf("error, when requesting renewed certificate. Error: %v", err)
    }
    return resource, nil
}

func (c *Client) StartChallengeProvider(port string) error {
    err := c.legoClient.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", port))
    if err != nil {
        return err
    }
    return nil
}

func (c *Client) generateKey(ctx context.Context) (*ecdsa.PrivateKey, error) {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        return nil, fmt.Errorf("error, when generating private key. Error: %v", err)
    }

    err = c.bucket.UploadLegoRegistrationKeyToBucket(ctx, privateKey)
    if err != nil {
        return nil, fmt.Errorf("error, when UploadLegoRegistrationKeyToBucket() for generateKey(). Error: %v", err)
    }

    err = c.sqlite.MarkUserAsKeyGenerated()
    if err != nil {
        return nil, fmt.Errorf("error, when MarkUserAsKeyGenerated() for generateKey(). Error: %v", err)
    }
    return privateKey, nil
}
