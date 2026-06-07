package http

import (
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/report/app"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	app *app.Application
}

func NewHandler(app *app.Application) Handler {
	if app == nil {
		panic("nil application")
	}
	return Handler{app: app}
}

func InitRouter(h Handler, api huma.API) {
	reportsApi := data.WithTag(api, "reports")

	huma.Register(reportsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/reports", Summary: "List reports"}, h.SearchReports)
	huma.Register(reportsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/report/template", Summary: "Get report template by class"}, h.ReadReportTemplateByClass)
	huma.Register(reportsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/report", Summary: "Get report instance by class and period"}, h.ReadReportByClassAndPeriod)
	huma.Register(reportsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/report/generate", Summary: "Generate report instance"}, h.GenerateReport)
	huma.Register(reportsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/report/{reportId}/recalculate", Summary: "Recalculate report", DefaultStatus: 204}, h.RecalculateReport)
	huma.Register(reportsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/report/{reportId}/regenerate", Summary: "Regenerate report", DefaultStatus: 204}, h.RegenerateReport)
	huma.Register(reportsApi, huma.Operation{Method: "PATCH", Path: "/api/v1/sob/{sobId}/report/{reportId}", Summary: "Update report structure", DefaultStatus: 204}, h.UpdateReport)
}
