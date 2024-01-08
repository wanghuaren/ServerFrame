package gameuts

import (
	"baseutils/baseuts"
	"bufio"
	"net"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type DiscoveryConfig struct {
	ID      string
	Name    string
	Tags    []string
	Port    int64
	Address string
}

func RegisterService(dis DiscoveryConfig, consulAddress string) error {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	baseuts.ChkErr(err)
	registration := &consulapi.AgentServiceRegistration{
		ID:      dis.ID,
		Name:    dis.Name,
		Port:    int(dis.Port),
		Tags:    dis.Tags,
		Address: dis.Address,
	}
	// 启动tcp的健康检测，注意address不能使用127.0.0.1或者localhost，因为consul-agent在docker容器里，如果用这个的话，
	// consul会访问容器里的port就会出错，一直检查不到实例
	check := &consulapi.AgentServiceCheck{}
	portStr := strconv.FormatInt(dis.Port, 10)
	check.TCP = registration.Address + ":" + portStr
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "60s"
	registration.Check = check

	if err := client.Agent().ServiceRegister(registration); baseuts.ChkErrNormal(err) {
		return err
	}
	go startTcp()
	return nil
}

func startTcp() {
	ls, err := net.Listen("tcp", ":10111")
	if baseuts.ChkErr(err) {
		return
	}
	for {
		conn, err := ls.Accept()
		baseuts.ChkErr(err)
		go func(conn net.Conn) {
			_, err := bufio.NewWriter(conn).WriteString("hello consul")
			baseuts.ChkErr(err)
		}(conn)
	}
}

func Discovery(serviceName string, consulAddress string) []*consulapi.ServiceEntry {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	baseuts.ChkErr(err)
	service, _, err := client.Health().Service(serviceName, "", false, nil)
	baseuts.ChkErr(err)
	if len(service) < 1 {
		baseuts.Log(serviceName + "服务未发现,重试")
		time.Sleep(time.Second)
		return Discovery(serviceName, consulAddress)
	}
	return service
}
