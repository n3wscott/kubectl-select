# kubectl-select

A `kubectl` extension to select from local config via a TUI.

[![GoDoc](https://godoc.org/github.com/n3wscott/kubectl-select?status.svg)](https://godoc.org/github.com/n3wscott/kubectl-seect)
[![Go Report Card](https://goreportcard.com/badge/n3wscott/kubectl-select)](https://goreportcard.com/report/n3wscott/kubectl-select)


## Installation

`kubectl-select` can be installed via:

```shell
go get github.com/n3wscott/kubectl-select
```

`kubectl-select` is meant to be installed in a `PATH` location. 

To update your installation:

```shell
go get -u github.com/n3wscott/kubectl-select
```

## Usage

Use as a kubernetes extension, 

```shell
kubectl select
```

This will show a menu driven off the currently configured Kubernetes clients.
Select one by pressing `ENTER`. To cancel, `ESC` or `q`.

