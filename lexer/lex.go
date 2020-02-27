package lexer

import "strconv"

type Token struct {
    TypeOf TokenType
    Pos int // Starting position in file
    Length int
    Value interface{}
}

func (t Token) String() string {
    switch t.TypeOf {
        case Fn: return "fn"
        case If: return "if"
        case Else: return "else"
        case While: return "while"
        case Return: return "return"
        case Try: return "try"
        case True: return "true"
        case False: return "false"
        case LBrace: return "{"
        case RBrace: return "}"
        case LParen: return "("
        case RParen: return ")"
        case Comma: return ","
    }
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

    if bp.match("fn") {
        return Token{Fn, bp.pos-2, 2, nil}, nil
    } else if bp.match("return") {
        return Token{Return, bp.pos-6, 6, nil}, nil
    } else if bp.match("if") {
        return Token{If, bp.pos-2, 2, nil}, nil
    } else if bp.match("else") {
        return Token{Else, bp.pos-4, 4, nil}, nil
    } else if bp.match("while") {
        return Token{While, bp.pos-5, 5, nil}, nil
    } else if bp.match("try") {
        return Token{Try, bp.pos-3, 3, nil}, nil
    } else if bp.match("true") {
        return Token{True, bp.pos-4, 4, nil}, nil
    } else if bp.match("false") {
        return Token{False, bp.pos-5, 5, nil}, nil
    }

    char, eof := bp.readByte()
    if eof {
        return Token{}, EOF
    }

    switch(char) {
        // Symbols
        case '{':
            return Token{LBrace, bp.pos-1, 1, nil}, nil
        case '}':
            return Token{RBrace, bp.pos-1, 1, nil}, nil
        case '(':
            return Token{LParen, bp.pos-1, 1, nil}, nil
        case ')':
            return Token{RParen, bp.pos-1, 1, nil}, nil
        case ',':
            return Token{Comma, bp.pos-1, 1, nil}, nil

        // Values
        // TODO: add support for boolean and bitwise
        // TODO: add a operator type enum
        case '+', '-', '*', '/', '%':
            return Token{OperatorToken, bp.pos-1, 1, bp.buf[bp.pos-1]}, nil
        case '0','1','2','3','4','5','6','7','8','9':
            bp.unreadByte()
            return readNumber(bp)
        case '"':
            bp.unreadByte()
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
