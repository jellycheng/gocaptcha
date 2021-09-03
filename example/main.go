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
	//cp.SetFontName("SimSun18030.ttc")
	//cp.SetFontName("elephant.ttf")
	//cp.SetMode(gocaptcha.ModeArithmetic) // 设置为数学公式
	//cp.SetMode(gocaptcha.ModeStr) // 设置为字符串
	//cp.FontSize = 25
	//cp.Charstr = "123456789上海北京"
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

