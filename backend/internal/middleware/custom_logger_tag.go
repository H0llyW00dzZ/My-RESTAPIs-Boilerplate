// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// TODO: Implement more tag ?

package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// appNameTag is a fiber custom logger tag function that retrieves the name of the application.
//
// Note: Ignore the warning about unused parameters. Once it is bound to [WithLoggerCustomTags], the warning will disappear from the IDE.
func appNameTag(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	appName := c.App().Config().AppName
	return output.WriteString(appName)
}

// unixTimeTag is a fiber custom logger tag function that returns the current Unix timestamp.
//
// Note: Ignore the warning about unused parameters. Once it is bound to [WithLoggerCustomTags], the warning will disappear from the IDE.
func unixTimeTag(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	unixTime := strconv.FormatInt(time.Now().Unix(), 10)
	return output.WriteString(unixTime)
}
