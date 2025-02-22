package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// FT_List 获取所有索引的名称 ctx: 上下文
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

// FTAggregate 执行聚合查询 ctx: 上下文, index: 索引名称, query: 查询语句
func (r *Client) FTAggregate(ctx context.Context, index, query string) (val map[string]interface{}, err error) {
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

// FTAggregateWithArgs 执行聚合查询 ctx: 上下文, index: 索引名称, query: 查询语句, options: 聚合选项
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

// FTAliasAdd 添加别名 ctx: 上下文, index: 索引名称, alias: 别名
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

// FTAliasDel 删除别名 ctx: 上下文, alias: 别名
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

// FTAliasUpdate 更新别名 ctx: 上下文, index: 索引名称, alias: 别名
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

// FTAlter 修改索引 ctx: 上下文, index: 索引名称, skipInitialScan: 是否跳过初始扫描, definition: 定义
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

// FTConfigGet 获取配置 ctx: 上下文, option: 配置项
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

// FTConfigSet 设置配置 ctx: 上下文, option: 配置项, value: 配置值
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

// FTCreate 创建索引 ctx: 上下文, index: 索引名称, options: 创建选项, schema: 字段模式
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

// FTCursorDel 删除游标 ctx: 上下文, index: 索引名称, cursorId: 游标ID
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

// FTCursorRead 读取游标 ctx: 上下文, index: 索引名称, cursorId: 游标ID, count: 读取数量
func (r *Client) FTCursorRead(ctx context.Context, index string, cursorId, count int) (val map[string]interface{}, err error) {
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

// FTDictAdd 添加字典 ctx: 上下文, dict: 字典名称, term: 术语列表
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

// FTDictDel 删除字典 ctx: 上下文, dict: 字典名称, term: 术语列表
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

// FTDictDump 导出字典 ctx: 上下文, dict: 字典名称
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

// FTDropIndex 删除索引 ctx: 上下文, index: 索引名称
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

// FTDropIndexWithArgs 删除索引 ctx: 上下文, index: 索引名称, options: 删除选项
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

// FTExplain 解释查询 ctx: 上下文, index: 索引名称, query: 查询语句
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

// FTExplainWithArgs 解释查询 ctx: 上下文, index: 索引名称, query: 查询语句, options: 解释选项
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

// FTInfo 获取索引信息 ctx: 上下文, index: 索引名称
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

// FTSpellCheck 拼写检查 ctx: 上下文, index: 索引名称, query: 查询语句
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

// FTSpellCheckWithArgs 拼写检查 ctx: 上下文, index: 索引名称, query: 查询语句, options: 拼写检查选项
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

// FTSearch 执行搜索 ctx: 上下文, index: 索引名称, query: 查询语句
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

// FTSearchWithArgs 执行搜索 ctx: 上下文, index: 索引名称, query: 查询语句, options: 搜索选项
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

// FTSynDump 导出同义词 ctx: 上下文, index: 索引名称
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

// FTSynUpdate 更新同义词 ctx: 上下文, index: 索引名称, synGroupId: 同义词组ID, terms: 术语列表
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

// FTSynUpdateWithArgs 更新同义词 ctx: 上下文, index: 索引名称, synGroupId: 同义词组ID, options: 更新选项, terms: 术语列表
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

// FTTagVals 获取标签值 ctx: 上下文, index: 索引名称, field: 字段名称
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
