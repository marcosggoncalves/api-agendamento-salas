package response

import "ApiSup/internal/models"

type Grade struct {
	Sala     *models.Sala              `json:"sala"`
	Horarios []models.SalaGradeHorario `json:"horarios"`
}
