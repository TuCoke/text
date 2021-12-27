package util

import (
	"bytes"
	"fmt"
	"github.com/beevik/etree"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"regexp"
)

func JadeSecret() map[string]string {
	// 初始化根节点
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("file/jade-secret.xml"); err != nil {
		g.Log().Panicf(fmt.Sprintf("xml解析失败：%s",err.Error()))
	}
	var secret=make(map[string]string)
	root := doc.SelectElement("ttFont")
	for _,t:=range root.SelectElement("cmap").SelectElement("cmap_format_4").SelectElements("map"){
		code:=t.SelectAttrValue("code","")
		name:=t.SelectAttrValue("name","")
		secret[code]=name
	}
	secret["0x2e"]="."
	return secret
}
//解析字体文件中数字
func GetWoffNum(jadeSecret map[string]string,data string) string {
	//data:="￥&#593;&#2539;&#594;&#9281;&#9283;&#46;&#2537;&#1407;"
	var value bytes.Buffer
	for _,v:=range regexp.MustCompile(`&#(.*?);`).FindAllStringSubmatch(data,-1){
		if w,ok:=jadeSecret[fmt.Sprintf("0x%x",gconv.Int(v[1]))];ok {
			value.WriteString(w)
		}
	}
	return value.String()
}