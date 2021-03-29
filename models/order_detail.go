package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type OrderDetail struct {
	Id            int       `orm:"column(id);pk" description:"아이디"`
	OrderId       *Order    `orm:"column(order_id);rel(fk)" description:"주문"`
	GoodsTypeCd   string    `orm:"column(goods_type_cd);size(8)" description:"상품 유형"`
	GoodsTypeCdNm string    `orm:"column(goods_type_cd_nm);size(48);null" description:"상품 유형"`
	CategoryId    int64     `orm:"column(category_id)" description:"카테고리"`
	CategoryNm    string    `orm:"column(category_nm);size(48);null" description:"카테고리"`
	GoodsId       int64     `orm:"column(goods_id)" description:"상품"`
	GoodsNm       string    `orm:"column(goods_nm);size(48);null" description:"상품"`
	StaffId       *Staff    `orm:"column(staff_id);rel(fk)" description:"담당자"`
	StaffNm       string    `orm:"column(staff_nm);size(45);null" description:"담당자"`
	PayTypeCd     string    `orm:"column(pay_type_cd);size(8)" description:"결제 유형"`
	PayTypeCdNm   string    `orm:"column(pay_type_cd_nm);size(48)" description:"결제 유형"`
	Price         int       `orm:"column(price)" description:"상품금액"`
	DiscountAmt   int       `orm:"column(discount_amt);null" description:"할인액"`
	SalePrice     int       `orm:"column(sale_price)" description:"판매액"`
	SaleCnt       int16     `orm:"column(sale_cnt)" description:"판매수량"`
	TotalPrice    int       `orm:"column(total_price)" description:"합계액"`
	PointAmt      int       `orm:"column(point_amt);null" description:"적립포인트"`
	DepositAmt    int       `orm:"column(deposit_amt);null" description:"선불액"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy     string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy     string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *OrderDetail) TableName() string {
	return "order_detail"
}

func init() {
	orm.RegisterModel(new(OrderDetail))
}

// AddOrderDetail insert a new OrderDetail into database and returns
// last inserted Id on success.
func AddOrderDetail(m *OrderDetail) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOrderDetailById retrieves OrderDetail by Id. Returns error if
// Id doesn't exist
func GetOrderDetailById(id int) (v *OrderDetail, err error) {
	o := orm.NewOrm()
	v = &OrderDetail{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllOrderDetail retrieves all OrderDetail matches certain condition. Returns empty list if
// no records exist
func GetAllOrderDetail(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(OrderDetail))
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

	var l []OrderDetail
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

// UpdateOrderDetail updates OrderDetail by Id and returns error if
// the record to be updated doesn't exist
func UpdateOrderDetailById(m *OrderDetail) (err error) {
	o := orm.NewOrm()
	v := OrderDetail{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteOrderDetail deletes OrderDetail by Id and returns error if
// the record to be deleted doesn't exist
func DeleteOrderDetail(id int) (err error) {
	o := orm.NewOrm()
	v := OrderDetail{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&OrderDetail{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
