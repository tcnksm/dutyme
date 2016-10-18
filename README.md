# dutyme

> You should receive alerts while you're doing something on production. Take responsibility for it.

`dutyme` assigns PagerDuty on-call to yourself temporarily while operation. It creates [override layer](https://support.pagerduty.com/hc/en-us/articles/202830170-Creating-and-Deleting-Overrides) on the existing schedule. 

*NOTE*: `dutyme` is still under development. command interface maybe updated in future.

## Usage

To assign, use `start` command,

```bash
$ dutyme start
```

It asks all necessary infomation to override (your PagerDuty email address or schedule name to override) 
and creates override layer. After executing, all infomation will be saved on disk so you can skip input from next time.

By default, it overrides 1 hour. You can change it via `-working` flag.

```bash
$ dutyme start -working 30m
```

## Token

To use dutyme command, you need a PagerDuty API v2 token. The token must have full access to read, write, update, and delete. Only account administrators have the ability to generate token. (See more about token on official doc https://goo.gl/VPvlwB)

## Install

To install, use `go get`:

```bash
$ go get github.com/tcnksm/dutyme
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
