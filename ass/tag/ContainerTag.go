package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
)

type ContainerTag struct {
	Tags        []Tag
	Transitions map[int64][]Tag

	BakedTransform types.Option[math.MatrixTransform]
}

func (t *ContainerTag) HasColor() bool {
	for _, tag := range t.Tags {
		if colorTag, ok := tag.(ColorTag); ok {
			if colorTag.HasColor() {
				return true
			}
		}
	}
	for _, tags := range t.Transitions {
		for _, tag := range tags {
			if colorTag, ok := tag.(ColorTag); ok {
				if colorTag.HasColor() {
					return true
				}
			}
		}
	}
	return false
}

func (t *ContainerTag) TransitionColor(event Event, transform math.ColorTransform) ColorTag {
	container := t.Clone(false)

	index := event.GetEnd() - event.GetStart()

	//TODO: same color is added
	for _, tag := range container.Tags {
		if colorTag, ok := tag.(ColorTag); ok {
			newTag := colorTag.TransitionColor(event, transform)
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

func (t *ContainerTag) TransitionMatrixTransform(event Event, transform math.MatrixTransform) PositioningTag {
	if bakedTransform, ok := t.BakedTransform.Some(); ok {
		//Do not allow matrix changes, except moves
		if !transform.EqualsWithoutTranslation(bakedTransform, math.TransformCompareEpsilon) {
			return nil
		}
	}

	container := t.Clone(true)

	index := event.GetEnd() - event.GetStart()

	for i, tag := range container.Tags {
		if colorTag, ok := tag.(PositioningTag); ok {
			newTag := colorTag.TransitionMatrixTransform(event, transform)
			if newTag == nil {
				return nil
			}
			if !newTag.Equals(tag) {
				//Special case
				if _, ok := newTag.(*PositionTag); ok {
					container.Tags[i] = newTag
				} else {
					container.Transitions[index] = append(container.Transitions[index], newTag)
				}
			}
		}
	}
	return container
}

func (t *ContainerTag) TransitionStyleRecord(event Event, record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	container := t.Clone(false)

	index := event.GetEnd() - event.GetStart()

	for _, tag := range container.Tags {
		if colorTag, ok := tag.(StyleTag); ok {
			newTag := colorTag.TransitionStyleRecord(event, record, transform)
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

func (t *ContainerTag) TransitionShape(event Event, shape shapes.Shape) PathTag {
	container := t.Clone(false)

	index := event.GetEnd() - event.GetStart()

	for _, tag := range container.Tags {
		if colorTag, ok := tag.(PathTag); ok {
			newTag := colorTag.TransitionShape(event, shape)
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

func (t *ContainerTag) TransitionClipPath(event Event, clip *shapes.ClipPath) ClipPathTag {
	container := t.Clone(false)

	index := event.GetEnd() - event.GetStart()

	for _, tag := range container.Tags {
		if colorTag, ok := tag.(ClipPathTag); ok {
			newTag := colorTag.TransitionClipPath(event, clip)
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

func (t *ContainerTag) ApplyColorTransform(transform math.ColorTransform) ColorTag {
	panic("not supported")
}

func (t *ContainerTag) FromMatrixTransform(transform math.MatrixTransform) PositioningTag {
	panic("not supported")
}

func (t *ContainerTag) FromStyleRecord(record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	panic("not supported")
}

func (t *ContainerTag) Clone(cloneTags bool) *ContainerTag {
	bakedTransform := types.SomeWith(t.BakedTransform.Some())

	transitions := make(map[int64][]Tag, len(t.Transitions))
	for k := range t.Transitions {
		transitions[k] = slices.Clone(t.Transitions[k])
	}
	if cloneTags {
		return &ContainerTag{
			Tags:           slices.Clone(t.Tags),
			Transitions:    transitions,
			BakedTransform: bakedTransform,
		}
	} else {
		return &ContainerTag{
			Tags:           t.Tags,
			Transitions:    transitions,
			BakedTransform: bakedTransform,
		}
	}
}

func (t *ContainerTag) Equals(tag Tag) bool {
	if o, ok := tag.(*ContainerTag); ok && len(t.Tags) == len(o.Tags) {
		if o == t {
			return true
		}
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

func (t *ContainerTag) Encode(event time.EventTime) string {
	text := make([]string, 0, len(t.Tags)*2)
	var lastTransitionText []string
	for _, tag := range t.Tags {
		if _, ok := tag.(DrawingTag); !ok {
			thisText := tag.Encode(event)
			text = append(text, thisText)
			lastTransitionText = append(lastTransitionText, thisText)
		}
	}
	keys := maps.Keys(t.Transitions)
	slices.Sort(keys)
	for _, index := range keys {
		if len(t.Transitions[index]) == 0 {
			continue
		}
		if index > event.Duration {
			continue
		}

		var startTime, endTime int64
		//TODO: clone line?
		//TODO: animations with smoothing really don't play well. maybe allow them when only one animation "direction" exists, or smooth them manually?
		//Or just don't animate MatrixTransform / do it in a single tick
		if settings.GlobalSettings.ASSSmoothTransitions {
			startTime = event.GetDurationFromStartOffset(index-1).Milliseconds() + 1
			endTime = event.GetDurationFromStartOffset(index+1).Milliseconds() - 1
		} else {
			startTime = event.GetDurationFromStartOffset(index).Milliseconds() - 1
			endTime = event.GetDurationFromStartOffset(index).Milliseconds()
		}
		transitionText := make([]string, 0, len(t.Transitions[index]))

		for _, ttag := range t.Transitions[index] {
			thisTransitionText := ttag.Encode(event)
			//TODO: make this better with per-tag lookup for previous transition
			if !slices.Contains(lastTransitionText, thisTransitionText) {
				transitionText = append(transitionText, thisTransitionText)
			}
		}
		if len(transitionText) > 0 {
			text = append(text, fmt.Sprintf("\\t(%d,%d,%s)", startTime, endTime, strings.Join(transitionText, "")))
			lastTransitionText = transitionText
		}
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
		return
	}
	panic("tag is nil")
}

func ContainerTagFromPathEntry(path shapes.DrawPath, clip types.Option[shapes.ClipPath], colorTransform math.ColorTransform, matrixTransform math.MatrixTransform, bakeMatrixTransforms bool) *ContainerTag {
	container := &ContainerTag{
		Transitions: make(map[int64][]Tag),
	}

	if !matrixTransform.IsIdentity() {
		if bakeMatrixTransforms {
			path = path.ApplyMatrixTransform(matrixTransform, false)
		}
	}

	if settings.GlobalSettings.ASSBakeClips {
		clip.With(func(clip shapes.ClipPath) {
			//Clip is given in absolute coordinates. path is relative to translation
			//TODO: is this true for ClipPath???
			//TODO: this is broken
			translationTransform := math.TranslateTransform(matrixTransform.GetTranslation().Multiply(-1))
			path = shapes.DrawPath{
				Style: path.Style, //TODO: apply transform to Style?
				Shape: clip.ApplyMatrixTransform(translationTransform, true).ClipShape(path.Shape, true),
			}
		})
	} else {
		if clip, ok := clip.Some(); ok {
			container.TryAppend(NewClipTag(types.Some(clip.GetShape()), settings.GlobalSettings.ASSDrawingScale))
		} else {
			container.TryAppend(NewClipTag(types.None[shapes.Shape](), settings.GlobalSettings.ASSDrawingScale))
		}
	}

	/*
		//TODO Convert to fill????
		if($path->style instanceof LineStyleRecord){ //Convert to fill

		}
	*/

	if len(path.Shape) == 0 {
		return nil
	}

	if bakeMatrixTransforms {
		container.BakedTransform = types.Some(matrixTransform)
		container.TryAppend((&PositionTag{}).FromMatrixTransform(matrixTransform))
	} else {
		container.TryAppend((&PositionTag{}).FromMatrixTransform(matrixTransform))
		container.TryAppend((&MatrixTransformTag{}).FromMatrixTransform(matrixTransform))
	}

	container.TryAppend((&BorderTag{}).FromStyleRecord(path.Style, matrixTransform))

	container.TryAppend((&BlurGaussianTag{}).FromStyleRecord(path.Style, matrixTransform))

	{
		lineColorTag := &LineColorTag{}
		lineColorTag.FromStyleRecord(path.Style, matrixTransform)
		container.TryAppend(lineColorTag.ApplyColorTransform(colorTransform))
	}

	{
		fillColorTag := &FillColorTag{}
		fillColorTag.FromStyleRecord(path.Style, matrixTransform)
		container.TryAppend(fillColorTag.ApplyColorTransform(colorTransform))
	}

	container.TryAppend(NewDrawTag(path.Shape, settings.GlobalSettings.ASSDrawingScale))

	return container
}
