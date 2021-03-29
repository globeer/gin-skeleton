package models

import "time"

type CustomerDebt struct {
	CustomerId    *Customer `orm:"column(customer_id);rel(fk)" description:"고객"`
	GoodsTypeCd   string    `orm:"column(goods_type_cd);size(8)" description:"상품 유형"`
	GoodsTypeCdNm string    `orm:"column(goods_type_cd_nm);size(48);null" description:"상품 유형"`
	AccAmt        int       `orm:"column(acc_amt)" description:"적립금액"`
	UseAmt        int       `orm:"column(use_amt)" description:"사용금액"`
	RestAmt       int       `orm:"column(rest_amt)" description:"잔여"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy     string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy     string    `orm:"column(updated_by);size(48)" description:"수정자"`
}
