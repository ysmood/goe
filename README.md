# Overview

It provide some common helpers to load environment variables.

- No struct tag magic, only composable functions, better IDE autocomplete.
- `.env` file encryption and decryption.
- Use default value as return type (with [go generics](https://go.dev/blog/intro-generics)).
- Recursively search for the `.env` file, easier sub package unit testing.
- Auto [expand](https://pkg.go.dev/os#Expand) the environment variables in the `.env` file.
- Auto parse base64 encoded value.
- Auto read file content if the value is a file path.
- Customizable parser.
- Composable functions, easy to extend.

For usage check the [example](example/basic.go).

About the format of `.env`: [link](https://pkg.go.dev/github.com/hashicorp/go-envparse)

## Safely share .env file with team members

`.env` file usually contains sensitive information, like database password, API key, etc.
It's not recommended to commit it to the version control system, but it's usually required to run the project.
This package provides a simple way to encrypt and decrypt the `.env` file, so only selected team members can access it.

You need to add a `GOE_ENV_VIEWERS` env variable to the `.env` file, it's a comma separated list of addresses.
Each address can be a github user id, or a https public key url ([syntax details](https://github.com/ysmood/whisper)). For example:

```bash
GOE_ENV_VIEWERS="@ysmood,@https://test.com/jack.pub"
```

Then encrypt the `.env` file:

```bash
go run github.com/ysmood/goe/encrypt@latest
```

Then you can safely commit the generated `.env.goe` file to the version control system.

When a team member clones the project, they can directly use the `.env.goe` file without any extra steps,
they can just run the code like the [example](example/basic.go).

## Video Demo

[Video Link](https://youtu.be/vDTpzN9B4Nc)
