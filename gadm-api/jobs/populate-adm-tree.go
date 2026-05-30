package jobs

import (
	"context"
	"gadm-api/infra/pg"
	"gadm-api/logger"
	"gadm-api/models/adm"
)

const MAX_PG_CONNS = int32(45)

func PopulateAdmTreeJob() {
	dbPool := pg.InitPgPool(MAX_PG_CONNS)
	defer dbPool.Close()

	admRepo := adm.NewAdmRepo(dbPool)
	admService := adm.NewAdmService(admRepo)
	err := admService.PopulateAdmTree(context.Background())
	if err != nil {
		logger.Fatal("failed_to_populate_adm_tree %v", err)
	}
}

func PopulateAdmNeighborsJob() {
	dbPool := pg.InitPgPool(MAX_PG_CONNS)
	defer dbPool.Close()

	logger.Info("populate_adm_neighbors_job started")

	admRepo := adm.NewAdmRepo(dbPool)
	admService := adm.NewAdmService(admRepo)
	err := admService.PopulateAdmNeighbors(context.Background())
	if err != nil {
		logger.Fatal("failed_to_populate_adm_neighbors %v", err)
	}
}
