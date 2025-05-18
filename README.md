# Identicon

This repository provides a Go implementation of GitHub's [identicon](https://en.wikipedia.org/wiki/Identicon) algorithm.

## Usage

```
import "path/to/identicon/go"

func main() {
    // Generate an identicon for the number 480938 and save it
    // to the file hubot.png
    if err := identicon.GenerateIdenticon(480938, "hubot.png"); err != nil {
        log.Fatal(err)
    }
}
```

## Development

```
go test ./...
```
