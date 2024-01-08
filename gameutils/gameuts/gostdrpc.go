package gameuts

import (
	"baseutils/baseuts"
	"errors"
	"gameutils/common"
	"gameutils/common/microuts"
	"gameutils/common/pbuts"
	"gameutils/pbstruct"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"strconv"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

type IMicroServer interface {
	RegisiterFunc(string, func(*pbstruct.MicroTrans, *[]byte))
	RegisiterFuncJson(string, func(*pbstruct.MicroTrans, *string))
	StartMicroJson(string)
	Run()
}

type MicroServer struct {
	serverName      string
	server          *rpc.Server
	listener        net.Listener
	recvMapFunc     map[string]func(*pbstruct.MicroTrans, *[]byte)
	recvMapFuncJson map[string]func(*pbstruct.MicroTrans, *string)
	host            string
	port            string
}

func FactoryServer(host string, port string, serverName string) IMicroServer {
	var err error
	_server := MicroServer{}
	_server.host = host
	_server.port = port
	_server.serverName = serverName
	_server.server = rpc.NewServer()
	_server.recvMapFunc = map[string]func(*pbstruct.MicroTrans, *[]byte){}
	_server.recvMapFuncJson = map[string]func(*pbstruct.MicroTrans, *string){}
	_server.server.Register(&_server)
	_server.listener, err = net.Listen("tcp", host+":"+port)
	baseuts.Log(serverName, "Micro创建中", host, port)
	if baseuts.ChkErr(err) {
		baseuts.LogF(serverName, "Micro启动失败,重试", err)
		time.Sleep(time.Second)
		return FactoryServer(host, port, serverName)
	}
	return &_server
}

func (s *MicroServer) StartMicroJson(port string) {
	var err error
	_server := MicroServer{}
	_server.serverName = s.serverName
	_server.server = rpc.NewServer()
	_server.recvMapFuncJson = s.recvMapFuncJson
	_server.server.Register(s)
	_server.listener, err = net.Listen("tcp", s.host+":"+port)
	baseuts.Log(s.serverName, "Micro Json创建中", s.host, port)
	if baseuts.ChkErr(err) {
		baseuts.LogF(s.serverName, "Micro Json创建失败,重试", err)
		time.Sleep(time.Second)
		s.StartMicroJson(port)
	} else {
		baseuts.Log(s.serverName, "Micro Json启动成功")
		for {
			conn, err := _server.listener.Accept()
			if baseuts.ChkErr(err) {
				baseuts.Log(_server.serverName, "Micro Json 关闭")
				break
			}
			_server.server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}

func (s *MicroServer) Micro(req *[]byte, rsp *[]byte) error {
	defer baseuts.ChkRecover()
	mt := pbstruct.MicroTrans{}
	pbuts.ProtoUnMarshal(*req, &mt)
	if f, _ok := s.recvMapFunc[mt.Key]; _ok {
		// _postmsg := microuts.DeMTFromMTBuf(req.Protobuf)
		f(&mt, rsp)
	}
	return nil
}

func (s *MicroServer) MicroJSON(req string, rsp *string) error {
	var mt = pbstruct.MicroTrans{}
	err := protojson.Unmarshal([]byte(req), &mt)
	if !baseuts.ChkErr(err) {
		if f, _ok := s.recvMapFuncJson[mt.Key]; _ok {
			f(&mt, rsp)
		}
	}
	return nil
}

func (s *MicroServer) RegisiterFunc(key string, f func(*pbstruct.MicroTrans, *[]byte)) {
	s.recvMapFunc[key] = f
}

func (s *MicroServer) RegisiterFuncJson(key string, f func(*pbstruct.MicroTrans, *string)) {
	s.recvMapFuncJson[key] = f
}

func (s *MicroServer) Run() {
	// go runJsonRPC(s)
	baseuts.Log(s.serverName, "Micro启动成功")
	s.server.Accept(s.listener)
}

type IMicroClient interface {
	CallFunc(interface{}) *pbstruct.MicroTrans
}
type microClient struct {
	clientName string
	serverName string
	client     *rpc.Client
	conn       *net.Conn
	host       string
	port       string
}

func FactoryClient(host string, port string, serverName string, clientName string) IMicroClient {
	_client := microClient{}
	_client.clientName = clientName
	_client.serverName = serverName
	_client.host = host
	_client.port = port
	_client.client = createRPCHttp(&_client, host, port)
	return &_client
}

func createRPCHttp(c *microClient, host string, port string) *rpc.Client {
	conn, err := net.DialTimeout("tcp", host+":"+port, time.Millisecond*time.Duration(common.MICRO_TIMEOUT_MILLISECOND))
	c.conn = &conn
	// conn, err := rpc.DialHTTP("tcp", host+":"+port)

	if baseuts.ChkErrNormal(err, c.clientName+" 连接 "+c.serverName+" 失败", host, port) {
		return nil
	} else {
		baseuts.Log(c.clientName+" 连接 "+c.serverName+" 成功", host, port)
		_client := rpc.NewClient(conn)
		return _client
		// return conn
	}
}

func (c *microClient) CallFunc(mtTypeDat interface{}) *pbstruct.MicroTrans {
	return microCallFunc(c, false, false, mtTypeDat)
}

func (c *microClient) CallFuncJson(mtTypeDat interface{}) *pbstruct.MicroTrans {
	return microCallFunc(c, false, true, mtTypeDat)
}

func (c *microClient) CallFuncSync(mtTypeDat interface{}) *pbstruct.MicroTrans {
	return microCallFunc(c, true, false, mtTypeDat)
}

func (c *microClient) CallFuncSyncJson(mtTypeDat interface{}) *pbstruct.MicroTrans {
	return microCallFunc(c, true, true, mtTypeDat)
}

func microCallFunc(c *microClient, isSync bool, isJson bool, mtTypeDat interface{}) *pbstruct.MicroTrans {
	var _postmsg = &pbstruct.MicroTrans{}
	if isJson {
		// TODO
		// trans = microuts.GetMTFromMTTypeDatJson(mtTypeDat)
	} else {
		reqPb := microuts.GetMTFromMTTypeDat(mtTypeDat)
		if c.client == nil {
			c.client = createRPCHttp(c, c.host, c.port)
		}
		var errStr = "micro 连接 " + c.host + ":" + c.port + " 错误"
		if isSync {
			errStr = "sync " + errStr
		}
		err := errors.New(errStr)
		var rsp = &[]byte{}
		if c.client != nil {
			reqPb.Key = c.clientName
			req, _ := pbuts.ProtoMarshal(reqPb)
			if c.conn != nil && *c.conn != nil {
				(*c.conn).SetDeadline(time.Now().Add(time.Millisecond * time.Duration(common.MICRO_TIMEOUT_MILLISECOND)))
			}
			if isSync {
				replyDat := c.client.Go("MicroServer.Micro", req, rsp, nil)
				replyRsp := <-replyDat.Done
				baseuts.LogDebug("异步micro call", replyRsp.Reply)
				err = replyRsp.Error
			} else {
				err = c.client.Call("MicroServer.Micro", req, rsp)
			}
		}

		microError := common.MICRO_ERROR
		if baseuts.ChkErrNormal(err, c.clientName, c.serverName, c.host, c.port) {
			if c.conn != nil && *c.conn != nil {
				(*c.conn).Close()
				c.conn = nil
			}
			if c.client != nil {
				c.client.Close()
				c.client = nil
			}
			if err.Error() == "connection is shut down" {
				baseuts.Log("重试", isSync, isJson, mtTypeDat)
				return microCallFunc(c, isSync, isJson, mtTypeDat)
			}
			microError = common.MICRO_ERROR_TIMEOUT
		} else {
			// if rsp == nil {
			// 	microError = common.MICRO_ERROR_PROTOBUF_NIL
			// } else {
			pbuts.ProtoUnMarshal(*rsp, _postmsg)
			if _postmsg.ProtoBuf == nil && _postmsg.ProtoBufStr == "" {
				microError = common.MICRO_ERROR_PROTOBUF_NIL
			} else {
				microError = 0
			}
			// }
		}
		if microError != 0 {
			baseuts.LogF(c.clientName + "->" + c.serverName + " error:" + strconv.Itoa(int(microError)))
			_postmsg = nil
		}
	}
	return _postmsg
}
