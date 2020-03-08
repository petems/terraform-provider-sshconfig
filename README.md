# terraform-provider-sshconfig

[![Build Status](https://travis-ci.org/petems/terraform-provider-sshconfig.svg?branch=master)](https://travis-ci.org/petems/terraform-provider-sshconfig)

Terraform provider for reading and configurating SSH Config files

## Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)

## Installing the Provider
Follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing it into your plugins directory, run `terraform init` to initialize it.

## Usage

```hcl
data "sshconfig_host" "github_ssh_config" {
  path = pathexpand("~/.ssh/config")
  host = "github.com"
}

output "github_config" {
  value = "${data.sshconfig_host.github_ssh_config.rendered}"
}

output "host_map" {
  value = "${data.sshconfig_host.github_ssh_config.host_map}"
}

```

Gives the result:

```shell
$ terraform apply
data.sshconfig_host.github_ssh_config: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

github_config =
Host github.com
  ControlMaster auto
  ControlPath ~/.ssh/ssh-%r@%h:%p
  ControlPersist yes
  User git
  IdentityFile ~/.ssh/github_ssh_key_10_Mar_2020

host_map = {
  "ControlMaster" = "auto"
  "ControlPath" = "~/.ssh/ssh-%r@%h:%p"
  "ControlPersist" = "yes"
  "IdentityFile" = "~/.ssh/github_ssh_key_10_Mar_2020"
  "User" = "git"
}
```

Examples are under [/examples](/examples).

## Building the Provider
Clone and build the repository

```sh
go get github.com/petems/terraform-provider-sshconfig
make build
```

Symlink the binary to your terraform plugins directory:

```sh
mkdir -p ~/.terraform.d/plugins/
ln -s $GOPATH/bin/terraform-provider-sshconfig ~/.terraform.d/plugins/
```

## Updating the Provider

```sh
go get -u github.com/petems/terraform-provider-sshconfig
make build
```

## TODO

* Add configuration of the consensus timing (ie. how long it will wait to resolve)
* Add option of getting ipv6 or ipv4 ipaddress

## Contributing
* Install project dependencies: `go get github.com/kardianos/govendor`
* Run tests: `make test`
* Build the binary: `make build`