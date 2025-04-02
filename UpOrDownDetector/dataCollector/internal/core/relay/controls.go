package relay

func (r *Relay) Run() {
	go r.process()
}

func (r *Relay) Close() {
	<-r.done
}
