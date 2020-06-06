package tenet

import (
	"bytes"
	"encoding/binary"
	"strings"
)

func DeSerializeByte(r *bytes.Reader) (d TnByte) {
	binary.Read(r, binary.LittleEndian, &d)
	return
}

func DeSerializeShort(r *bytes.Reader) (d TnShort) {
	binary.Read(r, binary.LittleEndian, &d)
	return
}

func DeSerializeInt(r *bytes.Reader) (d TnInt) {
	binary.Read(r, binary.LittleEndian, &d)
	return
}

func DeSerializeFloat(r *bytes.Reader) (d TnFloat) {
	binary.Read(r, binary.LittleEndian, &d)
	return
}

func DeSerializeColor(r *bytes.Reader) (c TnColor) {
	binary.Read(r, binary.LittleEndian, &c.R)
	binary.Read(r, binary.LittleEndian, &c.G)
	binary.Read(r, binary.LittleEndian, &c.B)
	return
}

func DeSerializeVector(r *bytes.Reader) (v TnVector) {
	binary.Read(r, binary.LittleEndian, &v.X)
	binary.Read(r, binary.LittleEndian, &v.Y)
	return
}

func DeSerializeString(r *bytes.Reader) (d TnString) {
	tmp, _ := r.ReadByte()
	ln := int(tmp)
	start := 1
	var fln int
	if ln <= 127 {
		fln = ln
	} else {
		fln = ln - 128
		for {
			start++
			tmp, _ := r.ReadByte()
			ln := int(tmp)
			if ln <= 127 {
				fln += ln * (1 << ((start-1)*8 - 1))
				break
			} else {
				ln -= 128
				fln += ln * (1 << ((start-1)*8 - 1))
			}

		}
	}
	buf := make([]byte, fln)
	r.Read(buf)
	d = TnString(buf)
	return
}

func DeSerializeNetString(r *bytes.Reader) (d TnNetString) {
	tmp := DeSerializeByte(r)
	if tmp == 0 {
		d = TnNetString(DeSerializeString(r))
	} else {
		text := DeSerializeString(r)
		subsCount := int(DeSerializeByte(r))
		subs := make([]string, subsCount)
		for i := 0; i < subsCount; i++ {
			subs[i] = string(DeSerializeNetString(r))
		}
		if tmp == 2 {
			split := strings.Split(string(text), ".")
			if len(split) == 2 {
				cat, found := lang[split[0]]
				if found {
					val, found := cat[split[1]]
					if found {
						text = TnString(val)
					}
				}

			}
		}
		d = TnNetString(formatLang(string(text), subs))
	}
	return
}
