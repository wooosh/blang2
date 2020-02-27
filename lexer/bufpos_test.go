package lexer

import (
    "testing"
)

func TestBufPos(t *testing.T) {
    t.Run("Read", func(t *testing.T) {
        bp := bufpos{[]byte("12"), 0}
        char1, eof1 := bp.readByte()
        char2, eof2 := bp.readByte()
        _, eof3 := bp.readByte()

        if char1 != '1' || eof1 || char2 != '2' || eof2 || !eof3 {
            t.Fatal("Recieved unexpected eof or char")
        }
    })

    t.Run("Peek", func(t *testing.T) {
        bp := bufpos{[]byte("1"), 0}
        char1, eof1 := bp.peek()
        char2, eof2 := bp.readByte()
        _, eof3 := bp.peek()

        if char1 != '1' || char2 != '1' || eof1 || eof2 || !eof3 {
            t.Fatal("Recieved unexpected eof or char")
        }
    })

    t.Run("Unread", func(t *testing.T) {
        bp := bufpos{[]byte("1"), 0}
        char1, eof1 := bp.readByte()
        bp.unreadByte()
        char2, eof2 := bp.readByte()

        if char1 != '1' || char2 != '1' || eof1 || eof2 {
            t.Fatal("Recieved unexpected eof or char")
        }
    })
}
