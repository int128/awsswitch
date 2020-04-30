# awsswitch [![CircleCI](https://circleci.com/gh/int128/awsswitch.svg?style=shield)](https://circleci.com/gh/int128/awsswitch)

This is a command to export the following credentials variables to switch a role with MFA.

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_SESSION_TOKEN`

Key features:

- Single binary
- Interoperable config with AWS CLI (`~/.aws/config`)
- Interoperable token cache with AWS CLI (`~/.aws/cli/cache`)


## Getting Started

Install awsswitch.

```sh
# Homebrew (macOS)
brew install int128/awsswitch/awsswitch

# GitHub Releases
curl -LO https://github.com/int128/awsswitch/releases/download/v0.2.0/awsswitch_linux_amd64.zip
unzip awsswitch_linux_amd64.zip

# Go
go get github.com/int128/awsswitch
```

Set up your `.aws/config` to switch a role. For example,

```ini
[profile USERNAME]

[profile USERNAME_administrator]
mfa_serial = arn:aws:iam::1234567890:mfa/USERNAME
role_arn = arn:aws:iam::1234567890:role/AdministratorMFA
source_profile = USERNAME
duration_seconds = 43200
```

Run the command in your terminal.

```console
% eval $(awsswitch --profile=USERNAME_administrator)
Enter MFA code:
you got a valid token until 2020-04-19 21:43:38 +0000 UTC
```

It will export the credentials variables.
Now you can run tools such as AWS CLI and Terraform.

```console
% aws s3 ls

% terraform apply
```

It attempts to read the token cache in `~/.aws/cli/cache`.
You do not need to enter a MFA code if the token is valid.
This behavior is interoperable with AWS CLI.


## How to set up the role switching

### 1. Set up a role

Create an IAM role to switch to.
You need to set up a trusted relationship to an AWS account or IAM user.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::1234567890:root"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "Bool": {
          "aws:MultiFactorAuthPresent": "true"
        }
      }
    }
  ]
}
```

### 2. Set up a user

Create an IAM user.

You need to set up an assume role.
See [document](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use_permissions-to-switch.html) for more.

```json
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "sts:AssumeRole",
    "Resource": "arn:aws:iam::1234567890:role/AdministratorMFA"
  }
}
```

Set your credentials to `~/.aws/credentials`.

```console
% aws configure --profile=USERNAME
```

Add a profile to `.aws/config` to switch to the role.

```ini
[profile USERNAME]

[profile USERNAME_administrator]
mfa_serial = arn:aws:iam::1234567890:mfa/USERNAME
role_arn = arn:aws:iam::1234567890:role/AdministratorMFA
source_profile = USERNAME
duration_seconds = 43200
```


## Contributions

This is an open source software. Feel free to open issues and pull requests.
