package util

import (
	"compress/gzip"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpItem struct {
	Method string
	Url string
	Header map[string]string
	PostData string
	Proxy string
}

type HttpResult struct {
	Code int
	Html string
	Document *goquery.Document
}

func (this *HttpItem) GetHtml() (*HttpResult,error) {
	client := &http.Client{
		//CheckRedirect:nil,
	}
	//超时
	client.Timeout=time.Second*20
	//代理
	if len(this.Proxy)>0 {
		proxy, err := url.Parse(this.Proxy)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}
	httpreq, err := http.NewRequest(this.Method, this.Url, strings.NewReader(this.PostData))
	if err!=nil {
		return nil,err
	}
	if this.Header!=nil {
		for k,v:=range this.Header {
			httpreq.Header.Set(k,v)
		}
	}
	resp, err:= client.Do(httpreq)
	if err!=nil {
		return nil,err
	}
	// get converter to utf-8
	//b,err:=ioutil.ReadAll(resp.Body)
	var bodyStr string
	if resp.Header.Get("Content-Encoding") == "gzip" {
		bodyStr = this.changeCharsetEncodingAutoGzipSupport(resp.Header.Get("Content-Type"), resp.Body)
	} else {
		bodyStr = this.changeCharsetEncodingAuto(resp.Header.Get("Content-Type"), resp.Body)
	}
	result:=&HttpResult{Code:resp.StatusCode, Html:bodyStr}
	defer resp.Body.Close()
	return result,err
}

func (this *HttpItem) GetDocument()(*HttpResult,error){
	result,err:= this.GetHtml()
	if err!=nil {
		return nil,err
	}
	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(strings.NewReader(result.Html)); err != nil {
		return nil,err
	}
	result.Document=doc
	return result,err
}

func (this *HttpItem) changeCharsetEncodingAutoGzipSupport(contentTypeStr string, sor io.ReadCloser) string {
	var err error
	gzipReader, err := gzip.NewReader(sor)
	if err != nil {
		return ""
	}
	defer gzipReader.Close()
	destReader, err := charset.NewReader(gzipReader, contentTypeStr)

	if err != nil {
		destReader = sor
	}

	var sorbody []byte
	if sorbody, err = ioutil.ReadAll(destReader); err != nil {
	}
	//e,name,certain := charset.DetermineEncoding(sorbody,contentTypeStr)
	bodystr := string(sorbody)
	return bodystr
}

// Charset auto determine. Use golang.org/x/net/html/charset. Get page body and change it to utf-8
func (this *HttpItem) changeCharsetEncodingAuto(contentTypeStr string, sor io.ReadCloser) string {
	var err error
	destReader, err := charset.NewReader(sor, contentTypeStr)
	if err != nil {
		destReader = sor
	}
	var sorbody []byte
	if sorbody, err = ioutil.ReadAll(destReader); err != nil {
	}
	//e,name,certain := charset.DetermineEncoding(sorbody,contentTypeStr)
	bodystr := string(sorbody)
	return bodystr
}