package services

import (
	"ApiSup/internal/models"
	"ApiSup/internal/repositories"
	"ApiSup/pkg/pagination"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

type SalaService interface {
	Listagem(c echo.Context) (*pagination.Pagination, error)
	ListagemSimples() ([]models.SalaListagem, error)
	Novo(user *models.Sala) error
	Editar(updated *models.Sala) (*models.Sala, error)
	Deletar(id int) error
	Detalhar(hash string, data string) (*models.Sala, error)

	// Metodos Salvar e Excluir "Nova Reserva"]
	SalasDisponiveis(data string) ([]models.Sala, error)
	VisualizarAgenda() ([]models.Reserva, error)
	VisualizarReserva(id uint) (*models.SalaGradeHorarioReserva, error)
	NovaReserva(reserve *models.EfetuarReserva) error
	DeletarReserva(id int) error
}

type salaService struct {
	repository repositories.SalaRepository
}

func NewSalaService(repository repositories.SalaRepository) SalaService {
	return &salaService{
		repository: repository,
	}
}

func criarGradeHorarios(sala *models.Sala) ([]models.SalaGradeHorario, error) {
	format := "15:04:05"

	horarioInicial, err := time.Parse(format, sala.HorarioIniFuncionamento)
	if err != nil {
		return nil, err
	}

	horarioFinal, err := time.Parse(format, sala.HorarioFimFuncionamento)
	if err != nil {
		return nil, err
	}

	var horarios []models.SalaGradeHorario

	for horarioInicialGrade := horarioInicial; horarioInicialGrade.Before(horarioFinal); horarioInicialGrade = horarioInicialGrade.Add(time.Duration(sala.IntervaloPorAgendamento) * time.Minute) {

		horarioFinalGrade := horarioInicialGrade.Add(time.Duration(sala.IntervaloPorAgendamento) * time.Minute)

		if horarioFinalGrade.After(horarioFinal) {
			horarioFinalGrade = horarioFinal
		}

		horarios = append(horarios, models.SalaGradeHorario{
			SalaID:  sala.ID,
			Inicial: horarioInicialGrade.Format(format),
			Final:   horarioFinalGrade.Format(format),
		})
	}

	horarios = append(horarios, models.SalaGradeHorario{
		SalaID:     sala.ID,
		Inicial:    sala.HorarioIniFuncionamento,
		Final:      sala.HorarioFimFuncionamento,
		DiaInteiro: 1,
	})

	return horarios, nil
}

func (s *salaService) Listagem(c echo.Context) (*pagination.Pagination, error) {
	return s.repository.Listagem(c)
}

func (s *salaService) ListagemSimples() ([]models.SalaListagem, error) {
	return s.repository.ListagemSimples()
}

func (s *salaService) SalasDisponiveis(data string) ([]models.Sala, error) {
	return s.repository.SalasDisponiveis(data)
}

func (s *salaService) Novo(sala *models.Sala) error {
	/// Montar grade de horários
	horarios, err := criarGradeHorarios(sala)
	if err != nil {
		return err
	}

	/// Atruibuir grade de horários na sala cadastrada
	sala.Grade = horarios

	return s.repository.Novo(sala)
}

func (s *salaService) Editar(updated *models.Sala) (*models.Sala, error) {
	/// Buscar informações da sala
	sala, err := s.repository.DetalharByID(updated.ID)
	if err != nil {
		return nil, err
	}

	if updated.GerarGrade == "generate" {
		updated.ID = sala.ID
		horarios, err := criarGradeHorarios(updated)
		if err != nil {
			return nil, fmt.Errorf("failed to create schedule: %w", err)
		}
		updated.Grade = horarios
	}

	return s.repository.Editar(updated)
}

func (s *salaService) Deletar(id int) error {
	return s.repository.Deletar(id)
}

func (s *salaService) Detalhar(hash string, data string) (*models.Sala, error) {
	return s.repository.DetalharByHashData(hash, data)
}

// Metodos Salvar e Excluir "Nova Reserva"
func (s *salaService) NovaReserva(reserve *models.EfetuarReserva) error {
	err := s.repository.NovaReserva(reserve)
	if err != nil {
		return err
	}

	return nil
}

func (s *salaService) DeletarReserva(id int) error {
	return s.repository.DeletarReserva(id)
}

func (s *salaService) VisualizarReserva(id uint) (*models.SalaGradeHorarioReserva, error) {
	return s.repository.VisualizarReserva(id)
}

func (s *salaService) VisualizarAgenda() ([]models.Reserva, error) {
	return s.repository.VisualizarAgenda()
}
