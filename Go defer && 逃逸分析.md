# Go defer && 逃逸分析 

> -- 预留模板 可去除 --
>
> 视频连接：b站/油管链接
>
> 视频中PPT链接：
>
> 视频中代码仓库地址：
>
> 作者介绍：
>
> 博客地址/github地址：
>
> 宣传推广语：
>
> 总结语/概述语：通过本视频/本文，你可以收获到xxx
>
> 公众号（如有，公众号ID 或 QR code）



#### 个人介绍

来自罗胖——逻辑思维，适合在第一次场合介绍自己



分四部分

- 个人简介：我是谁，现在在哪里，背景条件，公众号、github、开源贡献
- 因缘（为什么来、如何了解到这个平台）
- 我能提供的：能力、近期研究、历史成果/经验
- 我希望获得：工作机会、粉丝、公众号打赏

----

#### Go defer

> 什么是 defer

- defer是 Go 语言的一种用于注册延迟调用的机制，使得函数或语句可以在当前函数执行完毕后执行。



> 为什么需要 defer

- Go提供的语法糖，减少资源泄露的发生。



> 如何使用 defer

- 在创建资源语句的附件，使用  defer 语句释放资源。



#### 代码示例

##### defer-code-1.go

```go
package main
import (
	"fmt"
)

/*func f1() (r int) {
  t := 5
  
  // 1.赋值指令
  r = t
  
  // 2. defer被插入到赋值与返回之间执行，这个例子中返回值 t 没有被修改过
	func() {
    t = t + 5
  } ()
  
  // 3. 空的 return 指令
  return 
}*/

func f1() (r int) {
  t := 5
  defer func() {
    t = t + 5
  } ()
  return t
}

func main() {
  fmt.Println(f1())
}

运行结果：
➜  golang夜读 git:(master) ✗ go run defer-code-1.go
5
```

在汇编中。vx2指令把值拷贝到栈上，然后一个空的ret指令返回，**defer指令** 实际生效 就在 赋值指令 和 return 指令 之间。



##### defer-code-2.go

```go
package main
import (
	"fmt"
)

/*func f2() (r int) {
  // 1. 赋值
  r = 1
  
  // 2. 这里改的r是之前传值进去的局部变量r，不会改变要返回的那个r值
	func(r int) {
    r = r + 5
  } (r)
  
  // 3. 空的 return 指令
  return 
}*/

func f2() (r int) {
  defer func(r int) {
    r = r + 5
  } (r)
  return 1
}

func main() {
  fmt.Println(f2())
}

运行结果：
➜  golang夜读 git:(master) ✗ go run defer-code-2.go
1
```

使用局部变量，未使用指针，值未修改



##### defer-code-3.go

```go
package main
import (
	"fmt"
)

/*func f2() (r int) {
  // 1. 赋值
  r = 1
  
  // 2. 这里改的r是之前传引用传进去的r，会改变要返回的那个r值
	func(r int) {
    r = r + 5
  } (r)
  
  // 3. 空的 return 指令
  return 
}*/

func f3() (r int) {
  defer func(r *int) {
    *r = *r + 5
  } (&r)
  return 1
}

func main() {
  fmt.Println(f3())
}

运行结果：
➜  golang夜读 git:(master) ✗ go run defer-code-3.go
6
```



##### defer-code-4.go

```go
package main
import (
  "fmt"
  "errors"
)

func e1() {
  var err error
  
  defer fmt.Println(err)
  
  err = errors.New("defer1 error")
  return
}

func e2() {
  var err error
  
  defer func() {
    fmt.Println(err)
  }()
  
  err = errors.New("defer2 error")
  return
}

func e3() {
  var err error
  
  defer func(err error) {
    fmt.Println(err)
  }(err)
  
  err = errors.New("defer3 error")
  return
}

func main() {
  e1()
  e2()
  e3()
}

运行结果：
➜  golang夜读 git:(master) ✗ go run defer-code-4.go
<nil>
defer2 error
<nil>
```



defer 语句：先入后出，后进先出。

为什么输出 “defer2 error”？golang 闭包函数捕获了 上文中的 err 变量引用



##### defer-code-5.go

```go
package main

import (
  "fmt"
)

/*
闭包引用了 x 变量，a,b 可以看作2个不同的实例，实例之间互相不影响。实例内部，x变量是同一个地址，因此具有”累加效应“
*/

func main() {
  var a = Accumulator()
  
  fmt.Printf("%d\n", a(1))
  fmt.Printf("%d\n", a(10))
  fmt.Printf("%d\n", a(100))
  
  fmt.Println("%-----------------")
  var b = Accumulator()
  
  fmt.Printf("%d\n", b(1))
  fmt.Printf("%d\n", b(10))
  fmt.Printf("%d\n", b(100))
  
}

func Accumulator() func(int) int {
  var x int
  
  return func(delta int) int {
    fmt.Printf("(%+v, %+v) - ", &x, x)
    x += delta
    return x
  }
}

运行结果：
➜  golang夜读 git:(master) ✗ go run defer-code-5.go
(0xc00001a0b0, 0) - 1
(0xc00001a0b0, 1) - 11
(0xc00001a0b0, 11) - 111
%-----------------
(0xc00001a0e8, 0) - 1
(0xc00001a0e8, 1) - 11
(0xc00001a0e8, 11) - 111

可以观察到 一组 相同函数内，x变量已放到堆上，所以 x变量 指针地址不变（使用同一个x变量），一直处于累加状态
```

​	

##### defer-code-6.go

```go
package main

import (
  "fmt"
  "time"
)

func main() {
  fmt.Println("The start of main function.")
  defer fmt.Println("defer main")
  var user = ""
  
  go func() {
    defer func() {
      fmt.Println("defer for goroutine caller.")
      if err := recover(); err !=nil {
        fmt.Println("recover success. err:", err)
      }
    }()
    
    func() {
      defer func() {
        fmt.Println("panic occurs here. then we have panic defer.")
      }()
      
      if user == "" {
        panic("should set user env.")
      }
      
      //此处不会执行（终止于 panic）
      fmt.Println("after panic will be here~~~")
    }()
  }()
  
  time.Sleep(1000)
  
  fmt.Println("The end of main function.")
}


运行结果：（由于运行快慢不同，运行多次，效果不同）
（运行过快，gc完成）
➜  golang夜读 git:(master) ✗ go run defer-code-6.go
The start of main function.
The end of main function.
defer main
（下面才是希望能看到的）
➜  golang夜读 git:(master) ✗ go run defer-code-6.go
The start of main function.
The end of main function.
defer main
panic occurs here. then we have panic defer.
defer for goroutine caller.
recover success. err: should set user env.
```



通过上面的例子，可以了解 defer、goroutine、panic、main 组合使用过程中，各个函数执行阶段。（和kubernetes类似，对程序的把控力和对 分布式系统/容器把控过程，都是了解清楚系统原理中各阶段生命周期的顺序、系统/人为行为、限制，然后基于这些已有条件进行优化、改造和扩展）



main进程创建的 goroutine 线程发生panic，main进程无法recover 该 goroutine 线程的panic。所以每有个匿名函数，最好有一个对应的defer函数来recover对应的panic



----

##### defer-code-7.go

```go
package main

import (
  "fmt"
)

func main() {
  for i :=0; i<5; i++ {
    defer fmt.Println(i, 1)
  }
  
  for i :=0; i<5; i++ {
    defer func() {
      fmt.Println(i, 2)
    }()
  }
  
  for i :=0; i<5; i++ {
    defer func() {
      j := i
      fmt.Println(j, 3)
      //fmt.Println(j, &j, i, &i, 3)
      //5 0xc00008e020 5 0xc000086028 3
      //5 0xc000086078 5 0xc000086028 3
      //5 0xc000086090 5 0xc000086028 3
      //5 0xc00008e038 5 0xc000086028 3
      //5 0xc0000860a8 5 0xc000086028 3
    }()
  }
  
  for i :=0; i<5; i++ {
    j := i
    defer fmt.Println(j, 4)
  }
  
  for i :=0; i<5; i++ {
    j := i
    defer func() {
      fmt.Println(j, 5)
    }()
  }
  
  for i :=0; i<5; i++ {
    defer func(j int) {
      fmt.Println(j, 6)
    }(i)
  }
}

运行结果：
➜  golang夜读 git:(master) ✗ go run defer-code-7.go
4 6
3 6
2 6
1 6
0 6
4 5
3 5
2 5
1 5
0 5
4 4
3 4
2 4
1 4
0 4
5 3
5 3
5 3
5 3
5 3
5 2
5 2
5 2
5 2
5 2
4 1
3 1
2 1
1 1
0 1
```

----

#### Go 逃逸分析

> 什么是逃逸分析

- Go语言编译器执行静态代码分析后，决定哪些变量逃逸到堆上。（决定程序运行时的变量内存分配，只有编译器能证明变量在函数返回后就没有引用的，才会存在栈上，保证函数调用返回被销毁即可）



> 为什么需要 逃逸分析

- 尽可能地将变量分配到栈上。



> 逃逸分析如何进行

- 只有在编译器能证明变量在函数返回后不再被引用，才能分配到栈上，其他情况分配到堆上。



堆——动态分配内存，栈——静态分配内存。动态分配开销 > 静态分配。

举例:

new 动态分配，遍历分配器的内存，会涉及内存碎片、涉及内存垃圾回收；静态分配，两个CPU指令即可



#### 代码演示

##### escape-code-1.go

```go
/*
go build -gcflags '-m -l' escape-code-1.go

没有逃逸

值传递，直接在栈上分配

GO语言函数传递都是通过值的，调用函数时，直接在栈上COPY出一份参数，不存在逃逸
*/
package main

type T1 struct {}

func main() {
  var x T1
  _ = identity1(x)
}

func identity1(x T1) T1 {
  return x
}

运行结果（无输出，即没有逃逸）
➜  golang夜读 git:(master) ✗ go build -gcflags '-m -l' escape-code-1.go
➜  golang夜读 git:(master) ✗
```



##### escape-code-2.go

```go
/*
go build -gcflags '-m -l' escape-code-2.go

x 没有发生逃逸

identity函数的输入直接当成返回值了，因为没有对 z 作引用，所以z没有逃逸.
对x的引用也没有逃出 main 函数的作用域，因此 x 也没有发生逃逸。
*/
package main

type T2 struct {}

func main() {
  var x T2
  y := &x
  _ = *identity1(y)
}

func identity1(z *T2) *T2 {
  return z
}

运行结果：
➜  golang夜读 git:(master) ✗ go build -gcflags '-m -l' escape-code-2.go
# command-line-arguments
./escape-code-2.go:11:16: leaking param: z to result ~r1 level=0
./escape-code-2.go:7:7: main &x does not escape
```



##### escape-code-3.go

```go
/*
go build -gcflags '-m -l' escape-code-3.go

z 发生逃逸

z 是对 x 的拷贝， ref3 函数中对 z 取了引用，所以 z 不能放在栈上。否则在 ref3 函数之外，通过引用不能找到 z，所以z必须要逃逸到堆上。尽管在main函数中，直接丢弃了 ref3 的结果，但是 Go编译器还没有那么智能，分析不出这种情况。
而对 x 从来没有取引用，所以 x 不会发生逃逸。
*/
package main

type T3 struct {}

func main() {
  var x T3
  _ = *ref3(x)
}

func ref3(z T3) *T3 {
  return &z
}

运行结果：
➜  golang夜读 git:(master) ✗ go build -gcflags='-m -l' escape-code-3.go
# command-line-arguments
./escape-code-3.go:11:9: &z escapes to heap
./escape-code-3.go:10:11: moved to heap: z
```



##### escape-code-4.go

```go
/*
go build -gcflags '-m -l' escape-code-4.go

y 发生逃逸

ref4 函数对 y 取了引用，所以 y 发生了逃逸。
*/
package main

type T4 struct {
  M *int
}

func main() {
	var i int
  ref4(i)
}

func ref4(y int) (z T4){
  z.M = &y
  return z
}

运行结果
➜  golang夜读 git:(master) ✗ go build -gcflags '-m -l' escape-code-4.go
# command-line-arguments
./escape-code-4.go:13:8: &y escapes to heap
./escape-code-4.go:12:11: moved to heap: y
```



##### escape-code-5.go

```go
/*
go build -gcflags '-m -l' escape-code-5.go

i 没有发生逃逸

在main函数里对 i 取了引用，把它传给ref4 函数，i的引用一直在 main 函数的作用域里，因此 i 没有发生逃逸。

和上个代码例子相比，i的写法有了点小差别，但是i的命运是不同的，导致程序运行效果不同：例子4中，变量i先在main的栈帧中分配，之后又在ref3函数栈帧中分配，然后又逃逸到堆上，到堆上分配了一次，共三次分配。
本例中，i 只分配了一次，然后通过引用传递。
*/
package main

type T5 struct {
  M *int
}

func main() {
	var i int
  ref5(&i)
}

func ref5(y *int) (z T5){
  z.M = y
  return z
}

运行结果：
➜  golang夜读 git:(master) ✗ go build -gcflags '-m -l' escape-code-5.go
# command-line-arguments
./escape-code-5.go:12:11: leaking param: y to result z level=0
./escape-code-5.go:9:7: main &i does not escape
```



##### escape-code-6.go

```go
/*
go build -gcflags '-m -l' escape-code-6.go

i 逃逸，x未逃逸。x的作用域一直在main函数中

本例 i 发生了逃逸，按照前面例子5的分析，i不会发生逃逸。但是此例y赋值给了一个输入的struct，go不能从输出反推到输入。
两个例子的区别是例子5中的T5是在返回值里的，输入只能”流入“到输出，
本例中S6是在输入参数中，所以逃逸分析失败，i要逃逸到堆上。 
*/
package main

type T6 struct {
  M *int
}

func main() {
  var x T6
	var i int
  ref6(&i, &x)
}

func ref6(y *int, z *T6) {
  z.M = y
}

运行结果：
➜  golang夜读 git:(master) ✗ go build -gcflags '-m -l' escape-code-6.go
# command-line-arguments
./escape-code-6.go:13:11: leaking param: y
./escape-code-6.go:13:19: ref6 z does not escape
./escape-code-6.go:10:7: &i escapes to heap
./escape-code-6.go:9:6: moved to heap: i
./escape-code-6.go:10:11: main &x does not escape
```



----

汇编查看 defer 的实现原理



```go
/*
compile: go build -o test asm.go
check assembly: go tool objdump -s "main\.main" test
*/

package main

func main() {
  defer println(0x11)
}

运行结果：
➜  golang夜读 git:(master) ✗ go build -o test asm.go
➜  golang夜读 git:(master) ✗ go tool objdump -s "main\.main" test
TEXT main.main(SB) /Users/pauljobs/Documents/golang夜读/asm.go
  asm.go:3		0x104ea20		65488b0c2530000000	MOVQ GS:0x30, CX
  asm.go:3		0x104ea29		483b6110		CMPQ 0x10(CX), SP
  asm.go:3		0x104ea2d		7653			JBE 0x104ea82
  asm.go:3		0x104ea2f		4883ec20		SUBQ $0x20, SP
  asm.go:3		0x104ea33		48896c2418		MOVQ BP, 0x18(SP)
  asm.go:3		0x104ea38		488d6c2418		LEAQ 0x18(SP), BP
  asm.go:4		0x104ea3d		c7042408000000		MOVL $0x8, 0(SP)
  asm.go:4		0x104ea44		488d05d5290200		LEAQ go.func.*+61(SB), AX
  asm.go:4		0x104ea4b		4889442408		MOVQ AX, 0x8(SP)
  asm.go:4		0x104ea50		48c744241011000000	MOVQ $0x11, 0x10(SP)
  asm.go:4		0x104ea59		e8c231fdff		CALL runtime.deferproc(SB)
  asm.go:4		0x104ea5e		85c0			TESTL AX, AX
  asm.go:4		0x104ea60		7510			JNE 0x104ea72
  asm.go:5		0x104ea62		90			NOPL
  asm.go:5		0x104ea63		e8483afdff		CALL runtime.deferreturn(SB)
  asm.go:5		0x104ea68		488b6c2418		MOVQ 0x18(SP), BP
  asm.go:5		0x104ea6d		4883c420		ADDQ $0x20, SP
  asm.go:5		0x104ea71		c3			RET
  asm.go:4		0x104ea72		90			NOPL
  asm.go:4		0x104ea73		e8383afdff		CALL runtime.deferreturn(SB)
  asm.go:4		0x104ea78		488b6c2418		MOVQ 0x18(SP), BP
  asm.go:4		0x104ea7d		4883c420		ADDQ $0x20, SP
  asm.go:4		0x104ea81		c3			RET
  asm.go:3		0x104ea82		e8d984ffff		CALL runtime.morestack_noctxt(SB)
  asm.go:3		0x104ea87		eb97			JMP main.main(SB)
  :-1			0x104ea89		cc			INT $0x3
  :-1			0x104ea8a		cc			INT $0x3
  :-1			0x104ea8b		cc			INT $0x3
  :-1			0x104ea8c		cc			INT $0x3
  :-1			0x104ea8d		cc			INT $0x3
  :-1			0x104ea8e		cc			INT $0x3
  :-1			0x104ea8f		cc			INT $0x3
```



defer调用返回——deferreturn

源代码：

https://github.com/golang/go/blob/2517f4946b42b8deedb864c884f1b41311d45850/src/runtime/panic.go#L524

```go
// Run a deferred function if there is one.
// The compiler inserts a call to this at the end of any
// function which calls defer.
// If there is a deferred function, this will call runtime·jmpdefer,
// which will jump to the deferred function such that it appears
// to have been called by the caller of deferreturn at the point
// just before deferreturn was called. The effect is that deferreturn
// is called again and again until there are no more deferred functions.
//
// Declared as nosplit, because the function should not be preempted once we start
// modifying the caller's frame in order to reuse the frame to call the deferred
// function.
//
// The single argument isn't actually used - it just has its address
// taken so it can be matched against pending defers.
//go:nosplit
func deferreturn(arg0 uintptr) {
	gp := getg()
	d := gp._defer
	if d == nil {
		return
	}
	sp := getcallersp()
	if d.sp != sp {
		return
	}
	if d.openDefer {
		done := runOpenDeferFrame(gp, d)
		if !done {
			throw("unfinished open-coded defers in deferreturn")
		}
		gp._defer = d.link
		freedefer(d)
		return
	}

	// Moving arguments around.
	//
	// Everything called after this point must be recursively
	// nosplit because the garbage collector won't know the form
	// of the arguments until the jmpdefer can flip the PC over to
	// fn.
	switch d.siz {
	case 0:
		// Do nothing.
	case sys.PtrSize:
		*(*uintptr)(unsafe.Pointer(&arg0)) = *(*uintptr)(deferArgs(d))
	default:
		memmove(unsafe.Pointer(&arg0), deferArgs(d), uintptr(d.siz))
	}
	fn := d.fn
	d.fn = nil
	gp._defer = d.link
	freedefer(d)
	// If the defer function pointer is nil, force the seg fault to happen
	// here rather than in jmpdefer. gentraceback() throws an error if it is
	// called with a callback on an LR architecture and jmpdefer is on the
	// stack, because the stack trace can be incorrect in that case - see
	// issue #8153).
	_ = fn.fn
	jmpdefer(fn, uintptr(unsafe.Pointer(&arg0)))
}
```



deferproc ——defer函数调用入口

源代码：

https://github.com/golang/go/blob/2517f4946b42b8deedb864c884f1b41311d45850/src/runtime/panic.go#L223

```go
// Create a new deferred function fn with siz bytes of arguments.
// The compiler turns a defer statement into a call to this.
//go:nosplit
func deferproc(siz int32, fn *funcval) { // arguments of fn follow fn
	gp := getg()
	if gp.m.curg != gp {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}

	// the arguments of fn are in a perilous state. The stack map
	// for deferproc does not describe them. So we can't let garbage
	// collection or stack copying trigger until we've copied them out
	// to somewhere safe. The memmove below does that.
	// Until the copy completes, we can only call nosplit routines.
	sp := getcallersp()
	argp := uintptr(unsafe.Pointer(&fn)) + unsafe.Sizeof(fn)
	callerpc := getcallerpc()

	d := newdefer(siz)
	if d._panic != nil {
		throw("deferproc: d.panic != nil after newdefer")
	}
	d.link = gp._defer
	gp._defer = d
	d.fn = fn
	d.pc = callerpc
	d.sp = sp
	switch siz {
	case 0:
		// Do nothing.
	case sys.PtrSize:
		*(*uintptr)(deferArgs(d)) = *(*uintptr)(unsafe.Pointer(argp))
	default:
		memmove(deferArgs(d), unsafe.Pointer(argp), uintptr(siz))
	}

	// deferproc returns 0 normally.
	// a deferred func that stops a panic
	// makes the deferproc return 1.
	// the code the compiler generates always
	// checks the return value and jumps to the
	// end of the function if deferproc returns != 0.
	return0()
	// No code can go here - the C return register has
	// been set and must not be clobbered.
}

// ! 本人额外标注，与defer函数相关栈
// deferprocStack queues a new deferred function with a defer record on the stack.
// The defer record must have its siz and fn fields initialized.
// All other fields can contain junk.
// The defer record must be immediately followed in memory by
// the arguments of the defer.
// Nosplit because the arguments on the stack won't be scanned
// until the defer record is spliced into the gp._defer list.
//go:nosplit
func deferprocStack(d *_defer) {
	gp := getg()
	if gp.m.curg != gp {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}
	// siz and fn are already set.
	// The other fields are junk on entry to deferprocStack and
	// are initialized here.
	d.started = false
	d.heap = false
	d.openDefer = false
	d.sp = getcallersp()
	d.pc = getcallerpc()
	d.framepc = 0
	d.varp = 0
	// The lines below implement:
	//   d.panic = nil
	//   d.fd = nil
	//   d.link = gp._defer
	//   gp._defer = d
	// But without write barriers. The first three are writes to
	// the stack so they don't need a write barrier, and furthermore
	// are to uninitialized memory, so they must not use a write barrier.
	// The fourth write does not require a write barrier because we
	// explicitly mark all the defer structures, so we don't need to
	// keep track of pointers to them with a write barrier.
	*(*uintptr)(unsafe.Pointer(&d._panic)) = 0
	*(*uintptr)(unsafe.Pointer(&d.fd)) = 0
	*(*uintptr)(unsafe.Pointer(&d.link)) = uintptr(unsafe.Pointer(gp._defer))
	*(*uintptr)(unsafe.Pointer(&gp._defer)) = uintptr(unsafe.Pointer(d))

	return0()
	// No code can go here - the C return register has
	// been set and must not be clobbered.
}
```





deferproc相关介绍

源代码地址：

https://github.com/golang/go/blob/28e549dec3954b36d0c83442be913d8709d7e5ae/src/runtime/stack.go#L48

```go
The bottom StackGuard - StackSmall bytes are important: there has
to be enough room to execute functions that refuse to check for
stack overflow, either because they need to be adjacent to the
actual caller's frame (deferproc) or because they handle the imminent
stack overflow (morestack).
For example, deferproc might call malloc, which does one of the
above checks (without allocating a full frame), which might trigger
a call to morestack.  This sequence needs to fit in the bottom
section of the stack.  On amd64, morestack's frame is 40 bytes, and
deferproc's frame is 56 bytes.  That fits well within the
StackGuard - StackSmall bytes at the bottom.
```



同文件（src/runtime/stack.go）下代码

GoExit（goroutine退出调用函数）

https://github.com/golang/go/blob/2517f4946b42b8deedb864c884f1b41311d45850/src/runtime/panic.go#L579

```go
// Goexit terminates the goroutine that calls it. No other goroutine is affected.
// Goexit runs all deferred calls before terminating the goroutine. Because Goexit
// is not a panic, any recover calls in those deferred functions will return nil.
//
// Calling Goexit from the main goroutine terminates that goroutine
// without func main returning. Since func main has not returned,
// the program continues execution of other goroutines.
// If all other goroutines exit, the program crashes.
func Goexit() {
	// Run all deferred functions for the current goroutine.
	// This code is similar to gopanic, see that implementation
	// for detailed comments.
	gp := getg()

	// Create a panic object for Goexit, so we can recognize when it might be
	// bypassed by a recover().
	var p _panic
	p.goexit = true
	p.link = gp._panic
	gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

	addOneOpenDeferFrame(gp, getcallerpc(), unsafe.Pointer(getcallersp()))
	for {
		d := gp._defer
		if d == nil {
			break
		}
		if d.started {
			if d._panic != nil {
				d._panic.aborted = true
				d._panic = nil
			}
			if !d.openDefer {
				d.fn = nil
				gp._defer = d.link
				freedefer(d)
				continue
			}
		}
		d.started = true
		d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))
		if d.openDefer {
			done := runOpenDeferFrame(gp, d)
			if !done {
				// We should always run all defers in the frame,
				// since there is no panic associated with this
				// defer that can be recovered.
				throw("unfinished open-coded defers in Goexit")
			}
			if p.aborted {
				// Since our current defer caused a panic and may
				// have been already freed, just restart scanning
				// for open-coded defers from this frame again.
				addOneOpenDeferFrame(gp, getcallerpc(), unsafe.Pointer(getcallersp()))
			} else {
				addOneOpenDeferFrame(gp, 0, nil)
			}
		} else {

			// Save the pc/sp in reflectcallSave(), so we can "recover" back to this
			// loop if necessary.
			reflectcallSave(&p, unsafe.Pointer(d.fn), deferArgs(d), uint32(d.siz))
		}
		if p.aborted {
			// We had a recursive panic in the defer d we started, and
			// then did a recover in a defer that was further down the
			// defer chain than d. In the case of an outstanding Goexit,
			// we force the recover to return back to this loop. d will
			// have already been freed if completed, so just continue
			// immediately to the next defer on the chain.
			p.aborted = false
			continue
		}
		if gp._defer != d {
			throw("bad defer entry in Goexit")
		}
		d._panic = nil
		d.fn = nil
		gp._defer = d.link
		freedefer(d)
		// Note: we ignore recovers here because Goexit isn't a panic
	}
	goexit1()
}

```



----

### 总结

> 第一次自我介绍，四点法

- 基本信息（硬信息，如个人信息）
- 因缘
- 你能提供的
- 你希望收获的



> defer

- 每次 defer 语句执行的时候，会把函数”压栈“，函数参数会被拷贝下来；
- 当外层函数（非代码块，如一个 for 循环）退出时，defer函数按照定义的逆序执行；
- 如果defer执行的函数为nil，那么会在最终调用函数调用的产生panic；



> 逃逸分析

- 动态内存分配（堆上）比静态内存分配（栈上）开销大得多（可以后续个人研究比较）；
- 如果变量在函数外部没有引用，则优先放到栈中；如果在函数外部存在引用，则必定放到堆中；

