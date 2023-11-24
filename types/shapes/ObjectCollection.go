package shapes

import (
	"golang.org/x/exp/maps"
)

type ObjectCollection map[uint16]ObjectDefinition

func (o ObjectCollection) Clone() ObjectCollection {
	m := make(ObjectCollection)
	maps.Copy(m, o)
	return m
}

func (o ObjectCollection) Get(objectId uint16) ObjectDefinition {
	return o[objectId]
}

func (o ObjectCollection) Add(def ObjectDefinition) {
	if _, ok := o[def.GetObjectId()]; ok {
		panic("object already exists")
	}
	o[def.GetObjectId()] = def
}
