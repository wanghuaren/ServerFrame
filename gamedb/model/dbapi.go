package model

import (
	"baseutils/baseuts"
	"gameutils/common/dbuts"
	"gameutils/dbstruct"
	"strings"
	"time"

	"reflect"

	"github.com/beego/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func initDB() {
	err := orm.RegisterDataBase("default", "mysql", baseuts.MysqlAccount+":"+baseuts.MysqlPwd+"@tcp(127.0.0.1:"+baseuts.MysqlPort+")/"+baseuts.MysqlDBName+"?charset=utf8&parseTime=true&loc=Local")
	if !ChkErr(err, "连接数据库失败:") {
		for _, v := range dbstruct.DBTabName {
			rs := v.MStruct()
			orm.RegisterModel(rs)
		}
		Log("连接数据库成功:表注册完成")
		baseuts.AutoBackupMysql()
	} else {
		Log("重试")
		time.Sleep(time.Second)
		initDB()
	}
	orm.Debug = IsDebug
}

func DBAdd(install_data any, fields ...[]string) interface{} {
	if fields == nil || len(fields) < 2 || fields[1] == nil || len(fields[1]) < 1 {
		Log("DBAdd", "参数不足")
		return nil
	}

	tabName, _ := dbuts.GetTabNameIDFromDBDat(install_data)
	o := orm.NewOrm()
	editFields := fields[1]
	sql := "insert into " + tabName + "(" + strings.Join(editFields, ",") + ") values ("
	values := []interface{}{}
	for _, _v := range editFields {
		sql += "?,"
		_type := reflect.TypeOf(install_data).Elem()
		_value := reflect.ValueOf(install_data).Elem()

		for i := 0; i < _type.NumField(); i++ {
			key := _type.Field(i).Name
			if key == _v {
				value := _value.Field(i)
				switch value.Kind() {
				case reflect.Int:
					values = append(values, value.Int())
				case reflect.Int8:
					values = append(values, int8(value.Int()))
				case reflect.Int16:
					values = append(values, int16(value.Int()))
				case reflect.Int32:
					values = append(values, int32(value.Int()))
				case reflect.Int64:
					values = append(values, int64(value.Int()))
				case reflect.String:
					values = append(values, value.String())
				case reflect.Bool:
					values = append(values, value.Bool())
				case reflect.Float32, reflect.Float64:
					values = append(values, value.Float())
				}
				break
			}
		}
	}
	sql = sql[:len(sql)-1] + ")"

	_result, err := o.Raw(sql, values).Exec()
	if ChkErr(err) {
		return nil
	}
	_, err = _result.LastInsertId()
	if ChkErr(err) {
		return nil
	}
	_resultList := DBFind(install_data, editFields)
	return _resultList
}

// func DBDel(del_data any, fields ...[]string) bool {
func DBDel(del_data any, fields ...[]string) interface{} {
	if fields == nil || len(fields) < 1 || fields[0] == nil || len(fields[0]) < 1 {
		Log("DBDel", "参数不足")
		return nil
	}
	o := orm.NewOrm()
	_resultList := DBFind(del_data, fields[0])
	_, err := o.Delete(del_data, fields[0]...)
	if ChkErr(err) {
		return nil
	}
	return _resultList
}
func DBEdit(edit_data any, fields ...[]string) interface{} {
	if fields == nil || len(fields) < 2 || fields[0] == nil || len(fields[0]) < 1 || fields[1] == nil || len(fields[1]) < 1 {
		Log("DBEdit", "参数不足")
		return nil
	}
	tabName, _ := dbuts.GetTabNameIDFromDBDat(edit_data)

	searchFields := fields[0]
	editFields := fields[1]

	// "UPDATE user SET name = ? WHERE name = ?", "testing", "slene"
	sql := "update " + tabName + " set "
	values := []interface{}{}

	for _, v := range editFields {
		sql += v + " = ? and"
	}
	sql = sql[:len(sql)-4]

	if len(searchFields) > 0 {
		sql += " where "
		for _, v := range searchFields {
			sql += v + " = ? and"
		}
		sql = sql[:len(sql)-4]
	}

	allFields := append(editFields, searchFields...)
	for _, _v := range allFields {
		_type := reflect.TypeOf(edit_data).Elem()
		_value := reflect.ValueOf(edit_data).Elem()

		for i := 0; i < _type.NumField(); i++ {
			key := _type.Field(i).Name
			if key == _v {
				value := _value.Field(i)
				switch value.Kind() {
				case reflect.Int:
					values = append(values, value.Int())
				case reflect.Int8:
					values = append(values, int8(value.Int()))
				case reflect.Int16:
					values = append(values, int16(value.Int()))
				case reflect.Int32:
					values = append(values, int32(value.Int()))
				case reflect.Int64:
					values = append(values, int64(value.Int()))
				case reflect.String:
					values = append(values, value.String())
				case reflect.Bool:
					values = append(values, value.Bool())
				case reflect.Float32, reflect.Float64:
					values = append(values, value.Float())
				}
				break
			}
		}
	}

	o := orm.NewOrm()
	_result, err := o.Raw(sql, values).Exec()

	if ChkErr(err) {
		return nil
	}
	_, err = _result.RowsAffected()
	if ChkErr(err) {
		return nil
	}
	_resultProto := DBFind(edit_data, searchFields)
	return _resultProto
}
func DBFind(search_data any, fields ...[]string) interface{} {
	_result, _resultProto := dbuts.CreateStructFromDBDat(search_data)
	o := orm.NewOrm()
	_querySet := o.QueryTable(_result)
	if fields == nil || len(fields) < 1 || fields[0] == nil || len(fields[0]) < 1 {
		ok, _resultProtoData := dbuts.TabDatIsArray(_resultProto)
		if ok {
			_querySet.All(_resultProtoData)
		}
	} else {
		for _, _v := range fields[0] {
			_type := reflect.TypeOf(search_data).Elem()
			_value := reflect.ValueOf(search_data).Elem()

			for i := 0; i < _type.NumField(); i++ {
				key := _type.Field(i).Name
				if key == _v {
					value := _value.Field(i)
					switch value.Kind() {
					case reflect.Int:
						_querySet = _querySet.Filter(key, value.Int())
					case reflect.Int8:
						_querySet = _querySet.Filter(key, int8(value.Int()))
					case reflect.Int16:
						_querySet = _querySet.Filter(key, int16(value.Int()))
					case reflect.Int32:
						_querySet = _querySet.Filter(key, int32(value.Int()))
					case reflect.Int64:
						_querySet = _querySet.Filter(key, int64(value.Int()))
					case reflect.String:
						_querySet = _querySet.Filter(key, value.String())
					case reflect.Bool:
						_querySet = _querySet.Filter(key, value.Bool())
					case reflect.Float32, reflect.Float64:
						_querySet = _querySet.Filter(key, value.Float())
					}
					break
				}
			}
		}
		ok, _resultProtoData := dbuts.TabDatIsArray(_resultProto)
		if ok {
			_querySet.All(_resultProtoData)
		}
	}
	if dbuts.TabDatIsEmpty(_resultProto) {
		_resultProto = nil
	}
	return _resultProto
}
