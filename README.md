# go-blankstrip

复制粘贴改格式改到手酸？这玩意儿一行就搞定。

清理文本里的多余空白：行尾空格、空行、制表符。零依赖。

## 装

```
go build -o go-blankstrip .
```

## 用

```
go-blankstrip -t notes.txt              # 去行尾空白
go-blankstrip -s notes.txt              # 连续空行压成一行
go-blankstrip -b notes.txt              # 删掉所有空行
go-blankstrip -t -e -s notes.txt        # 组合着用
go-blankstrip -tabs 4 code.py           # 制表符换 4 个空格

go-blankstrip -t -w src/main.go         # 原地改写
cat messy.txt | go-blankstrip -b        # 管道
```

## 选项

| 选项 | 说明 |
|---|---|
| `-t` | 去行尾空白 |
| `-l` | 去行首空白 |
| `-b` | 删掉所有空行，只含空格的行也算 |
| `-s` | 连续空行压成一行 |
| `-e` | 去掉文件开头和结尾的空行，中间的不动 |
| `-tabs N` | 制表符换成 N 个空格 |
| `-w` | 原地改写文件 |

## 说明

- `-b` 和 `-s` 同时给时 `-b` 优先，空行都删了也没得压。
- 制表符展开在去行首空白之前做，所以 `-tabs 2 -l` 能去掉 tab 缩进。
- `-w` 是先写临时文件再改名，中途出错不会把原文件弄坏，文件权限也保留。
- 首行的 UTF-8 BOM 会去掉。
- 单行放宽到 10MB。
