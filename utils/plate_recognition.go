package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"parking/config"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type PlateResult struct {
	Success   bool   `json:"success"`
	PlateNo   string `json:"plate_no"`
	PlateType string `json:"plate_type"`
	PlateColor string `json:"plate_color"`
	Confidence float64 `json:"confidence"`
	Message   string `json:"message"`
}

type AliyunResponse struct {
	RequestId string `json:"RequestId"`
	Data      struct {
		Plates []struct {
			PlateNumber string `json:"PlateNumber"`
			PlateType   string `json:"PlateType"`
			PlateColor  string `json:"PlateColor"`
			Confidence  float64 `json:"Confidence"`
		} `json:"Plates"`
	} `json:"Data"`
}

func RecognizePlateFromImage(imagePath string) (*PlateResult, error) {
	imgData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("读取图片失败: %v", err)
	}
	base64Img := base64.StdEncoding.EncodeToString(imgData)
	return RecognizePlate(base64Img)
}

func RecognizePlate(base64Img string) (*PlateResult, error) {
	config := config.AppConfig.Aliyun
	
	if config.AccessKeyID == "your_access_key_id" || config.AccessKeySecret == "your_access_key_secret" {
		return &PlateResult{
			Success:   false,
			Message:   "阿里云API未配置，请先在config.yaml中配置AccessKey",
		}, nil
	}

	client, err := sdk.NewClientWithAccessKey("cn-shanghai", config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云客户端失败: %v", err)
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = config.Endpoint
	request.Version = "2019-12-30"
	request.ApiName = "RecognizeLicensePlate"
	request.QueryParams["ImageURL"] = ""
	request.FormParams["ImageBase64"] = base64Img

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, fmt.Errorf("调用阿里云API失败: %v", err)
	}

	var result AliyunResponse
	err = json.Unmarshal(response.GetHttpContentBytes(), &result)
	if err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v", err)
	}

	if len(result.Data.Plates) > 0 {
		plate := result.Data.Plates[0]
		return &PlateResult{
			Success:   true,
			PlateNo:   plate.PlateNumber,
			PlateType: plate.PlateType,
			PlateColor: plate.PlateColor,
			Confidence: plate.Confidence,
			Message:   "识别成功",
		}, nil
	}

	return &PlateResult{
		Success:   false,
		Message:   "未识别到车牌",
	}, nil
}

func RecognizePlateWithMock() (*PlateResult, error) {
	plates := []string{
		"京A12345", "沪B67890", "粤C11111", "川D22222", "鄂E33333",
	}
	
	randPlate := plates[int(time.Now().Unix())%len(plates)]
	
	return &PlateResult{
		Success:   true,
		PlateNo:   randPlate,
		PlateType: "蓝牌",
		PlateColor: "蓝",
		Confidence: 99.5,
		Message:   "模拟识别成功",
	}, nil
}

func GenerateRecordNo() string {
	now := time.Now()
	return fmt.Sprintf("RK%s%06d", 
		now.Format("20060102150405"), 
		now.Nanosecond()/1000%1000000)
}

func GenerateFinanceNo() string {
	now := time.Now()
	return fmt.Sprintf("FJ%s%06d", 
		now.Format("20060102150405"), 
		now.Nanosecond()/1000%1000000)
}

func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip != "" {
			ips := strings.Split(ip, ",")
			if len(ips) > 0 {
				ip = strings.TrimSpace(ips[0])
			}
		}
	}
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}
