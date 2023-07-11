package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kubernetes-csi/csi-driver-nfs/pkg/nhn/client"
	gcfg "gopkg.in/gcfg.v1"
	"io"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"os"
	"reflect"
)

type Config struct {
	Global client.AuthOpts
}

type TokenInfo struct {
	Token struct {
		Methods []string `json:"methods"`
		Roles   []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"roles"`
		Project struct {
			Domain struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
		} `json:"project"`
		Catalog []struct {
			Endpoints []struct {
				Url       string `json:"url"`
				Region    string `json:"region"`
				RegionId  string `json:"region_id"`
				Interface string `json:"interface"`
			} `json:"endpoints"`
			Type string `json:"type"`
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"catalog"`
		User struct {
			Domain struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			}
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		ExpiresAt string   `json:"expires_at"`
		IssuedAt  string   `json:"issued_at"`
		AuditIds  []string `json:"audit_ids"`
	} `json:"token"`
}

type AuthBody struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	Identity Identity `json:"identity"`
	Scope    Scope    `json:"scope"`
}

type Identity struct {
	Methods  []string `json:"methods"`
	PassWord PassWord `json:"password"`
}

type PassWord struct {
	User User `json:"user"`
}

type User struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

type Scope struct {
	OsTrust OsTrust `json:"OS-TRUST:trust"`
}

type OsTrust struct {
	Id string `json:"id"`
}

const (
	configFilePath = "/etc/config/secret-nfsplugin"
)

func ReadConfig(config io.Reader) (Config, error) {
	if config == nil {
		return Config{}, fmt.Errorf("no cloud provider config file given")
	}

	var cfg Config

	err := gcfg.FatalOnly(gcfg.ReadInto(&cfg, config))
	if err != nil {
		return Config{}, err
	}

	//klog.V(5).Infof("Config, loaded from the config file:")
	client.LogCfg(cfg.Global)

	return cfg, err
}

func GetToken() (TokenInfo, error) {
	f, err := os.Open(configFilePath)

	if err != nil {
		return TokenInfo{}, err
	}

	cfg, err := ReadConfig(f)

	if err != nil {
		return TokenInfo{}, err
	}

	client := &http.Client{}

	authBody := AuthBody{
		Auth: Auth{
			Identity: Identity{
				Methods: []string{"password"},
				PassWord: PassWord{
					User: User{
						Id:       cfg.Global.TrustID,
						Password: cfg.Global.Password,
					},
				},
			},
			Scope: Scope{
				OsTrust: OsTrust{
					Id: cfg.Global.TrustID,
				},
			},
		},
	}

	authBodyJson, _ := json.Marshal(authBody)
	buff := bytes.NewBuffer(authBodyJson)

	req, err := http.NewRequest("POST", cfg.Global.AuthURL, buff)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	token := resp.Header.Get("X-Subject-Token")
	print(reflect.TypeOf(token))

	klog.Warningf("KKJ token get %v", token)

	if err != nil {
		log.Fatalln(err)
	}

	return token, nil
}
