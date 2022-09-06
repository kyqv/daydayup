package lc

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 198
// 你是一个专业的小偷，计划偷窃沿街的房屋。每间房内都藏有一定的现金，
// 影响你偷窃的唯一制约因素就是相邻的房屋装有相互连通的防盗系统，如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警。
// 给定一个代表每个房屋存放金额的非负整数数组，计算你 不触动警报装置的情况下 ，一夜之内能够偷窃到的最高金额。
// 示例 1：
// 输入：[1,2,3,1]
// 输出：4
// 解释：偷窃 1 号房屋 (金额 = 1) ，然后偷窃 3 号房屋 (金额 = 3)。
//      偷窃到的最高金额 = 1 + 3 = 4 。
// 示例 2：

// 输入：[2,7,9,3,1]
// 输出：12
// 解释：偷窃 1 号房屋 (金额 = 2), 偷窃 3 号房屋 (金额 = 9)，接着偷窃 5 号房屋 (金额 = 1)。
//      偷窃到的最高金额 = 2 + 9 + 1 = 12 。
func rob(money []int) int {
	if len(money) == 0 {
		return 0
	}

	if len(money) == 1 {
		return money[0]
	}

	//dp[i] 偷前i个最大金额
	dp := make([]int, len(money))
	dp[0] = money[0]
	dp[1] = max(money[0], money[1])

	for i := 2; i < len(money); i++ {
		dp[i] = max(dp[i-2]+money[i], dp[i-1])
	}

	return dp[len(dp)-1]
}

// 这个地方所有的房屋都 围成一圈
func rob2(money []int) int {
	if len(money) == 0 {
		return 0
	}

	if len(money) == 1 {
		return money[0]
	}

	return max(rob(money[0:len(money)-1]), rob(money[1:]))
}

func TestRob198(t *testing.T) {
	a := []int{1, 2, 3, 1}
	fmt.Println(rob(a))
	a = []int{2, 7, 9, 3, 1}
	fmt.Println(rob(a))

	a = []int{1, 2, 3, 1}
	fmt.Println(rob2(a))
	a = []int{2, 3, 2}
	fmt.Println(rob2(a))
}

// 单调栈
// 496. 下一个更大元素 I
// 给你两个 没有重复元素 的数组 nums1 和 nums2 ，其中nums1 是 nums2 的子集。
// 请你找出 nums1 中每个元素在 nums2 中的下一个比其大的值。
// nums1 中数字 x 的下一个更大元素是指 x 在 nums2 中对应位置的右边的第一个比 x 大的元素。如果不存在，对应位置输出 -1 。
// 输入: nums1 = [4,1,2], nums2 = [1,3,4,2].
// 输出: [-1,3,-1]
// 解释:
//     对于 num1 中的数字 4 ，你无法在第二个数组中找到下一个更大的数字，因此输出 -1 。
//     对于 num1 中的数字 1 ，第二个数组中数字1右边的下一个较大数字是 3 。
//     对于 num1 中的数字 2 ，第二个数组中没有下一个更大的数字，因此输出 -1
func nextGreaterElement(nums1 []int, nums2 []int) []int {
	//记录每个元素下个最大的
	result := map[int]int{}
	stack := make([]int, len(nums2))

	for i := len(nums2) - 1; i >= 0; i-- {
		num := nums2[i]

		//1.弹出所有比它小的
		for j := len(stack) - 1; j >= 0; j-- {
			if num < stack[j] {
				break
			}

			stack = stack[:len(stack)-1]
		}

		fmt.Printf("index:%v,stack:%v\n", i, stack)

		//当前栈中元素就是下一个最大的数
		if len(stack) == 0 {
			result[num] = -1
		} else {
			result[num] = stack[len(stack)-1]
		}

		//这个元素追加到栈
		stack = append(stack, num)
		fmt.Printf("index:%v,stack:%v\n", i, stack)
	}

	stack = stack[:0]
	for i := range nums1 {
		stack = append(stack, result[nums1[i]])
	}

	return stack
}

func TestNextGreaterElement496(t *testing.T) {
	a := []int{4, 1, 2}
	b := []int{1, 3, 4, 2}
	fmt.Println(nextGreaterElement(a, b))
}

func TestChannel(t *testing.T) {
	out := make(chan int, 100)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			out <- i
		}
		close(out)
		fmt.Println("close channel")
	}()
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("start")
		defer wg.Done()
		for i := range out {
			fmt.Println(i)
		}

		fmt.Println("fuck")
	}()
	wg.Wait()
}

func minimumDifference(nums []int) int64 {
	m := len(nums)
	n := m / 3
	minPQ := minHeap{nums[m-n:]}
	heap.Init(&minPQ)
	sum := 0
	for _, v := range nums[m-n:] {
		sum += v
	}
	sufMax := make([]int, m-n+1) // 后缀最大和
	sufMax[m-n] = sum
	for i := m - n - 1; i >= n; i-- {
		if v := nums[i]; v > minPQ.IntSlice[0] {
			sum += v - minPQ.IntSlice[0]
			minPQ.IntSlice[0] = v
			heap.Fix(&minPQ, 0)
		}
		sufMax[i] = sum
	}

	maxPQ := maxHeap{nums[:n]}
	heap.Init(&maxPQ)
	preMin := 0 // 前缀最小和
	for _, v := range nums[:n] {
		preMin += v
	}
	ans := preMin - sufMax[n]
	for i := n; i < m-n; i++ {
		if v := nums[i]; v < maxPQ.IntSlice[0] {
			preMin += v - maxPQ.IntSlice[0]
			maxPQ.IntSlice[0] = v
			heap.Fix(&maxPQ, 0)
		}
		ans = min(ans, preMin-sufMax[i+1])
	}
	return int64(ans)
}

type minHeap struct{ sort.IntSlice }

func (minHeap) Push(interface{})     {}
func (minHeap) Pop() (_ interface{}) { return }

type maxHeap struct{ sort.IntSlice }

func (h maxHeap) Less(i, j int) bool { return h.IntSlice[i] > h.IntSlice[j] }
func (maxHeap) Push(interface{})     {}
func (maxHeap) Pop() (_ interface{}) { return }

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func TestString(t *testing.T) {
	s := "我是中国人"
	fmt.Println(s)
	fmt.Println([]rune(s))
	fmt.Println(string([]rune(s)))

	fmt.Println(".....:", strings.Index(s, "是"))
	fmt.Println(".....:", strings.IndexRune(s, '是'))
	fmt.Println(string([]rune(s)[3]))
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now())
	now := time.Now()
	fmt.Println(now.Weekday().String())

	fmt.Println(now.Format("2006/01/02 15:04 PM"))
}

func TestFile(t *testing.T) {
	f, err := os.OpenFile("./1.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	f.WriteString("test a string!")

	buf := make([]byte, 1024)
	var str string
	f.Seek(0, io.SeekStart)
	for {
		n, ferr := f.Read(buf)
		if ferr != nil && ferr != io.EOF {
			fmt.Println(ferr.Error())
			break
		}
		if n == 0 {
			break
		}
		fmt.Println(n)
		str += string(buf[0:n])
	}
	fmt.Println(str)

	content, err := ioutil.ReadAll(f)
	fmt.Println(string(content))
}

func TestDomainI(t *testing.T) {
	cfg := map[string]map[string]int64{
		"play-ws.vod.shopee.com": {
			"157.185.156.117": 60,
			"138.113.218.16":  30,
			"210.61.180.5":    10,
		},
	}
	str, _ := json.Marshal(&cfg)

	println(string(str))
}
