package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type LetterAligo struct {
	Id          int       `orm:"column(id);pk" description:"아이디"`
	LetterId    *Letter   `orm:"column(letter_id);rel(fk)" description:"문자"`
	Key         string    `orm:"column(key);size(48)" description:"인증용 API Key"`
	UserId      string    `orm:"column(user_id);size(48)" description:"사용자 아이디"`
	Sender      string    `orm:"column(sender);size(16)" description:"발신자 전화번호"`
	Receiver    string    `orm:"column(receiver);size(16)" description:"수신자 전화번호"`
	Msg         string    `orm:"column(msg);size(2000)" description:"메시지 내용"`
	MsgType     string    `orm:"column(msg_type);size(8);null" description:"SMS, LMS, MMS"`
	Title       string    `orm:"column(title);size(44);null" description:"문자제목"`
	Destination string    `orm:"column(destination);size(32);null" description:"치환용 입력"`
	Rdate       string    `orm:"column(rdate);size(8);null" description:"예약일"`
	Rtime       string    `orm:"column(rtime);size(4);null" description:"예약시간"`
	Image       string    `orm:"column(image);null" description:"첨부이미지"`
	TestmodeYn  string    `orm:"column(testmode_yn);size(1)" description:"테스트여부"`
	ResultCode  int       `orm:"column(result_code);null" description:"결과 코드"`
	Message     string    `orm:"column(message);size(256);null" description:"결과 메세지"`
	MsgId       int       `orm:"column(msg_id);null" description:"메세지 고유 ID"`
	SuccessCnt  int       `orm:"column(success_cnt);null" description:"요청성공 건수"`
	ErrorCnt    int       `orm:"column(error_cnt);null" description:"요청실패 건수"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy   string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy   string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *LetterAligo) TableName() string {
	return "letter_aligo"
}

func init() {
	orm.RegisterModel(new(LetterAligo))
}

// AddLetterAligo insert a new LetterAligo into database and returns
// last inserted Id on success.
func AddLetterAligo(m *LetterAligo) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetLetterAligoById retrieves LetterAligo by Id. Returns error if
// Id doesn't exist
func GetLetterAligoById(id int) (v *LetterAligo, err error) {
	o := orm.NewOrm()
	v = &LetterAligo{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLetterAligo retrieves all LetterAligo matches certain condition. Returns empty list if
// no records exist
func GetAllLetterAligo(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(LetterAligo))
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

	var l []LetterAligo
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

// UpdateLetterAligo updates LetterAligo by Id and returns error if
// the record to be updated doesn't exist
func UpdateLetterAligoById(m *LetterAligo) (err error) {
	o := orm.NewOrm()
	v := LetterAligo{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLetterAligo deletes LetterAligo by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLetterAligo(id int) (err error) {
	o := orm.NewOrm()
	v := LetterAligo{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&LetterAligo{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
