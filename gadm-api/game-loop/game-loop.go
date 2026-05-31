package gameloop

import (
	"context"
	"gadm-api/logger"
	"gadm-api/utils"
	"time"
)

var TOUR_DURATION = 5 * time.Second

func GameLoop() {
	ctx := context.Background()
	tourTimer := getTourTimer(ctx, TOUR_DURATION)
	loopCount := 0
	for {

		logger.Info("game_loop_iteration", "loopCount", loopCount)
		tourTimer()
		loopCount++
	}
}

func getTourTimer(ctx context.Context, loopDuration time.Duration) func() {
	lastTourTime := time.Now()
	tourTimer := func() {
		for {
			if time.Now().Sub(lastTourTime) > loopDuration {
				lastTourTime = time.Now()
				break
			}
			utils.Sleep(ctx, 1*time.Second)
		}

	}
	return tourTimer
}
