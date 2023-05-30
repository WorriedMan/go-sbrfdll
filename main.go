package main

import (
	"encoding/binary"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/pkg/errors"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"strconv"
)

// FileName файл, в который будут выводиться результаты операций. Структуру файла см. в функции outputResultToFile
const FileName = "result.bin"

func main() {
	initLogger()
	sber, err := getObject()
	if err != nil {
		log.Fatal("Unable to get object", err)
	}
	runOperation(sber, os.Args)
}

// initLogger инициализирует логгер, который пишет в файл sber.log
func initLogger() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "sber.log",
		MaxSize:    1,
		MaxBackups: 2,
		MaxAge:     28,
	})
}

// runOperation запускает выполнение операции с аргументами, указанными в params
// для оплаты: [1] - id операции, [2] - pay/return, [3] - сумма
// для возврата: [1] - id операции, [2] - pay/return, [3] - сумма (опционально 0), [4] - rrn (опционально)
// для отмены: [1] - id операции, [2] - cancel, [3] - сумма (опционально 0), [4] - rrn
// для закрытия смены: [1] - id операции, [2] - close
func runOperation(sber Terminal, params []string) {
	u64, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		log.Fatal("Unable to convert operation id", params[1], err)
	}
	id := uint32(u64)
	operation := params[2]
	var result operationResult
	switch operation {
	case "pay":
		total, err := strconv.Atoi(params[3])
		if err != nil {
			log.Fatal("Unable to get total", err)
		}
		result, err = payOperation(sber, total)
		if err != nil {
			log.Fatal("Unable to complete pay operation", err)
		}
	case "return":
		total, err := strconv.Atoi(params[3])
		if err != nil {
			log.Fatal("Unable to get total", err)
		}
		var rrn string
		if len(params) > 4 {
			rrn = params[4]
		}
		result, err = returnOperation(sber, total, rrn)
		if err != nil {
			log.Fatal("Unable to complete return operation", err)
		}
	case "cancel":
		total, err := strconv.Atoi(params[3])
		result, err = cancelOperation(sber, total, params[4])
		if err != nil {
			log.Fatal("Unable to complete cancel operation", err)
		}
	case "close":
		result, err = shiftCloseOperation(sber)
		if err != nil {
			log.Fatal("Unable to complete close operation", err)
		}
	default:
		log.Fatal("Unknown operation", operation)
	}
	err = outputResultToFile(id, result)
	if err != nil {
		log.Fatal("Unable to output result to file", err)
	}
}

// outputResultToFile выводит результат выполнения операции в файл, в котором первые 4 байта - ид операции,
// вторые 4 байта - код результата, следующие 12 байт - RRN, остальные байты - сообщение об ошибке
// (если код результата не 0) либо банковский слип
func outputResultToFile(operationId uint32, result operationResult) error {
	f, err := os.OpenFile(FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal("Unable to close file", err)
		}
	}()
	err = binary.Write(f, binary.BigEndian, operationId)
	if err != nil {
		return errors.Wrap(err, "write operation id")
	}
	err = binary.Write(f, binary.BigEndian, result.code)
	if err != nil {
		return errors.Wrap(err, "write code")
	}
	rrn := make([]byte, 12)
	copy(rrn, result.rrn)
	_, err = f.Write(rrn)
	_, err = f.Write([]byte(result.getMessage()))
	if err != nil {
		return errors.Wrap(err, "write sleep")
	}
	return nil
}

// getObject инициализирует и возвращает OLE объект Сбера
func getObject() (TerminalOle, error) {
	var t TerminalOle
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		return t, errors.Wrap(err, "initialize")
	}
	unknown, err := oleutil.CreateObject("SBRFSRV.Server")
	if err != nil {
		return t, errors.Wrap(err, "create object")
	}
	sber, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return t, errors.Wrap(err, "query interface")
	}
	_, err = oleutil.CallMethod(sber, "Clear")
	if err != nil {
		return t, errors.Wrap(err, "clear")
	}
	t.dispatch = sber
	return t, nil
}

// runTerminalFunction вызывает NFun терминала и возвращает результат операции.
// Если rrn == true, то в результат также загружается RRN операции
func runTerminalFunction(sber Terminal, nfun int, rrn bool) (operationResult, error) {
	var r operationResult
	resPay, err := sber.CallMethod("NFun", nfun)
	if err != nil {
		return r, errors.Wrap(err, "call nfun")
	}
	r.code = uint32(resPay.GetVal())
	resCheque, err := sber.CallMethod("GParamString", "Cheque1251")
	if err != nil {
		return r, errors.Wrap(err, "getting cheque")
	}
	r.sleep = resCheque.ToString()
	if rrn {
		resRrn, err := sber.CallMethod("GParamString", "RRN")
		if err != nil {
			return r, errors.Wrap(err, "getting rrn")
		}
		r.rrn = resRrn.ToString()
	}
	_, err = sber.CallMethod("Clear")
	if err != nil {
		return r, errors.Wrap(err, "clear")
	}
	log.Println("Got code:", r.code)
	if r.code == 0 {
		log.Println("Got sleep:\n", r.sleep)
		if rrn {
			log.Println("Got RRN:", r.rrn)
		}
	}
	return r, nil
}

// payOperation вызывает на терминале операцию оплаты с суммой total (в копейках)
func payOperation(sber Terminal, total int) (operationResult, error) {
	var r operationResult
	_, err := sber.CallMethod("SParam", "Amount", total)
	if err != nil {
		return r, errors.Wrap(err, "call amount")
	}
	log.Println("-----Pay-----")
	r, err = runTerminalFunction(sber, 4000, true)
	if err != nil {
		return r, errors.Wrap(err, "run terminal function")
	}
	return r, nil
}

// returnOperation вызывает на терминале операцию возврата с суммой total (в копейках) и RRN
func returnOperation(sber Terminal, total int, rrn string) (operationResult, error) {
	var r operationResult
	_, err := sber.CallMethod("SParam", "Amount", total)
	if err != nil {
		return r, errors.Wrap(err, "call amount")
	}
	if rrn != "" {
		_, err = sber.CallMethod("SParam", "RRN", rrn)
		if err != nil {
			return r, errors.Wrap(err, "call rrn")
		}
		_, err = sber.CallMethod("SParam", "Track2", "QSELECT")
		if err != nil {
			return r, errors.Wrap(err, "call track2")
		}
	}
	log.Println("-----Return-----")
	r, err = runTerminalFunction(sber, 4002, true)
	if err != nil {
		return r, errors.Wrap(err, "run terminal function")
	}
	return r, nil
}

// cancelOperation вызывает на терминале операцию отмены операции по RRN
func cancelOperation(sber Terminal, amount int, rrn string) (operationResult, error) {
	var r operationResult
	_, err := sber.CallMethod("SParam", "Amount", amount)
	if err != nil {
		return r, errors.Wrap(err, "call amount")
	}
	_, err = sber.CallMethod("SParam", "RRN", rrn)
	if err != nil {
		return r, errors.Wrap(err, "call rrn")
	}
	_, err = sber.CallMethod("SParam", "Track2", "QSELECT")
	if err != nil {
		return r, errors.Wrap(err, "call track2")
	}
	log.Println("-----Cancel-----")
	r, err = runTerminalFunction(sber, 4003, true)
	if err != nil {
		return r, errors.Wrap(err, "run terminal function")
	}
	return r, nil
}

// shiftCloseOperation вызывает на терминале операцию сверки итогов
func shiftCloseOperation(sber Terminal) (operationResult, error) {
	var r operationResult
	log.Println("-----Shift close-----")
	r, err := runTerminalFunction(sber, 6000, false)
	if err != nil {
		return r, errors.Wrap(err, "run terminal function")
	}
	return r, nil
}
