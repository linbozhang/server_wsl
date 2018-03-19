package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Artarget struct {
	Id          int       `orm:"column(id);auto"`
	VuforiaId   string    `orm:"column(vuforia_id);size(64);null"`
	WslName     string    `orm:"column(wsl_name);size(128);null"`
	IconUrl     string    `orm:"column(icon_url);size(512);null"`
	ShowId      int       `orm:"column(show_id);null"`
	ShowType    int8      `orm:"column(show_type);null"`
	ShowUrl     string    `orm:"column(show_url);size(512);null"`
	ProductName string    `orm:"column(product_name);size(128);null"`
	SortId      int       `orm:"column(sort_id);null"`
	BrandId     int       `orm:"column(brand_id);null"`
	Createtime  time.Time `orm:"column(createtime);type(datetime);null"`
	ActiveFlag  int8      `orm:"column(active_flag);null"`
	RedirectUrl string    `orm:"column(redirect_url);size(512);null"`
}

func (t *Artarget) TableName() string {
	return "artarget"
}

func init() {
	orm.RegisterModel(new(Artarget))
}

// AddArtarget insert a new Artarget into database and returns
// last inserted Id on success.
func AddArtarget(m *Artarget) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetArtargetById retrieves Artarget by Id. Returns error if
// Id doesn't exist
func GetArtargetById(id int) (v *Artarget, err error) {
	o := orm.NewOrm()
	v = &Artarget{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllArtarget retrieves all Artarget matches certain condition. Returns empty list if
// no records exist
func GetAllArtarget(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Artarget))
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

	var l []Artarget
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

// UpdateArtarget updates Artarget by Id and returns error if
// the record to be updated doesn't exist
func UpdateArtargetById(m *Artarget) (err error) {
	o := orm.NewOrm()
	v := Artarget{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteArtarget deletes Artarget by Id and returns error if
// the record to be deleted doesn't exist
func DeleteArtarget(id int) (err error) {
	o := orm.NewOrm()
	v := Artarget{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Artarget{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
