package gocaptcha

import (
	"crypto/rand"
	"github.com/golang/freetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"math/big"
	mathrand "math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	ModeStr int = 0
	ModeArithmetic int = 1
)

const (
	defaultcharstr  = "23456789ABCDEFGHJKMNPQRSTWXYabcdefghjkmnpqrstwxy"
	operator        = "+-*/"
	defaultLen      = 4
	defaultFontSize = 25
	defaultDpi      = 72
	defaultMode = ModeStr
)

// 图形验证码结构体
type Captcha struct {
	W int // 图片宽度
	H int // 图片高度
	CodeLen int //验证码长度，默认长度4
	FontSize float64 //字体大小，默认25
	Dpi int   //清晰度，默认72
	FontPath string //字体目录
	FontName string //字体名字
	Mode int //验证模式 0：普通字符串，1：数学公式
	Charstr string
}

// 返回输出
func (me *Captcha) OutPut() (string, *image.RGBA) {
	img := me.initCanvas()
	return me.doImage(img)
}

// 初始化画布
func (me *Captcha) initCanvas() *image.RGBA {
	dest := image.NewRGBA(image.Rect(0, 0, me.W, me.H))

	r := uint8(me.RangeRand(157, 255))
	g := uint8(me.RangeRand(157, 255))
	b := uint8(me.RangeRand(157, 255))
	r = uint8(255)
	g = uint8(255)
	b = uint8(255)
	// 填充背景色
	for x := 0; x < me.W; x++ {
		for y := 0; y < me.H; y++ {
			dest.Set(x, y, color.RGBA{r, g, b, 255}) //设定alpha图片的透明度
		}
	}

	return dest
}

// 处理图像
func (me *Captcha) doImage(dest *image.RGBA) (string, *image.RGBA) {
	gc := draw2dimg.NewGraphicContext(dest)
	defer gc.Close()
	defer gc.FillStroke()

	me.setFont(gc)
	me.doPoint(gc)
	me.doLine(gc)
	me.doSinLine(gc)

	var codeStr string
	if me.Mode == ModeArithmetic {
		ret, arithData := me.getArithmeticData()
		codeStr = ret
		me.doArithmetic(gc, arithData)
	} else {
		codeStr = me.GetRandCode()
		me.doCode(gc, codeStr)
	}

	return codeStr, dest
}

// 设置字体
func (me *Captcha) setFont(gc *draw2dimg.GraphicContext) {
	if me.FontPath == "" {
		return
	}
	if me.FontName == "" {
		return
	}
	// 字体文件
	fontFile := strings.TrimRight(me.FontPath, "/") + "/" + strings.TrimLeft(me.FontName, "/")
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return
	}
	// 设置自定义字体相关信息
	gc.FontCache = draw2d.NewSyncFolderFontCache(me.FontPath)
	gc.FontCache.Store(draw2d.FontData{Name: me.FontName, Family: 0, Style: draw2d.FontStyleNormal}, font)
	gc.SetFontData(draw2d.FontData{Name: me.FontName, Style: draw2d.FontStyleNormal})

	//设置清晰度
	if me.Dpi <= 0 {
		me.Dpi = defaultDpi
	}
	gc.SetDPI(me.Dpi)

	// 设置字体大小
	if me.FontSize <= 0 {
		me.FontSize = defaultFontSize
	}
	gc.SetFontSize(me.FontSize)
}

// 增加干扰点
func (me *Captcha) doPoint(gc *draw2dimg.GraphicContext) {
	for n := 0; n < 50; n++ {
		gc.SetLineWidth(float64(me.RangeRand(1, 3)))

		// 随机色
		r := uint8(me.RangeRand(0, 255))
		g := uint8(me.RangeRand(0, 255))
		b := uint8(me.RangeRand(0, 255))

		gc.SetStrokeColor(color.RGBA{r, g, b, 255})

		x := me.RangeRand(0, int64(me.W)+10) + 1
		y := me.RangeRand(0, int64(me.H)+5) + 1

		gc.MoveTo(float64(x), float64(y))
		gc.LineTo(float64(x+me.RangeRand(1, 2)), float64(y+me.RangeRand(1, 2)))

		gc.Stroke()
	}
}

// 增加干扰线
func (me *Captcha) doLine(gc *draw2dimg.GraphicContext) {
	// 设置干扰线
	for n := 0; n < 5; n++ {
		// gc.SetLineWidth(float64(captcha.RangeRand(1, 2)))
		gc.SetLineWidth(1)

		// 随机背景色
		r := uint8(me.RangeRand(0, 255))
		g := uint8(me.RangeRand(0, 255))
		b := uint8(me.RangeRand(0, 255))

		gc.SetStrokeColor(color.RGBA{r, g, b, 255})

		// 初始化位置
		gc.MoveTo(float64(me.RangeRand(0, int64(me.W)+10)), float64(me.RangeRand(0, int64(me.H)+5)))
		gc.LineTo(float64(me.RangeRand(0, int64(me.W)+10)), float64(me.RangeRand(0, int64(me.H)+5)))

		gc.Stroke()
	}
}

// 增加正弦干扰线
func (me *Captcha) doSinLine(gc *draw2dimg.GraphicContext) {
	h1 := me.RangeRand(-12, 12)
	h2 := me.RangeRand(-1, 1)
	w2 := me.RangeRand(5, 20)
	h3 := me.RangeRand(5, 10)

	h := float64(me.H)
	w := float64(me.W)

	// 随机色
	r := uint8(me.RangeRand(128, 255))
	g := uint8(me.RangeRand(128, 255))
	b := uint8(me.RangeRand(128, 255))

	gc.SetStrokeColor(color.RGBA{r, g, b, 255})
	gc.SetLineWidth(float64(me.RangeRand(2, 4)))

	var i float64
	for i = -w / 2; i < w/2; i = i + 0.1 {
		y := h/float64(h3)*math.Sin(i/float64(w2)) + h/2 + float64(h1)

		gc.LineTo(i+w/2, y)

		if h2 == 0 {
			gc.LineTo(i+w/2, y+float64(h2))
		}
	}

	gc.Stroke()
}

// 验证码字符设置到图像上
func (me *Captcha) doCode(gc *draw2dimg.GraphicContext, code string) {
	var xBase = me.W / me.CodeLen
	tmpCode := []rune(code)
	for l,codeV := range tmpCode {
		randY := me.RangeRand(1, int64(me.H/2))
		randX := me.RangeRand(1, 5)
		x := xBase * l + int(randX)

		r := uint8(me.RangeRand(0, 255))
		g := uint8(me.RangeRand(0, 255))
		b := uint8(me.RangeRand(0, 255))

		gc.SetFillColor(color.RGBA{r, g, b, 255})

		gc.FillStringAt(string(codeV), float64(x), float64(me.H/2)+float64(randY))
		gc.Stroke()
	}
}

// 获取算术运算公式
func (me *Captcha) getArithmeticData() (string, []string) {
	num1 := int(me.RangeRand(11, 20))
	num2 := int(me.RangeRand(1, 10))
	opArr := []rune(operator)
	opRand := opArr[me.RangeRand(0, 3)]

	strNum1 := strconv.Itoa(num1)
	strNum2 := strconv.Itoa(num2)

	var ret int
	var opRet string
	switch string(opRand) {
	case "+":
		ret = num1 + num2
		opRet = "+"
	case "-":
		ret = num1 - num2
		opRet = "-"
	case "*":
		ret = num1 * num2
		opRet = "×"
	case "/":
		strNum1 = strconv.Itoa(num1 * num2)
		ret = num1
		opRet = "÷"
	}

	return strconv.Itoa(ret), []string{strNum1, opRet, strNum2, "=", "?"}
}


// 验证码字符设置到图像上
func (me *Captcha) doArithmetic(gc *draw2dimg.GraphicContext, arithArr []string) {
	var xBase = me.W / me.CodeLen
	for l := 0; l < len(arithArr); l++ {
		randY := me.RangeRand(0, 10)
		randX := me.RangeRand(5, 10)
		x := xBase * l + int(randX)

		r := uint8(me.RangeRand(10, 255))
		g := uint8(me.RangeRand(10, 255))
		b := uint8(me.RangeRand(10, 255))

		gc.SetFillColor(color.RGBA{r, g, b, 255})

		gc.FillStringAt(arithArr[l], float64(x), float64(me.H/2)+float64(randY))
		gc.Stroke()
	}
}

// 获取区间[min, max]的随机数
func (me *Captcha) RangeRand(min, max int64) int64 {
	if min > max {
		min,max = max,min
	}
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

// 获取随机字符串
func (me *Captcha) GetRandCode() string {
	if me.CodeLen <= 0 {
		me.CodeLen = defaultLen
	}
	letterStr := []rune(me.Charstr)
	lenStr := len(letterStr)
	b := make([]rune, me.CodeLen)
	for i := range b {
		mathrand.Seed(time.Now().UnixNano())
		b[i] = letterStr[mathrand.Intn(lenStr)]
	}
	return string(b)
}


// 设置模式
func (me *Captcha) SetMode(mode int) {
	me.Mode = mode
}

// 设置字体路径
func (me *Captcha) SetFontPath(fontPath string) {
	me.FontPath = fontPath
}

// 设置字体名称
func (me *Captcha) SetFontName(fontName string) {
	me.FontName = fontName
}

// 设置字体大小
func (me *Captcha) SetFontSize(fontSize float64) {
	me.FontSize = fontSize
}

// 设置字符串
func (me *Captcha) SetCharstr(charstr string) {
	me.Charstr = charstr
}

// 实例化验证码
func NewCaptcha(w, h int) *Captcha {
	captchaObj := &Captcha{W: w, H: h}
	captchaObj.CodeLen = defaultLen
	captchaObj.FontSize = defaultFontSize
	captchaObj.Dpi = defaultDpi
	captchaObj.Mode = defaultMode
	captchaObj.Charstr = defaultcharstr
	return captchaObj
}

