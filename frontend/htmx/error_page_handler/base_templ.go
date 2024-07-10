// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
// Copyright (c) 2024 H0llyW00dz All rights reserved.

//

// License: BSD 3-Clause License

package htmx

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Base(titlePage, cfheader string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
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
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"utf-8\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(titlePage)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `frontend/htmx/error_page/base.templ`, Line: 12, Col: 21}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><style>\r\n\t\t\tbody {\r\n\t\t\t\tbackground-color: #f5f5f5;\r\n\t\t\t\tmargin-top: 8%;\r\n\t\t\t\tcolor: #5d5d5d;\r\n\t\t\t\tfont-family: -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, \"Helvetica Neue\", Arial,\r\n\t\t\t\t\t\"Noto Sans\", sans-serif, \"Apple Color Emoji\", \"Segoe UI Emoji\", \"Segoe UI Symbol\",\r\n\t\t\t\t\t\"Noto Color Emoji\";\r\n\t\t\t\ttext-shadow: 0px 1px 1px rgba(255, 255, 255, 0.75);\r\n\t\t\t\ttext-align: center;\r\n\t\t\t}\r\n\r\n\t\t\th1 {\r\n\t\t\t\tfont-size: 2.45em;\r\n\t\t\t\tfont-weight: 700;\r\n\t\t\t\tcolor: #5d5d5d;\r\n\t\t\t\tletter-spacing: -0.02em;\r\n\t\t\t\tmargin-bottom: 30px;\r\n\t\t\t\tmargin-top: 30px;\r\n\t\t\t}\r\n\r\n\t\t\t.container {\r\n\t\t\t\twidth: 100%;\r\n\t\t\t\tmargin-right: auto;\r\n\t\t\t\tmargin-left: auto;\r\n\t\t\t}\r\n\r\n\t\t\t.animate__animated {\r\n\t\t\t\tanimation-duration: 1s;\r\n\t\t\t\tanimation-fill-mode: both;\r\n\t\t\t}\r\n\r\n\t\t\t.animate__fadeIn {\r\n\t\t\t\tanimation-name: fadeIn;\r\n\t\t\t}\r\n\r\n\t\t\t.info {\r\n\t\t\t\tcolor: #5594cf;\r\n\t\t\t\tfill: #5594cf;\r\n\t\t\t}\r\n\r\n\t\t\t.error {\r\n\t\t\t\tcolor: #c92127;\r\n\t\t\t\tfill: #c92127;\r\n\t\t\t}\r\n\r\n\t\t\t.warning {\r\n\t\t\t\tcolor: #ffcc33;\r\n\t\t\t\tfill: #ffcc33;\r\n\t\t\t}\r\n\r\n\t\t\t.success {\r\n\t\t\t\tcolor: #5aba47;\r\n\t\t\t\tfill: #5aba47;\r\n\t\t\t}\r\n\r\n\t\t\t.icon-large {\r\n\t\t\t\theight: 132px;\r\n\t\t\t\twidth: 132px;\r\n\t\t\t}\r\n\r\n\t\t\t.description-text {\r\n\t\t\t\tcolor: #707070;\r\n\t\t\t\tletter-spacing: -0.01em;\r\n\t\t\t\tfont-size: 1.25em;\r\n\t\t\t\tline-height: 20px;\r\n\t\t\t}\r\n\r\n\t\t\t.footer {\r\n\t\t\t\tmargin-top: 40px;\r\n\t\t\t\tfont-size: 0.7em;\r\n\t\t\t}\r\n\r\n\t\t\t.animate__delay-1s {\r\n\t\t\t\tanimation-delay: 1s;\r\n\t\t\t}\r\n\r\n\t\t\t@keyframes fadeIn {\r\n\t\t\t\tfrom {\r\n\t\t\t\t\topacity: 0;\r\n\t\t\t\t}\r\n\t\t\t\tto {\r\n\t\t\t\t\topacity: 1;\r\n\t\t\t\t}\r\n\t\t\t}\r\n\t\t</style><script src=\"/styles/js/htmx.min.js\"></script></head><body><main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if cfheader != "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"description-text animate__animated animate__fadeIn animate__delay-1s\"><section class=\"footer\"><strong>Ray ID:</strong> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(cfheader)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `frontend/htmx/error_page/base.templ`, Line: 107, Col: 64}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</section></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}