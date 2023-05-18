# goma

Control SESAME SmartLock by SESAME v2 API

## Installation

As a library

```shell
go get github.com/ouest/goma
```

or if you want to use it as a bin command

go >= 1.17
```shell
go install github.com/ouest/goma/cmd/goma@latest
```

go < 1.17
```shell
go get github.com/ouest/goma/cmd/goma
```

## Usage

Add your application configuration to your `.env` file in the root of your project:

```shell
SESAME_UUID=YOUR_SESAME_UUID
SESAME_API_KEY=YOUR_SESAME_API_KEY
SESAME_SECRET_KEY=YOUR_SERSAME_SECRET_KEY
```

Then in your Go app you can do something like

```go
package main

import (
    "github.com/ouest/goma"
)

func main() {
    // get goma state
    state := goma.State()
    // get goma history
    history := goma.History(0, 10)
    // lock goma
    goma.Lock("account_name")
    // toggle goma
    goma.Toggle("account_name")
    // unlock goma
    goma.Unlock("account_name")
}
```

### Command Mode

Assuming you've installed the command as above and you've got `$GOPATH/bin` in your `$PATH`

```
goma state
goma history -p 0 -n 10
goma lock account_name
goma toggle account_name
goma unlock account_name
```

## License

See [LICENSE](LICENSE) © [ouest](https://github.com/ouest/)
