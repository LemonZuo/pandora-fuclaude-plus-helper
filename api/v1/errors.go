package v1

var (
	// ErrSuccess common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "参数错误")
	ErrUnauthorized        = newError(401, "未授权")
	ErrForbidden           = newError(403, "禁止访问")
	ErrNotFound            = newError(404, "不存在")
	ErrMethodNotAllowed    = newError(405, "不允许的方法")
	ErrInternalServerError = newError(500, "内部服务器错误")

	// ErrEmailAlreadyUse more biz errors
	ErrEmailAlreadyUse   = newError(1001, "该电子邮件已被使用。")
	ErrCannotRefresh     = newError(1002, "无法刷新帐户。")
	ErrAccessTokenEmpty  = newError(1003, "访问令牌为空。")
	ErrLoginFailed       = newError(1004, "登录失败")
	ErrCannotDeleteToken = newError(1005, "已有关联账户，请先删除关联账户。")
)
