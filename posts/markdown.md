[title:: MARKDOWN]
(date:: 2026-10-12)  
(author:: ink)
## 标题2
共有1-6级标题
### 换行
文本1后敲两个空格space再回车enter写入文本2
### 换段
文本1后敲两次回车enter写入文本2

> 引用cite/quote
>> 嵌套
## 列表list

- 无序list
  - 嵌套
  > 引用也可用于list

1. 有序list
  2. 嵌套

### task
- [ ] undone
- [x] done

## 代码块code

行内`code`代码

```go
func main() {
    a := 1
    fmt.Println(a)
}
```

#### 数学公式latex
用4个$框住
$$
\frac{\partial f}{\partial x} = 2\sqrt{a}x
$$

### 表格table
|k1|k2|k3|
|:---|:---:|---:|
|v1|v2|v3|
|haaaa|soooo|goooo|

### 脚注
脚注链接[^短内容]

...    
[^短内容]: 详细解释内容，可放在最后。

#### 分割线
***
### 链接 link

[link-rename](link.com)

[跳转到标题2](##标题2)

### 图片

![rename](img路径)

***haha字体测试***
**haha字体测试**
*haha字体测试*
==highlight高亮==
~~del-text~~

