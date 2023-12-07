# Overview

Load the `.env` file from current directory or one of the parent directories, then use it to set the process environment variables.

It also provides a function to parse the value of environment variables, it uses generics so you don't need to cast the value manually.

For usage check the [example](example/basic.go).

About the format of `.env`: [link](https://docs.docker.com/compose/environment-variables/env-file/)
