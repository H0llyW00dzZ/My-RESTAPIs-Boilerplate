// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package htmxheader

// HTMX Request Headers
const (
	// HXBoosted represents the MIME type for the HX-Boosted header.
	//
	// The HX-Boosted header indicates that the request is via an element using hx-boost.
	HXBoosted = "HX-Boosted"

	// HXCurrentURL represents the MIME type for the HX-Current-URL header.
	//
	// The HX-Current-URL header is used by HTMX to indicate the current URL of the browser.
	// It provides the full URL of the current page, which can be useful for server-side processing.
	HXCurrentURL = "HX-Current-URL"

	// HXHistoryRestoreRequest represents the MIME type for the HX-History-Restore-Request header.
	//
	// The HX-History-Restore-Request header is set to "true" if the request is for history restoration
	// after a miss in the local history cache.
	HXHistoryRestoreRequest = "HX-History-Restore-Request"

	// HXPrompt represents the MIME type for the HX-Prompt header.
	//
	// The HX-Prompt header is used by HTMX to specify the user response to an hx-prompt.
	HXPrompt = "HX-Prompt"

	// HXRequest represents the MIME type for the HX-Request header.
	//
	// The HX-Request header is always set to "true" for HTMX requests.
	HXRequest = "HX-Request"

	// HXTarget represents the MIME type for the HX-Target header.
	//
	// The HX-Target header is used by HTMX to specify the id of the target element if it exists.
	HXTarget = "HX-Target"

	// HXTriggerName represents the MIME type for the HX-Trigger-Name header.
	//
	// The HX-Trigger-Name header is used by HTMX to specify the name of the triggered element if it exists.
	HXTriggerName = "HX-Trigger-Name"

	// HXTrigger represents the MIME type for the HX-Trigger header.
	//
	// The HX-Trigger header is used by HTMX to specify the id of the triggered element if it exists.
	HXTrigger = "HX-Trigger"
)

// HTMX Response Headers
const (
	// HXLocation represents the MIME type for the HX-Location response header.
	//
	// The HX-Location header allows a client-side redirect that does not do a full page reload.
	HXLocation = "HX-Location"

	// HXPushURL represents the MIME type for the HX-Push-Url response header.
	//
	// The HX-Push-Url header pushes a new URL into the history stack.
	HXPushURL = "HX-Push-Url"

	// HXRedirect represents the MIME type for the HX-Redirect response header.
	//
	// The HX-Redirect header can be used to do a client-side redirect to a new location.
	HXRedirect = "HX-Redirect"

	// HXRefresh represents the MIME type for the HX-Refresh response header.
	//
	// If the HX-Refresh header is set to "true", the client-side will do a full refresh of the page.
	HXRefresh = "HX-Refresh"

	// HXReplaceURL represents the MIME type for the HX-Replace-Url response header.
	//
	// The HX-Replace-Url header replaces the current URL in the location bar.
	HXReplaceURL = "HX-Replace-Url"

	// HXReswap represents the MIME type for the HX-Reswap response header.
	//
	// The HX-Reswap header allows specifying how the response will be swapped.
	// See hx-swap for possible values.
	HXReswap = "HX-Reswap"

	// HXRetarget represents the MIME type for the HX-Retarget response header.
	//
	// The HX-Retarget header is a CSS selector that updates the target of the content update
	// to a different element on the page.
	HXRetarget = "HX-Retarget"

	// HXReselect represents the MIME type for the HX-Reselect response header.
	//
	// The HX-Reselect header is a CSS selector that allows choosing which part of the response
	// is used to be swapped in. It overrides an existing hx-select on the triggering element.
	HXReselect = "HX-Reselect"

	// HXTriggerResponse represents the MIME type for the HX-Trigger response header.
	//
	// The HX-Trigger header allows triggering client-side events.
	HXTriggerResponse = "HX-Trigger"

	// HXTriggerAfterSettle represents the MIME type for the HX-Trigger-After-Settle response header.
	//
	// The HX-Trigger-After-Settle header allows triggering client-side events after the settle step.
	HXTriggerAfterSettle = "HX-Trigger-After-Settle"

	// HXTriggerAfterSwap represents the MIME type for the HX-Trigger-After-Swap response header.
	//
	// The HX-Trigger-After-Swap header allows triggering client-side events after the swap step.
	HXTriggerAfterSwap = "HX-Trigger-After-Swap"
)
