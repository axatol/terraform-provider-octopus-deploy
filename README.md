# Octopus Deploy Terraform Provider

Provider for Octopus Deploy, my little project to learn about Terraform providers.

Original framework readme [here](./framework-readme.md).

## Developing

You can build the provider and configure Terraform to prefer your local version

```bash
# Build and install the provider
go install

# Find where go put the provider binary
GOBIN="$(go env GOPATH)/bin"

# Add the override
cat <<EOF > $HOME/.terraformrc
provider_installation {
  dev_overrides {
    "registry.terraform.io/axatol/octopusdeploycontrib" = "${GOBIN}"
  }

  direct {}
}
EOF
```

At this point, you can use the provider like so:

```terraform
terraform {
  required_providers {
    octopusdeploycontrib = {
      source = "registry.terraform.io/axatol/octopusdeploycontrib"
    }
  }
}

provider "octopusdeploycontrib" {}

data "octopusdeploycontrib_project" "test" {}
```
