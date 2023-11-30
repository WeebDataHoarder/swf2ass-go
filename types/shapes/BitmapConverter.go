package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"github.com/ctessum/polyclip-go"
	"github.com/nfnt/resize"
	"golang.org/x/exp/maps"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
)

func ImageToPNG(im image.Image, fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, im)
	if err != nil {
		return err
	}
	return nil
}

func ConvertBitmapBytesToDrawPathList(imageData []byte, alphaData []byte) (DrawPathList, error) {
	im, err := subtypes.DecodeImageBitsJPEG(imageData, alphaData)
	if err != nil {
		return nil, err
	}

	drawPathList := ConvertBitmapToDrawPathList(im)
	return drawPathList, nil
}

func QuantizeBitmap(i image.Image) image.Image {
	if settings.GlobalSettings.BitmapPaletteSize == 0 {
		return i
	}
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
	if settings.GlobalSettings.BitmapMaxDimension > 0 {
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
					contour := polyclip.Contour{
						{float64(iX), float64(y)},
						{float64(iX), float64(y + 1)},
						{float64(iX + 1), float64(y + 1)},
						{float64(iX + 1), float64(y)},
					}

					myResults[p] = append(myResults[p], contour)
					/*
						if _, ok := myResults[p]; ok {
							//u := existingColor.Construct(polyclip.UNION, poly)
						} else {
							myResults[p] = poly
						}*/
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
