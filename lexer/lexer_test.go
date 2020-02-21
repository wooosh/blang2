package lexer

import (
    "reflect"
    "testing"
)

// TODO: integration testing
func testToken(in string, tokType TokenType, expectedValue interface{}) func (*testing.T) {
    return func(t *testing.T) {
        tokens, err := Lex([]byte(in))
        if err != nil {
            t.Fatal(err)
        }
        if len(tokens) != 1 {
            t.Fatal("Expected 1 token, recieved", len(tokens))
        }
        tok := tokens[0]

        if tokType != tok.TypeOf {
            t.Fatal("Expected token type", tokType, "recieved", tok.TypeOf)
        }

        if !reflect.DeepEqual(expectedValue, tok.Value) {
            t.Fatal("Expected value and recived value do not match.")
        }

    }
}

func TestNumbers(t *testing.T) {
    t.Run("Decimal Integers", testToken("1234567890", NumberToken, 1234567890))
    t.Run("Hexadecimal Integers", testToken("0x1234567890abcdef", NumberToken, 0x1234567890abcdef))
    t.Run("Binary Integers", testToken("0b1100110101", NumberToken, 0b1100110101))
}

func TestStrings(t *testing.T) {
    // TODO: literal strings (ignore escape sequences)
    t.Run("Regular String", testToken(`"abcdefg1234"`, StringToken, []byte("abcdefg1234")))
    escapeSequences := map[string][]byte{
        `\n`: []byte{'\n'},
        `\t`: []byte{'\t'},
        `\r`: []byte{'\r'},
        `\e`: []byte{'\x1b'},
        `\\`: []byte{'\\'},
        `\"`: []byte{'"'},
    }
    for k,v := range escapeSequences {
        t.Run("Escape Sequence '" + k + "'", testToken(`"` + k + `"`, StringToken, v))
    }
}
