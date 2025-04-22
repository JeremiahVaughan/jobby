package lego

import (
    "fmt"
    "crypto"
    "github.com/JeremiahVaughan/jobby/config" 
    "github.com/go-acme/lego/v4/registration"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/certcrypto"
    "github.com/go-acme/lego/v4/lego"
)

type Client struct {
    config config.Lego
	email        string
	registration *registration.Resource
	key          crypto.PrivateKey
    legoClient *lego.Client
}

func (c *Client) GetEmail() string {
    // todo implement
	return c.email
}
func (c Client) GetRegistration() *registration.Resource {
    // todo implement
	return c.registration
}
func (c *Client) GetPrivateKey() crypto.PrivateKey {
    // todo implement
	return c.key
}

func New(config config.Lego) (*Client, error) {
    client := Client{
		email: config.Email,
    }

	legoConfig := lego.NewConfig(&client)

    // Empty string means we want to use the production url
    if legoConfig.CADirURL != "" {
        legoConfig.CADirURL = legoConfig.CADirURL
    }
    // todo look into providing the newer smaller certificates as well should clients support it
	legoConfig.Certificate.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
    var err error
	client.legoClient, err = lego.NewClient(legoConfig)
	if err != nil {
        return nil, fmt.Errorf("error, when creating new lego client. Error: %v", err)
	}

    // todo check if registration is required

	reg, err := client.legoClient.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
        return nil, fmt.Errorf("error, when registering new lego user. Error: %v", err)
	}
	client.registration = reg
    return nil, nil
}

func (c *Client) Obtain() error {
    request := certificate.ObtainRequest{
		Domains: []string{"mydomain.com"},
		Bundle:  true,
	}
	_, err := c.legoClient.Certificate.Obtain(request)
	if err != nil {
        return fmt.Errorf("error, when requesting new certificate. Error: %v", err)
	}
    return nil
}
