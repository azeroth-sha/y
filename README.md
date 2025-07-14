# y - 一个自用的工具库

y 是一个 Go 语言编写的多功能工具库，旨在为开发人员提供一系列常用的基础功能。这个库包含了数据类型转换、文件操作、随机数生成等多种实用功能。

## 特性

- **模块化设计**：每个功能独立封装为包，便于扩展和维护。
- **高性能**：利用 Go 的并发特性和高效的标准库实现高性能操作。
- **简洁易用**：API 设计直观简单，易于使用。
- **广泛的功能覆盖**：涵盖数据类型转换、文件处理、随机数生成等多个方面。

## 安装

要安装 y 工具库，请确保已经安装了 [Go](https://golang.org/dl/) (1.22+) 和 Git：

```bash
go get github.com/azeroth-sha/y
```


## 快速入门

以下是一些常见用法示例：

### 数据类型转换

```go
package main

import (
    "fmt"
    "github.com/azeroth-sha/y/yconv"
)

func main() {
    // 将字符串转换为布尔值
    b, _ := yconv.Bool("true")
    fmt.Println(b) // 输出: true

    // 将整数转换为字符串
    s := yconv.MustString(123)
    fmt.Println(s) // 输出: "123"

    // 将字节切片转换为字符串（无内存拷贝）
    bs := []byte("hello")
    str := yconv.ToString(bs)
    fmt.Println(str) // 输出: "hello"
}
```


### 文件操作

```go
package main

import (
    "fmt"
    "github.com/azeroth-sha/y/yfile"
)

func main() {
    // 检查文件是否存在
    exists := yfile.IsExist("test.txt")
    fmt.Println(exists) // 输出: false 或 true

    // 复制文件
    err := yfile.FileCopy("source.txt", "destination.txt")
    if err != nil {
        panic(err)
    }
}
```


### 随机数生成

```go
package main

import (
    "fmt"
    "github.com/azeroth-sha/y/yrand"
)

func main() {
    // 生成随机 uint64
    randUint64 := yrand.Uint64()
    fmt.Println(randUint64)

    // 生成包含数字和字母的随机字符串
    randStr := yrand.StringBy(10, yrand.AlphaNum)
    fmt.Println(randStr)
}
```


### 缓存操作

```go
package main

import (
    "fmt"
    "time"
    "github.com/azeroth-sha/y/ycache"
)

func main() {
    // 创建一个新的缓存实例
    cache := ycache.New()

    // 设置键值对并指定过期时间
    cache.Set("key", "value", ycache.WithItemDur(5*time.Second))

    // 获取键值
    val, ok := cache.Get("key")
    fmt.Println(val.(string), ok) // 输出: value true

    // 等待键过期
    time.Sleep(6 * time.Second)

    // 检查键是否已过期
    has, exp := cache.DelExpired("key")
    fmt.Println(has, exp) // 输出: true true
}
```


### 锁机制

```go
package main

import (
    "sync"
    "time"
    "github.com/azeroth-sha/y/ylock"
)

func main() {
    // 创建一个名称锁
    locker := ylock.NewNameLocker()
    key := "my_key"

    // 获取锁
    locker.Lock(key)
    
    // 在另一个 goroutine 中尝试获取同一个锁
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        locker.Lock(key)
        fmt.Println("成功获取锁")
        locker.Unlock(key)
    }()

    // 休眠一段时间后释放锁
    time.Sleep(2 * time.Second)
    locker.Unlock(key)

    wg.Wait() // 等待 goroutine 完成
}
```


### 日志记录

```go
package main

import (
    "github.com/azeroth-sha/y/logger"
    "github.com/sirupsen/logrus"
)

func main() {
    // 初始化日志记录器
    log := logrus.New()
    logger.SetDefault(log)

    // 记录信息级别的日志
    logger.DefaultLog().Info("这是一个信息级别的日志")

    // 记录错误级别的日志
    logger.DefaultLog().Error("这是一个错误级别的日志")
}
```


## 目录结构

```
.
├── internal           # 内部使用的功能
│   ├── calc.go        # 条件选择函数
│   ├── host.go        # 主机标识符相关功能
│   └── unsafe.go      # 不安全操作，用于优化性能
├── logger             # 日志功能
│   └── logger.go
├── ybuff              # 缓冲区操作
│   └── buff.go
├── ycache             # 缓存功能
│   ├── cache.go       # 缓存接口定义及实现
│   ├── item.go        # 缓存项定义
│   ├── option.go      # 缓存选项配置
│   └── shard.go       # 分片逻辑实现
├── yconst             # 常量定义
│   └── number.go      # 数字常量定义
├── yconv              # 数据类型转换
│   ├── bool.go        # 布尔值转换
│   ├── digit.go       # 数字类型转换
│   ├── string.go      # 字符串转换
│   └── unsafe.go      # 不安全转换
├── yfile              # 文件处理
│   └── file.go        # 文件操作函数
├── ygrace             # 优雅启动和关闭服务
│   ├── grace.go       # Grace 接口及其实现
│   ├── option.go      # 服务选项配置
│   └── service.go     # 服务接口定义
├── yguid              # GUID 生成
│   └── guid.go        # GUID 类型定义及操作
├── yhttp              # HTTP 客户端
│   └── http.go        # HTTP 客户端实现
├── ylock              # 锁机制
│   ├── define.go      # 锁接口定义
│   ├── hash_mutex.go  # 基于哈希的互斥锁
│   ├── mutex_test.go  # 锁机制单元测试
│   ├── name_mutex.go  # 基于名称的互斥锁
│   └── pool_mutex.go  # 基于池的互斥锁
├── yrand              # 随机数生成
│   └── rand.go        # 随机数生成函数
├── ysum               # 校验和计算
│   ├── crc16.go       # CRC16 校验和实现
│   └── tools.go       # 校验和工具函数
├── ytime              # 时间处理
│   ├── calc.go        # 时间计算辅助函数
│   ├── consts.go      # 时间格式常量
│   └── conv.go        # 时间转换函数
├── define.go          # 全局定义
├── go.mod             # Go Module 配置
├── README.md          # 当前文档
└── tools.go           # 工具函数
```


## 贡献指南

我们欢迎任何形式的贡献！如果您有兴趣参与开发或改进 `y` 库，请遵循以下指导原则：

1. **编码规范**：遵循 Go 官方编码规范。
2. **单元测试**：对于每一个 PR，我们都要求有充分的单元测试来验证代码的正确性。
3. **文档更新**：如果您的更改影响到了公共 API 或者添加了新功能，请记得更新相关的文档。

## 协议

本项目采用 Apache License 2.0 协议进行许可。有关详细信息，请参阅 LICENSE 文件。

## 联系方式

如果你有任何问题或者想要提出建议，请通过 GitHub Issues 页面与我们联系。