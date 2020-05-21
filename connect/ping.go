package connect

import (
	"time"
)

// Ping 连接的 Ping 信息
// 可以主动设置发送时间，或者使用内部的三次握手实现双端 Ping
type Ping struct {
	syn       int32
	serverSeq int32
	clientSeq int32
	// 时延 毫秒
	rtt      uint64
	req1Time uint64
	req2Time uint64
}

// RecordSend 通过发送/接收判断延迟
func (p *Ping) RecordSend() {
	p.req1Time = uint64(time.Now().UnixNano()) / 1000000
}

// RecordRecv 通过发送/接收判断延迟
func (p *Ping) RecordRecv() {
	p.req2Time = uint64(time.Now().UnixNano()) / 1000000
	p.rtt = p.req2Time - p.req1Time
}

// OnRecv 当收到ping请求时，根据ping状态，可以自动判断该ping请求是否需要pong，
// 计算出需要返回的信息（如果需要pong），计算ping时延
func (p *Ping) OnRecv(syn, ack, seq int32) (int32, int32, int32) {
	if p.syn == 0 {
		p.onReq1(syn, ack, seq)
		return p.getRes()
	}
	p.onReq2(syn, ack, seq)
	return 0, 0, 0
}

func (p *Ping) onReq1(syn, ack, seq int32) {
	if syn == 0 {
		p.clearStatus()
		return
	}
	p.syn = syn
	p.clientSeq = seq
	p.req1Time = uint64(time.Now().UnixNano()) / 1000000
}

func (p *Ping) onReq2(syn, ack, seq int32) {
	if syn != 0 {
		p.clearStatus()
		return
	}
	if seq != p.clientSeq+1 {
		return
	}
	if ack != p.serverSeq+1 {
		return
	}
	p.syn = syn
	p.clientSeq = seq
	p.req2Time = uint64(time.Now().UnixNano()) / 1000000
	p.rtt = p.req2Time - p.req1Time
}

// getRes 获取syn ack seq
func (p *Ping) getRes() (int32, int32, int32) {
	if p.syn == 0 {
		return 0, 0, 0
	}
	p.serverSeq++
	return p.syn, p.clientSeq + 1, p.serverSeq
}

// getReq 获取syn ack seq
func (p *Ping) getReq() (int32, int32, int32) {
	p.clientSeq++
	return 1, p.clientSeq + 1, p.serverSeq
}

// RTT 该ping信息的上一次延迟时间
func (p *Ping) RTT() uint64 {
	return p.rtt
}

func (p *Ping) clearStatus() {
	p.clientSeq = 0
	p.syn = 0
	p.req1Time = 0
	p.req2Time = 0
}
