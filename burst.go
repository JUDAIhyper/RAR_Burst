package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	pwdPath string
	rarPath string

	password = make(chan string) //接收密码
	isOver   = make(chan bool)   //判断是否退出
)

func main() {
	flag.StringVar(&pwdPath, "p", "./password.txt", "密码本路径")
	flag.StringVar(&rarPath, "f", "./1.rar", "文件路径")
	flag.Parse()

	go passTxt(pwdPath)

Loop:
	for {
		select {
		case rarpwd := <-password:
			go cmdShell(rarPath, rarpwd)
		case <-time.After(time.Second * time.Duration(10)):
			break Loop
		case <-isOver:
			break Loop
		}
	}
}

func cmdShell(rarPath string, pwd string) {
	cmd := exec.Command("unrar", "e", "-p"+pwd, rarPath, "./test")
	out, _ := cmd.Output()

	log.Println("pass: ", pwd)
	// fmt.Println("out: ", string(out))

	if len(out) >= 200 {
		fmt.Printf("密码为: %s\n", pwd)
		isOver <- true
	}
}

func passTxt(pwdPath string) {
	fp, _ := os.OpenFile(pwdPath, os.O_RDONLY, 6)
	defer fp.Close()

	// 创建文件的缓存区
	r := bufio.NewReader(fp)
	for {
		pass, err2 := r.ReadBytes('\n')
		if err2 == io.EOF { //文件末尾
			break
		}
		pass = pass[:len(pass)-2] // 去除末尾 /n
		password <- string(pass)
	}
}
