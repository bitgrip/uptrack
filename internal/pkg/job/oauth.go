package job

import (
	"time"
)

type Oauth struct {
	AuthUrl string            `json:"auth_url,omitempty" yaml:"auth_url,omitempty"`
	Params  map[string]string `json:"params,omitempty" yaml:"params,omitempty"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type OauthResponse struct {
	AccessToken  string `json:"access_token,omitempty" yaml:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty" yaml:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty" yaml:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty" yaml:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty" yaml:"scope,omitempty"`
	ExpiresAt    time.Time
	Refresh      bool
}
