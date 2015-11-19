package operator

func SubjectIsOwner(extra map[string]interface{}, ctx *Context) bool {
	if _, ok := extra[SubjectKey]; !ok {
		return false
	} else if ctx == nil {
		return false
	}
	return ctx.Owner == extra[SubjectKey]
}

func SubjectIsNotOwner(extra map[string]interface{}, ctx *Context) bool {
	if _, ok := extra[SubjectKey]; !ok {
		return false
	} else if ctx == nil {
		return false
	}
	return ctx.Owner != extra[SubjectKey]
}
