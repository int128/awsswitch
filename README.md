# awsswitch

This is a command to export the credentials variables to switch a role with MFA code.


## Getting Started

Install awsswitch.

```sh
# Go
go get github.com/int128/awsswitch
```

### Set up a user

Create an IAM user on AWS management console.
You do not need to attach any IAM policy.

Set your credentials to `~/.aws/credentials`.

```console
% aws configure --profile=USERNAME
```

### Set up a role to switch

Create an IAM role on AWS management console.
You need to set up an assume role.
See [document](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use_permissions-to-switch.html) for more.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::1234567890:user/USERNAME"
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

Add a profile to `.aws/config` to switch a role.

```
[profile USERNAME]

[profile USERNAME_administrator]
mfa_serial = arn:aws:iam::1234567890:mfa/USERNAME
role_arn = arn:aws:iam::1234567890:role/AdministratorMFA
source_profile = USERNAME
duration_seconds = 43200
```

### Switch a role

You can export `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_SESSION_TOKEN`.

```console
% export AWS_PROFILE=YOURNAME
% eval $(awsswitch)
```


## Contributions

This is an open source software. Feel free to open issues and pull requests.
