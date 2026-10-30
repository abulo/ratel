// Package change 实现 toolkit 的更新日志（CHANGELOG）生成命令。
//
// 该包读取当前 Git 仓库的版本标签、标签日期与标签区间内的提交记录，按约定式提交
// 类型（feat、fix、docs、perf 等）分类整理，生成带版本区间链接的 Markdown 更新日志，
// 写入项目根目录 CHANGELOG.md。
package change
