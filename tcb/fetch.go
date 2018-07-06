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

package tcb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

var _ = goquery.NewDocumentFromReader

// Award 是双色球开奖结果
type Award struct {
	Term         uint32
	AwardDate    time.Time
	DeadlineDate time.Time
	Number       []uint8
	SalesVolume  uint64
	RemainBonus  uint64
	Pieces       []Piece
}

// AwardLevel 是奖项等级
type AwardLevel uint8

const (
	// FirstAward 是一等奖
	FirstAward = AwardLevel(1)
	// SecondAward 是二等奖
	SecondAward = AwardLevel(2)
	// ThirdAward 是三等奖
	ThirdAward = AwardLevel(3)
	// FourthAward 是四等奖
	FourthAward = AwardLevel(4)
	// FifthAward 是五等奖
	FifthAward = AwardLevel(5)
	// SixthAward 是六等奖
	SixthAward = AwardLevel(6)
)

// Piece 是单注开奖详情
type Piece struct {
	Level AwardLevel
	Count uint32
	Bonus uint32
}

const (
	url = "http://kaijiang.500.com/shtml/ssq/18077.shtml"
)

// Fetch will fetch award data from url
func Fetch() ([]Award, error) {
	awards := make([]Award, 0)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	v := mahonia.NewDecoder("gb18030")
	bodyStr := v.ConvertString(string(body))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		return nil, err
	}

	s1 := doc.Find(".kj_tablelist02")

	s1.Each(func(i int, s *goquery.Selection) {
		// fmt.Println(s.Text())
	})

	s2 := s1.Find("tr tr")

	var s3 *goquery.Selection
	s2.Each(func(i int, s *goquery.Selection) {
		if s3 == nil {
			if strings.Contains(s.Text(), "出球顺序") {
				s3 = s
			}
		}
		fmt.Println(i)
	})

	s4 := s3.Find("td")
	for i := s4.Length() - 1; i >= 0; i-- {
		n := s4.Get(i)
		// s5 := goquery.NewDocumentFromNode(n)
		// fmt.Println(s5.Text())
		fmt.Println(n.FirstChild.Data)
	}

	// fmt.Println(len(s4.Nodes))

	// fmt.Println(bodyStr)

	return awards, nil
}
