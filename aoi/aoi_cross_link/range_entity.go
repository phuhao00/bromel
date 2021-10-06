package crosslink

import (
	"math"
	"unsafe"
)

// EntityListNodeImp
type EntityListNodeImp interface {
	getAOINode()
}

// AOIReceiver

// EntityListNode
type EntityListNode struct {
	CLNode
	aoiNode *EntityAOINode
	posX    CLPosValType
	posZ    CLPosValType
}

func (thisNode *EntityListNode) x() CLPosValType {
	return thisNode.posX
}

func (thisNode *EntityListNode) z() CLPosValType {
	return thisNode.posZ
}

var g_offsetEntityListNode uintptr

func initOffsetEntityListNode() {
	dummy := (*EntityListNode)(unsafe.Pointer(&g_offsetEntityListNode))
	g_offsetEntityListNode = uintptr(unsafe.Pointer(&dummy.CLNode)) - uintptr(unsafe.Pointer(dummy))
}

func (thisNode *EntityListNode) getEntityID() EntityIDValType {
	return thisNode.aoiNode.entID
}

func (thisNode *EntityListNode) removeMyself(oldZ CLPosValType) {
	thisNode.posZ = math.MaxFloat32
	shuffleZ(thisNode, thisNode.posX, oldZ)
	thisNode.removeFromRangeList()
}

func (thisNode *EntityListNode) moveToPos(tgtX CLPosValType, tgtZ CLPosValType) {
	oldX, oldZ := thisNode.posX, thisNode.posZ
	thisNode.posX, thisNode.posZ = tgtX, tgtZ
	shuffleXThenZ(thisNode, oldX, oldZ)
}

// EntityAOINode
type EntityAOINode struct {
	entID       EntityIDValType
	entListNode *EntityListNode
	triggers    map[RangeIDValType]*RangeTrigger
	nextRangeID RangeIDValType
	space       *AOISpaceCL
	entityImp   AOIEntityImp //entityIF
}

func newEntityListNode(entAOINode *EntityAOINode, x CLPosValType, z CLPosValType) *EntityListNode {
	eln := new(EntityListNode)
	eln.aoiNode = entAOINode
	eln.posX = x
	eln.posZ = z
	eln.nodeType = CLNODE_ENTITY
	return eln
}

func newEntityAOINode(space *AOISpaceCL, entityImp AOIEntityImp, x CLPosValType, z CLPosValType) *EntityAOINode {
	aoiNode := new(EntityAOINode)
	aoiNode.entID = entityImp.AoiCLEntityID()
	aoiNode.entListNode = newEntityListNode(aoiNode, x, z)
	aoiNode.triggers = make(map[RangeIDValType]*RangeTrigger)
	aoiNode.space = space
	aoiNode.entityImp = entityImp
	return aoiNode
}

func (thisNode *EntityAOINode) onEntityEnterRange(entID EntityIDValType, rangeID RangeIDValType) {
	thisNode.space.onEntityEnterRange(thisNode.entID, entID, rangeID)
}

func (thisNode *EntityAOINode) onEntityLeaveRange(entID EntityIDValType, rangeID RangeIDValType) {
	thisNode.space.onEntityLeaveRange(thisNode.entID, entID, rangeID)
}

func (thisNode *EntityAOINode) addRange(rangeX CLPosValType, rangeZ CLPosValType, eventType EventValType) {
	trigger := newRangeTrigger(thisNode.entListNode, rangeX, rangeZ, thisNode.nextRangeID, eventType)
	trigger.insert()
	thisNode.triggers[thisNode.nextRangeID] = trigger
	thisNode.nextRangeID++
}

func (thisNode *EntityAOINode) removeMyself() {
	// after calling this, this node is not usable anymore
	oldZ := thisNode.entListNode.posZ
	thisNode.entListNode.removeMyself(oldZ)
	for _, trigger := range thisNode.triggers {
		trigger.removeMyself()
	}
	thisNode.triggers = nil
}

func (thisNode *EntityAOINode) moveToPos(tgtX CLPosValType, tgtZ CLPosValType) {
	thisNode.entListNode.moveToPos(tgtX, tgtZ)

	for _, trigger := range thisNode.triggers {
		trigger.moveCenterToPos(tgtX, tgtZ)
	}
}
