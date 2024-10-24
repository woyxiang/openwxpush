package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/joho/godotenv"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	bot.SyncCheckCallback = nil                      //关闭心跳检测
	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("./storage.json")

	defer push("程序结束", "1", "程序结束") // 确保无论程序如何结束，都推送一条消息

	defer func(reloadStorage io.ReadWriteCloser) {
		err := reloadStorage.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(reloadStorage)

	// 登录
	if err := bot.HotLogin(reloadStorage); err != nil {
		fmt.Println("热登陆失败，尝试免扫码登录")
		err := bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())
		if err != nil {
			return
		}
	}

	fmt.Println("登陆成功")
	defaultPriority := "5"
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
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
			var group *openwechat.User
			var err error
			if msg.IsSendByGroup() {
				if msg.IsSendBySelf() {
					group, err = msg.Receiver()
				} else {
					group, err = msg.Sender()
				}
			}
			var groupName string
			if group != nil && group.NickName != "" {
				groupName = group.NickName
			}

			groupSender, err := msg.SenderInGroup()

			if err != nil {
				fmt.Println(err)
				return
			}
			groupNamesToReceive := strings.Split(os.Getenv("GROUP_NAME"), ";")
			//只接收指定群组和@所有人的消息
			if contains(groupName, groupNamesToReceive) || strings.Contains(msg.Content, "@所有人") {
				fmt.Println(groupSender.NickName, ":", msg.Content)
				fmt.Println(push(groupSender.NickName, defaultPriority, msg.Content))
			}

		}
	}

	_ = bot.Block()

}

func push(pushTitle string, pushPriority string, pushMessage string) string {

	gotify := exec.Command("gotify", "push", "-t", pushTitle, "-p", pushPriority, pushMessage)
	output, err := gotify.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	return string(output)
}
func contains(item string, slice []string) bool {
	// 确保切片是有序的
	sort.Strings(slice)

	// 使用 SearchStrings 查找元素
	i := sort.SearchStrings(slice, item)
	return i < len(slice) && slice[i] == item
}
