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

package main

import (
	"fmt"

	"github.com/lsytj0413/tyche/pkg/lottery/tcb"
)

func main() {
	termList, err := tcb.FetchTermList()
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = termList

	_, err = tcb.FetchFromTerm(18077)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("%+v\n", awards)
	fmt.Println("tyche")
}
