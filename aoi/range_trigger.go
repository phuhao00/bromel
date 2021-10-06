package aoi

import "github.com/golang/geo/r3"

type RangeTrigger struct {
	rangeID        uint16
	oldPos, newPos *r3.Vector
	ownerEntity    EntityIF
	upperBound     *RangeTriggerNode
	lowerBound     *RangeTriggerNode
}

func (t *RangeTrigger) name() {

}

type RangeTriggerNode struct {
	node       *Node
	Cur, Old   r3.Vector
	isPositive bool
}

func (n *RangeTriggerNode) moveToPrevX() {
	panic("implement me")
}

func (n *RangeTriggerNode) moveToNextX() {
	panic("implement me")
}

func (n *RangeTriggerNode) isTriggerNode() bool {
	panic("implement me")
}

func (n *RangeTriggerNode) order() uint8 {
	panic("implement me")
}

func (n *RangeTriggerNode) onEnter(entityIF EntityIF) {
	panic("implement me")
}

func (n *RangeTriggerNode) onLeave(entityIF EntityIF) {
	panic("implement me")
}

func (n *RangeTriggerNode) GetEntityID() EntityID {
	panic("implement me")
}

func (n *RangeTriggerNode) GetRadius() float64 {
	panic("implement me")
}

func (n *RangeTriggerNode) moveTo(vector *r3.Vector) {
	panic("implement me")
}

func (n *RangeTriggerNode) removeSelf() {
	panic("implement me")
}

func (n *RangeTriggerNode) addRange(vector *r3.Vector, entityType EntityType) {
	panic("implement me")
}

func (n *RangeTriggerNode) GetNode() *Node {
	panic("implement me")
}

func (n *RangeTriggerNode) GetPosition() *r3.Vector {
	panic("implement me")
}

func (n *RangeTriggerNode) crossed(otherEntity EntityIF, positiveCross bool) {
	panic("implement me")
}

func (n *RangeTriggerNode) initialShuffle() {
	n.Old.X = 0
	n.Old.Z = 0
	shuffle(n)
	n.Old.X = n.Cur.X
	n.Old.Z = n.Cur.Z
}
