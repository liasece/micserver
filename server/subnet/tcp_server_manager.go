package subnet

import (
	"github.com/liasece/micserver/log"
	// "bytes"
	"errors"
	"fmt"
	"github.com/liasece/micserver/util"
	"net"
)

// 该函数会阻塞
func (this *SubnetManager) BindMyTCPServer(server IServerHandler) error {
	// serverid := this.moudleConf.Myserverinfo.Serverid
	// if serverid == 0 {
	// 	log.Error("[BindMyTCPServer] 本服务器ID为0 无法绑定本机ServerPort")
	// 	return errors.New("server id is 0")
	// }
	serverip := ""
	serverport := 0
	if serverport > 0 {
		portstr := fmt.Sprintf("%d", serverport)
		return this.BindTCPServer(serverip, portstr, server)
	} else {
		log.Error("[BindMyTCPServer] 本服务器serverport为0 " +
			"无法绑定本机ServerPort")
		return errors.New("server port is 0")
	}
	return nil
}

// 接口对象必须用new的
// 绑定服务器
func (this *SubnetManager) BindTCPServer(serverip string, serverport string,
	server IServerHandler) error {
	serverinfo := serverip + ":" + serverport
	netlisten, err := net.Listen("tcp", serverinfo)
	if err != nil {
		log.Error("[SubNetManager.BindTCPServer] "+
			"服务器绑定失败 IP[%s] Port[%s] Err[%s]",
			serverip,
			serverport, err.Error())
		return err
	}
	// myservertype := this.moudleConf.Myserverinfo.Servertype
	log.Debug("[SubNetManager.BindTCPServer] "+
		"服务器绑定成功 IP[%s] Port[%s]", serverip, serverport)

	this.taskManager.serverhandler = server
	go this.TCPServerListenerProcess(netlisten, server)
	return nil
}

func (this *SubnetManager) TCPServerListenerProcess(listener net.Listener,
	server IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[TCPServerListenerProcess] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	defer listener.Close()
	for true {
		this.mTCPServerListener(listener, server)
	}
}

func (this *SubnetManager) mTCPServerListener(listener net.Listener,
	server IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			log.Error("mTCPServerListener "+
				"Panic: ErrName[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	for true {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("[SubNetManager.TCPServerListenerProcess] "+
				"服务器端口监听异常 Err[%s]",
				err.Error())
			continue
		}
		log.Debug("[SubNetManager.BindTCPServer] "+
			"收到新的TCP连接 Addr[%s]",
			conn.RemoteAddr().String())
		tcptask := this.taskManager.AddTCPTask(conn)
		if tcptask != nil {
			server.OnCreateTCPConnect(tcptask)
			go this.taskManager.handleConnection(tcptask, server)
		}
	}
}
