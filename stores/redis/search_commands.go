package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// FT_List 返回当前全部已创建搜索索引的名称。
func (r *Client) FT_List(ctx context.Context) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FT_List(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// FTAggregate 在 index 上执行聚合查询 query，返回聚合结果 map。
func (r *Client) FTAggregate(ctx context.Context, index, query string) (val map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTAggregate(getCtx(ctx), index, query).Result()
		return err
	}, acceptable)
	return
}

// FTAggregateWithArgs 按 FTAggregateOptions 执行聚合查询，可配置分组、排序、分页与加载字段。
func (r *Client) FTAggregateWithArgs(ctx context.Context, index, query string, options *redis.FTAggregateOptions) (val *redis.FTAggregateResult, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTAggregateWithArgs(getCtx(ctx), index, query, options).Result()
		return err
	}, acceptable)
	return
}

// FTAliasAdd 为索引 index 添加别名 alias，查询时可使用别名访问该索引。
func (r *Client) FTAliasAdd(ctx context.Context, index, alias string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTAliasAdd(getCtx(ctx), index, alias).Result()
		return err
	}, acceptable)
	return
}

// FTAliasDel 删除别名 alias。
func (r *Client) FTAliasDel(ctx context.Context, alias string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTAliasDel(getCtx(ctx), alias).Result()
		return err
	}, acceptable)
	return
}

// FTAliasUpdate 将别名 alias 重新指向索引 index。
func (r *Client) FTAliasUpdate(ctx context.Context, index, alias string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTAliasUpdate(getCtx(ctx), index, alias).Result()
		return err
	}, acceptable)
	return
}

// FTAlter 为索引 index 追加字段定义；skipInitialScan 指示是否跳过对现有文档的初始扫描。
func (r *Client) FTAlter(ctx context.Context, index string, skipInitialScan bool, definition []any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTAlter(getCtx(ctx), index, skipInitialScan, definition).Result()
		return err
	}, acceptable)
	return
}

// FTConfigGet 返回 option 指定的全文检索配置项；option 为空或通配时返回全部配置。
func (r *Client) FTConfigGet(ctx context.Context, option string) (val map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTConfigGet(getCtx(ctx), option).Result()
		return err
	}, acceptable)
	return
}

// FTConfigSet 将 option 配置项设置为 value。
func (r *Client) FTConfigSet(ctx context.Context, option string, value any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTConfigSet(getCtx(ctx), option, value).Result()
		return err
	}, acceptable)
	return
}

// FTCreate 按 FTCreateOptions 与字段 schema 创建名为 index 的搜索索引。
func (r *Client) FTCreate(ctx context.Context, index string, options *redis.FTCreateOptions, schema ...*redis.FieldSchema) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTCreate(getCtx(ctx), index, options, schema...).Result()
		return err
	}, acceptable)
	return
}

// FTCursorDel 删除聚合查询游标 index 上的 cursorId，释放服务端资源。
func (r *Client) FTCursorDel(ctx context.Context, index string, cursorId int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTCursorDel(getCtx(ctx), index, cursorId).Result()
		return err
	}, acceptable)
	return
}

// FTCursorRead 从聚合游标 cursorId 读取 count 条后续结果。
func (r *Client) FTCursorRead(ctx context.Context, index string, cursorId, count int) (val map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTCursorRead(getCtx(ctx), index, cursorId, count).Result()
		return err
	}, acceptable)
	return
}

// FTDictAdd 向停用词/扩展词典 dict 添加 term 术语，返回新增术语数量。
func (r *Client) FTDictAdd(ctx context.Context, dict string, term ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTDictAdd(getCtx(ctx), dict, term...).Result()
		return err
	}, acceptable)
	return
}

// FTDictDel 从词典 dict 删除 term 术语。
func (r *Client) FTDictDel(ctx context.Context, dict string, term ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTDictDel(getCtx(ctx), dict, term...).Result()
		return err
	}, acceptable)
	return
}

// FTDictDump 导出词典 dict 中的全部术语。
func (r *Client) FTDictDump(ctx context.Context, dict string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTDictDump(getCtx(ctx), dict).Result()
		return err
	}, acceptable)
	return
}

// FTDropIndex 删除索引 index；索引内文档默认保留。
func (r *Client) FTDropIndex(ctx context.Context, index string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTDropIndex(getCtx(ctx), index).Result()
		return err
	}, acceptable)
	return
}

// FTDropIndexWithArgs 按 FTDropIndexOptions 删除索引，选项可指定是否同时删除文档。
func (r *Client) FTDropIndexWithArgs(ctx context.Context, index string, options *redis.FTDropIndexOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTDropIndexWithArgs(getCtx(ctx), index, options).Result()
		return err
	}, acceptable)
	return
}

// FTExplain 返回 index 上 query 的执行计划文本，用于分析检索流程。
func (r *Client) FTExplain(ctx context.Context, index, query string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTExplain(getCtx(ctx), index, query).Result()
		return err
	}, acceptable)
	return
}

// FTExplainWithArgs 按 FTExplainOptions 返回查询执行计划。
func (r *Client) FTExplainWithArgs(ctx context.Context, index, query string, options *redis.FTExplainOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTExplainWithArgs(getCtx(ctx), index, query, options).Result()
		return err
	}, acceptable)
	return
}

// FTInfo 返回索引 index 的字段、文档数、索引大小等统计信息。
func (r *Client) FTInfo(ctx context.Context, index string) (val redis.FTInfoResult, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTInfo(getCtx(ctx), index).Result()
		return err
	}, acceptable)
	return
}

// FTSpellCheck 对 query 做拼写检查，返回各词的候选纠正结果。
func (r *Client) FTSpellCheck(ctx context.Context, index, query string) (val []redis.SpellCheckResult, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTSpellCheck(getCtx(ctx), index, query).Result()
		return err
	}, acceptable)
	return
}

// FTSpellCheckWithArgs 按 FTSpellCheckOptions 执行拼写检查，可限定词典与距离阈值。
func (r *Client) FTSpellCheckWithArgs(ctx context.Context, index, query string, options *redis.FTSpellCheckOptions) (val []redis.SpellCheckResult, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTSpellCheckWithArgs(getCtx(ctx), index, query, options).Result()
		return err
	}, acceptable)
	return
}

// FTSearch 在 index 上执行 query 全文检索，返回命中文档与高亮、打分结果。
func (r *Client) FTSearch(ctx context.Context, index, query string) (val redis.FTSearchResult, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTSearch(getCtx(ctx), index, query).Result()
		return err
	}, acceptable)
	return
}

// FTSearchWithArgs 按 FTSearchOptions 执行检索，可配置过滤、排序、分页、返回字段与摘要。
func (r *Client) FTSearchWithArgs(ctx context.Context, index, query string, options *redis.FTSearchOptions) (val redis.FTSearchResult, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTSearchWithArgs(getCtx(ctx), index, query, options).Result()
		return err
	}, acceptable)
	return
}

// FTSynDump 导出索引 index 的全部同义词组及其术语。
func (r *Client) FTSynDump(ctx context.Context, index string) (val []redis.FTSynDumpResult, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTSynDump(getCtx(ctx), index).Result()
		return err
	}, acceptable)
	return
}

// FTSynUpdate 将 terms 加入索引 index 的同义词组 synGroupId。
func (r *Client) FTSynUpdate(ctx context.Context, index string, synGroupId any, terms []any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTSynUpdate(getCtx(ctx), index, synGroupId, terms).Result()
		return err
	}, acceptable)
	return
}

// FTSynUpdateWithArgs 按 FTSynUpdateOptions 更新同义词组，可配置是否清空现有术语。
func (r *Client) FTSynUpdateWithArgs(ctx context.Context, index string, synGroupId any, options *redis.FTSynUpdateOptions, terms []any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTSynUpdateWithArgs(getCtx(ctx), index, synGroupId, options, terms).Result()
		return err
	}, acceptable)
	return
}

// FTTagVals 返回索引 index 标签字段 field 中出现过的全部标签值。
func (r *Client) FTTagVals(ctx context.Context, index, field string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FTTagVals(getCtx(ctx), index, field).Result()
		return err
	}, acceptable)
	return
}
