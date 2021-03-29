package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Letter struct {
	Id          int       `orm:"column(id);pk" description:"아이디"`
	ShopId      *Shop     `orm:"column(shop_id);rel(fk)" description:"가맹점"`
	ShopNm      string    `orm:"column(shop_nm);size(48)" description:"가맹점 명"`
	CustId      int64     `orm:"column(cust_id)" description:"고객"`
	CustTelNo   string    `orm:"column(cust_tel_no);size(16)" description:"핸드폰 번호"`
	LetterTelNo string    `orm:"column(letter_tel_no);size(16)" description:"발송 전화 번호"`
	Subject     string    `orm:"column(subject);size(48);null" description:"제목"`
	Msg         string    `orm:"column(msg);size(2000)" description:"문자 내용"`
	MsgType     string    `orm:"column(msg_type);size(8)" description:"문자 유형"`
	ResultCode  string    `orm:"column(result_code);size(8);null" description:"결과 코드"`
	ResultMsg   string    `orm:"column(result_msg);size(512);null" description:"결과 메시지"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy   string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy   string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *Letter) TableName() string {
	return "letter"
}

func init() {
	orm.RegisterModel(new(Letter))
}

// AddLetter insert a new Letter into database and returns
// last inserted Id on success.
func AddLetter(m *Letter) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetLetterById retrieves Letter by Id. Returns error if
// Id doesn't exist
func GetLetterById(id int) (v *Letter, err error) {
	o := orm.NewOrm()
	v = &Letter{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLetter retrieves all Letter matches certain condition. Returns empty list if
// no records exist
func GetAllLetter(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Letter))
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

	var l []Letter
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

// UpdateLetter updates Letter by Id and returns error if
// the record to be updated doesn't exist
func UpdateLetterById(m *Letter) (err error) {
	o := orm.NewOrm()
	v := Letter{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLetter deletes Letter by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLetter(id int) (err error) {
	o := orm.NewOrm()
	v := Letter{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Letter{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
