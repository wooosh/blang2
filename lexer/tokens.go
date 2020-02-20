package lexer

type TokenType uint
const (
    NumberToken TokenType = iota
    OperatorToken // TODO: replace operators with identifiers
    StringToken
)
