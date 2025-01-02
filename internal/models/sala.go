package models

import (
	"ApiSup/pkg/hashing"

	"gorm.io/gorm"
)

type Sala struct {
	BaseModel
	Hash                    string             `json:"hash"`
	GerarGrade              string             `json:"gerar_grade" gorm:"-"`
	Nome                    string             `json:"nome" validate:"required"`
	Color                   string             `json:"color" validate:"required"`
	Descricao               string             `json:"descricao" validate:"required"`
	IntervaloPorAgendamento int                `json:"intervalo_por_agendamento" validate:"required"`
	HorarioIniFuncionamento string             `json:"horario_ini_funcionamento" validate:"required"`
	HorarioFimFuncionamento string             `json:"horario_fim_funcionamento" validate:"required"`
	Grade                   []SalaGradeHorario `json:"horarios" gorm:"foreignKey:SalaID"`
}

type SalaListagem struct {
	ID         int    `json:"id"`
	Hash       string `json:"hash"`
	GerarGrade string `json:"gerar_grade" gorm:"-"`
	Nome       string `json:"nome" validate:"required"`
	Color      string `json:"color" validate:"required"`
}

func (SalaListagem) TableName() string {
	return "sala"
}

func (Sala) TableName() string {
	return "sala"
}

func (u *Sala) BeforeCreate(tx *gorm.DB) (err error) {
	hashToken, err := hashing.GenerateToken(16)
	if err != nil {
		return err
	}

	u.Hash = hashToken

	return nil
}
