package main

import (
    "bytes"
    "fmt"
    "github.com/divan/gorilla-xmlrpc/xml"
    "log"
    "net/http"
)

func XmlRpcCall(method string, args struct{ Who string }) (reply struct{ Message string }, err error) {
    arg := struct {
        A struct {
            A string
            B string
        }
    }{}

    buf, _ := xml.EncodeClientRequest(method, &arg)

    fmt.Printf("%s\n", buf)

    resp, err := http.Post("http://localhost:1234/RPC2", "text/xml", bytes.NewBuffer(buf))
    if err != nil {
        return
    }
    defer resp.Body.Close()

    err = xml.DecodeClientResponse(resp.Body, &reply)
    return
}

func main() {
    reply, err := XmlRpcCall("HelloService.Say", struct{ Who string }{"User 1"})
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Response: %s\n", reply.Message)
}
