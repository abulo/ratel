package ratelimit

import "io"

type reader struct {
	r      io.Reader
	bucket *Bucket
}

// Reader 返回由给定令牌桶限流的 io.Reader，桶中每个令牌代表 1 个字节；
// 限流在本次读取完成后按实际读到的字节数等待令牌，读取错误原样透传。
func Reader(r io.Reader, bucket *Bucket) io.Reader {
	return &reader{
		r:      r,
		bucket: bucket,
	}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}
	r.bucket.Wait(int64(n))
	return n, err
}

type writer struct {
	w      io.Writer
	bucket *Bucket
}

// Writer 返回由给定令牌桶限流的 io.Writer，桶中每个令牌代表 1 个字节；
// 写出前先等待与缓冲区长度对应的令牌，再委托底层 Writer 执行写出。
func Writer(w io.Writer, bucket *Bucket) io.Writer {
	return &writer{
		w:      w,
		bucket: bucket,
	}
}

func (w *writer) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}
