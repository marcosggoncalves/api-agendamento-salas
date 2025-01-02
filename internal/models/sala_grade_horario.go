package models

type Status string

const (
	Reservado Status = "Reservado"
	Vencido   Status = "Reserva Vencida"
	Cancelado Status = "Cancelado"
)

type SalaGradeHorario struct {
	BaseModel
	SalaID     uint   `json:"sala_id"`
	Inicial    string `json:"inicial"`
	Final      string `json:"final"`
	DiaInteiro int    `json:"dia_inteiro"`
}

func (SalaGradeHorario) TableName() string {
	return "sala_grade_horario"
}
