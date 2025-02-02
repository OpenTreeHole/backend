package sensitive

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentreehole/backend/common"
	"github.com/spf13/viper"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/image/v5"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/image/v5/check"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/label"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/label/request"
	v5 "github.com/yidun/yidun-golang-sdk/yidun/service/antispam/text"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/text/v5/check/sync/single"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	//TypeHole   = "Hole"
	//TypeFloor  = "Floor"
	//TypeTag    = "Tag"
	TypeImage  = "Image"
	TypeReview = "Review"
	TypeTitle  = "Title"
)

var checkTypes = []string{TypeImage, TypeReview, TypeTitle}

type ParamsForCheck struct {
	Content  string
	Id       int64
	TypeName string
}

type ResponseForCheck struct {
	Pass   bool   `json:"content"`
	Labels []int  `json:"labels,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func CheckSensitive(params ParamsForCheck) (resp *ResponseForCheck, err error) {
	images, clearContent, err := findImagesInMarkdownContent(params.Content)
	if err != nil {
		return nil, err
	}
	if len(images) != 0 {
		for _, img := range images {
			imgUrl, err := url.Parse(img)
			if err != nil {
				return nil, err
			}
			host := viper.GetString(common.EnvExternalImageHost)
			if host != "" {
				imgUrl.Host = host
			}
			ret, err := checkSensitiveImage(ParamsForCheck{
				Content:  imgUrl.String(),
				Id:       time.Now().UnixNano(),
				TypeName: TypeImage,
			})
			if err != nil {
				return nil, err
			}
			if !ret.Pass {
				return ret, nil
			}
		}
	}

	contained, reason := containsUnsafeURL(clearContent)
	if contained {
		return &ResponseForCheck{
			Pass:   false,
			Labels: nil,
			Detail: "不允许使用外部链接" + reason,
		}, nil
	}
	params.Content = strings.TrimSpace(removeIDReprInContent(clearContent))
	if params.Content == "" {
		return &ResponseForCheck{
			Pass:   true,
			Labels: nil,
			Detail: "",
		}, nil
	}

	return CheckSensitiveText(params)
}

func CheckSensitiveText(params ParamsForCheck) (resp *ResponseForCheck, err error) {
	if !checkType(params) {
		return nil, fmt.Errorf("invalid type for sensitive check")
	}

	request := single.NewTextCheckRequest(viper.GetString(common.EnvYiDunBusinessIdText))
	var textCheckClient *v5.TextClient
	if viper.GetString(common.EnvProxyUrl) != "" {
		var proxyUrl *url.URL
		proxyUrl, err := url.Parse(viper.GetString(common.EnvProxyUrl))
		if err != nil {
			return nil, err
		}
		textCheckClient = v5.NewTextClientWithAccessKeyWithProxy(viper.GetString(common.EnvYiDunSecretId), viper.GetString(common.EnvYiDunSecretKey), http.ProxyURL(proxyUrl))
	} else {
		textCheckClient = v5.NewTextClientWithAccessKey(viper.GetString(common.EnvYiDunSecretId), viper.GetString(common.EnvYiDunSecretKey))
	}

	request.SetDataID(strconv.FormatInt(params.Id, 10) + "_" + params.TypeName)
	request.SetContent(params.Content)
	request.SetTimestamp(time.Now().UnixMilli())

	response, err := textCheckClient.SyncCheckText(request)
	if err != nil {
		// 处理错误并打印日志
		common.RequestLog(fmt.Sprintf("sync request error:%+v", err.Error()), params.TypeName, params.Id, false)
		return &ResponseForCheck{Pass: false}, nil
	}

	resp = &ResponseForCheck{}
	if response.GetCode() == 200 {

		if *response.Result.Antispam.Suggestion == 0 {
			common.RequestLog("Sensitive text check response code is 200", params.TypeName, params.Id, true)
			resp.Pass = true
			return
		}

		common.RequestLog("Sensitive text check response code is 200", params.TypeName, params.Id, false)
		resp.Pass = false
		var sensitiveDetailBuilder strings.Builder
		sensitiveLabelMap.RLock()
		defer sensitiveLabelMap.RUnlock()
		for _, label := range response.Result.Antispam.Labels {
			if label.Label == nil {
				continue
			}
			resp.Labels = append(resp.Labels, *label.Label)
			// response != nil && response.Result != nil && response.Result.Antispam != nil &&
			//if response.Result.Antispam.SecondLabel != nil && response.Result.Antispam.ThirdLabel != nil {
			//	str := *response.Result.Antispam.SecondLabel + " " + *response.Result.Antispam.ThirdLabel
			//}
			labelNumber := *label.Label
			if sensitiveLabelMap.data[labelNumber] != nil {
				sensitiveDetailBuilder.WriteString("{")
				sensitiveDetailBuilder.WriteString(sensitiveLabelMap.label[labelNumber])
				sensitiveDetailBuilder.WriteString("}")
			}

			if label.SubLabels != nil {
				for _, subLabel := range label.SubLabels {
					if sensitiveLabelMap.data[labelNumber] != nil {
						if subLabel.SubLabel != nil {
							sensitiveDetailBuilder.WriteString("[" + sensitiveLabelMap.data[labelNumber][*subLabel.SubLabel] + "]")
						}
						if subLabel.SecondLabel != nil {
							sensitiveDetailBuilder.WriteString("[" + sensitiveLabelMap.data[labelNumber][*subLabel.SecondLabel] + "]")
						}
						if subLabel.ThirdLabel != nil {
							sensitiveDetailBuilder.WriteString("[" + sensitiveLabelMap.data[labelNumber][*subLabel.ThirdLabel] + "]")
						}
					}

					if subLabel.Details != nil && subLabel.Details.HitInfos != nil {
						for _, hitInfo := range subLabel.Details.HitInfos {
							if hitInfo.Value == nil {
								continue
							}
							if sensitiveDetailBuilder.Len() != 0 {
								sensitiveDetailBuilder.WriteString("\n")
							}
							sensitiveDetailBuilder.WriteString(*hitInfo.Value)
						}
					}
				}
			}
		}
		if sensitiveDetailBuilder.Len() == 0 {
			sensitiveDetailBuilder.WriteString("文本敏感，未知原因")
		}
		resp.Detail = sensitiveDetailBuilder.String()
		return
	}

	common.RequestLog("Sensitive text check http response code is not 200", params.TypeName, params.Id, false)
	resp.Pass = false
	return
}

func checkSensitiveImage(params ParamsForCheck) (resp *ResponseForCheck, err error) {
	// 设置易盾内容安全分配的businessId
	imgUrl := params.Content

	request := check.NewImageV5CheckRequest(viper.GetString(common.EnvYiDunBusinessIdImage))

	// 实例化一个textClient，入参需要传入易盾内容安全分配的secretId，secretKey
	var imageCheckClient *image.ImageClient
	if viper.GetString(common.EnvProxyUrl) != "" {
		var proxyUrl *url.URL
		proxyUrl, err := url.Parse(viper.GetString(common.EnvProxyUrl))
		if err != nil {
			return nil, err
		}
		imageCheckClient = image.NewImageClientWithAccessKeyWithProxy(viper.GetString(common.EnvYiDunSecretId), viper.GetString(common.EnvYiDunSecretKey), http.ProxyURL(proxyUrl))
	} else {
		imageCheckClient = image.NewImageClientWithAccessKey(viper.GetString(common.EnvYiDunSecretId), viper.GetString(common.EnvYiDunSecretKey))
	}

	imageInst := check.NewImageBeanRequest()
	imageInst.SetData(imgUrl)
	imageInst.SetName(strconv.FormatInt(params.Id, 10) + "_" + params.TypeName)
	// 设置图片数据的类型，1：图片URL，2:图片BASE64
	imageInst.SetType(1)

	imageBeans := []check.ImageBeanRequest{*imageInst}
	request.SetImages(imageBeans)

	response, err := imageCheckClient.ImageSyncCheck(request)
	if err != nil {
		// 处理错误并打印日志
		common.RequestLog(fmt.Sprintf("sync request error:%+v", err.Error()), params.TypeName, params.Id, false)
		// TODO: 通知管理员
		return &ResponseForCheck{Pass: false}, nil
	}

	resp = &ResponseForCheck{}
	if response.GetCode() == 200 {
		if len(*response.Result) == 0 {
			return nil, fmt.Errorf("sensitive image check returns empty response")
		}

		if *((*response.Result)[0].Antispam.Suggestion) == 0 {
			common.RequestLog("Sensitive image check response code is 200", params.TypeName, params.Id, true)
			resp.Pass = true
			return
		}

		common.RequestLog("Sensitive image check response code is 200", params.TypeName, params.Id, false)
		resp.Pass = false
		for _, label := range *((*response.Result)[0].Antispam.Labels) {
			resp.Labels = append(resp.Labels, *label.Label)
		}
		var sensitiveDetailBuilder strings.Builder
		sensitiveLabelMap.RLock()
		defer sensitiveLabelMap.RUnlock()
		for _, result := range *response.Result {
			if result.Antispam != nil && result.Antispam.Labels != nil {
				for _, label := range *result.Antispam.Labels {
					if label.Label == nil || label.SubLabels == nil {
						continue
					}

					labelNumber := *label.Label
					if sensitiveLabelMap.data[labelNumber] != nil {
						sensitiveDetailBuilder.WriteString("{")
						sensitiveDetailBuilder.WriteString(sensitiveLabelMap.label[labelNumber])
						sensitiveDetailBuilder.WriteString("}")
					}

					if label.SubLabels != nil {
						for _, subLabel := range *label.SubLabels {
							if sensitiveLabelMap.data[labelNumber] != nil {
								if subLabel.SubLabel != nil {
									sensitiveDetailBuilder.WriteString("[" + sensitiveLabelMap.data[labelNumber][*subLabel.SubLabel] + "]")
								}
								if subLabel.SecondLabel != nil {
									sensitiveDetailBuilder.WriteString("[" + sensitiveLabelMap.data[labelNumber][*subLabel.SecondLabel] + "]")
								}
								if subLabel.ThirdLabel != nil {
									sensitiveDetailBuilder.WriteString("[" + sensitiveLabelMap.data[labelNumber][*subLabel.ThirdLabel] + "]")
								}
							}

							if subLabel.Details != nil && subLabel.Details.HitInfos != nil {
								for _, hitInfo := range *subLabel.Details.HitInfos {
									if hitInfo.Group != nil {
										sensitiveDetailBuilder.WriteByte(' ')
										sensitiveDetailBuilder.WriteString(*hitInfo.Group)
									}
									if hitInfo.Value != nil {
										sensitiveDetailBuilder.WriteByte(' ')
										sensitiveDetailBuilder.WriteString(*hitInfo.Value)
									}
									if hitInfo.Word != nil {
										sensitiveDetailBuilder.WriteByte(' ')
										sensitiveDetailBuilder.WriteString(*hitInfo.Word)
									}
								}
							}
						}
					}
				}
			}
			if result.Ocr != nil {
				if result.Ocr.Details != nil {
					for _, detail := range *result.Ocr.Details {
						if detail.Content == nil {
							continue
						}
						if sensitiveDetailBuilder.Len() != 0 {
							sensitiveDetailBuilder.WriteString("\n")
						}
						sensitiveDetailBuilder.WriteString(*detail.Content)
					}
				}
			}
			if result.Face != nil {
				if result.Face.Details != nil {
					for _, detail := range *result.Face.Details {
						if detail.FaceContents != nil {
							for _, faceContent := range *detail.FaceContents {
								if faceContent.Name == nil {
									continue
								}
								if sensitiveDetailBuilder.Len() != 0 {
									sensitiveDetailBuilder.WriteString("\n")
								}
								sensitiveDetailBuilder.WriteString(*faceContent.Name)
							}
						}
					}
				}
			}
		}
		if sensitiveDetailBuilder.Len() == 0 {
			sensitiveDetailBuilder.WriteString("图片敏感，未知原因")
		}
		resp.Detail = sensitiveDetailBuilder.String()
		return
	}

	common.RequestLog("Sensitive image check http response code is not 200", params.TypeName, params.Id, false)
	resp.Pass = false
	return
}

var sensitiveLabelMap struct {
	sync.RWMutex
	label      map[int]string
	data       map[int]map[string]string
	lastLength int
}

func InitSensitiveLabelMap() {
	// skip when bench

	// || viper.GetString(common.AuthUrl) == ""
	if viper.GetString(common.EnvMode) == "bench" {
		return
	}

	// 创建一个LabelQueryRequest实例
	request := request.NewLabelQueryRequest()

	// 实例化Client，入参需要传入易盾内容安全分配的AccessKeyId，AccessKeySecret
	labelClient := label.NewLabelClientWithAccessKey(viper.GetString(common.EnvYiDunAccessKeyId), viper.GetString(common.EnvYiDunAccessKeySecret))

	// 传入请求参数
	//设置返回标签的最大层级
	request.SetMaxDepth(3)
	//指定业务类型
	// request.SetBusinessTypes(&[]string{"1", "2"})
	//制定业务
	// request.SetBusinessID("SetBusinessID")
	// request.SetClientID("YOUR_CLIENT_ID")
	// request.SetLanguage("en")

	response, err := labelClient.QueryLabel(request)
	if err != nil {
		// 	log.Err(err).Str("model", "get admin").Msg("error sending auth server")
		common.RequestLog("Sensitive label init error", "label error", -1, false)
		return
	}

	if response.GetCode() != 200 {
		// log.Error().Str("model", "get admin").Msg("auth server response failed" + res.Status)
		common.RequestLog("Sensitive label init http response code is not 200", "label error", -1, false)
		return
	}

	responseByte, err := json.Marshal(response)
	if err != nil {
		common.RequestLog("Sensitive label Marshal error", "label error", -1, false)
		return
	}

	if sensitiveLabelMap.lastLength == len(responseByte) {
		common.RequestLog("Sensitive label unchanged", "label unchanged", 1, false)
		return
	}

	sensitiveLabelMap.Lock()
	defer sensitiveLabelMap.Unlock()
	sensitiveLabelMap.lastLength = len(responseByte)
	sensitiveLabelMap.label = make(map[int]string)
	sensitiveLabelMap.data = make(map[int]map[string]string)
	data := response.Data

	for _, label := range data {
		if label.Label == nil || label.Name == nil {
			continue
		}
		sensitiveLabelMap.label[*label.Label] = *label.Name
		labelNumber := *label.Label
		labelMap := make(map[string]string)
		for _, subLabel := range label.SubLabels {
			if subLabel.Code == nil || subLabel.Name == nil {
				continue
			}
			labelMap[*subLabel.Code] = *subLabel.Name
			for _, subSubLabel := range subLabel.SubLabels {
				if subSubLabel.Code == nil || subSubLabel.Name == nil {
					continue
				}
				labelMap[*subLabel.Code] = *subLabel.Name
			}
		}
		sensitiveLabelMap.data[labelNumber] = labelMap
	}
}

func UpdateSensitiveLabelMap(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			InitSensitiveLabelMap()
		}
	}
}
