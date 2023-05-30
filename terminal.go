package main

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

//go:generate mockery --name Result
type Result interface {
	ToString() string
	GetVal() int64
}

//go:generate mockery --name Terminal
type Terminal interface {
	CallMethod(name string, params ...interface{}) (result Result, err error)
}

type TerminalOle struct {
	dispatch *ole.IDispatch
}

func (t TerminalOle) CallMethod(name string, params ...interface{}) (Result, error) {
	var variant *ole.VARIANT
	variant, err := oleutil.CallMethod(t.dispatch, name, params...)
	r := ResultOle{variant: variant}
	if err != nil {
		return r, fmt.Errorf("call ole method %v, %w", name, err)
	}
	return r, nil
}

type ResultOle struct {
	variant *ole.VARIANT
}

func (r ResultOle) ToString() string {
	return r.variant.ToString()
}

func (r ResultOle) GetVal() int64 {
	return r.variant.Val
}

type operationResult struct {
	code       uint32
	rrn, sleep string
}

func (o operationResult) getMessage() string {
	if o.code == 0 {
		return o.sleep
	}
	return codeToString(o.code)
}
