package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf-go/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"github.com/ctessum/polyclip-go"
	"github.com/nfnt/resize"
	"golang.org/x/exp/maps"
	"image"
	"image/color"
	"image/png"
	math2 "math"
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
				//reduce alpha resolution a bit
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

	var hasAlpha bool
	colors := make(map[math.PackedColor]polyclip.Polygon)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			r, g, b, a := i.At(x, y).RGBA()
			p := math.NewPackedColor(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
			contour := polyclip.Contour{
				{float64(x), float64(y)},
				{float64(x), float64(y + 1)},
				{float64(x + 1), float64(y + 1)},
				{float64(x + 1), float64(y)},
			}
			if p.Alpha() < math2.MaxUint8 {
				hasAlpha = true
			}
			colors[p] = append(colors[p], contour)
		}
	}

	keys := maps.Keys(colors)

	var n atomic.Uint64

	type result struct {
		Color   math.PackedColor
		Polygon polyclip.Polygon
	}

	results := make([][]result, runtime.NumCPU())

	for cpuN := 0; cpuN < runtime.NumCPU(); cpuN++ {
		wg.Add(1)
		go func(cpuN int) {
			defer wg.Done()
			for {
				i := n.Add(1) - 1
				if i >= uint64(len(keys)) {
					break
				}

				results[cpuN] = append(results[cpuN], result{
					Color:   keys[i],
					Polygon: colors[keys[i]].Simplify(),
				})
			}
		}(cpuN)
	}
	wg.Wait()

	for _, cs := range results {
		for _, v := range cs {
			colors[v.Color] = v.Polygon
		}
	}

	//Sort from the highest size to lowest
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
		//make a rectangle covering the whole first area to optimize this case
		r = append(r, DrawPathFill(&FillStyleRecord{
			Fill: keys[0].Color(),
		}, Rectangle[float64]{
			TopLeft:     math.NewVector2[float64](0, 0),
			BottomRight: math.NewVector2(float64(size.X+1), float64(size.Y+1)),
		}.Draw()))

		keys = keys[1:]
	}

	for _, k := range keys {
		if k.Alpha() == 0 {
			//Skip fully transparent pixels
			continue
		}
		pol := colors[k]
		//Draw resulting shape
		r = append(r, DrawPathFill(&FillStyleRecord{
			Fill: k.Color(),
		}, ComplexPolygon{
			Pol: pol,
		}.GetShape()))
	}

	scale := math.ScaleTransform(math.NewVector2(ratioX, ratioY))
	r2 := r.ApplyMatrixTransform(scale, true)
	return r2.(DrawPathList)
}
