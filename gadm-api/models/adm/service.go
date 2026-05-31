package adm

import (
	"context"
	"fmt"
	"gadm-api/logger"
	"gadm-api/utils"
	"time"

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
	batchSize := 20
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
		g.SetLimit(20)

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

		utils.Sleep(gctx, 2*time.Minute)

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
	batchSize := 100
	processedCount := 0

	processBatch := func(ctx context.Context, batch []Adm) error {
		g, gctx := errgroup.WithContext(ctx)
		g.SetLimit(2)

		for _, adm := range batch {
			g.Go(func() error {
				neighbors, err := utils.Retry(
					gctx,
					func(ctx context.Context) ([]Adm, error) {
						return service.repo.GetNeighbors(ctx, adm.ID)
					},
					3,
					1*time.Second,
					10*time.Second)
				if err != nil {
					if gctx.Err() != nil {
						return gctx.Err()
					}
					logger.Error("skip_adm_neighbors: adm_id=%s err=%v", adm.ID, err)
					return nil
				}

				logger.Info("neighbors for adm_id=%s: %d", adm.ID, len(neighbors))

				neighborIds := make([]string, 0, len(neighbors))
				for _, neighbor := range neighbors {
					if neighbor.ID == adm.ID {
						continue
					}
					neighborIds = append(neighborIds, neighbor.ID)
				}

				if err := service.repo.UpsertAdmNeighborsBatch(gctx, adm.ID, neighborIds); err != nil {
					return fmt.Errorf(
						"failed_to_upsert_neighbors_batch: adm_id=%s: %w",
						adm.ID, err)
				}

				processedCount++
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		lastId := batch[len(batch)-1].ID
		logger.Info("populate_adm_neighbors_progress processed=%d last_id=%s", processedCount, lastId)
		return nil
	}

	err := service.repo.IterateLeafAdms(ctx, batchSize, processBatch)
	if err != nil {
		return err
	}

	logger.Info("populate_adm_neighbors_done processed=%d", processedCount)
	return nil
}
