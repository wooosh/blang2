package lexer

type TokenType uint
const (
    NumberToken TokenType = iota
    OperatorToken
)
