package pagination

import (
	"bytes"
	"math"
	"strconv"

	"github.com/abulo/ratel/v2/util"
)

type Pager struct {
	totalItems     int64 //总记录数
	currentPage    int64 //当前页码
	itemsPerPage   int64 //每页多少数据
	numPages       int64 //总页数
	maxPagesToShow int64 //
	Url            string
}

type page struct {
	num     string
	url     string
	current bool
}

// NewPage 构造实列
func NewPage(items, curPage, perNum int64, url string) *Pager {
	pager := &Pager{
		totalItems:     items,
		currentPage:    curPage,
		itemsPerPage:   perNum,
		Url:            url,
		maxPagesToShow: 9,
	}

	pager.numPages = int64(math.Ceil(float64(pager.totalItems) / float64(pager.itemsPerPage)))

	return pager
}

//SetMaxPagesToShow 设置最大页面显示
func (pager *Pager) SetMaxPagesToShow(maxPagesToShow int64) {
	pager.maxPagesToShow = maxPagesToShow
}

//HTML 转 html
func (pager *Pager) HTML() string {
	//总页数

	data := pager.getPages()
	var pageStr = bytes.Buffer{}
	pageStr.WriteString("<div class='layui-box layui-laypage layui-laypage-default' style='margin:0px;'>")

	for _, item := range data {
		if item.url != "" {
			if item.current {
				pageStr.WriteString("<span class='layui-laypage-curr' ><em class='layui-laypage-em'></em><em>")
				pageStr.WriteString(item.num)
				pageStr.WriteString("</em></span>")
			} else {
				pageStr.WriteString("<a href='" + item.url + "'>" + item.num + "</a>")
			}
		} else {
			pageStr.WriteString("<span>" + item.num + "</span>")
		}
	}

	pageStr.WriteString("<span class='layui-laypage-skip'>")
	pageStr.WriteString("转到<input type='text' id='page_num' value='" + strconv.FormatInt(pager.currentPage, 10) + "' class='layui-input'>页")
	pageStr.WriteString("/每页<input type='text' id='per_num' value='" + strconv.FormatInt(pager.itemsPerPage, 10) + "' class='layui-input'>条")
	pageStr.WriteString("<button type='button' class='layui-laypage-btn' id='layui-laypage-btn'>GO</button>")
	pageStr.WriteString("</span>")
	pageStr.WriteString("</div>")
	return pageStr.String()
}

func (pager *Pager) getPages() []page {
	var data []page
	if pager.numPages <= 1 {
		return data
	}
	if pager.numPages <= pager.maxPagesToShow {
		var i int64
		for i = 1; i <= pager.numPages; i++ {
			data = append(data, pager.createPage(strconv.FormatInt(i, 10), i == pager.currentPage))
		}
	} else {
		numAdjacents := int64((pager.maxPagesToShow - 3) / 2)
		var slidingStart int64
		var slidingEnd int64
		if (pager.currentPage + numAdjacents) > pager.numPages {
			slidingStart = pager.numPages - pager.maxPagesToShow + 2
		} else {
			slidingStart = pager.currentPage - numAdjacents
		}
		if slidingStart < 2 {
			slidingStart = 2
		}
		slidingEnd = slidingStart + pager.maxPagesToShow - 3

		if slidingEnd >= pager.numPages {
			slidingEnd = pager.numPages - 1
		}
		data = append(data, pager.createPage("1", 1 == pager.currentPage))

		if slidingStart > 2 {
			data = append(data, pager.createPageEllipsis())
		}
		for i := slidingStart; i <= slidingEnd; i++ {
			data = append(data, pager.createPage(strconv.FormatInt(i, 10), i == pager.currentPage))
		}
		if slidingEnd < (pager.numPages - 1) {
			data = append(data, pager.createPageEllipsis())
		}
		data = append(data, pager.createPage(strconv.FormatInt(pager.numPages, 10), pager.numPages == pager.currentPage))
	}

	return data

}

func (pager *Pager) createPage(pageNum string, current bool) page {
	return page{
		num:     pageNum,
		current: current,
		url:     pager.getPageURL(pageNum),
	}
}

func (pager *Pager) createPageEllipsis() page {

	return page{
		num:     "...",
		url:     "",
		current: false,
	}
}

func (pager *Pager) getPageURL(pageNum string) string {
	return util.StrReplace(":num", pageNum, pager.Url, -1)
}
