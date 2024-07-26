package v1

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	Type      int    `json:"type" binding:"required" example:"1"`
	AccountId int64  `json:"accountId" binding:"required" example:"1"`
	Password  string `json:"password" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	LoginType   int                    `json:"type"`
	AccessToken string                 `json:"accessToken"`
	User        map[string]interface{} `json:"user"`
	LoginUrl    string                 `json:"loginUrl"`
}

type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}

type AddUserRequest struct {
	UniqueName     string `json:"uniqueName"`
	Password       string `json:"password"`
	ExpirationTime string `json:"expirationTime"`
	Enable         int    `json:"enable"`
	Openai         int    `json:"openai"`
	OpenaiToken    int64  `json:"openaiToken"`
	Claude         int    `json:"claude"`
	ClaudeToken    int64  `json:"claudeToken"`
}

type UpdateUserRequest struct {
	ID             int64  `json:"id"`
	UniqueName     string `json:"uniqueName"`
	Password       string `json:"password"`
	ExpirationTime string `json:"expirationTime"`
	Enable         int    `json:"enable"`
	Openai         int    `json:"openai"`
	OpenaiToken    int64  `json:"openaiToken"`
	Claude         int    `json:"claude"`
	ClaudeToken    int64  `json:"claudeToken"`
}

type SearchUserRequest struct {
	UniqueName string `json:"uniqueName"`
}
