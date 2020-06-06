package temap

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"strings"
)

type Handle struct {
	reader *bufio.Reader
	handle *os.File
}

func (h *Handle) LoadFile(name string) bool {
	f, err := os.Open(name)
	if err == nil {
		h.handle = f
		h.reader = bufio.NewReader(f)
		return false
	} else {
		// log.Println("invalid file", err)
		return true
	}
}

func (h *Handle) LoadString(data string) {
	h.reader = bufio.NewReader(strings.NewReader(data))
}

func (h *Handle) LoadBytes(data []byte) {
	h.reader = bufio.NewReader(bytes.NewReader(data))
}

func (h *Handle) Close() {
	if h.handle != nil {
		h.handle.Close()
	}
	h.reader = nil
	h.handle = nil
}

func (h *Handle) Bool() (result bool) {
	var u8 uint8
	binary.Read(h.reader, binary.LittleEndian, &u8)
	result = u8 != 0
	return
}

func (h *Handle) U64() (result uint64) {
	binary.Read(h.reader, binary.LittleEndian, &result)
	return
}

func (h *Handle) U32() (result uint32) {
	binary.Read(h.reader, binary.LittleEndian, &result)
	return
}

func (h *Handle) U16() (result uint16) {
	binary.Read(h.reader, binary.LittleEndian, &result)
	return
}

func (h *Handle) U8() (result uint8) {
	binary.Read(h.reader, binary.LittleEndian, &result)
	return
}

func (h *Handle) F() (result float32) {
	binary.Read(h.reader, binary.LittleEndian, &result)
	return
}

func (h *Handle) D() (result float64) {
	binary.Read(h.reader, binary.LittleEndian, &result)
	return
}

func (h *Handle) S() (result string) {
	ln := 0
	s := 0
	for true {
		u7 := int(h.U8())
		ln |= (u7 & 0x7f) << s
		s += 7
		if (u7 & 0x80) == 0 {
			break
		}
	}
	buf := make([]byte, ln)
	h.reader.Read(buf)
	result = string(buf)
	return
}

func (h *Handle) Skip(cnt int) {
	h.reader.Discard(cnt)
}

func (h *Handle) Seek(offset int64, whence int) {
	h.handle.Seek(offset, whence)
	h.reader = bufio.NewReader(h.handle)
}

func (h *Handle) Read(cnt int) []byte {
	buf := make([]byte, cnt)
	io.ReadFull(h.reader, buf)
	return buf
}
