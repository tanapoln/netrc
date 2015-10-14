# netrc
A simple command line tool for managing your netrc file.

## Install

```sh
$ go get github.com/naaman/netrc
```

## Usage

List netrc entries:

```sh
$ netrc
api.heroku.com
github.com
```

Show logins:

```sh
$ netrc -l
api.heroku.com user@heroku.com
github.com user@github.com
```

Show passwords:

```sh
$ netrc -p
api.heroku.com 1234...
github.com 1234...
```

Show a password for a machine entry:

```sh
$ netrc -p -n api.heroku.com
1234...
```
