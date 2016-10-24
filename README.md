# dutyme [![Build Status](http://img.shields.io/travis/tcnksm/dutyme.svg?style=flat-square)][travis] [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[travis]: https://travis-ci.org/tcnksm/dutyme
[license]: https://github.com/tcnksm/dutyme/blob/master/LICENSE

> You should receive alerts while operation on production. Take responsibility for it.

`dutyme` assigns PagerDuty on-call to yourself while operation. It creates [override layer](https://support.pagerduty.com/hc/en-us/articles/202830170-Creating-and-Deleting-Overrides) on the existing schedule. *NOTE*: `dutyme` is still under development. command interface maybe updated in future.

![](/doc/dutyme.gif)

## Requirement

To use dutyme command, you need a PagerDuty API v2 token. 

The token must have full access to read, write, update, and delete. Only account administrators have the ability to generate token (See more about token on official doc https://goo.gl/VPvlwB).

## Usage

To assign, use `start` command,

```bash
$ dutyme start
```

It asks all necessary infomation to override (your PagerDuty email address or schedule name) and creates a override layer. You can create multiple overrides on the same term (the latest one has priority). After executing, all infomation will be saved on disk so you can skip input from next time. By default, it overrides 1 hour. You can change it via `-working` flag. See more usage by `-help` flag.

## Install

To install, you can use `go get` or `brew`:

```bash
$ brew tap tcnksm/dutyme
$ brew install dutyme
```

## Contribution

1. Fork ([https://github.com/tcnksm/dutyme/fork](https://github.com/tcnksm/dutyme/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[Taichi Nakashima](https://github.com/tcnksm)
