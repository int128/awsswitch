package tokencache

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type content struct {
	Credentials contentCredentials `json:"Credentials"`
}

type contentCredentials struct {
	AccessKeyID     string    `json:"AccessKeyId"`
	SecretAccessKey string    `json:"SecretAccessKey"`
	SessionToken    string    `json:"SessionToken"`
	Expiration      time.Time `json:"Expiration"`
}

func newFromAWSCredentials(credentials aws.Credentials) content {
	return content{
		Credentials: contentCredentials{
			AccessKeyID:     credentials.AccessKeyID,
			SecretAccessKey: credentials.SecretAccessKey,
			SessionToken:    credentials.SessionToken,
			Expiration:      credentials.Expires,
		},
	}
}

func (c *content) awsCredentials() aws.Credentials {
	return aws.Credentials{
		AccessKeyID:     c.Credentials.AccessKeyID,
		SecretAccessKey: c.Credentials.SecretAccessKey,
		SessionToken:    c.Credentials.SessionToken,
		Expires:         c.Credentials.Expiration,
		CanExpire:       true,
	}
}

func decodeFrom(name string) (*content, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", name, err)
	}
	defer f.Close()
	var c content
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, fmt.Errorf("cannot decode json from %s: %w", name, err)
	}
	return &c, nil
}

func encodeTo(name string, c content) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", name, err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(&c); err != nil {
		return fmt.Errorf("cannot encode json to %s: %w", name, err)
	}
	return nil
}
