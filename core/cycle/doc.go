// Package cycle 提供并发生命周期控制器，将一组后台任务纳入统一的启动、等待与关闭流程。
//
// 任务通过 Run 追加为后台协程；Done 与 DoneAndClose 等待全部任务结束，
// Wait 返回接收任务错误的通道，Close 用于关闭退出通道。
package cycle
