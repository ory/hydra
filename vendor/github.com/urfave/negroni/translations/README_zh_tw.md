# Negroni(尼格龍尼) [![GoDoc](https://godoc.org/github.com/codegangsta/negroni?status.svg)](http://godoc.org/github.com/codegangsta/negroni) [![wercker status](https://app.wercker.com/status/13688a4a94b82d84a0b8d038c4965b61/s "wercker status")](https://app.wercker.com/project/bykey/13688a4a94b82d84a0b8d038c4965b61)

尼格龍尼符合Go的web 中介器特性. 精簡、非侵入式、鼓勵使用 `net/http`  Handler.

如果你喜歡[Martini](http://github.com/go-martini/martini)，但覺得這其中包太多神奇的功能，那麼尼格龍尼會是你的最佳選擇。

## 入門

安裝完Go且設好[GOPATH](http://golang.org/doc/code.html#GOPATH)，建立你的第一個`.go`檔。可以命名為`server.go`。

~~~ go
package main

import (
  "github.com/codegangsta/negroni"
  "net/http"
  "fmt"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  })

  n := negroni.Classic()
  n.UseHandler(mux)
  n.Run(":3000")
}
~~~

安裝尼格龍尼套件 (最低需求為**go 1.1**或更高版本):
~~~
go get github.com/codegangsta/negroni
~~~

執行伺服器:
~~~
go run server.go
~~~

你現在起了一個Go的net/http網頁伺服器在`localhost:3000`.

## 有問題?
如果你有問題或新功能建議，[到這郵件群組討論](https://groups.google.com/forum/#!forum/negroni-users)。尼格龍尼在GitHub上的issues專欄是專門用來回報bug跟pull requests。

## 尼格龍尼是個framework嗎?
尼格龍尼**不是**framework，是個設計用來直接使用net/http的library。

## 路由?
尼格龍尼是BYOR (Bring your own Router，帶給你自訂路由)。在Go社群已經有大量可用的http路由器, 尼格龍尼試著做好完全支援`net/http`，例如與[Gorilla Mux](http://github.com/gorilla/mux)整合:

~~~ go
router := mux.NewRouter()
router.HandleFunc("/", HomeHandler)

n := negroni.New(中介器1, 中介器2)
// Or use a 中介器 with the Use() function
n.Use(中介器3)
// router goes last
n.UseHandler(router)

n.Run(":3000")
~~~

## `negroni.Classic()`
`negroni.Classic()` 提供一些好用的預設中介器:

* `negroni.Recovery` - Panic 還原中介器
* `negroni.Logging` - Request/Response 紀錄中介器
* `negroni.Static` - 在"public"目錄下的靜態檔案服務

尼格龍尼的這些功能讓你開發變得很簡單。

## 處理器(Handlers)
尼格龍尼提供一個雙向中介器的機制，介面為`negroni.Handler`:

~~~ go
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
~~~

如果中介器沒有寫入ResponseWriter，會呼叫通道裡面的下個`http.HandlerFunc`讓給中介處理器。可以被用來做良好的應用:

~~~ go
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // 在這之前做一些事
  next(rw, r)
  // 在這之後做一些事
}
~~~

然後你可以透過`Use`函數對應到處理器的通道:

~~~ go
n := negroni.New()
n.Use(negroni.HandlerFunc(MyMiddleware))
~~~

你也可以應原始的舊`http.Handler`:

~~~ go
n := negroni.New()

mux := http.NewServeMux()
// map your routes

n.UseHandler(mux)

n.Run(":3000")
~~~

## `Run()`
尼格龍尼有一個很好用的函數`Run`，`Run`接收addr字串辨識[http.ListenAndServe](http://golang.org/pkg/net/http#ListenAndServe)。

~~~ go
n := negroni.Classic()
// ...
log.Fatal(http.ListenAndServe(":8080", n))
~~~

## 路由特有中介器
如果你有一群路由需要執行特別的中介器，你可以簡單的建立一個新的尼格龍尼實體當作路由處理器。

~~~ go
router := mux.NewRouter()
adminRoutes := mux.NewRouter()
// add admin routes here

// 為管理中介器建立一個新的尼格龍尼
router.Handle("/admin", negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(adminRoutes),
))
~~~

## 第三方中介器

以下為目前尼格龍尼兼容的中介器清單。如果你自己做了一個中介器請自由放入你的中介器互換連結:

| 中介器 | 作者 | 說明 |
| -----------|--------|-------------|
| [RestGate](https://github.com/pjebs/restgate) | [Prasanga Siripala](https://github.com/pjebs) | REST API入口的安全認證 |
| [Graceful](https://github.com/stretchr/graceful) | [Tyler Bunnell](https://github.com/tylerb) | 優雅的HTTP關機 |
| [secure](https://github.com/unrolled/secure) | [Cory Jacobsen](https://github.com/unrolled) | 檢疫安全功能的中介器 |
| [JWT Middleware](https://github.com/auth0/go-jwt-middleware) | [Auth0](https://github.com/auth0) | 檢查JWT的中介器用來解析傳入請求的`Authorization` header |
| [binding](https://github.com/mholt/binding) | [Matt Holt](https://github.com/mholt) | 將HTTP請求轉到structs的資料綁定 |
| [logrus](https://github.com/meatballhat/negroni-logrus) | [Dan Buch](https://github.com/meatballhat) | 基於Logrus的紀錄器 |
| [render](https://github.com/unrolled/render) | [Cory Jacobsen](https://github.com/unrolled) | 渲染JSON、XML、HTML的樣板 |
| [gorelic](https://github.com/jingweno/negroni-gorelic) | [Jingwen Owen Ou](https://github.com/jingweno) | Go執行中的New Relic agent |
| [gzip](https://github.com/phyber/negroni-gzip) | [phyber](https://github.com/phyber) | GZIP回應壓縮 |
| [oauth2](https://github.com/goincremental/negroni-oauth2) | [David Bochenski](https://github.com/bochenski) | oAuth2中介器 |
| [sessions](https://github.com/goincremental/negroni-sessions) | [David Bochenski](https://github.com/bochenski) | Session管理 |
| [permissions2](https://github.com/xyproto/permissions2) | [Alexander Rødseth](https://github.com/xyproto) | Cookies與使用者權限 |
| [onthefly](https://github.com/xyproto/onthefly) | [Alexander Rødseth](https://github.com/xyproto) | 快速產生TinySVG、HTM、CSS |
| [cors](https://github.com/rs/cors) | [Olivier Poitrey](https://github.com/rs) | [Cross Origin Resource Sharing](http://www.w3.org/TR/cors/) 支援(CORS) |
| [xrequestid](https://github.com/pilu/xrequestid) | [Andrea Franz](https://github.com/pilu) | 在每個request指定一個隨機X-Request-Id header的中介器 |
| [VanGoH](https://github.com/auroratechnologies/vangoh) | [Taylor Wrobel](https://github.com/twrobel3) | Configurable [AWS-Style](http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html) HMAC 授權中介器 |
| [stats](https://github.com/thoas/stats) | [Florent Messa](https://github.com/thoas) | 儲存關於你的網頁應用資訊 (回應時間之類) |

## 範例
[mooseware](https://github.com/xyproto/mooseware)是用來寫尼格龍尼中介處理器的骨架，由[Alexander Rødseth](https://github.com/xyproto)建立。

## 即時程式重載?
[gin](https://github.com/codegangsta/gin)和[fresh](https://github.com/pilu/fresh)兩個尼格龍尼即時重載的應用。

## Go & 尼格龍尼初學者必讀

* [使用Context將資訊從中介器送到處理器](http://elithrar.github.io/article/map-string-interface/)
* [理解中介器](http://mattstauffer.co/blog/laravel-5.0-middleware-replacing-filters)

## 關於

尼格龍尼正是[Code Gangsta](http://codegangsta.io/)的執著設計。
譯者: Festum Qin (Festum@G.PL)
