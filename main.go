package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"golang.org/x/crypto/ssh/terminal"
)

const mfaCodePrompt = "Enter MFA code: "

type options struct {
	profile string
}

func main() {
	log.SetFlags(0)
	var o options
	flag.StringVar(&o.profile, "profile", "", "Use a specific profile from your credential file.")
	flag.Parse()
	if err := run(context.Background(), o); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, o options) error {
	credentials, err := getCredentials(ctx, o.profile)
	if err != nil {
		return fmt.Errorf("cannot get credentials: %w", err)
	}
	log.Printf("you got a valid token until %s", credentials.Expires)
	fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", strconv.Quote(credentials.AccessKeyID))
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", strconv.Quote(credentials.SecretAccessKey))
	fmt.Printf("export AWS_SESSION_TOKEN=%s\n", strconv.Quote(credentials.SessionToken))
	return nil
}

func getCredentials(ctx context.Context, profile string) (*aws.Credentials, error) {
	cfg, err := external.LoadDefaultAWSConfig(
		external.WithSharedConfigProfile(profile),
		external.WithRegion(endpoints.UsEast1RegionID),
		external.WithMFATokenFunc(readMFACode),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot load AWS config: %w", err)
	}
	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve credentials: %w", err)
	}
	return &credentials, nil
}

func readMFACode() (string, error) {
	if _, err := fmt.Fprint(os.Stderr, mfaCodePrompt); err != nil {
		return "", fmt.Errorf("cannot write to stderr: %w", err)
	}
	b, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("cannot read from stdin: %w", err)
	}
	if _, err := fmt.Fprintln(os.Stderr); err != nil {
		return "", fmt.Errorf("cannot write to stderr: %w", err)
	}
	return string(b), nil
}
