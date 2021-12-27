package util


import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	SeleniumPath string `yaml:"seleniumPath"`
	ChromeDriverPath string `yaml:"chromeDriverPath"`
	Port int `yaml:port`
}

//获取配置文件
func GetConf() (*Conf,error) {
	yamlFile, err := ioutil.ReadFile("file/config.yaml")
	if err != nil {
		return nil,err
	}
	//var c *Conf
	//err = yaml.Unmarshal(yamlFile, &c)
	c:=&Conf{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil,err
	}
	return c,nil
}
