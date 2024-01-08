package protobuff

import (
	"baseutils/baseuts"
	"bufio"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	//go install github.com/golang/protobuf/protoc-gen-go@latest
	//go install github.com/micro/micro/v3/cmd/protoc-gen-micro@latest
	//npm install -g protoc-gen-ts
)

// var copyTo = map[string][]string{"go": {"gameutils"}, "micro": {"gameutils"}}
// var outPaths = []string{"go", "cshap", "micro", "ts", "python"}
var copyTo = map[string][]string{"go": {"gameutils"}}
var outPaths = []string{"go", "cshap", "ts", "python"}
var rootPath string

var protoImportInfo = map[string]map[string][]string{}

var protoCount = 1

var protoIDCshapStr = ""
var protoIDPythonStr = ""
var protoUtsCshapStr = ""
var protoUtsGetIDCshapStr = ""
var protoUtsGetPBCshapStr = ""

var protoIDTSStr = ""
var protoUtsTSStr = ""
var protoUtsCodeStr = ""
var protoUtsCodeImportStr = map[string][]string{}

var protoIDGoStr = ""
var protoTypeGoStr = ""
var protoIDGoServerStr = ""
var protoTypeGoServerStr = ""
var protoUtilsGoStr = ""
var protoUtilsGOMarshalStr = ""
var protoUtilsGOUnMarshalStr = ""
var protoUtilsGoDeCTBuf2CTdat = ""
var protoUtilsGoType2ProtoStr = ""
var protoUtilsGoName2ProtoStr = ""
var protoUtilsGoType2IDStr = ""
var protoUtilsGoCallEvtStr = ""

func CreatePBStart() {
	pwd, err := os.Getwd()
	rootPath = pwd

	protoCount = 1
	protoUtsCshapStr = `using Google.Protobuf;
using PBStruct;
class PBUts
{
	public static PBUts I = new PBUts();
	public int getIDFromPB(IMessage protobuf)
    {` + "\n"
	protoIDCshapStr = `class PBStructID
{` + "\n"
	protoIDPythonStr = "class PBStructID:\n"
	protoIDTSStr = "export class PBID {\n"
	protoIDGoStr = "package pbstruct\n\n"
	protoTypeGoStr = "type ProtoType interface {\n	"
	protoIDGoServerStr = "package pbstruct\n\n"
	protoTypeGoServerStr = "type MicroType interface {\n	"
	protoUtilsGoStr = "package pbstruct\n\n" + `import (
	"baseutils/baseuts"
	_ "unsafe"

	"google.golang.org/protobuf/proto"
)` + "\n\n"
	protoUtilsGOMarshalStr = ""
	protoUtilsGOUnMarshalStr = ""
	protoUtilsGoDeCTBuf2CTdat = ""
	protoUtilsGoType2ProtoStr = ""
	protoUtilsGoName2ProtoStr = ""
	protoUtilsGoType2IDStr = ""
	protoUtilsGoCallEvtStr = ""
	protoUtsGetIDCshapStr = ""
	protoUtsGetPBCshapStr = ""

	protoUtsTSStr = ""

	if !baseuts.ChkErr(err) {
		// for i := range outPaths {
		// 	_path := "protobuff/pb" + outPaths[i]
		// 	err = os.RemoveAll(_path)
		// 	baseuts.ChkErr(err)
		// }

		for _, v := range copyTo {
			for i := range v {
				err = os.RemoveAll("../" + v[i] + "/pbstruct")
				baseuts.ChkErr(err)
			}
		}

		for i := range outPaths {
			if outPaths[i] != "python" {
				_path := "protobuff/pb" + outPaths[i]
				if _, err := os.Stat(_path); err == nil || os.IsExist(err) {
					err1 := os.RemoveAll(_path)
					baseuts.ChkErr(err1)
				}
				err2 := os.Mkdir(_path, os.ModePerm)
				baseuts.ChkErr(err2)
			}
		}
		os.Mkdir("protobuff/pbpython/pbstruct", os.ModePerm)

		getAllProtoFile(pwd + "/protobuff/proto")

		protoIDPythonStr += "  ProtoMax = " + strconv.Itoa(protoCount) + "\n"
		baseuts.SaveFile("protobuff/pbpython/pbstruct/proto_id.py", protoIDPythonStr)
		protoIDCshapStr += "}"
		baseuts.SaveFile("protobuff/pbcshap/ProtoID.cs", protoIDCshapStr)
		protoUtsGetIDCshapStr += `    	return 0;
	}` + "\n"
		protoUtsCshapStr += protoUtsGetIDCshapStr

		protoUtsCshapStr += `	public IMessage getPBFromID(int protoID)
	{` + "\n"
		protoUtsGetPBCshapStr += `    	return null;
	}` + "\n"
		protoUtsCshapStr += protoUtsGetPBCshapStr
		protoUtsCshapStr += `}`
		baseuts.SaveFile("protobuff/pbcshap/PBUts.cs", protoUtsCshapStr)
		protoTypeGoStr = protoTypeGoStr[:len(protoTypeGoStr)-2] + `
}` + "\n"
		protoIDGoStr += "const ProtoMax int32 = " + strconv.Itoa(protoCount) + "\n"
		protoIDGoStr += "\n" + protoTypeGoStr
		baseuts.SaveFile("protobuff/pbgo/protoid.go", protoIDGoStr)

		protoIDTSStr += "    static ProtoMax: number = " + strconv.Itoa(protoCount) + "\n}"
		baseuts.SaveFile("protobuff/pbts/protoid.ts", protoIDTSStr)

		for k, v := range protoUtsCodeImportStr {
			protoUtsTSStr += "import {"
			for i := range v {
				protoUtsTSStr += " " + v[i] + ","
			}
			protoUtsTSStr = protoUtsTSStr[:len(protoUtsTSStr)-1]
			protoUtsTSStr += " } from \"./" + strings.ReplaceAll(k, ".proto", "") + "\";\n"
		}
		protoUtsTSStr += `
export class PBUtils {
	private static codePBDat(pbDatID: number, pbDat: object, pbBuf: Uint8Array): object {
		var result
		switch (pbDatID) {` + "\n"
		protoUtsCodeStr += `        }
        return result
    }` + "\n"
		protoUtsTSStr += protoUtsCodeStr
		protoUtsTSStr += "}"
		baseuts.SaveFile("protobuff/pbts/pbutils.ts", protoUtsTSStr)

		protoTypeGoServerStr = protoTypeGoServerStr[:len(protoTypeGoServerStr)-2] + `
}` + "\n"
		protoIDGoServerStr += "\n" + protoTypeGoServerStr
		baseuts.SaveFile("protobuff/pbgo/protoidserver.go", protoIDGoServerStr)

		protoUtilsGoStr += "//go:linkname getCPFromProtoName gameutils/common/pbuts.GetCPFromProtoName\n"
		protoUtilsGoStr += `func getCPFromProtoName(_protoName string) interface{} {
		switch _protoName {` + "\n"
		protoUtilsGoStr += protoUtilsGoName2ProtoStr
		protoUtilsGoStr += `	}
	return nil
}` + "\n\n"
		protoUtilsGoStr += "//go:linkname getCPFromProtoID gameutils/common/pbuts.GetCPFromProtoID\n"
		protoUtilsGoStr += `func getCPFromProtoID(protoType int32, buffes ...[]byte) interface{} {
	switch protoType {` + "\n"
		protoUtilsGoStr += protoUtilsGoType2ProtoStr
		protoUtilsGoStr += `	}
	return nil
}` + "\n\n"
		protoUtilsGoStr += "//go:linkname getProtoIDNameFromCP gameutils/common/pbuts.GetProtoIDNameFromCP\n"
		protoUtilsGoStr += `func getProtoIDNameFromCP(clientPBData interface{}) (int32, string) {
	switch clientPBData.(type) {` + "\n"
		protoUtilsGoStr += protoUtilsGoType2IDStr
		protoUtilsGoStr += `	}
	return -1, ""
}` + "\n\n"
		protoUtilsGoStr += "//go:linkname callProtoEvt gameutils/common/pbuts.CallProtoEvt\n"
		protoUtilsGoStr += `func callProtoEvt(callFunc interface{}, param interface{}, rsp *[]byte, token string, jsonRsp ...*string) {
	switch param := param.(type) {` + "\n"
		protoUtilsGoStr += protoUtilsGoCallEvtStr
		protoUtilsGoStr += `	}
}` + "\n\n"
		protoUtilsGoStr += "//go:linkname protoMarshal gameutils/common/pbuts.ProtoMarshal\n"
		protoUtilsGoStr += `func protoMarshal(pbDat interface{}) ([]byte, int32) {
	var b []byte = nil
	var err error = nil
	var pbDatID int32 = -1
	switch pbDat := pbDat.(type) {` + "\n"
		protoUtilsGoStr += protoUtilsGOMarshalStr
		protoUtilsGoStr += `	}
	baseuts.ChkErr(err)
	return b, pbDatID
}` + "\n\n"
		protoUtilsGoStr += "//go:linkname protoUnMarshal gameutils/common/pbuts.ProtoUnMarshal\n"
		protoUtilsGoStr += `func protoUnMarshal(pbBuf []byte, dstPbDat interface{}) int32 {
	var err error = nil
	var pbDatID int32 = -1
	switch dstPbDat := dstPbDat.(type) {` + "\n"
		protoUtilsGoStr += protoUtilsGOUnMarshalStr
		protoUtilsGoStr += `	}
	baseuts.ChkErr(err)
	return pbDatID
}` + "\n\n"
		protoUtilsGoStr += "//go:linkname deCTBuf2CTdat gameutils/common/pbuts.DeCTBuf2CTdat\n"
		protoUtilsGoStr += `func deCTBuf2CTdat(buf []byte, result ...interface{}) *ClientTrans {
	var clentTrans ClientTrans
	proto.Unmarshal(buf, &clentTrans)
	if len(result) > 0 {
		var ptr = result[0]
		switch ptr := ptr.(type) {` + "\n"
		protoUtilsGoStr += protoUtilsGoDeCTBuf2CTdat
		protoUtilsGoStr += `		}
	}
	return &clentTrans
}` + "\n\n"
		baseuts.SaveFile("protobuff/pbgo/protoutils.go", protoUtilsGoStr)
		copyGODir()
		copyTSDir("protobuff/pbts/protobuff")
		copyPythonDir("protobuff/pbpython/protobuff")

		baseuts.CopyFile("protobuff/pbpython/pbstruct/protobuff", "protobuff/pbpython/protobuff")
		os.RemoveAll("protobuff/pbpython/protobuff")

		os.RemoveAll("protobuff/pbts/protobuff")
		os.RemoveAll("protobuff/pbgo")
		// os.RemoveAll("protobuff/pbmicro")
		baseuts.Log("proto buff created!!")
	}
}

func copyGODir() {
	for k, v := range copyTo {
		for i := range v {
			baseuts.CopyFile("../"+v[i]+"/pbstruct", "protobuff/pb"+k, func(f *os.File) string {
				allStr := ""
				reader := bufio.NewReader(f) // 读取文本数据
				for {
					str, err := reader.ReadString('\n')
					if err == io.EOF {
						break
					}
					str = strings.Replace(str, "github.com/micro/micro/v3/service/api", "github.com/asim/go-micro/v3/api", -1)
					str = strings.Replace(str, "github.com/micro/micro/v3/service/client", "github.com/asim/go-micro/v3/client", -1)
					str = strings.Replace(str, "github.com/micro/micro/v3/service/server", "github.com/asim/go-micro/v3/server", -1)
					allStr += str
				}
				return allStr
			})
		}
	}
}
func copyTSDir(src string) {
	fileInfoList, err := os.ReadDir(src)
	if err != nil {

	} else {
		for i := range fileInfoList {
			file := fileInfoList[i]
			file_path := src + "/" + file.Name()
			if file.IsDir() {
				copyTSDir(file_path)
			} else if path.Ext(file_path) == ".ts" {
				baseuts.CopyFile("protobuff/pbts/", file_path)
			}
		}
	}
	baseuts.CopyFile("protobuff/pbts/", "protobuff/proto/protobase.ts")
	checkImportProto()
}
func copyPythonDir(src string) {
	fileInfoList, err := os.ReadDir(src)
	if err != nil {

	} else {
		if len(fileInfoList) > 0 {
			baseuts.SaveFile(src+"/__init__.py", "")
		}
		for i := range fileInfoList {
			file := fileInfoList[i]
			file_path := src + "/" + file.Name()
			if file.IsDir() {
				copyPythonDir(file_path)
			}
		}
	}
}
func checkImportProto() {
	for k, v := range protoImportInfo {
		var pathBase = path.Base(k)
		tsPathRoot := rootPath + "/protobuff/pbts/" + strings.ReplaceAll(pathBase, ".proto", ".ts")
		importStr := "import { Long, writeBytes, readVarint64, readBytes, writeByte, readByte, popByteBuffer, toUint8Array, ByteBuffer, writeVarint32, writeVarint64, intToLong, writeString, wrapByteBuffer, isAtEnd, readVarint32, readString, skipUnknownField, writeDouble, readDouble, writeByteBuffer, pushByteBuffer, pushTemporaryLength } from \"./protobase\"\n"
		if _, err := os.Stat(tsPathRoot); (err == nil || os.IsExist(err)) && len(v) > 0 {
			importStr += "import { "
			for k1, v1 := range v {
				tsPath := rootPath + "/protobuff/pbts/" + strings.ReplaceAll(path.Base(k1), ".proto", ".ts")
				tsContentStr := ""
				f, _ := os.Open(tsPath)
				reader := bufio.NewReader(f) // 读取文本数据
				for {
					str, err := reader.ReadString('\n')
					for _, v2 := range v1 {
						funcNameDe := "function _decode" + v2 + "("
						funcNameEn := "function _encode" + v2 + "("
						if strings.Contains(str, funcNameDe) || strings.Contains(str, funcNameEn) {
							str = "export " + str
							if !strings.Contains(importStr, "_decode"+v2) {
								importStr += v2 + ", _decode" + v2 + ", _encode" + v2 + ","
							}
							break
						}
					}
					tsContentStr += str
					if err == io.EOF {
						break
					}
				}
				baseuts.SaveFile(tsPath, tsContentStr)
			}
			importStr = importStr[:len(importStr)-1] + " } from \"./tables\";\n"
		}
		// if strings.Contains(importStr, ",") {
		f, err := os.Open(tsPathRoot)
		if !baseuts.ChkErrNormal(err) {
			reader := bufio.NewReader(f) // 读取文本数据
			for {
				str, err := reader.ReadString('\n')
				if strings.Contains(str, "export interface Long {") {
					break
				}
				importStr += str
				if err == io.EOF {
					break
				}
			}
			baseuts.SaveFile(tsPathRoot, importStr)
		}
		// }
	}
	baseuts.Log(protoImportInfo)
}

func getAllProtoFile(folder_path string) {
	fileInfoList, err := os.ReadDir(folder_path)
	if !baseuts.ChkErr(err, folder_path) {
		for i := range fileInfoList {
			file := fileInfoList[i]
			file_path := folder_path + "/" + file.Name()
			if file.IsDir() {
				getAllProtoFile(file_path)
			} else if path.Ext(file_path) == ".proto" {
				// if strings.LastIndex(file_path, "/micro/") >= 0 {
				// 	createMicroPB(file_path)
				// } else
				if strings.LastIndex(file_path, "/client/") >= 0 {
					createTSPB(file_path)
					createCshapPB(file_path)
					createGoPB(file_path)
					createPythonPB(file_path)
					createProtoID(file_path)
				} else if strings.LastIndex(file_path, "/server/") >= 0 {
					createGoPB(file_path)
					createProtoID(file_path)
					createPythonPB(file_path)
				} else {
					createGoPB(file_path)
					createPythonPB(file_path)
				}
			}
		}
	}
}

func createProtoID(proto_path string) {
	var protoImport = map[string][]string{}
	f, _ := os.Open(proto_path)
	rootFileName := path.Base(proto_path)
	reader := bufio.NewReader(f) // 读取文本数据
	for {
		str, err := reader.ReadString('\n')

		regFindImport := regexp.MustCompile(".*import \"(.*?)\".*")
		matchArrayImportPath := regFindImport.FindStringSubmatch(str)

		if len(matchArrayImportPath) > 1 {
			var importPath = matchArrayImportPath[1]
			if _, ok := protoImport[importPath]; !ok {
				protoImport[importPath] = []string{}
			}
		} else {
			regFindProp := regexp.MustCompile(".* (.*) .* = .*")
			matchArrayProp := regFindProp.FindStringSubmatch(str)
			if len(matchArrayProp) > 1 {
				propNamr := matchArrayProp[1]
				for key, _ := range protoImport {
					if checkProtoHasStruct(key, propNamr) {
						protoImport[key] = append(protoImport[key], propNamr)
					}
				}
			}
			// else {
			regFind := regexp.MustCompile(".*message (.*?) .*")
			matchArray := regFind.FindStringSubmatch(str)
			if len(matchArray) > 1 {
				_structName := matchArray[1]
				protoUtilsGoName2ProtoStr += `	case "` + _structName + `", "` + strings.ToLower(_structName) + `":
		result := ` + baseuts.UnderScoreCase2CamelCase(_structName) + `{}
		return &result` + "\n"

				runElse := true
				_structID := strconv.Itoa(protoCount)
				if strings.Contains(matchArray[1], "SC") || strings.Contains(matchArray[1], "CS") || strings.Contains(matchArray[1], "Micro") || rootFileName == "tables.proto" {
					if strings.Contains(matchArray[1], "Micro") {
						runElse = false

						protoIDPythonStr += "  " + _structName + "_ID = " + _structID + "\n"
						protoIDGoServerStr += "const " + _structName + "_ID int32 = " + _structID + "\n"
						protoTypeGoServerStr += _structName + " |"
					} else if rootFileName == "tables.proto" {
						runElse = false

						if strings.Contains(matchArray[1], "SC") || strings.Contains(matchArray[1], "CS") {
							runElse = true
						}

						_structNameChange := baseuts.UnderScoreCase2CamelCase(_structName)

						if !runElse {
							protoIDGoServerStr += "const " + _structNameChange + "_ID int32 = " + _structID + "\n"
						}

						protoUtilsGOMarshalStr += `	case *` + _structNameChange + `:
		b, err = proto.Marshal(pbDat)
		pbDatID = ` + _structNameChange + "_ID \n"
						protoUtilsGOUnMarshalStr += `	case *` + _structNameChange + `:
		err = proto.Unmarshal(pbBuf, dstPbDat)
		pbDatID = ` + _structNameChange + "_ID \n"
						protoUtilsGoType2ProtoStr += `	case ` + _structID + `:
		result := ` + _structNameChange + `{}
		if len(buffes) > 0 {
			buf := buffes[0]
			err := proto.Unmarshal(buf, &result)
			if err != nil {
				baseuts.LogF(err.Error())
			}
		}
		return &result` + "\n"
						protoUtilsGoDeCTBuf2CTdat += `		case *` + _structNameChange + `:
			err := proto.Unmarshal(clentTrans.Protobuff, ptr)
			if err != nil {
				baseuts.LogF(err.Error())
			}` + "\n"

						protoUtilsGoType2IDStr += `	case *` + _structNameChange + `, []*` + _structNameChange + `, *[]*` + _structNameChange + `:
		return ` + _structID + `, "` + _structNameChange + `"` + "\n"
						protoUtilsGoCallEvtStr += `	case ` + _structNameChange + `:
		callFunc.(func(*` + _structNameChange + `, *[]byte, string, ...*string))(&param, rsp, token, jsonRsp...)
	case *` + _structNameChange + `:
		callFunc.(func(*` + _structNameChange + `, *[]byte, string, ...*string))(param, rsp, token, jsonRsp...)` + "\n"

					}
					if runElse {
						if _, ok := protoUtsCodeImportStr[rootFileName]; !ok {
							protoUtsCodeImportStr[rootFileName] = []string{}
						}
						protoUtsCodeImportStr[rootFileName] = append(protoUtsCodeImportStr[rootFileName], "encode"+_structName, "decode"+_structName)
						protoUtsCodeStr += `            case ` + _structID + `:
				if (pbDat != null) {
					result = encode` + _structName + `(pbDat)
				} else if (pbBuf != null) {
					result = decode` + _structName + `(pbBuf)
				}
				break` + "\n"

						protoUtsGetIDCshapStr += `		if (protobuf is ` + _structName + `)
		{
			return ` + _structID + `;
		}` + "\n"
						protoUtsGetPBCshapStr += `		if (protoID == ` + _structID + `)
		{
			return new ` + _structName + `();
		}` + "\n"
						protoIDCshapStr += "    public const int " + _structName + "_ID = " + _structID + ";\n"
						protoIDPythonStr += "  " + _structName + "_ID = " + _structID + "\n"
						protoIDGoStr += "const " + _structName + "_ID int32 = " + _structID + "\n"
						protoIDTSStr += "    static " + _structName + "_ID: number = " + _structID + "\n"

						protoTypeGoStr += _structName + " |"
					}
					if rootFileName != "tables.proto" {
						protoUtilsGOMarshalStr += `	case *` + _structName + `:
		b, err = proto.Marshal(pbDat)
		pbDatID = ` + _structName + "_ID \n"
						protoUtilsGOUnMarshalStr += `	case *` + _structName + `:
		err = proto.Unmarshal(pbBuf, dstPbDat)
		pbDatID = ` + _structName + "_ID \n"
						protoUtilsGoDeCTBuf2CTdat += `		case *` + _structName + `:
			err := proto.Unmarshal(clentTrans.Protobuff, ptr)
			if err != nil {
				baseuts.LogF(err.Error())
			}` + "\n"

						protoUtilsGoType2ProtoStr += `	case ` + _structID + `:
		result := ` + _structName + `{}
		if len(buffes) > 0 {
			buf := buffes[0]
			err := proto.Unmarshal(buf, &result)
			if err != nil {
				baseuts.LogF(err.Error())
			}
		}
		return &result` + "\n"
						protoUtilsGoType2IDStr += `	case *` + _structName + `, []*` + _structName + `, *[]*` + _structName + `:
		return ` + _structID + `, "` + _structName + `"` + "\n"
						protoUtilsGoCallEvtStr += `	case ` + _structName + `:
		callFunc.(func(*` + _structName + `, *[]byte, string, ...*string))(&param, rsp, token, jsonRsp...)
	case *` + _structName + `:
		callFunc.(func(*` + _structName + `, *[]byte, string, ...*string))(param, rsp, token, jsonRsp...)` + "\n"
					}
					protoCount++
				}
			}
		}

		if err == io.EOF {
			break
		}
	}
	// if len(protoImport) > 0 {
	rel, _ := filepath.Rel(rootPath, proto_path)
	rel = strings.ReplaceAll(rel, "\\", "/")
	protoImportInfo[rel] = protoImport
	// }
}

func checkProtoHasStruct(filePath string, structName string) bool {
	var b = make([]byte, 4096)
	f, _ := os.Open(filePath)
	reader := bufio.NewReader(f) // 读取文本数据
	reader.Read(b)
	allStr := string(b)
	return strings.Contains(allStr, "message "+structName+" ")
}

func createGoPB(proto_path string) {
	if runtime.GOOS == "windows" {
		runProtoC(proto_path, "--go_out=protobuff/pbgo/", "--plugin=protoc-gen-go=../../bin/protoc-gen-go.exe")
	} else {
		runProtoC(proto_path, "--go_out=protobuff/pbgo/", "--plugin=protoc-gen-go=../../bin/protoc-gen-go", "--proto_path=./")
	}
}

func createCshapPB(proto_path string) {
	runProtoC(proto_path, "--csharp_out=protobuff/pbcshap/")
}

func createPythonPB(proto_path string) {
	runProtoC(proto_path, "--python_out=protobuff/pbpython/")
}

func createTSPB(proto_path string) {
	var fileName = path.Base(proto_path)
	var fileNameSuffix = path.Ext(proto_path)
	runCommand("pbjs", proto_path, "--ts", "protobuff/pbts/"+fileName[:len(fileName)-len(fileNameSuffix)]+".ts")
}

func createMicroPB(proto_path string) {
	if runtime.GOOS == "windows" {
		// runProtoC(proto_path, "--micro_out=protobuff/pbmicro/", "--plugin=protoc-gen-micro=../../bin/protoc-gen-micro.exe")
		runProtoC(proto_path, "--go_out=protobuff/pbmicro/", "--plugin=protoc-gen-go=../../bin/protoc-gen-go.exe")
	} else {
		// runProtoC(proto_path, "--micro_out=protobuff/pbmicro/", "--plugin=protoc-gen-micro=../../bin/protoc-gen-micro")
		runProtoC(proto_path, "--go_out=protobuff/pbmicro/", "--plugin=protoc-gen-go=../../bin/protoc-gen-go")
	}
}

func runProtoC(input string, output string, other_args ...string) {
	input = strings.Replace(input, rootPath+"/", "", -1)
	other_args = append(other_args, output)
	other_args = append(other_args, input)

	var out []byte
	var err error
	if runtime.GOOS == "windows" {
		out, err = runCommand("executable/protoc-3.5.1-win32/bin/protoc", other_args...)
	} else {
		out, err = runCommand("protoc", other_args...)
	}
	baseuts.ChkErr(err, input+","+output+","+string(out))
}
func runCommand(exePath string, other_args ...string) ([]byte, error) {
	var cmd *exec.Cmd = exec.Command(exePath, other_args...)
	return cmd.CombinedOutput()
}
