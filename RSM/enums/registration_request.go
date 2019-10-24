package enums

type Registration_Request int

const (
	Registration_Request_start Registration_Request = iota
	Registration_Request_stop
	Registration_Request_partial_stop
	Registration_Request_add
)
