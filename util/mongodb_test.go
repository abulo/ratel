package util

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestPipeline(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want mongo.Pipeline
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pipeline(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertBson(t *testing.T) {
	type args struct {
		items bson.D
		item  bson.E
	}
	tests := []struct {
		name string
		args args
		want bson.D
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertBson(tt.args.items, tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertBson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indirect(t *testing.T) {
	type args struct {
		a interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := indirect(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("indirect() = %v, want %v", got, tt.want)
			}
		})
	}
}
