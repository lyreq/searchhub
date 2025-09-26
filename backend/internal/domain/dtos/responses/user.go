package responses

import "lytemp/internal/domain/models"

type UserDTO struct {
	ID              uint              `json:"id"`
	Name            string            `json:"name"`
	Surname         string            `json:"surname"`
	ShortName       string            `json:"short_name"`
	LongName        string            `json:"long_name"`
	Phone           string            `json:"phone"`
	Email           string            `json:"email"`
	Username        string            `json:"username"`
	Status          models.StatusType `json:"status"`
	IsRandevuActive bool              `json:"is_randevu_active"`
	AccountStatus   bool              `json:"account_status"`
	IsVerified      bool              `json:"is_verified"`
	Credit          uint              `json:"credit"`
	Debt            float64           `json:"debt"`
}
