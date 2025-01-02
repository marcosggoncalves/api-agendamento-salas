package services

import (
	"ApiSup/internal/models"
	"ApiSup/internal/repositories"
	"ApiSup/pkg/mapear/constants"
	"ApiSup/pkg/mapear/request"
	"ApiSup/pkg/pagination"
	"errors"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UsuarioService interface {
	Authenticate(body request.Login) (*models.Usuario, error)
	Listagem(c echo.Context) (*pagination.Pagination, error)
	Novo(user *models.Usuario) error
	Editar(updated *models.Usuario) (*models.Usuario, error)
	Deletar(id int) error
}

type usuarioService struct {
	repository repositories.UsuarioRepository
}

func NewUsuarioService(repository repositories.UsuarioRepository) UsuarioService {
	return &usuarioService{repository: repository}
}

func (s *usuarioService) Authenticate(body request.Login) (*models.Usuario, error) {
	user, err := s.repository.GetUserByCPF(body.CPF)
	if err != nil {
		return nil, errors.New(constants.USUARIO_ENCONTRADO)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Senha), []byte(body.Senha))
	if err != nil {
		return nil, errors.New(constants.SENHA_INVALIDA)
	}

	return user, nil
}

func (s *usuarioService) Listagem(c echo.Context) (*pagination.Pagination, error) {
	return s.repository.Listagem(c)
}

func (s *usuarioService) Novo(user *models.Usuario) error {
	return s.repository.Novo(user)
}

func (s *usuarioService) Editar(updated *models.Usuario) (*models.Usuario, error) {
	if updated.Senha != "" {
		hashedSenha, err := models.Criptografia(updated.Senha)
		if err != nil {
			return nil, err
		}

		updated.Senha = hashedSenha
	}

	return s.repository.Editar(updated)
}

func (s *usuarioService) Deletar(id int) error {
	return s.repository.Deletar(id)
}
