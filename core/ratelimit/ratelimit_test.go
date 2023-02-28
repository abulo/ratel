package ratelimit

import (
	"reflect"
	"testing"
	"time"
)

func TestNewBucket(t *testing.T) {
	type args struct {
		fillInterval time.Duration
		capacity     int64
	}
	tests := []struct {
		name string
		args args
		want *Bucket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBucket(tt.args.fillInterval, tt.args.capacity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBucket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBucketWithClock(t *testing.T) {
	type args struct {
		fillInterval time.Duration
		capacity     int64
		clock        Clock
	}
	tests := []struct {
		name string
		args args
		want *Bucket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBucketWithClock(tt.args.fillInterval, tt.args.capacity, tt.args.clock); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBucketWithClock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBucketWithRate(t *testing.T) {
	type args struct {
		rate     float64
		capacity int64
	}
	tests := []struct {
		name string
		args args
		want *Bucket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBucketWithRate(tt.args.rate, tt.args.capacity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBucketWithRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBucketWithRateAndClock(t *testing.T) {
	type args struct {
		rate     float64
		capacity int64
		clock    Clock
	}
	tests := []struct {
		name string
		args args
		want *Bucket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBucketWithRateAndClock(tt.args.rate, tt.args.capacity, tt.args.clock); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBucketWithRateAndClock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextQuantum(t *testing.T) {
	type args struct {
		q int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextQuantum(tt.args.q); got != tt.want {
				t.Errorf("nextQuantum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBucketWithQuantum(t *testing.T) {
	type args struct {
		fillInterval time.Duration
		capacity     int64
		quantum      int64
	}
	tests := []struct {
		name string
		args args
		want *Bucket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBucketWithQuantum(tt.args.fillInterval, tt.args.capacity, tt.args.quantum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBucketWithQuantum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBucketWithQuantumAndClock(t *testing.T) {
	type args struct {
		fillInterval time.Duration
		capacity     int64
		quantum      int64
		clock        Clock
	}
	tests := []struct {
		name string
		args args
		want *Bucket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBucketWithQuantumAndClock(tt.args.fillInterval, tt.args.capacity, tt.args.quantum, tt.args.clock); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBucketWithQuantumAndClock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Wait(t *testing.T) {
	type args struct {
		count int64
	}
	tests := []struct {
		name string
		tb   *Bucket
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tb.Wait(tt.args.count)
		})
	}
}

func TestBucket_WaitMaxDuration(t *testing.T) {
	type args struct {
		count   int64
		maxWait time.Duration
	}
	tests := []struct {
		name string
		tb   *Bucket
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tb.WaitMaxDuration(tt.args.count, tt.args.maxWait); got != tt.want {
				t.Errorf("Bucket.WaitMaxDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Take(t *testing.T) {
	type args struct {
		count int64
	}
	tests := []struct {
		name string
		tb   *Bucket
		args args
		want time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tb.Take(tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bucket.Take() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_TakeMaxDuration(t *testing.T) {
	type args struct {
		count   int64
		maxWait time.Duration
	}
	tests := []struct {
		name  string
		tb    *Bucket
		args  args
		want  time.Duration
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.tb.TakeMaxDuration(tt.args.count, tt.args.maxWait)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bucket.TakeMaxDuration() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Bucket.TakeMaxDuration() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBucket_TakeAvailable(t *testing.T) {
	type args struct {
		count int64
	}
	tests := []struct {
		name string
		tb   *Bucket
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tb.TakeAvailable(tt.args.count); got != tt.want {
				t.Errorf("Bucket.TakeAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_takeAvailable(t *testing.T) {
	type args struct {
		now   time.Time
		count int64
	}
	tests := []struct {
		name string
		tb   *Bucket
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tb.takeAvailable(tt.args.now, tt.args.count); got != tt.want {
				t.Errorf("Bucket.takeAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Available(t *testing.T) {
	tests := []struct {
		name string
		tb   *Bucket
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tb.Available(); got != tt.want {
				t.Errorf("Bucket.Available() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_available(t *testing.T) {
	type args struct {
		now time.Time
	}
	tests := []struct {
		name string
		tb   *Bucket
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tb.available(tt.args.now); got != tt.want {
				t.Errorf("Bucket.available() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Capacity(t *testing.T) {
	tests := []struct {
		name string
		tb   *Bucket
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tb.Capacity(); got != tt.want {
				t.Errorf("Bucket.Capacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Rate(t *testing.T) {
	tests := []struct {
		name string
		tb   *Bucket
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tb.Rate(); got != tt.want {
				t.Errorf("Bucket.Rate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_take(t *testing.T) {
	type args struct {
		now     time.Time
		count   int64
		maxWait time.Duration
	}
	tests := []struct {
		name  string
		tb    *Bucket
		args  args
		want  time.Duration
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.tb.take(tt.args.now, tt.args.count, tt.args.maxWait)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bucket.take() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Bucket.take() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBucket_currentTick(t *testing.T) {
	type args struct {
		now time.Time
	}
	tests := []struct {
		name string
		tb   *Bucket
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tb.currentTick(tt.args.now); got != tt.want {
				t.Errorf("Bucket.currentTick() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_adjustavailableTokens(t *testing.T) {
	type args struct {
		tick int64
	}
	tests := []struct {
		name string
		tb   *Bucket
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tb.adjustavailableTokens(tt.args.tick)
		})
	}
}

func Test_realClock_Now(t *testing.T) {
	tests := []struct {
		name string
		r    realClock
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Now(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("realClock.Now() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_realClock_Sleep(t *testing.T) {
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name string
		r    realClock
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Sleep(tt.args.d)
		})
	}
}
