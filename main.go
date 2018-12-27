package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strconv"
)

const RobotToken = ""

type Robot struct {
	bot *tgbotapi.BotAPI
}

// 初始化机器人
func (r *Robot) Init() {
	bot, err := tgbotapi.NewBotAPI(RobotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	r.bot = bot
}

func marshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func readData() []byte {
	fp, err := os.OpenFile("data.json", os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		return []byte("{}")
		//log.Fatal(err)
	}
	data := make([]byte, 100)
	n, err := fp.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data[:n]))
	return data[:n]
}

func writeData(data []byte) {
	//fp, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0755)
	fp, err := os.OpenFile("data.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

// 获取订阅用户列表
func getUsers() map[string]bool {
	data := readData()
	var jsonV map[string]bool
	json.Unmarshal(data, &jsonV)
	return jsonV
}

// 现在用户
func addUser(uid string) {
	data := readData()
	var jsonV map[string]bool
	json.Unmarshal(data, &jsonV)
	jsonV[uid] = true
	bt := marshal(jsonV)
	writeData(bt)
}

// 删除用户
func delUser(uid string) {
	data := readData()
	var jsonV map[string]bool
	json.Unmarshal(data, &jsonV)
	delete(jsonV, uid)
	bt := marshal(jsonV)
	writeData(bt)
}

// 首页
func Index(c *gin.Context) {
	c.JSON(0, gin.H{
		"version": "1.0",
		"msg":     "telegram robot",
	})
}

// 监听telegram消息
func (r Robot) Pulling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := r.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		chatId := update.Message.Chat.ID
		msgText := update.Message.Text

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		println("收到消息" + strconv.FormatInt(chatId, 10) + msgText)

		if msgText == "天王盖地虎" {
			addUser(strconv.FormatInt(chatId, 10))
			println("订阅成功")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "订阅成功")
			//msg.ReplyToMessageID = update.Message.MessageID
			r.bot.Send(msg)
		} else if msgText == "TD" || msgText == "td" {
			delUser(strconv.FormatInt(chatId, 10))
			println("退订成功")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "退订成功")
			//msg.ReplyToMessageID = update.Message.MessageID
			r.bot.Send(msg)
		} else {
			println("非法用户" + strconv.FormatInt(chatId, 10) + msgText)
		}

	}
}

// 通知接口
// 发消息给所有订阅用户
func (r Robot) Notify(c *gin.Context) {

	//msg := c.PostForm("msg")
	msg, _ := c.GetQuery("msg")

	println("notify msg = " + msg)

	users := getUsers()

	for uid := range users {
		//println(uid)
		chatId, _ := strconv.ParseInt(uid, 10, 64)
		println("chatid = " + strconv.FormatInt(chatId, 10))
		msg := tgbotapi.NewMessage(chatId, msg)
		r.bot.Send(msg)
	}

	c.JSON(0, gin.H{
		"msg": "ok",
	})
}

func main() {

	r := Robot{}
	r.Init()
	go r.Pulling()

	router := gin.Default()

	router.GET("/", Index)
	router.Any("/notify", r.Notify)

	fmt.Println("启动telegram robot")

	router.Run(":9008") // listen and serve on 0.0.0.0:8080

}

//
