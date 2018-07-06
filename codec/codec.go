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

package codec

import (
	"io"
	"io/ioutil"

	"github.com/axgle/mahonia"
)

var (
	decoder mahonia.Decoder
)

// ToUtf8 convert gb18030 inbuf to utf8 string
func ToUtf8(r io.Reader) (string, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return decoder.ConvertString(string(buf)), nil
}

func init() {
	decoder = mahonia.NewDecoder("gb18030")
	if decoder == nil {
		panic("No gb18030 Decoder")
	}
}
