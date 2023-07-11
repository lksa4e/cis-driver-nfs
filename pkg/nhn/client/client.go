package client

import "k8s.io/klog/v2"

type AuthOpts struct {
	AuthURL  string `gcfg:"auth-url" mapstructure:"auth-url" name:"os-authURL" dependsOn:"os-password|os-trustID|os-applicationCredentialSecret|os-clientCertPath"`
	UserID   string `gcfg:"user-id" mapstructure:"user-id" name:"os-userID" value:"optional" dependsOn:"os-password"`
	Password string `name:"os-password" value:"optional" dependsOn:"os-domainID|os-domainName,os-projectID|os-projectName,os-userID|os-userName"`
	TrustID  string `gcfg:"trust-id" mapstructure:"trust-id" name:"os-trustID" value:"optional"`
	CAFile   string `gcfg:"ca-file" mapstructure:"ca-file" name:"os-certAuthorityPath" value:"optional"`
}

func LogCfg(authOpts AuthOpts) {
	klog.V(5).Infof("AuthURL: %s", authOpts.AuthURL)
	klog.V(5).Infof("UserID: %s", authOpts.UserID)
	klog.V(5).Infof("TrustID: %s", authOpts.TrustID)
	klog.V(5).Infof("CAFile: %s", authOpts.CAFile)
}
