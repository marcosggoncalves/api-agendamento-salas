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

type SalaController struct {
	service services.SalaService
}

func NewSalaController(salaService services.SalaService) *SalaController {
	return &SalaController{
		service: salaService,
	}
}

func (controller *SalaController) Listagem(c echo.Context) error {
	salas, err := controller.service.Listagem(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.ERRO_LISTAGEM_REGISTRO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, salas)
}

func (controller *SalaController) ListagemSimples(c echo.Context) error {
	salas, err := controller.service.ListagemSimples()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.ERRO_LISTAGEM_REGISTRO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, salas)
}

func (controller *SalaController) Created(c echo.Context) error {
	var sala models.Sala
	if err := c.Bind(&sala); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.BODY_FALHA, Description: err.Error()})
	}

	if err := c.Validate(sala); err != nil {
		return config.ValidationErrors(c, err)
	}

	if err := controller.service.Novo(&sala); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.CADASTRO_FALHA_INSERCAO, Description: err.Error()})
	}

	return c.JSON(http.StatusCreated, response.Success{Message: constants.CADASTRO_REALIZADO})
}

func (controller *SalaController) Updated(c echo.Context) error {
	var updated models.Sala
	if err := c.Bind(&updated); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.BODY_FALHA, Description: err.Error()})
	}

	if err := c.Validate(updated); err != nil {
		return config.ValidationErrors(c, err)
	}

	sala, err := controller.service.Editar(&updated)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error{Message: constants.REGISTRO_NAO_ENCONTRADO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, response.SuccessBody{Message: constants.CADASTRO_ALTERADO, Body: sala})
}

func (controller *SalaController) Detalhar(c echo.Context) error {
	hash := c.Param("hash")
	if hash == "" {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.ID_NAO_INFORMADO, Description: constants.DATA_NAO_INFORMADO})
	}

	data := c.Param("data")
	if data == "" {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.DATA_NAO_INFORMADO, Description: constants.DATA_NAO_INFORMADO})
	}

	sala, err := controller.service.Detalhar(hash, data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.ERRO_LISTAGEM_REGISTRO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, sala)
}

func (controller *SalaController) SalasDisponiveis(c echo.Context) error {
	data := c.Param("data")
	if data == "" {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.DATA_NAO_INFORMADO, Description: constants.DATA_NAO_INFORMADO})
	}

	sala, err := controller.service.SalasDisponiveis(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.ERRO_LISTAGEM_REGISTRO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, sala)
}

func (controller *SalaController) Deleted(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.ID_NAO_INFORMADO, Description: err.Error()})
	}

	if err := controller.service.Deletar(id); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.CADASTRO_FALHA_EXCLUSAO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, response.Success{Message: constants.CADASTRO_EXCLUIDO})
}

// Metodos Salvar e Excluir "Nova Reserva"
func (controller *SalaController) NovaReserva(c echo.Context) error {
	var reserva models.EfetuarReserva
	if err := c.Bind(&reserva); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.BODY_FALHA, Description: err.Error()})
	}

	if err := c.Validate(reserva); err != nil {
		return config.ValidationErrors(c, err)
	}

	if err := controller.service.NovaReserva(&reserva); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.CADASTRO_FALHA_INSERCAO, Description: err.Error()})
	}

	return c.JSON(http.StatusCreated, response.Success{Message: constants.RESERVA_REALIZADA})
}

func (controller *SalaController) DeletedReserva(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.ID_NAO_INFORMADO, Description: err.Error()})
	}

	if err := controller.service.DeletarReserva(id); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.CADASTRO_FALHA_EXCLUSAO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, response.Success{Message: constants.RESERVA_CANCELADA})
}

func (controller *SalaController) VisualizarReserva(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Message: constants.ID_NAO_INFORMADO, Description: err.Error()})
	}

	reserva, err := controller.service.VisualizarReserva(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.CADASTRO_FALHA_EXCLUSAO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, reserva)
}

func (controller *SalaController) VisualizarAgenda(c echo.Context) error {
	agendamentos, err := controller.service.VisualizarAgenda()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Message: constants.ERRO_LISTAGEM_REGISTRO, Description: err.Error()})
	}

	return c.JSON(http.StatusOK, agendamentos)
}
