package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Reserve struct {
	Id         int       `orm:"column(id);pk" description:"아이디"`
	ShopId     *Shop     `orm:"column(shop_id);rel(fk)" description:"가맹점"`
	ShopNm     string    `orm:"column(shop_nm);size(48);null" description:"가맹점"`
	CustomerId int64     `orm:"column(customer_id)" description:"고객"`
	ReserveDe  string    `orm:"column(reserve_de);size(10)" description:"예약일"`
	ReserveHm  string    `orm:"column(reserve_hm);size(8)" description:"예약시간"`
	CategoryId int64     `orm:"column(category_id)" description:"카테고리"`
	CategoryNm string    `orm:"column(category_nm);size(48);null" description:"카테고리 명"`
	GoodsId    int64     `orm:"column(goods_id)" description:"상품"`
	GoodsNm    string    `orm:"column(goods_nm);size(48);null" description:"상품"`
	CompleteYn string    `orm:"column(complete_yn);size(1);null" description:"완료"`
	CreatedAt  time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy  string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt  time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy  string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *Reserve) TableName() string {
	return "reserve"
}

func init() {
	orm.RegisterModel(new(Reserve))
}

// AddReserve insert a new Reserve into database and returns
// last inserted Id on success.
func AddReserve(m *Reserve) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetReserveById retrieves Reserve by Id. Returns error if
// Id doesn't exist
func GetReserveById(id int) (v *Reserve, err error) {
	o := orm.NewOrm()
	v = &Reserve{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllReserve retrieves all Reserve matches certain condition. Returns empty list if
// no records exist
func GetAllReserve(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Reserve))
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

	var l []Reserve
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

// UpdateReserve updates Reserve by Id and returns error if
// the record to be updated doesn't exist
func UpdateReserveById(m *Reserve) (err error) {
	o := orm.NewOrm()
	v := Reserve{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteReserve deletes Reserve by Id and returns error if
// the record to be deleted doesn't exist
func DeleteReserve(id int) (err error) {
	o := orm.NewOrm()
	v := Reserve{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Reserve{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
