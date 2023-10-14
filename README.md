# Octopus Deploy Terraform Provider

Provider for Octopus Deploy, my little project to learn about Terraform providers.

Original framework readme [here](./framework-readme.md).

## Developing

Setting up:

```bash
# Build and install the provider
go install

# Find where go put the binary
GOBIN="$(go env GOPATH)/bin"

# Add the override, replacing <path> with GOBIN
cat <<EOF > $HOME/.terraformrc
provider_installation {
  dev_overrides {
    "registry.terraform.io/axatol/octopusdeploy" = "<path>"
  }

  direct {}
}
EOF
```
