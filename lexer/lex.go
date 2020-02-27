package lexer

import "strconv"

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
        case []byte:
            return string(v)
        case int:
            return strconv.Itoa(int(v))
        default:
            return "Unknown token"
    }
}

// TODO: return multiple errors
func Lex(in []byte) ([]Token, []error) {
    bp := bufpos{in, 0}
    var tokens []Token
    var errors []error
    for {
        tok, err := readToken(&bp);
        if err == EOF {
            break
        }
        if err != nil {
            errors = append(errors, err)
            if _, ok := err.(*InvalidTokenError); ok {
                break
            }

        }
        tokens = append(tokens, tok)
    }

    return tokens, errors
}

// TODO: allow returning multiple errors
func readToken(bp *bufpos) (Token, error) {
    readWhile(bp, []byte(" \n\r\t")) // Cut whitespace
    char, eof := bp.peek()
    if eof {
        return Token{}, EOF
    }

    switch(char) {
        case '+', '-', '*':
            bp.readByte()
            return Token{OperatorToken, bp.pos-1, 1, bp.buf[bp.pos-1]}, nil
        case '0','1','2','3','4','5','6','7','8','9':
            return readNumber(bp)
        case '"':
            return readString(bp)
    }

    // TODO: fix error handling
    return Token{}, (*InvalidTokenError)(bp)
}

func readByteWithEscapeCode(bp *bufpos) (byte, error) {
    // TODO: octal + hex support
    char, eof := bp.readByte()
    if eof {
        panic("unexpected EOF")
    }
    if char == '\\' {
        var val byte
        char, eof = bp.readByte()
        if eof {
            // TODO: proper unexpected eof error
            panic("unexpected eof")
        }
        switch(char) {
            case '\\':
                val = '\\'
            case 'n':
                val = '\n'
            case 't':
                val = '\t'
            case 'e':
                val = '\x1b'
            case 'r':
                val = '\r'
            case '"':
                val = '"'
            default:
                // TODO: recovery to skip to try to read rest of token
                return '?', &InvalidEscapeCodeError{bp.copyAt(bp.pos-2), bp.buf[bp.pos-2:bp.pos]}
            }
            return val, nil

    } else {
        return char, nil
    }
}

func readString(bp *bufpos) (Token, error) {
    // TODO: bounds checking for each ++
    var strbuf []byte
    initialPos := bp.pos
    bp.readByte() // Skip quote
    for bp.len() > 0 {
        if c, _ := bp.peek(); c == '"' {
            bp.readByte()
            return Token{StringToken, initialPos, len(strbuf), strbuf}, nil
        }

        char, err := readByteWithEscapeCode(bp)
        if err != nil {
            return Token{}, err
        }

        strbuf = append(strbuf, char)
    }
    return Token{}, (*NonTerminatedStringError)(bp.copyAt(initialPos))
}

// TODO: float support
func readNumber(bp *bufpos) (Token, error) {
    initialPos := bp.pos
    num := readWhile(bp, []byte("1234567890abcdefx"))
    i, err := strconv.ParseInt(string(num), 0, 0)
    if err != nil {
        return Token{}, &NumberSyntaxError{bp.copyAt(initialPos), err.(*strconv.NumError)}
    } else {
        return Token{NumberToken, initialPos, bp.pos - initialPos, int(i)}, err
    }
}
