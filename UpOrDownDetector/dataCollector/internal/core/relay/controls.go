package relay

func (r *Relay) LoggingEnabled() bool {

	return r.logging
}

func (r *Relay) ResendEnabled() bool {

	return r.resend
}

func (r *Relay) Run() {
	go r.process()
}

func (r *Relay) Close() {
	<-r.done
}
