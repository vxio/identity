# Identity

[![GoDoc](https://godoc.org/github.com/moov-io/identity?status.svg)](https://godoc.org/github.com/moov-io/identity)
[![Build Status](https://travis-ci.com/moov-io/identity.svg?branch=master)](https://travis-ci.com/moov-io/identity)
[![Coverage Status](https://codecov.io/gh/moov-io/identity/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/identity)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/identity)](https://goreportcard.com/report/github.com/moov-io/identity)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/identity/master/LICENSE)

Moov Identity handles management for our users who will be managing authentication and authorization our systems.

## Project Status

This project is currently under development and could introduce breaking changes to reach a stable status. We are looking for community feedback so please try out our code or give us feedback!

## Getting Started

Identity is primarily a Go based HTTP server with unit tests and code fuzzing the code to help ensure our code is production ready for everyone. Identity uses [Go Modules](https://github.com/golang/go/wiki/Modules) to manage dependencies and suggests Go 1.14 or greater.

[API](api/identityapi.yml)

To clone our code and verify our tests on your system run:

```
$ git clone git@github.com:moov-io/identity.git
$ cd identity

$ go test ./...
ok   	github.com/moov-io/identity	0.710s	coverage: 98.1% of statements
```

## Getting Help

 channel | info
 ------- | -------
 [Project Documentation](https://docs.moov.io/) | Our project documentation available online.
 Google Group [moov-users](https://groups.google.com/forum/#!forum/moov-users)| The Moov users Google group is for contributors other people contributing to the Moov project. You can join them without a google account by sending an email to [moov-users+subscribe@googlegroups.com](mailto:moov-users+subscribe@googlegroups.com). After receiving the join-request message, you can simply reply to that to confirm the subscription.
Twitter [@moov_io](https://twitter.com/moov_io)	| You can follow Moov.IO's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channels to have an interactive discussion about the development of the project.

## Supported and Tested Platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## Contributing

Yes please! Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](https://github.com/moov-io/ach/blob/master/CODE_OF_CONDUCT.md) to get started! [Checkout our issues](https://github.com/moov-io/identity/issues) for something to help out with.

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) and uses Go 1.14 or higher. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/identity/releases/latest) as well. We highly recommend you use a tagged release for production.

## License

Apache License 2.0 See [LICENSE](LICENSE) for details.
