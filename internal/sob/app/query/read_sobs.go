package query

import (
	"context"

	"github.com/pkg/errors"
)

type SobsReadModel interface {
	AllSobs(ctx context.Context) ([]Sob, error)
	SobById(ctx context.Context, sobId string) (Sob, error)
}

type ReadSobsHandler struct {
	readModel SobsReadModel
}

func NewReadSobsHandler(readModel SobsReadModel) ReadSobsHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	return ReadSobsHandler{
		readModel: readModel,
	}
}

func (r ReadSobsHandler) HandleReadAll(ctx context.Context) ([]Sob, error) {
	sobs, err := r.readModel.AllSobs(ctx)
	if err != nil {
		return []Sob{}, errors.Wrapf(err, "read all sobs failed")
	}
	return sobs, nil
}

func (r ReadSobsHandler) HandleReadById(ctx context.Context, sobId string) (Sob, error) {
	sob, err := r.readModel.SobById(ctx, sobId)
	if err != nil {
		return Sob{}, errors.Wrapf(err, "read sob failed")
	}
	return sob, nil
}
