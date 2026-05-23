package adm

import "context"

type Service struct {
	repo *Repo
}

func NewAdmService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (service *Service) GetAdmNeighbors(ctx context.Context, admId string) ([]Adm, error) {
	result, err := service.repo.GetAdmNeighbors(ctx, admId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
