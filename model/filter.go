package model

import (
	"errors"
	"strconv"
)

func integerParseFilter(c *Context) {
	i, err := strconv.Atoi(c.input)
	if err != nil {
		c.errOutput = "parse index error"
		c.err = err
		return
	}
	c.inputIndex = i
}

func intRangeFilter(c *Context, min int, max int) pipeFunc {
	return func(context *Context) {
		i, err := strconv.Atoi(c.input)
		if err != nil {
			c.errOutput = "parse index error"
			c.err = err
			return
		}
		if i < min || i > max {
			c.errOutput = "index not in range"
			c.err = errors.New(c.errOutput)
			return
		}
		c.inputIndex = i
	}
}
