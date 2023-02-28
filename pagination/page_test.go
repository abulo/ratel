package pagination

import (
	"reflect"
	"testing"
)

func TestNewPage(t *testing.T) {
	type args struct {
		items   int64
		curPage int64
		perNum  int64
		url     string
	}
	tests := []struct {
		name string
		args args
		want *Pager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPage(tt.args.items, tt.args.curPage, tt.args.perNum, tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPager_SetMaxPagesToShow(t *testing.T) {
	type args struct {
		maxPagesToShow int64
	}
	tests := []struct {
		name  string
		pager *Pager
		args  args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pager.SetMaxPagesToShow(tt.args.maxPagesToShow)
		})
	}
}

func TestPager_HTML(t *testing.T) {
	tests := []struct {
		name  string
		pager *Pager
		want  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pager.HTML(); got != tt.want {
				t.Errorf("Pager.HTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPager_getPages(t *testing.T) {
	tests := []struct {
		name  string
		pager *Pager
		want  []page
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pager.getPages(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pager.getPages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPager_createPage(t *testing.T) {
	type args struct {
		pageNum string
		current bool
	}
	tests := []struct {
		name  string
		pager *Pager
		args  args
		want  page
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pager.createPage(tt.args.pageNum, tt.args.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pager.createPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPager_createPageEllipsis(t *testing.T) {
	tests := []struct {
		name  string
		pager *Pager
		want  page
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pager.createPageEllipsis(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pager.createPageEllipsis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPager_getPageURL(t *testing.T) {
	type args struct {
		pageNum string
	}
	tests := []struct {
		name  string
		pager *Pager
		args  args
		want  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pager.getPageURL(tt.args.pageNum); got != tt.want {
				t.Errorf("Pager.getPageURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
