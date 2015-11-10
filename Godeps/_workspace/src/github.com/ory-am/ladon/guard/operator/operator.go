package operator

import "time"

const SubjectKey = "subject"

type Context struct {
	Owner     string    `json:"owner"`
	ClientIP  string    `json:"clientIP"`
	Timestamp time.Time `json:"timestamp"`
	UserAgent string    `json:"userAgent"`
}

type Operator func(extra map[string]interface{}, ctx *Context) bool
