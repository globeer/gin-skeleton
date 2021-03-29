package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type SetLetter struct {
	Id             int       `orm:"column(id);pk" description:"아이디"`
	ShopId         *Shop     `orm:"column(shop_id);rel(fk)" description:"가맹점"`
	ShopNm         string    `orm:"column(shop_nm);size(48)" description:"가맹점"`
	LetterTypeCd   string    `orm:"column(letter_type_cd);size(8)" description:"문자 유형"`
	LetterTypeCdNm string    `orm:"column(letter_type_cd_nm);size(48)" description:"문자 유형"`
	BeforeDay      int16     `orm:"column(before_day)" description:"몇일전"`
	BeforeHour     int16     `orm:"column(before_hour)" description:"몇시간전"`
	UseYn          string    `orm:"column(use_yn);size(1)" description:"사용"`
	CreatedAt      time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy      string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy      string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *SetLetter) TableName() string {
	return "set_letter"
}

func init() {
	orm.RegisterModel(new(SetLetter))
}

// AddSetLetter insert a new SetLetter into database and returns
// last inserted Id on success.
func AddSetLetter(m *SetLetter) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetSetLetterById retrieves SetLetter by Id. Returns error if
// Id doesn't exist
func GetSetLetterById(id int) (v *SetLetter, err error) {
	o := orm.NewOrm()
	v = &SetLetter{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSetLetter retrieves all SetLetter matches certain condition. Returns empty list if
// no records exist
func GetAllSetLetter(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SetLetter))
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

	var l []SetLetter
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

// UpdateSetLetter updates SetLetter by Id and returns error if
// the record to be updated doesn't exist
func UpdateSetLetterById(m *SetLetter) (err error) {
	o := orm.NewOrm()
	v := SetLetter{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSetLetter deletes SetLetter by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSetLetter(id int) (err error) {
	o := orm.NewOrm()
	v := SetLetter{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SetLetter{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
