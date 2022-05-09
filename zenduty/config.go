package zenduty

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zenduty/zenduty-go-sdk/client"
)

type Config struct {
	Token   string
	BaseURL string
}

const invalidCreds = `
No valid credentials found for zenduty provider.
`

func (c *Config) Client() (*client.Client, error) {
	if c.Token == "" {
		return nil, fmt.Errorf(invalidCreds)
	}

	var httpClient *http.Client
	httpClient = http.DefaultClient

	config := &client.Config{
		BaseURL:    c.BaseURL,
		HTTPClient: httpClient,
		Token:      c.Token,
	}

	client, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] zenduty client configured")

	return client, nil
}
