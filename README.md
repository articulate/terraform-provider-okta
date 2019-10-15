# Terraform Provider Okta

- [![Build Status](https://travis-ci.org/terraform-providers/terraform-provider-okta.svg?branch=master)](https://travis-ci.org/terraform-providers/terraform-provider-okta)
- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

## Usage

This plugin requires two inputs to run: the okta organization name and the okta api token. The okta base url is not required and will default to "okta.com" if left out.

You can specify the inputs in your tf plan:

```
provider "okta" {
  org_name  = <okta instance name, e.g. dev-XXXXXX>
  api_token = <okta instance api token with the Administrator role>
  base_url  = <okta base url, e.g. oktapreview.com>

  // Optional settings, https://en.wikipedia.org/wiki/Exponential_backoff
  max_retries      = <number of retries on api calls, default: 5>
  backoff          = <enable exponential backoff strategy for rate limits, default = true>
  min_wait_seconds = <min number of seconds to wait on backoff, default: 30>
  max_wait_seconds = <max number of seconds to wait on backoff, default: 300>
}
```

OR you can specify environment variables:

```
OKTA_ORG_NAME=<okta instance name, e.g. dev-XXXXXX>
OKTA_API_TOKEN=<okta instance api token with the Administrator role>
OKTA_BASE_URL=<okta base url, e.g. oktapreview.com>
```

## Examples

As we build out resources we build concomitant acceptance tests that require use to create resource config that actually creates and modifies real resources. We decided to put these test fixtures to good use and provide them [as examples here.](./examples)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-okta`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-okta
```

Enter the provider directory and build the provider. Ensure you have Go Modules enabled, depending on the version of Go you are using, you may have to flip it on with `GO111MODULE=on`.

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-okta
$ make build
```

## Using the provider

Example terraform plan:

```
provider "okta" {
  org_name  = "dev-XXXXX"
  api_token = "XXXXXXXXXXXXXXXXXXXXXXXX"
  base_url  = "oktapreview.com"
}

resource "okta_user" "blah" {
  first_name = "blah"
  last_name  = "blergh"
  email      = "XXXXX@XXXXXXXX.XXX"
  login      = "XXXXX@XXXXXXXX.XXX"
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-okta
...
```

In order to test the provider, you can simply run `make test`. The acceptance tests require an API token and a corresponding Okta org, if you use dotenv, you can `cp .env.sample .env` and add your Okta settings there, and prefix make test with `dotenv`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
