package adm

import (
	"context"
	"fmt"
	"gadm-api/logger"
	"gadm-api/utils"

	"golang.org/x/sync/errgroup"
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

func (service *Service) GetAdmNeighborsForPoint(ctx context.Context, point utils.Point) ([]Adm, error) {
	result, err := service.repo.GetAdmForPoint(ctx, point)
	if err != nil {
		return nil, err
	}
	neighbors, err := service.repo.GetAdmNeighbors(ctx, result.ID)
	if err != nil {
		return nil, err
	}
	return neighbors, nil
}

func (service *Service) PopulateAdmTree(ctx context.Context) error {
	startAfterId := ""
	batchSize := 1000
	processedCount := 0

	for {
		adms, err := service.repo.GetAdms(ctx, startAfterId, batchSize)
		if err != nil {
			return fmt.Errorf("fetch_adms_batch_after: start_after_id=%q: %w", startAfterId, err)
		}
		if len(adms) == 0 {
			break
		}

		g, gctx := errgroup.WithContext(ctx)
		g.SetLimit(45) // no more than total postgres connections in pool

		for _, adm := range adms {
			g.Go(func() error {
				children, err := service.repo.
					GetAdmsDirectChildrenForId(gctx, adm.ID, adm.Level)
				if err != nil {
					return fmt.Errorf("failed_to_get_direct_children: for_adm_id=%s: %w", adm.ID, err)
				}
				if len(children) == 0 {
					return nil
				}

				childIds := make([]string, len(children))
				for i, child := range children {
					childIds[i] = child.ID
				}

				if err := service.repo.UpsertAdmTreeRelationships(gctx, adm.ID, childIds); err != nil {
					return fmt.Errorf("failed_to_upsert_tree: adm_id=%s: %w", adm.ID, err)
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		processedCount += len(adms)
		lastId := adms[len(adms)-1].ID
		logger.Info("populate_adm_tree_progress processed=%d last_id=%s", processedCount, lastId)

		if len(adms) < batchSize {
			break
		}
		startAfterId = lastId
	}

	logger.Info("populate_adm_tree_done processed=%d", processedCount)
	return nil
}

func (service *Service) PopulateAdmNeighbors(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
