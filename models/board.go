package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Board struct {
	Id            int       `orm:"column(id);pk" description:"아이디"`
	Subject       string    `orm:"column(subject);size(48)" description:"제목"`
	BoardTypeCd   string    `orm:"column(board_type_cd);size(8)" description:"게시판 유형"`
	BoardTypeCdNm string    `orm:"column(board_type_cd_nm);size(48);null" description:"게시판 유형"`
	ReadCnt       int       `orm:"column(read_cnt)" description:"조회수"`
	CommentCnt    int       `orm:"column(comment_cnt)" description:"댓글수"`
	FileCnt       int16     `orm:"column(file_cnt)" description:"첨부파일수"`
	StartDe       string    `orm:"column(start_de);size(10);null" description:"게시시작일"`
	EndDe         string    `orm:"column(end_de);size(10);null" description:"게시종료일"`
	Content       string    `orm:"column(content)" description:"내용"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy     string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy     string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *Board) TableName() string {
	return "board"
}

func init() {
	orm.RegisterModel(new(Board))
}

// AddBoard insert a new Board into database and returns
// last inserted Id on success.
func AddBoard(m *Board) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetBoardById retrieves Board by Id. Returns error if
// Id doesn't exist
func GetBoardById(id int) (v *Board, err error) {
	o := orm.NewOrm()
	v = &Board{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBoard retrieves all Board matches certain condition. Returns empty list if
// no records exist
func GetAllBoard(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Board))
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

	var l []Board
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

// UpdateBoard updates Board by Id and returns error if
// the record to be updated doesn't exist
func UpdateBoardById(m *Board) (err error) {
	o := orm.NewOrm()
	v := Board{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBoard deletes Board by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBoard(id int) (err error) {
	o := orm.NewOrm()
	v := Board{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Board{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
