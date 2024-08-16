// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package htmxheader provides constants for HTMX headers to ensure consistency and avoid rewriting.
//
// HTMX is a library that allows access to modern browser features directly from HTML, without the need for JavaScript.
// It uses a set of custom headers to enable communication between the client and server.
//
// This package defines constants for both request and response headers used by HTMX.
//
// Request Headers (src https://htmx.org/docs/#request-headers):
//   - HXBoosted: Indicates that the request is via an element using hx-boost.
//   - HXCurrentURL: The current URL of the browser.
//   - HXHistoryRestoreRequest: Set to "true" if the request is for history restoration after a miss in the local history cache.
//   - HXPrompt: The user response to an hx-prompt.
//   - HXRequest: Always set to "true" for HTMX requests.
//   - HXTarget: The id of the target element if it exists.
//   - HXTriggerName: The name of the triggered element if it exists.
//   - HXTrigger: The id of the triggered element if it exists.
//
// Response Headers (src https://htmx.org/docs/#response-headers):
//   - HXLocation: Allows a client-side redirect that does not do a full page reload.
//   - HXPushURL: Pushes a new URL into the history stack.
//   - HXRedirect: Can be used to do a client-side redirect to a new location.
//   - HXRefresh: If set to "true", the client-side will do a full refresh of the page.
//   - HXReplaceURL: Replaces the current URL in the location bar.
//   - HXReswap: Allows specifying how the response will be swapped. See hx-swap for possible values.
//   - HXRetarget: A CSS selector that updates the target of the content update to a different element on the page.
//   - HXReselect: A CSS selector that allows choosing which part of the response is used to be swapped in. Overrides an existing hx-select on the triggering element.
//   - HXTriggerResponse: Allows triggering client-side events.
//   - HXTriggerAfterSettle: Allows triggering client-side events after the settle step.
//   - HXTriggerAfterSwap: Allows triggering client-side events after the swap step.
//
// By using these constants, consistency can be ensured in the codebase and the need to rewrite the header names multiple times can be avoided.
package htmxheader
