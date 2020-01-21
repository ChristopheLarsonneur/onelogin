# OneLogin API written in GO

## How to use it

```go
package main

import (
    "github.com/clarsonneur/onelogin"
    "github.com/op/go-logging"
)

func main() {
    ol := onelogin.NewService("eu", string(clientID), string(clientSecret), "myCompany", logging.INFO)
    olAPI, err := ol.GetAPI() // Get API object with access token.

    if err != nil {
        log.Fatalf("%s", err)
    }

    user := api.NewGetUserByID()
    user.Get(olAPI, 12345)

    if roles, err := ol.GetRoles() ; err != nil {
        log.Fatalf("%s", err)
    } else {
        ...
    }
    ...
}
```
