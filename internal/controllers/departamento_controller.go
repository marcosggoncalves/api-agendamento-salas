package controllers

import (
	"ApiSup/internal/config"
	"ApiSup/internal/models"
	"ApiSup/internal/services"
	"ApiSup/pkg/mapear/constants"
	"ApiSup/pkg/mapear/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DepartamentoController struct {
	service services.DepartamentoService
}

func NewDepartamentoController(departamentoService services.DepartamentoService) *DepartamentoController {
	return &DepartamentoController{
		service: departamentoService,
	}
}

func (controller *DepartamentoController) Listagem(c echo.Context) error {
	departamentos, err := controller.service.Listagem(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.ERRO_LISTAGEM_REGISTRO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, departamentos)
}

func (controller *DepartamentoController) Created(c echo.Context) error {
	var departamento models.Departamento
	if err := c.Bind(&departamento); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.BODY_FALHA, Description: err.Error()})
	}

	if err := c.Validate(departamento); err != nil {
		return config.ValidationErrors(c, err)
	}

	if err := controller.service.Novo(&departamento); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.CADASTRO_FALHA_INSERCAO, Description: err.Error()})
	}

	return c.JSON(http.StatusCreated, response.Success{Message: constants.CADASTRO_REALIZADO})
}

func (controller *DepartamentoController) Updated(c echo.Context) error {
	var updated models.Departamento
	if err := c.Bind(&updated); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.BODY_FALHA, Description: err.Error()})
	}

	if err := c.Validate(updated); err != nil {
		return config.ValidationErrors(c, err)
	}

	departamento, err := controller.service.Editar(&updated)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error{Message: constants.REGISTRO_NAO_ENCONTRADO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, response.SuccessBody{Message: constants.CADASTRO_ALTERADO, Body: departamento})
}

func (controller *DepartamentoController) Deleted(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.ID_NAO_INFORMADO, Description: err.Error()})
	}

	if err := controller.service.Deletar(id); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.CADASTRO_FALHA_EXCLUSAO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, response.Success{Message: constants.CADASTRO_EXCLUIDO})
}
