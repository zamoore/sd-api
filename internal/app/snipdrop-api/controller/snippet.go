package controller

import (
	"net/http"
	"snipdrop-rest-api/internal/app/snipdrop-api/model"
	"snipdrop-rest-api/internal/app/snipdrop-api/repository"
	"snipdrop-rest-api/internal/app/snipdrop-api/service"

	"github.com/gin-gonic/gin"
)

type SnippetController struct {
	Service *service.SnippetService
}

func (c *SnippetController) ListSnippets(ctx *gin.Context) {
	var params repository.SnippetQueryParams

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snippets, err := c.Service.ListSnippets(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, snippets)
}

func (c *SnippetController) CreateSnippet(ctx *gin.Context) {
	var snippet model.Snippet
	if err := ctx.ShouldBindJSON(&snippet); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.Service.CreateSnippet(ctx, snippet)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, snippet)
}

func (c *SnippetController) GetSnippet(ctx *gin.Context) {
	id := ctx.Param("id")
	snippet, err := c.Service.GetSnippet(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, snippet)
}

func (c *SnippetController) DeleteSnippet(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.Service.DeleteSnippet(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Snippet deleted"})
}
