package tag

import (
	"bytes"
	"compress/zlib"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"image"
	"image/color"
	"io"
)

type DefineBitsLossless struct {
	_              struct{} `swfFlags:"root"`
	CharacterId    uint16
	Format         uint8
	Width, Height  uint16
	ColorTableSize uint8 `swfCondition:"HasColorTableSize()"`
	ZlibData       types.Bytes
}

func (t *DefineBitsLossless) HasColorTableSize(ctx types.ReaderContext) bool {
	return t.Format == 3
}

func (t *DefineBitsLossless) GetImage() image.Image {
	r, err := zlib.NewReader(bytes.NewReader(t.ZlibData))
	if err != nil {
		return nil
	}
	defer r.Close()

	var buf [3]byte

	switch t.Format {
	case 3: // 8-bit colormapped image

		var palette color.Palette
		for i := 0; i < (int(t.ColorTableSize) + 1); i++ {
			_, err = io.ReadFull(r, buf[:])
			if err != nil {
				return nil
			}
			palette = append(palette, color.RGBA{R: buf[0], G: buf[1], B: buf[2], A: 255})
		}

		im := image.NewPaletted(image.Rectangle{
			Min: image.Point{},
			Max: image.Point{X: int(t.Width), Y: int(t.Height)},
		}, palette)
		for y := 0; y < int(t.Height); y++ {
			for x := 0; x < int(t.Width); x++ {
				_, err = io.ReadFull(r, buf[:1])
				if err != nil {
					return nil
				}
				im.SetColorIndex(x, y, buf[0])
			}
		}
		return im
	case 4: // 15-bit RGB image
		panic("15-bit RGB not supported")
	case 5: // 24-bit RGB image
		im := image.NewRGBA(image.Rectangle{
			Min: image.Point{},
			Max: image.Point{X: int(t.Width), Y: int(t.Height)},
		})

		for y := 0; y < int(t.Height); y++ {
			for x := 0; x < int(t.Width); x++ {
				_, err = io.ReadFull(r, buf[:])
				if err != nil {
					return nil
				}
				im.SetRGBA(x, y, color.RGBA{R: buf[0], G: buf[1], B: buf[2], A: 255})
			}
		}
		return im
	default:
		panic("not supported")
	}
}

func (t *DefineBitsLossless) Code() Code {
	return RecordDefineBitsLossless
}
