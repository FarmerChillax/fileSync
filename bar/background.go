package bar

import "fmt"

type Bar struct {
	percent int64
	cur     int64
	total   int64
	rate    string
	graph   string
}

// 初始化
// 自定义显示图案
func (bar *Bar) NewOptionWithGraph(start, total int64, graph string) {
	bar.graph = graph
	bar.NewOption(start, total)
}

// 初始化
func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph //初始化进度条位置
	}
}

// 获取百分比
func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

// 显示
func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%4 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-25s]%3d%%  %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *Bar) Finish() {
	fmt.Println()
}
