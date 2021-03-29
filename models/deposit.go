package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Deposit struct {
	Id          int       `orm:"column(id);pk" description:"아이디"`
	ShopId      *Shop     `orm:"column(shop_id);rel(fk)" description:"가맹점"`
	ShopNm      string    `orm:"column(shop_nm);size(48)" description:"가맹점"`
	DepositAmt  int       `orm:"column(deposit_amt)" description:"예치금"`
	UseYn       string    `orm:"column(use_yn);size(1)" description:"사용"`
	GoodsNm     string    `orm:"column(goods_nm);size(128);null" description:"상품"`
	PaymentCd   string    `orm:"column(payment_cd);size(8)" description:"결제 유형"`
	PaymentCdNm string    `orm:"column(payment_cd_nm);size(48)" description:"결제 유형"`
	ResultCd    string    `orm:"column(result_cd);size(8);null" description:"결과 코드"`
	ResultMsg   string    `orm:"column(result_msg);size(128);null" description:"결과 메시지"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy   string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy   string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *Deposit) TableName() string {
	return "deposit"
}

func init() {
	orm.RegisterModel(new(Deposit))
}

// AddDeposit insert a new Deposit into database and returns
// last inserted Id on success.
func AddDeposit(m *Deposit) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetDepositById retrieves Deposit by Id. Returns error if
// Id doesn't exist
func GetDepositById(id int) (v *Deposit, err error) {
	o := orm.NewOrm()
	v = &Deposit{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllDeposit retrieves all Deposit matches certain condition. Returns empty list if
// no records exist
func GetAllDeposit(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Deposit))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Deposit
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateDeposit updates Deposit by Id and returns error if
// the record to be updated doesn't exist
func UpdateDepositById(m *Deposit) (err error) {
	o := orm.NewOrm()
	v := Deposit{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDeposit deletes Deposit by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDeposit(id int) (err error) {
	o := orm.NewOrm()
	v := Deposit{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Deposit{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
