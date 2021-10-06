package crosslink

import (
	"fmt"
	"testing"
)

type gameEntity struct {
	x float32
	z float32

	id uint32
}

func (thisEnt *gameEntity) AoiCLX() CLPosValType {
	return CLPosValType(thisEnt.x)
}

func (thisEnt *gameEntity) AoiCLZ() CLPosValType {
	return CLPosValType(thisEnt.z)
}

func (thisEnt *gameEntity) AoiCLEntityID() EntityIDValType {
	return EntityIDValType(thisEnt.id)
}

func (thisEnt *gameEntity) onEnterRange(enteringID EntityIDValType, rangeID RangeIDValType) {
	fmt.Println("onEnterRange", thisEnt.id, enteringID, rangeID)
}

func (thisEnt *gameEntity) onLeaveRange(leavingID EntityIDValType, rangeID RangeIDValType) {
	fmt.Println("onLeaveRange", thisEnt.id, leavingID, rangeID)
}

func TestALL(t *testing.T) {
	entities := []*gameEntity{}
	var testEnt *gameEntity
	for i := 40; i < 60; i++ {
		for j := 40; j < 60; j++ {
			newEnt := new(gameEntity)
			newEnt.x = float32(i)
			newEnt.z = float32(j)
			newEnt.id = uint32(1000000 + i*1000 + j)
			entities = append(entities, newEnt)
			if i == 50 && j == 50 {
				testEnt = newEnt
			}
		}
	}
	fmt.Println("step: entities created", len(entities))

	aoiSpace := NewAOISpaceCL()

	for _, ent := range entities {
		aoiSpace.AddEntity(ent)
	}

	fmt.Println("step: entities added", len(entities))

	for _, ent := range entities {
		aoiSpace.AddRangeOfEntity(ent, 1.1, 1.1, EVENT_ALL)
	}

	fmt.Println("step: range added", len(entities))

	ids, err := aoiSpace.EntitiesInRange(testEnt, 1.1, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("EntitiesInRange", testEnt.id, ids)

	aoiSpace.MoveEntity(testEnt, testEnt.AoiCLX()+1, testEnt.AoiCLZ()+1)
	ids, err = aoiSpace.EntitiesInRange(testEnt, 1.1, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("EntitiesInRange after move", testEnt.id, ids)

	aoiSpace.RemoveEntity(testEnt)
	fmt.Println("RemoveEntity done")

	aoiSpace.AddEntity(testEnt)
	fmt.Println("Re-AddEntity done")

	ids, err = aoiSpace.EntitiesInRange(testEnt, 1.1, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("EntitiesInRange after re-add", testEnt.id, ids)
}
