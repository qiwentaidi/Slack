// 基于 https://github.com/HG-ha/ICP_Query 项目的接口服务数据获取
package icp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/utils/randutil"
	"strings"
	"time"

	"github.com/qiwentaidi/clients"
)

// 接口返回结构
type WebApiResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Params struct {
		PageNum  int       `json:"pageNum"`
		PageSize int       `json:"pageSize"`
		Total    int       `json:"total"`
		NextPage int       `json:"nextPage"`
		HasNext  bool      `json:"hasNextPage"`
		List     []WebData `json:"list"`
	} `json:"params"`
	Success bool `json:"success"`
}

// 每条记录的结构体
type WebData struct {
	Domain           string `json:"domain"`
	DomainID         int64  `json:"domainId"`
	MainID           int64  `json:"mainId"`
	MainLicence      string `json:"mainLicence"`
	ServiceID        int64  `json:"serviceId"`
	ServiceLicence   string `json:"serviceLicence"`
	UnitName         string `json:"unitName"`
	UpdateRecordTime string `json:"updateRecordTime"`
	NatureName       string `json:"natureName"`
	LimitAccess      string `json:"limitAccess"`
}

// const defaultApiAddress = "http://127.0.0.1:16181"

func FetchWebInfo(ctx context.Context, apiAddress string, query string) (*WebApiResponse, error) {
	var allList []WebData
	pageSize := 10
	pageNum := 1
	totalPages := 1
	var lastResponse WebApiResponse
	retryCount := 0
	maxRetries := 5

	for {
		queryURL := fmt.Sprintf("%s/query/web?search=%s&pageNum=%d&pageSize=%d", apiAddress, url.QueryEscape(query), pageNum, pageSize)
		resp, err := clients.SimpleGet(queryURL, clients.NewRestyClient(nil, true))
		if err != nil {
			retryCount++
			if retryCount > maxRetries {
				gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询备案信息重试超过%d次，返回空结果", query, maxRetries))
				return &WebApiResponse{}, nil
			}
			gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询备案信息出现错误: %v，重试第%d次", query, err, retryCount))
			time.Sleep(randutil.SleepRandTime(2))
			continue
		}

		bodyStr := string(resp.Body())
		// 如果不包含params参数时代码查询失败，需要重新查询
		if !strings.Contains(bodyStr, "params") {
			retryCount++
			if retryCount > maxRetries {
				gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询备案信息重试超过%d次，返回空结果", query, maxRetries))
				return &WebApiResponse{}, nil
			}
			gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询备案信息时出现错误: %s，重试第%d次", query, bodyStr, retryCount))
			time.Sleep(randutil.SleepRandTime(2))
			continue
		}

		var response WebApiResponse
		if err := json.Unmarshal(resp.Body(), &response); err != nil {
			retryCount++
			if retryCount > maxRetries {
				gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询备案信息重试超过%d次，返回空结果", query, maxRetries))
				return &WebApiResponse{}, nil
			}
			gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s解析备案信息时出现错误: %v，重试第%d次", query, err, retryCount))
			time.Sleep(randutil.SleepRandTime(2))
			continue
		}
		if !response.Success {
			break
		}

		if pageNum == 1 {
			lastResponse = response
			total := response.Params.Total
			if total > 0 {
				totalPages = (total + pageSize - 1) / pageSize
			}
		}

		allList = append(allList, response.Params.List...)

		if pageNum >= totalPages {
			break
		}

		pageNum++
		retryCount = 0 // 重置重试计数器
		time.Sleep(3000 * time.Millisecond)
	}

	lastResponse.Params.List = allList
	return &lastResponse, nil
}

type AppApiResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Params struct {
		EndRow           int           `json:"endRow"`
		FirstPage        int           `json:"firstPage"`
		HasNextPage      bool          `json:"hasNextPage"`
		HasPreviousPage  bool          `json:"hasPreviousPage"`
		IsFirstPage      bool          `json:"isFirstPage"`
		IsLastPage       bool          `json:"isLastPage"`
		LastPage         int           `json:"lastPage"`
		List             []structs.App `json:"list"`
		NavigatePages    int           `json:"navigatePages"`
		NavigatepageNums []int         `json:"navigatepageNums"`
		NextPage         int           `json:"nextPage"`
		PageNum          int           `json:"pageNum"`
		PageSize         int           `json:"pageSize"`
		Pages            int           `json:"pages"`
		PrePage          int           `json:"prePage"`
		Size             int           `json:"size"`
		StartRow         int           `json:"startRow"`
		Total            int           `json:"total"`
	} `json:"params"`
	Success bool `json:"success"`
}

func FetchAppInfo(ctx context.Context, apiAddress string, companyName string) (*AppApiResponse, error) {
	var allList []structs.App
	pageSize := 10
	pageNum := 1
	totalPages := 1
	var lastResponse AppApiResponse
	retryCount := 0
	maxRetries := 5

	for {
		queryURL := fmt.Sprintf("%s/query/app?search=%s&pageNum=%d&pageSize=%d", apiAddress, url.QueryEscape(companyName), pageNum, pageSize)
		resp, err := clients.SimpleGet(queryURL, clients.NewRestyClient(nil, true))
		if err != nil {
			retryCount++
			if retryCount > maxRetries {
				gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询App信息重试超过%d次，返回空结果", companyName, maxRetries))
				return &AppApiResponse{}, nil
			}
			gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询App信息出现错误: %v，重试第%d次", companyName, err, retryCount))
			time.Sleep(randutil.SleepRandTime(2))
			continue
		}

		bodyStr := string(resp.Body())
		// 如果不包含params参数时代码查询失败，需要重新查询
		if !strings.Contains(bodyStr, "params") {
			retryCount++
			if retryCount > maxRetries {
				gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询App信息重试超过%d次，返回空结果", companyName, maxRetries))
				return &AppApiResponse{}, nil
			}
			gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询App信息时出现错误: %s，重试第%d次", companyName, bodyStr, retryCount))
			time.Sleep(randutil.SleepRandTime(2))
			continue
		}

		var response AppApiResponse
		if err := json.Unmarshal(resp.Body(), &response); err != nil {
			retryCount++
			if retryCount > maxRetries {
				gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询App信息重试超过%d次，返回空结果", companyName, maxRetries))
				return &AppApiResponse{}, nil
			}
			gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s解析App信息时出现错误: %v，重试第%d次", companyName, err, retryCount))
			time.Sleep(randutil.SleepRandTime(2))
			continue
		}
		if !response.Success {
			break
		}

		if pageNum == 1 {
			lastResponse = response
			total := response.Params.Total
			if total > 0 {
				totalPages = (total + pageSize - 1) / pageSize
			}
		}

		allList = append(allList, response.Params.List...)

		if pageNum >= totalPages {
			break
		}

		pageNum++
		retryCount = 0 // 重置重试计数器
		time.Sleep(3000 * time.Millisecond)
	}

	lastResponse.Params.List = allList
	return &lastResponse, nil
}

type AppletApiResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Params struct {
		EndRow           int              `json:"endRow"`
		FirstPage        int              `json:"firstPage"`
		HasNextPage      bool             `json:"hasNextPage"`
		HasPreviousPage  bool             `json:"hasPreviousPage"`
		IsFirstPage      bool             `json:"isFirstPage"`
		IsLastPage       bool             `json:"isLastPage"`
		LastPage         int              `json:"lastPage"`
		List             []structs.Applet `json:"list"`
		NavigatePages    int              `json:"navigatePages"`
		NavigatepageNums []int            `json:"navigatepageNums"`
		NextPage         int              `json:"nextPage"`
		PageNum          int              `json:"pageNum"`
		PageSize         int              `json:"pageSize"`
		Pages            int              `json:"pages"`
		PrePage          int              `json:"prePage"`
		Size             int              `json:"size"`
		StartRow         int              `json:"startRow"`
		Total            int              `json:"total"`
	} `json:"params"`
	Success bool `json:"success"`
}

func FetchAppletInfo(ctx context.Context, apiAddress string, companyName string) (*AppletApiResponse, error) {
	var allList []structs.Applet
	pageSize := 10
	pageNum := 1
	totalPages := 1
	var lastResponse AppletApiResponse
	retryCount := 0
	maxRetries := 5

	for {
		queryURL := fmt.Sprintf("%s/query/mapp?search=%s&pageNum=%d&pageSize=%d", apiAddress, url.QueryEscape(companyName), pageNum, pageSize)
		resp, err := clients.SimpleGet(queryURL, clients.NewRestyClient(nil, true))
		if err != nil {
			retryCount++
			if retryCount > maxRetries {
				gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询小程序信息重试超过%d次，返回空结果", companyName, maxRetries))
				return &AppletApiResponse{}, nil
			}
			gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询小程序信息出现错误: %v，重试第%d次", companyName, err, retryCount))
			time.Sleep(randutil.SleepRandTime(2))
			continue
		}

		bodyStr := string(resp.Body())
		// 如果不包含params参数时代码查询失败，需要重新查询
		if !strings.Contains(bodyStr, "params") {
			retryCount++
			if retryCount > maxRetries {
				gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询小程序信息重试超过%d次，返回空结果", companyName, maxRetries))
				return &AppletApiResponse{}, nil
			}
			gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询小程序信息时出现错误: %s，重试第%d次", companyName, bodyStr, retryCount))
			time.Sleep(randutil.SleepRandTime(2))
			continue
		}

		var response AppletApiResponse
		if err := json.Unmarshal(resp.Body(), &response); err != nil {
			retryCount++
			if retryCount > maxRetries {
				gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s查询小程序信息重试超过%d次，返回空结果", companyName, maxRetries))
				return &AppletApiResponse{}, nil
			}
			gologger.DualLog(ctx, gologger.Level_DEBUG, fmt.Sprintf("[icp] %s解析小程序信息时出现错误: %v，重试第%d次", companyName, err, retryCount))
			time.Sleep(randutil.SleepRandTime(2))
			continue
		}
		if !response.Success {
			break
		}

		if pageNum == 1 {
			lastResponse = response
			total := response.Params.Total
			if total > 0 {
				totalPages = (total + pageSize - 1) / pageSize
			}
		}

		allList = append(allList, response.Params.List...)

		if pageNum >= totalPages {
			break
		}

		pageNum++
		retryCount = 0 // 重置重试计数器
		time.Sleep(3000 * time.Millisecond)
	}

	lastResponse.Params.List = allList
	return &lastResponse, nil
}
