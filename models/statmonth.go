package models

import "time"

type Statmonth struct {
	ShopId       *Shop     `orm:"column(shop_id);rel(fk)" description:"가맹점"`
	SaleMonth    string    `orm:"column(sale_month);size(7)" description:"영업월"`
	ShopNm       string    `orm:"column(shop_nm);size(48)" description:"가맹점"`
	TotalSaleCnt int16     `orm:"column(total_sale_cnt)" description:"판매건수"`
	TotalSaleAmt int       `orm:"column(total_sale_amt)" description:"판매금액"`
	SaleSvcCnt   int16     `orm:"column(sale_svc_cnt)" description:"시술 판매건수"`
	SaleSvcAmt   int       `orm:"column(sale_svc_amt)" description:"시술 판매액"`
	SaleMebsCnt  int16     `orm:"column(sale_mebs_cnt)" description:"회원권 판매건수"`
	SaleMebsAmt  int       `orm:"column(sale_mebs_amt)" description:"회원권 판매액"`
	SaleDepoCnt  int16     `orm:"column(sale_depo_cnt)" description:"선불권 판매건수"`
	SaleDepoAmt  int       `orm:"column(sale_depo_amt)" description:"선불권 판매액"`
	SaleGoodsCnt int16     `orm:"column(sale_goods_cnt)" description:"상품 판매건수"`
	SaleGoodsAmt int       `orm:"column(sale_goods_amt)" description:"상품 판매액"`
	PayCashCnt   int16     `orm:"column(pay_cash_cnt)" description:"현금 결제건수"`
	PayCashAmt   int       `orm:"column(pay_cash_amt)" description:"현금 결제액"`
	PayCardCnt   int16     `orm:"column(pay_card_cnt)" description:"카드 결제건수"`
	PayCardAmt   int       `orm:"column(pay_card_amt)" description:"카드 결제액"`
	PayMembCnt   int16     `orm:"column(pay_memb_cnt)" description:"회원권 결제건수"`
	PayMembAmt   int       `orm:"column(pay_memb_amt)" description:"회원권 결제액"`
	PayDepoCnt   int16     `orm:"column(pay_depo_cnt)" description:"선불권 결제건수"`
	PayDepoAmt   int       `orm:"column(pay_depo_amt)" description:"선불권 결제액"`
	PayAcctCnt   int16     `orm:"column(pay_acct_cnt)" description:"계좌이체 결제건수"`
	PayAcctAmt   int       `orm:"column(pay_acct_amt)" description:"계좌이체 결제액"`
	PayPointCnt  int16     `orm:"column(pay_point_cnt)" description:"포인트 결제건수"`
	PayPointAmt  int       `orm:"column(pay_point_amt)" description:"포인트 결제액"`
	AccPointCnt  int16     `orm:"column(acc_point_cnt)" description:"포인트 적립건수"`
	AccPointAmt  int       `orm:"column(acc_point_amt)" description:"포인트 적립액"`
	CreatedAt    time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	UpdatedAt    time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
}
