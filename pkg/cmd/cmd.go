// Package cmd provides the command line interface.
package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/int128/awsswitch/pkg/prompt"
	"github.com/int128/awsswitch/pkg/tokencache"
)

type options struct {
	profile string
}

func (o *options) register(f *flag.FlagSet) {
	f.StringVar(&o.profile, "profile", "", "Use a specific profile from your credential file.")
}

// Run parses the command line arguments and run the use-case.
func Run(ctx context.Context, osArgs []string) error {
	var o options
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	o.register(f)
	if err := f.Parse(osArgs[1:]); err != nil {
		return fmt.Errorf("invalid argument: %w", err)
	}

	credentials, err := getCredentials(ctx, o.profile)
	if err != nil {
		return fmt.Errorf("cannot get credentials: %w", err)
	}
	fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", strconv.Quote(credentials.AccessKeyID))
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", strconv.Quote(credentials.SecretAccessKey))
	fmt.Printf("export AWS_SESSION_TOKEN=%s\n", strconv.Quote(credentials.SessionToken))
	return nil
}

func getCredentials(ctx context.Context, profile string) (*aws.Credentials, error) {
	cfg, err := external.LoadDefaultAWSConfig(
		external.WithSharedConfigProfile(profile),
		external.WithRegion(endpoints.UsEast1RegionID),
		external.WithMFATokenFunc(func() (string, error) {
			return prompt.ReadPassword("Enter MFA code: ")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot load AWS config: %w", err)
	}
	for _, cs := range cfg.ConfigSources {
		if sc, ok := cs.(external.SharedConfig); ok {
			credentials, err := getCredentialsUsingSharedConfig(ctx, cfg, sc)
			if err != nil {
				log.Printf("skipped the config: %s", err)
				continue
			}
			return credentials, err
		}
	}
	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve credentials from aws: %w", err)
	}
	log.Printf("you got a valid token until %s", credentials.Expires)
	return &credentials, nil
}

func getCredentialsUsingSharedConfig(ctx context.Context, cfg aws.Config, sc external.SharedConfig) (*aws.Credentials, error) {
	cache, err := tokencache.Load(sc)
	if err == nil {
		if !cache.Expired() {
			log.Printf("you already have a valid token until %s", cache.Expires)
			return cache, nil
		}
		log.Printf("token cache has expired at %s", cache.Expires)
	}
	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve credentials from aws: %w", err)
	}
	if err := tokencache.Save(sc, credentials); err != nil {
		return nil, fmt.Errorf("cannot save cache: %w", err)
	}
	log.Printf("you got a valid token until %s", credentials.Expires)
	return &credentials, nil
}
