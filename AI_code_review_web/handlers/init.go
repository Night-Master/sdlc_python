package handlers

import (
	"ai_code_review/utils"
	"database/sql"
	"sync"
)

// 初始化审计相关变量
func Initialize(cfg *utils.Config, database *sql.DB) {
	// 设置全局变量
	auditMutex = sync.Mutex{}
	progressMap = make(map[string]int)
	totalFiles = make(map[string]int)
	statusMap = make(map[string]string)
	config = cfg
	db = database
}
