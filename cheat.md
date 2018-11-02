## 判断是否为空

https://blog.csdn.net/phantom_111/article/details/54670598

使用”“判断string变量是否为空。

使用nil判断结构体的指针是否为空。 

len(s) == 0用于求数组、切片和字典的长度。 

## 类型转换
### string to uint32
 
package main

import (
    "fmt"
    "strconv"
)

func main() {
    width := "42"
    u64, err := strconv.ParseUint(width, 10, 32)
    if err != nil {
        fmt.Println(err)
    }
    wd := uint(u64)
    fmt.Println(wd)
}


### uint32 to string

down vote
accepted
I would simply use Sprintf or even just Sprint:

var n uint32 = 42
str := fmt.Sprint(n)
println(str)
 
### string到int 
int,err:=strconv.Atoi(string)
### string到int64 
int64, err := strconv.ParseInt(string, 10, 64)
### int到string 
string:=strconv.Itoa(int)
### int64到string 
string:=strconv.FormatInt(int64,10)
 




