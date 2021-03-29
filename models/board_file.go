package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type BoardFile struct {
	Id               int    `orm:"column(id);pk" description:"아이디"`
	BoardId          *Board `orm:"column(board_id);rel(fk)" description:"게시판"`
	OriginalFilename string `orm:"column(original_filename);size(48)" description:"파일명"`
	Size             int64  `orm:"column(size)" description:"크기"`
	ContentType      string `orm:"column(content_type);size(32)" description:"유형"`
	BucketName       string `orm:"column(bucket_name);size(48)" description:"버킷"`
	ObjectName       string `orm:"column(object_name);size(256)" description:"오브젝트"`
}

func (t *BoardFile) TableName() string {
	return "board_file"
}

func init() {
	orm.RegisterModel(new(BoardFile))
}

// AddBoardFile insert a new BoardFile into database and returns
// last inserted Id on success.
func AddBoardFile(m *BoardFile) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetBoardFileById retrieves BoardFile by Id. Returns error if
// Id doesn't exist
func GetBoardFileById(id int) (v *BoardFile, err error) {
	o := orm.NewOrm()
	v = &BoardFile{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBoardFile retrieves all BoardFile matches certain condition. Returns empty list if
// no records exist
func GetAllBoardFile(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(BoardFile))
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

	var l []BoardFile
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

// UpdateBoardFile updates BoardFile by Id and returns error if
// the record to be updated doesn't exist
func UpdateBoardFileById(m *BoardFile) (err error) {
	o := orm.NewOrm()
	v := BoardFile{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBoardFile deletes BoardFile by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBoardFile(id int) (err error) {
	o := orm.NewOrm()
	v := BoardFile{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&BoardFile{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
