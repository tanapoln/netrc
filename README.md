# netrc
A simple command line tool for managing your netrc file.

## Install

```sh
$ go get github.com/naaman/netrc
```

## Usage

List netrc entries:

```sh
$ netrc list
api.heroku.com 1234...
github.com 1234...
```

Show a password for a machine entry:

```sh
$ netrc password api.heroku.com
1234...
```
