package mappers

import (
	"lytemp/internal/domain/dtos/responses"
	"lytemp/internal/domain/models"
)

func MapUserToDTO(u models.User) responses.UserDTO {
	return responses.UserDTO{
		ID:      u.ID,
		Name:    u.Name,
		Surname: u.Surname,

		Phone:    u.Phone,
		Email:    u.Email,
		Username: u.Username,
		Status:   models.StatusType(u.Status),
	}
}
