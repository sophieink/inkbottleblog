[title:: Go]
(date:: 2026-10-12)  
(author:: ink)  
## Go?
Go(Golang)是谷歌(Google)开发的静态强类型、编译型、并发型，自带垃圾回收的编译语言。
## Go结构
- Package Declaration (命名包声明) 
- Import Package (引入包)
- Functions(函数)
- Variables(变量)
- Statements and Expressions  (表达式)
- Comments (注释)
***
## demo
go -> package + import + func
- comment 
  - //单line  
  - /*  
    多line  
    */  
- var 变量
- const 常量
- output
- datatypes 数据类型
  - bool
  - number >> int float complex
  - string
  - derived
    - [ ] pointer 指针
    - [x] array 数组
    - [x] structure 结构体
    - [x] function 函数
    - [x] slice 切片
    - [ ] interface 接口
    - [x] map 映射
    - [ ] channel 通道
- operator 运算符
- if 条件判断
- for 循环
- switch 
  - continue
  - break
- range 遍历
  - return 返回
***
### int-type
- int
    - uint8 (0-255)
    - uint16 (0-65535)
    - uint32 (0-4294967295)
    - uint64 (0-18446744073709551615)
    - int8 (-128-127)
    - int16 (-32768-32767)
    - int32 int64
- float
    - float32 + float64
- complex
    - complex64 + complex128
- other
    - byte == uint8
    - rune == int32
    - uint == 32/64 == int
    - uintptr == 无符号整数，用于储存指针值的uniterpreted位
***
### var
- var ?
  - `var v_list go_data_type` `v_list`多变量，`go_data_type`go支持的数据类型
```go
var i, j, k   int;
var c, ch byte;

//具体的值和类型>>静态类型变量声明
var a int;
a = 3;
//只有值，让编译器自识别>>动态类型变量声明
b := 5; // 在`func`内部支持，但外部要使用`b = 5`
// 多变量支持
var e, f, h = 7, 9, "hi"
```
- var type
  - byte 字节
  - int 整数
  - float32 单精度浮点值
***
### const
> 程序在执行过程中不得更改的固定值
> `const v int`
#### 转义符
|转义序列|含义|
|:---:|:---:|
|`\\`|`\`|
|`\'`|'|
|`\?`|?|
|`\b`|space|
|`\f`|换页page|
|`\n`|换行|
|`\r`|回车enter|

