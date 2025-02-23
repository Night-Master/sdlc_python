package utils

// AuditResult 表示AI审计的结果
type AuditResult struct {
	Result []Vulnerability `json:"result"`
}
