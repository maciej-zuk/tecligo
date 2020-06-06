package temap

func (w *World) pt(x, y uint32) uint32 {
	return y*w.tilesWide + x
}

func (w *World) getTreeVariant(offset uint32) int {
	t := w.tiles[offset].Type
	if t == 23 {
		return 0
	} else if t == 60 {
		if offset <= uint32(w.header.WorldSurface)*w.tilesWide {
			return 1
		} else {
			return 5
		}
	} else if t == 70 {
		return 6
	} else if t == 109 {
		return 2
	} else if t == 147 {
		return 3
	} else if t == 199 {
		return 4
	}
	return -1
}

func (w *World) treeStyle(x int) int {
	xs := w.header.TreeX
	i := 0
	for _, ix := range xs {
		if x <= int(ix) {
			break
		}
		i += 1
	}
	styles := w.header.TreeStyle
	style := styles[i]
	switch style {
	case 0:
		return 0
	case 5:
		return 10
	default:
		return int(style) + 5
	}
}

func (w *World) findTreeStyle(x, y uint32) int {
	var style int
	var snow int
	stride := w.tilesWide
	offset := y*stride + x

	for i := 0; i < 100; i++ {
		switch w.tiles[offset].Type {
		case 2: // grass
			return w.treeStyle(int(x))
		case 23: // corrupt grass
			return 1
		case 70: // mushroom grass
			return 14
		case 60: // jungle grass
			style = 2
			if int(w.header.Styles[2]) == 1 {
				style = 11
			}
			if offset > uint32(w.header.WorldSurface)*stride {
				style = 13
			}
			return style
		case 109: // hallowed grass
			return 3
		case 147: // snow
			style = 4
			snow = int(w.header.Styles[3])
			if snow == 0 {
				style = 12
				if x%10 == 0 {
					style = 18
				}
			}
			if snow == 2 || snow == 3 || snow == 4 || snow == 32 || snow == 42 {
				if snow&1 != 0 {
					if x <= w.tilesWide/2 {
						style = 17
					} else {
						style = 16
					}
				} else {
					if x >= w.tilesWide/2 {
						style = 17
					} else {
						style = 16
					}
				}
			}
			return style
		case 199: // flesh grass
			return 5
		}
		offset += stride
	}
	return 0
}

func (w *World) findBranchStyle(x, y uint32) int {
	var style int
	stride := w.tilesWide
	offset := y*stride + x
	for i := 0; i < 100; i++ {
		switch w.tiles[offset].Type {
		case 2: // grass
			return w.treeStyle(int(x))
		case 23: // corrupt grass
			return 1
		case 70: // mushroom grass
			return 14
		case 60: // jungle grass
			style = 2
			if offset > uint32(w.header.WorldSurface)*stride {
				style = 13
			}
			return style
		case 109: // hallowed grass
			return 3
		case 147: // snow
			style = 4
			if int(w.header.Styles[3]) == 0 {
				style = 12
			}
			return style
		case 199: // flesh grass
			return 5
		}
		offset += stride

	}
	return 0
}

func (w *World) findPalmTreeStyle(x, y uint32) int {
	for i := 0; i < 100; i++ {
		switch w.tiles[w.pt(x, y+uint32(i))].Type {
		case 53: // sand
			return 0
		case 112: // ebonsand
			return 3
		case 116: // pearlsand
			return 2
		case 234: // crimsand
			return 1
		}
	}
	return 0

}

// Get world size
func (w *World) GetSize() (uint32, uint32) {
	return w.tilesWide, w.tilesHigh
}
