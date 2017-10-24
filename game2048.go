package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"math/rand"
	"fmt"
)

//内存保存所有用户游戏状态
var game = make(map[string]*G2048)

//反射映射函数，少写些条件判断
var AC = map[string]interface{}{
	"left":  left,
	"right": right,
	"up":    up,
	"down":  down,
}

const suf = "game2048"
const auth = "Authorization"

type Start struct {
	Username string `用户名`
	Size     int    `游戏规模`
}

type Action struct {
	Dir string `left right up down`
}

type G2048 struct {
	Gmap  [][]int `矩阵`
	Score int     `当前得分`
	Size  int     `矩阵规模`
}

//解析开始的参数
func PraseStart(req *http.Request) (d Start, err error) {
	var data Start
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}

//解析用户动作
func PraseAction(req *http.Request) (a Action, err error) {
	var data Action
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}

//初始化游戏
func gameInit(n int) *G2048 {
	var m = [][]int{}
	for i := 0; i < n; i++ {
		mt := make([]int, 0, n)
		for j := 0; j < n; j++ {
			mt = append(mt, 0)
		}
		m = append(m, mt)
	}
	var g2048 = G2048{
		Gmap:  m,
		Score: 0,
		Size:  n,
	}
	Rand2(&g2048, 2, n)
	return &g2048
}

func getToken(name string) string {
	return name + suf
}

func end(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get(auth)
	delete(game, token)
	AllowOrigin(res)
}

func AllowOrigin(res http.ResponseWriter) {
	res.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	res.Header().Set("content-type", "application/json")             //返回数据格式是json
}

func start(res http.ResponseWriter, req *http.Request) {
	data, err := PraseStart(req)
	if err != nil || data.Size <= 2 || data.Size >= 10 {
		res.WriteHeader(400)
		return
	}
	token := getToken(data.Username)
	if _, ok := game[token]; !ok {
		game[token] = gameInit(data.Size)
	}
	jso, err := json.Marshal(game[token])
	AllowOrigin(res)
	res.Header().Set(auth, token)
	res.Write(jso)
}

//获取当前方格状态
func status(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get(auth)
	AllowOrigin(res)
	jso, err := json.Marshal(game[token])
	if err != nil {
		res.WriteHeader(404)
	} else {
		res.Write(jso)
	}
}

//玩游戏
func play(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get(auth)
	g := game[token]
	ac, err := PraseAction(req)
	if err != nil {
		res.WriteHeader(400)
		return
	}
	meth, ok := AC[ac.Dir]
	if !ok {
		res.WriteHeader(400)
		return
	}
	meth.(func(g2048 *G2048))(g)
	status(res, req)
}

//向左移动
func left(g *G2048) {
	c := 0
	n := g.Size
	for i := 0; i < n; i++ {
		for j := 1; j < n; j++ {
			t := j
			if g.Gmap[i][j] <= 0 {
				continue
			}
			for t >= 1 && g.Gmap[i][t-1] <= 0 {
				t--
			}
			g.Gmap[i][j], g.Gmap[i][t] = g.Gmap[i][t], g.Gmap[i][j]
		}

		for j := 1; j < n; j++ {
			if g.Gmap[i][j] == g.Gmap[i][j-1] && g.Gmap[i][j] != 0 {
				g.Gmap[i][j-1] += g.Gmap[i][j]
				c += g.Gmap[i][j-1]
				g.Gmap[i][j] = 0
			}
		}
		for j := 1; j < n; j++ {
			t := j
			if g.Gmap[i][j] <= 0 {
				continue
			}
			for t >= 1 && g.Gmap[i][t-1] <= 0 {
				t--
			}
			g.Gmap[i][j], g.Gmap[i][t] = g.Gmap[i][t], g.Gmap[i][j]
		}
	}
	g.Score += c
	Rand2(g, 1, n)
}

//向右移动
func right(g *G2048) {
	c := 0
	n := g.Size
	for i := 0; i < n; i++ {
		for j := n - 2; j >= 0; j-- {
			t := j
			if g.Gmap[i][j] <= 0 {
				continue
			}
			for t <= n-2 && g.Gmap[i][t+1] <= 0 {
				t++
			}
			g.Gmap[i][j], g.Gmap[i][t] = g.Gmap[i][t], g.Gmap[i][j]
		}

		for j := n - 2; j >= 0; j-- {
			if g.Gmap[i][j] == g.Gmap[i][j+1] && g.Gmap[i][j] != 0 {
				g.Gmap[i][j+1] += g.Gmap[i][j]
				c += g.Gmap[i][j+1]
				g.Gmap[i][j] = 0
			}
		}
		for j := n - 2; j >= 0; j-- {
			t := j
			if g.Gmap[i][j] <= 0 {
				continue
			}
			for t <= n-2 && g.Gmap[i][t+1] <= 0 {
				t++
			}
			g.Gmap[i][j], g.Gmap[i][t] = g.Gmap[i][t], g.Gmap[i][j]
		}
	}
	g.Score += c
	Rand2(g, 1, n)
}

//向上移动
func up(g *G2048) {
	c := 0
	n := g.Size
	for j := 0; j < n; j++ {
		for i := 1; i < n; i++ {
			t := i
			if g.Gmap[i][j] <= 0 {
				continue
			}
			for t >= 1 && g.Gmap[t-1][j] <= 0 {
				t--
			}
			g.Gmap[i][j], g.Gmap[t][j] = g.Gmap[t][j], g.Gmap[i][j]
		}

		for i := 1; i < n; i++ {
			if g.Gmap[i][j] == g.Gmap[i-1][j] && g.Gmap[i][j] != 0 {
				g.Gmap[i-1][j] += g.Gmap[i][j]
				c += g.Gmap[i-1][j]
				g.Gmap[i][j] = 0
			}
		}
		for i := 1; i < n; i++ {
			t := i
			if g.Gmap[i][j] <= 0 {
				continue
			}
			for t >= 1 && g.Gmap[t-1][j] <= 0 {
				t--
			}
			g.Gmap[i][j], g.Gmap[t][j] = g.Gmap[t][j], g.Gmap[i][j]
		}
	}
	g.Score += c
	Rand2(g, 1, n)
}

//向下移动
func down(g *G2048) {
	c := 0
	n := g.Size
	for j := 0; j < n; j++ {
		for i := n - 2; i >= 0; i-- {
			t := i
			if g.Gmap[i][j] <= 0 {
				continue
			}
			for t <= n-2 && g.Gmap[t+1][j] <= 0 {
				t++
			}
			g.Gmap[i][j], g.Gmap[t][j] = g.Gmap[t][j], g.Gmap[i][j]
		}

		for i := n - 2; i >= 0; i-- {
			if g.Gmap[i][j] == g.Gmap[i+1][j] && g.Gmap[i][j] != 0 {
				g.Gmap[i+1][j] += g.Gmap[i][j]
				c += g.Gmap[i+1][j]
				g.Gmap[i][j] = 0
			}
		}
		for i := n - 2; i >= 0; i-- {
			t := i
			if g.Gmap[i][j] <= 0 {
				continue
			}
			for t <= n-2 && g.Gmap[t+1][j] <= 0 {
				t++
			}
			g.Gmap[i][j], g.Gmap[t][j] = g.Gmap[t][j], g.Gmap[i][j]
		}
	}
	g.Score += c
	Rand2(g, 1, n)
}

//每次移动后随机生成 2
func Rand2(g *G2048, k, n int) {
	for k > 0 {
		i := rand.Intn(n) % n
		j := rand.Intn(n) % n
		if g.Gmap[i][j] == 0 {
			g.Gmap[i][j] = 2
			k--
		}
	}
}

//调度器
func Game(res http.ResponseWriter, req *http.Request) {
	meth := req.Method
	if meth == "POST" {
		start(res, req)
	} else if meth == "DELETE" {
		end(res, req)
	} else if meth == "GET" {
		status(res, req)
	} else if meth == "PUT" {
		play(res, req)
	}
}

func main() {
	fmt.Println("游戏服务启动中...")
	gameServer := http.NewServeMux()
	gameServer.HandleFunc("/game2048", Game)
	log.Fatal(http.ListenAndServe(":8080", gameServer))
}
