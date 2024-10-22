package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"github.com/eatmoreapple/openwechat"
)

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("./storage.json")

	defer reloadStorage.Close()

	// 登录
	if err := bot.HotLogin(reloadStorage); err != nil {
		fmt.Println("热登陆失败，尝试免扫码登录")
		bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())
	}

	fmt.Println("登陆成功")
	defaultPriority := "5"

	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsSendBySelf() { //自己发送的消息
			//跳过
			return
		} else if msg.IsSendByFriend() { //好友发送的消息
			friendSender, err := msg.Sender()
			if err != nil {
				fmt.Println(err)
				return
			}

			friendSenderName := friendSender.RemarkName
			if len(friendSender.RemarkName) == 0 {
				friendSenderName = friendSender.NickName
			}

			if msg.IsText() {
				fmt.Println(friendSenderName, ":", msg.Content)
				fmt.Println(push(friendSenderName, defaultPriority, msg.Content))
			} else if msg.IsPicture() {
				fmt.Println(friendSenderName, ":", "[图片]")
				fmt.Println(push(friendSenderName, defaultPriority, "[图片]"))
			} else if msg.IsVoice() {
				fmt.Println(friendSenderName, ":", "[语音]")
				fmt.Println(push(friendSenderName, defaultPriority, "[语音]"))
			} else if msg.IsVideo() {
				fmt.Println(friendSenderName, ":", "[视频]")
				fmt.Println(push(friendSenderName, defaultPriority, "[视频]"))
			} else if msg.IsEmoticon() {
				fmt.Println(friendSenderName, ":", "[动画表情]")
				fmt.Println(push(friendSenderName, defaultPriority, "[动画表情]"))
			} else {
				fmt.Println(friendSenderName, ":", msg.Content)
				fmt.Println(push(friendSenderName, defaultPriority, msg.Content))
			}
		} else { //群聊发送的消息
			groupSender, err := msg.SenderInGroup()
			if err != nil {
				fmt.Println(err)
				return
			}
			if msg.IsText() {
				if strings.Contains(msg.Content, "@所有人") || strings.Comtaions()
				fmt.Println(groupSender.NickName, ":", msg.Content)
				fmt.Println(push(groupSender.NickName, defaultPriority, msg.Content))

			}
		}
	}

	bot.Block()
}

func push(pushTitle string, pushPriority string, pushMessage string) string {

	gotify := exec.Command("gotify", "push", "-t", pushTitle, "-p", pushPriority, pushMessage)
	output, err := gotify.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return string(output)
}
