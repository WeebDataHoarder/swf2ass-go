package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"golang.org/x/exp/maps"
	"slices"
)

type ShapeConverter struct {
	Collection ObjectCollection
	Styles     StyleList

	FillStyle0 *ActivePath
	FillStyle1 *ActivePath
	LineStyle  *ActivePath

	Position math.Vector2[types.Twip]

	Fills, Strokes PendingPathMap

	Commands DrawPathList

	Finished bool
}

func NewShapeConverter(collection ObjectCollection, styles StyleList) *ShapeConverter {
	return &ShapeConverter{
		Collection: collection,
		Styles:     styles,
		Position:   math.NewVector2[types.Twip](0, 0),
		Fills:      make(PendingPathMap),
		Strokes:    make(PendingPathMap),
	}
}

func NewMorphShapeConverter(collection ObjectCollection, styles StyleList) *ShapeConverter {
	return &ShapeConverter{
		Collection: collection,
		Styles:     styles,
		Position:   math.NewVector2[types.Twip](0, 0),
		Fills:      make(PendingPathMap),
		Strokes:    make(PendingPathMap),
	}
}

func (c *ShapeConverter) Convert(elements subtypes.SHAPERECORDS) DrawPathList {
	if c.Finished {
		return nil
	}

	for _, e := range elements {
		c.HandleNode(e)
	}

	c.FlushLayer()

	return c.Commands
}

// ConvertMorph
// We step through both the start records and end records, interpolating edges pairwise.
// Fill style/line style changes should only appear in the start records.
// However, StyleChangeRecord move_to can appear it both start and end records,
// and not necessarily in matching pairs; therefore, we have to keep track of the pen position
// in case one side is missing a move_to; it will implicitly use the last pen position.
func (c *ShapeConverter) ConvertMorph(start, end subtypes.SHAPERECORDS) (startList, endList subtypes.SHAPERECORDS) {
	if c.Finished {
		return nil, nil
	}

	var startPos, endPos math.Vector2[types.Twip]

	updatePos := func(v math.Vector2[types.Twip], s subtypes.SHAPERECORD) math.Vector2[types.Twip] {
		switch s := s.(type) {
		case *subtypes.StraightEdgeRecord:
			v = v.AddVector(math.NewVector2(s.DeltaX, s.DeltaY))
		case *subtypes.CurvedEdgeRecord:
			v = v.AddVector(math.NewVector2(s.ControlDeltaX+s.AnchorDeltaX, s.ControlDeltaY+s.AnchorDeltaY))
		case *subtypes.StyleChangeRecord:
			if s.Flag.MoveTo {
				v = math.NewVector2(s.MoveDeltaX, s.MoveDeltaY)
			}
		}
		return v
	}

	for len(start) > 0 {
		startPtr := start[0]
		endPtr := end[0]

		if startPtr.RecordType() == endPtr.RecordType() {
			switch s := startPtr.(type) {
			case *subtypes.StyleChangeRecord:
				if s.Flag.MoveTo {
					startPos = math.NewVector2(s.MoveDeltaX, s.MoveDeltaY)
				}
				e := endPtr.(*subtypes.StyleChangeRecord)
				endRecord := *s
				endRecord.Flag.MoveTo = e.Flag.MoveTo
				endRecord.MoveDeltaX = e.MoveDeltaX
				endRecord.MoveDeltaY = e.MoveDeltaY
				if e.Flag.MoveTo {
					endPos = math.NewVector2(e.MoveDeltaX, e.MoveDeltaY)
				}
				startList = append(startList, s)
				endList = append(endList, &endRecord)

				start = start[1:]
				end = end[1:]
			default:
				startList = append(startList, startPtr)
				endList = append(endList, endPtr)

				start = start[1:]
				end = end[1:]
			}
		} else {
			if s, ok := startPtr.(*subtypes.StyleChangeRecord); ok {
				endRecord := *s
				if s.Flag.MoveTo {
					startPos = math.NewVector2(s.MoveDeltaX, s.MoveDeltaY)
					endRecord.MoveDeltaX = endPos.X
					endRecord.MoveDeltaY = endPos.Y
				}
				startList = append(startList, startPtr)
				endList = append(endList, &endRecord)
				startPos = updatePos(startPos, startPtr)
				start = start[1:]
			} else if e, ok := endPtr.(*subtypes.StyleChangeRecord); ok {
				startRecord := *e
				if e.Flag.MoveTo {
					endPos = math.NewVector2(e.MoveDeltaX, e.MoveDeltaY)
					startRecord.MoveDeltaX = startPos.X
					startRecord.MoveDeltaY = startPos.Y
				}
				startList = append(startList, &startRecord)
				endList = append(endList, endPtr)
				endPos = updatePos(endPos, startPtr)
				end = end[1:]
			} else {
				startList = append(startList, startPtr)
				endList = append(endList, endPtr)
				startPos = updatePos(startPos, startPtr)
				endPos = updatePos(endPos, endPtr)

				start = start[1:]
				end = end[1:]
			}
		}
	}

	if len(end) > 0 {
		panic("did not complete")
	}

	return
}

func (c *ShapeConverter) HandleNode(node subtypes.SHAPERECORD) {
	switch node := node.(type) {
	case *subtypes.StyleChangeRecord:
		if node.Flag.MoveTo {
			moveTo := math.NewVector2[types.Twip](node.MoveDeltaX, node.MoveDeltaY)
			c.Position = moveTo
			c.FlushPaths()
		}

		if node.Flag.NewStyles {
			c.FlushLayer()
			c.Styles = StyleListFromSWFItems(c.Collection, node.FillStyles, node.LineStyles)
		}

		if node.Flag.FillStyle1 {
			if c.FillStyle1 != nil {
				c.Fills.MergePath(c.FillStyle1, true)
			}

			if node.FillStyle1 > 0 {
				c.FillStyle1 = NewActivePath(int(node.FillStyle1), c.Position)
			} else {
				c.FillStyle1 = nil
			}
		}

		if node.Flag.FillStyle0 {
			if c.FillStyle0 != nil {
				if !c.FillStyle0.Segment.IsEmpty() {
					c.FillStyle0.Flip()
					c.Fills.MergePath(c.FillStyle0, true)
				}
			}

			if node.FillStyle0 > 0 {
				c.FillStyle0 = NewActivePath(int(node.FillStyle0), c.Position)
			} else {
				c.FillStyle0 = nil
			}
		}

		if node.Flag.LineStyle {
			if c.LineStyle != nil {
				c.Strokes.MergePath(c.LineStyle, false)
			}

			if node.LineStyle > 0 {
				c.LineStyle = NewActivePath(int(node.LineStyle), c.Position)
			} else {
				c.LineStyle = nil
			}
		}
	case *subtypes.StraightEdgeRecord:
		to := c.Position.AddVector(math.NewVector2[types.Twip](node.DeltaX, node.DeltaY))
		c.VisitPoint(to, false)
		c.Position = to
	case *subtypes.CurvedEdgeRecord:
		control := c.Position.AddVector(math.NewVector2[types.Twip](node.ControlDeltaX, node.ControlDeltaY))
		anchor := control.AddVector(math.NewVector2[types.Twip](node.AnchorDeltaX, node.AnchorDeltaY))
		c.VisitPoint(control, true)
		c.VisitPoint(anchor, false)
		c.Position = anchor
	case *subtypes.EndShapeRecord:
		c.Finished = true
	}
}

func (c *ShapeConverter) VisitPoint(pos math.Vector2[types.Twip], isBezierControlPoint bool) {
	point := VisitedPoint[types.Twip]{
		Pos:             pos,
		IsBezierControl: isBezierControlPoint,
	}

	if c.FillStyle0 != nil {
		c.FillStyle0.AddPoint(point)
	}

	if c.FillStyle1 != nil {
		c.FillStyle1.AddPoint(point)
	}

	if c.LineStyle != nil {
		c.LineStyle.AddPoint(point)
	}
}

func (c *ShapeConverter) FlushPaths() {
	if c.FillStyle1 != nil {
		c.Fills.MergePath(c.FillStyle1, true)
		c.FillStyle1 = NewActivePath(c.FillStyle1.StyleId, c.Position)
	}

	if c.FillStyle0 != nil {
		if !c.FillStyle0.Segment.IsEmpty() {
			c.FillStyle0.Flip()
			c.Fills.MergePath(c.FillStyle0, true)
		}
		c.FillStyle0 = NewActivePath(c.FillStyle0.StyleId, c.Position)
	}

	if c.LineStyle != nil {
		c.Strokes.MergePath(c.LineStyle, false)
		c.LineStyle = NewActivePath(c.LineStyle.StyleId, c.Position)
	}
}

func (c *ShapeConverter) FlushLayer() {
	c.FlushPaths()

	c.FillStyle0 = nil
	c.FillStyle1 = nil
	c.LineStyle = nil

	fillsKeys := maps.Keys(c.Fills)
	slices.Sort(fillsKeys)
	for _, styleId := range fillsKeys {
		path := c.Fills[styleId]
		if styleId <= 0 || styleId > len(c.Styles.FillStyles) {
			panic("should not happen")
		}

		style := c.Styles.GetFillStyle(styleId - 1)
		if style == nil {
			panic("should not happen")
		}
		c.Commands = append(c.Commands, DrawPathFill(style, path.GetShape()))
	}
	clear(c.Fills)

	strokesKeys := maps.Keys(c.Strokes)
	slices.Sort(strokesKeys)
	for _, styleId := range strokesKeys {
		path := c.Strokes[styleId]
		if styleId <= 0 || styleId > len(c.Styles.LineStyles) {
			panic("should not happen")
		}

		style := c.Styles.GetLineStyle(styleId - 1)
		if style == nil {
			panic("should not happen")
		}

		//wrap around all segments, even if closed. ASS does NOT like them otherwise. so we draw everything backwards to have border around the line, not just on one side
		//TODO: using custom line borders later using fills this can be removed
		var newSegments PendingPath[types.Twip]
		for _, segment := range *path {
			other := slices.Clone(*segment)
			other.Flip()
			segment.Merge(other)

			newSegments.MergePath(segment, false)
		}

		if len(newSegments) > 0 {
			//Reduce width of line style to account for double border
			//TODO: using custom line borders later using fills this can be removed
			fixedStyle := *style
			fixedStyle.Width /= 2
			c.Commands = append(c.Commands, DrawPathStroke(&fixedStyle, newSegments.GetShape()))
		}
		//TODO: leave this as-is and create a fill in renderer
	}
	clear(c.Strokes)
}
