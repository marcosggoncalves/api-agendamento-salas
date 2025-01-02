package repositories

import (
	"ApiSup/internal/config"
	"ApiSup/internal/models"
	"ApiSup/pkg/pagination"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UsuarioRepository interface {
	GetUserByCPF(cpf string) (*models.Usuario, error)
	Listagem(c echo.Context) (*pagination.Pagination, error)
	Novo(user *models.Usuario) error
	Editar(updated *models.Usuario) (*models.Usuario, error)
	Deletar(id int) error
}

type usuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository() UsuarioRepository {
	return &usuarioRepository{db: config.DB}
}

func (r *usuarioRepository) GetUserByCPF(cpf string) (*models.Usuario, error) {
	var usuario models.Usuario
	if err := r.db.Where("cpf = ?", cpf).First(&usuario).Error; err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *usuarioRepository) Listagem(c echo.Context) (*pagination.Pagination, error) {
	var usuarios []models.UsuarioView

	paginations, err := pagination.Paginate(c, r.db, &usuarios, nil)

	if err != nil {
		return nil, err
	}

	return paginations, nil
}

func (r *usuarioRepository) Novo(usuario *models.Usuario) error {
	return r.db.Save(usuario).Error
}

func (r *usuarioRepository) Editar(updated *models.Usuario) (*models.Usuario, error) {
	usuario := new(models.Usuario)
	if err := r.db.First(usuario, updated.ID).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&usuario).Updates(updated).Error; err != nil {
		return nil, err
	}
	return usuario, nil
}

func (r *usuarioRepository) Deletar(id int) error {
	return r.db.Delete(&models.Usuario{}, id).Error
}
