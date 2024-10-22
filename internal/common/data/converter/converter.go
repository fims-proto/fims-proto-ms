package converter

import "github.com/google/uuid"

func BOsToPOs[B any, P any](bos []*B, convertFn func(bo *B) P) []P {
	var pos []P
	for _, bo := range bos {
		po := convertFn(bo)
		pos = append(pos, po)
	}
	return pos
}

func POsToBOs[B any, P any](pos []P, convertFn func(po P) (*B, error)) ([]*B, error) {
	var bos []*B
	for _, po := range pos {
		bo, err := convertFn(po)
		if err != nil {
			return nil, err
		}
		bos = append(bos, bo)
	}
	return bos, nil
}

func POsToDTOs[P any, D any](pos []P, convert func(po P) D) []D {
	var dtos []D
	for _, po := range pos {
		dto := convert(po)
		dtos = append(dtos, dto)
	}
	return dtos
}

func DTOsToVOs[D any, V any](pos []D, convert func(po D) V) []V {
	// same logic as persistence object to data transfer object
	return POsToDTOs(pos, convert)
}

func UUIDToPtr(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func UUIDFromPtr(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}
