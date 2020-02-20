package lexer

import (
    "bytes"
    "errors"
    "strconv"
)

var EOF error
type InvalidTokenError bufpos
type NumberSyntaxError struct {
    bp *bufpos
    numError *strconv.NumError
}

// returns line and character at the buffers current position
func getLine(bp *bufpos) (int, int) {
    line := bytes.Count(bp.buf[:bp.pos], []byte("\n"))

    lastNewline := bytes.LastIndexByte(bp.buf[:bp.pos], '\n')
    if lastNewline == -1 {
        lastNewline = 0
    }

    character := bp.pos - lastNewline
    return line, character
}

// Returns a string of the buffers position represented as "line:character"
func lineStr(bp *bufpos) string {
    line, char := getLine(bp)
    return strconv.Itoa(line) + ":" + strconv.Itoa(char)
}

func (e *NumberSyntaxError) Error() string {
    return lineStr(e.bp) + " Cannot parse number '" + e.numError.Num + "'"
}

func (e *InvalidTokenError) Error() string {
    return lineStr((*bufpos)(e)) + " Does not match any token pattern"
}

func init() {
    EOF = errors.New("EOF")
}
