package main

import (
    "fmt"
    "blang/lexer"
)


func main() {
    toks, err := lexer.Lex([]byte("1f0+0xf -\n 0b0110"))
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(toks)
    }
}
