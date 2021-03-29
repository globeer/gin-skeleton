package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Customer struct {
	Id            int       `orm:"column(id);pk" description:"고객"`
	CustomerNm    string    `orm:"column(customer_nm);size(128)" description:"고객"`
	CustTelNo     string    `orm:"column(cust_tel_no);size(16);null" description:"핸드폰 번호"`
	CustTelNo4    string    `orm:"column(cust_tel_no4);size(4);null" description:"번호4자리"`
	ShopId        *Shop     `orm:"column(shop_id);rel(fk)" description:"가맹점"`
	ShopNm        string    `orm:"column(shop_nm);size(48)" description:"가맹점"`
	StaffId       *Staff    `orm:"column(staff_id);rel(fk)" description:"담당자"`
	StaffNm       string    `orm:"column(staff_nm);size(48);null" description:"담당자"`
	DepositAmt    int       `orm:"column(deposit_amt)" description:"예치금"`
	PointAmt      int       `orm:"column(point_amt)" description:"포인트"`
	BirthDe       string    `orm:"column(birth_de);size(10);null" description:"생일"`
	MarriageDe    string    `orm:"column(marriage_de);size(10);null" description:"결혼기념일"`
	LatestVisitDe string    `orm:"column(latest_visit_de);size(10);null" description:"최근 방문 일"`
	LetterYn      string    `orm:"column(letter_yn);size(1)" description:"문자 수신"`
	Memo          string    `orm:"column(memo);size(1024);null" description:"메모"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy     string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy     string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *Customer) TableName() string {
	return "customer"
}

func init() {
	orm.RegisterModel(new(Customer))
}

// AddCustomer insert a new Customer into database and returns
// last inserted Id on success.
func AddCustomer(m *Customer) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCustomerById retrieves Customer by Id. Returns error if
// Id doesn't exist
func GetCustomerById(id int) (v *Customer, err error) {
	o := orm.NewOrm()
	v = &Customer{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCustomer retrieves all Customer matches certain condition. Returns empty list if
// no records exist
func GetAllCustomer(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Customer))
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

	var l []Customer
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

// UpdateCustomer updates Customer by Id and returns error if
// the record to be updated doesn't exist
func UpdateCustomerById(m *Customer) (err error) {
	o := orm.NewOrm()
	v := Customer{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCustomer deletes Customer by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCustomer(id int) (err error) {
	o := orm.NewOrm()
	v := Customer{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Customer{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
