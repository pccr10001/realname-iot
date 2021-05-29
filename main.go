package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/pccr10001/realname-iot/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var client mqtt.Client
var config model.Config

func main() {

	e := gin.Default()
	e.LoadHTMLGlob("templates/**")
	e.GET("/qr/:id", GetQr)

	if f, err := ioutil.ReadFile("config.yml"); err != nil && yaml.Unmarshal(f, &config) != nil {
		log.Fatalln("failed to parse config.yml, err = " + err.Error())
	}

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", config.Mqtt.Server, config.Mqtt.Port)).
		SetClientID(config.Mqtt.ClientId).
		SetUsername(config.Mqtt.Username).
		SetPassword(config.Mqtt.Password).
		SetAutoReconnect(true).
		SetConnectRetry(true)

	client = mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	_ = e.Run("0.0.0.0:8080")
}

func GetQr(c *gin.Context) {
	id := c.Param("id")
	if _, err := strconv.Atoi(id); err != nil && len(id) != 15 {
		c.String(http.StatusForbidden, "Invalid place Id")
		return
	}

	go func() {
		if token := client.Publish(config.Mqtt.TopicPrefix+id, 1, false, id); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}()

	c.Redirect(http.StatusFound, "intent:smsto:1922#Intent;action=android.intent.action.SENDTO;S.sms_body=%E5%A0%B4%E6%89%80%E4%BB%A3%E7%A2%BC%EF%BC%9A"+id+"%0A%E6%9C%AC%E6%AC%A1%E5%AF%A6%E8%81%AF%E7%B0%A1%E8%A8%8A%E9%99%90%E9%98%B2%E7%96%AB%E7%9B%AE%E7%9A%84%E4%BD%BF%E7%94%A8%E3%80%82;end")
}
