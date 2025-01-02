package models

type Departamento struct {
	BaseModel
	Descricao string `json:"descricao" validate:"required"`
}

func (Departamento) TableName() string {
	return "departamento"
}
