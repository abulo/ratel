package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// GeoAdd 将指定的地理空间位置（纬度、经度、名称）添加到指定的key中
func (r *Client) GeoAdd(ctx context.Context, key any, geoLocation ...*redis.GeoLocation) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoAdd(getCtx(ctx), r.k(key), geoLocation...).Result()
		return err
	}, acceptable)
	return
}

// GeoPos 从key里返回所有给定位置元素的位置（经度和纬度）
func (r *Client) GeoPos(ctx context.Context, key any, members ...string) (val []*redis.GeoPos, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoPos(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// GeoRadius 以给定的经纬度为中心， 找出某一半径内的元素
func (r *Client) GeoRadius(ctx context.Context, key any, longitude, latitude float64, query *redis.GeoRadiusQuery) (val []redis.GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadius(getCtx(ctx), r.k(key), longitude, latitude, query).Result()
		return err
	}, acceptable)
	return
}

// GeoRadiusStore -> GeoRadius
func (r *Client) GeoRadiusStore(ctx context.Context, key any, longitude, latitude float64, query *redis.GeoRadiusQuery) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadiusStore(getCtx(ctx), r.k(key), longitude, latitude, query).Result()
		return err
	}, acceptable)
	return
}

// GeoRadiusByMember -> GeoRadius
func (r *Client) GeoRadiusByMember(ctx context.Context, key any, member string, query *redis.GeoRadiusQuery) (val []redis.GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadiusByMember(getCtx(ctx), r.k(key), member, query).Result()
		return err
	}, acceptable)
	return
}

// GeoRadiusByMemberStore 找出位于指定范围内的元素，中心点是由给定的位置元素决定
func (r *Client) GeoRadiusByMemberStore(ctx context.Context, key any, member string, query *redis.GeoRadiusQuery) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadiusByMemberStore(getCtx(ctx), r.k(key), member, query).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) GeoSearch(ctx context.Context, key any, q *redis.GeoSearchQuery) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoSearch(getCtx(ctx), r.k(key), q).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) GeoSearchLocation(ctx context.Context, key any, q *redis.GeoSearchLocationQuery) (val []redis.GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoSearchLocation(getCtx(ctx), r.k(key), q).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) GeoSearchStore(ctx context.Context, key any, store string, q *redis.GeoSearchStoreQuery) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoSearchStore(getCtx(ctx), r.k(key), store, q).Result()
		return err
	}, acceptable)
	return
}

// GeoDist 返回两个给定位置之间的距离
func (r *Client) GeoDist(ctx context.Context, key any, member1, member2, unit string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoDist(getCtx(ctx), r.k(key), member1, member2, unit).Result()
		return err
	}, acceptable)
	return

}

// GeoHash 返回一个或多个位置元素的 Geohash 表示
func (r *Client) GeoHash(ctx context.Context, key any, members ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoHash(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}
