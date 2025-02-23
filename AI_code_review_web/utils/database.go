package utils

import (
	"database/sql"

	_ "modernc.org/sqlite" // 使用纯Go的SQLite驱动
)

// Task 表示一个审计任务
type Task struct {
	TaskID       string         `json:"taskID"`
	Status       string         `json:"status"`
	StartTime    string         `json:"startTime"`
	EndTime      sql.NullString `json:"endTime,omitempty"`
	Language     string         `json:"language"`
	ThreadNum    int            `json:"threadNum"`
	AIPlatform   string         `json:"aiPlatform"`
	TargetPath   string         `json:"targetPath"`             // 目录或压缩包名
	WebFramework string         `json:"webFramework,omitempty"` // 新增的 Web 框架字段
}

// Vulnerability 表示一个漏洞
type Vulnerability struct {
	FuncName    string `json:"func_name"`
	Line        int    `json:"line"`
	Message     string `json:"message"`
	Severity    string `json:"severity"`
	CWE         string `json:"cwe"`
	FilePath    string `json:"file_path"`
	CodeSnippet string `json:"code_snippet"` // 新增的代码片段
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./report.db") // 使用 "sqlite" 驱动名称
	if err != nil {
		return nil, err
	}

	// 创建tasks表，新增字段 WebFramework
	createTasksTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
        task_id TEXT PRIMARY KEY,
        status TEXT,
        start_time DATETIME,
        end_time DATETIME,
        language TEXT,
        thread_num INTEGER,
        ai_platform_name TEXT,
        target_path TEXT,
        web_framework TEXT -- 新增字段
    );`
	_, err = db.Exec(createTasksTableSQL)
	if err != nil {
		return nil, err
	}

	// 确保web_framework字段存在（SQLite不支持IF NOT EXISTS语法）
	// 检查字段是否存在
	var count int
	err = db.QueryRow("PRAGMA table_info(tasks)").Scan(&count)
	if err != nil {
		// 处理错误
	}
	// 这里可以添加逻辑检查字段是否存在，并进行相应的迁移
	// 由于SQLite不支持直接检查列是否存在，我们可能需要采用其他方式

	return db, nil
}
