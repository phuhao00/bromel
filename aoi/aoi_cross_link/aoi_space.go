package crosslink

import (
	"errors"
	"math"
)

// EntityIDValType is a entity's ID
// the cross linked list manages many entities in a scene
type EntityIDValType uint32

// CLPosValType is the position's value type
type CLPosValType float32

// CLNodeValType is the type of any CLNode
type CLNodeValType uint8

const (
	CLNODE_TAIL CLNodeValType = iota
	CLNODE_ENTITY
	CLNODE_TRIGGER
)

// EventValType : RangeTrigger will trigg some specific event like enter or leave
type EventValType uint16

const (
	EVENT_ENTER EventValType = 1 << iota
	EVENT_LEAVE              = 1 << iota
	EVENT_ALL                = EVENT_ENTER | EVENT_LEAVE
)

// RangeIDValType : every RangeTrigger in a AOINode will have a unique rangeId
type RangeIDValType uint16

// Abs is a float64 Abs function wrapper
func Abs(f CLPosValType) CLPosValType {
	return CLPosValType(math.Abs(float64(f)))
}

// Max is a float64 function wrapper
func Max(m1 CLPosValType, m2 CLPosValType) CLPosValType {
	return CLPosValType(math.Max(float64(m1), float64(m2)))
}

const MAX_VALID_POS_IN_AOI CLPosValType = math.MaxFloat32

func IsValidAoiCLPosXZ(x CLPosValType, z CLPosValType) bool {
	return -MAX_VALID_POS_IN_AOI < x && x < MAX_VALID_POS_IN_AOI &&
		-MAX_VALID_POS_IN_AOI < z && z < MAX_VALID_POS_IN_AOI
}

type AOIEntityImp interface {
	// these function definitions make the position available
	AoiCLX() CLPosValType
	AoiCLZ() CLPosValType

	AoiCLEntityID() EntityIDValType

	onEnterRange(enteringID EntityIDValType, rangeID RangeIDValType)

	onLeaveRange(enteringID EntityIDValType, rangeID RangeIDValType)
}

type AOISpaceImp interface {
}

// AOISpaceCL manages a space, which has a independant AOI (area of interest) space, spaces won't affect each other
type AOISpaceCL struct {
	tailNode *CLNodeTail

	entityNodesMap map[EntityIDValType]*EntityAOINode
}

func initPtrs() {
	initOffsetEntityListNode()
	initOffsetRangeTriggerNode()
}

// NewAOISpaceCL creates a new AOISpaceCL with parameters
func NewAOISpaceCL() *AOISpaceCL {
	initPtrs()
	space := new(AOISpaceCL)
	space.tailNode = newCLNodeTail()
	space.entityNodesMap = make(map[EntityIDValType]*EntityAOINode)

	return space
}

func (thisSpace *AOISpaceCL) onEntityEnterRange(whoID EntityIDValType, enteringID EntityIDValType, rangeID RangeIDValType) {
	who, ok := thisSpace.entityNodesMap[whoID]
	if !ok {
		return
	}
	who.entityImp.onEnterRange(enteringID, rangeID)
}

func (thisSpace *AOISpaceCL) onEntityLeaveRange(whoID EntityIDValType, leavingID EntityIDValType, rangeID RangeIDValType) {
	who, ok := thisSpace.entityNodesMap[whoID]
	if !ok {
		return
	}
	who.entityImp.onLeaveRange(leavingID, rangeID)
}

// AddEntity : add an Entity to this AOI space
func (thisSpace *AOISpaceCL) AddEntity(aoiEntity AOIEntityImp) error {
	entID := aoiEntity.AoiCLEntityID()
	_, ok := thisSpace.entityNodesMap[entID]
	if ok {
		return errors.New("entity already in aoi")
	}

	x := aoiEntity.AoiCLX()
	z := aoiEntity.AoiCLZ()

	if !IsValidAoiCLPosXZ(x, z) {
		return errors.New("entity pos not valid")
	}

	aoiNode := newEntityAOINode(thisSpace, aoiEntity, x, z)
	thisSpace.entityNodesMap[entID] = aoiNode

	thisSpace.tailNode.insertBeforeX(&aoiNode.entListNode.CLNode)
	thisSpace.tailNode.insertBeforeZ(&aoiNode.entListNode.CLNode)
	shuffleXThenZ(aoiNode.entListNode, math.MaxFloat32, math.MaxFloat32)

	return nil
}

// AddRangeOfEntity : add range after add entity
func (thisSpace *AOISpaceCL) AddRangeOfEntity(aoiEntity AOIEntityImp, rangeX CLPosValType, rangeZ CLPosValType, eventType EventValType) error {
	entID := aoiEntity.AoiCLEntityID()
	aoiNode, ok := thisSpace.entityNodesMap[entID]
	if !ok {
		return errors.New("entity not in aoi")
	}

	aoiNode.addRange(rangeX, rangeZ, eventType)
	return nil
}

// RemoveEntity : remove an Entity from this AOI space
func (thisSpace *AOISpaceCL) RemoveEntity(aoiEntity AOIEntityImp) error {
	entID := aoiEntity.AoiCLEntityID()
	aoiNode, ok := thisSpace.entityNodesMap[entID]
	if !ok {
		return errors.New("entity not in aoi")
	}

	aoiNode.removeMyself()
	delete(thisSpace.entityNodesMap, entID)

	return nil
}

// MoveEntity : move a entity to a new position, and auto recalc the aoi
func (thisSpace *AOISpaceCL) MoveEntity(aoiEntity AOIEntityImp, tgtX CLPosValType, tgtZ CLPosValType) error {
	entID := aoiEntity.AoiCLEntityID()
	aoiNode, ok := thisSpace.entityNodesMap[entID]
	if !ok {
		return errors.New("entity not in aoi")
	}

	if !IsValidAoiCLPosXZ(tgtX, tgtZ) {
		return errors.New("entity pos not valid")
	}

	aoiNode.moveToPos(tgtX, tgtZ)
	return nil
}

// EntitiesInRange : get entities in specified range of this AOI space
func (thisSpace *AOISpaceCL) EntitiesInRange(aoiEntity AOIEntityImp, r CLPosValType, includeThis bool) ([]EntityIDValType, error) {
	entID := aoiEntity.AoiCLEntityID()
	aoiNode, ok := thisSpace.entityNodesMap[entID]
	if !ok {
		return nil, errors.New("entity not in aoi")
	}

	if r <= 0 {
		return nil, errors.New("r should be greater than 0")
	}

	listNode := aoiNode.entListNode
	centerPosX, centerPoxZ := listNode.posX, listNode.posZ

	res := []EntityIDValType{}
	if includeThis {
		res = append(res, entID)
	}

	var xDist, zDist CLPosValType

	xDist, zDist = 0, 0
	cursor := listNode.nextX()
	for {
		if cursor == nil || xDist > r {
			break
		}
		if cursor.isEntity() {
			xDist, zDist = Abs(cursor.x()-centerPosX), Abs(cursor.z()-centerPoxZ)
			if xDist <= r && zDist <= r {
				res = append(res, cursor.getEntityID())
			}
		}
		cursor = cursor.nextX()
	}

	xDist, zDist = 0, 0
	cursor = listNode.prevX()
	for {
		if cursor == nil || xDist > r {
			break
		}
		if cursor.isEntity() {
			xDist, zDist = Abs(cursor.x()-centerPosX), Abs(cursor.z()-centerPoxZ)
			if xDist <= r && zDist <= r {
				res = append(res, cursor.getEntityID())
			}
		}
		cursor = cursor.prevX()
	}
	return res, nil
}
