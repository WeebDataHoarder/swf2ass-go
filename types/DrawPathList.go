package types

type DrawPathList []DrawPath

func (l DrawPathList) Merge(b DrawPathList) DrawPathList {
	newList := make(DrawPathList, 0, len(l)+len(b))
	newList = append(newList, l...)
	newList = append(newList, b...)
	return newList
}
