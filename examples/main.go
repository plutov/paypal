package main

import (
    "paypalsdk"

    "fmt"
    "os"
)

func main() {
    client, err := paypalsdk.NewClient("123", "123", paypalsdk.APIBaseSandBox)
    if err == nil {
        fmt.Println("DEBUG: ClientID=" + client.ClientID + " APIBase=" + client.APIBase)
    } else {
        fmt.Println("ERROR: " + err.Error())
        os.Exit(1)
    }
}
