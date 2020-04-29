// Package tokencache provides access to the token cache in ~/.aws/cli/cache.
// This is interoperable with the token cache of AWS CLI.
package tokencache

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

// Load reads the token cache file and returns the credentials.
func Load(cfg external.SharedConfig) (*aws.Credentials, error) {
	name := computeFilename(cfg)
	cache, err := decodeFrom(name)
	if err != nil {
		return nil, fmt.Errorf("cannot load the credentials from cache: %w", err)
	}
	credentials := cache.awsCredentials()
	return &credentials, nil
}

// Save writes the credentials to the token cache file.
func Save(cfg external.SharedConfig, credentials aws.Credentials) error {
	name := computeFilename(cfg)
	if err := encodeTo(name, newFromAWSCredentials(credentials)); err != nil {
		return fmt.Errorf("cannot save the credentials to %s: %w", name, err)
	}
	return nil
}

func computeFilename(cfg external.SharedConfig) string {
	key := computeKey(cfg)
	dir := filepath.Join(filepath.Dir(external.DefaultSharedConfigFilename()), "cli", "cache")
	return filepath.Join(dir, key+".json")
}

// computeKey returns a key for the config.
// This is based on https://github.com/boto/botocore/blob/9f5afe26cc4a2695dcb62db453f9db7c299f3ffc/botocore/credentials.py#L714
func computeKey(cfg external.SharedConfig) string {
	var a []string
	// keys must be sorted
	if cfg.RoleDurationSeconds != nil {
		a = append(a, fmt.Sprintf(`"DurationSeconds": %d`, int(cfg.RoleDurationSeconds.Seconds())))
	}
	a = append(a, `"RoleArn": `+strconv.Quote(cfg.RoleARN))
	if cfg.RoleSessionName != "" {
		a = append(a, `"RoleSessionName": `+strconv.Quote(cfg.RoleSessionName))
	}
	a = append(a, `"SerialNumber": `+strconv.Quote(cfg.MFASerial))
	s := sha1.Sum([]byte(fmt.Sprintf("{%s}", strings.Join(a, ", "))))
	return hex.EncodeToString(s[:])
}
