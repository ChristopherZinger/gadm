package adm

import (
	"context"
	"gadm-api/utils"
)

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

func (service *Service) GetAdmForPoint(ctx context.Context, point utils.Point) (Adm, error) {
	result, err := service.repo.GetAdmForPoint(ctx, point)
	if err != nil {
		return Adm{}, err
	}
	return result, nil
}
