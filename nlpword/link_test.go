package nlpword

import (
	"reflect"
	"testing"
)

func TestLinkList_Push(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		list *LinkList
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.list.Push(tt.args.v)
		})
	}
}

func TestLinkList_Pop(t *testing.T) {
	tests := []struct {
		name string
		list *LinkList
		want any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.Pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LinkList.Pop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkList_Empty(t *testing.T) {
	tests := []struct {
		name string
		list *LinkList
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.Empty(); got != tt.want {
				t.Errorf("LinkList.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}
