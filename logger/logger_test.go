package logger_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/wiselike/revel/logger"
)

// 可记录调用次数与最后一条记录的简单 handler
type testRecorder struct {
	hit  int
	last *logger.Record
}

func (r *testRecorder) Log(rec *logger.Record) error {
	r.hit++
	r.last = rec
	return nil
}

var msg = "unit-test"
var re = regexp.MustCompile(`.*?` + msg + `\s*`)

func TestCompositeMultiHandler_AppendRecorder(t *testing.T) {
	t.Parallel()

	// 1.先得到 CompositeMultiHandler（Revel 默认会用）
	cmh, _ := logger.NewCompositeMultiHandler()

	// 2. 准备 handler
	newProbe := &testRecorder{} // 我们本次测试要追加的 probe

	// 3. 追加我们的 testRecorder，同样 replace=false
	cmh.SetHandler(newProbe, false, logger.LvlDebug)

	var buf bytes.Buffer
	cmh.SetHandler(logger.StreamHandler(&buf, logger.TerminalFormatHandler(true, true)), false, logger.LvlDebug)

	// 4. 把 cmh 挂到一个 logger 上
	l := logger.New()
	l.SetHandler(cmh)

	// 5. 触发一条 debug 日志
	l.Debug(msg, "x", 1, "y", 3.2, "equals", "=", "quote", `"`,
		"nil", nil, "carriage_return", "bang"+string('\r')+"foo", "tab", "bar	baz", "newline", "foo\nbar")

	// 6. 断言：handler 应该收到一次调用
	if newProbe.hit != 1 {
		t.Fatalf("probe handler hit=%d, want 1", newProbe.hit)
	}

	// 7. 再简单验证 probe 捕获到的消息内容
	if got := newProbe.last.Message; got != msg {
		t.Errorf("probe message=%q, want %q", got, msg)
	}

	// 8. 验证变量打印顺序
	got := re.ReplaceAllString(buf.String(), "")
	expected := `x=1 y=3.2000000 equals="=" quote="\"" nil=nil carriage_return="bang\rfoo" tab="bar\tbaz" newline="foo\nbar"` + "\n"
	if got != expected {
		t.Fatalf("Got `%s`, expected `%s`", got, expected)
	}
}
