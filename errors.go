package qgnet

import "fmt"

type QgError struct {
	Code    string `json:"code"`    // 错误码
	Message string `json:"message"` // 描述
}

func (e *QgError) Error() string {
	return fmt.Sprintf("QgError: [%s] %s", e.Code, e.Message)
}

func ErrorOf(code string, message string) *QgError {
	return &QgError{
		Code:    code,
		Message: message,
	}
}

type ResultRO struct {
	Code      string `json:"code"`       // 错误码
	Message   string `json:"message"`    // 描述
	RequestId string `json:"request_id"` // 请求ID
}

var (
	ErrCodeStatus = "STATUS_ERROR"           // 状态错误, 为没有错误码的CODE预留
	ErrCodeUnknow = "UNKNOW_ERROR"           // 未知错误码
	ErrSuccess    = ErrorOf("SUCCESS", "成功") // 请求成功

	ErrInternalError             = ErrorOf("INTERNAL_ERROR", "系统内部异常")
	ErrInvalidParameter          = ErrorOf("INVALID_PARAMETER", "参数错误")
	ErrInvalidKey                = ErrorOf("INVALID_KEY", "Key不存在或已过期")
	ErrUnavailableKey            = ErrorOf("UNAVAILABLE_KEY", "Key不可用")
	ErrAccessDeny                = ErrorOf("ACCESS_DENY", "Key没有此接口的权限")
	ErrApiAuthDeny               = ErrorOf("API_AUTH_DENY", "Api授权不通过")
	ErrKeyBlock                  = ErrorOf("KEY_BLOCK", "Key被封禁")
	ErrRequestLimitExceeded      = ErrorOf("REQUEST_LIMIT_EXCEEDED", "请求频率超出限制")
	ErrBalanceInsufficient       = ErrorOf("BALANCE_INSUFFICIENT", "Key余额不足")
	ErrNoResourceFound           = ErrorOf("NO_RESOURCE_FOUND", "资源不足")
	ErrFailedOperation           = ErrorOf("FAILED_OPERATION", "提取失败")
	ErrExtractLimitExceeded      = ErrorOf("EXTRACT_LIMIT_EXCEEDED", "超出提取配额")
	ErrDeleteLimitExceeded       = ErrorOf("DELETE_LIMIT_EXCEEDED", "释放频率超出限制")
	ErrIpWhitelistLimitExceeded  = ErrorOf("IP_WHITELIST_LIMIT_EXCEEDED", "白名单数量超出限制")
	ErrStaticDeleteTimeLimit     = ErrorOf("STATIC_DELETE_TIME_LIMIT", "静态资源需要24小时后才能释放")
	ErrMonopolyDeleteTimeLimit   = ErrorOf("MONOPOLY_DELETE_TIME_LIMIT", "独占资源需要12小时后才能释放")
	ErrMonopolyChangeIpTimeLimit = ErrorOf("MONOPOLY_CHANGEIP_TIME_LIMIT", "独占资源切换IP需要等待10秒")
	ErrNoAvailableChannel        = ErrorOf("NO_AVAILABLE_CHANNEL", "没有可用的空闲通道")

	ErrMap = map[string]*QgError{ // 错误码映射
		ErrCodeStatus:                     ErrInternalError,
		ErrCodeUnknow:                     ErrInternalError,
		ErrInternalError.Code:             ErrInternalError,
		ErrInvalidParameter.Code:          ErrInvalidParameter,
		ErrInvalidKey.Code:                ErrInvalidKey,
		ErrUnavailableKey.Code:            ErrUnavailableKey,
		ErrAccessDeny.Code:                ErrAccessDeny,
		ErrApiAuthDeny.Code:               ErrApiAuthDeny,
		ErrKeyBlock.Code:                  ErrKeyBlock,
		ErrRequestLimitExceeded.Code:      ErrRequestLimitExceeded,
		ErrBalanceInsufficient.Code:       ErrBalanceInsufficient,
		ErrNoResourceFound.Code:           ErrNoResourceFound,
		ErrFailedOperation.Code:           ErrFailedOperation,
		ErrExtractLimitExceeded.Code:      ErrExtractLimitExceeded,
		ErrDeleteLimitExceeded.Code:       ErrDeleteLimitExceeded,
		ErrIpWhitelistLimitExceeded.Code:  ErrIpWhitelistLimitExceeded,
		ErrStaticDeleteTimeLimit.Code:     ErrStaticDeleteTimeLimit,
		ErrMonopolyDeleteTimeLimit.Code:   ErrMonopolyDeleteTimeLimit,
		ErrMonopolyChangeIpTimeLimit.Code: ErrMonopolyChangeIpTimeLimit,
		ErrNoAvailableChannel.Code:        ErrNoAvailableChannel,
	}
)

func GetError(code string, message string) *QgError {
	if err, ok := ErrMap[code]; ok {
		// 返回已知错误码
		return err
	}
	if message != "" {
		// 返回未知错误码
		return ErrorOf(code, message)
	}
	// 返回未知错误码
	return ErrorOf(ErrCodeUnknow, "未知错误码: "+code)
}
