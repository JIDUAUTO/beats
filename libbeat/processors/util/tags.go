package util

const (
	BenchmarkPrefix = "00000000"
	// BenchmarkExcludePrefix 吉利调用集度接口时，传递的trace_id以16位0开头，为了防止这种非
	// 压测日志被过滤，特由此逻辑
	BenchmarkExcludePrefix = BenchmarkPrefix + BenchmarkPrefix
	MsgTag                 = "##JIDU##"
	MsgTagConcatenated     = "##JIDU####JIDU##"
)
