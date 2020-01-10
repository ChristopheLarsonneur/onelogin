# OneLogin API written in GO

## How to use it

```go
package main

import (
    "github.com/clarsonneur/onelogin"
    "github.com/op/go-logging"
)

func main() {
    ol := onelogin.NewAPI("eu", string(clientID), string(clientSecret), "myCompany", logging.INFO)
    err := ol.ObtainAPIAccess() // required to be authorized to call the API.

    
}
```
