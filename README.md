***
View in [[English](README-en.md)][[中文](README.md)]
***
# confilter
敏感词判定工具

[![Build Status](https://drone.io/github.com/hanguofeng/confilter/status.png)](https://drone.io/github.com/hanguofeng/confilter/latest)  [![Coverage Status](https://coveralls.io/repos/hanguofeng/confilter/badge.png)](https://coveralls.io/r/hanguofeng/confilter)


Feature
-------
* 支持加载多个词库，并返回是否包含各词库中的敏感词
* 使用[[Aho-Corasick算法](http://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_string_matching_algorithm)]

Useage
------
**安装**

	go get github.com/hanguofeng/confilter

**Quick Start**

参考 [confilter_test.go](confilter_test.go)

参考 [server/](server)

**文档**

[[confilter Wiki](https://github.com/hanguofeng/confilter/wiki)]

TODO
----
* 支持插件机制以实现对文本的预加工、排除干扰等

LICENCE
-------
confilter使用[[MIT许可协议](LICENSE)]

