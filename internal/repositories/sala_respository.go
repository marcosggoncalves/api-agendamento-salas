package repositories

import (
	"ApiSup/internal/config"
	"ApiSup/internal/models"
	"ApiSup/pkg/mapear/constants"
	"ApiSup/pkg/pagination"
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SalaRepository interface {
	Listagem(c echo.Context) (*pagination.Pagination, error)
	ListagemSimples() ([]models.SalaListagem, error)
	Novo(sala *models.Sala) error
	DetalharByHashData(hash string, data string) (*models.Sala, error)
	DetalharByID(id uint) (*models.Sala, error)
	Editar(updated *models.Sala) (*models.Sala, error)
	Deletar(id int) error

	// Metodos Salvar e Excluir "Nova Reserva"
	SalasDisponiveis(data string) ([]models.Sala, error)
	VisualizarAgenda() ([]models.Reserva, error)
	VisualizarReserva(id uint) (*models.SalaGradeHorarioReserva, error)
	NovaReserva(reserve *models.EfetuarReserva) error
	DeletarReserva(id int) error
}

type salaRepository struct {
	db *gorm.DB
}

func NewSalaRepository() SalaRepository {
	return &salaRepository{db: config.DB}
}

func (r *salaRepository) Listagem(c echo.Context) (*pagination.Pagination, error) {
	var salas []models.Sala

	paginations, err := pagination.Paginate(c, r.db, &salas, nil, "Grade")

	if err != nil {
		return nil, err
	}

	return paginations, nil
}

func (r *salaRepository) ListagemSimples() ([]models.SalaListagem, error) {
	var salas []models.SalaListagem

	if err := r.db.Order("nome asc").Find(&salas).Error; err != nil {
		return nil, err
	}

	return salas, nil
}

func (r *salaRepository) SalasDisponiveis(data string) ([]models.Sala, error) {
	var salas []models.Sala

	// Carregar as salas com a "Grade" associada
	if err := r.db.Find(&salas).Error; err != nil {
		return nil, err
	}

	for i := range salas {
		// Subconsulta para verificar se o dia inteiro está reservado
		diaInteiroReservado := r.db.Table("sala_grade_horario sg").
			Select("1").
			Joins("LEFT JOIN sala_grade_horario_reserva shr ON shr.sala_grade_horario_id = sg.id").
			Joins("LEFT JOIN sala s ON s.id = sg.sala_id").
			Where("s.id = ?", salas[i].ID).
			Where("sg.dia_inteiro = ?", 1).
			Where("shr.status = ?", models.Reservado).
			Where("shr.data_reserva = ?", data)

		// Subconsulta para verificar IDs não reservados para a data
		reservaNaoReservada := r.db.Table("sala_grade_horario_reserva sg").
			Select("sala_grade_horario_id").
			Joins("LEFT JOIN sala_grade_horario sh ON sh.id = sg.sala_grade_horario_id").
			Joins("LEFT JOIN sala s ON s.id = sh.sala_id").
			Where("sg.data_reserva = ?", data).
			Where("s.id = ?", salas[i].ID).
			Where("sg.status = ?", models.Reservado).
			Where("sg.deleted_at IS NULL")

		var diasJaReservados int64
		if err := reservaNaoReservada.Count(&diasJaReservados).Error; err != nil {
			return nil, err
		}

		var grade []models.SalaGradeHorario
		query := r.db.
			Where("NOT EXISTS (?)", diaInteiroReservado).
			Where("id NOT IN (?)", reservaNaoReservada).
			Where("sala_id = ?", salas[i].ID).
			Where("deleted_at IS NULL")

		// Caso já tenha algum dia reservado, não pode aparecer a possibilidade de agendar o "Dia Inteiro"
		if diasJaReservados > 0 {
			query.Where("dia_inteiro = ?", 0)
		}

		if err := query.Find(&grade).Error; err != nil {
			return nil, err
		}

		salas[i].Grade = grade
	}

	return salas, nil
}

func (r *salaRepository) DetalharByHashData(hash string, data string) (*models.Sala, error) {
	sala := new(models.Sala)

	// Subconsulta para verificar se o dia inteiro está reservado
	diaInteiroReservado := r.db.Table("sala_grade_horario sg").
		Select("1").
		Joins("LEFT JOIN sala_grade_horario_reserva shr ON shr.sala_grade_horario_id = sg.id").
		Joins("LEFT JOIN sala s ON s.id = sg.sala_id").
		Where("sg.dia_inteiro = ?", 1).
		Where("s.hash = ?", hash).
		Where("shr.status = ?", models.Reservado).
		Where("shr.data_reserva = ?", data)

	// Subconsulta para verificar IDs não reservados para a data
	reservaNaoReservada := r.db.Table("sala_grade_horario_reserva sg").
		Select("sala_grade_horario_id").
		Joins("JOIN sala_grade_horario sh ON sh.id = sg.sala_grade_horario_id").
		Joins("JOIN sala s ON s.id = sh.sala_id").
		Where("sg.data_reserva = ?", data).
		Where("s.hash = ?", hash).
		Where("sg.status = ?", models.Reservado).
		Where("sg.deleted_at IS NULL")

	var diasJaReservados int64
	if err := reservaNaoReservada.Count(&diasJaReservados).Error; err != nil {
		return nil, err
	}

	if err := r.db.Where("hash = ?", hash).
		Preload("Grade", func(db *gorm.DB) *gorm.DB {
			query := db.
				Where("NOT EXISTS (?)", diaInteiroReservado).
				Where("id NOT IN (?)", reservaNaoReservada).
				Where("deleted_at IS NULL")

			// Caso já tenha algum dia reservado, não pode aparecer a possibilidade de agendar o "Dia Inteiro"
			if diasJaReservados > 0 {
				query.Where("dia_inteiro = ?", 0)
			}

			return query
		}).First(&sala).Error; err != nil {
		return nil, err
	}

	return sala, nil
}

func (r *salaRepository) DetalharByID(id uint) (*models.Sala, error) {
	sala := new(models.Sala)

	if err := r.db.First(sala, id).Error; err != nil {
		return nil, err
	}

	return sala, nil
}

func (r *salaRepository) Novo(sala *models.Sala) error {
	return r.db.Save(sala).Error
}

func (r *salaRepository) Editar(updated *models.Sala) (*models.Sala, error) {
	sala := new(models.Sala)

	transacao := r.db.Begin()

	if err := transacao.First(sala, updated.ID).Error; err != nil {
		transacao.Rollback()
		return nil, err
	}

	if err := transacao.Model(sala).Updates(updated).Error; err != nil {
		transacao.Rollback()
		return nil, err
	}

	if updated.GerarGrade == "generate" {
		if err := transacao.
			Where("sala_id = ?", sala.ID).
			Unscoped().
			Delete(&models.SalaGradeHorario{}).Error; err != nil {
			transacao.Rollback()
			return nil, err
		}

		if err := transacao.Save(updated.Grade).Error; err != nil {
			transacao.Rollback()
			return nil, err
		}
	}

	if err := transacao.Commit().Error; err != nil {
		return nil, err
	}
	return sala, nil
}

func (r *salaRepository) Deletar(id int) error {
	return r.db.Delete(&models.Sala{}, id).Error
}

// Metodos Salvar e Excluir "Reservas"
func (r *salaRepository) NovaReserva(reserve *models.EfetuarReserva) error {
	transacao := r.db.Begin()

	horarioVerificar := new(models.SalaGradeHorario)

	var reservas []models.SalaGradeHorarioReserva

	for _, horario := range reserve.Horarios {
		query := transacao.
			Joins("join sala_grade_horario_reserva s on s.sala_grade_horario_id = sala_grade_horario.id").
			Where("status", models.Reservado).
			Where("data_reserva", reserve.DataReserva).
			Where("sala_grade_horario_id", horario).
			First(horarioVerificar)

		if err := query.Error; err == nil {
			transacao.Rollback()
			return fmt.Errorf(constants.RESERVA_JA_REALIZADA+":%q-%q", horarioVerificar.Inicial, horarioVerificar.Final)
		}

		reservas = append(reservas, models.SalaGradeHorarioReserva{
			SalaGradeHorarioID: horario,
			Nome:               reserve.Nome,
			DepartamentoID:     reserve.DepartamentoID,
			DataReserva:        reserve.DataReserva,
			Status:             models.Reservado,
		})
	}

	if err := transacao.Save(reservas).Error; err != nil {
		transacao.Rollback()
		return err
	}

	if err := transacao.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *salaRepository) DeletarReserva(id int) error {
	if err := r.db.Model(&models.SalaGradeHorarioReserva{}).
		Where("id", id).
		Update("status", models.Cancelado).Error; err != nil {
		return err
	}

	return r.db.Delete(&models.SalaGradeHorarioReserva{}, id).Error
}

func (r *salaRepository) VisualizarReserva(id uint) (*models.SalaGradeHorarioReserva, error) {
	var reserva *models.SalaGradeHorarioReserva

	if err := r.db.
		Preload("Usuario.Departamento").
		Preload("Horario").
		First(&reserva, id).Error; err != nil {
		return nil, err
	}

	return reserva, nil
}

func (r *salaRepository) VisualizarAgenda() ([]models.Reserva, error) {
	query := `
		SELECT
			sr.id, 
			CONCAT('Reserva(', s.nome, ' - ', d.descricao, ') - ', sr.nome) AS title, 
			CONCAT(sr.data_reserva, 'T', LPAD(HOUR(sl.inicial), 2, '0'), ':', LPAD(MINUTE(sl.inicial), 2, '0'), ':00') as start,
			CONCAT(sr.data_reserva, 'T', LPAD(HOUR(sl.final), 2, '0'), ':', LPAD(MINUTE(sl.final), 2, '0'), ':00')  AS end,
			s.color
		FROM sala_grade_horario sl
		INNER JOIN sala_grade_horario_reserva sr ON sr.sala_grade_horario_id = sl.id
		INNER JOIN sala s ON s.id = sl.sala_id
		INNER JOIN departamento d ON d.id = sr.departamento_id
		AND sr.status = "Reservado"
		ORDER BY sr.data_reserva, sl.inicial, sl.final ASC
	`
	var reservas []models.Reserva

	result := config.DB.Raw(query).Scan(&reservas)

	if result.Error != nil {
		return nil, result.Error
	}

	return reservas, nil
}

func ValidarReservas() error {
	query := `
        UPDATE sala_grade_horario_reserva sr
        JOIN sala_grade_horario sl ON sr.sala_grade_horario_id = sl.id
        SET sr.status = 'Reserva Vencida'
        WHERE sr.status = 'Reservado'
  		AND (CONCAT(sr.data_reserva, ' ', sl.inicial) <= CURRENT_TIMESTAMP 
		OR CONCAT(sr.data_reserva, ' ', sl.final) <= CURRENT_TIMESTAMP);
    `
	result := config.DB.Exec(query)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
