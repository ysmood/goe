package main

var USAGE = `
  Usage:

	goe [dotenv-file-path]
	goe -h | --help | -help


  Examples:

	# By default it will use .env file in current working directory to start a new shell.
	goe

	# Load file file/path/.env.dev as dotenv file.
	goe file/path/.env.dev

	# If there are arguments after the dotenv, they will be executed without starting a new shell.
	goe .env.dev node app.js


`
