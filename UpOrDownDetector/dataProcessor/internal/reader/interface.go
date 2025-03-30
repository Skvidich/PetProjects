package reader

import "dataProcessor/common"

type Reader interface {
	Next() (*common.ServiceStatus, error)
}
