package dbcreatestruct

import (
	"baseutils/baseuts"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// var tab_struct_id = 1

var db_utils_python_str = ""

// var tab_struct_str = ""
var tab_struct_type_str = ""
var tab_utils_str = ""
var tab_struct_is_array_str = ""
var tab_struct_is_empty_str = ""
var tab_utils_python_type_str = ""
var tab_field_name_struct_str = ""
var tab_struct_to_redis_str = ""
var tab_struct_factory_str = ""
var tab_struct_index_factory_str = ""
var tab_struct_get_name_str = ""
var tab_print_str = ""
var tab_create_struct_str = ""
var tab_struct_python_str = ""

func CreateDBStructStart() {
	tab_struct_is_array_str = ""
	tab_struct_is_empty_str = ""
	tab_struct_type_str = ""
	tab_utils_str = ""
	tab_field_name_struct_str = ""
	tab_struct_to_redis_str = ""
	tab_struct_factory_str = ""
	tab_struct_index_factory_str = ""
	tab_struct_get_name_str = ""
	tab_print_str = ""
	tab_create_struct_str = ""
	tab_struct_python_str = ""

	_path := "dbcreatestruct/dbstruct"
	if _, err := os.Stat(_path); err == nil || os.IsExist(err) {
		err := os.RemoveAll(_path)
		baseuts.ChkErr(err)
	}
	err := os.Mkdir(_path, os.ModePerm)
	baseuts.ChkErr(err)

	_path = "protobuff/pbpython"
	if _, err := os.Stat(_path); err == nil || os.IsExist(err) {
		err := os.RemoveAll(_path)
		baseuts.ChkErr(err)
	}
	err = os.MkdirAll(_path+"/dbstruct", os.ModePerm)
	baseuts.ChkErr(err)

	tableArray := baseuts.FindAllTable(baseuts.MysqlAccount, baseuts.MysqlPwd, baseuts.MysqlHost, baseuts.MysqlPort, baseuts.MysqlDBName)
	db_utils_python_str = "from dbstruct import "
	db_utils_str := `package dbstruct` + "\n\n"
	db_utils_str += `import (
	"baseutils/baseuts"
	"encoding/json"
	"gameutils/common/commuts"
	"gameutils/common/pbuts"
	"gameutils/pbstruct"
	"reflect"
	_ "unsafe"
)` + "\n\n"
	db_utils_str += `type DBTable struct{}` + "\n\n"
	db_utils_str += `type DBFieldName interface {
	MStruct() interface{}
}` + "\n\n"
	db_utils_str += `var DBInfo = DBTable{}` + "\n\n"
	db_utils_str += `var DBTabName = map[string]DBFieldName{` + "\n"
	for tabName, _ := range tableArray {
		_tableFieldInfos := baseuts.FindDBTableField(baseuts.MysqlAccount, baseuts.MysqlPwd, baseuts.MysqlHost, baseuts.MysqlPort, baseuts.MysqlDBName, tabName)
		createStruct(tabName, _tableFieldInfos)
	}
	db_utils_str += tab_field_name_struct_str
	db_utils_str += `}` + "\n\n"
	db_utils_str += tab_utils_str
	db_utils_str += "//go:linkname createTabAndListFromTabID gameutils/common/dbuts.CreateTabAndListFromTabID\n"
	db_utils_str += `func createTabAndListFromTabID(_index int32, MicroTransDB ...*pbstruct.MicroTransDB) (interface{}, interface{}) {
	var microST *pbstruct.MicroTransDB = nil
	var _dbBuf []byte = nil
	var _dbBufStr = ""
	if len(MicroTransDB) > 0 {
		microST = MicroTransDB[0]
		_index = microST.Id
		_dbBuf = microST.Table
		_dbBufStr = microST.TableStr
	}
	switch _index {` + "\n"
	db_utils_str += tab_struct_index_factory_str
	db_utils_str += `	}
	return nil, nil
}` + "\n\n"
	db_utils_str += `//go:linkname createTabAndListFromTabName gameutils/common/dbuts.CreateTabAndListFromTabName
func createTabAndListFromTabName(_dbName string) (interface{}, interface{}) {
	switch _dbName {` + "\n"
	db_utils_str += tab_struct_factory_str
	db_utils_str += `	}
	return nil, nil
}` + "\n\n"
	db_utils_str += `//go:linkname tabDatIsArray gameutils/common/dbuts.TabDatIsArray
func tabDatIsArray(_dbData interface{}) (bool, interface{}) {
	switch _dbData := _dbData.(type) {` + "\n"
	db_utils_str += tab_struct_is_array_str
	db_utils_str += `	}
	return false, nil
}` + "\n\n"
	db_utils_str += `//go:linkname tabDatIsEmpty gameutils/common/dbuts.TabDatIsEmpty
func tabDatIsEmpty(_dbData interface{}) bool {
	_dbDatRef := reflect.ValueOf(_dbData)
	if _dbDatRef.Kind() == reflect.Ptr && _dbDatRef.IsNil() {
		return true
	} 
	switch _dbData := _dbData.(type) {` + "\n"
	db_utils_str += tab_struct_is_empty_str
	db_utils_str += `	}
	return false
}` + "\n\n"
	db_utils_str += "//go:linkname createStructFromDBDat gameutils/common/dbuts.CreateStructFromDBDat\n"
	db_utils_str += `func createStructFromDBDat(_dbData interface{}) (interface{}, interface{}) {
	switch _dbData.(type) {` + "\n"
	db_utils_str += tab_create_struct_str
	db_utils_str += `	}
	return nil, nil
}` + "\n\n"
	db_utils_str += "//go:linkname getTabNameIDFromDBDat gameutils/common/dbuts.GetTabNameIDFromDBDat\n"
	db_utils_str += `func getTabNameIDFromDBDat(_dbData interface{}) (string, int32) {
	switch _dbData.(type) {
` + tab_struct_get_name_str + `	}
	return "", -1
}` + "\n\n"
	db_utils_str += "//go:linkname getAddMapFromDBDat gameutils/common/dbuts.GetAddMapFromDBDat\n"
	db_utils_str += `func getAddMapFromDBDat(_dbData interface{}, _addElem ...interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	switch _dbData := _dbData.(type) {
` + tab_struct_to_redis_str + `	}
	return result
}` + "\n\n"
	db_utils_str += "//go:linkname printDBResult gameutils/common/dbuts.PrintDBResult\n"
	db_utils_str += `func printDBResult(result interface{}) {
	isPrinted := true
	switch result := result.(type) {
` + tab_print_str + `	default:
		baseuts.Log("找不到类型打印", result, &result)
	}
	if isPrinted {
		baseuts.Log("-----------------")
	}
}` + "\n\n"
	db_utils_str += "//go:linkname struct2Map gameutils/common/commuts.Struct2Map\n"
	db_utils_str += `func struct2Map(_struct interface{}) map[string]interface{} {
	_type := reflect.TypeOf(_struct)
	_value := reflect.ValueOf(_struct)
	_itemMap := commuts.StructReflect2Map(_type, _value)
	return _itemMap
}` + "\n\n"
	// 	db_utils_str += "//go:linkname enGobDB gameutils/common/commuts.EnGobDB\n"
	// 	db_utils_str += `func enGobDB(_data interface{}) []byte {
	// 	var _buf bytes.Buffer
	// 	dec := gob.NewEncoder(&_buf)
	// 	err := dec.Encode(_data)
	// 	baseuts.ChkErr(err)
	// 	return _buf.Bytes()
	// }` + "\n\n"
	// 	db_utils_str += "//go:linkname deGobDB gameutils/common/commuts.DeGobDB\n"
	// 	db_utils_str += `func deGobDB(_data interface{}, _b []byte) {
	// 	if _b != nil {
	// 		var _buf bytes.Buffer
	// 		_buf.Write(_b)
	// 		dec := gob.NewDecoder(&_buf)
	// 		err := dec.Decode(_data)
	// 		baseuts.ChkErr(err)
	// 	}
	// }`
	baseuts.SaveFile("dbcreatestruct/dbstruct/dbutils.go", db_utils_str)

	db_struct_str := `package dbstruct` + "\n\n"
	db_struct_str += `import "gameutils/pbstruct"` + "\n\n"
	db_struct_str += `type TableType interface {
	` + tab_struct_type_str[:len(tab_struct_type_str)-1] + `
}` + "\n\n"
	// db_struct_str += tab_struct_str
	baseuts.SaveFile("dbcreatestruct/dbstruct/dbstruct.go", db_struct_str)
	baseuts.SaveFile("protobuff/pbpython/dbstruct/dbstruct.py", tab_struct_python_str)
	db_utils_python_str = db_utils_python_str[:len(db_utils_python_str)-1] + "\n"
	db_utils_python_str += "def getTabNameIDFromDBDat(_dbData):\n"
	db_utils_python_str += tab_utils_python_type_str
	baseuts.SaveFile("protobuff/pbpython/dbstruct/dbutils.py", db_utils_python_str)
	baseuts.CopyFile("../gameutils/dbstruct", "dbcreatestruct/dbstruct")
	os.RemoveAll("dbcreatestruct/dbstruct")
	baseuts.Log("database struct created!!")
}

func createStruct(tableName string, fieldInfos []baseuts.Field) {
	// _arr_table_name := strings.Split(tableName, "_")
	// _struct_name := ""
	// for _, v := range _arr_table_name {
	// 	_struct_name += strings.ToUpper(v[:1]) + v[1:]
	// }
	_struct_name := strings.ToUpper(tableName[:1]) + tableName[1:]
	db_utils_python_str += _struct_name + ","
	tab_utils_python_type_str += `	if isinstance(_dbData,` + _struct_name + `):
		return "` + tableName + `", ` + _struct_name + "_ID\n"
	tab_field_name_struct_str += `	"` + tableName + `": ` + _struct_name + `FieldName{` + "\n"
	var _table_struct_str = ""
	var _tab_struct_python_str = ""
	// var _tab_struct_python_tojson_str = `	def toJson(self):
	// 	return {`
	var _field_name_struct_str = ""
	for _, v := range fieldInfos {
		// _arr_field_dame := strings.Split(v.FieldName, "_")
		// _field_name := ""
		// for _, vv := range _arr_field_dame {
		// 	_field_name += strings.ToUpper(vv[:1]) + vv[1:]
		// }
		_field_name := strings.ToUpper(v.FieldName[:1]) + v.FieldName[1:]
		_table_struct_str += `	` + _field_name + ` ` + typeDB2Go(v.DataType) + "\n"
		_defaultValue := "None"
		if v.FieldDefault != "" {
			_defaultValue = v.FieldDefault
		}
		_tab_struct_python_str += "		self." + _field_name + " = " + _defaultValue + "\n"
		// _tab_struct_python_tojson_str += "\"" + _field_name + "\":self." + _field_name + ","
		_field_name_struct_str += `	` + _field_name + ` string` + "\n"
		tab_field_name_struct_str += `		` + _field_name + `:    "` + v.FieldName + `",` + "\n"
	}

	var pbstructName = "pbstruct." + baseuts.UnderScoreCase2CamelCase(_struct_name)

	tab_struct_to_redis_str += `	case *` + pbstructName + `:
		result = append(result, struct2Map(_dbData))
	case *` + pbstructName + `Result:
		for _, _v := range _addElem {
			_dbData.Data = append(_dbData.Data, _v.(*` + pbstructName + `))
		}
		for _, value := range _dbData.Data {
			result = append(result, struct2Map(value))
		}` + "\n"
	tab_struct_get_name_str += `	case *` + pbstructName + `:
		return "` + tableName + `", ` + pbstructName + `_ID
	case *` + pbstructName + `Result:
		return "` + tableName + `", ` + pbstructName + "Result_ID\n"
	tab_struct_index_factory_str += `	case ` + pbstructName + "_ID, " + pbstructName + `Result_ID:	
		result := ` + pbstructName + `{}
		resultList := ` + pbstructName + `Result{}
		if _dbBufStr != "" {
			var err error
			if _index == pbstruct.UserTab_ID {
				err = json.Unmarshal([]byte(_dbBufStr), &result)
			} else {
				err = json.Unmarshal([]byte(_dbBufStr), &resultList)
			}
			if baseuts.ChkErr(err) {
				return nil, nil
			}
		} else if _dbBuf != nil {
			if _index == pbstruct.UserTab_ID {
				pbuts.ProtoUnMarshal(_dbBuf, &result)
			} else {
				pbuts.ProtoUnMarshal(_dbBuf, &resultList)
			}
		}
		return &result, &resultList` + "\n"
	tab_field_name_struct_str += `	},` + "\n"
	tab_struct_is_array_str += `	case *` + pbstructName + `Result:
		return true, &_dbData.Data` + "\n"
	tab_struct_is_empty_str += `	case *` + pbstructName + `Result:
		if len(_dbData.Data) < 1 {
			return true
		}` + "\n"
	tab_create_struct_str += `	case *` + pbstructName + `, *` + pbstructName + `Result:
		_result := ` + pbstructName + `{}
		_resultList := ` + pbstructName + `Result{Data: []*` + pbstructName + `{}}
		return &_result, &_resultList` + "\n"
	tab_print_str += `	case *` + pbstructName + `:
		baseuts.Log(struct2Map(result))
	case *` + pbstructName + `Result:
		isPrinted = false
		for _, v := range result.Data {
			isPrinted = true
			baseuts.Log(struct2Map(v))
		}` + "\n"
	tab_struct_factory_str += `	case "` + _struct_name + `", "` + strings.ToLower(_struct_name) + `":
		result := ` + pbstructName + `{}
		resultList := ` + pbstructName + `Result{}
		return &result, &resultList` + "\n"
	tab_struct_type_str += "*" + pbstructName + " | *" + pbstructName + "Result |"

	tab_struct_python_str += "class " + _struct_name + "(object):\n"
	tab_struct_python_str += `	"""docstring for ` + _struct_name + `"""
	def __init__(self):
		super(` + _struct_name + `, self).__init__()` + "\n"
	tab_struct_python_str += _tab_struct_python_str + "\n"
	// tab_struct_python_str += _tab_struct_python_tojson_str[:len(_tab_struct_python_tojson_str)-1] + "}"
	// tab_struct_python_str += "\n\n"

	// tab_struct_str += `type ` + _struct_name + ` struct {` + "\n"
	// tab_struct_str += _table_struct_str
	// tab_struct_str += `}` + "\n\n"
	tab_utils_str += `type ` + _struct_name + `FieldName struct {` + "\n"
	tab_utils_str += _field_name_struct_str
	tab_utils_str += `}` + "\n\n"
	tab_utils_str += `func (` + _struct_name + `FieldName) MStruct() interface{} {
	return &pbstruct.` + baseuts.UnderScoreCase2CamelCase(_struct_name) + `{}
}` + "\n\n"
	tab_utils_str += `func (DBTable) ` + _struct_name + `(get_name ...bool) any {
	var result any
	str := "` + tableName + `"
	if len(get_name) > 0 && get_name[0] {
		result = str
	} else {
		result = DBTabName[str]
	}
	return result
}` + "\n\n"
	// tab_struct_id++
}

func typeDB2Go(db_type string) string {
	switch db_type {
	case "varchar", "langtext":
		return "string"
	case "bigint":
		return "int64"
	case "int":
		return "int"
	case "tinyint":
		return "int8"
	case "double":
		return "float64"
	}
	return db_type
}
