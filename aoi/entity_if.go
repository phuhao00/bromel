package aoi

import "github.com/golang/geo/r3"

type CallBackFn func()

type EntityIF interface {
	onEnter(entityIF EntityIF)
	onLeave(entityIF EntityIF)
	GetEntityID() EntityID
	GetRadius() float64
	moveTo(vector *r3.Vector)
	removeSelf()
	addRange(vector *r3.Vector, entityType EntityType)
	GetNode() *Node
	GetPosition() *r3.Vector
	crossed(otherEntity EntityIF, positiveCross bool)
	order() uint8
	isTriggerNode() bool
	moveToPrevX()
	moveToNextX()
}

type EntityID uint64

type EntityType uint16
