package routes

import (
	controller "ApiSup/internal/controllers"
	"ApiSup/internal/middlewares"
	"ApiSup/internal/repositories"
	"ApiSup/internal/services"

	"github.com/labstack/echo/v4"
)

func InitializeRoutes(e *echo.Echo) {
	usuariosRepository := repositories.NewUsuarioRepository()
	usuariosService := services.NewUsuarioService(usuariosRepository)
	usuariosController := controller.NewUsuarioController(usuariosService)

	salaRepository := repositories.NewSalaRepository()
	salaService := services.NewSalaService(salaRepository)
	salaController := controller.NewSalaController(salaService)

	departamentoRepository := repositories.NewDepartamentoRepository()
	departamentoService := services.NewDepartamentoService(departamentoRepository)
	departamentoController := controller.NewDepartamentoController(departamentoService)

	auth := e.Group("/auth")
	auth.POST("/login", usuariosController.Login)

	publico := e.Group("/publico")
	publico.GET("/sala-agendamento/:hash/:data", salaController.Detalhar)
	publico.POST("/sala-agendamento", salaController.NovaReserva)
	publico.GET("/departamentos", departamentoController.Listagem)
	publico.GET("/salas", salaController.ListagemSimples)

	protected := e.Group("/api")
	protected.Use(middlewares.VerifyTokenHandler())

	protected.POST("/usuarios", usuariosController.Created)
	protected.GET("/usuarios", usuariosController.Listagem)
	protected.PUT("/usuarios", usuariosController.Updated)
	protected.DELETE("/usuarios/:id", usuariosController.Deleted)

	protected.POST("/salas", salaController.Created)
	protected.GET("/salas", salaController.Listagem)
	protected.GET("/salas-disponiveis/:data", salaController.SalasDisponiveis)
	protected.PUT("/salas", salaController.Updated)
	protected.DELETE("/salas/:id", salaController.Deleted)

	protected.GET("/salas-visualizar-agenda", salaController.VisualizarAgenda)
	protected.GET("/sala-visualizar/:id", salaController.VisualizarReserva)
	protected.DELETE("/sala-reserva/:id", salaController.DeletedReserva)

	protected.POST("/departamentos", departamentoController.Created)
	protected.GET("/departamentos", departamentoController.Listagem)
	protected.PUT("/departamentos", departamentoController.Updated)
	protected.DELETE("/departamentos/:id", departamentoController.Deleted)
}
