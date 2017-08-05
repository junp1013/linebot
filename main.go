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
				//find sender's username
				userid := event.Source.UserID
				if res, err := bot.GetProfile(userid).Do(); err != nil {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(userid+" said:"+message.Text)).Do()
					log.Print(err)
				} else {
					replytoken := event.ReplyToken
					username := res.DisplayName
					orimsg := message.Text
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

//my line id:Uf016e10434dee6b3f864be761f5f723c
func replytext(userid, replytoken, username, orimsg string) {
	if userid == "Uf016e10434dee6b3f864be761f5f723c" {
		bot.ReplyMessage(replytoken, linebot.NewTextMessage(username+"最棒!")).Do()
		bot.ReplyMessage(replytoken, linebot.NewStickerMessage("2", "172")).Do()
	} else {
		bot.ReplyMessage(replytoken, linebot.NewTextMessage(username+" said:"+orimsg)).Do()
	}

	//	if _, err = bot.ReplyMessage(replytoken, linebot.NewTextMessage(username+" said:"+orimsg)).Do(); err != nil {
	//		log.Print(err)
	//	}
	//	switch str {
	//		}
}
