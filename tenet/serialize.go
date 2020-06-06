package tenet

import (
	"bytes"
	"encoding/binary"
)

func (b TnByte) Serialize(w *bytes.Buffer) {
	w.WriteByte(byte(b))
}

func (s TnShort) Serialize(w *bytes.Buffer) {
	binary.Write(w, binary.LittleEndian, s)
}
func (i TnInt) Serialize(w *bytes.Buffer) {
	binary.Write(w, binary.LittleEndian, i)
}

func (f TnFloat) Serialize(w *bytes.Buffer) {
	binary.Write(w, binary.LittleEndian, f)
}

func (c TnColor) Serialize(w *bytes.Buffer) {
	w.WriteByte(c.R)
	w.WriteByte(c.G)
	w.WriteByte(c.B)
}
func (v TnVector) Serialize(w *bytes.Buffer) {
	binary.Write(w, binary.LittleEndian, v.X)
	binary.Write(w, binary.LittleEndian, v.Y)
}

func (a TnArray) Serialize(w *bytes.Buffer) {
	for i := 0; i < a.Count; i++ {
		a.Seed.Serialize(w)
	}
}

func (s TnString) Serialize(w *bytes.Buffer) {
	ln := len(s)
	for ln > 127 {
		w.WriteByte(byte(ln & 255))
		ln = ln >> 7
	}
	w.WriteByte(byte(ln))
	w.WriteString(string(s))
}

func (s TnNetString) Serialize(w *bytes.Buffer) {

}
