package response

import "github.com/shopspring/decimal"

type SystemSummaryResponse SystemSummary

// 1 总用户数
// 2 总卡数
// 3 总消费金额
// 4 总充值金额
// 5 总用户提现金额
// 6 总商户提现金额
// 7 总商户冻结金额

type SystemSummary struct {
	TotalUserCount                int64           `json:"total_user_count"`                 //1 总用户数
	TotalCardCount                int64           `json:"total_card_count"`                 //2 总卡数
	TotalBillAmount               decimal.Decimal `json:"total_bill_amount"`                //3 总消费金额
	TotalUserTopupAmount          decimal.Decimal `json:"total_user_topup_amount"`          //4 总充值金额
	TotalUserWithdrawalAmount     decimal.Decimal `json:"total_user_withdrawal_amount"`     //5 总用户提现金额
	TotalMerchantWithdrawalAmount decimal.Decimal `json:"total_merchant_withdrawal_amount"` //5 总商户提现金额
	TotalMerchantFrozenAmount     decimal.Decimal `json:"total_merchant_frozen_amount"`     //6 总商户冻结金额
}
