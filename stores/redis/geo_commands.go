package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// GeoAdd 将一组地理坐标（经度、纬度、成员名）写入有序集合 key，返回新增成员数量；已存在成员更新坐标。
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

// GeoPos 返回给定 members 的经纬度坐标切片；不存在的成员对应位置为 nil。
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

// GeoRadius 以给定经纬度为中心，按 query 条件查询半径范围内的成员及其距离等信息。
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

// GeoRadiusStore 按给定经纬度查询半径范围内成员，并将结果存入 query 指定的目标键。
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

// GeoRadiusByMember 以成员 member 的坐标为中心点，按 query 条件查询其范围内的成员。
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

// GeoRadiusByMemberStore 以成员 member 的坐标为中心查询范围内成员，并将结果存入指定目标键。
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

// GeoSearch 按 GeoSearchQuery 在 key 中搜索满足条件的地理成员，返回成员名称。
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

// GeoSearchLocation 按条件搜索地理成员，并返回成员名称、坐标与距离等完整位置信息。
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

// GeoSearchStore 按条件搜索 key 中的地理成员，并将结果存入 store 目标键，返回命中成员数量。
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

// GeoDist 返回成员 member1 与 member2 之间的距离，unit 指定单位（m/km/ft/mi）；任一成员不存在时返回 nil。
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

// GeoHash 返回给定 members 的 Geohash 编码字符串切片，顺序与 members 一致。
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
