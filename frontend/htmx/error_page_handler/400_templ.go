// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.865
// Copyright (c) 2024 H0llyW00dz All rights reserved.

//

// By accessing or using this software, you agree to be bound by the terms

// of the License Agreement, which you can find at LICENSE files.

package htmx

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func PageBadRequest400(v viewData) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"container min-h-screen flex items-center justify-center dark:bg-gray-900\"><div class=\"bg-white dark:bg-gray-800 p-4 sm:p-8 md:p-12 lg:p-20 rounded-lg shadow-lg flex flex-col items-center text-center\"><div class=\"animate__animated animate__fadeIn mb-8 flex flex-col items-center\"><svg class=\"bad-request icon-large w-20 h-20 sm:w-30 sm:h-30 md:w-40 md:h-40 lg:w-50 lg:h-50\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 512 512\"><g><path style=\"opacity:1\" fill=\"#304056\" d=\"M 495.5,62.5 C 495.5,78.8333 495.5,95.1667 495.5,111.5C 400.5,111.5 305.5,111.5 210.5,111.5C 203.046,106.553 195.713,101.386 188.5,96C 153.5,95.3333 118.5,95.3333 83.5,96C 82.6667,96.8333 81.8333,97.6667 81,98.5C 80.6667,102.167 80.3333,105.833 80,109.5C 79.2917,110.381 78.4584,111.047 77.5,111.5C 56.8333,111.5 36.1667,111.5 15.5,111.5C 15.5,95.1667 15.5,78.8333 15.5,62.5C 18.1929,37.3003 31.8596,21.467 56.5,15C 189.167,14.3333 321.833,14.3333 454.5,15C 479.14,21.467 492.807,37.3003 495.5,62.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#fefefe\" d=\"M 290.5,46.5 C 341.834,46.3333 393.168,46.5 444.5,47C 455.174,52.2507 457.174,59.7507 450.5,69.5C 448.869,71.2974 446.869,72.4641 444.5,73C 392.833,73.6667 341.167,73.6667 289.5,73C 277.29,63.6494 277.624,54.8161 290.5,46.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#fafafb\" d=\"M 67.5,46.5 C 80.3609,45.8668 86.8609,51.8668 87,64.5C 84.3027,73.7674 78.1361,77.934 68.5,77C 57.7033,72.241 54.5366,64.4077 59,53.5C 61.4412,50.5464 64.2745,48.2131 67.5,46.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#fafafb\" d=\"M 120.5,46.5 C 135.661,46.4825 141.828,53.8158 139,68.5C 132.814,77.5651 124.98,79.3985 115.5,74C 106.291,62.8324 107.958,53.6657 120.5,46.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#f9fafa\" d=\"M 173.5,46.5 C 182.755,45.3757 189.255,49.0424 193,57.5C 194.037,70.7973 187.87,77.2973 174.5,77C 163.695,72.2263 160.528,64.393 165,53.5C 167.441,50.5464 170.275,48.2131 173.5,46.5 Z\"></path></g> <g><path style=\"opacity:0.995\" fill=\"#dadada\" d=\"M 210.5,111.5 C 305.665,112.499 400.999,112.833 496.5,112.5C 496.667,224.5 496.5,336.5 496,448.5C 492.845,473.656 479.011,489.489 454.5,496C 321.833,496.667 189.167,496.667 56.5,496C 31.979,489.48 18.1457,473.647 15,448.5C 14.5,336.5 14.3333,224.5 14.5,112.5C 35.6733,112.831 56.6733,112.497 77.5,111.5C 78.4584,111.047 79.2917,110.381 80,109.5C 80.3333,105.833 80.6667,102.167 81,98.5C 81.8333,97.6667 82.6667,96.8333 83.5,96C 118.5,95.3333 153.5,95.3333 188.5,96C 195.713,101.386 203.046,106.553 210.5,111.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#8b939c\" d=\"M 15.5,62.5 C 15.5,78.8333 15.5,95.1667 15.5,111.5C 36.1667,111.5 56.8333,111.5 77.5,111.5C 56.6733,112.497 35.6733,112.831 14.5,112.5C 14.1702,95.6583 14.5035,78.9916 15.5,62.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#89919b\" d=\"M 495.5,62.5 C 496.497,78.9916 496.83,95.6583 496.5,112.5C 400.999,112.833 305.665,112.499 210.5,111.5C 305.5,111.5 400.5,111.5 495.5,111.5C 495.5,95.1667 495.5,78.8333 495.5,62.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#010101\" d=\"M 43.5,138.5 C 64.5026,138.333 85.5026,138.5 106.5,139C 109.167,141 109.167,143 106.5,145C 85.5,145.667 64.5,145.667 43.5,145C 40.9387,142.844 40.9387,140.677 43.5,138.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#010101\" d=\"M 121.5,138.5 C 133.171,138.333 144.838,138.5 156.5,139C 159.167,141 159.167,143 156.5,145C 144.833,145.667 133.167,145.667 121.5,145C 118.939,142.844 118.939,140.677 121.5,138.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#498ff6\" d=\"M 247.5,161.5 C 261.523,159.109 272.357,163.775 280,175.5C 322.333,254.167 364.667,332.833 407,411.5C 410.393,430.409 402.893,442.242 384.5,447C 297.833,447.667 211.167,447.667 124.5,447C 107.497,441.173 100.664,429.339 104,411.5C 146.333,332.833 188.667,254.167 231,175.5C 235.198,169.123 240.698,164.456 247.5,161.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#121212\" d=\"M 44.5,163.5 C 51.8409,163.334 59.1742,163.5 66.5,164C 69.859,166.034 70.1924,168.368 67.5,171C 59.5,171.667 51.5,171.667 43.5,171C 40.9327,168.178 41.266,165.678 44.5,163.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#101010\" d=\"M 83.5,163.5 C 99.8367,163.333 116.17,163.5 132.5,164C 135.167,166.333 135.167,168.667 132.5,171C 116.167,171.667 99.8333,171.667 83.5,171C 81.011,168.566 81.011,166.066 83.5,163.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#101010\" d=\"M 43.5,189.5 C 62.8362,189.333 82.1695,189.5 101.5,190C 104.167,192.333 104.167,194.667 101.5,197C 82.1667,197.667 62.8333,197.667 43.5,197C 40.9142,194.506 40.9142,192.006 43.5,189.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#fdfefe\" d=\"M 252.5,200.5 C 255.97,199.875 258.804,200.875 261,203.5C 298.438,272.711 335.605,342.044 372.5,411.5C 372.24,413.187 371.573,414.687 370.5,416C 332.25,417.161 293.916,417.661 255.5,417.5C 217.833,417.333 180.167,417.167 142.5,417C 140.079,416.002 138.745,414.168 138.5,411.5C 175.733,340.698 213.733,270.365 252.5,200.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#060606\" d=\"M 251.5,255.5 C 257.094,254.743 261.261,256.743 264,261.5C 264.683,289.554 263.683,317.554 261,345.5C 257.74,350.285 254.073,350.619 250,346.5C 247.502,318.547 246.169,290.547 246,262.5C 247.371,259.734 249.204,257.401 251.5,255.5 Z\"></path></g> <g><path style=\"opacity:1\" fill=\"#090909\" d=\"M 252.5,363.5 C 261.982,363.48 266.482,368.147 266,377.5C 263.386,384.157 258.553,386.657 251.5,385C 244.304,380.672 242.804,374.838 247,367.5C 248.812,366.023 250.645,364.69 252.5,363.5 Z\"></path></g></svg> <span class=\"text-4xl sm:text-6xl md:text-8xl font-mono text-gray-800 dark:text-white mt-4\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(v.httpStatus)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `frontend/htmx/error_page_handler/400.templ`, Line: 36, Col: 111}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</span></div><h1 class=\"animate__animated animate__fadeIn text-2xl sm:text-3xl md:text-4xl font-mono text-gray-800 dark:text-white mb-4 mt-8\">Bad Request</h1><p class=\"animate__animated animate__fadeIn text-gray-600 dark:text-gray-400 text-base sm:text-lg Roboto mb-8\">Oops! Something went wrong. There was a problem with your request. Please try again later.  If the problem persists, contact administrator.</p><button type=\"button\" class=\"animate__animated animate__fadeIn line-block py-2 sm:py-3 px-4 sm:px-6 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-semibold dark:bg-blue-600 dark:hover:bg-blue-700\" hx-get=\"/\" hx-swap=\"outerHTML\" hx-indicator=\"#spinner\">Go back to the homepage</button><div id=\"spinner\" class=\"htmx-indicator fixed inset-0 z-50 flex items-center justify-center bg-gray-700 bg-opacity-50 hidden\"><div class=\"bg-white p-4 sm:p-8 rounded-lg shadow-lg flex flex-col items-center space-y-4\"><svg aria-hidden=\"true\" class=\"w-8 h-8 sm:w-12 sm:h-12 text-blue-600 animate-spin\" viewBox=\"0 0 100 101\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z\" fill=\"currentColor\"></path> <path d=\"M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z\" fill=\"currentColor\"></path></svg> <span class=\"htmx-indicator-text text-gray-800\">Loading...</span></div></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = Base(v).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
