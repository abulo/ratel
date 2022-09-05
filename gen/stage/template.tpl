package {{Package}}

import (
	"cloud/dao"
	"cloud/module/{{Module}}"
	"net/http"

	"github.com/abulo/ratel/v3/gin"
	"github.com/spf13/cast"
)

// Create 创建数据
func Create(ctx *gin.Context) {
	//实例化
	var obj {{Dao}}
	//绑定解析
	if err := ctx.ShouldBind(&obj); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	//调用
	if _, err := {{Module}}.Create(ctx, obj); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	//成功
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "创建成功",
	})
}

// Update 更新数据
func Update(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	//实例化
	var obj {{Dao}}
	//绑定解析
	if err := ctx.ShouldBind(&obj); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	//调用
	if _, err := {{Module}}.Update(ctx, id, obj); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	//成功
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "更新成功",
	})
}

// Show 获取数据
func Show(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	//实例化
	data, err := {{Module}}.Show(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	//成功
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": data,
	})
}

// Delete 删除数据
func Delete(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	//实例化
	if _, err := {{Module}}.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	//成功
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "删除成功",
	})
}

// List 获取列表
func List(ctx *gin.Context) {
	//当前页面
	pageNum := cast.ToInt64(ctx.DefaultQuery("pageNum", "1"))
	//每页多少数据
	pageSize := cast.ToInt64(ctx.DefaultQuery("pageSize", "10"))
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := pageSize * (pageNum - 1)
	//查询条件
	var condition map[string]interface{}
	total, err := {{Module}}.Total(ctx, condition)
	if err != nil || total < 1 {
		//返回数据
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "",
			"data": gin.H{
				"list":  dao.NilSilce(),
				"total": 0,
			},
		})
		return
	}
	list, err := {{Module}}.List(ctx, condition, offset, pageSize)
	if err != nil {
		//返回数据
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "",
			"data": gin.H{
				"list":  dao.NilSilce(),
				"total": total,
			},
		})
		return
	}
	//返回数据
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": gin.H{
			"list":  list,
			"total": total,
		},
	})
}
