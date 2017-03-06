package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func PrintLocalDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	conn, err := dial.Dial(network, addr)
	if err != nil {
		return conn, err
	}

	fmt.Println("connect done, use", conn.LocalAddr().String())
	return conn, err
}

func doGet(url string, id int) {
	/*
		client := &http.Client{
			Transport: &http.Transport{
				Dial: PrintLocalDial,
			},
		}
	*/

	//resp, err := client.Get(url)
	resp, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader("name=cjb"))
	if err != nil {
		fmt.Println(err)
		return
	}

	//_, err = ioutil.ReadAll(resp.Body)
	buf, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%d: %s -- %v\n", id, string(buf), err)
	if err := resp.Body.Close(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	//Url := "http://127.0.0.1:4000/str"
	Url := "http://106.38.255.199:9098/buddy/full_sync/"
	for {
		doGet(Url, 1)
		//go doGet(client, Url, 2)
		time.Sleep(1 * time.Second)
	}
}
