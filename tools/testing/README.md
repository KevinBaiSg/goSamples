
## 1. Go 原生测试工具

Go 语言在设计之初就非常重视程序设计，并且在官方的包中提供了非常丰富的测试工具。配合 Go 的 struct 可以方便的实现
表驱动测试。

### 1.1 约束  
Go 对测试函数名称和签名的约束：
* 功能测试函数: 名称必须以Test为前缀，且参数列表中只应有一个*testing.T类型的参数声明.    
`func TestXxx(*testing.T)`  
* 基准测试函数: 名称必须以Benchmark为前缀，并且唯一参数的类型必须是*testing.B类型的。  
`func BenchmarkXxx(*testing.B)`     
* 示例测试函数: 名称必须以 Example 为前缀，但对函数的参数列表没有强制规定。  

一般情况下我们会把实现文件和测试代码文件放在同一个目录下，并且文件名在结尾处增加 `_test`。
例如本利中实现文件 `zigZag.go` 是我编写的 LeetCode 问题 `6` 的实现，其测试文件为 `zigZig_test.go`。

```         
The file will be excluded from regular package builds but will be included when the “go test” command is run
For more detail, run “go help test” and “go help testflag”.
```     

### 1.2 单元测试    
go 的 testing 包可以基于包进测试，意思是执行 `go test` 默认单位是包范围，运行`go test` 可以自动执行当前包下面的所有
`func TestXxx(*testing.T)` 格式的测试方法。（`Xxx` 可以是任意字母，但最好和需要测试的方法对应)

例如本例： `func TestZigZagConversion(t *testing.T)` 

单元测试中，传递给测试函数的参数是 `*testing.T` 类型。它用于管理测试状态并支持格式化测试日志。
测试日志会在执行测试的过程中不断累积，并在测试完成时转储至标准输出。

如果代码支持并行运行，可以在测试代码中加上 `t.Parallel()`, 也可以在运行测试时加上 `-race` 检查。

#### 1.2.1 表驱动测试
文章支持就提到过这种测试方法，其实就是很好的利用了 Go 的 struct 特性。

```go
cases := []struct {
    s 		string
    numRows int
    want 	string
}{
    {"", 10, ""},
    {"LEETCODEISHIRING", 0, ""},
    {"LEETCODEISHIRING", 1, "LEETCODEISHIRING"},
    {"LEETCODEISHIRING", 3, "LCIRETOESIIGEDHN"},
    {"LEETCODEISHIRING", 4, "LDREOEIIECIHNTSG"},
}

for _, c := range cases {
    zigZag := convert(c.s, c.numRows)
    if zigZag != c.want {
        t.Errorf("Reverse(%q) == %q, want %q", c.s, zigZag, c.want)
    }
}
```

#### 1.2.2 常见的测试日志输出对比  

带 f 的是格式化的，格式化语法参考 fmt 包    

|    API   | case 是否失败 | case 是否中断 | 说明 |
|----------|:-------------:|:------:|:------:|
| Fail | 是 | 否 | | 
| FailNow | 是 | 是 | 内部通过调用 Goexit 中断测试 |
| SkipNow | 是 | 是 | 只跳过测试，不标识错误，内部通过调用 Goexit 中断测试 |
| Log/Logf | 是 | 否 | 输出信息/输出格式化的信息 |
| Skip/Skipf | 是 | 是 | Log/Logf + SkipNow |
| Error/Error | 是 | 否 | Log/Logf + Fail |
| Fatal/Fatalf | 是 | 否 | Log/Logf + FailNow |

### 1.3 基准测试
本例中的基准测试：    
`func BenchmarkZigZagConversion(b *testing.B)`  

基准测试以Benchmark为前缀，必须要执行 b.N 次，b.N 的值是系统根据实际情况去调整的，从而保证测试的稳定性。

如果在测试代码运行之前有很长时间的初始化工作，可以在测试代码之前增加 `b.ResetTimer()`，在这之前的处理不会放到执行时间里，
也不会输出到报告中。除了 `ResetTimer` 还有两个 api 可以管理 timer：`StartTimer` 和 `StopTimer`。 `StartTimer` 在
基准测试开始时就会自动被调用，当然我们可以在 `StopTimer` 被调用后手动调用。

基准测试并不会默认执行，他需要增加 `-bench` 参数，例如 `go test -bench=.` 

基准测试也可以开启并行测试，需要执行 `b.RunParallel(func(pb *testing.PB)` 方法，默认会以逻辑 CPU 个数来进行并行测试。
当然也可以通过指令来指定 cpu 个数 `go test -bench=. -cpu 1`    

#### 1.3.1 测试结果 
```bash
goos: darwin
goarch: amd64
pkg: github.com/KevinBaiSg/goSamples/tools/testing
BenchmarkZigZagConversion-8      2000000               720 ns/op
PASS
ok      github.com/KevinBaiSg/goSamples/tools/testing   2.184s
```
测试结果说明
BenchmarkZigZagConversion-8： 中的 8 最大 P 数量为8。
2000000： 代表测试次数
720 ns/op： 代表 720 ns 每 loop

#### 1.3.2 并行测试
通过 RunParallel 方法以并行的方式执行给定的基准测试。Run Parallel 会创建出多个 goroutine，
并将 b.N 分配给这些 goroutine 执行，其中 goroutine 数量的默认值为 GOMAXPROCS。
用户如果想要增加非CPU受限（non-CPU-bound）基准测试的并行性，那么可以在 RunParallel 
之前调用 SetParallelism（如 SetParallelism(2)，则 goroutine 数量为 2*GOMAXPROCS）。
RunParallel 通常会与 -cpu 标志一同使用。

body 函数将在每个 goroutine 中执行，这个函数需要设置所有 goroutine 本地的状态，
并迭代直到 pb.Next 返回 false 值为止。因为 StartTimer、StopTime 和 ResetTimer 
这三个方法都带有全局作用，所以 body 函数不应该调用这些方法； 除此之外，body 函数也不应该调用 Run 方法。

```go
func BenchmarkZigZagConversionRunParallel(b *testing.B)  {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			convert("LEETCODEISHIRING", 4)
		}
	})
}
```

#### 1.3.3 内存统计
`ReportAllocs` 方法用于打开当前基准测试的内存统计功能， 与 `go test` 使用 `-benchmem` 标志类似，
但 `ReportAllocs` 只影响那些调用了该函数的基准测试。

```go
func BenchmarkZigZagConversion(b *testing.B)  {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		convert("LEETCODEISHIRING", 4)
	}
}
```

### 1.4 示例验证
示例测试既可以被当做文档来使用，也可以当做测试来运行（可以不用为了测试代码写个 `main() {...}` 了）.

#### 1.4.1 输出
```go
func ExampleZigZagConversion() {
	convert("LEETCODEISHIRING", 4)
	// Output: LDREOEIIECIHNTSG
}
```
注意上例中的 `Output`，测试代码在运行后会与 `Output` 后边的值进行比较。
有时候在测试中我们没办法判断输出顺序，这时候可以用 `Unordered output` 代替。
其是这时候也注意到示例验证有两个要求，一是函数以 `Example` 开头，二是代码中包含 `Output` 或 `Unordered output` 的注释。
如果没有包含`Output` 或 `Unordered output`，测试代码只会被编译，不会运行。

#### 1.4.2 命名约定
约定中符号：函数 F， 类型 T，类型 T 中的 方法 M
```go
func Example() { ... }
func ExampleF() { ... }
func ExampleT() { ... }
func ExampleT_M() { ... }
```

有时，我们想要给 包/类型/函数/方法 提供多个示例，可以通过在示例函数名称后附加一个不同的后缀来实现，但这种后缀必须以小写字母开头，如：
```go
func Example_suffix() { ... }
func ExampleF_suffix() { ... }
func ExampleT_suffix() { ... }
func ExampleT_M_suffix() { ... }
```
除了创建用于测试的示例外，示例还用于显示在生成的文档中。

### 1.5 子测试和子基准测试
还可以通过 `T` 和 `B` 的 Run 方法开启子测试和子基准测试，主要是可以共享公共的设置和资源清除的管理。

每个子测试都有一个唯一的名字，以父测试用/隔开来唯一表示，运行的时候使用 `-run regexp`
指定测试和 `-bench regexp` 来指定基准测试，`.` 表示所有。

下面是一个模板
```go
func TestFoo(t *testing.T) {
    // <setup code>
    t.Run("A=1", func(t *testing.T) { ... })
    t.Run("A=2", func(t *testing.T) { ... })
    t.Run("B=1", func(t *testing.T) { ... })
    // <tear-down code>
}
```

下面是测试指令的模板
```bash
go test -run ''      # Run 所有测试。
go test -run Foo     # Run 匹配 "Foo" 的顶层测试，例如 "TestFooBar"。
go test -run Foo/A=  # 匹配顶层测试 "Foo"，运行其匹配 "A=" 的子测试。
go test -run /A=1    # 运行所有匹配 "A=1" 的子测试。
```

在并行子测试完成之前，Run 方法不会返回。

### 1.6 Main 测试
在写测试时，有时需要在测试之前或之后进行额外的设置（setup）或拆卸（teardown）；
有时，测试还需要控制在主线程上运行的代码。为了支持这些需求，testing 提供了 TestMain 函数:
`func TestMain(m *testing.M)`   

如果测试文件中包含该函数，那么生成的测试将调用 `TestMain(m)`，而不是直接运行测试。
TestMain 运行在主 goroutine 中, 可以在调用 m.Run 前后做任何设置和拆卸。
注意，在 TestMain 函数的最后，应该使用 m.Run 的返回值作为参数调用 `os.Exit`。

另外，在调用 TestMain 时, flag.Parse 并没有被调用。
所以，如果 TestMain 依赖于 command-line 标志 (包括 testing 包的标记), 
则应该显示的调用 flag.Parse。注意，这里说的依赖，说的是如果 TestMain 函数内依赖 flag，
则必须显示调用 flag.Parse，否则不需要，因为 m.Run 中调用 flag.Parse。

一个包含 TestMain 的例子如下：
```go
package mytestmain

import (  
    "flag"
    "fmt"
    "os"
    "testing"
)

var db struct {  
    Dns string
}

func TestMain(m *testing.M) {
    db.Dns = os.Getenv("DATABASE_DNS")
    if db.Dns == "" {
        db.Dns = "root:123456@tcp(localhost:3306)/?charset=utf8&parseTime=True&loc=Local"
    }

    flag.Parse()
    exitCode := m.Run()

    db.Dns = ""

    // 退出
    os.Exit(exitCode)
}

func TestDatabase(t *testing.T) {
    fmt.Println(db.Dns)
}
```

### 1.7 测试覆盖率
由单元测试的代码，触发运行到的被测试代码的代码行数占所有代码行数的比例，被称为测试覆盖率。
Go 为我们提供了测试覆盖率的相关的工具（`go test -cover` 和 `go tool cover`）。
下面介绍几种常用指令

#### 1.7.1 查看代码覆盖率
```bash
$go test -cover // 

PASS
coverage: 95.5% of statements
ok  	github.com/KevinBaiSg/goSamples/tools/testing	0.005s
``` 

#### 1.7.2 测试中指定代码覆盖率生成文件 c.out
```bash
$go test -v -coverprofile=c.out // 

=== RUN   TestZigZagConversion
--- PASS: TestZigZagConversion (0.00s)
=== RUN   ExampleZigZagConversion
--- PASS: ExampleZigZagConversion (0.00s)
PASS
coverage: 95.5% of statements
ok  	github.com/KevinBaiSg/goSamples/tools/testing	0.006s
``` 

#### 1.7.3 测试代码覆盖，并生成文件
```bash
$go tool cover -html=c.out -o=tag.html 
``` 


## 其他测试框架
### testify  
[github](https://github.com/stretchr/testify)   
[godoc](https://godoc.org/github.com/stretchr/testify/mock)  
[Mocking in Golang Using Testify](https://blog.lamida.org/mocking-in-golang-using-testify/)  

## Reference：  
[1. Package testing](https://golang.org/pkg/testing/)  
[2. Go 语言的几种测试方法](https://blog.csdn.net/lastsweetop/article/details/78469507?utm_campaign=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com)  
[3. 测试的基本规则和流程-极客时间](https://time.geekbang.org/column/article/41036)      
[4. Go 测试，go test 工具的具体指令 flag](https://deepzz.com/post/the-command-flag-of-go-test.html)   
[5. Go 语言中文网](https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter09/09.2.html)  
[6. Go 测试，go test 工具的具体指令 flag](https://deepzz.com/post/the-command-flag-of-go-test.html)  
