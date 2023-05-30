package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"testing"
)

const TestRrn = "123456789012"

func getMockTerminal(nfun int, sleep string, t *testing.T) *MockTerminal {
	mockTerminal := NewMockTerminal(t)
	mockResultPay := NewMockResult(t)
	mockResultPay.On("GetVal").Return(int64(0))
	mockTerminal.EXPECT().CallMethod("NFun", nfun).Return(mockResultPay, nil)
	mockResultSleep := NewMockResult(t)
	mockResultSleep.On("ToString").Return(sleep)
	mockTerminal.EXPECT().CallMethod("GParamString", "Cheque1251").Return(mockResultSleep, nil)
	mockTerminal.EXPECT().CallMethod("Clear").Return(nil, nil)
	return mockTerminal
}

func checkResultFileContent(operationId uint32, result operationResult, t *testing.T) {
	f, err := os.Open(FileName)
	if err != nil {
		t.Errorf("open file error = %v", err)
	}
	defer f.Close()
	fstat, err := f.Stat()
	if err != nil {
		t.Errorf("stat file error = %v", err)
	}
	fmt.Println("fstat.Size() = ", fstat.Size())
	b1 := make([]byte, 4)
	n1, err := f.Read(b1)
	if err != nil {
		t.Errorf("read operation id error = %v", err)
	}
	if n1 != 4 {
		t.Errorf("read operation id n1 = %v, want %v", n1, 4)
	}
	if binary.BigEndian.Uint32(b1) != operationId {
		t.Errorf("read operation id b3 = %v, want %v", binary.BigEndian.Uint32(b1), operationId)
	}
	b2 := make([]byte, 4)
	n2, err := f.Read(b2)
	if err != nil {
		t.Errorf("read code error = %v", err)
	}
	if n2 != 4 {
		t.Errorf("read code n2 = %v, want %v", n2, 4)
	}
	if binary.BigEndian.Uint32(b2) != result.code {
		t.Errorf("read code b3 = %v, want %v", binary.BigEndian.Uint32(b2), result.code)
	}
	b3 := make([]byte, 12)
	n3, err := f.Read(b3)
	if err != nil {
		t.Errorf("read RRN error = %v", err)
	}
	if n3 != 12 {
		t.Errorf("read RRN n3 = %v, want %v", n3, 12)
	}
	if result.rrn != "" && string(b3) != result.rrn {
		t.Errorf("read RRN b3 = %v, want %v", string(b3), result.rrn)
	}
	b4 := make([]byte, int(fstat.Size())-n1-n2-n3)
	_, err = f.Read(b4)
	if err != nil {
		t.Errorf("read sleep error = %v", err)
	}
	if string(b4) != result.sleep {
		t.Errorf("read sleep b4 = %v, want %v", string(b4), result.sleep)
	}
}

func Test_outputResultToFile(t *testing.T) {
	var operationId uint32 = 10
	result := operationResult{
		code:  0,
		sleep: "Тестовый слип",
		rrn:   "123456789012",
	}
	err := outputResultToFile(operationId, result)
	if err != nil {
		t.Errorf("outputResultToFile() error = %v", err)
	}
	checkResultFileContent(operationId, result, t)
}

func Test_result_getMessage(t *testing.T) {
	r1 := operationResult{
		code: 99,
	}
	if r1.getMessage() != "Пинпад не подключен" {
		t.Errorf("bad message on code 99 = %v, want 'Пинпад не подключен'", r1.getMessage())
	}
	r2 := operationResult{
		code:  0,
		sleep: "Тестовый слип",
	}
	if r2.getMessage() != "Тестовый слип" {
		t.Errorf("bad message on code 0 = %v, want 'Тестовый слип'", r2.getMessage())
	}
}

func Test_payOperation(t *testing.T) {
	amount := 500
	sleep := fmt.Sprintf("Тестовый слип на %v", amount)
	mockTerminal := getMockTerminal(4000, sleep, t)
	mockTerminal.EXPECT().CallMethod("SParam", "Amount", amount).Return(nil, nil)
	mockResultRrn := NewMockResult(t)
	mockResultRrn.On("ToString").Return(TestRrn)
	mockTerminal.EXPECT().CallMethod("GParamString", "RRN").Return(mockResultRrn, nil)
	result, err := payOperation(mockTerminal, amount)
	if err != nil {
		t.Errorf("payOperation() error = %v", err)
	}
	if result.sleep != sleep {
		t.Errorf("payOperation() result.sleep = %v, want %v", result.sleep, sleep)
	}
	if result.code != 0 {
		t.Errorf("payOperation() result.code = %v, want %v", result.code, 0)
	}
	if result.rrn != TestRrn {
		t.Errorf("payOperation() result.rrn = %v, want %v", result.rrn, TestRrn)
	}
}

func Test_runOperation_pay(t *testing.T) {
	sleep := "Тестовый слип на 400"
	params := []string{"", "20", "pay", "400"}
	mockTerminal := getMockTerminal(4000, sleep, t)
	mockTerminal.EXPECT().CallMethod("SParam", "Amount", 400).Return(nil, nil)
	mockResultRrn := NewMockResult(t)
	mockResultRrn.On("ToString").Return(TestRrn)
	mockTerminal.EXPECT().CallMethod("GParamString", "RRN").Return(mockResultRrn, nil)
	runOperation(mockTerminal, params)
	result := operationResult{
		code:  0,
		sleep: sleep,
		rrn:   TestRrn,
	}
	checkResultFileContent(20, result, t)
}

func Test_returnOperation(t *testing.T) {
	amount := 600
	rrn := "123456789013"
	sleep := fmt.Sprintf("Тестовый слип на %v", amount)
	mockTerminal := getMockTerminal(4002, sleep, t)
	mockTerminal.EXPECT().CallMethod("SParam", "Amount", amount).Return(nil, nil)
	mockTerminal.EXPECT().CallMethod("SParam", "RRN", rrn).Return(nil, nil)
	mockTerminal.EXPECT().CallMethod("SParam", "Track2", "QSELECT").Return(nil, nil)
	mockResultRrn := NewMockResult(t)
	mockResultRrn.On("ToString").Return(TestRrn)
	mockTerminal.EXPECT().CallMethod("GParamString", "RRN").Return(mockResultRrn, nil)
	result, err := returnOperation(mockTerminal, amount, rrn)
	if err != nil {
		t.Errorf("returnOperation() error = %v", err)
	}
	if result.sleep != sleep {
		t.Errorf("returnOperation() result.sleep = %v, want %v", result.sleep, sleep)
	}
	if result.code != 0 {
		t.Errorf("returnOperation() result.code = %v, want %v", result.code, 0)
	}
	if result.rrn != TestRrn {
		t.Errorf("returnOperation() result.rrn = %v, want %v", result.rrn, TestRrn)
	}
}

func Test_runOperation_return_withRrn(t *testing.T) {
	rrn := "123456789013"
	amount := 600
	params := []string{"", "24", "return", strconv.Itoa(amount), rrn}
	sleep := fmt.Sprintf("Тестовый слип на %v", amount)
	mockTerminal := getMockTerminal(4002, sleep, t)
	mockTerminal.EXPECT().CallMethod("SParam", "Amount", amount).Return(nil, nil)
	mockTerminal.EXPECT().CallMethod("SParam", "RRN", rrn).Return(nil, nil)
	mockTerminal.EXPECT().CallMethod("SParam", "Track2", "QSELECT").Return(nil, nil)
	mockResultRrn := NewMockResult(t)
	mockResultRrn.On("ToString").Return(TestRrn)
	mockTerminal.EXPECT().CallMethod("GParamString", "RRN").Return(mockResultRrn, nil)
	runOperation(mockTerminal, params)
	result := operationResult{
		code:  0,
		sleep: sleep,
		rrn:   TestRrn,
	}
	checkResultFileContent(24, result, t)
}

func Test_runOperation_return_woRrn(t *testing.T) {
	amount := 600
	params := []string{"", "25", "return", strconv.Itoa(amount), ""}
	sleep := fmt.Sprintf("Тестовый слип на %v", amount)
	mockTerminal := getMockTerminal(4002, sleep, t)
	mockTerminal.EXPECT().CallMethod("SParam", "Amount", amount).Return(nil, nil)
	mockResultRrn := NewMockResult(t)
	mockResultRrn.On("ToString").Return(TestRrn)
	mockTerminal.EXPECT().CallMethod("GParamString", "RRN").Return(mockResultRrn, nil)
	runOperation(mockTerminal, params)
	result := operationResult{
		code:  0,
		sleep: sleep,
		rrn:   TestRrn,
	}
	checkResultFileContent(25, result, t)
}

func Test_runOperation_cancel(t *testing.T) {
	rrn := "123456789013"
	amount := 605
	params := []string{"", "30", "cancel", strconv.Itoa(amount), rrn}
	sleep := fmt.Sprintf("Тестовый слип отмены")
	mockTerminal := getMockTerminal(4003, sleep, t)
	mockTerminal.EXPECT().CallMethod("SParam", "Amount", amount).Return(nil, nil)
	mockTerminal.EXPECT().CallMethod("SParam", "RRN", rrn).Return(nil, nil)
	mockTerminal.EXPECT().CallMethod("SParam", "Track2", "QSELECT").Return(nil, nil)
	mockResultRrn := NewMockResult(t)
	mockResultRrn.On("ToString").Return(TestRrn)
	mockTerminal.EXPECT().CallMethod("GParamString", "RRN").Return(mockResultRrn, nil)
	runOperation(mockTerminal, params)
	result := operationResult{
		code:  0,
		sleep: sleep,
		rrn:   TestRrn,
	}
	checkResultFileContent(30, result, t)
}

func Test_runOperation_close(t *testing.T) {
	params := []string{"", "30", "close"}
	sleep := "Смена закрыта"
	terminal := getMockTerminal(6000, sleep, t)
	runOperation(terminal, params)
	result := operationResult{sleep: sleep, rrn: ""}
	checkResultFileContent(30, result, t)
}
