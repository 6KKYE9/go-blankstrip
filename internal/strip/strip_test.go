package strip

import (
	"bytes"
	"strings"
	"testing"
)

func run(t *testing.T, in string, opt Options) string {
	t.Helper()
	var out bytes.Buffer
	if err := Run(strings.NewReader(in), &out, opt); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func TestTrailing(t *testing.T) {
	got := run(t, "a   \nb\t\n", Options{Trailing: true})
	if got != "a\nb\n" {
		t.Fatalf("行尾空白没去掉: %q", got)
	}
}

func TestLeading(t *testing.T) {
	got := run(t, "   a\n\tb\n", Options{Leading: true})
	if got != "a\nb\n" {
		t.Fatalf("行首空白没去掉: %q", got)
	}
}

func TestBlankRemovesAll(t *testing.T) {
	got := run(t, "a\n\n\nb\n", Options{Blank: true})
	if got != "a\nb\n" {
		t.Fatalf("空行没删干净: %q", got)
	}
}

// 只有空格的行也该算空行
func TestBlankCountsWhitespaceOnly(t *testing.T) {
	got := run(t, "a\n   \n\t\nb\n", Options{Blank: true})
	if got != "a\nb\n" {
		t.Fatalf("纯空白行应算空行: %q", got)
	}
}

func TestSqueeze(t *testing.T) {
	got := run(t, "a\n\n\n\nb\n", Options{Squeeze: true})
	if got != "a\n\nb\n" {
		t.Fatalf("连续空行应压成一行: %q", got)
	}
}

// 两个选项都开时 Blank 优先，不能互相打架
func TestBlankBeatsSqueeze(t *testing.T) {
	got := run(t, "a\n\n\nb\n", Options{Blank: true, Squeeze: true})
	if got != "a\nb\n" {
		t.Fatalf("Blank 应优先: %q", got)
	}
}

func TestTabs(t *testing.T) {
	got := run(t, "\ta\n", Options{Tabs: 4})
	if got != "    a\n" {
		t.Fatalf("制表符没展开: %q", got)
	}
}

// 展开制表符要在去行首之前做，否则 tab 缩进去不掉
func TestTabsThenLeading(t *testing.T) {
	got := run(t, "\t\ta\n", Options{Tabs: 2, Leading: true})
	if got != "a\n" {
		t.Fatalf("展开后应能去掉缩进: %q", got)
	}
}

func TestTrimEnds(t *testing.T) {
	got := run(t, "\n\na\nb\n\n\n", Options{TrimEnds: true})
	if got != "a\nb\n" {
		t.Fatalf("首尾空行没去掉: %q", got)
	}
}

// 中间的空行不能被 TrimEnds 误伤
func TestTrimEndsKeepsMiddle(t *testing.T) {
	got := run(t, "\na\n\nb\n\n", Options{TrimEnds: true})
	if got != "a\n\nb\n" {
		t.Fatalf("中间空行不该动: %q", got)
	}
}

// 全是空行时 TrimEnds 不能越界
func TestTrimEndsAllBlank(t *testing.T) {
	got := run(t, "\n\n\n", Options{TrimEnds: true})
	if got != "" {
		t.Fatalf("全空行应输出空: %q", got)
	}
}

func TestBOMStripped(t *testing.T) {
	got := run(t, "\ufeffa\nb\n", Options{})
	if got != "a\nb\n" {
		t.Fatalf("BOM 没去掉: %q", got)
	}
}

func TestEmpty(t *testing.T) {
	if got := run(t, "", Options{Trailing: true, Blank: true}); got != "" {
		t.Fatalf("空输入应输出空: %q", got)
	}
}

func TestNoOptionsPassthrough(t *testing.T) {
	in := "a  \n\n  b\n"
	if got := run(t, in, Options{}); got != in {
		t.Fatalf("不开选项应原样输出: %q", got)
	}
}

func TestLongLine(t *testing.T) {
	long := strings.Repeat("x", 200*1024)
	got := run(t, long+"   \n", Options{Trailing: true})
	if got != long+"\n" {
		t.Fatal("长行处理出错")
	}
}
