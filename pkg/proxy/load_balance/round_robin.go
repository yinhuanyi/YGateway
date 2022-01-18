/**
 * @Author：Robby
 * @Date：2022/1/18 18:34
 * @Function：
 **/

package load_balance

import (
	"errors"
)

type RoundRobinBalance struct {
	// 当前下游服务器的index值
	curIndex int
	// 下游服务器列表
	rss      []string
}

func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}


func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	lens := len(r.rss) //5
	if r.curIndex >= lens {
		r.curIndex = 0
	}
	// 基于当前的curIndex，获取下游服务器地址
	curAddr := r.rss[r.curIndex]
	// 然后更新当前的curIndex的值，取模是保证curIndex索引不越界
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

func (r *RoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}