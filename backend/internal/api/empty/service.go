package empty

import (
	"context"
	"lytemp/internal/domain/models"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Inject( /*userService IUserService*/ ) {
	// s.userService = userService
}

// This should be in the internal/domain/dtos/requests package
// Its here to not pollute the requests directory
type CreateCity struct { // requests.CreateCity
	Name string `json:"name" validate:"required"`
}

func (s *Service) Example(ctx context.Context, request CreateCity) (models.User, error) {
	city := models.User{
		Name: request.Name,
	}

	return city, s.repository.Create(ctx, &city)
}
