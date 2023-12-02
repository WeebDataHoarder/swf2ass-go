package ass

import (
	"fmt"
	"slices"
	"strconv"
	"sync/atomic"
)

type rendererStatsEntry struct {
	Frames *atomic.Uint64
	Lines  *atomic.Uint64
	Size   *atomic.Uint64
}

type RendererStatistics []RendererStatistic

func (s RendererStatistics) SortByFrames() {
	slices.SortStableFunc(s, func(a, b RendererStatistic) int {
		return int(b.Frames - a.Frames)
	})
}

func (s RendererStatistics) SortByLines() {
	slices.SortStableFunc(s, func(a, b RendererStatistic) int {
		return int(b.Lines - a.Lines)
	})
}

func (s RendererStatistics) SortBySize() {
	slices.SortStableFunc(s, func(a, b RendererStatistic) int {
		return int(b.Size - a.Size)
	})
}

func (s RendererStatistics) Reverse() {
	slices.Reverse(s)
}

func (s RendererStatistics) Strings() (r []string) {
	maxObjectId := len("Object Id")
	maxFrames := len("Frames")
	maxLines := len("Lines")
	maxSize := len("Lines")
	for _, e := range s {
		oId, frames, lines, size := e.Strings()

		maxObjectId = max(maxObjectId, len(oId))
		maxFrames = max(maxFrames, len(frames))
		maxLines = max(maxLines, len(lines))
		maxSize = max(maxSize, len(size))
	}

	maxIndex := len(strconv.FormatUint(uint64(len(s)), 10))

	for i, e := range s {
		oId, frames, lines, size := e.Strings()
		i := strconv.FormatUint(uint64(i), 10)

		r = append(r, fmt.Sprintf("| %*s | %*s | %*s | %*s | %*s |", maxIndex, i, maxObjectId, oId, maxFrames, frames, maxLines, lines, maxSize, size))
	}
	slices.Insert(r, 0, fmt.Sprintf("| %*s | %*s | %*s | %*s | %*s |", maxIndex, "#", maxObjectId, "Object Id", maxFrames, "Frames", maxLines, "Lines", maxSize, "Size"))
	return r
}

type RendererStatistic struct {
	ObjectId uint16
	Frames   uint64
	Lines    uint64
	Size     uint64
}

func (s RendererStatistic) Strings() (string, string, string, string) {
	oId := strconv.FormatUint(uint64(s.ObjectId), 10)
	frames := strconv.FormatUint(s.Frames, 10)
	lines := strconv.FormatUint(s.Lines, 10)
	size := strconv.FormatFloat(float64(s.Size)/(1<<20), 'f', 2, 64) + " MiB"
	return oId, frames, lines, size
}

func (r *Renderer) AggregateStatistics() (s RendererStatistics) {
	for objectId, stats := range r.Statistics {
		s = append(s, RendererStatistic{
			ObjectId: objectId,
			Frames:   stats.Frames.Load(),
			Lines:    stats.Lines.Load(),
			Size:     stats.Size.Load(),
		})
	}
	return s
}
