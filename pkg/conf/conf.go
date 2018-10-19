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

// Package conf provide iploc application config object
package conf

// Config is iploc server config instance
type Config struct {
	// 配置项
	Name                   string `json:"name"`
	DefaultListenClientURL string `json:"listenClientUrl"`
	IsDebug                bool
	IsPprof                bool

	// 客户端证书
	ClientTLSInfo TLSInfo
	IsTLSEnable   bool

	WxAppID          string
	WxAppSecret      string
	WxToken          string
	WxEncodingAESKey string
}

// TLSInfo is tls certificate info
type TLSInfo struct {
	CertFile       string
	KeyFile        string
	TrustedCAFile  string
	ClientCertAuth bool
	CRLFile        string

	InsecureSkipVerify bool
}

const (
	defaultName = "svs"
)

// New will construct a Config instance
func New() *Config {
	c := &Config{
		Name: defaultName,
	}

	return c
}
