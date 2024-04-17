// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"github.com/dustin/go-humanize"
)

func getFinesTotal(fines []FineWithPlayer) float64 {
	if len(fines) > 0 {
		var total float64 = 0
		for _, f := range fines {
			total += f.Fine.Amount
		}
		return total
	}
	return 0
}

func fineList(fines []FineWithPlayer, page int, presetFineUpdated uint, isFineMaster bool) templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"m-2 bg-gray-200 h-full shadow-xl p-2 h-full mt-10\" id=\"fine-list-container\"><div class=\"text-center\"><div class=\"flex justify-center items-center mb-4 cursor-pointer\" hx-get=\"/fines\" hx-target=\"#fine-list-container\" hx-swap=\"outerHTML\" hx-trigger=\"click\"><span class=\"flex-grow text-center font-bold\">")
		if err != nil {
			return err
		}
		var_2 := `Recent Fine List`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span><span class=\"text-3xl ml-2\">")
		if err != nil {
			return err
		}
		var_3 := `↻`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></div>")
		if err != nil {
			return err
		}
		if isFineMaster {
			_, err = templBuffer.WriteString("<div>")
			if err != nil {
				return err
			}
			var var_4 string = fmt.Sprintf("$%v", getFinesTotal(fines))
			_, err = templBuffer.WriteString(templ.EscapeString(var_4))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<table id=\"fine-list\" class=\"min-w-full\"><tbody class=\"divide-y divide-gray-900\">")
		if err != nil {
			return err
		}
		for _, f := range fines {
			err = fineRow(isFineMaster, f).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			if presetFineUpdated == f.Fine.ID {
				err = success("updated!").Render(ctx, templBuffer)
				if err != nil {
					return err
				}
			}
		}
		_, err = templBuffer.WriteString("</tbody></table></div></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func fineRow(isFineMaster bool, f FineWithPlayer) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_5 := templ.GetChildren(ctx)
		if var_5 == nil {
			var_5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<tr id=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("fr-%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"bg-gray-200 p-2\"><td class=\"p-2 text-gray-900 flex flex-col text-wrap\"><div class=\"text-bold text-3xl\">")
		if err != nil {
			return err
		}
		var var_6 string = f.Player.Name
		_, err = templBuffer.WriteString(templ.EscapeString(var_6))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"text-2xl\">")
		if err != nil {
			return err
		}
		var var_7 string = f.Fine.Reason
		_, err = templBuffer.WriteString(templ.EscapeString(var_7))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"italic\">")
		if err != nil {
			return err
		}
		if f.Fine.Approved {
			var var_8 string = fmt.Sprintf("$%v - ", f.Fine.Amount)
			_, err = templBuffer.WriteString(templ.EscapeString(var_8))
			if err != nil {
				return err
			}
		}
		if len(f.Match.Opponent) > 0 {
			var var_9 string = f.Match.Opponent
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
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			if f.Match.StartTime != nil {
				var var_11 string = niceDate(f.Match.StartTime)
				_, err = templBuffer.WriteString(templ.EscapeString(var_11))
				if err != nil {
					return err
				}
			}
		} else {
			var_12 := `(`
			_, err = templBuffer.WriteString(var_12)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_13 string = humanize.Time(f.Fine.FineAt)
			_, err = templBuffer.WriteString(templ.EscapeString(var_13))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_14 := `)`
			_, err = templBuffer.WriteString(var_14)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div></td><td><div class=\"text-lg text-gray-900 text-wrap\">")
		if err != nil {
			return err
		}
		var var_15 string = f.Fine.Context
		_, err = templBuffer.WriteString(templ.EscapeString(var_15))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <div>")
		if err != nil {
			return err
		}
		var var_16 = []any{smPri}
		err = templ.RenderCSSItems(ctx, templBuffer, var_16...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button hx-get=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/edit/%d?isContext=true", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"innerHTML\" hx-target=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("#fr-%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_16).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		if len(f.Fine.Context) == 0 {
			var_17 := `Add Context	`
			_, err = templBuffer.WriteString(var_17)
			if err != nil {
				return err
			}
		} else {
			var_18 := `Edit Context`
			_, err = templBuffer.WriteString(var_18)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</button></div></div></td><td>")
		if err != nil {
			return err
		}
		if isFineMaster {
			if !f.Fine.Approved {
				_, err = templBuffer.WriteString("<div><form hx-post=\"/fines/approve\" hx-swap=\"outerHTML\" ethod=\"POST\"><input type=\"hidden\" name=\"fid\" value=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", f.Fine.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"> ")
				if err != nil {
					return err
				}
				var_19 := `$`
				_, err = templBuffer.WriteString(var_19)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" <input type=\"number\" name=\"amount\" id=\"amount-input-3\" class=\"px-2 py-1 border rounded\" placeholder=\"Set amount\"")
				if err != nil {
					return err
				}
				if f.Fine.Amount > 0 {
					_, err = templBuffer.WriteString(" value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", f.Fine.Amount)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\"")
					if err != nil {
						return err
					}
				} else {
					_, err = templBuffer.WriteString(" value=\"2\"")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				var var_20 = []any{bigAdd}
				err = templ.RenderCSSItems(ctx, templBuffer, var_20...)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_20).String()))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var_21 := `Approve`
				_, err = templBuffer.WriteString(var_21)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</button></form></div>")
				if err != nil {
					return err
				}
			}
		}
		if len(f.Fine.Contest) == 0 {
			_, err = templBuffer.WriteString("<div class=\"max-w-96\">")
			if err != nil {
				return err
			}
			var var_22 = []any{smPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_22...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button hx-get=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/edit/%d?isContest=true", f.Fine.ID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_22).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_23 := `Contest`
			_, err = templBuffer.WriteString(var_23)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div>")
			if err != nil {
				return err
			}
			var var_24 string = f.Fine.Contest
			_, err = templBuffer.WriteString(templ.EscapeString(var_24))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</td>")
		if err != nil {
			return err
		}
		if isFineMaster {
			_, err = templBuffer.WriteString("<td>")
			if err != nil {
				return err
			}
			var var_25 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_25...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button hx-get=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/edit/%d?isEdit=true", f.Fine.ID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-target=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("#fr-%d", f.Fine.ID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-target-error=\"#any-errors\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_25).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_26 := `edit`
			_, err = templBuffer.WriteString(var_26)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></td>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</tr>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func fineContextRow(f FineWithPlayer, matches []Match) templ.Component {
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
		_, err = templBuffer.WriteString("<td colspan=\"7\"><div class=\"border rounded-lg flex flex-col items-center p-4 space-y-4 w-full max-w-4xl mx-auto\"><div class=\"text-center w-full\">")
		if err != nil {
			return err
		}
		var var_28 string = f.Player.Name
		_, err = templBuffer.WriteString(templ.EscapeString(var_28))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_29 := `- `
		_, err = templBuffer.WriteString(var_29)
		if err != nil {
			return err
		}
		var var_30 string = fmt.Sprintf("$%v - ", f.Fine.Amount)
		_, err = templBuffer.WriteString(templ.EscapeString(var_30))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_31 := `-  `
		_, err = templBuffer.WriteString(var_31)
		if err != nil {
			return err
		}
		var var_32 string = f.Fine.Reason
		_, err = templBuffer.WriteString(templ.EscapeString(var_32))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><p class=\"text-sm w-full text-gray-700\">")
		if err != nil {
			return err
		}
		var_33 := `Add (optional) context for this fine:`
		_, err = templBuffer.WriteString(var_33)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p><div class=\"w-full\"><label class=\"block w-full\">")
		if err != nil {
			return err
		}
		var_34 := `Context:`
		_, err = templBuffer.WriteString(var_34)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <input type=\"text\" name=\"context\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Fine.Context))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"px-4 py-2 border border-gray-300 rounded-lg w-full focus:ring-blue-500 focus:border-blue-500\" placeholder=\"Reason\"></label></div><div class=\"mt-2\"><div hx-get=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/match-list?type=select&matchId=%d", f.Fine.MatchId)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-trigger=\"load once\"></div></div><div>")
		if err != nil {
			return err
		}
		var_35 := `OR`
		_, err = templBuffer.WriteString(var_35)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div class=\"w-full\" id=\"dateInputDiv\"><label class=\"block\">")
		if err != nil {
			return err
		}
		var_36 := `Date/Time:`
		_, err = templBuffer.WriteString(var_36)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <input type=\"datetime-local\" id=\"fineAt\" name=\"fineAt\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Fine.FineAt.Format("2006-01-02T15:04")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"px-2 py-1 border rounded\"><p class=\"italic text-md\">")
		if err != nil {
			return err
		}
		var_37 := `(defaults to create time)`
		_, err = templBuffer.WriteString(var_37)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p></label></div><input type=\"hidden\" name=\"fid\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><button class=\"px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600\" hx-post=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString("/fines/context"))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-include=\"closest tr\" type=\"submit\">")
		if err != nil {
			return err
		}
		var_38 := `Save`
		_, err = templBuffer.WriteString(var_38)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button><div class=\"w-full text-center mt-4\">")
		if err != nil {
			return err
		}
		var var_39 = []any{sec}
		err = templ.RenderCSSItems(ctx, templBuffer, var_39...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<a href=\"")
		if err != nil {
			return err
		}
		var var_40 templ.SafeURL = templ.SafeURL("/#fine-list-container")
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_40)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"")
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
		var_41 := `Close`
		_, err = templBuffer.WriteString(var_41)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a></div></div></td>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func fineContestRow(f FineWithPlayer) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_42 := templ.GetChildren(ctx)
		if var_42 == nil {
			var_42 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"border rounded-lg flex flex-col items-center p-4 space-y-4\" id=\"contest-form\"><p class=\"text-sm w-full text-gray-700\">")
		if err != nil {
			return err
		}
		var_43 := `Contest fine:`
		_, err = templBuffer.WriteString(var_43)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p><div class=\"w-full\"><input type=\"text\" name=\"contest\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Fine.Contest))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"px-4 py-2 border border-gray-300 rounded-lg w-full focus:ring-blue-500 focus:border-blue-500\" placeholder=\"Why do you contest this fine?\"></div><input type=\"hidden\" name=\"fid\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var var_44 = []any{add}
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
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/contest")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-include=\"closest #contest-form\" type=\"submit\">")
		if err != nil {
			return err
		}
		var_45 := `Save`
		_, err = templBuffer.WriteString(var_45)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button><div class=\"w-full text-center\">")
		if err != nil {
			return err
		}
		var var_46 = []any{sec}
		err = templ.RenderCSSItems(ctx, templBuffer, var_46...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<a href=\"")
		if err != nil {
			return err
		}
		var var_47 templ.SafeURL = templ.SafeURL("/#fine-list-container")
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_47)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"")
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
		var_48 := `Cancel`
		_, err = templBuffer.WriteString(var_48)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a></div></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func fineEditRow(f FineWithPlayer) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_49 := templ.GetChildren(ctx)
		if var_49 == nil {
			var_49 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<tr id=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("fr-%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"bg-white divide-y divide-gray-200\"><td class=\"px-6 py-4\"><label for=\"reason\">")
		if err != nil {
			return err
		}
		var_50 := `Reason`
		_, err = templBuffer.WriteString(var_50)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"reason\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Fine.Reason))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"px-2 py-1 border rounded w-full\" placeholder=\"Reason\"><input type=\"text\" name=\"context\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Fine.Context))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"px-2 py-1 border rounded w-full\" placeholder=\"Context\"><input type=\"playerId\" name=\"playerId\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", f.Player.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"invisible\"><input type=\"fid\" name=\"fid\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"invisible\"></td><td class=\"px-6 py-4\"><input type=\"number\" name=\"amount\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", f.Fine.Amount)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"px-2 py-1 border rounded w-full\" placeholder=\"Amount\"></td><td>")
		if err != nil {
			return err
		}
		var var_51 string = f.Player.Name
		_, err = templBuffer.WriteString(templ.EscapeString(var_51))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</td><td><select name=\"approved\" class=\"px-2 py-1 border rounded\"><option value=\"true\"")
		if err != nil {
			return err
		}
		if f.Fine.Approved {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_52 := `Approved`
		_, err = templBuffer.WriteString(var_52)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"false\"")
		if err != nil {
			return err
		}
		if !f.Fine.Approved {
			_, err = templBuffer.WriteString(" selected")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var_53 := `Not Approved`
		_, err = templBuffer.WriteString(var_53)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option></select></td><td><input type=\"date\" id=\"fineAt\" name=\"fineAt\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Fine.FineAt.Format("2006-01-02")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"px-2 py-1 border rounded w-full\"></td><td>")
		if err != nil {
			return err
		}
		var var_54 = []any{bigDel}
		err = templ.RenderCSSItems(ctx, templBuffer, var_54...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_54).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-confirm=\"Are you sure you want to delete the fine by this player?\" hx-delete=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines?fid=%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_55 := `Delete`
		_, err = templBuffer.WriteString(var_55)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button>")
		if err != nil {
			return err
		}
		var var_56 = []any{pri}
		err = templ.RenderCSSItems(ctx, templBuffer, var_56...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_56).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-post=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/edit/%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-target=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("#fr-%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-include=\"closest tr\" type=\"submit\">")
		if err != nil {
			return err
		}
		var_57 := `Save`
		_, err = templBuffer.WriteString(var_57)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button>")
		if err != nil {
			return err
		}
		var var_58 = []any{sec}
		err = templ.RenderCSSItems(ctx, templBuffer, var_58...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button hx-get=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/edit/%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-target=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("#fr-%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" type=\"button\" class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_58).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_59 := `Cancel`
		_, err = templBuffer.WriteString(var_59)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></td></tr>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
