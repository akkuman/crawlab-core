package controllers

import (
	"errors"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var ExportController ActionController

func getExportActions() []Action {
	ctx := newExportContext()
	return []Action{
		{
			Method:      http.MethodPost,
			Path:        "/:type",
			HandlerFunc: ctx.postExport,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:type/:id",
			HandlerFunc: ctx.getExport,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:type/:id/download",
			HandlerFunc: ctx.getExportDownload,
		},
	}
}

type exportContext struct {
	csvSvc interfaces.ExportService
}

func (ctx *exportContext) postExport(c *gin.Context) {
	exportType := c.Param("type")
	exportTarget := c.Query("target")
	exportFilter, _ := GetFilter(c)

	var exportId string
	var err error
	switch exportType {
	case constants.ExportTypeCsv:
		exportId, err = ctx.csvSvc.Export(exportType, exportTarget, exportFilter)
	default:
		HandleErrorBadRequest(c, errors.New(fmt.Sprintf("invalid export type: %s", exportType)))
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, exportId)
}

func (ctx *exportContext) getExport(c *gin.Context) {
	exportType := c.Param("type")
	exportId := c.Param("id")

	var export interfaces.Export
	var err error
	switch exportType {
	case constants.ExportTypeCsv:
		export, err = ctx.csvSvc.GetExport(exportId)
	default:
		HandleErrorBadRequest(c, errors.New(fmt.Sprintf("invalid export type: %s", exportType)))
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, export)
}

func (ctx *exportContext) getExportDownload(c *gin.Context) {
	exportType := c.Param("type")
	exportId := c.Param("id")

	var export interfaces.Export
	var err error
	switch exportType {
	case constants.ExportTypeCsv:
		export, err = ctx.csvSvc.GetExport(exportId)
	default:
		HandleErrorBadRequest(c, errors.New(fmt.Sprintf("invalid export type: %s", exportType)))
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", export.GetDownloadPath()))
	c.Header("Content-Length", strconv.Itoa(len(export.GetDownloadPath())))
	c.File(export.GetDownloadPath())
}

func newExportContext() *exportContext {
	return &exportContext{}
}
