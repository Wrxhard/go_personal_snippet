package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

// Custom function
func add(a, b int) int {
	return a + b
}

// Function actually create a copy of the variable not operate on the original variable itself
func add2(x int) int {
	// a temporary variable is created to store the result of x + 1
	x++
	// if we want to capture the result, we need to assign it back to x
	return x
}

// What happen if we don't return anything
func add3(x int) {
	// a temporary variable is created to store the result of x + 1
	x++
	// if we dont capture and assign a new value x wont change
}

// Multiple return
func minus(a, b int) (int, int, int) {
	return a, b, a - b
}

// Named return for better readability
func minus2(a, b int) (result int) {
	result = a - b
	return result
}

// Struct is basically a class in other language and the same as C
type Person struct {
	name string
	age  int
}

type work interface {
	hour(h int) int
	salary(s int) int
}

func implementsWorkInterface(worker interface{}) bool {
	if _, ok := worker.(work); ok {
		return true
	} else {
		return false
	}
}

func (w Worker) hour(h int) int {
	return h
}

func (w Worker) salary(s int) int {
	return s
}

// Nested struct
type Worker struct {
	name    string
	Person2 Person
}

// Embeded struct
type Student struct {
	name string
	Person
}

// Add function to struct
func (p Student) toString() string {
	return p.name + " " + strconv.Itoa(p.age)
}

// Pass a slice to a function
func sum(a ...int) int {
	total := 0
	for _, num := range a {
		total += num
	}
	return total
}

func add_pointer(a, b *int) {
	*a++
	*b++
}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func player(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	// Send a value to notify that we're done.
	done <- true
}

func check_database(tokenChan chan struct{}) {
	<-tokenChan
	fmt.Println("Database exist")
}
func initDatabase(version int) chan struct{} {
	tokenChan := make(chan struct{})

	//Perform some database initialization
	//...

	// return the token channel to indicate that the database is ready
	return tokenChan
}

func addEmailsToQueue(emails []string) chan string {
	emailChan := make(chan string, len(emails))

	go func() {
		for _, email := range emails {
			emailChan <- email
		}
		close(emailChan)
	}()

	return emailChan
}

func readEmailsFromQueue(batchsize int, emailChan chan string) {
	for email := range emailChan {
		fmt.Println(email)
	}

	// We can also close the channel if it has done all it job
	// close(emailChan)
	// check if the channel is close we can use
	/*
		if _, ok := <-emailChan; ok {
			//true indicate that the channel is still open

			//false indicate that the channel is close
		}
	*/
}

func fibonacci(n int, chInts chan int) {
	if n <= 1 {
		chInts <- n
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		chInts <- a
		a, b = b, a+b
	}
	close(chInts)
}

func concurrentFibonacci(n int) {
	chInts := make(chan int)

	go func() {
		fibonacci(n, chInts)
	}()

	//We can also use range to literate over the channel
	for i := range chInts {
		fmt.Println(i)
	}

}

func server1(ch chan string) {
	time.Sleep(time.Second * 1)
	ch <- "from server1"
}

func server2(ch chan string) {
	time.Sleep(time.Second * 2)
	ch <- "from server2"
	close(ch)
}

func logServer(server1 chan string, server2 chan string) {
	// Listen contienously for value from the channel
	for {
		select {
		case s1, ok := <-server1:
			if ok {
				fmt.Println(s1)
			} else {
				//exit the function and for loop if one of the channel is close
				return
			}
		case s2, ok := <-server2:
			if ok {
				fmt.Println(s2)
			} else {
				return
			}
		}
	}
}

func readOnlyChannel(ch <-chan string) {
	msg := <-ch
	fmt.Println(msg)
}

func writeOnlyChannel(ch chan<- string) {
	// We can only write to the channel
	ch <- "Hello"
	close(ch)
}

// Mutex is use to avoid race condition between goroutine
type Counter struct {
	value int
	mux   *sync.Mutex
}

func (c *Counter) Increment() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.value++
}

type SafeCounter struct {
	//Use RWMutex to allow multiple read but only one write at a time
	mux *sync.RWMutex
	v   map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	c.v[key]++
	c.mux.Unlock()
}

func (c *SafeCounter) Value(key string) int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.v[key]
}

// Generic
func addGeneric[T any](s []T) ([]T, error) {
	if len(s) == 0 {
		return nil, errors.New("Empty slice")
	}
	return s, nil
}

// Specific type of generic
func addGeneric2[T int | float64 | string](s []T) ([]T, error) {
	if len(s) == 0 {
		return nil, errors.New("Empty slice")
	}
	return s, nil
}

func main() {

	//auto detect type
	intNum := 1
	fmt.Printf("intNum: %d\n", intNum)

	//Go also support type conversion just like Python. Ex: intNum2 := int8(intNum)
	floatNum := float64(intNum)
	fmt.Printf("floatNum: %f\n", floatNum)

	//Explicitly declare type int8 mean it will take 8 bit of memory
	var intNum2 int8 = 1
	fmt.Printf("intNum2: %d\n", intNum2)

	//String
	str := "Hello World"
	str = "Hey" //string is immutable, so we need to reassign the value to change it

	/*Rune stand for character, In order to get the length of string,
	we need to use RuneCountInString from unicode/utf8 package */
	fmt.Printf("String Length: %d %s\n", utf8.RuneCountInString(str), str)

	//Normal if else
	len := utf8.RuneCountInString(str)
	if len > 0 {
		fmt.Printf("String Length: %d %s\n", len, str)
	} else {
		fmt.Println("Empty String")
	}

	//Better way of wrtting if else with short statement.
	//The variable declared in the if statement is wrap in the if scope
	if len := utf8.RuneCountInString(str); len > 0 {
		fmt.Printf("String Length: %d %s\n", len, str)
	}

	//Use a function
	fmt.Printf("add: %d\n", add(1, 2))

	x := 1
	//Capture the result of add2
	x = add2(x)
	fmt.Printf("add2: %d\n", x)

	add3(x)
	fmt.Printf("add3: %d\n", x)

	// Anoyomous struct
	person := struct {
		name string
		age  int
	}{
		name: "John",
		age:  30,
	}
	fmt.Printf("person: %v\n", person)

	// Initialize nested struct
	worker := Worker{
		name: "John",
		Person2: Person{
			age: 40,
		},
	}
	// Get value from nested struct
	fmt.Printf("worker.age: %d\n", worker.Person2.age)

	// Initialize embeded struct
	struct2 := Student{
		name: "John",
		Person: Person{
			age: 30,
		},
	}
	// Get value from embeded struct
	fmt.Printf("person.age: %d\n", struct2.age)

	fmt.Printf("person.toString(): %s\n", struct2.toString())

	y := 34

	z := &y

	*z = 4

	a := 3
	b := 4

	mem_str := "dadsds"
	fmt.Println(&mem_str)

	mem_str = "Hello"

	fmt.Println(y)

	add_pointer(&a, &b)
	fmt.Println(a, b)

	mem_str2 := &mem_str

	mem_str = "Heyd"

	if implementsWorkInterface(worker) {
		fmt.Println("true")
	}

	fmt.Println(*mem_str2)

	//Array
	ex_array := [3]int{1, 2, 3}
	fmt.Println(ex_array)

	//Slice is like an array but it can be resized it manipulate the original array inside memory under the hood and dont need to have fixed size
	ex_slice := []int{1, 2, 3}
	fmt.Println(ex_slice)
	//Slice of an array
	ex_slice2 := ex_array[1:2]
	fmt.Println(ex_slice2)
	//For each in slice
	for i, num := range ex_slice {
		fmt.Println(i, num)
	}
	//Sum of slice
	//The advantage of using ...int is that it provide more flexibility to the function for example we can pass sum(1,2,3) instead of sum([]int{1,2,3})
	fmt.Println(sum(1, 2, 3))
	//Spread operator slice... just like it name it spread the slice into individual element
	name := []int{1, 2}
	fmt.Println(sum(name...))
	//Append to slice
	name = append(name, 4)
	fmt.Println(name)
	//2D Matrix with slice
	matrix := [][]int{{1, 2}, {3, 4}}
	fmt.Println(matrix)
	//We can increase the speed by preallocate the slice
	matrix2 := make([][]int, 2)
	matrix2[0] = []int{1, 2}
	matrix2[1] = []int{3, 4}
	fmt.Println(matrix2)
	//We can append to the slice even though it is preallocate
	matrix2 = append(matrix2, []int{1, 2})
	fmt.Println(matrix2)
	//Map
	ex_map := map[string]int{"a": 1, "b": 2}
	fmt.Println(ex_map)

	// Now is the fun stuff Goroutine!!
	go say("world")
	go say("hello")

	//Anoyomous goroutine function
	go func() {
		fmt.Println("Running in a goroutine")
	}() //() in the end mean we are calling the function

	// Channels are the pipes that connect concurrent goroutines.
	// Start a worker goroutine, giving it the channel to notify on.
	done := make(chan bool, 1)
	go player(done)

	// Block until we receive a notification from the worker on the channel.
	bol := <-done
	fmt.Println(bol)

	//Buffer Channel
	emails_slice := []string{"Hello", "This", "Is", "An", "Example", "Of", "Buffer", "Channel"}
	batchsize := cap(emails_slice)
	// Batch size indicate how many email we want to send at once
	// the channel will hold all the email until it reach the batch size then send it
	emails_chan := addEmailsToQueue(emails_slice)
	// Read the email from the channel
	readEmailsFromQueue(batchsize, emails_chan)

	// Range over channel
	concurrentFibonacci(10)

	//Select is like switch case for channel what it will do is it will check ii any channel send out a value and execute the case
	output1 := make(chan string)
	output2 := make(chan string)

	go server1(output1)
	go server2(output2)

	logServer(output1, output2)

	//Specific use case for read only and write only channel
	spc_chan := make(chan string)
	go writeOnlyChannel(spc_chan)
	readOnlyChannel(spc_chan)

	//Mutex
	//Create a wait group
	var wg sync.WaitGroup

	//Create a counter with initial value of 0
	counter := Counter{
		mux:   &sync.Mutex{},
		value: 0,
	}

	for i := 0; i < 100; i++ {
		//Add 1 to the wait group mean we are adding 1 goroutine to the wait group
		wg.Add(1)
		go func() {
			//after the code is done executing from the goroutine we will remove 1 goroutine from the wait group
			defer wg.Done()
			//Increment the counter
			counter.Increment()
		}()
	}
	//Wait for all the goroutine to finish executing
	wg.Wait()
	//Print the value of the counter
	fmt.Println(counter.value)

	//RWMutex
	c := SafeCounter{
		mux: &sync.RWMutex{},
		v:   make(map[string]int),
	}
	var wg2 sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			c.Inc("key")
		}()
	}

	wg2.Wait()
	fmt.Println(c.Value("key"))

	//Generic
	s := []int{1, 2, 3}
	s2 := []string{"a", "b", "c"}

	s_res, _ := addGeneric(s)
	s2_res, _ := addGeneric(s2)
	fmt.Println(s_res)
	fmt.Println(s2_res)

	//Specific Generic
	s3 := []int{1, 2, 3}
	s4 := []string{"a", "b", "c"}

	s3_res, _ := addGeneric(s3)
	s4_res, _ := addGeneric(s4)
	fmt.Println(s3_res)
	fmt.Println(s4_res)

}
