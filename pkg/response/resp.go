package response

import (
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
)

// http请求响应封装
type Resp struct {
	Code int         `json:"code"` // 错误代码代码
	Data interface{} `json:"data"` // 数据内容
	Msg  string      `json:"msg"`  // 消息提示
}

// 分页封装
type PageInfo struct {
	PageNum      uint `json:"pageNum" form:"pageNum"`           // 当前页码
	PageSize     uint `json:"pageSize" form:"pageSize"`         // 每页显示条数
	Total        uint `json:"total"`                            // 数据总条数
	NoPagination bool `json:"noPagination" form:"noPagination"` // 不使用分页
}

// 带分页数据封装
type PageData struct {
	PageInfo
	List interface{} `json:"list"` // 数据列表
}

// 计算limit/offset, 如果需要用到返回的PageSize, PageNum, 务必保证Total值有效
func (s *PageInfo) GetLimit() (limit uint, offset uint) {
	// 传入参数可能不合法, 设置默认值
	// 每页显示条数不能小于1
	if s.PageSize < 1 {
		s.PageSize = 10
	}
	// 页码不能小于1
	if s.PageNum < 1 {
		s.PageNum = 1
	}

	// 如果偏移量比总条数还多
	if s.Total > 0 {
		if s.PageSize > s.Total {
			s.PageSize = s.Total
		}
		if s.PageNum > s.Total {
			s.PageNum = s.Total
		}
	}

	// 计算最大页码
	maxPageNum := s.Total/s.PageSize + 1
	if s.Total%s.PageSize == 0 {
		maxPageNum = s.Total / s.PageSize
	}
	// 页码不能小于1
	if maxPageNum < 1 {
		maxPageNum = 1
	}

	// 超出最后一页
	if s.PageNum > maxPageNum {
		s.PageNum = maxPageNum
	}

	limit = s.PageSize
	offset = limit * (s.PageNum - 1)
	return
}

func Result(code int, msg string, data interface{}) {
	// 结果以panic异常的形式抛出, 交由异常处理中间件处理
	panic(Resp{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func Success() {
	Result(Ok, CustomError[Ok], map[string]interface{}{})
}

func SuccessWithData(data interface{}) {
	Result(Ok, CustomError[Ok], data)
}

func SuccessWithMsg(msg string) {
	Result(Ok, msg, map[string]interface{}{})
}

func Fail() {
	FailWithCode(NotOk)
}

func FailWithMsg(msg string) {
	Result(NotOk, msg, map[string]interface{}{})
}

func FailWithCode(code int) {
	// 查找给定的错误码存在对应的错误信息, 默认使用NotOk
	msg := CustomError[NotOk]
	if val, ok := CustomError[code]; ok {
		msg = val
	}
	Result(code, msg, map[string]interface{}{})
}

// 写入json返回值
func JSON(c *gin.Context, code int, resp interface{}) {
	// 调用gin写入json
	c.JSON(code, resp)
	// 保存响应对象到context, Operation Log会读取到
	c.Set(global.Conf.System.OperationLogKey, resp)
}
