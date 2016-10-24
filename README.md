# dutyme [![Build Status](http://img.shields.io/travis/tcnksm/dutyme.svg?style=flat-square)][travis] [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[travis]: https://travis-ci.org/tcnksm/dutyme
[license]: https://github.com/tcnksm/dutyme/blob/master/LICENSE

> You should receive alerts while operation on production. Take responsibility for it.

`dutyme` assigns PagerDuty on-call to you while you do operation or deployment.

## Why?

Normally, the on-call persion is fixed for a certain time (e.g., 1 week or 2 weeks) and rotated after that period. This works fine but has some issues. In working time, we deploy and operate a lot of time even though we are not the on-call. And sometimes that operation triggers alerts (because of components or because of type of job). This means no matter who the operator is, alerts are sent to the on-call persion.

Who adds changes or does something can fix issue fast (because he/she knows better about that). So when alerts are fired, its operator should receive alerts and handle it first. In addition to that, even in on-call, we nomarlly focus on our own task if no incident happens. If we receive alerts, we are disturbed. I want to avoid to disturb the on-call person by my operation and be disturbed by someones operation, too. That's why I made this tool.

## DEMO

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

*NOTE*: `dutyme` uses [override](https://support.pagerduty.com/hc/en-us/articles/202830170-Creating-and-Deleting-Overrides), which allows you to make one-time adjustments to on-call schedules (It doesn't modify the existing schedules). 


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
