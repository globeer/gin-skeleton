package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Category struct {
	Id            int       `orm:"column(id);pk" description:"카테고리"`
	GoodsTypeCd   string    `orm:"column(goods_type_cd);size(8);null" description:"상품 유형"`
	GoodsTypeCdNm string    `orm:"column(goods_type_cd_nm);size(48);null" description:"상품 유형"`
	ShopId        *Shop     `orm:"column(shop_id);rel(fk)" description:"가맹점"`
	ShopNm        string    `orm:"column(shop_nm);size(48)" description:"가맹점"`
	CategoryNm    string    `orm:"column(category_nm);size(48)" description:"카테고리 명"`
	SortSeq       int16     `orm:"column(sort_seq)" description:"순번"`
	UseYn         string    `orm:"column(use_yn);size(1)" description:"사용"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy     string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy     string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *Category) TableName() string {
	return "category"
}

func init() {
	orm.RegisterModel(new(Category))
}

// AddCategory insert a new Category into database and returns
// last inserted Id on success.
func AddCategory(m *Category) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCategoryById retrieves Category by Id. Returns error if
// Id doesn't exist
func GetCategoryById(id int) (v *Category, err error) {
	o := orm.NewOrm()
	v = &Category{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCategory retrieves all Category matches certain condition. Returns empty list if
// no records exist
func GetAllCategory(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Category))
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

	var l []Category
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

// UpdateCategory updates Category by Id and returns error if
// the record to be updated doesn't exist
func UpdateCategoryById(m *Category) (err error) {
	o := orm.NewOrm()
	v := Category{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCategory deletes Category by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCategory(id int) (err error) {
	o := orm.NewOrm()
	v := Category{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Category{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
