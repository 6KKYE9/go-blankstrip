package strip

import (
	"bufio"
	"io"
	"strings"
)

// Options 控制清理哪些空白。
type Options struct {
	Trailing bool // 去行尾空白
	Leading  bool // 去行首空白
	Blank    bool // 删掉所有空行
	Squeeze  bool // 连续空行压成一行
	Tabs     int  // >0 时把制表符换成这么多空格
	TrimEnds bool // 去掉文件开头和结尾的空行
}

// Run 逐行清理后写出。
func Run(r io.Reader, w io.Writer, opt Options) error {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	var lines []string
	first := true
	for sc.Scan() {
		line := sc.Text()
		if first {
			line = strings.TrimPrefix(line, "\ufeff")
			first = false
		}
		if opt.Tabs > 0 {
			line = strings.ReplaceAll(line, "\t", strings.Repeat(" ", opt.Tabs))
		}
		if opt.Leading {
			line = strings.TrimLeft(line, " \t")
		}
		if opt.Trailing {
			line = strings.TrimRight(line, " \t")
		}
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		return err
	}

	lines = filterBlank(lines, opt)
	if opt.TrimEnds {
		lines = trimEnds(lines)
	}

	bw := bufio.NewWriter(w)
	defer bw.Flush()
	for _, l := range lines {
		if _, err := bw.WriteString(l + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// filterBlank 处理空行。Blank 优先于 Squeeze——都删了就没得压了。
func filterBlank(lines []string, opt Options) []string {
	if opt.Blank {
		out := lines[:0:0]
		for _, l := range lines {
			if !isBlank(l) {
				out = append(out, l)
			}
		}
		return out
	}
	if opt.Squeeze {
		out := lines[:0:0]
		prevBlank := false
		for _, l := range lines {
			b := isBlank(l)
			if b && prevBlank {
				continue
			}
			out = append(out, l)
			prevBlank = b
		}
		return out
	}
	return lines
}

func trimEnds(lines []string) []string {
	i := 0
	for i < len(lines) && isBlank(lines[i]) {
		i++
	}
	j := len(lines)
	for j > i && isBlank(lines[j-1]) {
		j--
	}
	return lines[i:j]
}

// isBlank 只有空白字符的行也算空行。
func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}
