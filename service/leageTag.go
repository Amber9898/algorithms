package service

import (
	"algorithm/utils"
	"fmt"
)

// SortLeague 签位排序
//固定轮转编排
//固定轮转法也叫常规轮转法，是我国传统的编排方法。
//左边第一号固定不动，逆时针转动，逐一排出。
type round struct {
	p1 int
	p2 int
}

func SortLeague(n int) (res [][]*round) {
	total := n //比赛轮次
	if n%2 == 0 {
		total--
	}
	var queue []int//队列
	//0号不放进去
	for i := 1; i < n; i++ {
		queue = append(queue, i)
	}

	//长度奇数，补-1代表轮空
	if len(queue)+1%2 != 0 {
		queue = append(queue, -1)
	}
	fmt.Println(queue)

	for r := 0; r < total; r++ {
		noLast := false
		var tmp []*round
		l, r := 0, len(queue)-1 //左右指针
		if queue[r] != -1{
			tmp = append(tmp, &round{p1: 0, p2: queue[r]})
		}
		r--
		for l < r {
			if queue[l] == -1 || queue[r] == -1{
				if queue[l] == n-1  || queue[r] == n-1{
					noLast = true
				}
				l++
				r--
				continue
			}
			tmp = append(tmp, &round{p1: queue[l], p2: queue[r]})
			l++
			r--
		}
		//循环队列
		queue =  append([]int{queue[len(queue)-1]}, queue[0:len(queue)-1]...)
		if !noLast{
			res = append(res, tmp)
		}
	}

	var strList []string
	for _, v := range res{
		str := "["
		for _, vv := range v{
			tmpStr := fmt.Sprintf("[%d, %d],", vv.p1, vv.p2)
			str += tmpStr
		}
		str = str[:len(str)-1]+"]"
		strList = append(strList, str)
	}

	utils.WriteSortExcel(strList)

	return
}


