package response

const (
	ErrCodeSuccess      = 20001
	ErrCodeParamInvalid = 20003

	ErrInvalidToken     = 30001
	ErrInvalidOTP       = 30002
	ErrSendEmailOtp     = 30003
	ErrCodeAuthFailed   = 40005
	ErrCodeUserHasExist = 50001
	// Err login
	ErrCodeOtpNotExist      = 60009
	ErrCodeUserOtpNotExists = 60008

	ErrCodeTwoFactorSetupFailed  = 80001
	ErrCodeOtpMismatch           = 60010
	ErrCodeOtpDeleteFailed       = 60011
	ErrCodeTwoFactorUpdateFailed = 60012
	ErrCodePostFailed            = 70000
	ErrCodeInternal              = 80000
	ErrCodeNotFound              = 80008
	ErrCodeComment               = 90000
	ErrCodeCreateRoom            = 101
	ErrCodeGetMessage            = 102
)

var msg = map[int]string{
	ErrCodeSuccess:               "Success",
	ErrCodeParamInvalid:          "Email is invalid",
	ErrInvalidToken:              "Invalid token",
	ErrCodeUserHasExist:          "User has exist",
	ErrInvalidOTP:                "OTP error",
	ErrSendEmailOtp:              "Failed to send email OTP",
	ErrCodeOtpNotExist:           "OTP exists but not registed",
	ErrCodeUserOtpNotExists:      "ErrCodeUserOtpNotExists",
	ErrCodeAuthFailed:            "ErrCodeAuthFailed",
	ErrCodeTwoFactorSetupFailed:  "ErrCodeTwoFactorSetupFailed",
	ErrCodeOtpMismatch:           "OTP does not match",
	ErrCodeOtpDeleteFailed:       "Failed to delete OTP",
	ErrCodeTwoFactorUpdateFailed: "Failed to update two-factor authentication status",
	ErrCodePostFailed:            "ErrCodePostFailed",
	ErrCodeInternal:              "ErrCodeInternal",
	ErrCodeNotFound:              "ErrCodeNotFound",
	ErrCodeComment:               "ErrCodeCommentFailed",
	ErrCodeCreateRoom:            "ErrCodeCreateRoom",
	ErrCodeGetMessage:            "ErrCodeGetMessage",
}
