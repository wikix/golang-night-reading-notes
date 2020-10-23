## 带着问题去阅读golang 源码

问题1：golang 怎么建立一个一一映射，现在的映射map是可以多对一的吧？

A：建立两个 map 互查。 segment tree ？



问题2：golang中make和new关键字的区别和使用场景？

```

代码来源：
https://github.com/golang/go/blob/master/src/builtin/builtin.go
go/src/builtin/builtin.go

// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
//	Slice: The size specifies the length. The capacity of the slice is
//	equal to its length. A second integer argument may be provided to
//	specify a different capacity; it must be no smaller than the
//	length. For example, make([]int, 0, 10) allocates an underlying array
//	of size 10 and returns a slice of length 0 and capacity 10 that is
//	backed by this underlying array.
//	Map: An empty map is allocated with enough space to hold the
//	specified number of elements. The size may be omitted, in which case
//	a small starting size is allocated.
//	Channel: The channel's buffer is initialized with the specified
//	buffer capacity. If zero, or the size is omitted, the channel is
//	unbuffered.
func make(t Type, size ...IntegerType) Type

// The new built-in function allocates memory. The first argument is a type,
// not a value, and the value returned is a pointer to a newly
// allocated zero value of that type.
func new(Type) *

```

`new` 它接受一个参数，这个参数是类型，不是一个值，分配好内存后，返回一个指向该类型内存地址的指针，同时请注意：它把分配的内存内置为 零，也就是类型的 零值。



new 不常用：new 在一些需要实例化接口的地方用的比较多，但是可以用 &A{}替代。



但是 new 和 &A{} 也是有差别的，主要差别在于 &A{} 显示执行堆分配。



make 也是用于内存分配的，和 new 不同，它只用于 channel、map以及slice的内存创建，而且它返回的类型就这三个类型本身，而不是它们的指针类型，因为这三个类型就是引用类型，所以就没有必要返回它们的指针了。



	> 注意：这三种类型是引用类型，所以必须得初始化，但是否置为零值，这个和new是不一样的



举例说明

```go
package main

import (
  "fmt"
)

func main(){
  p := new([]int) // p == nil; with len and cap 0
  fmt.Println(p)
  
  v := make([]int, 10, 50) // v is intialized with len 10, cap 50
  fmt.Println(v)
  
  /***** Output *****/
  &[]
  [0 0 0 0 0 0 0 0 0]
  *******************/
  
  (*p)[0] = 18 // panic: runtime error: index out of range
               // because p is a nil pointer, with len and cap 0
  
  v[1] = 18  // ok
}

运行结果：
➜  code git:(master) ✗ go run test-new-make.go
&[]
[0 0 0 0 0 0 0 0 0 0]
panic: runtime error: index out of range

goroutine 1 [running]:
main.main()
	/Users/pauljobs/code/test-new-make.go:19 +0x15d
```



Effective Go对 make 和 new 的 区别解释

```go
Allocation with new
Go has two allocation primitives, the built-in functions new and make. They do different things and apply to different types, which can be confusing, but the rules are simple. Let's talk about new first. It's a built-in function that allocates memory, but unlike its namesakes in some other languages it does not initialize the memory, it only zeros it. That is, new(T) allocates zeroed storage for a new item of type T and returns its address, a value of type *T. In Go terminology, it returns a pointer to a newly allocated zero value of type T.

Since the memory returned by new is zeroed, it's helpful to arrange when designing your data structures that the zero value of each type can be used without further initialization. This means a user of the data structure can create one with new and get right to work. For example, the documentation for bytes.Buffer states that "the zero value for Buffer is an empty buffer ready to use." Similarly, sync.Mutex does not have an explicit constructor or Init method. Instead, the zero value for a sync.Mutex is defined to be an unlocked mutex.

The zero-value-is-useful property works transitively. Consider this type declaration.

type SyncedBuffer struct {
    lock    sync.Mutex
    buffer  bytes.Buffer
}
Values of type SyncedBuffer are also ready to use immediately upon allocation or just declaration. In the next snippet, both p and v will work correctly without further arrangement.

p := new(SyncedBuffer)  // type *SyncedBuffer
var v SyncedBuffer      // type  SyncedBuffer
```





----



#### stackoverflow 

question



In Golang you can allocate memory for a slice with the below syntax :

```golang
my_slice := make( []int, 0 )
```

And then later on I can add elements with the built-in append function as :

```golang
my_slice := append(my_slice, 23)
```

My question is, what's the difference between giving that zero ( or 2 or 5 or whatever) when "making" the slice if later on we can keep adding items as long as we wan?

Is there a performance bonus by trying to guess the capacity that slice will end up having?





#### answer 1

> My question is, what's the difference between giving that zero ( or 2 or 5 or whatever) when "making" the slice if later on we can keep adding items as long as we wan?

A not appropriate allocation can cause re-allocation.

> Is there a performance bonus by trying to guess the capacity that slice will end up having?

yes, re-allocation will lead to additional copy on slice.

> Slice: The size specifies the length. The capacity of the slice is equal to its length. A second integer argument may be provided to specify a different capacity; it must be no smaller than the length.
>
> For example, make([]int, 0, 10) allocates an underlying array of size 10 and returns a slice of length 0 and capacity 10 that is backed by this underlying array.

copy from go doc.

In my view. Slice is **like a pointer with length and capacity** in c/c++. Append to a slice will append elements after the offset on that pointer. The capacity is the total sequential space size. Once the `capacity`- `allocated` space is not enough, append will cause a re-allocation and a copy on the slice.

> ```
> make(s, 1)
> ```

go will do more than you think:

1. Allocate a sequential space larger than the size you give (same in c++ vector) to avoid the re-allocation which can cause low performance.
2. Initialize the size you give in make. (RAII)
3. Once reallocation happened, go will allocate a twice size sequential space and copy the old slice to that place. This will also low down the performance.

To avoid the reallocation happen, we can give the optional capacity arguments in `make` to tell the go we need a bigger space.



----

#### answer 2

The difference is that the memory for the slice is allocated upfront and `len(mySlice)` returns the total slice length.

Performance wise it is beneficial to allocate the size upfront because when you call `a = append(a, n)` the following occurs:

- It calls the builtin append function and for that it first copies the `a` slice (slice header, backing array is not part of the header), and it has to create a temporary slice for the variadic parameter which will contain the value `n`.
- Then it has to reslice `a` if it has enough capacity like `a = a[:len(a)+1]` - which involves assigning the new slice to `a` inside the append function. If a would not have big enough capacity to do the append "in-place" then a new array would have to be allocated, content from slice copied, and then the assign / append be executed.
- Then assigns `n` to a `[len(a)-1]`.
- Then returns the new slice from the append function, and this new slice is assigned to the local variable `a`.

Compared to `a[i] = n` which is a simple assignment.

- https://golang.org/ref/spec#Assignments
- https://blog.golang.org/slices



问题验证代码

```go
package main

import (
	"fmt"
)

func main() {
	v1 := make([]int, 0, 10) // v is intialized with len 10, cap 50
	fmt.Println(v1)

	v2 := make([]int, 10, 10)
	fmt.Println(v2)

	v1 = append(v1, 18)
	fmt.Println(v1)

	v2 = append(v2, 18)
	fmt.Println(v2)
}

运行结果：
➜  code git:(master) ✗ go run test-new-make.go
[]
[0 0 0 0 0 0 0 0 0 0]
[18]
[0 0 0 0 0 0 0 0 0 0 18]
```



上述问题根本原因：append函数源码

```go
// The append built-in function appends elements to the end of a slice. If
// it has sufficient capacity, the destination is resliced to accommodate the
// new elements. If it does not, a new underlying array will be allocated.
// Append returns the updated slice. It is therefore necessary to store the
// result of append, often in the variable holding the slice itself:
//	slice = append(slice, elem1, elem2)
//	slice = append(slice, anotherSlice...)
// As a special case, it is legal to append a string to a byte slice, like this:
//	slice = append([]byte("hello "), "world"...)
func append(slice []Type, elems ...Type) []
```

----



#### 新问题：slice 自动扩容如何扩？按怎样的维度扩容？

>  slice，小于1024，两倍增长（申请内存）；大于1024，四分之一增长（申请内存）
>
> 即大家所说
>
> slice扩容，cap不够1024的，直接翻倍；cap超过1024的，新cap变为老cap的1.25倍。



验证代码

```go
package main

import (
        "fmt"
)

func main() {
        slice1 := make([]int, 2)
        fmt.Println("cap of slice1", cap(slice1))
        slice1 = append(slice1, 1)
        fmt.Println("cap of slice1", cap(slice1))
        slice1 = append(slice1, 2)
        fmt.Println("cap of slice1", cap(slice1))

        fmt.Println("-------------")

        slice1024 := make([]int, 1024)
        fmt.Println("cap of slice1024", cap(slice1024))
        slice1024 = append(slice1024, 1)
        fmt.Println("cap of slice1024", cap(slice1024))
        slice1024 = append(slice1024, 2)
        fmt.Println("cap of slice1024", cap(slice1024))
        for i := 0; i < (1280 - 1024 - 1); i++ {
                slice1024 = append(slice1024, i)
        }
        fmt.Println("cap of slice1024", cap(slice1024))
}

输出结果：
➜  code git:(master) ✗ go run test-slice.go
cap of slice1 2
cap of slice1 4
cap of slice1 4
-------------
cap of slice1024 1024
cap of slice1024 1280
cap of slice1024 1280
cap of slice1024 1696

1696的原因：内存对齐，管理机制原理是？
```



相关阅读

```go
golang slice源码
https://github.com/golang/go/blob/master/src/runtime/slice.go#L125

go runtime slice
相关推荐阅读文章
https://eddycjy.com/posts/go/slice/2018-12-11-slice/
https://jodezer.github.io/2017/05/golangSlice%E7%9A%84%E6%89%A9%E5%AE%B9%E8%A7%84%E5%88%99
https://juejin.im/post/6844903812331732999
```



----

Go的逃逸分析 决定了 是分配到堆上还是栈上



参考链接：

www.zenlife.tk/go-allocated-on-heap-or-stack.md

https://studygolang.com/articles/12444



验证代码

```go
package main

type user struct {
	name  string
	email string
}

func main() {
	u1 := createUserV1()
	u2 := createUserV2()

	println("u1", &u1, "u2", u2)
}

//go:noinline
func createUserV1() user {
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V1", &u)
	return u
}

//go:noinline
func createUserV2() *user {
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V2", &u)
	return &u
}


运行结果：
➜  code git:(master) ✗ go run test-escape.go
V1 0xc000032710
V2 0xc00000c060
u1 0xc000032768 u2 0xc00000c060

```





Go的变量内存分配到底是在堆还是栈？

```go
示例代码
package main

import (
	"fmt"
)

const Width, Height = 640, 480

type Cursor struct {
	X, Y int
}

func Center(c *Cursor) {
	c.X += Width / 2
	c.Y += Height / 2
}

func CenterCursor() {
	c := new(Cursor)
	Center(c)
	fmt.Println(c.X, c.Y)
}

运行结果：
➜  code git:(master) ✗ go tool compile -m test-var.go
test-var.go:13:6: can inline Center
test-var.go:20:8: inlining call to Center
test-var.go:21:13: inlining call to fmt.Println
test-var.go:13:13: Center c does not escape
test-var.go:21:15: c.X escapes to heap
test-var.go:21:20: c.Y escapes to heap
test-var.go:21:13: io.Writer(os.Stdout) escapes to heap
test-var.go:19:10: CenterCursor new(Cursor) does not escape
test-var.go:21:13: CenterCursor []interface {} literal does not escape
<autogenerated>:1: os.(*File).close .this does not escape
```



compile -S 是打印汇编中间代码

compile -m 是打印出编译优化。从输出上看，它说	`new(Cursor)` 没有escape，于是在栈上分配了。等价于 C 的写法：



```c
void CenterCursor(){
    struct Cursor c;
    Center(&c);
}
```



再看下面一段代码

```go
package main

func main(){
  var a [1]int
  c := a[:]
  println(c)
}
```





#### 逃逸分析的用处（For 性能）

- 减少gc的压力，不逃逸的对象分配在栈上，当函数返回时就回收了资源，不需要gc标记清除
- 逃逸分析完后，可以确定哪些变量分配在栈上，栈的分配比堆快，性能好
- 同步消除，如果你定义的对象的方法上有同步锁，但在运行时，却只有一个线程在访问，此时逃逸分析后的机器码，会去掉同步锁运行



#### go消除了 堆和栈的区别

一定程度上消除了区别，因为go在编译的时候进行逃逸分析，来决定一个对象放在栈上还是放在堆上，不逃逸的对象放在栈上，逃逸的放在堆上。



#### ！判断是否需要逃逸的标准？

一个对象，对这个对象所有的引用，都未超出当前函数以外（生命周期仅在当前函数内），则无需堆栈分配，仅分配栈上即可，无需逃逸；如果变量 或 变量的指针在函数内声明，但是在函数外被引用，说明可能被函数外其他程序访问，所以变量在此时会逃逸，被分配在堆上。



Done。

