package model

var (
	DASHBOARD_PERMISSION = map[string]interface{}{
		"id":        "9710971640510357",
		"parentId":  "",
		"label":     "sys.menu.analysis",
		"name":      "Analysis",
		"type":      1,
		"route":     "home",
		"icon":      "ic-analysis",
		"order":     1,
		"component": "/dashboard/analysis/index.tsx",
	}

	OPENAI_PERMISSION = map[string]interface{}{
		"id":       "9100714781927721",
		"parentId": "",
		"label":    "sys.menu.token",
		"name":     "Token",
		"icon":     "ph:key",
		"type":     0,
		"route":    "token",
		"order":    2,
		"children": []map[string]interface{}{
			{
				"id":        "84269992294009655",
				"parentId":  "9100714781927721",
				"label":     "sys.menu.openai-token-management",
				"name":      "OpenaiToken",
				"type":      1,
				"route":     "openai-token",
				"component": "/token/openai/token/index.tsx",
			},
			{
				"id":        "84269992294009656",
				"parentId":  "9100714781927721",
				"hide":      false,
				"label":     "sys.menu.openai-account-management",
				"name":      "OpenaiAccount",
				"type":      1,
				"route":     "openai-account",
				"component": "/token/openai/account/index.tsx",
			},
			{
				"id":        "84269992294009657",
				"parentId":  "9100714781927721",
				"label":     "sys.menu.claude-token-management",
				"name":      "ClaudeToken",
				"type":      1,
				"route":     "claude-token",
				"component": "/token/claude/token/index.tsx",
			},
			{
				"id":        "84269992294009658",
				"parentId":  "9100714781927721",
				"hide":      false,
				"label":     "sys.menu.claude-account-management",
				"name":      "ClaudeAccount",
				"type":      1,
				"route":     "claude-account",
				"component": "/token/claude/account/index.tsx",
			},
			{
				"id":        "84269992294009659",
				"parentId":  "9100714781927721",
				"hide":      false,
				"label":     "sys.menu.user-management",
				"name":      "User",
				"type":      1,
				"route":     "user",
				"component": "/token/user/index.tsx",
			},
		},
	}

	PERMISSION_LIST = []map[string]interface{}{
		DASHBOARD_PERMISSION,
		OPENAI_PERMISSION,
	}

	ADMIN_ROLE = map[string]interface{}{
		"id":         "4281707933534332",
		"name":       "Admin",
		"label":      "Admin",
		"status":     1,
		"order":      1,
		"desc":       "Super Admin",
		"permission": PERMISSION_LIST,
	}
)
