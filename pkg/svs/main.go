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
	"fmt"
	"os"
)

// Main is entrance of Server, it will block until Server is closed
func Main() {
	s, _ := New()
	ch, err := s.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error At Server Start: %s", err.Error())
		os.Exit(1)
	}

	<-ch
}
