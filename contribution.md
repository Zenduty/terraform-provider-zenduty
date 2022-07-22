# Contributing

Welcome to terraform provide of zenduty.

# Steps

For issues, please raise a Pull request.

#

## Steps to build and test [local development]

```
go mod init
go fmt
go mod tidy
go build -o terraform-provider-zenduty
```

Copy the build `terraform-provider-zenduty` file to appropriate location for local testing. 

For example: 

## For Mac
```
mkdir -p ../localtesting/.terraform/plugins/zenduty.com/zenduty/zenduty/0.0.1/darwin_arm64
cp terraform-provider-zenduty ../localtesting/.terraform/plugins/zenduty.com/zenduty/zenduty/0.0.1/darwin_arm64
```

## For Linux
```
mkdir -p ../localtesting/.terraform/plugins/zenduty.com/zenduty/zenduty/0.0.1/linux_amd64
cp terraform-provider-zenduty ../localtesting/.terraform/plugins/zenduty.com/zenduty/0.0.1/linux_amd64
```

```
cd ../localtesting
touch versions.tf
touch terraform.tf
touch main.tf
```

Content of `versions.tf` file:
```
terraform {
  required_providers {
    zenduty = {
      version = "= 0.0.1"
      source  = "zenduty.com/zenduty/zenduty"
    }
  }
}
```
Content of `terraform.tf` file:
```

provider "zenduty" {
  # Configuration options
    token = "<your API token>"
}

```

Content of `main.tf` file:
```
data "zenduty_user" "me" {
    email = "deepak@abc.com"
}
output "user" {
  value = data.zenduty_user.me.users[0]
}
```

Run the terraform script like:

```
terraform init -plugin-dir=.terraform/plugins/
terraform plan
terraform apply
terraform destroy
```

Keep removing the `.terraform.lock.hcl` file to avoid hash conflicts.

#
Releases are planned according to the features and Blug fixes.

For api documentation of zenduty, please ref: https://apidocs.zenduty.com/