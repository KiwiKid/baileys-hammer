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

func secureFineMasterbaseUrl(finemasterBaseUrl string, pass string) string {
	return fmt.Sprintf("%s/%s", finemasterBaseUrl, pass)
}

func finemasterNav(finemasterBaseUrl string) templ.Component {
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
		_, err = templBuffer.WriteString("<nav class=\"fixed inset-x-0 bottom-0 bg-gray-800 text-white pb-12\"><div class=\"flex justify-between\"><a href=\"")
		if err != nil {
			return err
		}
		var var_2 templ.SafeURL = makeSafeUrlWithAnchor(finemasterBaseUrl, true, false, false, false, "fine-add")
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_2)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"flex-1 text-center py-3 hover:bg-gray-700\">")
		if err != nil {
			return err
		}
		var_3 := `Add`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><button _=\"on click toggle .hidden on .quick-fine then go to top of #quick-finer\" class=\"flex-1 text-center py-3 hover:bg-gray-700\">")
		if err != nil {
			return err
		}
		var_4 := `Quick Fine`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button><a href=\"")
		if err != nil {
			return err
		}
		var var_5 templ.SafeURL = makeSafeUrlWithAnchor(finemasterBaseUrl, false, false, false, true, "preset-fine")
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_5)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"flex-1 text-center py-3 hover:bg-gray-700\">")
		if err != nil {
			return err
		}
		var_6 := `Standard Fines`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><a href=\"")
		if err != nil {
			return err
		}
		var var_7 templ.SafeURL = makeSafeUrlWithAnchor(finemasterBaseUrl, false, false, true, false, "players-manage")
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_7)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"flex-1 text-center py-3 hover:bg-gray-700\">")
		if err != nil {
			return err
		}
		var_8 := `Players`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><a href=\"")
		if err != nil {
			return err
		}
		var var_9 templ.SafeURL = makeSafeUrlWithAnchor(finemasterBaseUrl, false, false, false, false, "fine-list-container")
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_9)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"flex-1 text-center py-3 hover:bg-gray-700\">")
		if err != nil {
			return err
		}
		var_10 := `Recent`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a></div></nav>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func finemaster(pass string, players []PlayerWithFines, pFines []PresetFine, qp FineMasterQueryParams) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_11 := templ.GetChildren(ctx)
		if var_11 == nil {
			var_11 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<html hx-boost=\"true\">")
		if err != nil {
			return err
		}
		err = header().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		err = tomSelectLinks().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body><div class=\"bg-gray-900 text-center p-5\"><h1 class=\"text-xl md:text-3xl lg:text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-teal-400\">")
		if err != nil {
			return err
		}
		var_12 := `Fine Master Zone`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1></div><div class=\"grid grid-cols-1 gap-2\">")
		if err != nil {
			return err
		}
		err = fineAddV2(secureFineMasterbaseUrl(finemasterBaseUrl, pass), qp.FinesOpen, players, pFines, true).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div>")
		if err != nil {
			return err
		}
		err = fineAdd(secureFineMasterbaseUrl(finemasterBaseUrl, pass), qp.FinesOpen, players, pFines, true).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"container mx-auto p-4\" id=\"quick-finer\"><div class=\"flex items-center justify-center bg-gray-100 mx-auto\"><ul>")
		if err != nil {
			return err
		}
		for _, p := range players {
			_, err = templBuffer.WriteString("<li class=\"mb-2\">")
			if err != nil {
				return err
			}
			var var_13 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_13...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<div _=\"on click toggle .hidden on next &lt;section/&gt;\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_13).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var var_14 string = p.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_14))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_15 := `- `
			_, err = templBuffer.WriteString(var_15)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_16 string = fmt.Sprintf("$%d (%d)", p.TotalFines, p.TotalFineCount)
			_, err = templBuffer.WriteString(templ.EscapeString(var_16))
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
				var var_17 string = f.Reason
				_, err = templBuffer.WriteString(templ.EscapeString(var_17))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var_18 := `- `
				_, err = templBuffer.WriteString(var_18)
				if err != nil {
					return err
				}
				var var_19 string = fmt.Sprintf("$%.0f", f.Amount)
				_, err = templBuffer.WriteString(templ.EscapeString(var_19))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</p></div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div></div></section><section class=\"hidden quick-fine\">")
			if err != nil {
				return err
			}
			for _, pf := range pFines {
				if pf.Approved && !pf.NotQuickFine {
					_, err = templBuffer.WriteString("<form class=\" inline-flex mx-2 space-y-2\">")
					if err != nil {
						return err
					}
					var var_20 = []any{fmt.Sprintf("fine-group-%d-%d", pf.ID, p.ID)}
					err = templ.RenderCSSItems(ctx, templBuffer, var_20...)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("<div hidden class=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_20).String()))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\"><input type=\"hidden\" name=\"playerId\" value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", p.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\"><input type=\"hidden\" name=\"presetFineId\" value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", pf.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\"><input type=\"hidden\" name=\"approved\" value=\"on\"><input type=\"hidden\" name=\"dontRedirect\" value=\"true\"></div>")
					if err != nil {
						return err
					}
					var var_21 = []any{bigAdd}
					err = templ.RenderCSSItems(ctx, templBuffer, var_21...)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("<button hx-post=\"/fines\" hx-swap=\"this\" hx-include=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf(".fine-group-%d-%d", pf.ID, p.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" class=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_21).String()))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var var_22 string = fmt.Sprintf("%s ($%v)", pf.Reason, pf.Amount)
					_, err = templBuffer.WriteString(templ.EscapeString(var_22))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</button></form>")
					if err != nil {
						return err
					}
				}
			}
			_, err = templBuffer.WriteString("</section></li>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</ul></div>")
		if err != nil {
			return err
		}
		var var_23 = []any{bigPri}
		err = templ.RenderCSSItems(ctx, templBuffer, var_23...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div _=\"on click toggle .hidden on .quick-fine\" class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_23).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_24 := `Toggle Quick Fines`
		_, err = templBuffer.WriteString(var_24)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div></div>")
		if err != nil {
			return err
		}
		err = playersManage(secureFineMasterbaseUrl(finemasterBaseUrl, pass), players, qp.PlayerOpen).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		err = presetFines(fmt.Sprintf("%s/%s", finemasterBaseUrl, pass), qp.PresetFinesOpen, pFines).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div id=\"fine-list-container\" class=\"h-screen\" hx-get=\"/fines\" hx-trigger=\"load once\" hx-swap=\"OuterHTML\" class=\"w-full text-center\">")
		if err != nil {
			return err
		}
		var_25 := `loading latest..`
		_, err = templBuffer.WriteString(var_25)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"mt-96\"></div>")
		if err != nil {
			return err
		}
		err = finemasterNav(fmt.Sprintf("%s/%s", finemasterBaseUrl, pass)).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</body><script src=\"https://unpkg.com/htmx.org\">")
		if err != nil {
			return err
		}
		var_26 := ``
		_, err = templBuffer.WriteString(var_26)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script>")
		if err != nil {
			return err
		}
		err = tomSelectLinks().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</html>")
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
		var_27 := templ.GetChildren(ctx)
		if var_27 == nil {
			var_27 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if isOpen {
			_, err = templBuffer.WriteString("<div class=\"px-8 py-6 text-left bg-white shadow-xl m-10\" id=\"preset-fine\"><h3 class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			var_28 := `Add or Approve Fines`
			_, err = templBuffer.WriteString(var_28)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3><div class=\"w-full flex justify-center items-center py-2\"><p>")
			if err != nil {
				return err
			}
			var_29 := `Approve fines submitted, or add new fines`
			_, err = templBuffer.WriteString(var_29)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p></div><div class=\"grid grid-cols-1 md:grid-cols-2 gap-4\"><div class=\"text-2xl\"><h1 class=\"font-bold text-center\">")
			if err != nil {
				return err
			}
			var_30 := `Approve/Existing Preset Fines`
			_, err = templBuffer.WriteString(var_30)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1><p>")
			if err != nil {
				return err
			}
			var_31 := `(These will be added to the quick toggle list)`
			_, err = templBuffer.WriteString(var_31)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p>")
			if err != nil {
				return err
			}
			for _, pfine := range presetFines {
				if !pfine.Approved {
					_, err = templBuffer.WriteString("<div class=\"mt-2 text-center\"><div>")
					if err != nil {
						return err
					}
					var var_32 = []any{bigAdd}
					err = templ.RenderCSSItems(ctx, templBuffer, var_32...)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("<button hx-post=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/preset-fines/approve?pfid=%d", pfine.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" class=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_32).String()))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var_33 := `☐ ✨ `
					_, err = templBuffer.WriteString(var_33)
					if err != nil {
						return err
					}
					var var_34 string = pfine.Reason
					_, err = templBuffer.WriteString(templ.EscapeString(var_34))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(" ")
					if err != nil {
						return err
					}
					var var_35 string = fmt.Sprintf("$%.0f", pfine.Amount)
					_, err = templBuffer.WriteString(templ.EscapeString(var_35))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</button>")
					if err != nil {
						return err
					}
					var var_36 = []any{bigDel}
					err = templ.RenderCSSItems(ctx, templBuffer, var_36...)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("<button class=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_36).String()))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" hx-delete=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/preset-fines?pfid=%d", pfine.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" hx-confirm=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("Remove %s from standard fines?", pfine.Reason)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var_37 := `Delete`
					_, err = templBuffer.WriteString(var_37)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</button></div></div>")
					if err != nil {
						return err
					}
				}
			}
			for _, pfine := range presetFines {
				if pfine.Approved {
					_, err = templBuffer.WriteString("<div class=\"mt-2 text-center\"><div>")
					if err != nil {
						return err
					}
					var var_38 string = pfine.Reason
					_, err = templBuffer.WriteString(templ.EscapeString(var_38))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(" ")
					if err != nil {
						return err
					}
					var_39 := `- `
					_, err = templBuffer.WriteString(var_39)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(" ")
					if err != nil {
						return err
					}
					var var_40 string = fmt.Sprintf("$%.0f", pfine.Amount)
					_, err = templBuffer.WriteString(templ.EscapeString(var_40))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(" ")
					if err != nil {
						return err
					}
					var var_41 = []any{bigDel}
					err = templ.RenderCSSItems(ctx, templBuffer, var_41...)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("<button class=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_41).String()))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" hx-delete=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/preset-fines?pfid=%d", pfine.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" hx-confirm=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("Remove %s from standard fines?", pfine.Reason)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var_42 := `Delete`
					_, err = templBuffer.WriteString(var_42)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</button>")
					if err != nil {
						return err
					}
					if pfine.NotQuickFine {
						var_43 := `Hidden`
						_, err = templBuffer.WriteString(var_43)
						if err != nil {
							return err
						}
					} else {
						var var_44 = []any{bigPri}
						err = templ.RenderCSSItems(ctx, templBuffer, var_44...)
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString("<button class=\"")
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_44).String()))
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString("\" hx-post=\"")
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/preset-fines/hide?pfid=%d", pfine.ID)))
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString("\" hx-confirm=\"")
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("Hide %s from quick fines? (delete to remove from \"Add Fine\" drop-down too)", pfine.Reason)))
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString("\">")
						if err != nil {
							return err
						}
						var_45 := `Hide from Quick Fines`
						_, err = templBuffer.WriteString(var_45)
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString("</button>")
						if err != nil {
							return err
						}
					}
					_, err = templBuffer.WriteString("</div></div>")
					if err != nil {
						return err
					}
				}
			}
			_, err = templBuffer.WriteString("</div><div><form hx-post=\"/preset-fines\" class=\"mt-4\"><div><label for=\"reason\" class=\"block\"><input required type=\"text\" name=\"reason\" id=\"reason\" placeholder=\"Reason for the fine\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></label></div><div class=\"mt-4\"><label for=\"amount\" class=\"block\">")
			if err != nil {
				return err
			}
			var_46 := `Amount ($)`
			_, err = templBuffer.WriteString(var_46)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input required type=\"number\" step=\"0.01\" name=\"amount\" id=\"amount\" placeholder=\"Amount\" class=\"p-2 w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></div><div class=\"mt-4\"><label class=\"block\"><input type=\"checkbox\" checked=\"checked\" name=\"approved\" class=\"text-2xl m-2 py-2 mt-2 p-2 \">")
			if err != nil {
				return err
			}
			var_47 := `Approved`
			_, err = templBuffer.WriteString(var_47)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label></div><div class=\"flex items-center justify-between mt-4\">")
			if err != nil {
				return err
			}
			var var_48 = []any{bigAdd}
			err = templ.RenderCSSItems(ctx, templBuffer, var_48...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_48).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_49 := `Add Preset Fine`
			_, err = templBuffer.WriteString(var_49)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div></form><div class=\"flex justify-center w-full\">")
			if err != nil {
				return err
			}
			var var_50 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_50...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_51 templ.SafeURL = makeSafeUrl(baseUrl, false, false, false, false)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_51)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_50).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_52 := `Close`
			_, err = templBuffer.WriteString(var_52)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></div></div></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"flex justify-center w-full\">")
			if err != nil {
				return err
			}
			var var_53 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_53...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_54 templ.SafeURL = makeSafeUrlWithAnchor(baseUrl, false, false, false, true, "preset-fine")
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_54)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_53).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_55 := `Manage Standard Fines`
			_, err = templBuffer.WriteString(var_55)
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

func playersManage(baseUrl string, players []PlayerWithFines, isOpen bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_56 := templ.GetChildren(ctx)
		if var_56 == nil {
			var_56 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if isOpen {
			_, err = templBuffer.WriteString("<div class=\"flex items-center justify-center bg-gray-100  p-3\" id=\"players-manage\"><div class=\"px-8 py-6 text-left bg-white shadow-xl m-10\"><!--")
			if err != nil {
				return err
			}
			var_57 := ` Section for Adding New Player `
			_, err = templBuffer.WriteString(var_57)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("--><h3 class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			var_58 := `Add New Player`
			_, err = templBuffer.WriteString(var_58)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3><form hx-post=\"/players\" method=\"POST\" class=\"mt-4\"><div><label for=\"name\" class=\"block\">")
			if err != nil {
				return err
			}
			var_59 := `Name`
			_, err = templBuffer.WriteString(var_59)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input required type=\"text\" name=\"name\" id=\"name\" placeholder=\"Name\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></div><div class=\"flex items-center w-full\">")
			if err != nil {
				return err
			}
			var var_60 = []any{bigAdd}
			err = templ.RenderCSSItems(ctx, templBuffer, var_60...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_60).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_61 := `Add Player`
			_, err = templBuffer.WriteString(var_61)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div></form><div class=\"p-6 max-w-36\"><!--")
			if err != nil {
				return err
			}
			var_62 := ` Section for Deleting Existing Player `
			_, err = templBuffer.WriteString(var_62)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("--><h3 class=\"text-2xl font-bold text-center mt-8\">")
			if err != nil {
				return err
			}
			var_63 := `Delete Player`
			_, err = templBuffer.WriteString(var_63)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3>")
			if err != nil {
				return err
			}
			for _, p := range players {
				var var_64 = []any{bigDel}
				err = templ.RenderCSSItems(ctx, templBuffer, var_64...)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("<button type=\"submit\" hx-delete=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/players?playerId=%d", p.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" hx-confirm=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("Are you sure you want to delete %s?", p.Name)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" class=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_64).String()))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var_65 := `Delete `
				_, err = templBuffer.WriteString(var_65)
				if err != nil {
					return err
				}
				var var_66 string = p.Name
				_, err = templBuffer.WriteString(templ.EscapeString(var_66))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</button>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div><div class=\"flex justify-center w-full mt-4\">")
			if err != nil {
				return err
			}
			var var_67 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_67...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_68 templ.SafeURL = makeSafeUrl(baseUrl, false, false, false, false)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_68)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-transition=\"true\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_67).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_69 := `Close`
			_, err = templBuffer.WriteString(var_69)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></div></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"flex justify-center w-full p-3\">")
			if err != nil {
				return err
			}
			var var_70 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_70...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_71 templ.SafeURL = makeSafeUrlWithAnchor(baseUrl, false, false, true, false, "players-manage")
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_71)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-transition=\"true\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_70).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_72 := `Manage Players`
			_, err = templBuffer.WriteString(var_72)
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
