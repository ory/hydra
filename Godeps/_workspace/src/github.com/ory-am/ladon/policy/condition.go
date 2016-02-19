package policy

type Condition interface {
	GetOperator() string
	GetExtra() map[string]interface{}
}

type DefaultCondition struct {
	Operator string                 `json:"op"`
	Extra    map[string]interface{} `json:"data"`
}

func (c *DefaultCondition) GetOperator() string {
	return c.Operator
}

func (c *DefaultCondition) GetExtra() map[string]interface{} {
	if c.Extra == nil {
		c.Extra = make(map[string]interface{})
	}
	return c.Extra
}
