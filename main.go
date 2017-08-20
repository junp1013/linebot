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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var bot *linebot.Client

func main() {
	var err error
	//line bot
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
			replytoken := event.ReplyToken
			switch message := event.Message.(type) {
			//text message
			case *linebot.TextMessage:
				orimsg := message.Text
				//leave the chat
				if strings.Contains(strings.ToUpper(orimsg), "LCYBYE") {
					//get group id
					srcty := event.Source.Type
					if srcty == linebot.EventSourceTypeGroup {
						grpid := event.Source.GroupID
						if _, err := bot.LeaveGroup(grpid).Do(); err != nil {
							log.Print(err)
						}
					}
					if srcty == linebot.EventSourceTypeRoom {
						rid := event.Source.RoomID
						if _, err := bot.LeaveRoom(rid).Do(); err != nil {
							log.Print(err)
						}
					}

				}
				//button menu
				if strings.Contains(strings.ToUpper(orimsg), "LCYMENU") {
					imageURL := "https://pkget.com/images/skill/136.png"
					template := linebot.NewButtonsTemplate(
						imageURL, "Menu", "選單",
						linebot.NewURITemplateAction("開啟測試網址", "http://www.pokemon.com/us/pokedex/eevee"),
						//												linebot.NewPostbackTemplateAction("在哪?", "where", ""),
						linebot.NewPostbackTemplateAction("在哪?", "where", "你在哪?"),
						linebot.NewMessageTemplateAction("誇一下", "我好棒棒!"),
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
				//location message
			case *linebot.LocationMessage:
				lalo := fmt.Sprint(message.Latitude) + "," + fmt.Sprint(message.Longitude)
				//google search anything nearby
				nearbysearch(lalo, replytoken)

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
	if strings.Contains(orimsg, "哪") {
		bot.ReplyMessage(replytoken, linebot.NewLocationMessage("訊息位置", "111台北市士林區忠誠路二段55號", 25.111857, 121.531367)).Do()
	}
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

func nearbysearch(lalo, replytoken string) {
	var (
		//apiKey = flag.String("key", "", "API Key for using Google Maps API.")
		//	clientID  = flag.String("client_id", "", "ClientID for Maps for Work API access.")
		//	signature = flag.String("signature", "", "Signature for Maps for Work API access.")
		location = flag.String("location", lalo, "The latitude/longitude around which to retrieve place information. This must be specified as latitude,longitude.")
		radius   = flag.Uint("radius", 500, "Defines the distance (in meters) within which to bias place results. The maximum allowed radius is 50,000 meters.")
		//	keyword   = flag.String("keyword", "", "Specifies the language in which to return results. Optional.")
		language = flag.String("language", "zh-TW", "The language in which to return results.")
		//	minPrice  = flag.String("minprice", "", "Restricts results to only those places within the specified price level.")
		//	maxPrice  = flag.String("maxprice", "", "Restricts results to only those places within the specified price level.")
		//	name      = flag.String("name", "", "One or more terms to be matched against the names of places, separated with a space character.")
		openNow   = flag.Bool("open_now", true, "Restricts results to only those places that are open for business at the time the query is sent.")
		rankBy    = flag.String("rankby", "distance", "Specifies the order in which results are listed. Valid values are prominence or distance.")
		placeType = flag.String("type", "restaurant", "Restricts the results to places matching the specified type.")
		pageToken = flag.String("pagetoken", "", "Set to retrieve the next next page of results.")
	)
	flag.Parse()
	//google maps
	gmaps, err := maps.NewClient(maps.WithAPIKey(os.Getenv("ApiKey")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.NearbySearchRequest{
		Radius: *radius,
		//		Keyword:   *keyword,
		Language: *language,
		//		Name:      *name,
		OpenNow:   *openNow,
		PageToken: *pageToken,
	}

	parseLocation(*location, r)
	//	parsePriceLevels(*minPrice, *maxPrice, r)
	parseRankBy(*rankBy, r)
	parsePlaceType(*placeType, r)

	resp, err := gmaps.NearbySearch(context.Background(), r)

	result := resp.Results[0]
	if result.Name != "" {
		rname := result.Name
		radd := result.FormattedAddress
		//		rla := result.Geometry.Location.Lat
		//		rlo := result.Geometry.Location.Lng
		bot.ReplyMessage(replytoken, linebot.NewTextMessage("google map: "+rname+","+radd)).Do()
	}
}

func parseLocation(location string, r *maps.NearbySearchRequest) {
	if location != "" {
		actual, err := maps.ParseLatLng(location)
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}
		r.Location = &actual
	}
}

func parseRankBy(rankBy string, r *maps.NearbySearchRequest) {
	switch rankBy {
	case "prominence":
		r.RankBy = maps.RankByProminence
		return
	case "distance":
		r.RankBy = maps.RankByDistance
		return
	default:
		return
	}
}

func parsePlaceType(placeType string, r *maps.NearbySearchRequest) {
	if placeType != "" {
		t, err := maps.ParsePlaceType(placeType)
		if err != nil {
			log.Print(err)
		}
		r.Type = t
	}
}
