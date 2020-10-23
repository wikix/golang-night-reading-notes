## 如何阅读 `Golang` 源代码?

杨文



（开发人员）方法论：

1、以终为始：带着问题（Issue）找答案（PR）

2、Main/Test 入口函数

3、pkg

4、Exported Function 层级调用

5、兴趣

6、变更 Commit

例如

image commit history

github golang/go image package的所有历史

https://github.com/yangwenmai/learning-growth/blob/master/04.researching/history/image_history.md



----

Felix

学习论：

#### 需要明确自己的目标：要从什么角度出发读代码



- 学习优雅代码
  - 架构优雅：简洁优雅的架构/抽象，如何恰当使用语言机制实现具体功能
    - 这种目的的代码阅读，需要有一定的背景知识，了解目标库的具体业务场景，从具体场景出发，了解目标库的代码实现优雅之美，自己带着问题思考，如果自己来做同样的功能/逻辑的实现，能否有不同的方法，再和目标库的实现作对比
  - 性能优雅：如何提示代码的性能，这部分一般偏向并发编程 & 底层编程
    - 这种目的的代码阅读，需要有系统底层相关的背景，包括但不局限于：操作系统、文件系统、IO异步处理、网络编程、并发处理机制等等。带着相关背景去看代码，结合具体系统负载查看目标库是如何做技术选型的，建议在过程中不断用 benchmark 验证想法
- 学习编写健壮的代码
  - 这种目标的代码阅读，需要从具体问题出发，一般从 issue 开始，看某些具体场景下会出现什么问题，后面又有什么相关的讨论，之后结合pullrequest又是做出了怎样的修复。从这个完整的链路出发，了解到：什么场景下容易出现问题；有哪些常见的问题；针对这些问题一般有哪些常见的问题定位思路；如果前面几步都很顺利的话，最后的解决方案几乎是最简单的步骤了。整个链路中最重要的部分，是 `issue`上下文的阅读 以及 PR讨论的阅读，代码的具体实现倒并不是最最重要的。
- 了解自己常使用库的实现
  - 这种目的出发的代码阅读目的性就很明确，就是为了看自己使用库是否存在潜在的“坑”？是否能满足自己的使用场景？这种情况可以结合自己的具体使用场景，针对库里的不同方法，一点一点做 test 和 benchmark，验证正确性和性能，保证逻辑和性能都满足自己的需求。



----

Felix 代码

#### 实际阅读方法

阅读环境：

- 网页端（github 配合 SourceGraph，可以代码搜索、代码关键字高亮、代码调整等基础功能）
- 本地IDE，随时可以添加代码跑 test 验证想法，可以设置热键



#### 阅读代码从哪里来？

- 首先，Go官网 https://blog.golang.org/，介绍go最佳实践，如对背后具体实现逻辑感兴趣，可以读相关代码
- 其次，Go邮件组 https://groups.google.com/forum/#!forum/golang-nuts，Go主题讨论，Go core team的工程师也会在这里回答问题
- 关注库的github issue区的 求助&讨论
- 最后，Go的语言提案，https://github.com/golang/proposal/ 新的语言层级提案在此列出



----

欧神（changkun Ou） 经验

1、不抵触英文资料



2、初级资料

golang.org/doc

blog.golang.org

golang.org/pkg



3、进阶资料

dev.golang.org

github.com/golang/go

github.com/golang/proposal

github.com/golang/go/wiki

go-review.googlesource.com

groups.google.com/g/golang-nuts

groups.google.com/g/golang-dev

groups.google.com/g/golang-tools

reddit.com/r/golang



4、前沿资料

www.sigplan.org( Special Interest Group on Programming Languages )

dl.acm.org

arxiv.org/list/cs.PL/recent

scholar.google.de 



----

饶全成

1、输出来源于输入：写文章——看源码——带着兴趣看源码

2、反复看源码，多读几遍

3、补充其他知识，扩展成知识体系



----



傲飞

1、初学者——从具体项目 入坑——第三方库入手，例如：gin/echo库入坑 net/http使用，然后入坑http

golang web方向 | golang 系统底层   等，如上



----

曹大（Xargin）



为什么要读源码？

解决问题



互联网线上项目的场景，例如：压力高的时候为什么为崩？程序为什么崩？

对源码的掌握程度，可以提供你调试、优化的手段。