# LoxGo

根据[crafting interpreters](https://craftinginterpreters.com/)中的教程一步步学习编译原理的知识，从零开始编写Lox语言，使用Golang语言编写。

# 进度
*标注有\*号为额外完成进度*
- [x] 词法分析
- [x] 语法分析
- [x] (*) 类型系统
  - [x] 内置类型及转换函数 (int, float, bool, string, array)
  - [x] 获取类型对应字符串函数 (type)
  - [x] 获取数组长度函数 (len)
- [ ] 解释器
  - [x] 表达式求值
  - [x] (*) 自增自减运算符 `a++; a--;++a;--a;`
  - [x] (*) 三目运算符 `var b = a < 5 ? 0 : 1;`
  - [x] 输出语句 `print a;`
  - [x] (*) (多)变量定义语句 `var a = 2, b = 3;`
  - [x] 变量赋值语句 `a = 5;`
  - [x] 控制流相关
    - [x] if `if (condition) {statments...} else {statments...}`
    - [x] while `while(condition) {statments...}`
    - [x] for `for (init;cond;inc) {statments...}`
    - [x] (*) break `break;`
    - [x] (*) continue `continue;`
  - [x] 函数相关 
    - [x] 函数定义 `fun demo(a, b) {statments...}`
    - [x] 函数调用 `demo(1, 2);`
    - [x] (*) return语句 `return a+b;` (并非使用异常来处理返回值)
    - [x] 闭包
    - [x] (*) 匿名函数 `demo(func(a) {statments...}, 0)`
  - [x] 静态解析
    - [x] 变量作用域
    - [x] 在块中重定义变量错误
    - [x] 在非函数中使用return语句错误
  - [ ] 类相关

- 使用
```bash
go run main.go # 进入交互模式
go run main.go ./tests/if.lox # 运行文件
```

# 参考
- [crafting interpreters](https://craftinginterpreters.com/contents.html)
- [crafting interpreters zh](https://github.com/GuoYaxiang/craftinginterpreters_zh)