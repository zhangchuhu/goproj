package main

import (
	"fmt"
	"log"
	"net/http"
	//"text/template"
	"html/template"
	_ "net/http/pprof"
)

type String string

type Struct struct {
	Greeting string
	Punct    string
	Who      string
}

func (h String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Form)
	fmt.Println(r.URL.Path)
	fmt.Println(r.Form["url_long"])
	fmt.Fprint(w, h)
}

func (h *Struct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, h.Greeting, h.Punct, h.Who)
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		if t.Execute(w, nil) != nil {
			log.Println("execute faild")
		}
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])

		// 转义标签
		//template.HTMLEscape(w, []byte(r.Form.Get("username")))

		// 不转义标签
		t, _ := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		//t.ExecuteTemplate(w, "T", "<script>alert('yes');</script>")   // text/template
		t.ExecuteTemplate(w, "T", template.HTML("<script>alert('yes');</script>"))
	}
}

func main() {
	http.Handle("/str", String("I'm a frayed knot."))
	http.Handle("/st", &Struct{"Hello", ":", "Gophers!"})
	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe("0.0.0.0:4000", nil))
}
