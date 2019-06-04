package tcpconn

import (
	"time"
)

type Ping struct {
	syn       int32
	serverSeq int32
	clientSeq int32
	// 时延 毫秒
	rtt      uint64
	req1Time uint64
	req2Time uint64
}

// 当收到ping请求时
func (this *Ping) OnRecv(syn, ack, seq int32) (int32, int32, int32) {
	if this.syn == 0 {
		this.onReq1(syn, ack, seq)
		return this.getRes()
	} else {
		this.onReq2(syn, ack, seq)
		return 0, 0, 0
	}
}

func (this *Ping) onReq1(syn, ack, seq int32) {
	if syn == 0 {
		this.clearStatus()
		return
	}
	this.syn = syn
	this.clientSeq = seq
	this.req1Time = uint64(time.Now().UnixNano()) / 1000000
}

func (this *Ping) onReq2(syn, ack, seq int32) {
	if syn != 0 {
		this.clearStatus()
		return
	}
	if seq != this.clientSeq+1 {
		return
	}
	if ack != this.serverSeq+1 {
		return
	}
	this.syn = syn
	this.clientSeq = seq
	this.req2Time = uint64(time.Now().UnixNano()) / 1000000
	this.rtt = this.req2Time - this.req1Time
}

// 获取syn ack seq
func (this *Ping) getRes() (int32, int32, int32) {
	if this.syn == 0 {
		return 0, 0, 0
	}
	this.serverSeq++
	return this.syn, this.clientSeq + 1, this.serverSeq
}

// 获取syn ack seq
func (this *Ping) getReq() (int32, int32, int32) {
	this.clientSeq++
	return 1, this.clientSeq + 1, this.serverSeq
}

func (this *Ping) RTT() uint64 {
	return this.rtt
}

func (this *Ping) clearStatus() {
	this.clientSeq = 0
	this.syn = 0
	this.req1Time = 0
	this.req2Time = 0
}
