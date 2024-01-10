package execluts

import (
	"baseutils/baseuts"
	"gametools/confige"
	"io/fs"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

var ef *excelize.File
var holder interface{}

func fileIsStatic(fileName string, findPath string, result *bool) {
	if *result {
		return
	}
	fileInfoList, err := os.ReadDir(findPath)
	if !baseuts.ChkErr(err, findPath) {
		for i := range fileInfoList {
			file := fileInfoList[i]
			file_path := findPath + "/" + file.Name()
			low_file_path := strings.ToLower(file_path)
			if file.IsDir() {
				fileIsStatic(fileName, file_path, result)
			} else if strings.Contains(low_file_path, ".xlsx") {
				if strings.Contains(low_file_path, "静态") && strings.Contains(low_file_path, "#"+fileName) {
					*result = true
					return
				}
			}
		}
	}
}

func SaveDB2Execl() {
	allTable := baseuts.FindAllTable(confige.MysqlAccount, confige.MysqlPwd, confige.MysqlHost, confige.MysqlPort, confige.MysqlDBName)
	os.RemoveAll("execluts/execlout")
	os.MkdirAll("execluts/execlout/动态配置表", fs.ModePerm)
	os.MkdirAll("execluts/execlout/静态配置表", fs.ModePerm)
	rowNum := 1
	for k, desc := range allTable {
		fields := baseuts.FindDBTableField(confige.MysqlAccount, confige.MysqlPwd, confige.MysqlHost, confige.MysqlPort, confige.MysqlDBName, k)
		f := excelize.NewFile()
		index, _ := f.NewSheet("Sheet1")
		var cellWidth uint8 = 18
		f.SetSheetProps("Sheet1", &excelize.SheetPropsOptions{BaseColWidth: &cellWidth})

		rowStyle := excelize.Style{}
		rowStyle.Fill.Pattern = 1
		rowStyle.Fill.Type = "pattern"
		rowStyle.Border = []excelize.Border{
			{
				Type:  "right",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "left",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			}}

		rowStyle.Font = &excelize.Font{
			Bold: true,
			// Italic: false,
			// Underline: "single",
			Size:   12,
			Family: "微软雅黑",
			// Strike:    true, // 删除线
			// Color: "#0000FF",
		}
		rowStyle.Fill.Color = []string{"#FFFF00"}
		rowStyleID1, _ := f.NewStyle(&rowStyle)

		rowStyle.Font = &excelize.Font{
			Size:   12,
			Family: "微软雅黑",
		}
		rowStyle.Fill.Color = []string{"#FFFF00"}
		rowStyleID2, _ := f.NewStyle(&rowStyle)

		rowStyle.Font = &excelize.Font{
			Bold:   true,
			Size:   12,
			Family: "微软雅黑",
		}
		rowStyle.Fill.Color = []string{"#AEAAAA"}
		rowStyleID3, _ := f.NewStyle(&rowStyle)

		rowStyle.Font = &excelize.Font{
			Size:   12,
			Family: "微软雅黑",
		}
		rowStyle.Fill.Color = []string{"#75BD42"}
		rowStyleID4, _ := f.NewStyle(&rowStyle)

		rowStyle.Font = &excelize.Font{
			Size:   12,
			Family: "微软雅黑",
		}
		rowStyle.Fill.Color = nil
		rowStyleID5, _ := f.NewStyle(&rowStyle)

		tabIsStatic := false
		fileIsStatic(k+".xlsx", "execluts/execlin", &tabIsStatic)
		num := 'A'
		for _, _v := range fields {
			fn := strings.ReplaceAll(strconv.QuoteRune(num), "'", "")
			colStyleStr := fiedTypeDB2Execl(_v.DataType)

			f.SetCellValue("Sheet1", fn+"1", _v.FieldName)
			f.SetCellStyle("Sheet1", fn+"1", fn+"1", rowStyleID1)

			f.SetCellValue("Sheet1", fn+"2", colStyleStr)
			f.SetCellStyle("Sheet1", fn+"2", fn+"2", rowStyleID2)

			f.SetCellValue("Sheet1", fn+"3", _v.FieldDesc)
			f.SetCellStyle("Sheet1", fn+"3", fn+"3", rowStyleID3)

			if !tabIsStatic {
				fn := strings.ReplaceAll(strconv.QuoteRune(num), "'", "")
				f.SetCellValue("Sheet1", fn+"4", baseuts.GetFieldsValue(_v, true))
				f.SetCellStyle("Sheet1", fn+"4", fn+"4", rowStyleID4)
			}
			num++
		}

		rowNum++

		allTableData := baseuts.FindDBTableData(confige.MysqlAccount, confige.MysqlPwd, confige.MysqlHost, confige.MysqlPort, confige.MysqlDBName, k)
		num = 'A'
		for _, _v := range fields {
			fn := strings.ReplaceAll(strconv.QuoteRune(num), "'", "")
			dataRowNum := 5
			if tabIsStatic {
				dataRowNum = 4
			}
			for _, mv := range allTableData {
				f.SetCellValue("Sheet1", fn+strconv.Itoa(dataRowNum), getCellValue(mv, _v))
				f.SetCellStyle("Sheet1", fn+strconv.Itoa(dataRowNum), fn+strconv.Itoa(dataRowNum), rowStyleID5)
				dataRowNum++
			}
			num++
		}

		f.SetActiveSheet(index)
		folderName := "动态配置表"
		if tabIsStatic {
			folderName = "静态配置表"
		}
		if err := f.SaveAs("execluts/execlout/" + folderName + "/" + desc + "#" + k + ".xlsx"); err != nil {
			println(err.Error())
		}
	}
}

func getCellValue(values map[string]string, fieldST baseuts.Field) interface{} {
	value := values[fieldST.FieldName]
	if fieldST.DataType == "int" {
		_v, err := strconv.ParseInt(value, 10, 64)
		if baseuts.ChkErrNormal(err) {
			return 0
		} else {
			return _v
		}
	} else if fieldST.DataType == "tinyint" {
		_v, err := strconv.ParseInt(value, 10, 64)
		if baseuts.ChkErrNormal(err) {
			return 0
		} else {
			return _v
		}
	} else if fieldST.DataType == "double" {
		_v, err := strconv.ParseFloat(value, 64)
		if baseuts.ChkErrNormal(err) {
			return 0
		} else {
			return _v
		}
	} else if fieldST.DataType == "varchar" {
		return value
	} else if fieldST.DataType == "bool" {
		_v, err := strconv.ParseBool(value)
		if baseuts.ChkErrNormal(err) {
			return false
		} else {
			return _v
		}
	}
	return ""
}

var rootPath string

func ImportExecl2DB() {
	baseuts.BackupMysql("root", "root2023", "xhhy")
	pwd, err := os.Getwd()
	if !baseuts.ChkErr(err) {
		rootPath = pwd
		exec2DB(rootPath + "/execluts/execlin")
		baseuts.LogF("import execl to mysql database success!")
	} else {
		baseuts.LogF("import execl to mysql database fail!")
	}
}

func CreateExecl2Proto() {
	pwd, err := os.Getwd()
	if !baseuts.ChkErr(err) {
		rootPath = pwd
		protoStr = `syntax = "proto3";
		
option go_package = "./;pbstruct";
option csharp_namespace = "PBStruct";` + "\n\n"
		proroDBMapStr = ""
		proroDBMapCount = 1

		os.Remove("protobuff/proto/client/tables.proto")
		execl2Proto(rootPath + "/execluts/execlin")

		if proroDBMapStr != "" {
			protoStr += "message CSStaticTab {}\n\n"
			protoStr += "message SCStaticTab {\n"
			protoStr += proroDBMapStr
			protoStr += "}\n"
		}

		if protoStr != "" {
			baseuts.SaveFile("protobuff/proto/client/tables.proto", protoStr)
		}
		baseuts.LogF("create proto from execl success!")
	} else {
		baseuts.LogF("create proto from execl fail!")
	}
}

var sqlProcess = ""

func exec2DB(_path string) {
	fileInfoList, err := os.ReadDir(_path)
	if !baseuts.ChkErr(err, _path) {
	continueLabel:
		for i := range fileInfoList {
			file := fileInfoList[i]
			file_path := _path + "/" + file.Name()
			if file.IsDir() {
				exec2DB(file_path)
			} else if path.Ext(file_path) == ".xlsx" {

				editSQL := ""
				addSQL := ""
				delSQL := ""
				createSQL := ""
				insertSQL := ""

				var field []string
				var fieldType []string
				var descStr []string
				var defaultValue []string

				var tableName = ""
				var tableDesc = ""
				tableInfo := strings.Split(file.Name(), "#")
				if len(tableInfo) < 2 {
					baseuts.Log("缺少少表名", tableInfo)
					continue
				}
				tableDesc = tableInfo[0]
				tableName = tableInfo[1]
				tableName = strings.ReplaceAll(tableName, ".xlsx", "")
				tableName = strings.ToLower(tableName)

				isStatic := strings.Contains(file_path, "静态")

				ef, err = excelize.OpenFile(file_path)
				if !baseuts.ChkErr(err, _path) {
					rows, err := ef.GetRows("sheet1")
					if !baseuts.ChkErr(err, _path) {
						limitRowNum := 4
						if isStatic {
							limitRowNum = 3
						}
						for rowNum, row := range rows {
							if rowNum < limitRowNum {
								var currArray *[]string = nil
								if rowNum == 0 {
									currArray = &field
								} else if rowNum == 1 {
									currArray = &fieldType
								} else if rowNum == 2 {
									currArray = &descStr
								} else if rowNum == 3 {
									currArray = &defaultValue
								}
								*currArray = append(*currArray, row...)
							} else {
								insertSQL += "("
								for colNum, _ := range field {
									if colNum < len(row) {
										cell := row[colNum]
										if fieldType[colNum] == "string" {
											insertSQL += "'" + cell + "',"
										} else {
											if cell == "" {
												cell = "NULL"
											}
											insertSQL += cell + ","
										}
									} else {
										insertSQL += "NULL,"
									}
								}
								insertSQL = insertSQL[:len(insertSQL)-1] + "),"
							}
						}

						tableFieldData := baseuts.FindDBTableField(confige.MysqlAccount, confige.MysqlPwd, confige.MysqlHost, confige.MysqlPort, confige.MysqlDBName, tableName)
						if CheckMysqlKeyword(tableName) {
							baseuts.Log("表名 " + tableName + " 名称不规范")
							continue
						}
						prevFieldName := ""
						for num, _ := range field {
							_addSQL := ""
							_editSQL := ""
							fValue := strings.ReplaceAll(strings.ToLower(field[num]), " ", "")
							if CheckMysqlKeyword(fValue) {
								baseuts.Log("表 " + tableName + " 中字段 " + fValue + " 名称不规范")
								continue continueLabel
							}
							if fValue == "id" {
								createSQL += "`id` int NOT NULL AUTO_INCREMENT,"
							} else {
								fValueType := strings.ToLower(fieldType[num])
								if len(defaultValue) <= num || defaultValue[num] == "" {
									_editSQL += " modify `" + fValue + "` " + fiedTypeExecl2DB(fValueType) + " DEFAULT NULL"
									_addSQL += " add `" + fValue + "` " + fiedTypeExecl2DB(fValueType) + " DEFAULT NULL"
									createSQL += "`" + fValue + "` " + fiedTypeExecl2DB(fValueType) + " DEFAULT NULL"
								} else {
									_editSQL += " modify `" + fValue + "` " + fiedTypeExecl2DB(fValueType) + " DEFAULT '" + defaultValue[num] + "'"
									_addSQL += "add `" + fValue + "` " + fiedTypeExecl2DB(fValueType) + " DEFAULT '" + defaultValue[num] + "'"
									createSQL += "`" + fValue + "` " + fiedTypeExecl2DB(fValueType) + " DEFAULT '" + defaultValue[num] + "'"
								}

								if len(descStr) <= num || descStr[num] == "" {
									_editSQL += ","
									_addSQL += ","
									createSQL += ","
								} else {
									_editSQL += " COMMENT '" + descStr[num] + "',"
									_addSQL += " COMMENT '" + descStr[num] + "'"
									createSQL += " COMMENT '" + descStr[num] + "',"
								}

								if prevFieldName == "" {
									_addSQL += ","
								} else {
									_addSQL += " AFTER " + prevFieldName + ","
								}

								if hasKey(tableFieldData, fValue) {
									_addSQL = ""
								} else {
									_editSQL = ""
								}

								editSQL += _editSQL
								addSQL += _addSQL
								prevFieldName = fValue
							}
						}
						if addSQL != "" {
							addSQL = addSQL[:len(addSQL)-1]
						}

						if editSQL != "" {
							editSQL = editSQL[:len(editSQL)-1]
						}

						createSQL += "PRIMARY KEY (`id`)"

						for _, _field := range tableFieldData {
							_delSQL := ""
							fValue := strings.ToLower(_field.FieldName)
							if fValue == "id" {
								//id不能删
							} else {
								_delSQL += "drop " + fValue + ","

								for i := range field {
									if strings.ToLower(field[i]) == fValue {
										_delSQL = ""
										break
									}
								}
								delSQL += _delSQL
							}
						}
						if len(delSQL) > 0 {
							delSQL = delSQL[:len(delSQL)-1]
						}

						sqlProcess = ""
						if isStatic {
							sqlProcess += "DROP TABLE IF EXISTS `" + tableName + "`;"
							sqlProcess += "CREATE TABLE `" + tableName + "`("
							sqlProcess += createSQL
							sqlProcess += ")ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci comment='" + tableDesc + "#static';"
						} else {
							if len(tableFieldData) < 1 {
								sqlProcess += "CREATE TABLE IF NOT EXISTS `" + tableName + "`("
								sqlProcess += createSQL
								sqlProcess += ")ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci comment='" + tableDesc + "';"
							} else {
								if addSQL != "" || editSQL != "" {
									sqlProcess += "ALTER TABLE " + tableName
									if addSQL == "" {
										sqlProcess += " " + editSQL + ";"
									} else if editSQL == "" {
										sqlProcess += addSQL + ";"
									} else {
										sqlProcess += " " + addSQL + "," + editSQL + ";"
									}
								}
							}
						}
						if len(insertSQL) > 1 && isStatic {
							sqlProcess += "LOCK TABLES `" + tableName + "` WRITE;"
							sqlProcess += "INSERT INTO `" + tableName + "` VALUES " + insertSQL[:len(insertSQL)-1] + ";"
							sqlProcess += "UNLOCK TABLES;"
						}
						// alter table tablename drop column1,drop column2
						if len(delSQL) > 0 {
							sqlProcess += "alter table " + tableName + " "
							sqlProcess += delSQL + ";"
						}
						// baseuts.Log(sqlProcess)
						if sqlProcess != "" {
							baseuts.ExecMySQLCommand(confige.MysqlAccount, confige.MysqlPwd, confige.MysqlHost, confige.MysqlPort, confige.MysqlDBName, sqlProcess)
						}
					}
				}
			}
		}
	}
}

func hasKey(dat []baseuts.Field, key string) bool {
	for _, v := range dat {
		if strings.EqualFold(v.FieldName, key) {
			return true
		}
	}
	return false
}

var protoStr = ""
var proroDBMapStr = ""
var proroDBMapCount = 1

func execl2Proto(_path string) {
	fileInfoList, err := os.ReadDir(_path)
	if !baseuts.ChkErr(err, _path) {
		for i := range fileInfoList {
			file := fileInfoList[i]
			file_path := _path + "/" + file.Name()
			if file.IsDir() {
				execl2Proto(file_path)
			} else if path.Ext(file_path) == ".xlsx" {
				var field []string
				var fieldType []string

				var tableName = ""
				var fieidsStr = ""

				ef, err = excelize.OpenFile(file_path)
				if !baseuts.ChkErr(err, _path) {
					rows, err := ef.GetRows("sheet1")
					if !baseuts.ChkErr(err, _path) {
						for rowNum, row := range rows {
							if rowNum < 2 {
								var currArray *[]string
								if rowNum == 0 {
									currArray = &field
								} else if rowNum == 1 {
									currArray = &fieldType
								}
								*currArray = append(*currArray, row...)
							}
						}
						tableNameInfo := strings.Split(file.Name(), "#")
						if len(tableNameInfo) < 2 {
							baseuts.Log("缺少表名称", tableNameInfo)
							continue
						}
						tableName = tableNameInfo[1]
						tableName = strings.ReplaceAll(tableName, ".xlsx", "")
						tableName = strings.ToUpper(tableName[:1]) + strings.ToLower(tableName[1:])

						for num, value := range field {
							// if num < len(field)-1 {
							fValue := strings.ToLower(value)
							fieidsStr += "  " + fiedTypeExecl2Proto(fieldType[num]) + " " + fValue + " = " + strconv.Itoa(num+1) + ";\n"
							// }
						}
					}
				}
				if fieidsStr != "" {
					protoStr += "message " + tableName + " {\n"
					protoStr += fieidsStr
					protoStr += "}\n\n"

					protoStr += "message " + tableName + "_result {\n"
					protoStr += "	repeated " + tableName + " data = 1;\n"
					protoStr += "}\n\n"
				}
				if strings.Contains(file_path, "静态") {
					proroDBMapStr += "  repeated " + tableName + " " + strings.ToLower(tableName) + " = " + strconv.Itoa(proroDBMapCount) + ";\n"
				}
				proroDBMapCount++
			}
		}
	}
}

func fiedTypeExecl2Proto(execlType string) string {
	switch execlType {
	case "int":
		return "int32"
	case "int64":
		return "int64"
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "double":
		return "double"
	}
	return ""
}

func fiedTypeExecl2DB(execlType string) string {
	switch execlType {
	case "int":
		return "int"
	case "int64":
		return "bigint(64)"
	case "double":
		return "double"
	case "string":
		return "varchar(255)"
	case "bool":
		return "tinyint"
	}
	return ""
}

func fiedTypeDB2Execl(dbDataType string) string {
	switch dbDataType {
	case "int":
		return "int"
	case "bigint":
		return "int64"
	case "varchar":
		return "string"
	case "tinyint":
		return "bool"
	}
	return ""
}

func ReadFirstRow(idx int) error {
	rows, err := ef.GetRows("Sheet1") // 所有行
	if err != nil {
		return err
	}
	row := rows[1]

	tp := reflect.TypeOf(holder).Elem().Elem().Elem() // 结构体的类型
	val := reflect.New(tp)                            // 创建一个新的结构体对象

	field := val.Elem().Field(idx) // 第idx个字段的反射Value
	cellValue := row[idx]          // 第idx个字段对应的Excel数据
	field.SetString(cellValue)     // 将Excel数据保存到结构体对象的对应字段中

	listV := reflect.ValueOf(holder)
	listV.Elem().Set(reflect.Append(listV.Elem(), val)) // 将结构体对象添加到holder中

	return nil
}

var fileMap = map[string]string{}

func GetJSONFile(fileName string) string {
	var result = fileMap[fileName]
	if result == "" {
		filePtr, err := os.ReadFile("./" + fileName + ".json")
		if !baseuts.ChkErr(err, "文件打开失败") {
			result = string(filePtr)
			fileMap[fileName] = result
		}
	}
	return result
}
