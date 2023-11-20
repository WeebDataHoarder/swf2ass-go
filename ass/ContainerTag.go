package ass

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
)

type ContainerTag struct {
	Tags        []Tag
	Transitions map[int64][]Tag

	BakeTransforms *types.MatrixTransform
}

func (t *ContainerTag) TransitionColor(line *Line, transform types.ColorTransform) ColorTag {
	container := t.Clone()

	index := line.End - line.Start

	for _, tag := range container.Tags {
		if colorTag, ok := tag.(ColorTag); ok {
			newTag := colorTag.TransitionColor(line, transform)
			if newTag == nil {
				return nil
			}
			if !newTag.Equals(tag) {
				container.Transitions[index] = append(container.Transitions[index], newTag)
			}
		}
	}
	return container
}

func (t *ContainerTag) TransitionMatrixTransform(line *Line, transform types.MatrixTransform) ColorTag {
	if t.BakeTransforms != nil {
		//Do not allow matrix changes, except moves
		if !transform.EqualsWithoutTranslation(*t.BakeTransforms, types.TransformCompareEpsilon) {
			return nil
		}
	}

	container := t.Clone()

	index := line.End - line.Start

	for _, tag := range container.Tags {
		if colorTag, ok := tag.(PositioningTag); ok {
			newTag := colorTag.TransitionMatrixTransform(line, transform)
			if newTag == nil {
				return nil
			}
			if !newTag.Equals(tag) {
				//Special case
				if _, ok := newTag.(*PositionTag); ok {
					container.Tags = append(container.Tags, newTag)
				} else {
					container.Transitions[index] = append(container.Transitions[index], newTag)
				}
			}
		}
	}
	return container
}

func (t *ContainerTag) TransitionStyleRecord(line *Line, record types.StyleRecord) StyleTag {
	container := t.Clone()

	index := line.End - line.Start

	for _, tag := range container.Tags {
		if colorTag, ok := tag.(StyleTag); ok {
			newTag := colorTag.TransitionStyleRecord(line, record)
			if newTag == nil {
				return nil
			}
			if !newTag.Equals(tag) {
				container.Transitions[index] = append(container.Transitions[index], newTag)
			}
		}
	}
	return container
}

func (t *ContainerTag) TransitionShape(line *Line, shape *types.Shape) PathTag {
	container := t.Clone()

	index := line.End - line.Start

	for _, tag := range container.Tags {
		if colorTag, ok := tag.(PathTag); ok {
			newTag := colorTag.TransitionShape(line, shape)
			if newTag == nil {
				return nil
			}
			if !newTag.Equals(tag) {
				container.Transitions[index] = append(container.Transitions[index], newTag)
			}
		}
	}
	return container
}

func (t *ContainerTag) TransitionClipPath(line *Line, clip *types.ClipPath) ClipPathTag {
	container := t.Clone()

	index := line.End - line.Start

	for _, tag := range container.Tags {
		if colorTag, ok := tag.(ClipPathTag); ok {
			newTag := colorTag.TransitionClipPath(line, clip)
			if newTag == nil {
				return nil
			}
			if !newTag.Equals(tag) {
				container.Transitions[index] = append(container.Transitions[index], newTag)
			}
		}
	}
	return container
}

func (t *ContainerTag) ApplyColorTransform(transform types.ColorTransform) ColorTag {
	panic("not supported")
}

func (t *ContainerTag) FromMatrixTransform(transform types.MatrixTransform) PositioningTag {
	panic("not supported")
}

func (t *ContainerTag) FromStyleRecord(record types.StyleRecord) StyleTag {
	panic("not supported")
}

func (t *ContainerTag) Clone() *ContainerTag {
	var transform *types.MatrixTransform
	if t.BakeTransforms != nil {
		t2 := *t.BakeTransforms
		transform = &t2
	}
	transitions := make(map[int64][]Tag, len(t.Transitions))
	for k := range t.Transitions {
		transitions[k] = slices.Clone(t.Transitions[k])
	}
	return &ContainerTag{
		Tags:           slices.Clone(t.Tags),
		Transitions:    transitions,
		BakeTransforms: transform,
	}
}

func (t *ContainerTag) Equals(tag Tag) bool {
	if o, ok := tag.(*ContainerTag); ok && len(t.Tags) == len(o.Tags) {
		//TODO: optimize this?
		tags := slices.Clone(t.Tags)
		otherTags := slices.Clone(o.Tags)

		for len(tags) > 0 {
			if !func() bool {
				i := len(tags) - 1
				t1 := tags[len(tags)-1]
				for j, t2 := range otherTags {
					if t1.Equals(t2) {
						tags = slices.Delete(tags, i, i+1)
						otherTags = slices.Delete(otherTags, j, j+1)
						return true
					}
				}
				return false
			}() {
				return false
			}

		}

		return len(tags) == 0 && len(otherTags) == 0
	}

	return false
}

func (t *ContainerTag) Encode(event EventTime) string {
	text := make([]string, len(t.Tags)*2)
	for _, tag := range t.Tags {
		if _, ok := tag.(DrawingTag); !ok {
			text = append(text, tag.Encode(event))
		}
	}
	keys := maps.Keys(t.Transitions)
	slices.Sort(keys)
	for _, index := range keys {
		if len(t.Transitions[index]) == 0 {
			continue
		}

		var startTime, endTime int64
		//TODO: clone line?
		//TODO: animations with smoothing really don't play well. maybe allow them when only one animation "direction" exists, or smooth them manually?
		//Or just don't animate MatrixTransform / do it in a single tick
		if GlobalSettings.SmoothTransitions {
			startTime = event.GetDurationFromStartOffset(index-1).Milliseconds() + 1
			endTime = event.GetDurationFromStartOffset(index+1).Milliseconds() - 1
		} else {
			startTime = event.GetDurationFromStartOffset(index).Milliseconds() - 1
			endTime = event.GetDurationFromStartOffset(index).Milliseconds()
		}
		transitionText := make([]string, 0, len(t.Transitions[index]))
		for _, ttag := range t.Transitions[index] {
			transitionText = append(transitionText, ttag.Encode(event))
		}
		text = append(text, fmt.Sprintf("\\t(%d,%d,%s)", startTime, endTime, strings.Join(transitionText, "")))
	}

	for _, tag := range t.Tags {
		if _, ok := tag.(DrawingTag); ok {
			text = append(text, tag.Encode(event))
		}
	}
	return strings.Join(text, "")
}

func (t *ContainerTag) TryAppend(tag Tag) {
	if tag != nil {
		t.Tags = append(t.Tags, tag)
	}
	panic("tag is nil")
}

var identityMatrixTransform = types.IdentityTransform()

func ContainerTagFromPathEntry(path types.DrawPath, clip *types.ClipPath, colorTransform types.ColorTransform, matrixTransform types.MatrixTransform, bakeTransforms bool) *ContainerTag {
	container := &ContainerTag{
		Transitions: make(map[int64][]Tag),
	}

	container.TryAppend(NewClipTag(clip, GlobalSettings.DrawingScale))

	/*
		//TODO Convert to fill????
		if($path->style instanceof LineStyleRecord){ //Convert to fill

		}
	*/

	container.TryAppend((&BorderTag{}).FromStyleRecord(path.Style))

	{
		lineColorTag := &LineColorTag{}
		lineColorTag.FromStyleRecord(path.Style)
		container.TryAppend(lineColorTag.ApplyColorTransform(colorTransform))
	}

	{
		fillColorTag := &FillColorTag{}
		fillColorTag.FromStyleRecord(path.Style)
		container.TryAppend(fillColorTag.ApplyColorTransform(colorTransform))
	}

	if bakeTransforms {
		container.BakeTransforms = &matrixTransform

		container.TryAppend((&PositionTag{}).FromMatrixTransform(matrixTransform))

		drawTag := DrawingTag(NewDrawTag(path.Commands, GlobalSettings.DrawingScale))
		if !matrixTransform.EqualsExact(identityMatrixTransform) {
			drawTag = drawTag.ApplyMatrixTransform(matrixTransform, false)
		}

		container.TryAppend(drawTag)
	} else {
		container.TryAppend((&PositionTag{}).FromMatrixTransform(matrixTransform))
		container.TryAppend((&MatrixTransformTag{}).FromMatrixTransform(matrixTransform))

		container.TryAppend(NewDrawTag(path.Commands, GlobalSettings.DrawingScale))
	}

	return container
}
