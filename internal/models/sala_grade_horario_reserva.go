package models

type SalaGradeHorarioReserva struct {
	BaseModel
	DepartamentoID     uint             `json:"departamento_id" validate:"required"`
	SalaGradeHorarioID uint             `json:"sala_grade_horario_id" validate:"required"`
	DataReserva        string           `json:"data_reserva" validate:"required"`
	Status             Status           `json:"status"`
	Nome               string           `json:"nome" validate:"required"`
	Departamento       Departamento     `json:"departamento" gorm:"foreignKey:DepartamentoID" validate:"-"`
	Horario            SalaGradeHorario `json:"horario" gorm:"foreignKey:SalaGradeHorarioID" validate:"-"`
}

type EfetuarReserva struct {
	BaseModel
	DepartamentoID uint   `json:"departamento_id" validate:"required"`
	Horarios       []uint `json:"horarios" validate:"required"`
	DataReserva    string `json:"data_reserva" validate:"required"`
	Status         Status `json:"status"`
	Nome           string `json:"nome" validate:"required" `
}

type Reserva struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Start      string `json:"start"`
	End        string `json:"end"`
	Color      string `json:"color"`
	DiaInteiro int    `json:"AllDay"`
}

func (SalaGradeHorarioReserva) TableName() string {
	return "sala_grade_horario_reserva"
}

func (EfetuarReserva) TableName() string {
	return "sala_grade_horario_reserva"
}
