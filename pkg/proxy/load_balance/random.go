/**
 * @Author：Robby
 * @Date：2022/1/18 18:30
 * @Function：
 **/

package load_balance

import (
	"errors"
	"math/rand"
)

type RandomBalance struct {
	// 当前下游服务器的index值
	curIndex int
	// 下游服务器列表
	rss      []string
}

// Add 添加服务器地址到负载均衡中
func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

// Next 从r.rss数组中，随机的获取下游服务器地址
func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

// Get 为了让RandomBalance实现接口的Get方法
func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}