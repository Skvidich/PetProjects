package StatusRelay

func (rel *StatusRelay) GetLogState() bool {

	return rel.isLog
}

func (rel *StatusRelay) GetResendState() bool {

	return rel.isResend
}

func (rel *StatusRelay) Run() {
	go rel.inputProcess()
}

func (rel *StatusRelay) Close() {
	<-rel.endChan
}
