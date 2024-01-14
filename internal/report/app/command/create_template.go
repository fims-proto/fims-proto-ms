package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/template"

	"github.com/google/uuid"
)

type CreateTemplateCmd struct {
	TemplateId uuid.UUID
	SobId      uuid.UUID
	Creater    uuid.UUID
	Class      string
	Title      string
	Tables     []TableCmd
}

type CreateTemplateHandler struct {
	repo domain.Repository
}

func NewcreateTemplateHandler(repo domain.Repository) CreateTemplateHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CreateTemplateHandler{repo: repo}
}

func (h CreateTemplateHandler) Handle(ctx context.Context, cmd CreateTemplateCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.createTemplate(txCtx, cmd)
	})
}

func (h CreateTemplateHandler) createTemplate(ctx context.Context, cmd CreateTemplateCmd) error {
	tables, err := prepareTables(ctx, h.repo, cmd.SobId, cmd.Tables)
	newTemplate, err := template.New(
		cmd.TemplateId,
		cmd.SobId,
		cmd.Class,
		cmd.Title,
		tables,
	)
	if err != nil {
		return err
	}
	return h.repo.CreateTemplate(ctx, newTemplate)
}
