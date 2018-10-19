// Copyright (c) 2018 soren yang
//
// Licensed under the MIT License
// you may not use this file except in complicance with the License.
// You may obtain a copy of the License at
//
//     https://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package svs

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *server) Version(c *gin.Context) (interface{}, error) {
	return map[string]string{
		"Version":     "1.0.0",
		"Name":        "svs",
		"Description": "svs website server.",
	}, nil
}

func (s *server) Index(c *gin.Context) (interface{}, error) {
	return s.Version(c)
}

func (s *server) WxVerify(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	params := []string{s.c.WxToken, timestamp, nonce}
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	if signature == fmt.Sprintf("%x", h.Sum(nil)) {
		c.Writer.Write([]byte(echostr))
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		return
	}

	c.Writer.WriteHeader(http.StatusBadRequest)
}

// TextMessage struct
type TextMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   uint64   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgID        uint64   `xml:"MsgId"`
}

// TextReply struct
type TextReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   uint64   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
}

func (s *server) WxEntry(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	text := TextMessage{}
	err = xml.Unmarshal(body, &text)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	reply := &TextReply{
		ToUserName:   text.FromUserName,
		FromUserName: text.ToUserName,
		CreateTime:   uint64(time.Now().Second()) * 1000,
		MsgType:      "text",
		Content:      text.Content + " reply from tyche",
	}
	replyByte, err := xml.Marshal(reply)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(replyByte)
	c.Writer.Write(replyByte)
	c.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.Writer.WriteHeader(http.StatusOK)
}
