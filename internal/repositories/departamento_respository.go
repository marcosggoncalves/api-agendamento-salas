package repositories

import (
	"ApiSup/internal/config"
	"ApiSup/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type DepartamentoRepository interface {
	Listagem(c echo.Context) ([]models.Departamento, error)
	Novo(departamento *models.Departamento) error
	Editar(updated *models.Departamento) (*models.Departamento, error)
	Deletar(id int) error
}

type departamentoRepository struct {
	db *gorm.DB
}

func NewDepartamentoRepository() DepartamentoRepository {
	return &departamentoRepository{db: config.DB}
}

func (r *departamentoRepository) Listagem(c echo.Context) ([]models.Departamento, error) {
	var departamentos []models.Departamento

	if err := r.db.Order("descricao asc").Find(&departamentos).Error; err != nil {
		return nil, err
	}

	return departamentos, nil
}

func (r *departamentoRepository) Novo(departamento *models.Departamento) error {
	return r.db.Save(departamento).Error
}

func (r *departamentoRepository) Editar(updated *models.Departamento) (*models.Departamento, error) {
	departamento := new(models.Departamento)
	if err := r.db.First(departamento, updated.ID).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(departamento).Updates(updated).Error; err != nil {
		return nil, err
	}
	return departamento, nil
}

func (r *departamentoRepository) Deletar(id int) error {
	return r.db.Delete(&models.Departamento{}, id).Error
}
