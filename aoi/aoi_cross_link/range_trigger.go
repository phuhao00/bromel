package crosslink

import (
	"math"
	"unsafe"
)

// RangeTriggerNode is a node that belongs to a RangeTrigger
// each RangeTrigger has 2 RangeTriggerNode, one positive and one negative, represents 2 sides of a range
// it implements CLPosImp as CLNode does
type RangeTriggerNode struct {
	CLNode

	rangeX    CLPosValType
	rangeZ    CLPosValType
	oldRangeX CLPosValType
	oldRangeZ CLPosValType

	myTrigger  *RangeTrigger
	isPositive bool
}

var g_offsetRangeTriggerNode uintptr

func initOffsetRangeTriggerNode() {
	dummy := (*RangeTriggerNode)(unsafe.Pointer(&g_offsetRangeTriggerNode))
	g_offsetRangeTriggerNode = uintptr(unsafe.Pointer(&dummy.CLNode)) - uintptr(unsafe.Pointer(dummy))
}

// RangeTrigger is a representation of a trigger range
type RangeTrigger struct {
	ownerEntNode *EntityListNode

	upperBound *RangeTriggerNode
	lowerBound *RangeTriggerNode

	posX CLPosValType
	posZ CLPosValType
	oldX CLPosValType
	oldZ CLPosValType

	rangeID RangeIDValType
	eventTp EventValType
}

// newRangeTriggerNode creates a new RangeTriggerNode
// isPositive means is upper bound, or the position value is bigger
// the given parameters rangeX, rangeZ should be greater than 0
func newRangeTriggerNode(rTrigger *RangeTrigger, isPositive bool, rangeX CLPosValType, rangeZ CLPosValType) *RangeTriggerNode {
	var node = new(RangeTriggerNode)
	node.isPositive = isPositive
	node.myTrigger = rTrigger

	if isPositive {
		node.rangeX = rangeX
		node.rangeZ = rangeZ
	} else {
		node.rangeX = -rangeX
		node.rangeZ = -rangeZ
	}
	node.nodeType = CLNODE_TRIGGER
	// now it's not in a proper cross linked list
	return node
}

func (thisNode *RangeTriggerNode) initialShuffle(oldX CLPosValType, oldZ CLPosValType) {
	// don't trigger anything
	thisNode.oldRangeX = 0
	thisNode.oldRangeZ = 0
	shuffleXThenZ(thisNode, oldX, oldZ)
	thisNode.oldRangeX = thisNode.rangeX
	thisNode.oldRangeZ = thisNode.rangeZ
}

func (thisNode *RangeTriggerNode) x() CLPosValType {
	return thisNode.myTrigger.x() + thisNode.rangeX
}

func (thisNode *RangeTriggerNode) z() CLPosValType {
	return thisNode.myTrigger.z() + thisNode.rangeZ
}

func (thisNode *RangeTriggerNode) getEntityID() EntityIDValType {
	return thisNode.myTrigger.ownerEntNode.aoiNode.entID
}

// this 2 methods overrides the methods of 'parent' CLNode
func (thisNode *RangeTriggerNode) crossedX(otherNode CLNodeImp, positiveCross bool, otherOldX CLPosValType, otherOldZ CLPosValType) {
	if thisNode.getEntityID() == otherNode.getEntityID() {
		return
	}

	wasInZRange := thisNode.myTrigger.wasInZRange(otherOldZ, Abs(thisNode.oldRangeZ))
	if !wasInZRange {
		return
	}

	isEnter := bool(thisNode.isPositive != positiveCross)
	if isEnter {
		if thisNode.myTrigger.isInXRange(otherNode.x(), Abs(thisNode.rangeX)) &&
			thisNode.myTrigger.isInZRange(otherNode.z(), Abs(thisNode.rangeZ)) {
			thisNode.myTrigger.triggerEnter(otherNode)
		}
	} else {
		if thisNode.myTrigger.wasInXRange(otherOldX, Abs(thisNode.oldRangeX)) {
			thisNode.myTrigger.triggerLeave(otherNode)
		}
	}
}

func (thisNode *RangeTriggerNode) crossedZ(otherNode CLNodeImp, positiveCross bool, otherOldX CLPosValType, otherOldZ CLPosValType) {
	if thisNode.getEntityID() == otherNode.getEntityID() {
		return
	}

	wasInXRange := thisNode.myTrigger.wasInXRange(otherOldX, Abs(thisNode.oldRangeX))
	if !wasInXRange {
		return
	}

	isEnter := bool(thisNode.isPositive != positiveCross)
	if isEnter {
		if thisNode.myTrigger.isInZRange(otherNode.z(), Abs(thisNode.rangeZ)) {
			thisNode.myTrigger.triggerEnter(otherNode)
		}
	} else {
		if thisNode.myTrigger.wasInXRange(otherOldX, Abs(thisNode.oldRangeX)) &&
			thisNode.myTrigger.wasInZRange(otherOldZ, Abs(thisNode.oldRangeZ)) {
			thisNode.myTrigger.triggerLeave(otherNode)
		}
	}
}

func (thisNode *RangeTriggerNode) setRange(rangeX CLPosValType, rangeZ CLPosValType) {
	oldX, oldZ := thisNode.x(), thisNode.z()
	thisNode.rangeX, thisNode.rangeZ = rangeX, rangeZ
	shuffleXThenZ(thisNode, oldX, oldZ)
	thisNode.oldRangeX, thisNode.oldRangeZ = thisNode.rangeX, thisNode.rangeZ
}

// newRangeTrigger is a method for creating new RangeTrigger
func newRangeTrigger(ownerNode *EntityListNode, rangeX CLPosValType, rangeZ CLPosValType, rangeID RangeIDValType, eventTp EventValType) *RangeTrigger {
	var trigger = new(RangeTrigger)
	trigger.ownerEntNode = ownerNode
	trigger.upperBound = newRangeTriggerNode(trigger, true, rangeX, rangeZ)
	trigger.lowerBound = newRangeTriggerNode(trigger, false, rangeX, rangeZ)
	trigger.posX, trigger.posZ = ownerNode.x(), ownerNode.z()
	trigger.oldX, trigger.oldZ = trigger.posX, trigger.posZ
	trigger.rangeID = rangeID
	trigger.eventTp = eventTp
	return trigger
}

func (thisTrigger *RangeTrigger) x() CLPosValType {
	return thisTrigger.posX
}

func (thisTrigger *RangeTrigger) z() CLPosValType {
	return thisTrigger.posZ
}

func (thisTrigger *RangeTrigger) owner() CLNodeImp {
	return thisTrigger.ownerEntNode
}

func (thisTrigger *RangeTrigger) isInXRange(x CLPosValType, rangeX CLPosValType) bool {
	return (thisTrigger.posX-rangeX) < x && x <= (thisTrigger.posX+rangeX)
}

func (thisTrigger *RangeTrigger) isInZRange(z CLPosValType, rangeZ CLPosValType) bool {
	return (thisTrigger.posZ-rangeZ) < z && z <= (thisTrigger.posZ+rangeZ)
}

func (thisTrigger *RangeTrigger) wasInXRange(x CLPosValType, rangeX CLPosValType) bool {
	return (thisTrigger.oldX-rangeX) < x && x <= (thisTrigger.oldX+rangeX)
}

func (thisTrigger *RangeTrigger) wasInZRange(z CLPosValType, rangeZ CLPosValType) bool {
	return (thisTrigger.oldZ-rangeZ) < z && z <= (thisTrigger.oldZ+rangeZ)
}

func (thisTrigger *RangeTrigger) triggerEnter(entering CLNodeImp) {
	thisTrigger.ownerEntNode.aoiNode.onEntityEnterRange(entering.getEntityID(), thisTrigger.rangeID)
}

func (thisTrigger *RangeTrigger) triggerLeave(leaving CLNodeImp) {
	thisTrigger.ownerEntNode.aoiNode.onEntityLeaveRange(leaving.getEntityID(), thisTrigger.rangeID)
}

func (thisTrigger *RangeTrigger) setRange(rangeX CLPosValType, rangeZ CLPosValType) {
	rangeX = Max(rangeX, 0.00000001)
	rangeZ = Max(rangeZ, 0.00000001)
	thisTrigger.upperBound.setRange(rangeX, rangeZ)
	thisTrigger.lowerBound.setRange(rangeX, rangeZ)
}

func (thisTrigger *RangeTrigger) insert() {
	cursor := thisTrigger.ownerEntNode.pNextX
	cursor.getCLNodePtr().insertBeforeX(&thisTrigger.lowerBound.CLNode)
	cursor.getCLNodePtr().insertBeforeX(&thisTrigger.upperBound.CLNode)

	cursor = thisTrigger.ownerEntNode.pNextZ
	cursor.getCLNodePtr().insertBeforeZ(&thisTrigger.lowerBound.CLNode)
	cursor.getCLNodePtr().insertBeforeZ(&thisTrigger.upperBound.CLNode)

	thisTrigger.upperBound.initialShuffle(thisTrigger.upperBound.x(), thisTrigger.upperBound.z())
	thisTrigger.lowerBound.initialShuffle(thisTrigger.lowerBound.x(), thisTrigger.lowerBound.z())
}

func (thisTrigger *RangeTrigger) removeMyself() {
	thisTrigger.posZ = math.MaxFloat32
	thisTrigger.shuffleZ()
	thisTrigger.upperBound.removeFromRangeList()
	thisTrigger.lowerBound.removeFromRangeList()
}

func (thisTrigger *RangeTrigger) moveCenterToPos(tgtX CLPosValType, tgtZ CLPosValType) {
	thisTrigger.posX, thisTrigger.posZ = tgtX, tgtZ
	thisTrigger.shuffleXThenZ()
}

func (thisTrigger *RangeTrigger) shuffleZ() {
	shuffleZ(thisTrigger.upperBound, thisTrigger.oldX+thisTrigger.upperBound.rangeX, thisTrigger.oldZ+thisTrigger.upperBound.rangeZ)
	shuffleZ(thisTrigger.lowerBound, thisTrigger.oldX+thisTrigger.lowerBound.rangeX, thisTrigger.oldZ+thisTrigger.lowerBound.rangeZ)
	thisTrigger.oldZ = thisTrigger.posZ
}

func (thisTrigger *RangeTrigger) shuffleXThenZ() {
	upOldX := thisTrigger.oldX + thisTrigger.upperBound.rangeX
	upOldZ := thisTrigger.oldZ + thisTrigger.upperBound.rangeZ
	loOldX := thisTrigger.oldX + thisTrigger.lowerBound.rangeX
	loOldZ := thisTrigger.oldZ + thisTrigger.lowerBound.rangeZ

	// let the range expand first then shrink
	if thisTrigger.oldX < thisTrigger.posX {
		shuffleX(thisTrigger.upperBound, upOldX, upOldZ)
		shuffleX(thisTrigger.lowerBound, loOldX, loOldZ)
	} else {
		shuffleX(thisTrigger.lowerBound, loOldX, loOldZ)
		shuffleX(thisTrigger.upperBound, upOldX, upOldZ)
	}

	if thisTrigger.oldZ < thisTrigger.posZ {
		shuffleZ(thisTrigger.upperBound, upOldX, upOldZ)
		shuffleZ(thisTrigger.lowerBound, loOldX, upOldZ)
	} else {
		shuffleZ(thisTrigger.lowerBound, loOldX, upOldZ)
		shuffleZ(thisTrigger.upperBound, upOldX, upOldZ)
	}

	thisTrigger.oldX, thisTrigger.oldZ = thisTrigger.posX, thisTrigger.posZ
}
