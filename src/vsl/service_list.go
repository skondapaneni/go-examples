package vsl

import (
	"collections"
	"fmt"
)

type ServiceList struct {
	alist *collections.ArrayList
}

func NewServiceList() *ServiceList {
	return &ServiceList{
		alist: collections.NewArrayList(),
	}
}

func (sl *ServiceList) Add(service *VSLConfig) {
	sl.alist.Add(service)
}

func (sl *ServiceList) Contains_i(b *VSLConfig) (int, bool) {
	return sl.alist.Contains_i(b)
}

func (sl *ServiceList) RemoveIndex(i int) {
	sl.alist.RemoveIndex(i)
}

func (sl *ServiceList) ReplaceIndex(b *VSLConfig, i int) {
	sl.alist.ReplaceIndex(b, i)
}

func (sl *ServiceList) Format() []byte {
	iter := sl.alist.NewIterator()
	rval := []byte("")
	for iter.HasNext() {
		service := iter.Next().(*VSLConfig)
		line := fmt.Sprintf("%s %s %s %s %s\n",
			service.Interface,
			service.Service,
			service.App,
			service.Role,
			service.Subnet.String())
		lbytes := []byte(line)
		rval = append(rval[:], lbytes[:]...)
	}
	iter.Close()
	return rval
}
