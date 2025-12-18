package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 指针1
	//var i int = 99
	//point1(&i)
	//fmt.Println(i)

	// 指针2
	//var nums []int = []int{1, 2, 3, 4, 5}
	//point2(&nums)
	//fmt.Println(nums)

	// Goroutine1
	//goroutine1()

	// Goroutine2
	//var f []func() = make([]func(), 3)
	//f[0] = task1
	//f[1] = task2
	//f[2] = task3
	//goroutine2(f)

	//oo1()
	//oo2()

	//channel1()
	//channel2()

	mutex_lock1()
	//mutex_lock2()
}

// 指针
// 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
// 考察点 ：指针的使用、值传递与引用传递的区别。
func point1(num *int) {
	*num += 10
}

// 指针
// 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
// 考察点 ：指针运算、切片操作。
func point2(num *[]int) {
	for i := 0; i < len(*num); i++ {
		(*num)[i] *= 2
	}
}

// Goroutine
// 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// 考察点 ： go 关键字的使用、协程的并发执行。
func goroutine1() {
	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Print(i, ",")
		}
		fmt.Println()
	}()

	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Print(i, ",")
		}
		fmt.Println()
	}()

	time.Sleep(time.Second * 2)
}

// Goroutine
// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 考察点 ：协程原理、并发任务调度。
func goroutine2(f []func()) {
	for i := 0; i < len(f); i++ {
		go f[i]()
	}
	time.Sleep(time.Second * 5)
}
func task1() {
	startTime := time.Now()
	time.Sleep(time.Second * 1)
	endTime := time.Now()
	fmt.Println("task1 耗时：", endTime.Sub(startTime))
}
func task2() {
	startTime := time.Now()
	time.Sleep(time.Second * 2)
	endTime := time.Now()
	fmt.Println("task2 耗时：", endTime.Sub(startTime))
}
func task3() {
	startTime := time.Now()
	time.Sleep(time.Second * 3)
	endTime := time.Now()
	fmt.Println("task3 耗时：", endTime.Sub(startTime))
}

// 面向对象
// 题目 ：定义 Shape 接口，包含 Area() 和 Perimeter() 两个方法。创建 Rectangle 和 Circle 结构体，实现 Shape 接口。 创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
// 考察点 ：接口的定义与实现、面向对象编程风格。
func oo1() {
	r := Rectangle{2, 3}
	fmt.Println("r矩形的面积：", r.Area())
	fmt.Println("r矩形的周长：", r.Perimeter())

	c := Circle{5}
	fmt.Println("c圆形的面积：", c.Area())
	fmt.Println("c圆形的周长：", c.Perimeter())
}

type Shape interface {
	Area() float64
	Perimeter() float64
}
type Rectangle struct {
	height float64
	width  float64
}

func (r Rectangle) Area() float64 {
	return r.height * r.width
}
func (r Rectangle) Perimeter() float64 {
	return (r.height + r.width) * 2
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}
func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.radius
}

// 面向对象
// 题目 ：使用组合的方式创建 Person 结构体，包含 Name 和 Age 字段，再创建 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
// 考察点 ：组合的使用、方法接收者。
func oo2() {
	p := Person{"张天下", 18}
	e := Employee{p, 9527}
	e.PrintInfo()
}

type Person struct {
	Name string
	Age  int
}
type Employee struct {
	Person     Person
	EmployeeID int
}

func (e Employee) PrintInfo() {
	fmt.Println("员工:", e.Person.Name, ", 年龄:", e.Person.Age, ", 编号:", e.EmployeeID)
}

// Channel
// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
// 考察点 ：通道的基本使用、协程间通信。
func channel1() {
	ch := make(chan int)
	go sendOnly(ch, 10)
	go receiveOnly(ch)
	timeout := time.After(time.Second * 10) // 防止死循环，限定退出时间

	for {
		select {
		case v, ok := <-ch:
			if ok {
				fmt.Println("主进程读取 channel 数据：", v)
			} else {
				fmt.Println("channel close")
				return
			}
		case <-timeout:
			fmt.Println("timeout exit")
			return
		}
	}
}
func sendOnly(ch chan<- int, l int) {
	for i := 1; i <= l; i++ {
		ch <- i
		fmt.Println("写入 channel 数据：", i)
	}
	close(ch)
}
func receiveOnly(ch <-chan int) {
	for v := range ch {
		fmt.Println("读取 channel 数据：", v)
	}
}

// Channel
// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
// 考察点 ：通道的缓冲机制。
func channel2() {
	ch := make(chan int, 1)
	go sendOnly(ch, 100)
	go receiveOnly(ch)
	timeout := time.After(time.Second * 10) // 防止死循环，限定退出时间

	for {
		select {
		case v, ok := <-ch:
			if ok {
				fmt.Println("主进程读取 channel 数据：", v)
			} else {
				fmt.Println("channel close")
				return
			}
		case <-timeout:
			fmt.Println("timeout exit")
			return
		}
	}
}

// 锁机制
// 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：sync.Mutex 的使用、并发数据安全。
func mutex_lock1() {
	mu := sync.Mutex{}
	sc := SafeCounter{0}

	for i := 0; i < 10; i++ {
		go func() {
			mu.Lock()
			defer mu.Unlock()
			for j := 0; j < 1000; j++ {
				sc.c++
			}
		}()
	}

	time.Sleep(time.Second * 5)
	fmt.Println(sc.c)
}

// 锁机制
// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：原子操作、并发数据安全。
func mutex_lock2() {
	sc := SafeCounter{0}
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				sc.c++
			}
		}()
	}

	time.Sleep(time.Second * 5)
	fmt.Println(sc.c)
}

type SafeCounter struct {
	c int
}
