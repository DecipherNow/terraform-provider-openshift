# Terraform Provider - OpenShift
This repository contains the source code for the Terraform OpenShift provider.

## Development
The following instructions provide general guidelines for developers working on this provider.

### Prerequisites
The following must be available on a development or build machine.

- Go 1.11.4
- Terraform 0.10+

### Recommendations
We recommend using [goenv](https://github.com/syndbg/goenv) to manage your installed Go versions and thus we include a `.go-version` file to support automatically selecting this version when installed with goenv. See the goenv repository for additional information on how this works.

This repository also uses Go modules for dependency management. Please read up on how Go modules work [here](), but we recommend cloning this repository outside of your `${GOPATH}` and turning Go modules on with `export GO111MODULE=on`.
