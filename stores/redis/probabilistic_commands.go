package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// BFAdd 将 element 加入布隆过滤器 key；首次加入返回 true，重复加入返回 false。
func (r *Client) BFAdd(ctx context.Context, key string, element any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFAdd(getCtx(ctx), key, element).Result()
		return err
	}, acceptable)
	return
}

// BFCard 返回布隆过滤器 key 中已加入元素的近似数量。
func (r *Client) BFCard(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFCard(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// BFExists 判断 element 是否可能存在于布隆过滤器 key；可能误报但不会漏报。
func (r *Client) BFExists(ctx context.Context, key string, element any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFExists(getCtx(ctx), key, element).Result()
		return err
	}, acceptable)
	return
}

// BFInfo 返回布隆过滤器 key 的容量、过滤器数量等统计信息。
func (r *Client) BFInfo(ctx context.Context, key string) (val redis.BFInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFInfo(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// BFInfoArg 按 option 指定字段返回布隆过滤器 key 的信息。
func (r *Client) BFInfoArg(ctx context.Context, key, option string) (val redis.BFInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFInfoArg(getCtx(ctx), key, option).Result()
		return err
	}, acceptable)
	return
}

// BFInfoCapacity 返回布隆过滤器 key 的配置容量。
func (r *Client) BFInfoCapacity(ctx context.Context, key string) (val redis.BFInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFInfoCapacity(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// BFInfoSize 返回布隆过滤器 key 当前占用的内存大小。
func (r *Client) BFInfoSize(ctx context.Context, key string) (val redis.BFInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFInfoSize(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// BFInfoFilters 返回布隆过滤器 key 内部过滤器（子数组）数量。
func (r *Client) BFInfoFilters(ctx context.Context, key string) (val redis.BFInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFInfoFilters(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// BFInfoItems 返回布隆过滤器 key 中已记录元素的近似数量。
func (r *Client) BFInfoItems(ctx context.Context, key string) (val redis.BFInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFInfoItems(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// BFInfoExpansion 返回布隆过滤器 key 的扩容因子。
func (r *Client) BFInfoExpansion(ctx context.Context, key string) (val redis.BFInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFInfoExpansion(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// BFInsert 按 BFInsertOptions 将 elements 批量插入布隆过滤器 key，返回各元素是否为新增。
func (r *Client) BFInsert(ctx context.Context, key string, options *redis.BFInsertOptions, elements ...any) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFInsert(getCtx(ctx), key, options, elements...).Result()
		return err
	}, acceptable)
	return
}

// BFMAdd 将 elements 批量加入布隆过滤器 key，返回各元素是否为新增。
func (r *Client) BFMAdd(ctx context.Context, key string, elements ...any) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFMAdd(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// BFMExists 批量判断 elements 是否可能存在于布隆过滤器 key，返回布尔切片顺序与 elements 一致。
func (r *Client) BFMExists(ctx context.Context, key string, elements ...any) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFMExists(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// BFReserve 以给定错误率 errorRate 与容量 capacity 创建布隆过滤器 key；键已存在时报错。
func (r *Client) BFReserve(ctx context.Context, key string, errorRate float64, capacity int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFReserve(getCtx(ctx), key, errorRate, capacity).Result()
		return err
	}, acceptable)
	return
}

// BFReserveExpansion 以错误率、容量与扩容因子 expansion 创建可扩容布隆过滤器 key。
func (r *Client) BFReserveExpansion(ctx context.Context, key string, errorRate float64, capacity, expansion int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFReserveExpansion(getCtx(ctx), key, errorRate, capacity, expansion).Result()
		return err
	}, acceptable)
	return
}

// BFReserveNonScaling 创建不可扩容的布隆过滤器 key，容量达到上限后不再自动增长。
func (r *Client) BFReserveNonScaling(ctx context.Context, key string, errorRate float64, capacity int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFReserveNonScaling(getCtx(ctx), key, errorRate, capacity).Result()
		return err
	}, acceptable)
	return
}

// BFReserveWithArgs 按 BFReserveOptions 创建布隆过滤器 key。
func (r *Client) BFReserveWithArgs(ctx context.Context, key string, options *redis.BFReserveOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFReserveWithArgs(getCtx(ctx), key, options).Result()
		return err
	}, acceptable)
	return
}

// BFScanDump 以 iterator 游标分批导出布隆过滤器 key 的内部数据块，供跨实例恢复使用。
func (r *Client) BFScanDump(ctx context.Context, key string, iterator int64) (val redis.ScanDump, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFScanDump(getCtx(ctx), key, iterator).Result()
		return err
	}, acceptable)
	return
}

// BFLoadChunk 将 BFScanDump 导出的数据块按 iterator 写回布隆过滤器 key。
func (r *Client) BFLoadChunk(ctx context.Context, key string, iterator int64, data any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BFLoadChunk(getCtx(ctx), key, iterator, data).Result()
		return err
	}, acceptable)
	return
}

// CFAdd 将 element 加入计数布隆过滤器（Cuckoo Filter）key，返回是否插入成功。
func (r *Client) CFAdd(ctx context.Context, key string, element any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFAdd(getCtx(ctx), key, element).Result()
		return err
	}, acceptable)
	return
}

// CFAddNX 仅当 element 不存在时加入 Cuckoo Filter key，已存在则不重复插入。
func (r *Client) CFAddNX(ctx context.Context, key string, element any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFAddNX(getCtx(ctx), key, element).Result()
		return err
	}, acceptable)
	return
}

// CFCount 返回 element 在 Cuckoo Filter key 中的近似计数。
func (r *Client) CFCount(ctx context.Context, key string, element any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFCount(getCtx(ctx), key, element).Result()
		return err
	}, acceptable)
	return
}

// CFDel 从 Cuckoo Filter key 中移除 element，返回是否移除成功。
func (r *Client) CFDel(ctx context.Context, key string, element any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFDel(getCtx(ctx), key, element).Result()
		return err
	}, acceptable)
	return
}

// CFExists 判断 element 是否可能存在于 Cuckoo Filter key。
func (r *Client) CFExists(ctx context.Context, key string, element any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFExists(getCtx(ctx), key, element).Result()
		return err
	}, acceptable)
	return
}

// CFInfo 返回 Cuckoo Filter key 的大小、桶数、容量等统计信息。
func (r *Client) CFInfo(ctx context.Context, key string) (val redis.CFInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFInfo(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// CFInsert 按 CFInsertOptions 将 elements 批量插入 Cuckoo Filter key，返回各元素插入结果。
func (r *Client) CFInsert(ctx context.Context, key string, options *redis.CFInsertOptions, elements ...any) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFInsert(getCtx(ctx), key, options, elements...).Result()
		return err
	}, acceptable)
	return
}

// CFInsertNX 仅插入不存在的元素，返回各元素插入前后的计数。
func (r *Client) CFInsertNX(ctx context.Context, key string, options *redis.CFInsertOptions, elements ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFInsertNX(getCtx(ctx), key, options, elements...).Result()
		return err
	}, acceptable)
	return
}

// CFMExists 批量判断 elements 是否可能存在于 Cuckoo Filter key，返回布尔切片顺序与 elements 一致。
func (r *Client) CFMExists(ctx context.Context, key string, elements ...any) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFMExists(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// CFReserve 以容量 capacity 创建 Cuckoo Filter key；键已存在时报错。
func (r *Client) CFReserve(ctx context.Context, key string, capacity int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFReserve(getCtx(ctx), key, capacity).Result()
		return err
	}, acceptable)
	return
}

// CFReserveWithArgs 按 CFReserveOptions 创建 Cuckoo Filter key。
func (r *Client) CFReserveWithArgs(ctx context.Context, key string, options *redis.CFReserveOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFReserveWithArgs(getCtx(ctx), key, options).Result()
		return err
	}, acceptable)
	return
}

// CFReserveExpansion 以容量与扩容因子 expansion 创建可扩容 Cuckoo Filter key。
func (r *Client) CFReserveExpansion(ctx context.Context, key string, capacity, expansion int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFReserveExpansion(getCtx(ctx), key, capacity, expansion).Result()
		return err
	}, acceptable)
	return
}

// CFReserveBucketSize 以容量与每桶元素数 bucketsize 创建 Cuckoo Filter key。
func (r *Client) CFReserveBucketSize(ctx context.Context, key string, capacity, bucketsize int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFReserveBucketSize(getCtx(ctx), key, capacity, bucketsize).Result()
		return err
	}, acceptable)
	return
}

// CFReserveMaxIterations 以容量与最大迭代次数 maxiterations 创建 Cuckoo Filter key。
func (r *Client) CFReserveMaxIterations(ctx context.Context, key string, capacity, maxiterations int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFReserveMaxIterations(getCtx(ctx), key, capacity, maxiterations).Result()
		return err
	}, acceptable)
	return
}

// CFScanDump 以 iterator 游标分批导出 Cuckoo Filter key 的内部数据块。
func (r *Client) CFScanDump(ctx context.Context, key string, iterator int64) (val redis.ScanDump, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFScanDump(getCtx(ctx), key, iterator).Result()
		return err
	}, acceptable)
	return
}

// CFLoadChunk 将 CFScanDump 导出的数据块按 iterator 写回 Cuckoo Filter key。
func (r *Client) CFLoadChunk(ctx context.Context, key string, iterator int64, data any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CFLoadChunk(getCtx(ctx), key, iterator, data).Result()
		return err
	}, acceptable)
	return
}

// CMSIncrBy 将 elements 批量加入计数草图（Count-Min Sketch）key 并累加计数，返回各元素更新后的计数。
func (r *Client) CMSIncrBy(ctx context.Context, key string, elements ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CMSIncrBy(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// CMSInfo 返回 Count-Min Sketch key 的宽度、高度与已记录元素统计信息。
func (r *Client) CMSInfo(ctx context.Context, key string) (val redis.CMSInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CMSInfo(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// CMSInitByDim 以宽度 width、高度 height 初始化 Count-Min Sketch key。
func (r *Client) CMSInitByDim(ctx context.Context, key string, width, height int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CMSInitByDim(getCtx(ctx), key, width, height).Result()
		return err
	}, acceptable)
	return
}

// CMSInitByProb 按误差 errorRate 与概率 probability 初始化 Count-Min Sketch key。
func (r *Client) CMSInitByProb(ctx context.Context, key string, errorRate, probability float64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CMSInitByProb(getCtx(ctx), key, errorRate, probability).Result()
		return err
	}, acceptable)
	return
}

// CMSMerge 将多个 sourceKeys 的 Count-Min Sketch 合并到 destkey。
func (r *Client) CMSMerge(ctx context.Context, destkey string, sourceKeys ...string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CMSMerge(getCtx(ctx), destkey, sourceKeys...).Result()
		return err
	}, acceptable)
	return
}

// CMSMergeWithWeight 按 sourceKeys 中给出的权重加权合并到 destkey。
func (r *Client) CMSMergeWithWeight(ctx context.Context, destkey string, sourceKeys map[string]int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CMSMergeWithWeight(getCtx(ctx), destkey, sourceKeys).Result()
		return err
	}, acceptable)
	return
}

// CMSQuery 查询 elements 在 Count-Min Sketch key 中的近似计数，返回切片顺序与 elements 一致。
func (r *Client) CMSQuery(ctx context.Context, key string, elements ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CMSQuery(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// TopKAdd 将 elements 加入 Top-K 草图 key，返回被挤出 Top-K 名单的元素。
func (r *Client) TopKAdd(ctx context.Context, key string, elements ...any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TopKAdd(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// TopKCount 查询 elements 在 Top-K 草图 key 中的近似计数。
func (r *Client) TopKCount(ctx context.Context, key string, elements ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TopKCount(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// TopKIncrBy 按 elements 给出的元素与增量更新 Top-K 草图 key，返回被挤出名单的元素。
func (r *Client) TopKIncrBy(ctx context.Context, key string, elements ...any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TopKIncrBy(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// TopKInfo 返回 Top-K 草图 key 的 k、宽度、深度、衰减因子等配置信息。
func (r *Client) TopKInfo(ctx context.Context, key string) (val redis.TopKInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TopKInfo(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TopKList 返回当前 Top-K 草图 key 维护的高频元素列表。
func (r *Client) TopKList(ctx context.Context, key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TopKList(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TopKListWithCount 返回高频元素及其近似计数组成的 map。
func (r *Client) TopKListWithCount(ctx context.Context, key string) (val map[string]int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TopKListWithCount(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TopKQuery 批量判断 elements 是否位于 Top-K 名单中，返回布尔切片顺序与 elements 一致。
func (r *Client) TopKQuery(ctx context.Context, key string, elements ...any) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TopKQuery(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// TopKReserve 创建保留前 k 个高频元素的 Top-K 草图 key。
func (r *Client) TopKReserve(ctx context.Context, key string, k int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TopKReserve(getCtx(ctx), key, k).Result()
		return err
	}, acceptable)
	return
}

// TopKReserveWithOptions 按 k、宽度 width、深度 depth 与衰减因子 decay 创建 Top-K 草图 key。
func (r *Client) TopKReserveWithOptions(ctx context.Context, key string, k, width, depth int64, decay float64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TopKReserveWithOptions(getCtx(ctx), key, k, width, depth, decay).Result()
		return err
	}, acceptable)
	return
}

// TDigestAdd 将 elements 样本加入 T-Digest 近似分位数结构 key。
func (r *Client) TDigestAdd(ctx context.Context, key string, elements ...float64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestAdd(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// TDigestByRank 返回给定 rank 排名对应的样本值（升序），切片顺序与 rank 一致。
func (r *Client) TDigestByRank(ctx context.Context, key string, rank ...uint64) (val []float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestByRank(getCtx(ctx), key, rank...).Result()
		return err
	}, acceptable)
	return
}

// TDigestByRevRank 返回给定反向排名对应的样本值（降序）。
func (r *Client) TDigestByRevRank(ctx context.Context, key string, rank ...uint64) (val []float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestByRevRank(getCtx(ctx), key, rank...).Result()
		return err
	}, acceptable)
	return
}

// TDigestCDF 返回给定样本值在 T-Digest 中的累积分布概率。
func (r *Client) TDigestCDF(ctx context.Context, key string, elements ...float64) (val []float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		cmd := conn.TDigestCDF(getCtx(ctx), key, elements...)
		val, err = cmd.Result()
		return err
	}, acceptable)
	return
}

// TDigestCreate 创建 T-Digest 结构 key，使用默认压缩参数。
func (r *Client) TDigestCreate(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		cmd := conn.TDigestCreate(getCtx(ctx), key)
		val, err = cmd.Result()
		return err
	}, acceptable)
	return
}

// TDigestCreateWithCompression 以压缩参数 compression 创建 T-Digest key，压缩值越大精度越高、占用内存越大。
func (r *Client) TDigestCreateWithCompression(ctx context.Context, key string, compression int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		cmd := conn.TDigestCreateWithCompression(getCtx(ctx), key, compression)
		val, err = cmd.Result()
		return err
	}, acceptable)
	return
}

// TDigestInfo 返回 T-Digest key 的压缩参数、已记录样本数量等信息。
func (r *Client) TDigestInfo(ctx context.Context, key string) (val redis.TDigestInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TDigestInfo(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TDigestMax 返回 T-Digest key 中记录的最大样本值。
func (r *Client) TDigestMax(ctx context.Context, key string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestMax(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TDigestMin 返回 T-Digest key 中记录的最小样本值。
func (r *Client) TDigestMin(ctx context.Context, key string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestMin(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TDigestMerge 按 TDigestMergeOptions 将多个 sourceKeys 的 T-Digest 合并到 destkey。
func (r *Client) TDigestMerge(ctx context.Context, destkey string, options *redis.TDigestMergeOptions, sourceKeys ...string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestMerge(getCtx(ctx), destkey, options, sourceKeys...).Result()
		return err
	}, acceptable)
	return
}

// TDigestQuantile 返回给定分位数（0~1）对应的样本值，切片顺序与 elements 一致。
func (r *Client) TDigestQuantile(ctx context.Context, key string, elements ...float64) (val []float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestQuantile(getCtx(ctx), key, elements...).Result()
		return err
	}, acceptable)
	return
}

// TDigestRank 返回给定样本值在 T-Digest 中的升序排名。
func (r *Client) TDigestRank(ctx context.Context, key string, values ...float64) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TDigestRank(getCtx(ctx), key, values...).Result()
		return err
	}, acceptable)
	return
}

// TDigestReset 清空 T-Digest key 中已记录的全部样本。
func (r *Client) TDigestReset(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestReset(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TDigestRevRank 返回给定样本值在 T-Digest 中的降序排名。
func (r *Client) TDigestRevRank(ctx context.Context, key string, values ...float64) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestRevRank(getCtx(ctx), key, values...).Result()
		return err
	}, acceptable)
	return
}

// TDigestTrimmedMean 剔除两端 lowCutQuantile、highCutQuantile 分位样本后返回中间部分的均值。
func (r *Client) TDigestTrimmedMean(ctx context.Context, key string, lowCutQuantile, highCutQuantile float64) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TDigestTrimmedMean(getCtx(ctx), key, lowCutQuantile, highCutQuantile).Result()
		return err
	}, acceptable)
	return
}
