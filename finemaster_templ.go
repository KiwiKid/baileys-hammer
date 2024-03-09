// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
)

var finemasterBaseUrl = "/finemaster"

func finemaster(pass string, players []PlayerWithFines, pFines []PresetFine, qp FineMasterQueryParams) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<html>")
		if err != nil {
			return err
		}
		err = header().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body><div class=\"bg-yellow-400 border-l-4 border-yellow-800 text-yellow-800 p-2 text-center\" role=\"alert\"><p class=\"font-bold\">")
		if err != nil {
			return err
		}
		var_2 := `🚧 Under Construction (Version 1) 🚧`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p></div><div class=\"bg-gray-900 text-center p-5\"><h1 class=\"text-xl md:text-3xl lg:text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-teal-400\">")
		if err != nil {
			return err
		}
		var_3 := `Fine Master Zone`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1></div><div><a hx-transition=\"true\" href=\"/\">")
		if err != nil {
			return err
		}
		var_4 := `Reset`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a></div><div class=\"container mx-auto p-4\"><h1 class=\"text-2xl font-bold mb-4\">")
		if err != nil {
			return err
		}
		var_5 := `Player Fines`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><div class=\"flex items-center justify-center bg-gray-100 mx-auto\"><ul>")
		if err != nil {
			return err
		}
		for _, p := range players {
			_, err = templBuffer.WriteString("<li class=\"mb-2\"><div _=\"on click toggle .hidden on next &lt;section/&gt;\" class=\"cursor-pointer p-2 bg-gray-200 rounded hover:bg-gray-300\">")
			if err != nil {
				return err
			}
			var var_6 string = p.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_6))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_7 := `- `
			_, err = templBuffer.WriteString(var_7)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_8 string = fmt.Sprintf("$%d (%d)", p.TotalFines, p.Number)
			_, err = templBuffer.WriteString(templ.EscapeString(var_8))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div><section class=\"fines-info hidden\"><div class=\"p-2\"><div class=\"p-2\">")
			if err != nil {
				return err
			}
			for _, f := range p.Fines {
				_, err = templBuffer.WriteString("<div class=\"mt-1\"><p>")
				if err != nil {
					return err
				}
				var var_9 string = f.Reason
				_, err = templBuffer.WriteString(templ.EscapeString(var_9))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var_10 := `- `
				_, err = templBuffer.WriteString(var_10)
				if err != nil {
					return err
				}
				var var_11 string = fmt.Sprintf("$%.0f", f.Amount)
				_, err = templBuffer.WriteString(templ.EscapeString(var_11))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</p></div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div></div></section></li>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</ul></div></div><div class=\"grid grid-cols-1 gap-2\">")
		if err != nil {
			return err
		}
		err = fineAdd(fmt.Sprintf("%s/%s", finemasterBaseUrl, pass), qp.FinesOpen, players, pFines).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		err = playersAdd(fmt.Sprintf("%s/%s", finemasterBaseUrl, pass), qp.PlayerOpen).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		err = presetFines(fmt.Sprintf("%s/%s", finemasterBaseUrl, pass), qp.PresetFinesOpen, pFines).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><script src=\"https://unpkg.com/htmx.org\">")
		if err != nil {
			return err
		}
		var_12 := ``
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func presetFines(baseUrl string, isOpen bool, presetFines []PresetFine) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_13 := templ.GetChildren(ctx)
		if var_13 == nil {
			var_13 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if isOpen {
			_, err = templBuffer.WriteString("<div class=\"px-8 py-6 mt-4 text-left bg-white shadow-lg\"><h3 class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			var_14 := `Preset Fines`
			_, err = templBuffer.WriteString(var_14)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3><div class=\"grid grid-cols-1 md:grid-cols-2 gap-4\"><div><form action=\"/preset-fines\" method=\"POST\" class=\"mt-4\"><div><label for=\"reason\" class=\"block\">")
			if err != nil {
				return err
			}
			var_15 := `Reason`
			_, err = templBuffer.WriteString(var_15)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input required type=\"text\" name=\"reason\" id=\"reason\" placeholder=\"Reason for the fine\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></div><div class=\"mt-4\"><label for=\"amount\" class=\"block\">")
			if err != nil {
				return err
			}
			var_16 := `Amount ($)`
			_, err = templBuffer.WriteString(var_16)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input required type=\"number\" step=\"0.01\" name=\"amount\" id=\"amount\" placeholder=\"Amount\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></div><div class=\"flex items-center justify-between mt-4\">")
			if err != nil {
				return err
			}
			var var_17 = []any{add}
			err = templ.RenderCSSItems(ctx, templBuffer, var_17...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_17).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_18 := `Add Preset Fine`
			_, err = templBuffer.WriteString(var_18)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div></form><div>")
			if err != nil {
				return err
			}
			var var_19 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_19...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_20 templ.SafeURL = makeSafeUrl(baseUrl, false, false, false)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_20)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_19).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_21 := `Close`
			_, err = templBuffer.WriteString(var_21)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></div><div><h1>")
			if err != nil {
				return err
			}
			var_22 := `Existing Preset Fines`
			_, err = templBuffer.WriteString(var_22)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1>")
			if err != nil {
				return err
			}
			for _, fine := range presetFines {
				_, err = templBuffer.WriteString("<div class=\"mt-2\"><p>")
				if err != nil {
					return err
				}
				var var_23 string = fine.Reason
				_, err = templBuffer.WriteString(templ.EscapeString(var_23))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var_24 := `- `
				_, err = templBuffer.WriteString(var_24)
				if err != nil {
					return err
				}
				var var_25 string = fmt.Sprintf("$%f", fine.Amount)
				_, err = templBuffer.WriteString(templ.EscapeString(var_25))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</p><div hx-delete=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/preset-fines?pfid=%d", fine.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var_26 := `🗑`
				_, err = templBuffer.WriteString(var_26)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div></div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div></div></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"flex justify-center w-full\">")
			if err != nil {
				return err
			}
			var var_27 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_27...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_28 templ.SafeURL = makeSafeUrl(baseUrl, false, false, true)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_28)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_27).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_29 := `View Preset Fines`
			_, err = templBuffer.WriteString(var_29)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div>")
			if err != nil {
				return err
			}
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func playersAdd(baseUrl string, isOpen bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_30 := templ.GetChildren(ctx)
		if var_30 == nil {
			var_30 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if isOpen {
			_, err = templBuffer.WriteString("<div class=\"flex items-center justify-center bg-gray-100\"><div class=\"px-8 py-6 mt-4 text-left bg-white shadow-lg\"><h3 class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			var_31 := `Add New Player`
			_, err = templBuffer.WriteString(var_31)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3><form hx-post=\"/players\" method=\"POST\" class=\"mt-4\"><div><label for=\"name\" class=\"block\">")
			if err != nil {
				return err
			}
			var_32 := `Name`
			_, err = templBuffer.WriteString(var_32)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"name\" id=\"name\" placeholder=\"Name\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></div><div class=\"mt-4\"><label for=\"number\" class=\"block\">")
			if err != nil {
				return err
			}
			var_33 := `Number`
			_, err = templBuffer.WriteString(var_33)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"number\" name=\"number\" id=\"number\" placeholder=\"Number\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></div><div class=\"flex items-center w-full\">")
			if err != nil {
				return err
			}
			var var_34 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_34...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_34).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_35 := `Add Player`
			_, err = templBuffer.WriteString(var_35)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div><div class=\"flex justify-center w-full\">")
			if err != nil {
				return err
			}
			var var_36 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_36...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_37 templ.SafeURL = makeSafeUrl(baseUrl, false, false, false)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_37)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-transition=\"true\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_36).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_38 := `Close`
			_, err = templBuffer.WriteString(var_38)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></form></div></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"flex justify-center w-full\">")
			if err != nil {
				return err
			}
			var var_39 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_39...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_40 templ.SafeURL = makeSafeUrl(baseUrl, false, true, false)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_40)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-transition=\"true\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_39).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_41 := `Add/Remove Players`
			_, err = templBuffer.WriteString(var_41)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div>")
			if err != nil {
				return err
			}
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
