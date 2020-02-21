package main

import (
    "fmt"
    "blang/lexer"
)


func main() {
    toks, errors := lexer.Lex([]byte(`0x32x324 0fff0x`))
    if errors != nil {
        for _, err := range errors {
            fmt.Println(err)
        }
    } else {
        fmt.Println(toks)
    }
}
