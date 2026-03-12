package services

import (
	"sync"
	"time"

	"quickBillController/app"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/params/response"
	"quickBillController/models"

	"github.com/shopspring/decimal"
)

type SystemService struct {
}

func NewSystemService() *SystemService {
	return &SystemService{}
}

func (system *SystemService) SystemSummary(request request.SystemSummaryRequest) (resp response.SystemSummaryResponse, err error) {
	// 将时间戳转换为 time.Time
	startTime := time.Unix(request.StartAt, 0)
	endTime := time.Unix(request.EndAt, 0)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var queryErr error

	// 1. 总用户数
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		if err := app.DB.Model(&models.User{}).Where("created_at >= ? AND created_at <= ?", startTime, endTime).Count(&count).Error; err != nil {
			mu.Lock()
			if queryErr == nil {
				queryErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		resp.TotalUserCount = count
		mu.Unlock()
	}()

	// 2. 总卡数
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		if err := app.DB.Model(&models.Card{}).Where("created_at >= ? AND created_at <= ?", startTime, endTime).Count(&count).Error; err != nil {
			mu.Lock()
			if queryErr == nil {
				queryErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		resp.TotalCardCount = count
		mu.Unlock()
	}()

	// 3. 总消费金额（根据时间范围）
	wg.Add(1)
	go func() {
		defer wg.Done()
		var result struct {
			Total decimal.Decimal
		}
		if err := app.DB.Model(&models.Bill{}).
			Where("created_at >= ? AND created_at <= ?", startTime, endTime).
			Select("COALESCE(SUM(amount), 0) as total").
			Scan(&result).Error; err != nil {
			mu.Lock()
			if queryErr == nil {
				queryErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		resp.TotalBillAmount = result.Total
		mu.Unlock()
	}()

	// 4. 总用户充值金额（根据时间范围）
	wg.Add(1)
	go func() {
		defer wg.Done()
		var result struct {
			Total decimal.Decimal
		}
		if err := app.DB.Model(&models.Topup{}).
			Where("created_at >= ? AND created_at <= ?", startTime, endTime).
			Select("COALESCE(SUM(amount), 0) as total").
			Scan(&result).Error; err != nil {
			mu.Lock()
			if queryErr == nil {
				queryErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		resp.TotalUserTopupAmount = result.Total
		mu.Unlock()
	}()

	// 5. 总用户提现金额（根据时间范围，user_id > 0 且 merchant_id = 0）
	wg.Add(1)
	go func() {
		defer wg.Done()
		var result struct {
			Total decimal.Decimal
		}
		if err := app.DB.Model(&models.Withdrawal{}).
			Where("user_id > 0 AND merchant_id = 0 AND created_at >= ? AND created_at <= ?", startTime, endTime).
			Select("COALESCE(SUM(amount), 0) as total").
			Scan(&result).Error; err != nil {
			mu.Lock()
			if queryErr == nil {
				queryErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		resp.TotalUserWithdrawalAmount = result.Total
		mu.Unlock()
	}()

	// 6. 总商户提现金额
	wg.Add(1)
	go func() {
		defer wg.Done()
		var result struct {
			Total decimal.Decimal
		}
		if err := app.DB.Model(&models.Withdrawal{}).
			Where("merchant_id > 0 AND created_at >= ? AND created_at <= ?", startTime, endTime).
			Select("COALESCE(SUM(amount), 0) as total").
			Scan(&result).Error; err != nil {
			mu.Lock()
			if queryErr == nil {
				queryErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		resp.TotalMerchantWithdrawalAmount = result.Total
		mu.Unlock()
	}()

	// 7. 总商户冻结金额（查询所有商户的冻结金额总和）
	wg.Add(1)
	go func() {
		defer wg.Done()
		var merchants []models.Merchant
		if err := app.DB.Model(&models.Merchant{}).Find(&merchants).Error; err != nil {
			mu.Lock()
			if queryErr == nil {
				queryErr = err
			}
			mu.Unlock()
			return
		}

		totalFrozenAmount := decimal.NewFromInt(0)
		for _, merchant := range merchants {
			merchantBalanceData, err := merchant.GetMerchantBalance(merchant.Id)
			if err != nil {
				mu.Lock()
				if queryErr == nil {
					queryErr = err
				}
				mu.Unlock()
				return
			}
			totalFrozenAmount = totalFrozenAmount.Add(merchantBalanceData.FrozenAmount)
		}

		mu.Lock()
		resp.TotalMerchantFrozenAmount = totalFrozenAmount
		mu.Unlock()
	}()

	// 等待所有 goroutine 完成
	wg.Wait()

	if queryErr != nil {
		return resp, queryErr
	}

	return resp, nil
}
