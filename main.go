package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"go-blankstrip/internal/strip"
)

func main() {
	var opt strip.Options
	flag.BoolVar(&opt.Trailing, "t", false, "去行尾空白")
	flag.BoolVar(&opt.Leading, "l", false, "去行首空白")
	flag.BoolVar(&opt.Blank, "b", false, "删掉所有空行")
	flag.BoolVar(&opt.Squeeze, "s", false, "连续空行压成一行")
	flag.BoolVar(&opt.TrimEnds, "e", false, "去掉文件首尾的空行")
	flag.IntVar(&opt.Tabs, "tabs", 0, "把制表符换成 N 个空格")
	inplace := flag.Bool("w", false, "原地改写文件")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		if *inplace {
			fmt.Fprintln(os.Stderr, "-w 需要指定文件")
			os.Exit(1)
		}
		if err := strip.Run(os.Stdin, os.Stdout, opt); err != nil {
			fmt.Fprintln(os.Stderr, "处理失败:", err)
			os.Exit(1)
		}
		return
	}

	for _, path := range args {
		if err := handle(path, opt, *inplace); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			os.Exit(1)
		}
	}
}

func handle(path string, opt strip.Options, inplace bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if !inplace {
		return strip.Run(f, os.Stdout, opt)
	}

	var buf bytes.Buffer
	if err := strip.Run(f, &buf, opt); err != nil {
		return err
	}
	f.Close()

	// 先写临时文件再改名，中途出错也不会把原文件弄坏
	st, err := os.Stat(path)
	mode := os.FileMode(0o644)
	if err == nil {
		mode = st.Mode()
	}
	tmp := filepath.Join(filepath.Dir(path), "."+filepath.Base(path)+".tmp")
	if err := os.WriteFile(tmp, buf.Bytes(), mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func usage() {
	fmt.Fprint(os.Stderr, `go-blankstrip — 清理文本里的多余空白

用法:
  go-blankstrip [选项] [文件...]      不给文件就读标准输入

选项:
  -t         去行尾空白
  -l         去行首空白
  -b         删掉所有空行（纯空格的行也算）
  -s         连续空行压成一行
  -e         去掉文件开头和结尾的空行
  -tabs N    把制表符换成 N 个空格
  -w         原地改写文件（需要指定文件）

例子:
  go-blankstrip -t -s notes.txt
  go-blankstrip -t -w src/*.go
  cat messy.txt | go-blankstrip -b
`)
}
