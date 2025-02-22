package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// GeoAdd 将指定的地理空间位置（纬度、经度、名称）添加到指定的key中
func (r *Client) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoAdd(getCtx(ctx), key, geoLocation...).Result()
		return err
	}, acceptable)
	return
}

// GeoPos 从key里返回所有给定位置元素的位置（经度和纬度）
func (r *Client) GeoPos(ctx context.Context, key string, members ...string) (val []*redis.GeoPos, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoPos(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// GeoRadius 以给定的经纬度为中心， 找出某一半径内的元素
func (r *Client) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) (val []redis.GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadius(getCtx(ctx), key, longitude, latitude, query).Result()
		return err
	}, acceptable)
	return
}

// GeoRadiusStore 找出位于指定范围内的元素，中心点是由给定的经纬度决定 ctx: 上下文 key: 键名 longitude: 经度 latitude: 纬度 query: 查询条件
func (r *Client) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadiusStore(getCtx(ctx), key, longitude, latitude, query).Result()
		return err
	}, acceptable)
	return
}

// GeoRadiusByMember 找出位于指定范围内的元素，中心点是由给定的位置元素决定 ctx: 上下文 key: 键名 member: 位置元素 query: 查询条件
func (r *Client) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) (val []redis.GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadiusByMember(getCtx(ctx), key, member, query).Result()
		return err
	}, acceptable)
	return
}

// GeoRadiusByMemberStore 找出位于指定范围内的元素，中心点是由给定的位置元素决定
func (r *Client) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadiusByMemberStore(getCtx(ctx), key, member, query).Result()
		return err
	}, acceptable)
	return
}

// GeoSearch 在指定key中搜索地理空间位置 ctx: 上下文 key: 键名 q: 搜索查询条件
func (r *Client) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoSearch(getCtx(ctx), key, q).Result()
		return err
	}, acceptable)
	return
}

// GeoSearchLocation 在指定key中搜索地理空间位置并返回位置信息 ctx: 上下文 key: 键名 q: 搜索查询条件
func (r *Client) GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) (val []redis.GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoSearchLocation(getCtx(ctx), key, q).Result()
		return err
	}, acceptable)
	return
}

// GeoSearchStore 在指定key中搜索地理空间位置并将结果存储到另一个key中 ctx: 上下文 key: 源键名 store: 目标键名 q: 搜索查询条件
func (r *Client) GeoSearchStore(ctx context.Context, key, store string, q *redis.GeoSearchStoreQuery) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoSearchStore(getCtx(ctx), key, store, q).Result()
		return err
	}, acceptable)
	return
}

// GeoDist 返回两个给定位置之间的距离 ctx: 上下文 key: 键名 member1: 位置1 member2: 位置2 unit: 距离单位(m|km|ft|mi)
func (r *Client) GeoDist(ctx context.Context, key string, member1, member2, unit string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoDist(getCtx(ctx), key, member1, member2, unit).Result()
		return err
	}, acceptable)
	return

}

// GeoHash 返回一个或多个位置元素的 Geohash 表示 ctx: 上下文 key: 键名 members: 位置元素列表
func (r *Client) GeoHash(ctx context.Context, key string, members ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoHash(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}
