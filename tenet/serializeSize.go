package tenet

func (b TnByte) Size() int {
	return 1
}

func (s TnShort) Size() int {
	return 2
}

func (i TnInt) Size() int {
	return 4
}

func (f TnFloat) Size() int {
	return 4
}

func (c TnColor) Size() int {
	return 3
}

func (v TnVector) Size() int {
	return 8
}

func (a TnArray) Size() int {
	return a.Seed.Size() * a.Count
}

func (s TnString) Size() int {
	ln := len(s)
	size := 1
	for ln > 127 {
		ln = ln >> 7
		size++
	}
	return size + len(s)
}

func (s TnNetString) Size() int {
	return 0
}
