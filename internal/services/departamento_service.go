package services

import (
	"ApiSup/internal/models"
	"ApiSup/internal/repositories"

	"github.com/labstack/echo/v4"
)

type DepartamentoService interface {
	Listagem(c echo.Context) ([]models.Departamento, error)
	Novo(departamento *models.Departamento) error
	Editar(updated *models.Departamento) (*models.Departamento, error)
	Deletar(id int) error
}

type departamentoService struct {
	repository repositories.DepartamentoRepository
}

func NewDepartamentoService(repository repositories.DepartamentoRepository) DepartamentoService {
	return &departamentoService{repository: repository}
}

func (s *departamentoService) Listagem(c echo.Context) ([]models.Departamento, error) {
	return s.repository.Listagem(c)
}

func (s *departamentoService) Novo(departamento *models.Departamento) error {
	return s.repository.Novo(departamento)
}

func (s *departamentoService) Editar(updated *models.Departamento) (*models.Departamento, error) {
	return s.repository.Editar(updated)
}

func (s *departamentoService) Deletar(id int) error {
	return s.repository.Deletar(id)
}
