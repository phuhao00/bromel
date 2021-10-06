package aoi

import "github.com/golang/geo/r3"

//EntityBase aoi entity base
type EntityBase struct {
	Vec           *r3.Vector
	AoiId         EntityID
	Radius        float64
	rangeTriggers map[uint16]RangeTrigger
	nextRangeID   uint16
	node          *Node
	nodeType      EntityType // 代表类型，主要是区分AOI矩形的边界和Entity本身

}

func (e *EntityBase) moveToPrevX() {
	panic("implement me")
}

func (e *EntityBase) moveToNextX() {
	panic("implement me")
}

func (e *EntityBase) isTriggerNode() bool {
	panic("implement me")
}

func (e *EntityBase) order() uint8 {
	panic("implement me")
}

func (e *EntityBase) crossed(otherEntity EntityIF, positiveCross bool) {
	if e.GetEntityID() == otherEntity.GetEntityID() {
		return
	}

	//if thisNode.getEntityID() == otherNode.getEntityID() {
	//		return
	//	}
	//
	//	wasInZRange := thisNode.myTrigger.wasInZRange(otherOldZ, Abs(thisNode.oldRangeZ))
	//	if !wasInZRange {
	//		return
	//	}
	//
	//	isEnter := bool(thisNode.isPositive != positiveCross)
	//	if isEnter {
	//		if thisNode.myTrigger.isInXRange(otherNode.x(), Abs(thisNode.rangeX)) &&
	//			thisNode.myTrigger.isInZRange(otherNode.z(), Abs(thisNode.rangeZ)) {
	//			thisNode.myTrigger.triggerEnter(otherNode)
	//		}
	//	} else {
	//		if thisNode.myTrigger.wasInXRange(otherOldX, Abs(thisNode.oldRangeX)) {
	//			thisNode.myTrigger.triggerLeave(otherNode)
	//		}
	//	}

	panic("implement me")
}

//Enter ...
func (e *EntityBase) onEnter(entityIF EntityIF) {
	panic("should  override this")
}

//Leave ...
func (e *EntityBase) onLeave(entityIF EntityIF) {
	panic("should  override this")
}

func (e *EntityBase) GetEntityID() EntityID {
	return e.AoiId
}

func (e *EntityBase) GetRadius() float64 {
	return e.Radius
}

func (e *EntityBase) moveTo(vector *r3.Vector) {
	panic("should  override this")

}

func (e *EntityBase) removeSelf() {
	panic("should  override this")

}

func (e *EntityBase) addRange(vector *r3.Vector, entityType EntityType) {
	panic("should  override this")

}

func (e *EntityBase) GetNode() *Node {
	return e.node
}

func (e *EntityBase) GetPosition() *r3.Vector {
	return e.Vec
}
