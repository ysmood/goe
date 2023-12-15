# Overview

It provide some common helpers to load environment variables.

- Loader function that uses generics so you don't need to cast the value manually.
- Load helper that recursively search for the `.env` file in parent folders until it finds one.

For usage check the [example](example/basic.go).

About the format of `.env`: [link](https://pkg.go.dev/github.com/hashicorp/go-envparse)
