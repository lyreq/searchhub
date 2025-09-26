package middleware

import (
	"fmt"
	"lytemp/pkg/lg"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Lexographics/logar"
	"github.com/labstack/echo/v4"
	"github.com/mileusna/useragent"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/logar") {
				return next(c)
			}

			requestId := c.Get("requestid")
			userIdAny := c.Get("userid")
			userId := uint(0)
			if userIdAny != nil {
				userId, _ = userIdAny.(uint)
			}

			ctx := c.Request().Context()
			ctx = lg.Get().PrepareContext(ctx, logar.Map{
				"url":        c.Request().Method + " " + c.Request().URL.Path,
				"request_id": requestId,
				"user_id":    userId,
			})
			c.SetRequest(c.Request().WithContext(ctx))

			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			userIdAny = c.Get("userid")
			if userIdAny != nil {
				userId, _ = userIdAny.(uint)
			}

			statusCode := http.StatusOK
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				if ok {
					statusCode = he.Code
				} else {
					statusCode = http.StatusInternalServerError
				}
			} else {
				statusCode = c.Response().Status
			}

			if statusCode > 400 {
				errText := "unknown error"
				if err != nil {
					errText = err.Error()
				}
				lg.WithContext(ctx).Error(lg.UserLogs, logar.Map{
					"error":       errText,
					"status_code": statusCode,
				}, "request")
			}

			contentLength, _ := strconv.ParseInt(c.Request().Header.Get(echo.HeaderContentLength), 10, 64)
			ua := useragent.Parse(c.Request().UserAgent())

			lg.Analytics().RegisterRequest(logar.RequestLog{
				Timestamp:  time.Now(),
				VisitorID:  fmt.Sprint(userId),
				Instance:   "green",
				Path:       c.Request().URL.Path,
				Latency:    duration,
				StatusCode: statusCode,
				UserAgent:  c.Request().UserAgent(),
				OS:         ua.OS,
				Browser:    ua.Name,
				Referer:    c.Request().Referer(),
				BytesSent:  c.Response().Size,
				BytesRecv:  contentLength,
			})

			// lg.Logger().Trace("user-trace", logar.Map{
			// 	"message":    "request successful",
			// 	"request_id": requestId,
			// 	"user_id":    userId,
			// 	"url":        c.Request().Method + " " + c.Request().URL.Path,
			// }, "request")

			return err
		}
	}
}
