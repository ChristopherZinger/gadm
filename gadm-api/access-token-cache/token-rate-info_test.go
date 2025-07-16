package accessTokenCache

import (
	"testing"
	"time"
)

type DateInfo struct {
	month int
	day   int
	year  int
}

func getCreatedAtFromDateInfo(dateInfo DateInfo) time.Time {
	return time.Now().AddDate(dateInfo.year, dateInfo.month, dateInfo.day)
}
func TestHandleHitExpiredToken(t *testing.T) {
	t.Logf("Test: isTokenExpired")

	expiredDateInfos := []DateInfo{
		{year: 0, month: -3, day: -1},
		{year: 0, month: -4, day: 0},
		{year: -1, month: 0, day: 0},
	}

	for _, expiredDate := range expiredDateInfos {
		if !isTokenExpired(getCreatedAtFromDateInfo(expiredDate)) {
			t.Errorf("failed to detect expired token: date: Y%d M%d D%d",
				expiredDate.year, expiredDate.month, expiredDate.day)
		}
	}

	datesBeforeExpiration := []DateInfo{
		{year: 0, month: -3, day: 1},
		{year: 0, month: 0, day: 0},
		{year: 0, month: -1, day: -3},
		{year: 0, month: -2, day: -29},
	}

	for _, v := range datesBeforeExpiration {
		if isTokenExpired(getCreatedAtFromDateInfo(v)) {
			t.Errorf("valid token date evaluated as expired: date: Y%d M%d D%d",
				v.year, v.month, v.day)
		}
	}
}
