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

About the format of `.env`: [link](https://pkg.go.dev/github.com/compose-spec/compose-go@v1.20.2/dotenv)

## CLI tool

Go to the [release page](https://github.com/ysmood/goe/releases) to download the CLI binary.

Usage examples:

```bash
# By default it will use .env file in current working directory to start a new shell.
goe

# Load file file/path/.env.dev as dotenv file.
goe file/path/.env.dev

# If there are arguments after the dotenv, they will be executed without starting a new shell.
goe .env.dev node app.js
```

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

### Decrypt .env int Github Actions

Generate a key pair for the CI:

```bash
whisper -gen-key id_ci
```

It will generate two files `id_ci` and `id_ci.pub`. When we encrypt `.env` we add the `id_ci.pub` as a recipient.

Add the `id_ci` to the [Github Action secret](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions).

Assume the Github secret name is `ID_CI`, then we can decrypt the file in the Github Action:

```yaml
name: Test

on: [push]

env:
  WHISPER_DEFAULT_KEY: ${{ secrets.ID_CI }}

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/setup-go@v4
        with:
          go-version: 1.21

      - name: setup .env
        run: go run github.com/ysmood/whisper@latest https://github.com/your-org/vault/blob/main/web.env.wsp > .env

      - name: test
        run: go test ./...
```

Usually you will use whisper to encrypt the `id_ci` file to the team's vault repo.
