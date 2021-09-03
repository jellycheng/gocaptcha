# gocaptcha
```
图片验证码
```
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/jellycheng/gocaptcha.svg)](https://pkg.go.dev/github.com/jellycheng/gocaptcha)

## 使用示例1
```
package main

import (
	"fmt"
	"image/png"
	"net/http"
	"github.com/jellycheng/gocaptcha"
)

func captchaHandler(w http.ResponseWriter, r *http.Request) {
	cp := gocaptcha.NewCaptcha(120, 40)
	cp.CodeLen = 4 // 字符串长度
	cp.SetFontPath("./font/") // 设置字体目录
	cp.SetFontName("luxisr.ttf") //字体文件名
	code, img := cp.OutPut()
	//图片验证码内容，这里补充你的代码逻辑todo
	fmt.Println(code)

	// 响应图片
	w.Header().Set("Content-Type", "image/png; charset=utf-8")
	png.Encode(w, img)
}

func main() {
	http.HandleFunc("/captcha", captchaHandler)

	_ = http.ListenAndServe("localhost:1234", nil)

}


```
![](http://captcha.nfangbian.com/captcha)


## gin框架集成图片验证码示例
```
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jellycheng/gocaptcha"
	"image/png"
)
func main()  {
	r := gin.Default()
	//设置图片验证码地址
	r.GET("/captcha", func(c *gin.Context) {
		cp := gocaptcha.NewCaptcha(120, 40)
		cp.CodeLen = 4 // 字符串长度
		cp.SetFontPath("./font/") // 设置字体目录
		cp.SetFontName("luxisr.ttf") //字体文件名
		code, img := cp.OutPut()
		//图片验证码内容，这里补充你的代码逻辑todo
		fmt.Println(code)

		// 响应图片
		c.Header("Content-Type", "image/png; charset=utf-8")
		png.Encode(c.Writer, img)
	})

	r.Run(":1234")
}

```
