一、[go语言的http包](https://my.oschina.net/u/943306/blog/151293)

摘要：

```go
package main

import "net/http"

func main() {
    http.ListenAndServe(":8080", nil)
}
```

查询[ListenAndServe的文档](http://gowalker.org/net/http#ListenAndServe)可知，第2个参数是一个[Hander](http://gowalker.org/net/http#Handler)是啥呢，它是一个接口。这个接口很简单，只要某个struct有`ServeHTTP(http.ResponseWriter, *http.Request)`这个方法，那这个struct就自动实现了[Hander](http://gowalker.org/net/http#Handler)接口。



> ServeHTTP方法，他需要2个参数，一个是http.ResponseWriter，另一个是*http.Request
> 往http.ResponseWriter写入什么内容，浏览器的网页源码就是什么内容
> *http.Request里面是封装了，浏览器发过来的请求（包含路径、浏览器类型等等）

```go
package main

import (
    "io"
    "net/http"
)

type a struct{}

func (*a) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path := r.URL.String()
    switch path {
    case "/":
        io.WriteString(w, "<h1>root</h1><a href=\"abc\">abc</a>")
    case "/abc":
        io.WriteString(w, "<h1>abc</h1><a href=\"/\">root</a>")
    }
}

func main() {
    http.ListenAndServe(":8080", &a{})//第2个参数需要实现Hander接口的struct，a满足
}
```

> ServeMux大致作用是，他有一张map表，map里的key记录的是r.URL.String()，而value记录的是一个方法，这个方法和ServeHTTP是一样的，这个方法有一个别名，叫[HandlerFunc](http://gowalker.org/net/http#HandlerFunc)
> ServeMux还有一个方法名字是[Handle](http://gowalker.org/net/http#ServeMux_Handle)，他是用来注册[HandlerFunc](http://gowalker.org/net/http#HandlerFunc) 的
> ServeMux还有另一个方法名字是ServeHTTP，这样ServeMux是实现Handler接口的，否者无法当http.ListenAndServe的第二个参数传输。
>
> 

```Go
package main

import (
    "net/http"
    "io"
)

type b struct{}

func (*b) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "hello")
}
func main() {
    mux := http.NewServeMux()
    mux.Handle("/h", &b{})
    http.ListenAndServe(":8080", mux)
}

mux := http.NewServeMux():新建一个ServeMux。
mux.Handle("/", &b{}):注册路由，把"/"注册给b这个实现Handler接口的struct，注册到map表中。
http.ListenAndServe(":8080", mux)第二个参数是mux。
运行时，因为第二个参数是mux，所以http会调用mux的ServeHTTP方法。
ServeHTTP方法执行时，会检查map表（表里有一条数据，key是“/h”，value是&b{}的ServeHTTP方法）
如果用户访问/h的话，mux因为匹配上了，mux的ServeHTTP方法会去调用&b{}的 ServeHTTP方法，从而打印hello
如果用户访问/abc的话，mux因为没有匹配上，从而打印404 page not found
```

**ServeMux就是个二传手！**

\##ServeMux的HandleFunc方法

发现了没有，b这个struct仅仅是为了装一个ServeHTTP而存在，所以能否跳过b呢，ServeMux说：可以[mux.HandleFunc](http://gowalker.org/net/http#ServeMux_HandleFunc)是用来注册func到map表中的

```go
package main

import (
    "net/http"
    "io"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/h", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "hello")
    })
    mux.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "byebye")
    })
    mux.HandleFunc("/hello", sayhello)
    http.ListenAndServe(":8080", mux)
}

func sayhello(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "hello world")
}
```

\##回到开头

回到开头，有让大家先忘掉`http.HandleFunc("/", sayhello)` *请先忘记引子里的http.HandleFunc("/", sayhello)，这个要到很后面才提到*

当`http.ListenAndServe(":8080", nil)`的第2个参数是nil时
**http内部会自己建立一个叫DefaultServeMux的ServeMux，因为这个ServeMux是http自己维护的，如果要向这个ServeMux注册的话，就要用http.HandleFunc这个方法啦**，现在看很简单吧





2.go可以使用一个未先被声明的变量？

[变量声明的顺序](https://studygolang.com/articles/4785)



3.函数也能返回局部变量的地址（指针）？

[先看C语言](http://blog.csdn.net/haiwil/article/details/6691854)

>  一般的来说，函数是可以返回局部变量的。 局部变量的作用域只在函数内部，**在函数返回后，局部变量的内存已经释放了。**因此，如果函数返回的是局部变量的值，不涉及地址，程序不会出错。但是如果返回的是局部变量的地址(指针)的话，程序运行后会出错。因为函数只是把指针复制后返回了，但是指针指向的内容已经被释放了，这样指针指向的内容就是不可预料的内容，调用就会出错。准确的来说，函数不能通过返回指向栈内存的指针(**注意这里指的是栈，返回指向堆内存的指针是可以的**)。



```C
#include <stdio.h>   
char *returnStr()   
{   
    char *p="hello world!";   
    return p;   
}   
int main()   
{   
    char *str;   
    str=returnStr();   
    printf("%s\n", str);   
    return 0;   
}  
```

这个没有任何问题，因为"hello world!"是一个字符串常量，存放在只读数据段，把该字符串常量存放的只读数据段的首地址赋值给了指针，**<u>所以returnStr函数退出时，该该字符串常量所在内存不会被回收，故能够通过指针顺利无误的访问。</u>**

```C
#include <stdio.h>   
char *returnStr()   
{   
    char p[]="hello world!";   
    return p;   
}   
int main()   
{   
    char *str;   
    str=returnStr();   
    printf("%s\n", str);   
    return 0;   
}   
```

"hello world!"是局部变量存放在栈中。当returnStr函数退出时，**<u>栈要清空，局部变量的内存也被清空了，所以这时的函数返回的是一个已被释放的内存地址，所以有可能打印出来的是乱码。</u>** 



**<u>数组是不能作为函数的返回值的，原因是编译器把数组名认为是局部变量（数组）的地址。返回一个数组一般用返回指向这个数组的指针代替，而且这个指针不能指向一个自动数组，因为函数结束后自动数组被抛弃</u>**，但可以返回一个指向静态局部数组的指针，因为静态存储期是从对象定义到程序结束的



4.如果你觉得Go的方法和接口容易忘

[方法和接口](https://studygolang.com/articles/4675)

**方法就是有接收者的函数。**

接口定义为一个方法的集合。方法包含实际的代码。换句话说，一个接口就是定义，而方法就是实现。因此，接收者不能定义为接口类型，这样做的话会引起 invalid receiver type … 的编译器错误。