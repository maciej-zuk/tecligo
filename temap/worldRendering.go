package temap

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (w *World) Render(x0 uint32, y0 uint32, x1 uint32, y1 uint32) *sdl.Surface {
	img, _ := sdl.CreateRGBSurfaceWithFormat(0, int32(16*(x1-x0)), int32(16*(y1-y0)), 32, sdl.PIXELFORMAT_ABGR8888)
	ctx := RenderingContext{x0, y0, x1, y1, 0, 0, img}
	w.RenderOnto(&ctx)
	return img
}

func (w *World) RenderOnto(ctx *RenderingContext) {
	w.drawBackground(ctx)
	w.drawWalls(ctx)
	w.drawTiles(ctx)
	w.drawNPCS(ctx)
	w.drawLiquids(ctx)
}

func (w *World) drawRepeated(x0 int32, y0 int32, sw int32, sh int32, tex *sdl.Surface, ctx *RenderingContext) {
	px0 := x0
	px1 := px0 + sw
	py0 := y0
	py1 := py0 + sh
	sx := tex.ClipRect.W
	sy := tex.ClipRect.H
	dr := sdl.Rect{}
	rx0 := int32(ctx.X0*16 - ctx.Xoff)
	ry0 := int32(ctx.Y0*16 - ctx.Yoff)
	rx1 := int32(ctx.X1 * 16)
	ry1 := int32(ctx.Y1 * 16)
	if px1 > rx1 {
		px1 = rx1
	}
	if py1 > ry1 {
		py1 = ry1
	}
	for x := px0; x < px1; x += sx {
		for y := py0; y < py1; y += sy {
			dr.X = x - rx0
			dr.Y = y - ry0
			dr.W = sx
			dr.H = sy
			tex.Blit(&tex.ClipRect, ctx.Img, &dr)
		}
	}
}

var backStyles = []int{
	66, 67, 68, 69, 128, 125, 185,
	70, 71, 68, 72, 128, 125, 185,
	73, 74, 75, 76, 134, 125, 185,
	77, 78, 79, 82, 134, 125, 185,
	83, 84, 85, 86, 137, 125, 185,
	83, 87, 88, 89, 137, 125, 185,
	121, 122, 123, 124, 140, 125, 185,
	153, 147, 148, 149, 150, 125, 185,
	146, 154, 155, 156, 157, 125, 185,
}

var trackUVs = []int16{
	0, 0, 0, 1, 0, 0, 2, 1, 1, 3, 1, 1, 0, 2, 8, 1, 2, 4,
	0, 1, 0, 1, 1, 0, 0, 3, 4, 1, 3, 8, 4, 1, 9, 5, 1, 5,
	6, 1, 1, 7, 1, 1, 2, 0, 0, 3, 0, 0, 4, 0, 8, 5, 0, 4,
	6, 0, 0, 7, 0, 0, 0, 4, 0, 1, 4, 0, 0, 5, 0, 1, 5, 0,
	2, 2, 2, 3, 2, 2, 4, 2, 10, 5, 2, 6, 6, 2, 2, 7, 2, 2,
	2, 3, 0, 3, 3, 0, 4, 3, 4, 5, 3, 8, 6, 3, 4, 7, 3, 8,
}

func (w *World) drawBackground(ctx *RenderingContext) {
	var tex *sdl.Surface
	groundLevel := int(w.header.WorldSurface)
	rockLevel := int(w.header.RockLayer)
	hellLevel := (int(w.tilesHigh-330) - groundLevel) / 6
	hellLevel = hellLevel*6 + groundLevel - 5
	hellBottom := (int(w.tilesHigh-200) - hellLevel) / 6
	hellBottom = hellBottom*6 + hellLevel - 5

	hellStyle := int(w.header.HellBackStyle)

	if ctx.Y0 < uint32(groundLevel) {
		ctx.Img.FillRect(&ctx.Img.ClipRect, 0xfff8aa83)
	}
	if ctx.Y1 > uint32(hellBottom) {
		ctx.Img.FillRect(&ctx.Img.ClipRect, 0xff19111e)
	}
	lastX := 0
	nextX := 0

	var style int
	for i := 0; i <= 3; i++ {
		style = int(w.header.CaveBackStyle[i]) * 7
		if style > len(backStyles)-1 {
			style = 0
		}
		if i == 3 {
			nextX = int(w.tilesWide)
		} else {
			nextX = int(w.header.CaveBackX[i])
		}
		tex = GetTextureCrop(Tex_Background|backStyles[style], 128, 16)
		w.drawRepeated(int32(lastX)*16, int32(groundLevel-1)*16, int32(nextX-lastX)*16, int32(1)*16, tex, ctx)
		tex = GetTextureCrop(Tex_Background|backStyles[style+1], 128, 96)
		w.drawRepeated(int32(lastX)*16, int32(groundLevel)*16, int32(nextX-lastX)*16, int32(rockLevel-groundLevel)*16, tex, ctx)
		tex = GetTextureCrop(Tex_Background|backStyles[style+2], 128, 16)
		w.drawRepeated(int32(lastX)*16, int32(rockLevel)*16, int32(nextX-lastX)*16, int32(1)*16, tex, ctx)
		tex = GetTextureCrop(Tex_Background|backStyles[style+3], 128, 96)
		w.drawRepeated(int32(lastX)*16, int32(rockLevel+1)*16, int32(nextX-lastX)*16, int32(hellLevel-(rockLevel+1))*16, tex, ctx)
		tex = GetTextureCrop(Tex_Background|backStyles[style+4]+hellStyle, 128, 16)
		w.drawRepeated(int32(lastX)*16, int32(hellLevel)*16, int32(nextX-lastX)*16, int32(1)*16, tex, ctx)
		tex = GetTextureCrop(Tex_Background|backStyles[style+5]+hellStyle, 128, 96)
		w.drawRepeated(int32(lastX)*16, int32(hellLevel+1)*16, int32(nextX-lastX)*16, int32(hellBottom-(hellLevel+1))*16, tex, ctx)
		tex = GetTextureCrop(Tex_Background|backStyles[style+6]+hellStyle, 128, 16)
		w.drawRepeated(int32(lastX)*16, int32(hellBottom)*16, int32(nextX-lastX)*16, int32(1)*16, tex, ctx)
		lastX = nextX
	}
}

func (w *World) drawTiles(ctx *RenderingContext) {
	for x := ctx.X0; x < ctx.X1; x++ {
		for y := ctx.Y0; y < ctx.Y1; y++ {
			tileData := &w.tiles[w.pt(x, y)]
			if !tileData.Active {
				continue
			}
			if tileData.U < 0 {
				fixTile(w, x, y)
			}
			w.putTile(ctx, x, y, tileData)
		}
	}
}

func (w *World) drawWalls(ctx *RenderingContext) {
	for x := ctx.X0; x < ctx.X1; x++ {
		for y := ctx.Y0; y < ctx.Y1; y++ {
			tileData := &w.tiles[w.pt(x, y)]
			if tileData.Wall > 0 {
				_x := (x - ctx.X0) * 16
				_y := (y - ctx.Y0) * 16
				if tileData.Wallu < 0 {
					fixWall(w, x, y)
				}
				tex := GetTexture(Tex_Wall | int(tileData.Wall))
				w.draw(ctx, tex, tileData.Wallu, tileData.Wallv, 32, 32, _x-8, _y-8)
			}
		}
	}
}

func (w *World) drawLiquids(ctx *RenderingContext) {
	for x := ctx.X0; x < ctx.X1; x++ {
		for y := ctx.Y0; y < ctx.Y1; y++ {
			tileData := &w.tiles[w.pt(x, y)]
			if tileData.Liquid < 1 {
				continue
			}
			_x := (x - ctx.X0) * 16
			_y := (y - ctx.Y0) * 16
			waterLevel := uint32((255 - tileData.Liquid) / 16)
			w.draw(ctx, helperTexture, (int16(tileData.LiquidType))*16, 0, 16, 16-waterLevel, _x, _y+waterLevel)
		}
	}
}

func (w *World) draw(ctx *RenderingContext, tex *sdl.Surface, sx int16, sy int16, dw uint32, dh uint32, x uint32, y uint32) {
	sr := sdl.Rect{int32(sx), int32(sy), int32(dw), int32(dh)}
	dr := sdl.Rect{int32(x - ctx.Yoff), int32(y - ctx.Yoff), int32(dw), int32(dh)}
	tex.Blit(&sr, ctx.Img, &dr)
}

func (w *World) putTile(ctx *RenderingContext, x uint32, y uint32, tile *Tile) {
	info := &w.info.Tiles[tile.Type]
	var texw, texh uint32
	texw = uint32(info.Width - 2)
	texh = uint32(info.Height - 2)
	if tile.Type == 5 && tile.U >= 22 && tile.V >= 198 {
		var variant uint32
		var style uint32
		switch tile.V {
		case 220:
			variant = 1
		case 242:
			variant = 2
		default:
			variant = 0
		}

		if tile.U == 22 {
			texw = uint32(80)
			texh = uint32(80)
			padx := uint32(32)
			style = uint32(w.findTreeStyle(x, y))
			if style == 2 || style == 11 || style == 13 {
				texw = 114
				texh = 96
				padx = 48
			}
			if style == 3 {
				texh = 140
				if x%3 == 1 {
					variant += 3
				} else if x%3 == 2 {
					variant += 6
				}
			}
			tex := GetTexture(Tex_TreeTops | int(style))
			w.draw(
				ctx,
				tex,
				int16(variant*(texw+2)),
				0,
				texw,
				texh,
				(x-ctx.X0)*16-padx,
				(y-ctx.Y0+1)*16-texh,
			)
		} else if tile.U == 44 {
			style = uint32(w.findBranchStyle(x+1, y))
			if style == 3 {
				if x%3 == 1 {
					variant += 3
				} else if x%3 == 2 {
					variant += 6
				}
			}
			tex := GetTexture(Tex_TreeBranches | int(style))
			w.draw(
				ctx,
				tex,
				0,
				int16(variant*42),
				40,
				40,
				(x-ctx.X0)*16-24,
				(y-ctx.Y0)*16-12,
			)
		} else if tile.U == 66 {
			style = uint32(w.findBranchStyle(x-1, y))
			if style == 3 {
				if x%3 == 1 {
					variant += 3
				} else if x%3 == 2 {
					variant += 6
				}
			}
			tex := GetTexture(Tex_TreeBranches | int(style))
			w.draw(
				ctx,
				tex,
				42,
				int16(variant*42),
				40,
				40,
				(x-ctx.X0)*16,
				(y-ctx.Y0)*16-12,
			)
		}
	} else if tile.Type == 323 && tile.U >= 88 && tile.U <= 132 {
		var variant uint32
		switch tile.U {
		case 110:
			variant = 1
			break
		case 132:
			variant = 2
			break
		default:
			variant = 0
			break
		}
		style := w.findPalmTreeStyle(x, y)
		tex := GetTexture(Tex_TreeTops | 15)
		w.draw(
			ctx,
			tex,
			int16(variant*82),
			int16(style*82),
			80,
			80,
			uint32(int(x-ctx.X0)*16-32+int(tile.V)),
			(y-ctx.Y0+1)*16-80,
		)
	} else if tile.Type == 72 && tile.U >= 36 {
		variant := 0
		switch tile.V {
		case 18:
			variant = 1
			break
		case 36:
			variant = 2
			break
		default:
			variant = 0
			break
		}
		tex := GetTexture(Tex_Shroom)
		w.draw(
			ctx,
			tex,
			int16(variant*62),
			0,
			60,
			42,
			(x-ctx.X0)*16-22,
			(y-ctx.Y0)*16-26,
		)
	} else if tile.Type == 314 {
		u := tile.U
		v := tile.V
		tex := GetTexture(Tex_Tile | 314)
		w.draw(
			ctx,
			tex,
			trackUVs[u*3]*18,
			trackUVs[u*3+1]*18,
			texw,
			texh,
			(x-ctx.X0)*16,
			(y-ctx.Y0)*16+uint32(info.Toppad),
		)
		if tile.V >= 0 {
			w.draw(
				ctx,
				tex,
				trackUVs[v*3]*18,
				trackUVs[v*3+1]*18,
				texw,
				texh,
				(x-ctx.X0)*16,
				(y-ctx.Y0)*16+uint32(info.Toppad),
			)
		}
		if (u >= 0 && u < 36) || (v >= 0 && v < 36) {
			mask := trackUVs[u*3+2]
			var mask2 int16 = 0
			if v >= 0 {
				mask2 = trackUVs[v*3+2]
			}
			if mask&8 != 0 || mask2&8 != 0 {
				w.draw(
					ctx,
					tex,
					0,
					108,
					texw,
					texh,
					(x-ctx.X0)*16,
					(y-ctx.Y0+1)*16+uint32(info.Toppad),
				)
			}
			if mask&4 != 0 || mask2&4 != 0 {
				w.draw(
					ctx,
					tex,
					18,
					108,
					texw,
					texh,
					(x-ctx.X0)*16,
					(y-ctx.Y0+1)*16+uint32(info.Toppad),
				)
			}
			if mask&2 != 0 || mask2&2 != 0 {
				w.draw(
					ctx,
					tex,
					18,
					126,
					texw,
					texh,
					(x-ctx.X0)*16,
					(y-ctx.Y0-1)*16+uint32(info.Toppad),
				)
			}
			if mask&1 != 0 || mask2&1 != 0 {
				w.draw(
					ctx,
					tex,
					0,
					126,
					texw,
					texh,
					(x-ctx.X0)*16,
					(y-ctx.Y0-1)*16+uint32(info.Toppad),
				)
			}
		}
	} else if tile.Type == 237 && tile.U == 0 && tile.V == 0 {
		tex := GetTexture(Tex_Tile | int(tile.Type))
		w.draw(
			ctx,
			tex,
			tile.U,
			tile.V,
			texw,
			texh,
			(x-ctx.X0)*16,
			(y-ctx.Y0)*16+uint32(info.Toppad),
		)
	} else if tile.Type == 323 {
		ry := y
		for w.tiles[w.pt(x, ry)].Active && w.tiles[w.pt(x, ry)].Type == 323 {
			ry++
		}
		variant := w.findPalmTreeStyle(x, ry)
		tex := GetTexture(Tex_Tile | int(tile.Type))
		w.draw(
			ctx,
			tex,
			tile.U,
			int16(22*variant),
			texw,
			texh,
			uint32(int(x-ctx.X0)*16+int(tile.V)),
			(y-ctx.Y0)*16+uint32(info.Toppad),
		)
	} else if tile.Type == 5 {
		toff := y*w.tilesWide + x
		if tile.U == 66 && tile.V <= 45 {
			toff++
		}
		if tile.U == 88 && tile.V >= 66 && tile.V <= 110 {
			toff--
		}
		if tile.U == 22 && tile.V >= 132 {
			toff--
		}
		if tile.U == 44 && tile.V >= 132 {
			toff++
		}
		for w.tiles[toff].Active && w.tiles[toff].Type == 5 {
			toff += w.tilesWide
		}
		variant := w.getTreeVariant(toff)
		var tex *sdl.Surface
		if variant == -1 {
			tex = GetTexture(Tex_Tile | int(tile.Type))
		} else {
			tex = GetTexture(Tex_Wood | int(variant))
		}
		w.draw(
			ctx,
			tex,
			tile.U,
			tile.V,
			16,
			16,
			(x-ctx.X0)*16,
			(y-ctx.Y0)*16,
		)
	} else {
		ry := (y-ctx.Y0)*16 + uint32(info.Toppad)
		if tile.Half {
			texh -= 8
			ry += 8
		}
		tex := GetTexture(Tex_Tile | int(tile.Type))
		w.draw(
			ctx,
			tex,
			tile.U,
			tile.V,
			texw,
			texh,
			(x-ctx.X0)*16-(texw-16)/2,
			ry,
		)
	}
}

func (w *World) drawNPCS(ctx *RenderingContext) {
	stride := w.tilesWide
	for _, npc := range w.npcs {
		if npc.Head != 0 {
			hx := npc.HomeX
			hy := npc.HomeY - 1
			offset := hy*w.tilesWide + hx
			for !w.tiles[offset].Active || !w.info.Tiles[w.tiles[offset].Type].Solid {
				hy--
				offset -= stride
				if hy < 10 {
					break
				}
			}
			hy++
			offset += stride
			if hx >= ctx.X0 && hx < ctx.X1 && hy >= ctx.Y0 && hy < ctx.Y1 {
				dy := uint32(18)
				if w.tiles[offset-stride].Type == 19 {
					dy -= 8
				}
				tex := GetTexture(Tex_Banner)
				w.draw(
					ctx, tex, 0, 0, 32, 40,
					(hx-ctx.X0)*16-uint32(tex.ClipRect.W/2),
					(hy-ctx.Y0)*16+dy-uint32(tex.ClipRect.H/2),
				)
				tex = GetTexture(Tex_NPCHead | int(npc.Head))
				w.draw(
					ctx, tex, 0, 0,
					uint32(tex.ClipRect.W), uint32(tex.ClipRect.H),
					(hx-ctx.X0)*16-uint32(tex.ClipRect.W/2),
					(hy-ctx.Y0)*16+dy-uint32(tex.ClipRect.H/2),
				)
			}
		}
	}
}
