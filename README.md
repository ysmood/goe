# Overview

It provide some common helpers to load environment variables.

- Loader function that uses generics so you don't need to cast the value manually.
- Load helper that recursively search for the `.env` file in parent folders until it finds one.
- Auto [expand](https://pkg.go.dev/os#Expand) the environment variables in the `.env` file.

For usage check the [example](example/basic.go).

About the format of `.env`: [link](https://pkg.go.dev/github.com/hashicorp/go-envparse)

## Video Demo

[demo](https://github.com/ysmood/goe/assets/1415488/b72cdfad-7123-4179-b2c3-839b7efc58e0)
