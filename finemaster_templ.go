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
		_, err = templBuffer.WriteString("<html hx-boost=\"true\">")
		if err != nil {
			return err
		}
		err = header().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body><div class=\"bg-gray-900 text-center p-5\"><h1 class=\"text-xl md:text-3xl lg:text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-teal-400\">")
		if err != nil {
			return err
		}
		var_2 := `Fine Master Zone`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1></div><div class=\"grid grid-cols-1 gap-2\">")
		if err != nil {
			return err
		}
		err = fineAdd(fmt.Sprintf("%s/%s", finemasterBaseUrl, pass), qp.FinesOpen, players, pFines, true).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		err = playersManage(fmt.Sprintf("%s/%s", finemasterBaseUrl, pass), qp.PlayerOpen).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"container mx-auto p-4\"><div class=\"flex items-center justify-center bg-gray-100 mx-auto\"><ul>")
		if err != nil {
			return err
		}
		for _, p := range players {
			_, err = templBuffer.WriteString("<li class=\"mb-2\">")
			if err != nil {
				return err
			}
			var var_3 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_3...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<div _=\"on click toggle .hidden on next &lt;section/&gt;\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_3).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var var_4 string = p.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_4))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_5 := `- `
			_, err = templBuffer.WriteString(var_5)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_6 string = fmt.Sprintf("$%d (%d)", p.TotalFines, p.TotalFineCount)
			_, err = templBuffer.WriteString(templ.EscapeString(var_6))
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
				var var_7 string = f.Reason
				_, err = templBuffer.WriteString(templ.EscapeString(var_7))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var_8 := `- `
				_, err = templBuffer.WriteString(var_8)
				if err != nil {
					return err
				}
				var var_9 string = fmt.Sprintf("$%.0f", f.Amount)
				_, err = templBuffer.WriteString(templ.EscapeString(var_9))
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
				if pf.Approved {
					_, err = templBuffer.WriteString("<form class=\" inline-flex mx-2 space-y-2\">")
					if err != nil {
						return err
					}
					var var_10 = []any{fmt.Sprintf("fine-group-%d-%d", pf.ID, p.PlayerID)}
					err = templ.RenderCSSItems(ctx, templBuffer, var_10...)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("<div hidden class=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_10).String()))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\"><input type=\"hidden\" name=\"playerId\" value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", p.PlayerID)))
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
					_, err = templBuffer.WriteString("\"><input type=\"hidden\" name=\"approved\" value=\"on\"></div>")
					if err != nil {
						return err
					}
					var var_11 = []any{bigAdd}
					err = templ.RenderCSSItems(ctx, templBuffer, var_11...)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("<button hx-post=\"/fines\" hx-include=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf(".fine-group-%d-%d", pf.ID, p.PlayerID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" class=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_11).String()))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var var_12 string = fmt.Sprintf("%s ($%v)", pf.Reason, pf.Amount)
					_, err = templBuffer.WriteString(templ.EscapeString(var_12))
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
		var var_13 = []any{bigPri}
		err = templ.RenderCSSItems(ctx, templBuffer, var_13...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div _=\"on click toggle .hidden on .quick-fine\" class=\"")
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
		var_14 := `Toggle Quick Fines`
		_, err = templBuffer.WriteString(var_14)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div></div><div id=\"fine-list-container\" hx-get=\"/fines\" hx-trigger=\"load once\" hx-swap=\"OuterHTML\" class=\"w-full text-center\">")
		if err != nil {
			return err
		}
		var_15 := `loading latest..`
		_, err = templBuffer.WriteString(var_15)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		err = presetFines(fmt.Sprintf("%s/%s", finemasterBaseUrl, pass), qp.PresetFinesOpen, pFines).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		err = pageFooter().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</body></html>")
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
		var_16 := templ.GetChildren(ctx)
		if var_16 == nil {
			var_16 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if isOpen {
			_, err = templBuffer.WriteString("<div class=\"px-8 py-6 text-left bg-white shadow-xl m-10\" id=\"preset-fine\"><h3 class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			var_17 := `Add or Approve Fines`
			_, err = templBuffer.WriteString(var_17)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3><div class=\"w-full flex justify-center items-center py-2\"><p>")
			if err != nil {
				return err
			}
			var_18 := `Approve fines submitted, or add new fines`
			_, err = templBuffer.WriteString(var_18)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p></div><div class=\"grid grid-cols-1 md:grid-cols-2 gap-4\"><div class=\"text-2xl\"><h1 class=\"font-bold text-center\">")
			if err != nil {
				return err
			}
			var_19 := `Approve/Existing Preset Fines`
			_, err = templBuffer.WriteString(var_19)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1>")
			if err != nil {
				return err
			}
			for _, fine := range presetFines {
				_, err = templBuffer.WriteString("<div class=\"mt-2 text-center\"><div>")
				if err != nil {
					return err
				}
				if !fine.Approved {
					_, err = templBuffer.WriteString("<button hx-post=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/preset-fines/approve?pfid=%d", fine.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var_20 := `☐ ✨`
					_, err = templBuffer.WriteString(var_20)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</button>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_21 string = fine.Reason
				_, err = templBuffer.WriteString(templ.EscapeString(var_21))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var_22 := `- `
				_, err = templBuffer.WriteString(var_22)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_23 string = fmt.Sprintf("$%.0f", fine.Amount)
				_, err = templBuffer.WriteString(templ.EscapeString(var_23))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" <button hx-delete=\"")
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
				var_24 := `🗑`
				_, err = templBuffer.WriteString(var_24)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</button></div></div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div><div><form hx-post=\"/preset-fines\" class=\"mt-4\"><div><label for=\"reason\" class=\"block\"><input required type=\"text\" name=\"reason\" id=\"reason\" placeholder=\"Reason for the fine\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></label></div><div class=\"mt-4\"><label for=\"amount\" class=\"block\">")
			if err != nil {
				return err
			}
			var_25 := `Amount ($)`
			_, err = templBuffer.WriteString(var_25)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input required type=\"number\" step=\"0.01\" name=\"amount\" id=\"amount\" placeholder=\"Amount\" class=\"p-2 w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></div><div class=\"mt-4\"><label class=\"block\"><input type=\"checkbox\" checked=\"checked\" name=\"approved\" class=\"text-2xl m-2 py-2 mt-2 p-2 \">")
			if err != nil {
				return err
			}
			var_26 := `Approved`
			_, err = templBuffer.WriteString(var_26)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label></div><div class=\"flex items-center justify-between mt-4\">")
			if err != nil {
				return err
			}
			var var_27 = []any{bigAdd}
			err = templ.RenderCSSItems(ctx, templBuffer, var_27...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
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
			var_28 := `Add Preset Fine`
			_, err = templBuffer.WriteString(var_28)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div></form><div class=\"flex justify-center w-full\">")
			if err != nil {
				return err
			}
			var var_29 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_29...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_30 templ.SafeURL = makeSafeUrl(baseUrl, false, false, false)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_30)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_29).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_31 := `Close`
			_, err = templBuffer.WriteString(var_31)
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
			var var_32 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_32...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_33 templ.SafeURL = makeSafeUrlWithAnchor(baseUrl, false, false, true, "preset-fine")
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_33)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"")
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
			var_34 := `Add Fines`
			_, err = templBuffer.WriteString(var_34)
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

func playersManage(baseUrl string, isOpen bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_35 := templ.GetChildren(ctx)
		if var_35 == nil {
			var_35 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if isOpen {
			_, err = templBuffer.WriteString("<div class=\"flex items-center justify-center bg-gray-100\" id=\"players-manage\"><div class=\"px-8 py-6 text-left bg-white shadow-xl m-10\"><!--")
			if err != nil {
				return err
			}
			var_36 := ` Section for Adding New Player `
			_, err = templBuffer.WriteString(var_36)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("--><h3 class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			var_37 := `Add New Player`
			_, err = templBuffer.WriteString(var_37)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3><form hx-post=\"/players\" method=\"POST\" class=\"mt-4\"><div><label for=\"name\" class=\"block\">")
			if err != nil {
				return err
			}
			var_38 := `Name`
			_, err = templBuffer.WriteString(var_38)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input required type=\"text\" name=\"name\" id=\"name\" placeholder=\"Name\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></div><div class=\"flex items-center w-full\">")
			if err != nil {
				return err
			}
			var var_39 = []any{bigAdd}
			err = templ.RenderCSSItems(ctx, templBuffer, var_39...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
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
			var_40 := `Add Player`
			_, err = templBuffer.WriteString(var_40)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div></form><!--")
			if err != nil {
				return err
			}
			var_41 := ` Section for Deleting Existing Player `
			_, err = templBuffer.WriteString(var_41)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("--><h3 class=\"text-2xl font-bold text-center mt-8\">")
			if err != nil {
				return err
			}
			var_42 := `Delete Player`
			_, err = templBuffer.WriteString(var_42)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3><form hx-delete=\"/players\" method=\"POST\" class=\"mt-4\"><div><label for=\"playerSelect\" class=\"block\">")
			if err != nil {
				return err
			}
			var_43 := `Select Player`
			_, err = templBuffer.WriteString(var_43)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><select name=\"playerId\" id=\"playerSelect\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"><option value=\"1\">")
			if err != nil {
				return err
			}
			var_44 := `Player 1`
			_, err = templBuffer.WriteString(var_44)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option><option value=\"2\">")
			if err != nil {
				return err
			}
			var_45 := `Player 2`
			_, err = templBuffer.WriteString(var_45)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option></select></div><div class=\"flex items-center w-full\">")
			if err != nil {
				return err
			}
			var var_46 = []any{bigDel}
			err = templ.RenderCSSItems(ctx, templBuffer, var_46...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button type=\"submit\" hx-confirm=\"Are you sure you want to delete this player?\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_46).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_47 := `Delete Player`
			_, err = templBuffer.WriteString(var_47)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div></form><div class=\"flex justify-center w-full mt-4\">")
			if err != nil {
				return err
			}
			var var_48 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_48...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_49 templ.SafeURL = makeSafeUrl(baseUrl, false, false, false)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_49)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-transition=\"true\" class=\"")
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
			var_50 := `Close`
			_, err = templBuffer.WriteString(var_50)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></div></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"flex justify-center w-full\">")
			if err != nil {
				return err
			}
			var var_51 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_51...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_52 templ.SafeURL = makeSafeUrlWithAnchor(baseUrl, false, true, false, "players-manage")
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_52)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-transition=\"true\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_51).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_53 := `Manage Players`
			_, err = templBuffer.WriteString(var_53)
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
