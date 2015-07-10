# YahooBoss for Go

YahooBoss is a [Go](http://golang.org) package for [Yahoo! Search BOSS](https://boss.yahoo.com/)

The package is no where near finished but it is a start

## Installation
    go get github.com/ThunApps/YahooBoss

## License

YahooBoss is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

## Documentation

```go
package main

import (
  "github.com/ThunApps/yahooboss"
)

func main() {
  bs := yahooboss.BossSearch{<token>,
                        <secret>,
                        "web"}

  bs.Search("Hello+World!")
}
```
