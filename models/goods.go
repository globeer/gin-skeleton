package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Goods struct {
	Id            int       `orm:"column(id);pk" description:"아이디"`
	GoodsTypeCd   string    `orm:"column(goods_type_cd);size(8)" description:"상품 유형"`
	GoodsTypeCdNm string    `orm:"column(goods_type_cd_nm);size(48);null" description:"상품 유형"`
	ShopId        *Shop     `orm:"column(shop_id);rel(fk)" description:"가맹점"`
	ShopNm        string    `orm:"column(shop_nm);size(48)" description:"가맹점"`
	CategoryId    *Category `orm:"column(category_id);rel(fk)" description:"카테고리"`
	CategoryNm    string    `orm:"column(category_nm);size(48)" description:"카테고리"`
	GoodsNm       string    `orm:"column(goods_nm);size(48)" description:"상품 명"`
	SortSeq       int16     `orm:"column(sort_seq);null" description:"순번"`
	Price         int       `orm:"column(price)" description:"상품금액"`
	DiscountRate  int16     `orm:"column(discount_rate)" description:"할인율"`
	DiscountAmt   int       `orm:"column(discount_amt)" description:"할인액"`
	SalePrice     int       `orm:"column(sale_price)" description:"판매액"`
	MembYn        string    `orm:"column(memb_yn);size(1);null" description:"회원권"`
	MembRate      int16     `orm:"column(memb_rate)" description:"회원권할인율"`
	PointRate     int16     `orm:"column(point_rate);null" description:"적립율"`
	PointAmt      int       `orm:"column(point_amt);null" description:"적립포인트"`
	DepositRate   int16     `orm:"column(deposit_rate);null" description:"적립율"`
	DepositAmt    int       `orm:"column(deposit_amt);null" description:"선불액"`
	UseYn         string    `orm:"column(use_yn);size(1)" description:"사용"`
	Memo          string    `orm:"column(memo);size(1024);null" description:"메모"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy     string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy     string    `orm:"column(updated_by);size(48)" description:"수정자"`
}

func (t *Goods) TableName() string {
	return "goods"
}

func init() {
	orm.RegisterModel(new(Goods))
}

// AddGoods insert a new Goods into database and returns
// last inserted Id on success.
func AddGoods(m *Goods) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetGoodsById retrieves Goods by Id. Returns error if
// Id doesn't exist
func GetGoodsById(id int) (v *Goods, err error) {
	o := orm.NewOrm()
	v = &Goods{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllGoods retrieves all Goods matches certain condition. Returns empty list if
// no records exist
func GetAllGoods(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Goods))
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

	var l []Goods
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

// UpdateGoods updates Goods by Id and returns error if
// the record to be updated doesn't exist
func UpdateGoodsById(m *Goods) (err error) {
	o := orm.NewOrm()
	v := Goods{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteGoods deletes Goods by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGoods(id int) (err error) {
	o := orm.NewOrm()
	v := Goods{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Goods{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
