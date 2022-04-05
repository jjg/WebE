# fzx

Sucessor of [jsfs](https://github.com/jjg/jsfs).  Part of the WebE project.

I'm not going to put much here for now. [The journal](../journal.md) is the best source of info at the moment.

## Usage

### Automated Tests

`go test ./...`

### Run it!

1. `go build`
2. `./fzx`

### Gokrazy

> Incomplete, currently just notes to myself.

`go get github.com/jjg/WebE/fzx`
`gokr-packer -overwrite=/dev/mmcblk0 -serial_console=disabled github.com/gokrazy/mkfs github.com/jjg/WebE/fzx`


## API

The API is based off [jsfs](https://github.com/jjg/jsfs#api).  It is nowhere near complete (again, see [the journal](../journal.md) for current status).

### Basic examples

* Upload a file: `curl -v --data-binary @fzx "http://localhost:7302/testing/fzx"`
* Get file info: `curl -v -I "http://localhost:7302/testing/fzx"`
* Download a file: `curl -v -o fzx3 "http://localhost:7302/testing/fzx"`

## TODO

* ~~Implement basic `HEAD`~~
* ~~Implement basic `GET`~~
* ~~Implement basic `POST`~~
* Implement basic `PUT`
* Implement basic `DELETE`
* Implement basic `EXECUTE`
* Implement complete responses (all headers, etc.)
* Implement authorization
* ~~Implement configuration (flags, config file, etc.)~~
* Implement compatibility mode (read/write legacy JSFS storage pools)
* Implement federation
* Refactor for performance
* Bugs
  + ~~`HEAD` doesn't close the connection when returning 404 for some reason?~~
  + Fix broken tests
