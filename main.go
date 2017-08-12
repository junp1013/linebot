// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			//text message
			case *linebot.TextMessage:
				replytoken := event.ReplyToken
				orimsg := message.Text
				//leave the group id
				if strings.Contains(strings.ToUpper(orimsg), "LCYBYE") {
					//get group id
					grpid := event.Source.GroupID
					if _, err := bot.LeaveRoom(grpid).Do(); err != nil {
						log.Print(err)
					}
				}
				//button menu
				if strings.Contains(strings.ToUpper(orimsg), "LCYMENU") {
					imageURL := "https://pkget.com/images/skill/136.png"
					template := linebot.NewButtonsTemplate(
						imageURL, "My button sample", "Hello, my button",
						linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
						//						linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
						//						linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは"),
						//						linebot.NewMessageTemplateAction("Say message", "Rice=米"),
					)
					//					bot.ReplyMessage(replytoken, linebot.NewTemplateMessage("Menu", template)).Do()
					if _, err = bot.ReplyMessage(replytoken, linebot.NewTemplateMessage("Menu", template)).Do(); err != nil {
						log.Print(err)
					}
				}
				//find sender's username
				userid := event.Source.UserID
				if res, err := bot.GetProfile(userid).Do(); err != nil {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(userid+" said:"+message.Text)).Do()
					log.Print(err)
				} else {
					username := res.DisplayName
					replytext(userid, replytoken, username, orimsg)
				}
				//sticker message
			case *linebot.StickerMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage("2", "179")).Do(); err != nil {
					log.Print(err)
				}
			default:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage("2", "149")).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func replytext(userid, replytoken, username, orimsg string) {
	//check coming message key words and reply
	if strings.Contains(orimsg, "哈") {
		bot.ReplyMessage(replytoken, linebot.NewStickerMessage("1", "110")).Do()
	}
	if strings.Contains(orimsg, "抱") {
		bot.ReplyMessage(replytoken, linebot.NewStickerMessage("2", "157")).Do()
	}
	textmsg := username + " said:" + orimsg
	//my line id:Uf016e10434dee6b3f864be761f5f723c
	//green line id:Ud517dfdbfd690d483692fc3efc234b37
	//send multiple mesages at once
	if userid == "Uf016e10434dee6b3f864be761f5f723c" {
		bot.ReplyMessage(replytoken, linebot.NewTextMessage(textmsg), linebot.NewStickerMessage("2", "144")).Do()
	}
	bot.ReplyMessage(replytoken, linebot.NewTextMessage(textmsg), linebot.NewStickerMessage("2", "167")).Do()
	//	if _, err = bot.ReplyMessage(replytoken, linebot.NewTextMessage(username+" said:"+orimsg)).Do(); err != nil {
	//		log.Print(err)
	//	}
	//	switch str {
	//		}
}
