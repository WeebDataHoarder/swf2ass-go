package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"golang.org/x/exp/maps"
	"slices"
)

type ShapeConverter struct {
	Styles StyleList

	FillStyle0 *ActivePath
	FillStyle1 *ActivePath
	LineStyle  *ActivePath

	Position math.Vector2[types.Twip]

	Fills, Strokes PendingPathMap

	Commands DrawPathList

	Finished bool

	FirstElement, SecondElement subtypes.SHAPERECORDS
}

func NewShapeConverter(element subtypes.SHAPERECORDS, styles StyleList) *ShapeConverter {
	return &ShapeConverter{
		Styles:       styles,
		Position:     math.NewVector2[types.Twip](0, 0),
		Fills:        make(PendingPathMap),
		Strokes:      make(PendingPathMap),
		FirstElement: element,
	}
}

func NewMorphShapeConverter(firstElement, secondElement subtypes.SHAPERECORDS, styles StyleList) *ShapeConverter {
	return &ShapeConverter{
		Styles:        styles,
		Position:      math.NewVector2[types.Twip](0, 0),
		Fills:         make(PendingPathMap),
		Strokes:       make(PendingPathMap),
		FirstElement:  firstElement,
		SecondElement: secondElement,
	}
}

func (c *ShapeConverter) Convert(flipElements bool) {
	if c.Finished {
		return
	}

	firstElement := c.FirstElement
	secondElement := c.SecondElement
	for {
		var a, b subtypes.SHAPERECORD
		if len(secondElement) > 0 {
			b = secondElement[0]
		}

		if len(firstElement) > 0 {
			a = firstElement[0]
		} else {
			if b != nil {
				panic("a finished, b did not")
			}
			break
		}

		if b == nil {
			if flipElements {
				panic("b finished, a did not")
			}
			c.HandleNode(a)
			firstElement = firstElement[1:]
			continue
		}

		if c.Finished {
			panic("more paths after end")
		}

		//TODO: check!

		if a.RecordType() == b.RecordType() {
			switch a := a.(type) {
			case *subtypes.StyleChangeRecord:
				bCopy := *b.(*subtypes.StyleChangeRecord)
				aCopy := *a

				if aCopy.Flag.NewStyles {
					bCopy.Flag.NewStyles = aCopy.Flag.NewStyles
					bCopy.FillStyles = aCopy.FillStyles
					bCopy.LineStyles = aCopy.LineStyles
				}
				if aCopy.Flag.LineStyle {
					bCopy.Flag.LineStyle = aCopy.Flag.LineStyle
					bCopy.LineStyle = aCopy.LineStyle
				}
				if aCopy.Flag.FillStyle0 {
					bCopy.Flag.FillStyle0 = aCopy.Flag.FillStyle0
					bCopy.FillStyle0 = aCopy.FillStyle0
				}
				if aCopy.Flag.FillStyle1 {
					bCopy.Flag.FillStyle1 = aCopy.Flag.FillStyle1
					bCopy.FillStyle1 = aCopy.FillStyle1
				}

				if !flipElements && !aCopy.Flag.MoveTo && bCopy.Flag.MoveTo {
					aCopy.Flag.MoveTo = bCopy.Flag.MoveTo
					aCopy.MoveDeltaX = c.Position.X
					aCopy.MoveDeltaY = c.Position.Y
				}

				if flipElements && aCopy.Flag.MoveTo && !bCopy.Flag.MoveTo {
					bCopy.Flag.MoveTo = aCopy.Flag.MoveTo
					bCopy.MoveDeltaX = c.Position.X
					bCopy.MoveDeltaY = c.Position.Y
				}

				if flipElements {
					c.HandleNode(&bCopy)
				} else {
					c.HandleNode(&aCopy)
				}
			}

			firstElement = firstElement[1:]
			secondElement = secondElement[1:]
		} else if a2, ok := a.(*subtypes.StyleChangeRecord); ok {
			bCopy := *a2

			if bCopy.Flag.MoveTo {
				bCopy.MoveDeltaX = c.Position.X
				bCopy.MoveDeltaY = c.Position.Y
			}

			if flipElements {
				c.HandleNode(&bCopy)
			} else {
				c.HandleNode(a2)
			}
			firstElement = firstElement[1:]
		} else if b2, ok := b.(*subtypes.StyleChangeRecord); ok {
			aCopy := *b2

			if aCopy.Flag.MoveTo {
				aCopy.MoveDeltaX = c.Position.X
				aCopy.MoveDeltaY = c.Position.Y
			}

			if flipElements {
				c.HandleNode(b2)
			} else {
				c.HandleNode(&aCopy)
			}
			secondElement = secondElement[1:]
		} else {
			//Curve/line records can differ

			if flipElements {
				c.HandleNode(b)
			} else {
				c.HandleNode(a)
			}

			firstElement = firstElement[1:]
			secondElement = secondElement[1:]
		}
	}

	c.FlushLayer()
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
			c.Styles = StyleListFromSWFItems(node.FillStyles, node.LineStyles)
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
	point := VisitedPoint{
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
		var newSegments PendingPath
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
			c.Commands = append(c.Commands, DrawPathStroke(&fixedStyle, path.GetShape()))
		}
	}
	clear(c.Strokes)
}
