package shapes

import (
	"bytes"
	"encoding/binary"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"github.com/ctessum/polyclip-go"
	"github.com/nfnt/resize"
	"golang.org/x/exp/maps"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
)

func ConvertBitmapBytesToDrawPathList(imageData []byte, alphaData []byte) (DrawPathList, error) {
	var im image.Image
	var err error

	for i, s := range bitmapHeaderFormats {
		if bytes.Compare(s, imageData[:len(s)]) == 0 {
			if i == 0 || i == 1 {
				//jpeg
				//remove invalid data
				jpegData := removeInvalidJPEGData(imageData)
				im, _, err = image.Decode(bytes.NewReader(jpegData))
				if im != nil {
					size := im.Bounds().Size()
					if len(alphaData) == size.X*size.Y {

						newIm := image.NewRGBA(im.Bounds())
						for x := 0; x < size.X; x++ {
							for y := 0; y < size.Y; y++ {
								rI, gI, bI, _ := im.At(x, y).RGBA()

								// The JPEG data should be premultiplied alpha, but it isn't in some incorrect SWFs.
								// This means 0% alpha pixels may have color and incorrectly show as visible.
								// Flash Player clamps color to the alpha value to fix this case.
								// Only applies to DefineBitsJPEG3; DefineBitsLossless does not seem to clamp.
								a := alphaData[y*size.X+x]
								if a != 0 {
									runtime.KeepAlive(a)
								}
								r := min(uint8(rI>>8), a)
								g := min(uint8(gI>>8), a)
								b := min(uint8(bI>>8), a)
								newIm.SetRGBA(x, y, color.RGBA{
									R: r,
									G: g,
									B: b,
									A: a,
								})
							}
						}
						im = newIm
					}
				}
			} else if i == 2 {
				//png
				im, _, err = image.Decode(bytes.NewReader(imageData))
			} else if i == 3 {
				//gif
				im, _, err = image.Decode(bytes.NewReader(imageData))
			}
			break
		}
	}
	if err != nil {
		return nil, err
	}

	drawPathList := ConvertBitmapToDrawPathList(im)
	return drawPathList, nil
}

func QuantizeBitmap(i image.Image) image.Image {
	size := i.Bounds().Size()

	palettedImage := image.NewPaletted(i.Bounds(), nil)
	quantizer := MedianCutQuantizer{
		NumColor: settings.GlobalSettings.BitmapPaletteSize,
	}
	quantizer.Quantize(palettedImage, i.Bounds(), i, image.Point{})

	// Restore alpha
	newIm := image.NewRGBA(i.Bounds())
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			r, g, b, _ := palettedImage.At(x, y).RGBA()
			_, _, _, a := i.At(x, y).RGBA()
			if a == 0 {
				newIm.SetRGBA(x, y, color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				})
			} else if (a >> 8) == 255 {
				newIm.SetRGBA(x, y, color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: 255,
				})
			} else {

				newIm.SetRGBA(x, y, color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8((a >> 9) << 1),
				})
			}
		}
	}
	return newIm
}

func ConvertBitmapToDrawPathList(i image.Image) (r DrawPathList) {

	size := i.Bounds().Size()

	ratioX := 1.0
	ratioY := 1.0
	maxDimension := max(size.X, size.Y)
	if maxDimension > settings.GlobalSettings.BitmapMaxDimension && size.X == maxDimension {

		ratio := float64(maxDimension) / float64(settings.GlobalSettings.BitmapMaxDimension)
		w, h := uint(settings.GlobalSettings.BitmapMaxDimension), uint(float64(size.Y)/ratio)

		i = resize.Resize(w, h, i, resize.Bicubic)
		ratioX = float64(size.X+1) / float64(w+1)
		ratioY = float64(size.Y+1) / float64(h+1)

	} else if maxDimension > settings.GlobalSettings.BitmapMaxDimension && size.Y == maxDimension {

		ratio := float64(maxDimension) / float64(settings.GlobalSettings.BitmapMaxDimension)
		w, h := uint(float64(size.X)/ratio), uint(settings.GlobalSettings.BitmapMaxDimension)

		i = resize.Resize(w, h, i, resize.Bicubic)
		ratioX = float64(size.X+1) / float64(w+1)
		ratioY = float64(size.Y+1) / float64(h+1)
	}

	i = QuantizeBitmap(i)

	size = i.Bounds().Size()

	var wg sync.WaitGroup
	var x atomic.Uint64

	results := make([]map[math.PackedColor]polyclip.Polygon, runtime.NumCPU())
	for n := 0; n < runtime.NumCPU(); n++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			myResults := make(map[math.PackedColor]polyclip.Polygon)
			for {
				iX := x.Add(1) - 1
				if iX >= uint64(size.X) {
					break
				}

				for y := 0; y < size.Y; y++ {
					r, g, b, a := i.At(int(iX), y).RGBA()

					p := math.NewPackedColor(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
					poly := polyclip.Polygon{{
						{float64(iX), float64(y)},
						{float64(iX), float64(y + 1)},
						{float64(iX + 1), float64(y + 1)},
						{float64(iX + 1), float64(y)},
					}}
					if existingColor, ok := myResults[p]; ok {
						u := existingColor.Construct(polyclip.UNION, poly).Simplify()
						myResults[p] = u
					} else {
						myResults[p] = poly
					}
				}
			}
			results[n] = myResults
		}(n)
	}
	wg.Wait()

	var hasAlpha bool

	colors := make(map[math.PackedColor]polyclip.Polygon)

	for _, r := range results {
		for k, c := range r {
			if k.Alpha() < 255 {
				hasAlpha = true
			}
			if k.Alpha() == 0 {
				//Skip fully transparent pixels
				continue
			}
			if existingColor, ok := colors[k]; ok {
				u := existingColor.Construct(polyclip.UNION, c).Simplify()
				colors[k] = u
			} else {
				colors[k] = c.Simplify()
			}
		}
	}

	//Sort from the highest size to lowest
	keys := maps.Keys(colors)
	getSize := func(p polyclip.Polygon) (r int) {
		for _, c := range p {
			r += c.Len()
		}
		return r
	}
	slices.SortFunc(keys, func(a, b math.PackedColor) int {
		sizeA := getSize(colors[a])
		sizeB := getSize(colors[b])
		if sizeA > sizeB {
			return -1
		} else if sizeB > sizeA {
			return 1
		} else {
			return 0
		}
	})

	// Full shape optimizations when alpha is not in use
	if !hasAlpha {

		/*
			for i, k := range keys {
				pol := colors[k]
				pol1 := pol
				//Iterate through all previous layers and merge
				for _, k2 := range keys[:i] {
					//Check each sub-polygon of the shape to see if it is within previous indicative of a good merge, merge only those
					for _, pol4 := range colors[k2].Union(pol).Polygons() {
						if pol4.Bounds().Within(pol) == geom.Inside {
							pol = pol.Union(pol4)
						}
					}
				}
				//Draw resulting shape
				r = append(r, DrawPathFill(&FillStyleRecord{
					Fill: k.Color(),
				}, ComplexPolygon{
					Pol: pol.Simplify(PolygonSimplifyTolerance).(geom.Polygonal),
				}.GetShape(), nil))
			}*/

		//make a rectangle covering the whole first area to optimize this case
		r = append(r, DrawPathFill(&FillStyleRecord{
			Fill: keys[0].Color(),
		}, Rectangle[float64]{
			TopLeft:     math.NewVector2[float64](0, 0),
			BottomRight: math.NewVector2(float64(size.X+1), float64(size.Y+1)),
		}.Draw()))

		for _, k := range keys[1:] {
			pol := colors[k]
			//Draw resulting shape
			r = append(r, DrawPathFill(&FillStyleRecord{
				Fill: k.Color(),
			}, ComplexPolygon{
				Pol: pol,
			}.GetShape()))
		}

	} else {
		for _, k := range keys {
			pol := colors[k]
			//Draw resulting shape
			r = append(r, DrawPathFill(&FillStyleRecord{
				Fill: k.Color(),
			}, ComplexPolygon{
				Pol: pol,
			}.GetShape()))
		}
	}

	scale := math.ScaleTransform(math.NewVector2(ratioX, ratioY))
	r2 := r.ApplyMatrixTransform(scale, true)
	return r2.(DrawPathList)
}

var bitmapHeaderJPEG = []byte{0xff, 0xd8}
var bitmapHeaderJPEGInvalid = []byte{0xff, 0xd9, 0xff, 0xd8}
var bitmapHeaderPNG = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
var bitmapHeaderGIF = []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}

var bitmapHeaderFormats = [][]byte{
	bitmapHeaderJPEG,
	bitmapHeaderJPEGInvalid,
	bitmapHeaderPNG,
	bitmapHeaderGIF,
}

// removeInvalidJPEGData
// SWF19 errata p.138:
// "Before version 8 of the SWF file format, SWF files could contain an erroneous header of 0xFF, 0xD9, 0xFF, 0xD8
// before the JPEG SOI marker."
// 0xFFD9FFD8 is a JPEG EOI+SOI marker pair. Contrary to the spec, this invalid marker sequence can actually appear
// at any time before the 0xFFC0 SOF marker, not only at the beginning of the data. I believe this is a relic from
// the SWF JPEGTables tag, which stores encoding tables separately from the DefineBits image data, encased in its
// own SOI+EOI pair. When these data are glued together, an interior EOI+SOI sequence is produced. The Flash JPEG
// decoder expects this pair and ignores it, despite standard JPEG decoders stopping at the EOI.
// When DefineBitsJPEG2 etc. were introduced, the Flash encoders/decoders weren't properly adjusted, resulting in
// this sequence persisting. Also, despite what the spec says, this doesn't appear to be version checked (e.g., a
// v9 SWF can contain one of these malformed JPEGs and display correctly).
// See https://github.com/ruffle-rs/ruffle/issues/8775 for various examples.
func removeInvalidJPEGData(data []byte) (buf []byte) {
	const SOF0 uint8 = 0xC0 // Start of frame
	const RST0 uint8 = 0xD0 // Restart (we shouldn't see this before SOS, but just in case)
	const RST1 uint8 = 0xD0
	const RST2 uint8 = 0xD0
	const RST3 uint8 = 0xD0
	const RST4 uint8 = 0xD0
	const RST5 uint8 = 0xD0
	const RST6 uint8 = 0xD0
	const RST7 uint8 = 0xD7
	const SOI uint8 = 0xD8 // Start of image
	const EOI uint8 = 0xD9 // End of image

	if bytes.HasPrefix(data, bitmapHeaderJPEGInvalid) {
		data = bytes.TrimPrefix(data, bitmapHeaderJPEGInvalid)
	} else {
		// Parse the JPEG markers searching for the 0xFFD9FFD8 marker sequence to splice out.
		// We only have to search up to the SOF0 marker.
		// This might be another case where eventually we want to write our own full JPEG decoder to match Flash's decoder.
		jpegData := data
		var pos int
		for {
			if len(jpegData) < 4 {
				break
			}

			var payloadLength int

			if bytes.Compare([]byte{0xFF, EOI, 0xFF, SOI}, jpegData[:4]) == 0 {
				// Invalid EOI+SOI sequence found, splice it out.
				data = slices.Delete(slices.Clone(data), pos, pos+4)
				break
			} else if bytes.Compare([]byte{0xFF, EOI}, jpegData[:2]) == 0 { // EOI, SOI, RST markers do not include a size.

			} else if bytes.Compare([]byte{0xFF, SOI}, jpegData[:2]) == 0 {

			} else if bytes.Compare([]byte{0xFF, RST0}, jpegData[:2]) == 0 {

			} else if bytes.Compare([]byte{0xFF, RST1}, jpegData[:2]) == 0 {

			} else if bytes.Compare([]byte{0xFF, RST2}, jpegData[:2]) == 0 {

			} else if bytes.Compare([]byte{0xFF, RST3}, jpegData[:2]) == 0 {

			} else if bytes.Compare([]byte{0xFF, RST4}, jpegData[:2]) == 0 {

			} else if bytes.Compare([]byte{0xFF, RST5}, jpegData[:2]) == 0 {

			} else if bytes.Compare([]byte{0xFF, RST6}, jpegData[:2]) == 0 {

			} else if bytes.Compare([]byte{0xFF, RST7}, jpegData[:2]) == 0 {

			} else if bytes.Compare([]byte{0xFF, SOF0}, jpegData[:2]) == 0 {
				// No invalid sequence found before SOF marker, return data as-is.
				break
			} else if jpegData[0] == 0xFF {
				// Other tags include a length.
				payloadLength = int(binary.BigEndian.Uint16(jpegData[2:]))
			} else {
				// All JPEG markers should start with 0xFF.
				// So this is either not a JPEG, or we screwed up parsing the markers. Bail out.
				break
			}

			if len(jpegData) < payloadLength+2 {
				break
			}

			jpegData = jpegData[payloadLength+2:]
			pos += payloadLength + 2
		}
	}

	// Some JPEGs are missing the final EOI marker (JPEG optimizers truncate it?)
	// Flash and most image decoders will still display these images, but jpeg-decoder errors.
	// Glue on an EOI marker if its not already there and hope for the best.
	if bytes.HasSuffix(data, []byte{0xff, EOI}) {
		return data
	} else {
		//JPEG is missing EOI marker and may not decode properly
		return append(slices.Clone(data), []byte{0xff, EOI}...)
	}
}
