package errors

import "fmt"

// ErrorType 定义错误类型枚举
type ErrorType int

const (
	ErrTypeInvalidURL ErrorType = iota // URL 格式无效
	ErrTypeUnsupportedURL               // 不支持的 URL 类型
	ErrTypeNotFound                     // 资源不存在 (404)
	ErrTypeRateLimit                    // API 限流
	ErrTypeAuthRequired                 // 需要认证
	ErrTypeNetwork                      // 网络故障
	ErrTypeFileExists                   // 文件已存在
)

// AppError 包含错误类型和原始错误
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %w", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 返回被包装的原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}
