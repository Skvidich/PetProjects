package StatusRelay

func (rel *StatusRelay) SetLogState(st bool) {
	rel.muLog.Lock()
	defer rel.muLog.Unlock()
	rel.isLog = st
}

func (rel *StatusRelay) GetLogState() bool {
	rel.muLog.Lock()
	defer rel.muLog.Unlock()
	return rel.isLog
}

func (rel *StatusRelay) SetResendState(st bool) {
	rel.muResend.Lock()
	defer rel.muResend.Unlock()
	rel.isResend = st
}

func (rel *StatusRelay) GetResendState() bool {
	rel.muResend.Lock()
	defer rel.muResend.Unlock()
	return rel.isResend
}

func (rel *StatusRelay) Run() {
	go rel.inputProcess()
}
