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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lsytj0413/tyche/util"
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
	url = "http://kaijiang.500.com/shtml/ssq/"
)

func termToString(term uint32) string {
	return fmt.Sprintf("%05d", term)
}

func requestDocumentContent() (content string, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	content, err = util.DoRequest(request)
	return
}

// FetchTermList will fetch all terms
func FetchTermList() (terms []uint32, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	content, err := util.DoRequest(request)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return
	}

	termNodes := doc.Find(".kj_main01_right .kjxq_box02 .iSelectBox .iSelectList a")
	termNodesLength := termNodes.Length()
	if termNodesLength < 1 {
		err = errors.New("[.kj_main01_right .kjxq_box02 .iSelectBox .iSelectList a] select length zero")
		return
	}

	terms = make([]uint32, termNodesLength)
	var termValue int
	for i := 0; i < termNodesLength; i++ {
		s := goquery.NewDocumentFromNode(termNodes.Get(i))
		termValue, err = strconv.Atoi(s.Text())
		if err != nil {
			htmlValue, _ := s.Html()
			err = fmt.Errorf("select node is unexpected format: %v, %v", htmlValue, err)
			return
		}

		terms[termNodesLength-i-1] = uint32(termValue)
	}

	return
}

// FetchFromTerm will fetch award data at term
func FetchFromTerm(term uint32) (award *Award, err error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s.shtml", url, termToString(term)), nil)
	if err != nil {
		return
	}

	content, err := util.DoRequest(request)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return
	}

	termTitleNode := doc.Find(".kj_main01_right .kj_tablelist02 .td_title01 span")
	if termTitleNode.Length() != 2 {
		err = fmt.Errorf("[.kj_main01_right .kj_tablelist02 .td_title01 span] select length doesnot equal 2")
		return
	}

	termValueNode := termTitleNode.Eq(0).Find("strong")
	termValue, err := strconv.Atoi(termValueNode.Text())
	if err != nil {
		htmlValue, _ := termValueNode.Html()
		err = fmt.Errorf("select node is unexpected format: %v, %v", htmlValue, err)
		return
	}
	if uint32(termValue) != term {
		err = fmt.Errorf("termValue[%d] from html doesnot equal to args term[%d]", termValue, term)
		return
	}

	// termTimeNodeText := termTitleNode.Eq(1).Text()

	return
}
