package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// BFAdd 添加元素到布隆过滤器 ctx: 上下文, key: 布隆过滤器键名, element: 要添加的元素
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

// BFCard 获取布隆过滤器的元素数量 ctx: 上下文, key: 布隆过滤器键名
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

// BFExists 检查元素是否存在于布隆过滤器 ctx: 上下文, key: 布隆过滤器键名, element: 要检查的元素
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

// BFInfo 获取布隆过滤器的信息 ctx: 上下文, key: 布隆过滤器键名
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

// BFInfoArg 获取布隆过滤器的信息 ctx: 上下文, key: 布隆过滤器键名, option: 信息选项
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

// BFInfoCapacity 获取布隆过滤器的容量 ctx: 上下文, key: 布隆过滤器键名
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

// BFInfoSize 获取布隆过滤器的大小 ctx: 上下文, key: 布隆过滤器键名
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

// BFInfoFilters 获取布隆过滤器的过滤器 ctx: 上下文, key: 布隆过滤器键名
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

// BFInfoItems 获取布隆过滤器的元素 ctx: 上下文, key: 布隆过滤器键名
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

// BFInfoExpansion 获取布隆过滤器的扩展 ctx: 上下文, key: 布隆过滤器键名
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

// BFInsert 插入元素到布隆过滤器 ctx: 上下文, key: 布隆过滤器键名, options: 插入选项, elements: 要插入的元素
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

// BFMAdd 添加多个元素到布隆过滤器 ctx: 上下文, key: 布隆过滤器键名, elements: 要添加的元素
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

// BFMExists 检查多个元素是否存在于布隆过滤器 ctx: 上下文, key: 布隆过滤器键名, elements: 要检查的元素
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

// BFReserve 创建布隆过滤器 ctx: 上下文, key: 布隆过滤器键名, errorRate: 错误率, capacity: 容量
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

// BFReserveExpansion 创建布隆过滤器 ctx: 上下文, key: 布隆过滤器键名, errorRate: 错误率, capacity: 容量, expansion: 扩展因子
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

// BFReserveNonScaling 创建布隆过滤器 ctx: 上下文, key: 布隆过滤器键名, errorRate: 错误率, capacity: 容量
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

// BFReserveWithArgs 创建布隆过滤器 ctx: 上下文, key: 布隆过滤器键名, options: 创建选项
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

// BFScanDump 获取布隆过滤器的迭代器 ctx: 上下文, key: 布隆过滤器键名, iterator: 迭代器
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

// BFLoadChunk 加载布隆过滤器的块 ctx: 上下文, key: 布隆过滤器键名, iterator: 迭代器, data: 数据块
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

// CFAdd 添加元素到计数过滤器 ctx: 上下文, key: 计数过滤器键名, element: 要添加的元素
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

// CFAddNX 添加元素到计数过滤器 ctx: 上下文, key: 计数过滤器键名, element: 要添加的元素
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

// CFCount 获取计数过滤器的元素数量 ctx: 上下文, key: 计数过滤器键名, element: 要计数的元素
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

// CFDel 删除计数过滤器的元素 ctx: 上下文, key: 计数过滤器键名, element: 要删除的元素
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

// CFExists 检查元素是否存在于计数过滤器 ctx: 上下文, key: 计数过滤器键名, element: 要检查的元素
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

// CFInfo 获取计数过滤器的信息 ctx: 上下文, key: 计数过滤器键名
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

// CFInsert 插入元素到计数过滤器 ctx: 上下文, key: 计数过滤器键名, options: 插入选项, elements: 要插入的元素
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

// CFInsertNX 插入元素到计数过滤器 ctx: 上下文, key: 计数过滤器键名, options: 插入选项, elements: 要插入的元素
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

// CFMExists 检查多个元素是否存在于计数过滤器 ctx: 上下文, key: 计数过滤器键名, elements: 要检查的元素
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

// CFReserve 创建计数过滤器 ctx: 上下文, key: 计数过滤器键名, capacity: 容量
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

// CFReserveWithArgs 创建计数过滤器 ctx: 上下文, key: 计数过滤器键名, options: 创建选项
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

// CFReserveExpansion 创建计数过滤器 ctx: 上下文, key: 计数过滤器键名, capacity: 容量, expansion: 扩展因子
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

// CFReserveBucketSize 创建计数过滤器 ctx: 上下文, key: 计数过滤器键名, capacity: 容量, bucketsize: 桶大小
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

// CFReserveMaxIterations 创建计数过滤器 ctx: 上下文, key: 计数过滤器键名, capacity: 容量, maxiterations: 最大迭代次数
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

// CFScanDump 获取计数过滤器的迭代器 ctx: 上下文, key: 计数过滤器键名, iterator: 迭代器
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

// CFLoadChunk 加载计数过滤器的块
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

// CMSIncrBy 增加元素到计数最小二乘 ctx: 上下文, key: 计数最小二乘键名, elements: 要增加的元素
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

// CMSInfo 获取计数最小二乘的信息 ctx: 上下文, key: 计数最小二乘键名
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

// CMSInitByDim 初始化计数最小二乘 ctx: 上下文, key: 计数最小二乘键名, width: 宽度, height: 高度
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

// CMSInitByProb 初始化计数最小二乘 ctx: 上下文, key: 计数最小二乘键名, errorRate: 错误率, probability: 概率
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

// CMSMerge 合并计数最小二乘 ctx: 上下文, destKey: 目标键名, sourceKeys: 源键名
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

// CMSMergeWithWeight 合并计数最小二乘 ctx: 上下文, destKey: 目标键名, sourceKeys: 源键名及其权重
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

// CMSQuery 查询计数最小二乘 ctx: 上下文, key: 计数最小二乘键名, elements: 要查询的元素
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

// TopKAdd 添加元素到TopK ctx: 上下文, key: TopK键名, elements: 要添加的元素
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

// TopKCount 获取TopK的元素数量 ctx: 上下文, key: TopK键名, elements: 要计数的元素
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

// TopKIncrBy 增加元素到TopK ctx: 上下文, key: TopK键名, elements: 要增加的元素
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

// TopKInfo 获取TopK的信息 ctx: 上下文, key: TopK键名
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

// TopKList 获取TopK的列表 ctx: 上下文, key: TopK键名
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

// TopKListWithCount 获取TopK的列表 ctx: 上下文, key: TopK键名
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

// TopKQuery 查询TopK ctx: 上下文, key: TopK键名, elements: 要查询的元素
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

// TopKReserve 创建TopK ctx: 上下文, key: TopK键名, k: 保留元素数量
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

// TopKReserveWithOptions 创建TopK ctx: 上下文, key: TopK键名, k: 保留元素数量, width: 宽度, depth: 深度, decay: 衰减因子
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

// TDigestAdd 添加元素到TDigest ctx: 上下文, key: TDigest键名, elements: 要添加的元素
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

// TDigestByRank 根据排名获取TDigest值 ctx: 上下文, key: TDigest键名, rank: 排名
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

// TDigestByRevRank 根据反向排名获取TDigest值 ctx: 上下文, key: TDigest键名, rank: 反向排名
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

// TDigestCDF 获取TDigest的累积分布函数值 ctx: 上下文, key: TDigest键名, elements: 要查询的值
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

// TDigestCreate 创建TDigest ctx: 上下文, key: TDigest键名
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

// TDigestCreateWithCompression 创建TDigest ctx: 上下文, key: TDigest键名, compression: 压缩参数
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

// TDigestInfo 获取TDigest的信息 ctx: 上下文, key: TDigest键名
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

// TDigestMax 获取TDigest的最大值 ctx: 上下文, key: TDigest键名
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

// TDigestMin 获取TDigest的最小值 ctx: 上下文, key: TDigest键名
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

// TDigestMerge 合并TDigest ctx: 上下文, destKey: 目标键名, options: 合并选项, sourceKeys: 源键名
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

// TDigestQuantile 获取TDigest的分位数 ctx: 上下文, key: TDigest键名, elements: 分位数值
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

// TDigestRank 获取TDigest的排名 ctx: 上下文, key: TDigest键名, values: 要查询的值
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

// TDigestReset 重置TDigest ctx: 上下文, key: TDigest键名
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

// TDigestRevRank 获取TDigest的反向排名 ctx: 上下文, key: TDigest键名, values: 要查询的值
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

// TDigestTrimmedMean 获取TDigest的修剪均值 ctx: 上下文, key: TDigest键名, lowCutQuantile: 下分位数, highCutQuantile: 上分位数
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
