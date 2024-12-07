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
	"github.com/gofiber/fiber/v2/utils"
)

// appNameTag is a fiber logger custom tag function that retrieves the name of the application.
//
// Note: Ignore the warning about unused parameters. Once it is bound to [WithLoggerCustomTags], the warning will disappear from the IDE.
func appNameTag(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	if c.App().Config().AppName == "" {
		return output.WriteString("-")
	}
	return output.WriteString(c.App().Config().AppName)
}

// unixTimeTag is a fiber logger custom tag function that returns the current Unix timestamp.
//
// Note: Ignore the warning about unused parameters. Once it is bound to [WithLoggerCustomTags], the warning will disappear from the IDE.
func unixTimeTag(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	return output.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
}

// hostNameTag is a fiber logger custom tag function that returns the current Hostname.
//
// Note: Ignore the warning about unused parameters. Once it is bound to [WithLoggerCustomTags], the warning will disappear from the IDE.
func hostNameTag(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	return output.WriteString(utils.CopyString(c.Hostname()))
}

// userAgentTag is a Fiber logger custom tag function that retrieves the User-Agent string from the incoming HTTP request.
func userAgentTag(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	if userAgent := c.Context().UserAgent(); userAgent != nil {
		return output.Write(userAgent)
	}
	return output.WriteString("-")
}

// proxyTag is a Fiber logger custom tag function that retrieves the remote IP address if the trusted proxy check is enabled.
func proxyTag(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	if c.App().Config().EnableTrustedProxyCheck {
		return output.WriteString(c.Context().RemoteIP().String())
	}
	return output.WriteString("-")
}

// chainPathError logs the current request path along with any error message stored in the context.
//
// This function is intended to be used as a custom logger function in a Fiber application.
// It retrieves an error message from the context's local storage and appends it to the request path
// for logging purposes.
//
// Note: Unlike the "${error}" tag, this function is suitable for custom error handling in Fiber.
func chainPathError(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	// Retrieve the error from the c.Locals
	if err, ok := c.Locals("error").(string); ok && err != "" {
		return output.WriteString(c.Path() + " - " + err)
	}
	return output.WriteString(c.Path())
}
