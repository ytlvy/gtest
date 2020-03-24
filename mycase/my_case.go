package mycase

import (
	"fmt"
	"sync"
)

type MyCase struct {
}

func (t *MyCase) TestRoutine() {
	a := make([]int, 0)

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			a = append(a, i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println(len(a))
}

func nilOrNot(v interface{}) bool {
	return v == nil
}

func (t *MyCase) TestInterfacePara() {
	type TestStruct struct{}

	var s *TestStruct
	fmt.Println(s == nil)
	fmt.Println(nilOrNot(s)) // 因为s传递过程中 转换为 interface{}结构 保留了原来的结构体定义信息
}

func (t *MyCase) TestRoutineSync() {
	a := make([]int, 0)
	ch := make(chan int)

	go func() {
		for v := range ch {
			a = append(a, v)
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			ch <- i
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println(len(a))
}
