package aoi

import (
	"errors"
	"math"

	"github.com/golang/geo/r3"
)

//Manager aoi manager
type Manager struct {
	tail      EntityIF
	entityMap map[EntityID]EntityIF
}

func (m *Manager) onEnterRange(ownerEntity, otherEntity EntityIF) {
	otherID := ownerEntity.GetEntityID()
	_, ok := m.entityMap[otherID]
	if !ok {
		return
	}
	ownerEntity.onEnter(otherEntity)
}

func (m *Manager) Move(entity EntityIF, point *r3.Vector) {
	entityID := entity.GetEntityID()
	_, ok := m.entityMap[entityID]
	if !ok {
		return
	}
	entity.moveTo(point)
}

func (m *Manager) onLeaveRange(ownerEntity, otherEntity EntityIF) {
	otherID := ownerEntity.GetEntityID()
	_, ok := m.entityMap[otherID]
	if !ok {
		return
	}
	ownerEntity.onLeave(otherEntity)
}

//Remove ...
func (m *Manager) Remove(entity EntityIF) {
	entityID := entity.GetEntityID()
	_, ok := m.entityMap[entityID]
	if !ok {
		return
	}
	entity.removeSelf()
	delete(m.entityMap, entity.GetEntityID())
}

//Add ...
func (m *Manager) Add(entity EntityIF) {
	entityID := entity.GetEntityID()
	_, ok := m.entityMap[entityID]
	if !ok {
		return
	}
	m.entityMap[entityID] = entity
}

//AddRangeOfEntity ...
func (m *Manager) AddRangeOfEntity(entity EntityIF, vector *r3.Vector, entityType EntityType) {
	entityID := entity.GetEntityID()
	_, ok := m.entityMap[entityID]
	if !ok {
		return
	}
	entity.addRange(vector, entityType)
}

//EntitiesInRange ...
func (m *Manager) EntitiesInRange(entity EntityIF, radius float64, isIncludeSelf bool) ([]EntityID, error) {
	entityID := entity.GetEntityID()
	_, ok := m.entityMap[entityID]
	if !ok {
		return nil, errors.New("entity not found ")
	}
	ret := make([]EntityID, 0)
	if isIncludeSelf {
		ret = append(ret, entityID)
	}
	centerPos := entity.GetPosition()

	getCursorFn := func(cursorTag int8) EntityIF {
		var cursor EntityIF
		if cursorTag == 1 {
			cursor = entity.GetNode().nextXEntity()
		}
		if cursorTag == 2 {
			cursor = entity.GetNode().preXEntity()
		}
		return cursor
	}

	fn := func(cursorTag int8) {
		cursor := getCursorFn(cursorTag)
		for {
			if cursor == nil {
				break
			}
			//todo check entity type equal
			pos := cursor.GetPosition()
			xDist, zDist := math.Abs(pos.X-centerPos.X), math.Abs(pos.Z-centerPos.Z)
			if xDist <= radius && zDist <= radius {
				ret = append(ret, cursor.GetEntityID())
			}
			cursor = getCursorFn(cursorTag)
		}
	}
	fn(1)
	fn(2)
	return ret, nil
}
