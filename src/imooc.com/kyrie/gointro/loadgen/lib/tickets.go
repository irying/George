package lib

import (
	"fmt"
	"errors"
)

type GoTickets interface {
	Take()
	Return()
	Active() bool
	Total() uint32
	Remainder() uint32
}

type myGoTickets struct {
	total uint32
	ticketCh chan byte
	active bool
}

func NewGoTickets(total uint32) (GoTickets, error) {
	gt := myGoTickets{}
	if !gt.init(total) {
		errMsg := fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}

	return &gt, nil
}

func (gt *myGoTickets) init(total uint32) bool {
	if gt.active {
		return false
	}

	if total == 0 {
		return false
	}
	ch := make(chan byte, total)
	n := int(total)
	for i := 0; i < n; i++ {
		ch<-1
	}

	gt.ticketCh = ch
	gt.total = total
	gt.active = true
	return true
}

func (gt *myGoTickets) Take()  {
	<-gt.ticketCh
}

func (gt *myGoTickets) Return()  {
	gt.ticketCh <- 1
}

func (gt *myGoTickets) Active() bool {
	return gt.active
}

func (gt *myGoTickets) Total() uint32 {
	return gt.total
}

func (gt *myGoTickets) Remainder() uint32 {
	return uint32(len(gt.ticketCh))
}