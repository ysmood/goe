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

We can use the [whisper](https://github.com/ysmood/whisper) tool to encrypt the `.env` file and commit it to the version control system.
Then only selected team members can access it.

First create a whisper batch config file `whisper.json` beside the `.env` file:

```json
{
  "files": {
    ".env": ["@alice", "@bob"]
  }
}
```

Here `alice` and `bob` are the github user ids.

Then run the command to encrypt the `.env` file:

```bash
go run github.com/ysmood/whisper@latest -be whisper.json
```

Then you can safely commit generated `.env.wsp` file to the version control system.

To decrypt the `.env` file:

```bash
go run github.com/ysmood/whisper@latest -bd whisper.json
```

If you have several teams to manage, you might want to create a dedicated git repo like `https://github.com/your-org/vault` to hold all the `.env` files for different services.
Then you can use command like this to decrypt the `web.env.wsp` file to your local:

```bash
go run github.com/ysmood/whisper@latest -d https://github.com/your-org/vault/blob/main/web.env.wsp > .env
```

## Video Demo

[Video Link](https://youtu.be/vDTpzN9B4Nc)
