package lexer

type TokenType uint
const (
    // Values
    NumberToken TokenType = iota
    StringToken
    IdentifierToken

    // Symbols
    LBrace
    RBrace
    LParen
    RParen
    Comma

    // Keywords
    Fn
    If
    Else
    While
    Try
    Return
    True
    False
)
