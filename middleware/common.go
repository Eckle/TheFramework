package middleware

import "net/http"

type Key string

type Middleware func(next http.Handler, arguments *map[string]interface{}) http.Handler

func CreateStack(arguments *map[string]interface{}, middlewares ...Middleware) Middleware {
	return func(next http.Handler, arguments *map[string]interface{}) http.Handler {
		for i := len(middlewares); i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(next, arguments)
		}
		return next
	}
}
