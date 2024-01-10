package model

import (
	"baseutils/baseuts"
	"errors"
	"fmt"
	"gameutils/common/commuts"
	"gameutils/common/dbuts"
	"gameutils/dbstruct"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/orm"
	"github.com/go-redis/redis"
)

var _rdb *redis.Client

var _rdbSlave *redis.Client

func rdb() *redis.Client {
	if checkRedis(_rdb) {
		return _rdb
	} else {
		Log("redis 异常")
		return nil
	}
}

func rdbSlave() *redis.Client {
	if checkRedis(_rdb) {
		if checkRedis(_rdbSlave) {
			return _rdbSlave
		} else {
			return _rdb
		}
	} else {
		Log("redis 异常")
		return nil
	}
}

func initRedis() {
	redisConn()
	redisSlaveConn()
}

func redisConn() {
	_rdb = redis.NewClient(&redis.Options{
		Addr:     Conf.String("redis_host") + ":" + Conf.String("redis_port"),
		Password: "root2023",
		//连接信息
		Network: "tcp", //网络类型，tcp or unix，默认tcp
		DB:      0,     // redis数据库index
		//连接池容量及闲置连接数量
		// PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		//超时
		// DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		// ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		// WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		// PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。
		//闲置连接检查包括IdleTimeout，MaxConnAge
		//IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		//IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		//MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
		//命令执行失败时的重试策略
		MaxRetries: 0, // 命令执行失败时，最多重试多少次，默认为0即不重试
		// MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		// MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//可自定义连接函数
		// Dialer: func() (net.Conn, error) {
		// 	netDialer := &net.Dialer{
		// 		Timeout:   5 * time.Second,
		// 		KeepAlive: 5 * time.Minute,
		// 	}
		// 	return netDialer.Dial("tcp", "127.0.0.1:6379")
		// },

		// //钩子函数
		// OnConnect: func(conn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
		// 	Log("conn=%v\n", conn)
		// 	return nil
		// },
	})
	if checkRedis(_rdb) {
		Log("创建Redis完成")
		getTableDefaultValue()
	} else {
		Log("Redis重试")
		time.Sleep(time.Second * 6)
		redisConn()
	}
}

func redisSlaveConn() {
	_rdbSlave = redis.NewClient(&redis.Options{
		Addr:     Conf.String("redis_slave_host") + ":" + Conf.String("redis_slave_port"),
		Password: "root2023",
		//连接信息
		Network:      "tcp", //网络类型，tcp or unix，默认tcp
		DB:           0,     // redis数据库index
		MinIdleConns: 10,
	})

	if checkRedis(_rdbSlave) {
		Log("创建Redis Slave完成")
	} else {
		Log("Redis重试")
		time.Sleep(time.Second * 6)
		redisSlaveConn()
	}
}

const minInterval int64 = 5

var prevCheckTime int64 = 0
var prevCheckRedis *redis.Client

func checkRedis(r *redis.Client) bool {
	currCheckTime := time.Now().Unix()
	if r == prevCheckRedis && currCheckTime-prevCheckTime < minInterval {
		prevCheckTime = currCheckTime
		prevCheckRedis = r
		return true
	}
	prevCheckTime = currCheckTime
	prevCheckRedis = r

	// 测试连接
	_, err := r.Ping().Result()
	return !baseuts.ChkErrNormal(err, "redis 检测异常")
}

var tabDefaultValue = map[string]map[string]baseuts.Field{}

func getTableDefaultValue() {
	allTable := baseuts.FindAllTable(Conf.String("db_account"), Conf.String("db_pwd"), "127.0.0.1", Conf.String("db_port"), Conf.String("db_name"))
	for tabName := range allTable {
		fields := baseuts.FindDBTableField(Conf.String("db_account"), Conf.String("db_pwd"), "127.0.0.1", Conf.String("db_port"), Conf.String("db_name"), tabName)
		var tabFields = map[string]baseuts.Field{}
		for _, v := range fields {
			if v.FieldName != "id" {
				tabFields[v.FieldName] = v
			}
		}
		tabDefaultValue[tabName] = tabFields
	}

}

var secondMainKey = map[string]string{"Inventory_id": "inventory_id"}
var keyLevel1 = map[string]string{"User_id": "user_id"}

func dB2Cache() {
	rdb().FlushAll()
	Log("DB2Cache 导入")
	// Expire 命令用于设置 key 的过期时间，key 过期后将不再可用。单位以秒计。
	// PERSIST 命令用于移除给定 key 的过期时间，使得 key 永不过期。
	// TTL 命令以秒为单位返回 key 的剩余过期时间

	for _tabName, v := range dbstruct.DBTabName {
		_ret := DBFind(v.MStruct())
		_r := dbuts.GetAddMapFromDBDat(_ret)

		var _incrKey int32 = 0
		for _, _map := range _r {

			_mKey := ""
			_mKey2 := ""
			_levelKey1 := ""

			for _k, _v := range _map {
				if _k == "Id" {
					_idValue := _v.(int32)
					_incrKey = max(_incrKey, _idValue)
					_mKey = strconv.FormatInt(int64(_idValue), 10)
				}

				if _tabKey, ok := secondMainKey[_k]; ok {
					switch _dat := _v.(type) {
					case int, int8, int16, int32, int64:
						_mKey2 = _tabKey + ":" + strconv.FormatInt(_dat.(int64), 10)
					case string:
						_mKey2 = _tabKey + ":" + _dat
					}
				}

				if _tabKey, ok := keyLevel1[_k]; ok {
					switch _dat := _v.(type) {
					case int, int8, int16, int32, int64:
						_levelKey1 = _tabKey + ":" + strconv.FormatInt(_dat.(int64), 10)
					case string:
						_levelKey1 = _tabKey + ":" + _dat
					}
				}
			}
			if _mKey == "" {
				baseuts.LogF("db2Cache", "表 "+_tabName+" 缺少索引Id")
			} else {
				_rootKey := _tabName + ":" + _mKey
				for _fk, _fv := range _map {
					// if IsDebug {
					// 	LogDebug("DB2Cache Set", _rootKey, _fk, _fv)
					// }
					_, err := rdb().HSet(_rootKey, strings.ToLower(_fk), _fv).Result()
					ChkErr(err)
				}
				rdb().RPush(_tabName, _rootKey)
				var err error
				if _mKey2 != "" {
					_, err = rdb().Set(_tabName+":"+_mKey2, _rootKey, 0).Result()
				} else if _levelKey1 != "" {
					_, err = rdb().RPush(_tabName+":"+_levelKey1, _rootKey).Result()
				}
				ChkErr(err)
			}
			_, err := rdb().Set(_tabName+"Incr", _incrKey, 0).Result()
			ChkErr(err, "incr")
		}
	}
	Log("DB2Cache 结束")

	go func() {
		defer baseuts.ChkRecover()
		RDDump()
	}()
}

func BackUp() {
	baseuts.BackupMysql(Conf.String("db_account"), Conf.String("db_pwd"), Conf.String("db_name"))
}

func Cache2DB() {
	BackUp()
	tableData := baseuts.FindAllTable(Conf.String("db_account"), Conf.String("db_pwd"), "127.0.0.1", Conf.String("db_port"), Conf.String("db_name"))
	for k, v := range dbstruct.DBTabName {
		_tabDesc := strings.Split(tableData[k], "#")
		if len(_tabDesc) > 1 && _tabDesc[1] == "static" {
			continue
		}
		o := orm.NewOrm()
		_, _err := o.Raw("delete from " + k).Exec()
		ChkErr(_err)
		_ret, _ := rdb().Keys(k + ":*").Result()
		for _, _v := range _ret {
			_keyArray := strings.Split(_v, ":")
			if len(_keyArray) > 2 {
				continue
			}
			__ret, err := rdb().HGetAll(_v).Result()
			if ChkErr(err) || len(__ret) < 1 {
				continue
			}
			_struct := v.MStruct()
			_type := reflect.TypeOf(_struct).Elem()
			_value := reflect.ValueOf(_struct).Elem()
			_editField := []string{}
			for i := 0; i < _value.NumField(); i++ {
				__type := _type.Field(i)
				__value := _value.Field(i)
				_mValue := __ret[strings.ToLower(__type.Name)]
				_editField = append(_editField, __type.Name)
				switch __value.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					_r, _ := strconv.ParseInt(_mValue, 10, 64)
					__value.SetInt(_r)
				case reflect.String:
					__value.SetString(_mValue)
				case reflect.Bool:
					_r, _ := strconv.ParseBool(_mValue)
					__value.SetBool(_r)
				case reflect.Float32, reflect.Float64:
					_r, _ := strconv.ParseFloat(_mValue, 64)
					__value.SetFloat(_r)
				}
			}
			DBAdd(_struct, []string{}, _editField)
		}
	}
	Log("redis to mysql success!!")
}

func RDAdd(data interface{}, fields ...[]string) interface{} {
	if fields == nil || len(fields) < 2 || fields[1] == nil || len(fields[1]) < 1 {
		Log("RDADD", "参数不足")
		return nil
	}
	_time := time.Now().UnixMilli()
	_result, _resultProto := dbuts.CreateStructFromDBDat(data)

	if rdb() != nil {
		editFields := fields[1]

		tabName, _ := dbuts.GetTabNameIDFromDBDat(data)

		tabFieldsDefaultValue := tabDefaultValue[tabName]
		_incr, _ := rdb().Incr(tabName + "Incr").Result()
		_id := strconv.FormatInt(_incr, 10)

		_rootKey := tabName + ":" + _id

		var _secondMainKey string
		var _secondMainKeyValue interface{}

		var _keyLevel1Key string
		var _keyLevel1KeyValue interface{}

		for k, v := range tabFieldsDefaultValue {
			_value := baseuts.GetFieldsValue(v, true)
			if _secondMainKey == "" || _keyLevel1Key == "" {
				_kToUpper := strings.ToUpper(k[:1]) + k[1:]
				if _dat, ok := secondMainKey[_kToUpper]; ok {
					_secondMainKey = _dat
					_secondMainKeyValue = _value
				} else if _dat, ok := keyLevel1[_kToUpper]; ok {
					_keyLevel1Key = _dat
					_keyLevel1KeyValue = _value
				}
			}
			_, err := rdb().HSet(_rootKey, k, _value).Result()
			ChkErr(err)
		}

		rdb().RPush(tabName, _rootKey)

		var _err error

		_type := reflect.TypeOf(data).Elem()
		_value := reflect.ValueOf(data).Elem()

		for i := 0; i < _type.NumField(); i++ {
			key := _type.Field(i).Name
			if key == "state" || key == "sizeCache" || key == "unknownFields" {
				continue
			}
			_hasKey := false
			for _, _v := range editFields {
				if strings.EqualFold(key, _v) {
					_hasKey = true
					break
				}
			}

			value := _value.Field(i)
			if key == "Id" {
				value.SetInt(_incr)
			} else if _hasKey {
				_valueStr := ""
				redisKey := strings.ToLower(key)
				switch value.Kind() {
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					_value := value.Int()
					_valueStr = strconv.FormatInt(_value, 10)
					_, _err = rdb().HSet(_rootKey, redisKey, _value).Result()
				case reflect.String:
					_valueStr = value.String()
					_, _err = rdb().HSet(_rootKey, redisKey, _valueStr).Result()
				case reflect.Bool:
					_value := value.Bool()
					_valueStr = strconv.FormatBool(_value)
					_, _err = rdb().HSet(_rootKey, redisKey, _value).Result()
				}
				if !ChkErr(_err) {
					if redisKey == strings.ReplaceAll(_secondMainKey, "_", "") {
						_secondMainKeyValue = _valueStr
					} else if redisKey == strings.ReplaceAll(_keyLevel1Key, "_", "") {
						_keyLevel1KeyValue = _valueStr
					}
				}
			}
		}
		_, _err = rdb().HSet(_rootKey, "id", _id).Result()
		ChkErr(_err)

		_mapRedis, _err := rdb().HGetAll(tabName + ":" + _id).Result()
		if ChkErr(_err) || len(_mapRedis) < 1 {
			_resultProto = nil
		} else {
			if _secondMainKey != "" {
				Log(tabName+":"+_secondMainKey, _secondMainKeyValue)
				rdb().Set(tabName+":"+_secondMainKey+":"+_secondMainKeyValue.(string), _rootKey, 0)
			} else if _keyLevel1Key != "" {
				Log(tabName+":"+_keyLevel1Key, _keyLevel1KeyValue)
				rdb().RPush(tabName+":"+_keyLevel1Key+":"+_keyLevel1KeyValue.(string), _rootKey)
			}
			commuts.Map2Struct(_result, _mapRedis, true)
			dbuts.GetAddMapFromDBDat(_resultProto, _result)
		}
	} else {
		_resultProto = DBAdd(data, fields...)
	}
	LogDebug("RDAdd 用时:", time.Now().UnixMilli()-_time)
	return _resultProto
}

func RDDel(data interface{}, fields ...[]string) interface{} {
	if fields == nil || len(fields) < 1 || fields[0] == nil || len(fields[0]) < 1 {
		Log("RDDel", "参数不足")
		return nil
	}
	_time := time.Now().UnixMilli()

	_result, _resultProto := dbuts.CreateStructFromDBDat(data)

	if rdb() != nil {
		tabName, _ := dbuts.GetTabNameIDFromDBDat(data)
		searchFields := fields[0]

		var _hasID = false

		for i := range searchFields {
			if strings.EqualFold("Id", searchFields[i]) {
				_hasID = true
				break
			}
		}

		_idValueStr := ""

		if _hasID {
			_type := reflect.TypeOf(data).Elem()
			_value := reflect.ValueOf(data).Elem()

			for i := 0; i < _type.NumField(); i++ {
				key := _type.Field(i).Name
				if strings.EqualFold("Id", key) {
					value := _value.Field(i)
					_idValueStr = strconv.FormatInt(value.Int(), 10)
					break
				}
			}

			_mapRedis, err := rdb().HGetAll(tabName + ":" + _idValueStr).Result()
			if !ChkErr(err) && len(_mapRedis) > 0 {
				commuts.Map2Struct(_result, _mapRedis, true)
				dbuts.GetAddMapFromDBDat(_resultProto, _result)
			} else {
				_resultProto = nil
			}
		} else {
			_ret, err := rdbSlave().LRange(tabName, 0, -1).Result()
			if !ChkErr(err) {
				if len(_ret) > 0 {
					_result, _ = findDataFromRedisKeys(_ret, _resultProto, data, searchFields)
				} else {
					_resultProto = nil
				}
			} else {
				_resultProto = nil
			}
		}

		if _result == nil {
			Log("RDDel", "没找到数据")
			_resultProto = nil
		} else {
			var _secondMainKey string
			var _secondMainKeyValue interface{}

			var _keyLevel1Key string
			var _keyLevel1KeyValue interface{}

			_type := reflect.TypeOf(_result).Elem()
			_value := reflect.ValueOf(_result).Elem()
			for i := 0; i < _type.NumField(); i++ {
				key := _type.Field(i).Name
				if key == "state" || key == "sizeCache" || key == "unknownFields" {
					continue
				}
				value := _value.Field(i)
				_valueStr := ""

				switch value.Kind() {
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					_value := value.Int()
					_valueStr = strconv.FormatInt(_value, 10)
				case reflect.String:
					_valueStr = value.String()
				case reflect.Bool:
					_value := value.Bool()
					_valueStr = strconv.FormatBool(_value)
				}

				if key == "Id" {
					_idValueStr = _valueStr
				} else if _secondMainKey == "" || _keyLevel1Key == "" {
					_kToUpper := strings.ToUpper(key[:1]) + key[1:]
					_kToUpper = baseuts.CamelCase2UnderScoreCase(_kToUpper)
					if _dat, ok := secondMainKey[_kToUpper]; ok {
						_secondMainKey = _dat
						_secondMainKeyValue = _valueStr
					} else if _dat, ok := keyLevel1[_kToUpper]; ok {
						_keyLevel1Key = _dat
						_keyLevel1KeyValue = _valueStr
					}
				}
				if _secondMainKey != "" && _keyLevel1Key != "" {
					break
				}
			}

			_, err := rdb().Del(tabName + ":" + _idValueStr).Result()
			if !ChkErr(err) {
				Log("删除", tabName, _idValueStr)
				if _secondMainKey != "" {
					_delKeyStr := tabName + ":" + _secondMainKey + ":" + _secondMainKeyValue.(string)
					_, err := rdb().Del(_delKeyStr).Result()
					if !ChkErr(err) {
						Log("删除Key数据", _delKeyStr, _idValueStr)
					}
				} else if _keyLevel1Key != "" {
					_delKeyStr := tabName + ":" + _keyLevel1Key + ":" + _keyLevel1KeyValue.(string)
					_, err := rdb().LRem(_delKeyStr, 0, tabName+":"+_idValueStr).Result()
					if !ChkErr(err) {
						Log("删除List中数据", _delKeyStr, _idValueStr)
					}
				}
			}
		}
	} else {
		_resultProto = DBDel(data, fields...)
	}

	LogDebug("RDDel 用时", time.Now().UnixMilli()-_time)
	return _resultProto
}
func RDEdit(data interface{}, fields ...[]string) interface{} {
	if fields == nil || len(fields) < 2 || fields[0] == nil || len(fields[0]) < 1 || fields[1] == nil || len(fields[1]) < 1 {
		Log("RDEdit", "参数不足")
		return nil
	}

	_time := time.Now().UnixMilli()

	searchFields := fields[0]
	editFields := fields[1]
	for _, v := range editFields {
		if strings.EqualFold("Id", v) {
			Log("RDEdit", "不能编辑主键ID")
			return nil
		}
	}
	_result, _resultProto := dbuts.CreateStructFromDBDat(data)

	if rdb() != nil {
		var _hasID = false
		for i := range searchFields {
			if strings.EqualFold("Id", searchFields[i]) {
				_hasID = true
				break
			}
		}
		var _idValueStr = ""
		var _mapRedis map[string]string
		var _err error
		_map := dbuts.GetAddMapFromDBDat(data)
		tabName, _ := dbuts.GetTabNameIDFromDBDat(data)
		if _hasID {
			var _idValueStr string
			_type := reflect.TypeOf(data).Elem()
			_value := reflect.ValueOf(data).Elem()

			for i := 0; i < _type.NumField(); i++ {
				key := _type.Field(i).Name
				if strings.EqualFold("Id", key) {
					value := _value.Field(i)
					_idValueStr = strconv.FormatInt(value.Int(), 10)
					break
				}
			}
			_mapRedis, _err = rdb().HGetAll(tabName + ":" + _idValueStr).Result()
		} else {
			_ret, err := rdbSlave().LRange(tabName, 0, -1).Result()
			if !ChkErr(err) {
				_result, _mapRedis = findDataFromRedisKeys(_ret, nil, data, searchFields)
			}
		}

		if ChkErr(_err) || len(_mapRedis) < 1 {
			Log("找不到要编辑的数据")
			_resultProto = nil
		} else {
			var _keyLevel1Key string
			var _keyLevel1KeyValueOld string
			var _keyLevel1KeyValueNew string

			for _, v := range editFields {
				if _dat, ok := keyLevel1[baseuts.CamelCase2UnderScoreCase(v)]; ok {
					_keyLevel1Key = _dat
					break
				}
			}
			if _keyLevel1Key != "" {
				_keyLevel1KeyValueOld = _mapRedis[_keyLevel1Key]
			}
			_idValueStr = _mapRedis["id"]
			for _, v := range editFields {
				setMapValue(tabName+":"+_idValueStr, v, _map[0])
			}
			_mapRedis, err := rdb().HGetAll(tabName + ":" + _idValueStr).Result()
			if ChkErr(err) || len(_mapRedis) < 1 {
				Log("RDEdit", "查询ID无数据返回", _idValueStr)
				_resultProto = nil
			} else {
				if _keyLevel1KeyValueOld != "" {
					_keyLevel1KeyValueNew = _mapRedis[_keyLevel1Key]
					if _keyLevel1KeyValueNew != "" && _keyLevel1KeyValueNew != _keyLevel1KeyValueOld {
						_, err := rdb().LRem(tabName+":"+_keyLevel1Key+":"+_keyLevel1KeyValueOld, 0, tabName+":"+_idValueStr).Result()
						if !ChkErr(err) {
							_, err := rdb().RPush(tabName+":"+_keyLevel1Key+":"+_keyLevel1KeyValueNew, tabName+":"+_idValueStr).Result()
							ChkErr(err)
						}
					}
				}
				commuts.Map2Struct(_result, _mapRedis, true)
				dbuts.GetAddMapFromDBDat(_resultProto, _result)
			}
		}
	} else {
		_resultProto = DBEdit(data, fields...)
	}
	LogDebug("RDEdit 用时", time.Now().UnixMilli()-_time)
	return _resultProto
}

func RDFind(data interface{}, fields ...[]string) interface{} {
	_time := time.Now().UnixMilli()
	_, _resultProto := dbuts.CreateStructFromDBDat(data)
	if rdbSlave() != nil {
		tabName, _ := dbuts.GetTabNameIDFromDBDat(data)
		_ret, err := rdbSlave().LRange(tabName, 0, -1).Result()
		if !ChkErr(err) {
			if fields == nil || fields[0] == nil {
				for _, _v := range _ret {
					_mapRedis, err := rdbSlave().HGetAll(_v).Result()
					if !ChkErr(err) && len(_mapRedis) > 0 {
						_result, _ := dbuts.CreateStructFromDBDat(data)
						commuts.Map2Struct(_result, _mapRedis, true)
						dbuts.GetAddMapFromDBDat(_resultProto, _result)
					}
				}
			} else {
				_map := dbuts.GetAddMapFromDBDat(data)
				searchFields := fields[0]
				_findKey := ""
				_multFindKey := ""
				if len(searchFields) == 1 {
					if _dat, ok := secondMainKey[baseuts.CamelCase2UnderScoreCase(searchFields[0])]; ok {
						var _findKeyValue = _map[0][searchFields[0]]
						switch _findKeyValue := _findKeyValue.(type) {
						case int, int8, int16, int32, int64:
							_findKey = _dat + ":" + strconv.FormatInt(_findKeyValue.(int64), 10)
						case string:
							_findKey = _dat + ":" + _findKeyValue
						}
					} else if _dat, ok := keyLevel1[baseuts.CamelCase2UnderScoreCase(searchFields[0])]; ok {
						var _findKeyValue = _map[0][searchFields[0]]
						switch _findKeyValue := _findKeyValue.(type) {
						case int, int8, int16, int32, int64:
							_multFindKey += _dat + ":" + strconv.FormatInt(int64(_findKeyValue.(int32)), 10)
						case string:
							_multFindKey += _dat + ":" + _findKeyValue
						}
					}

					if _findKey != "" {
						_currKey := tabName + ":" + _findKey
						_rootKey, err := rdbSlave().Get(_currKey).Result()
						if !baseuts.ChkErrNormal(err, "找不到Key:"+_currKey) {
							_mapRedis, err := rdbSlave().HGetAll(_rootKey).Result()
							if !ChkErr(err) && len(_mapRedis) > 0 {
								_result, _ := dbuts.CreateStructFromDBDat(data)
								commuts.Map2Struct(_result, _mapRedis, true)
								dbuts.GetAddMapFromDBDat(_resultProto, _result)
							}
						} else {
							_resultProto = nil
						}
					} else if _multFindKey != "" {
						_ret, _err := rdbSlave().LRange(tabName+":"+_multFindKey, 0, -1).Result()
						if !ChkErr(_err) {
							for _, _v := range _ret {
								_mapRedis, err := rdbSlave().HGetAll(_v).Result()
								if !ChkErr(err) {
									_result, _ := dbuts.CreateStructFromDBDat(data)
									commuts.Map2Struct(_result, _mapRedis, true)
									dbuts.GetAddMapFromDBDat(_resultProto, _result)
								}
							}
						} else {
							_resultProto = nil
						}
					} else {
						if searchFields[0] == "Id" {
							_rootKey := tabName + ":" + strconv.FormatInt(_map[0]["Id"].(int64), 10)
							_mapRedis, err := rdbSlave().HGetAll(_rootKey).Result()
							if !ChkErr(err) {
								_result, _ := dbuts.CreateStructFromDBDat(data)
								commuts.Map2Struct(_result, _mapRedis, true)
								dbuts.GetAddMapFromDBDat(_resultProto, _result)
							} else {
								_resultProto = nil
							}
						} else {
							findDataFromRedisKeys(_ret, _resultProto, data, searchFields)
						}
					}
				} else {
					findDataFromRedisKeys(_ret, _resultProto, data, searchFields)
				}
			}
		}
		if dbuts.TabDatIsEmpty(_resultProto) {
			_resultProto = nil
		}
	} else {
		_resultProto = DBFind(data, fields...)
	}
	LogDebug("RDFind 用时", time.Now().UnixMilli()-_time)
	return _resultProto
}

func RDFindAutoAdd(data interface{}, fields ...[]string) interface{} {
	_result := RDFind(data, fields...)
	if _result == nil {
		_result = RDAdd(data, fields...)
	}
	return _result
}

func findDataFromRedisKeys(keys []string, resultProto interface{}, searchDat interface{}, searchFields []string) (interface{}, map[string]string) {
	var _isMatch = true
	var _resultFirstDat interface{}
	var _resultFirstRedisDat map[string]string
	var _map = dbuts.GetAddMapFromDBDat(searchDat)
	for _, _v := range keys {
		_isMatch = true
		_mapRedis, _err := rdbSlave().HGetAll(_v).Result()
		if ChkErr(_err) {
			_isMatch = false
		} else {
			if len(_mapRedis) > 0 {
				for _, v := range searchFields {
					if !commuts.InterfaceCompareString(_map[0][v], _mapRedis[strings.ToLower(v)]) {
						_isMatch = false
						break
					}
				}
			} else {
				_isMatch = false
			}

		}
		if _isMatch {
			_result, _ := dbuts.CreateStructFromDBDat(searchDat)
			if resultProto != nil {
				commuts.Map2Struct(_result, _mapRedis, true)
				dbuts.GetAddMapFromDBDat(resultProto, _result)
			}
			if _resultFirstDat == nil {
				_resultFirstDat = _result
				_resultFirstRedisDat = _mapRedis
			}
		}
	}
	return _resultFirstDat, _resultFirstRedisDat
}

func RDDump(now ...bool) {
	if len(now) < 1 {
		time.Sleep(time.Hour * 2)
		Log("Redis 数据入库")
		Cache2DB()
		rdb().BgSave()
		RDDump()
	} else {
		RDDump()
	}
}

func setMapValue(key string, field string, mMap map[string]interface{}) {
	errStr := field + ": %v"
	errStr = fmt.Sprintf(errStr, mMap)
	if _mV, ok := mMap[field]; ok {
		var err error = errors.New("未知类型 " + errStr)
		_mVV := reflect.ValueOf(_mV)
		switch _mVV.Kind() {
		case reflect.Int:
			_, err = rdb().HSet(key, strings.ToLower(field), int(_mVV.Int())).Result()
		case reflect.Int8:
			_, err = rdb().HSet(key, strings.ToLower(field), int8(_mVV.Int())).Result()
		case reflect.Int16:
			_, err = rdb().HSet(key, strings.ToLower(field), int16(_mVV.Int())).Result()
		case reflect.Int32:
			_, err = rdb().HSet(key, strings.ToLower(field), int32(_mVV.Int())).Result()
		case reflect.Int64:
			_, err = rdb().HSet(key, strings.ToLower(field), _mVV.Int()).Result()
		case reflect.String:
			_, err = rdb().HSet(key, strings.ToLower(field), _mVV.String()).Result()
		case reflect.Bool:
			_, err = rdb().HSet(key, strings.ToLower(field), _mVV.Bool()).Result()
		}
		ChkErr(err)
	} else {
		Log("未知字段 " + errStr)
	}
}
