// Created by quicksandzn@gmail.com on 2018/7/25
package main

import (
	"log"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mail struct {
		From     string `yaml:"from"`
		To       string `yaml:"to"`
		Cc       string `yaml:"cc"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		UserName string `yaml:"username"`
		PassWord string `yaml:"password"`
	}
}

const (
	ConfigPath             = "./config/conf.yml"
	MailHeaderFieldFrom    = "From"
	MailHeaderFieldCc      = "Cc"
	MailHeaderFieldTo      = "To"
	MailHeaderFieldSubject = "Subject"
	MailBodyContentType    = "text/html"
)

func main() {
	config := Config{}
	buffer, err := ioutil.ReadFile(ConfigPath)
	failOnError(err, "read config error")
	err = yaml.Unmarshal(buffer, &config)
	failOnError(err, "yml convert error")

	m := gomail.NewMessage()
	m.SetHeader(MailHeaderFieldFrom, config.Mail.From)
	m.SetHeader(MailHeaderFieldTo, config.Mail.To)
	m.SetAddressHeader(MailHeaderFieldCc, config.Mail.Cc, config.Mail.Cc)
	m.SetHeader(MailHeaderFieldSubject, "Hello!")
	m.SetBody(MailBodyContentType, "Hello")
	//m.Attach("")

	d := gomail.NewDialer(config.Mail.Host, config.Mail.Port, config.Mail.UserName, config.Mail.PassWord)

	for   {
		err = d.DialAndSend(m)
	}

	failOnError(err, "send msg error")

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(err)
	}
}
