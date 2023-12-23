# Overview

It provide some common helpers to load environment variables.

- No struct tag magic, only composable functions, better IDE autocomplete.
- Use default value as return type (with [go generics](https://go.dev/blog/intro-generics)).
- Recursively search for the `.env` file, easier sub package unit testing.
- Auto [expand](https://pkg.go.dev/os#Expand) the environment variables in the `.env` file.
- Auto parse base64 encoded value.
- Auto read file content if the value is a file path.
- Customizable parser.

For usage check the [example](example/basic.go).

About the format of `.env`: [link](https://pkg.go.dev/github.com/hashicorp/go-envparse)

## Video Demo

[demo](https://github.com/ysmood/goe/assets/1415488/b72cdfad-7123-4179-b2c3-839b7efc58e0)
