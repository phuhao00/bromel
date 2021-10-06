package crosslink

import (
	"math"
	"unsafe"
)

// CLPosImp is a position x-z interface
type CLPosImp interface {
	x() CLPosValType
	z() CLPosValType
}

// CLNodeImp is a general node interface
type CLNodeImp interface {
	CLPosImp

	isTriggerNode() bool
	isEntity() bool
	order() uint8

	// here otherNode CLNodeImp should be a pointer to some struct, but not value of struct
	crossedX(otherNode CLNodeImp, positiveCross bool, otherOldX CLPosValType, otherOldZ CLPosValType)
	crossedZ(otherNode CLNodeImp, positiveCross bool, otherOldX CLPosValType, otherOldZ CLPosValType)

	moveToPrevX()
	moveToNextX()
	moveToPrevZ()
	moveToNextZ()

	clNodeType() CLNodeValType

	nodeID() unsafe.Pointer

	prevX() CLNodeImp
	prevZ() CLNodeImp
	nextX() CLNodeImp
	nextZ() CLNodeImp

	isTail() bool

	getCLNodePtr() *CLNode

	getEntityID() EntityIDValType
}

// CLNode is a base Node struct for the cross linked list
// it implements interface CLNodeImp
type CLNode struct {
	pPrevX   *CLNode
	pNextX   *CLNode
	pPrevZ   *CLNode
	pNextZ   *CLNode
	nodeType CLNodeValType
}

// implement interface CLNodeImp: begin

func getParentInst(ptr *CLNode) CLNodeImp {
	if ptr == nil {
		return nil
	}
	if ptr.nodeType == CLNODE_ENTITY {
		return (*EntityListNode)(unsafe.Pointer((uintptr)(unsafe.Pointer(ptr)) - g_offsetEntityListNode))
	} else if ptr.nodeType == CLNODE_TRIGGER {
		return (*RangeTriggerNode)(unsafe.Pointer((uintptr)(unsafe.Pointer(ptr)) - g_offsetEntityListNode))
	}
	return nil
}

func (thisNode *CLNode) prevX() CLNodeImp {
	return getParentInst(thisNode.pPrevX)
}

func (thisNode *CLNode) prevZ() CLNodeImp {
	return getParentInst(thisNode.pPrevZ)
}

func (thisNode *CLNode) nextX() CLNodeImp {
	return getParentInst(thisNode.pNextX)
}

func (thisNode *CLNode) nextZ() CLNodeImp {
	return getParentInst(thisNode.pNextZ)
}

func (thisNode *CLNode) clNodeType() CLNodeValType {
	return thisNode.nodeType
}

func (thisNode *CLNode) isTriggerNode() bool {
	return thisNode.nodeType == CLNODE_TRIGGER
}

func (thisNode *CLNode) isEntity() bool {
	return thisNode.nodeType == CLNODE_ENTITY
}

func (thisNode *CLNode) order() uint8 {
	return 0
}

func (thisNode *CLNode) crossedX(otherNode CLNodeImp, positiveCross bool, otherOldX CLPosValType, otherOldZ CLPosValType) {
	panic("override this func")
}

func (thisNode *CLNode) crossedZ(otherNode CLNodeImp, positiveCross bool, otherOldX CLPosValType, otherOldZ CLPosValType) {
	panic("override this func")
}

func (thisNode *CLNode) moveToPrevX() {
	if thisNode.pNextX != nil {
		thisNode.pNextX.pPrevX = thisNode.pPrevX
	}

	thisNode.pPrevX.pNextX = thisNode.pNextX

	thisNode.pNextX = thisNode.pPrevX
	thisNode.pPrevX = thisNode.pPrevX.pPrevX

	if thisNode.pPrevX != nil {
		thisNode.pPrevX.pNextX = thisNode
	}
	thisNode.pNextX.pPrevX = thisNode
}

func (thisNode *CLNode) moveToNextX() {

	if thisNode.pPrevX != nil {
		thisNode.pPrevX.pNextX = thisNode.pNextX
	}

	thisNode.pNextX.pPrevX = thisNode.pPrevX

	thisNode.pPrevX = thisNode.pNextX
	thisNode.pNextX = thisNode.pNextX.pNextX

	if thisNode.pNextX != nil {
		thisNode.pNextX.pPrevX = thisNode
	}
	thisNode.pPrevX.pNextX = thisNode
}

func (thisNode *CLNode) moveToPrevZ() {
	if thisNode.pNextZ != nil {
		thisNode.pNextZ.pPrevZ = thisNode.pPrevZ
	}

	thisNode.pPrevZ.pNextZ = thisNode.pNextZ

	thisNode.pNextZ = thisNode.pPrevZ
	thisNode.pPrevZ = thisNode.pPrevZ.pPrevZ

	if thisNode.pPrevZ != nil {
		thisNode.pPrevZ.pNextZ = thisNode
	}
	thisNode.pNextZ.pPrevZ = thisNode
}

func (thisNode *CLNode) moveToNextZ() {

	if thisNode.pPrevZ != nil {
		thisNode.pPrevZ.pNextZ = thisNode.pNextZ
	}

	thisNode.pNextZ.pPrevZ = thisNode.pPrevZ

	thisNode.pPrevZ = thisNode.pNextZ
	thisNode.pNextZ = thisNode.pNextZ.pNextZ

	if thisNode.pNextZ != nil {
		thisNode.pNextZ.pPrevZ = thisNode
	}
	thisNode.pPrevZ.pNextZ = thisNode
}

func (thisNode *CLNode) nodeID() unsafe.Pointer {
	return unsafe.Pointer(thisNode)
}

func (thisNode *CLNode) isTail() bool {
	return false
}

func (thisNode *CLNode) getCLNodePtr() *CLNode {
	return thisNode
}

// implement interface CLNodeImp: end

// removeFromRangeList remove this node from linked list, process pointers only
func (thisNode *CLNode) removeFromRangeList() {
	if thisNode.pPrevX != nil {
		thisNode.pPrevX.pNextX = thisNode.pNextX
	}

	if thisNode.pNextX != nil {
		thisNode.pNextX.pPrevX = thisNode.pPrevX
	}

	if thisNode.pPrevZ != nil {
		thisNode.pPrevZ.pNextZ = thisNode.pNextZ
	}

	if thisNode.pNextZ != nil {
		thisNode.pNextZ.pPrevZ = thisNode.pPrevZ
	}

	thisNode.pPrevX = nil
	thisNode.pNextX = nil
	thisNode.pPrevZ = nil
	thisNode.pNextZ = nil
}

func (thisNode *CLNode) insertBeforeX(newNode *CLNode) {
	if thisNode.pPrevX != nil {
		thisNode.pPrevX.pNextX = newNode
	}

	newNode.pPrevX = thisNode.pPrevX

	thisNode.pPrevX = newNode
	newNode.pNextX = thisNode
}

func (thisNode *CLNode) insertBeforeZ(newNode *CLNode) {
	if thisNode.pPrevZ != nil {
		thisNode.pPrevZ.pNextZ = newNode
	}

	newNode.pPrevZ = thisNode.pPrevZ

	thisNode.pPrevZ = newNode
	newNode.pNextZ = thisNode
}

// CLNodeTail is a terminator
type CLNodeTail struct {
	CLNode
}

func newCLNodeTail() *CLNodeTail {
	tail := new(CLNodeTail)
	tail.nodeType = CLNODE_TAIL
	return tail
}

func (thisNode *CLNodeTail) isTail() bool {
	return true
}

func (thisNode *CLNodeTail) x() CLPosValType {
	return math.MaxFloat32
}

func (thisNode *CLNodeTail) z() CLPosValType {
	return math.MaxFloat32
}

// shuffleX make this node in right position of X linked list
func shuffleX(thisNode CLNodeImp, oldX CLPosValType, oldZ CLPosValType) {
	thisPos := thisNode.x()
	for {
		prevNode := thisNode.prevX()
		if prevNode == nil {
			break
		}
		prevPos := prevNode.x()
		if thisPos < prevPos || (thisPos == prevPos && thisNode.order() < prevNode.order()) {
			if thisNode.isTriggerNode() && !prevNode.isTriggerNode() {
				thisNode.crossedX(prevNode, true, prevPos, prevNode.z())
			} else if !thisNode.isTriggerNode() && prevNode.isTriggerNode() {
				prevNode.crossedX(thisNode, false, oldX, oldZ)
			}
			thisNode.moveToPrevX()
		} else {
			break
		}
	}

	for {
		nextNode := thisNode.nextX()
		if nextNode == nil {
			break
		}
		nextPos := nextNode.x()
		if thisPos > nextPos || (thisPos == nextPos && thisNode.order() < nextNode.order()) {
			if thisNode.isTriggerNode() && !nextNode.isTriggerNode() {
				thisNode.crossedX(nextNode, false, nextPos, nextNode.z())
			} else if !thisNode.isTriggerNode() && nextNode.isTriggerNode() {
				nextNode.crossedX(thisNode, true, oldX, oldZ)
			}
			thisNode.moveToNextX()
		} else {
			break
		}
	}
}

// shuffleZ make this node in right position of Z linked list
func shuffleZ(thisNode CLNodeImp, oldX CLPosValType, oldZ CLPosValType) {
	thisPos := thisNode.z()
	for {
		prevNode := thisNode.prevZ()
		if prevNode == nil {
			break
		}
		prevPos := prevNode.z()
		if thisPos < prevPos || (thisPos == prevPos && thisNode.order() < prevNode.order()) {
			if thisNode.isTriggerNode() && !prevNode.isTriggerNode() {
				thisNode.crossedZ(prevNode, true, prevNode.x(), prevPos)
			} else if !thisNode.isTriggerNode() && prevNode.isTriggerNode() {
				prevNode.crossedZ(thisNode, false, oldX, oldZ)
			}
			thisNode.moveToPrevZ()
		} else {
			break
		}
	}

	for {
		nextNode := thisNode.nextZ()
		if nextNode == nil {
			break
		}
		nextPos := nextNode.z()
		if thisPos > nextPos || (thisPos == nextPos && thisNode.order() < nextNode.order()) {
			if thisNode.isTriggerNode() && !nextNode.isTriggerNode() {
				thisNode.crossedZ(nextNode, false, nextNode.x(), nextNode.z())
			} else if !thisNode.isTriggerNode() && nextNode.isTriggerNode() {
				nextNode.crossedZ(thisNode, true, oldX, oldZ)
			}
			thisNode.moveToNextZ()
		} else {
			break
		}
	}
}

// shuffleXThenZ make this node in right position of cross linked list
func shuffleXThenZ(thisNode CLNodeImp, oldX CLPosValType, oldZ CLPosValType) {
	shuffleX(thisNode, oldX, oldZ)
	shuffleZ(thisNode, oldX, oldZ)
}
