package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
)

type BackendApiMeResponse struct {
	Object                   string           `json:"object"`
	ID                       string           `json:"id"`
	Email                    string           `json:"email"`
	Name                     string           `json:"name"`
	Picture                  string           `json:"picture"`
	Created                  int64            `json:"created"`
	PhoneNumber              interface{}      `json:"phone_number"`
	MfaFlagEnabled           bool             `json:"mfa_flag_enabled"`
	Amr                      []interface{}    `json:"amr"`
	Groups                   []interface{}    `json:"groups"`
	Orgs                     BackendApiMeOrgs `json:"orgs"`
	HasPaygProjectSpendLimit bool             `json:"has_payg_project_spend_limit"`
}

type BackendApiMeOrgs struct {
	Object string                `json:"object"`
	Data   []BackendApiMeOrgData `json:"data"`
}

type BackendApiMeOrgData struct {
	Object                         string               `json:"object"`
	ID                             string               `json:"id"`
	Created                        int64                `json:"created"`
	Title                          string               `json:"title"`
	Name                           string               `json:"name"`
	Description                    string               `json:"description"`
	Personal                       bool                 `json:"personal"`
	Settings                       BackendApiMeSettings `json:"settings"`
	ParentOrgID                    interface{}          `json:"parent_org_id"`
	IsDefault                      bool                 `json:"is_default"`
	Role                           string               `json:"role"`
	IsScaleTierAuthorizedPurchaser interface{}          `json:"is_scale_tier_authorized_purchaser"`
	IsScimManaged                  bool                 `json:"is_scim_managed"`
	Projects                       Projects             `json:"projects"`
	Groups                         []interface{}        `json:"groups"`
	Geography                      interface{}          `json:"geography"`
}

type BackendApiMeSettings struct {
	ThreadsUIVisibility      string `json:"threads_ui_visibility"`
	UsageDashboardVisibility string `json:"usage_dashboard_visibility"`
	DisableUserAPIKeys       bool   `json:"disable_user_api_keys"`
}

type Projects struct {
	Object string        `json:"object"`
	Data   []interface{} `json:"data"`
}

// ProcessBackendApiMeResponse 处理 /backend-api/me 的响应体
func ProcessBackendApiMeResponse(body []byte) ([]byte, error) {
	var respData BackendApiMeResponse
	err := json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}

	email := respData.Email
	targetEmail := "admin@openai.com"

	// 修改 Orgs.Data 中每个组织的 Description 字段中的邮箱
	for i, org := range respData.Orgs.Data {
		respData.Orgs.Data[i].Description = strings.ReplaceAll(org.Description, email, targetEmail)
	}
	// 修改顶层的 email 字段
	respData.Email = targetEmail
	// 清空 Picture 字段
	respData.Picture = ""
	respData.Name = "admin"

	// 将修改后的数据重新编码为 JSON
	modifiedBody, err := json.Marshal(respData)
	if err != nil {
		fmt.Printf("Failed to marshal modified response body, %v", err)
		return nil, err
	}

	return modifiedBody, nil
}
