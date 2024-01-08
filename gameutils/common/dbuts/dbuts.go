package dbuts

import (
	"baseutils/baseuts"
	"encoding/json"
	"gameutils/common/pbuts"
	"gameutils/pbstruct"

	"reflect"
	_ "unsafe"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/encoding/protojson"
)

func CreateTabAndListFromTabName(_dbName string) (interface{}, interface{})

func CreateTabAndListFromTabID(_index int32, MicroTransDB ...*pbstruct.MicroTransDB) (interface{}, interface{})

func CreateStructFromDBDat(_dbData interface{}) (interface{}, interface{})

func GetTabNameIDFromDBDat(_dbData interface{}) (string, int32)

func GetAddMapFromDBDat(_dbData interface{}, _addElem ...interface{}) []map[string]interface{}

func TabDatIsArray(_dbData interface{}) (bool, interface{})

func TabDatIsEmpty(_dbData interface{}) bool

func PrintDBResult(result interface{})

// 数组relatedField 下标0:SearchField 下载1:EditField
func GetMTDBFromDBDat(_dbData interface{}, act string, relatedField ...[]string) *pbstruct.MicroTransDB {
	_data := getMTDBFromDBDatBase(false, _dbData, act, relatedField...)
	return _data
}

func GetMTDBFromDBDatJson(_dbData interface{}, act string, relatedField ...[]string) *pbstruct.MicroTransDB {
	_data := getMTDBFromDBDatBase(true, _dbData, act, relatedField...)
	return _data
}

func getMTDBFromDBDatBase(isJson bool, _dbData interface{}, act string, relatedField ...[]string) *pbstruct.MicroTransDB {
	var _id int32 = -1
	var _isArray = false

	var _type = reflect.TypeOf(_dbData)
	var _typeName = _type.Kind().String()

	if _type.Kind() == reflect.Pointer {
		_type = reflect.TypeOf(_dbData).Elem()
		_typeName = _type.Kind().String()
	}
	if _typeName == "slice" {
		_isArray = true
	}
	_, _id = GetTabNameIDFromDBDat(_dbData)

	_data := pbstruct.MicroTransDB{}
	_data.Id = _id
	_data.Act = act
	_data.IsArray = _isArray
	if isJson {
		_b, err := json.Marshal(_dbData)
		if baseuts.ChkErr(err) {
			_data.TableStr = string(_b)
		}
	} else {
		if _type.Kind() == reflect.Pointer {
			_data.Table, _ = pbuts.ProtoMarshal(&_dbData)
		} else {
			_data.Table, _ = pbuts.ProtoMarshal(_dbData)
		}
	}
	_data.SearchField = []string{}
	_data.EditField = []string{}
	if len(relatedField) > 0 {
		_data.SearchField = relatedField[0]
	}
	if len(relatedField) > 1 {
		_data.EditField = relatedField[1]
	}
	return &_data
}

func GetMTDBBufFromDBDat(_dbData interface{}, act string, relatedField ...[]string) []byte {
	b := getMTDBBufFromDBDatBase(false, _dbData, act, relatedField...)
	return b.([]byte)
}

func GetMTDBBufFromDBDatJson(_dbData interface{}, act string, relatedField ...[]string) string {
	jsonStr := getMTDBBufFromDBDatBase(true, _dbData, act, relatedField...)
	return jsonStr.(string)
}

func getMTDBBufFromDBDatBase(isJson bool, _dbData interface{}, act string, relatedField ...[]string) interface{} {
	if isJson {
		_data := GetMTDBFromDBDatJson(_dbData, act, relatedField...)
		str := protojson.Format(_data)
		return str
	} else {
		_data := GetMTDBFromDBDat(_dbData, act, relatedField...)
		b, _ := pbuts.ProtoMarshal(_data)
		return b
	}
}
