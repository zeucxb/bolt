package main

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type asyncError struct {
	RunnerName string
	Err        error
}

var wg sync.WaitGroup

func run(name string, eChan chan<- asyncError, f interface{}, params ...interface{}) {
	fValue := reflect.ValueOf(f)
	fType := fValue.Type()

	if fType.Kind() != reflect.Func {
		eChan <- asyncError{name, errors.New("We need a valid function to work")}
		return
	}

	err := make(chan error)
	errValue := reflect.ValueOf(err)

	in := []reflect.Value{errValue}

	for _, param := range params {
		pValue := reflect.ValueOf(param)
		in = append(in, pValue)
	}

	wg.Add(2)
	go fValue.Call(in)
	go func() {
		select {
		case e := <-err:
			defer wg.Done()
			if e != nil {
				supervisor(asyncError{name, e})
			}
		}
	}()
	wg.Wait()
}

func stop(eChan chan<- error, err error) {
	defer wg.Done()
	defer close(eChan)

	if err != nil {
		eChan <- err
	}
}

func supervisor(err asyncError) {
	fmt.Println(err.RunnerName)
	fmt.Println(err.Err)
}

func main() {
	err := make(chan asyncError)
	
	run("HARD TASK 1", err, hardTask, []int{1, 4, 6, 7})
	run("HARD TASK 2", err, hardTask, []int{1, 4, 6, 7})
	run("HARD TASK 3", err, hardTask, []int{1, 4, 6, 7})
	run("HARD TASK 4", err, hardTask, []int{1, 4, 6, 7})
	run("HARD TASK 5", err, hardTask, []int{1, 4, 6, 7})
}

func hardTask(eChan chan<- error, numbers []int) {
	fmt.Println(numbers)
	// stop(eChan, nil)
	stop(eChan, errors.New("EROOOOOR!!!"))
}
