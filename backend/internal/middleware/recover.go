package middleware

import (
	"fmt"
	"lytemp/pkg/lg"
	"runtime"
	"strings"

	"github.com/Lexographics/logar"
	"github.com/labstack/echo/v4"
)

func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					pc := make([]uintptr, 10)
					n := runtime.Callers(2, pc) // skip Recover() function and runtime.Callers

					if n == 0 {
						return
					}

					pc = pc[:n]
					frames := runtime.CallersFrames(pc)

					stack := []string{}

					for {
						frame, more := frames.Next()
						function := frame.Function

						// skip panic() or log.Panic()
						if strings.Contains(function, "runtime.gopanic") || strings.Contains(frame.File, "src/log/log.go") {
							continue
						}

						stack = append(stack, fmt.Sprintf("%s:%d", frame.File, frame.Line))

						if !more {
							break
						}
					}

					lg.WithContext(c.Request().Context()).Error(lg.UserLogs, logar.Map{
						"error": "panic recover: " + err.Error(),
						"stack": stack,
					}, "request")
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
