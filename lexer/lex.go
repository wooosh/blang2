package lexer

import (
    "bytes"
    "strconv"
)

type Token struct {
    TypeOf TokenType
    Pos int // Starting position in file
    Length int
    Value interface{}
}

func (t Token) String() string {
    switch v := t.Value.(type) {
        case byte:
            return string(v)
        case int64:
            return strconv.Itoa(int(v))
        default:
            return "Unknown token"
    }
}

type bufpos struct {
    buf []byte
    pos int
}

// TODO: return multiple errors
func Lex(in []byte) ([]Token, error) {
    bp := bufpos{in, 0}
    var tokens []Token
    for {
        tok, err := readToken(&bp);
        if err == EOF {
            break
        }
        if err != nil {
            return tokens, err
        }
        tokens = append(tokens, tok)
    }

    return tokens, nil
}

func readToken(bp *bufpos) (Token, error) {
    eatWhitespace(bp)
    if len(bp.buf) - bp.pos == 0 {
        return Token{}, EOF
    }

    switch(bp.buf[bp.pos]) {
        case '+', '-', '*':
            bp.pos++
            return Token{OperatorToken, bp.pos-1, 1, bp.buf[bp.pos-1]}, nil
        case '0','1','2','3','4','5','6','7','8','9':
            return readNumber(bp)
    }

    return Token{}, (*InvalidTokenError)(bp)
}

// TODO: float support
func readNumber(bp *bufpos) (Token, error) {
    initialPos := bp.pos
    num := readAny(bp, []byte("1234567890abcdefx"))
    i, err := strconv.ParseInt(string(num), 0, 0)
    if err != nil {
        return Token{}, &NumberSyntaxError{bp, err.(*strconv.NumError)}
    } else {
        return Token{NumberToken, initialPos, bp.pos - initialPos, int(i)}, err
    }
}


// Helper functions
func readAny(bp *bufpos, chars []byte) []byte {
    initialPos := bp.pos
    for bp.len() > 0 && bytes.ContainsRune(chars, rune(bp.buf[bp.pos])) {
        bp.pos++
    }
    return bp.buf[initialPos:bp.pos]
}

// TODO: replace with readAny
func eatWhitespace(bp *bufpos) {
    for bp.len() > 0 && bytes.ContainsRune([]byte(" \n\r\t"), rune(bp.buf[bp.pos])) {
        bp.pos++
    }
}

func (b *bufpos) len() int {
    return len(b.buf) - b.pos
}



