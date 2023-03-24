# TeeTimer

TeeTimer is a Go program designed to automate the process of reserving tee times at golf courses that use [ClubHouse Online](https://clubhouseonline-e3.com/) to manage their tee sheet. It reads your preferred tee times from a configuration file and reserves them as soon as they become available.

## Usage

```console
$ make build
$ ./teetimer -h
NAME:
   teetimer - Automate tee time reservations for clubhouselineline-e3.net golf courses

USAGE:
   teetimer [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --url value                 Tee sheet URL [$CLUBHOUSE_URL]
   --username value, -u value  Username used to log into the application [$CLUBHOUSE_USERNAME]
   --password value, -p value  Password used to log into the application [$CLUBHOUSE_PASSWORD]
   --help, -h                  show help
```
