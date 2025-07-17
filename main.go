package main

import "ReadBook/router"

// Book 定义小说数据结构

func main() {
	r := router.InitGin()
	r.Run(":8055")
}
