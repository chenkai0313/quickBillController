package errmes

type ErrorCode string

const (
	ErrInternalServer    ErrorCode = "internal_server_error"
	ErrInvalidRequest    ErrorCode = "invalid_request"
	ErrNotFound          ErrorCode = "not_found"
	ErrForbidden         ErrorCode = "forbidden"
	ErrUnauthorized      ErrorCode = "unauthorized"
	ErrTimeout           ErrorCode = "timeout"
	ErrRateLimitExceeded ErrorCode = "rate_limit_exceeded"
	ErrTokenInvalid      ErrorCode = "token_invalid"

	ErrGenerateNonce ErrorCode = "generate_nonce"
	ErrVerifyLogin   ErrorCode = "verify_login"
	ErrUploadFile    ErrorCode = "upload_file"

	// WebAuthn related error codes
	ErrWebAuthnRequired           ErrorCode = "webauthn_required"
	ErrWebAuthnVerificationFailed ErrorCode = "webauthn_verification_failed"
)
