package nas

import (
	"bytes"
	tokens "github.com/kubernetes-csi/csi-driver-nfs/pkg/nhn/token"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"log"
	"net/http"
)

type NasCreateResponse struct {
	Header struct {
		IsSuccessful  bool   `json:"isSuccessful"`
		ResultCode    int    `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
	} `json:"header"`
	Volume struct {
		Acl         []string `json:"acl"`
		CreatedAt   string   `json:"createdAt"`
		Description string   `json:"description"`
		Id          string   `json:"id"`
		Interfaces  []struct {
			Id       string `json:"id"`
			Path     string `json:"path"`
			Status   string `json:"status"`
			SubnetId string `json:"subnetId"`
		} `json:"interfaces"`
		Name           string `json:"name"`
		SizeGb         int    `json:"sizeGb"`
		SnapshotPolicy struct {
			MaxScheduledCount int `json:"maxScheduledCount"`
			ReservePercent    int `json:"reservePercent"`
			Schedule          struct {
				Time       string `json:"time"`
				TimeOffset string `json:"timeOffset"`
				Weekdays   []int  `json:"weekdays"`
			} `json:"schedule"`
		} `json:"snapshotPolicy"`
		Status    string `json:"status"`
		UpdatedAt string `json:"updatedAt"`
	} `json:"volume"`
}

type NasInfo struct {
	Volume struct {
		Acl         []string `json:"acl"`
		Description string   `json:"description"`
		Interfaces  []struct {
			SubnetId string `json:"subnetId"`
		} `json:"interfaces"`
		Name           string `json:"name"`
		SizeGb         int    `json:"sizeGb"`
		SnapshotPolicy struct {
			MaxScheduledCount int `json:"maxScheduledCount"`
			ReservePercent    int `json:"reservePercent"`
			Schedule          struct {
				Time       string `json:"time"`
				TimeOffset string `json:"timeOffset"`
				Weekdays   []int  `json:"weekdays"`
			} `json:"schedule"`
		} `json:"snapshotPolicy"`
	} `json:"volume"`
}

func checkReadFile(e error) {
	if e != nil {
		panic(e)
	}
}

func createVolume() {
	cfg := tokens.ReadConfig()

	klog.Warningf("KKJ Print inside paramNHNCloudNFS : %v %v %v %v", auth_url, trust_id, trustee_user_id, trustee_password)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pvc, err := clientset.CoreV1().PersistentVolumeClaims(parameters[pvcNamespaceKey]).Get(context.TODO(), parameters[pvcNameKey], v1.GetOptions{})

	klog.Warningf("KKJ pvc annotaion info %v", pvc.Annotations)
	klog.Warningf("KKJ pvc annotaion info nfs : %v", pvc.Annotations["nfs"])

	// data := make(map[string]TokenInfo)
	// data := TokenInfo{}
	// errr := json.Unmarshal([]byte(body), &data)

	//////////////////

	nasbody := pvc.Annotations["nfs"]

	nas_data := NasInfo{}
	nas_err := json.Unmarshal([]byte(nasbody), &nas_data)

	if nas_err != nil {
		log.Fatalln(nas_err)
	}

	nas_data_json, _ := json.Marshal(nas_data)
	klog.Warningf("KKJ marsharl data : %v", string(nas_data_json))
	buff = bytes.NewBuffer(nas_data_json)

	client = &http.Client{}

	req, err = http.NewRequest("POST", "http://stg-online-nas.kr2-t1.cloud.toastoven.net/v1/volumes", buff)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", token)

	resp, err = client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	data := NasCreateResponse{}
	errr := json.Unmarshal([]byte(body), &data)

	if errr != nil {
		log.Println(err)
	}

	datadd := string(body)
	klog.Warningf("KKJ output : %v", datadd)
	klog.Warningf("KKJ response output : %v", data.Header.IsSuccessful)

	defer resp.Body.Close()

	///////////////////////

}
