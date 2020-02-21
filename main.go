package main

import (
    "fmt"
    "blang/lexer"
)


func main() {
    toks, errors := lexer.Lex([]byte(`0xx "312`))
    if errors != nil {
        for _, err := range errors {
            fmt.Println(err)
        }
    } else {
        fmt.Println(toks)
    }
}
