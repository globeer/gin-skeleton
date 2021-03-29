package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Staff struct {
	Id          int       `orm:"column(id);pk" description:"직원"`
	StaffNm     string    `orm:"column(staff_nm);size(45)" description:"직원 명"`
	ShopId      *Shop     `orm:"column(shop_id);rel(fk)" description:"가맹점"`
	ShopNm      string    `orm:"column(shop_nm);size(48)" description:"가맹점"`
	TelNo       string    `orm:"column(tel_no);size(16);null" description:"전화 번호"`
	PaytypeCd   string    `orm:"column(paytype_cd);size(8)" description:"임금 유형"`
	PaytypeCdNm string    `orm:"column(paytype_cd_nm);size(48)" description:"임금 유형"`
	PayAmt      int       `orm:"column(pay_amt)" description:"임금"`
	StaffCd     string    `orm:"column(staff_cd);size(8);null" description:"직원 유형"`
	StaffCdNm   string    `orm:"column(staff_cd_nm);size(48);null" description:"직원 유형명"`
	BirthDe     string    `orm:"column(birth_de);size(10);null" description:"생일"`
	MarriageDe  string    `orm:"column(marriage_de);size(10);null" description:"결혼기념일"`
	EnterDe     string    `orm:"column(enter_de);size(10);null" description:"입사일"`
	RetireDe    string    `orm:"column(retire_de);size(10);null" description:"퇴사일"`
	UseYn       string    `orm:"column(use_yn);size(1)" description:"사용"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy   string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy   string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *Staff) TableName() string {
	return "staff"
}

func init() {
	orm.RegisterModel(new(Staff))
}

// AddStaff insert a new Staff into database and returns
// last inserted Id on success.
func AddStaff(m *Staff) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetStaffById retrieves Staff by Id. Returns error if
// Id doesn't exist
func GetStaffById(id int) (v *Staff, err error) {
	o := orm.NewOrm()
	v = &Staff{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllStaff retrieves all Staff matches certain condition. Returns empty list if
// no records exist
func GetAllStaff(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Staff))
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

	var l []Staff
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

// UpdateStaff updates Staff by Id and returns error if
// the record to be updated doesn't exist
func UpdateStaffById(m *Staff) (err error) {
	o := orm.NewOrm()
	v := Staff{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteStaff deletes Staff by Id and returns error if
// the record to be deleted doesn't exist
func DeleteStaff(id int) (err error) {
	o := orm.NewOrm()
	v := Staff{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Staff{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
