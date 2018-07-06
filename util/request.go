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

package util

import (
	"net/http"

	"github.com/lsytj0413/tyche/codec"
)

// DoRequest will process request and return response body
func DoRequest(request *http.Request) (content string, err error) {
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	content, err = codec.ToUtf8(resp.Body)
	return
}
