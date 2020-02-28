package lexer

import "bytes"

type bufpos struct {
    buf []byte
    pos int
}

// Reads a byte array until it doesnt match a char in the chars array
func readWhile(bp *bufpos, chars []byte) []byte {
    initialPos := bp.pos
    for {
        char, eof := bp.readByte()
        if !bytes.ContainsRune(chars, rune(char)) {
            bp.unreadByte()
            return bp.buf[initialPos:bp.pos]
        }
        if eof {
            return bp.buf[initialPos:bp.pos]
        }
    }
}

func (b *bufpos) copy() *bufpos {
    var b2 bufpos = *b
    return &b2
}

func (b *bufpos) copyAt(pos int) *bufpos {
    b2 := b.copy()
    b2.pos = pos
    return b2
}

// Returns the byte and a boolean that is true when it cannot read
func (bp *bufpos) readByte() (byte, bool) {
    if bp.len() > 0 {
        bp.pos++
        return bp.buf[bp.pos-1], false
    } else {
        bp.pos++
        return 0, true
    }
}

func (bp *bufpos) peek() (byte, bool) {
    if bp.len() > 0 {
        return bp.buf[bp.pos], false
    } else {
        return 0, true
    }

}

func (bp *bufpos) unreadByte() {
    bp.pos--
}

func (b *bufpos) len() int {
    return len(b.buf) - b.pos
}

// TODO: needs tests
// Reads if it matches the pattern, otherwise unreads
func (bp *bufpos) match(str string) bool {
    for i, v := range []byte(str) {
        char, eof := bp.readByte()
        if eof || char != v {
            bp.pos -= i + 1
            return false
        }
    }
    return true
}
