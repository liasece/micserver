package conf

// TBaseOptionCall use of config options
type TBaseOptionCall struct {
	calls []func(*BaseConfig)
}

// Apply of option caller
func (o *TBaseOptionCall) Apply(bf *BaseConfig) {
	for _, v := range o.calls {
		v(bf)
	}
}

// Append caller
func (o *TBaseOptionCall) Append(value func(*BaseConfig)) *TBaseOptionCall {
	o.calls = append(o.calls, value)
	return o
}

// Version option caller
func (o *TBaseOptionCall) Version(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(Version, value)
	})
	return o
}

// ProcessID option caller
func (o *TBaseOptionCall) ProcessID(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(ProcessID, value)
	})
	return o
}

// LogWholePath option caller
func (o *TBaseOptionCall) LogWholePath(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(LogWholePath, value)
	})
	return o
}

// LogLevel option caller
func (o *TBaseOptionCall) LogLevel(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(LogLevel, value)
	})
	return o
}

// SubnetTCPAddr option caller
func (o *TBaseOptionCall) SubnetTCPAddr(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(SubnetTCPAddr, value)
	})
	return o
}

// SubnetNoChan option caller
func (o *TBaseOptionCall) SubnetNoChan(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(SubnetNoChan, value)
	})
	return o
}

// GateTCPAddr option caller
func (o *TBaseOptionCall) GateTCPAddr(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(GateTCPAddr, value)
	})
	return o
}

// IsDaemon option caller
func (o *TBaseOptionCall) IsDaemon(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(IsDaemon, value)
	})
	return o
}

// MsgThreadNum option caller
func (o *TBaseOptionCall) MsgThreadNum(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(MsgThreadNum, value)
	})
	return o
}

// AsynchronousSyncRocbind option caller
func (o *TBaseOptionCall) AsynchronousSyncRocbind(value string) *TBaseOptionCall {
	o.calls = append(o.calls, func(bf *BaseConfig) {
		bf.set(AsynchronousSyncRocbind, value)
	})
	return o
}
