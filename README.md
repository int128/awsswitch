# awsswitch

This is a command to export the credentials variables to switch a role with MFA.


## Getting Started

Install awsswitch.

```sh
# Go
go get github.com/int128/awsswitch
```

### Set up a role to switch

Create an IAM role.
You need to set up a trusted relationship.

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

### Set up a user

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

Run awsswitch command in your terminal.

```console
% $(awsswitch --profile=USERNAME_administrator)
Enter MFA code:
you got a valid token until 2020-04-19 21:43:38 +0000 UTC
```

It will export `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_SESSION_TOKEN`.

Now you can run tools such as AWS CLI, Terraform and Ansible.

```console
% terraform apply
```


## Contributions

This is an open source software. Feel free to open issues and pull requests.
