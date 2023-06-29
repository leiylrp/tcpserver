package timer

import (
	"fmt"
	"reflect"
)

/*
   定义一个延迟调用函数
	延迟调用函数就是 时间定时器超时的时候，触发的事先注册好的
	回调函数
*/

type DelayFunc struct {
	// f 延迟调用函数原型
	f func(...interface{})

	// args 延迟调用函数传参
	args []interface{}
}

// NewDelayFunc 初始化延迟函数
func NewDelayFunc(f func(...interface{}), args []interface{}) *DelayFunc {
	fmt.Println("NewDelayFunc, args: ",  args)
	return &DelayFunc{
		f: f,
		args: args,
	}
}

// String 格式化记录信息到日志中
func (df *DelayFunc) String() string {
	return fmt.Sprintf("{DelayFunc:%s, args:%v}", reflect.TypeOf(df.f), df.args)
}

// Call 调用定时函数 如果执行失败，则抛出异常
func (df *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {

		}
	}()

	// 调用定时器超时函数

	df.f(df.args...)
}
