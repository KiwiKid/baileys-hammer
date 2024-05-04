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
		_, err = templBuffer.WriteString("<table id=\"fine-list\" class=\"min-w-full mb-36\"><tbody class=\"divide-y divide-gray-900\">")
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
		_, err = templBuffer.WriteString("\" class=\"bg-gray-200 p-2 \"><td class=\"p-2 text-gray-900 flex flex-col text-wrap\"><div class=\"text-bold text-3xl\">")
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
		_, err = templBuffer.WriteString("</div><div class=\"text-gray-900 text-wrap\">")
		if err != nil {
			return err
		}
		var var_8 string = f.Fine.Context
		_, err = templBuffer.WriteString(templ.EscapeString(var_8))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <div")
		if err != nil {
			return err
		}
		if len(f.Fine.Context) > 0 {
			_, err = templBuffer.WriteString(" class=\"w-3/5\"")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("></div></div><div class=\"italic\">")
		if err != nil {
			return err
		}
		if f.Fine.Approved {
			var var_9 string = fmt.Sprintf("$%v - ", f.Fine.Amount)
			_, err = templBuffer.WriteString(templ.EscapeString(var_9))
			if err != nil {
				return err
			}
		}
		if len(f.Match.Opponent) > 0 {
			var var_10 string = f.Match.Opponent
			_, err = templBuffer.WriteString(templ.EscapeString(var_10))
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var var_11 string = niceDate(&f.Fine.FineAt)
		_, err = templBuffer.WriteString(templ.EscapeString(var_11))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div></td><td><div class=\"m-2\">")
		if err != nil {
			return err
		}
		var var_12 = []any{smPri}
		err = templ.RenderCSSItems(ctx, templBuffer, var_12...)
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
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_12).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		if len(f.Fine.Context) == 0 {
			var_13 := `Add Context	`
			_, err = templBuffer.WriteString(var_13)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			if f.Fine.MatchId == 0 {
				_, err = templBuffer.WriteString("<span class=\"text-sm\">")
				if err != nil {
					return err
				}
				var_14 := `⚠️`
				_, err = templBuffer.WriteString(var_14)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
			}
		} else {
			var_15 := `Edit Context`
			_, err = templBuffer.WriteString(var_15)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</button></div><div class=\"m-2\">")
		if err != nil {
			return err
		}
		if len(f.Fine.Contest) == 0 {
			_, err = templBuffer.WriteString("<div class=\"max-w-96\">")
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
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/edit/%d?isContest=true&isFineMaster=%d", f.Fine.ID, isFineMaster)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" class=\"")
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
			var_17 := `Contest`
			_, err = templBuffer.WriteString(var_17)
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
			var var_18 string = f.Fine.Contest
			_, err = templBuffer.WriteString(templ.EscapeString(var_18))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div></td><td>")
		if err != nil {
			return err
		}
		if isFineMaster {
			if !f.Fine.Approved {
				_, err = templBuffer.WriteString("<form hx-post=\"/fines/approve\" hx-swap=\"outerHTML\" class=\"w-full\" method=\"POST\"><div class=\"grid grid-cols-2 gap-4\"><!--")
				if err != nil {
					return err
				}
				var_19 := ` Corrected class name and added gap for spacing `
				_, err = templBuffer.WriteString(var_19)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("--><div class=\"flex items-center\"><!--")
				if err != nil {
					return err
				}
				var_20 := ` Added flex layout to vertically center align items `
				_, err = templBuffer.WriteString(var_20)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("--><input type=\"hidden\" name=\"fid\" value=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", f.Fine.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"><input type=\"hidden\" name=\"approved\" value=\"on\"><div class=\"mt-2\"><label for=\"amount-input-3\" class=\"block text-lg text-gray-700 text-lg font-semibold\">")
				if err != nil {
					return err
				}
				var_21 := `Amount:`
				_, err = templBuffer.WriteString(var_21)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</label><div class=\"mt-1 flex rounded-md shadow-sm\"><span class=\"inline-flex items-center px-3 rounded-l-md border border-r-0 border-gray-300 bg-gray-50 text-gray-500 text-lg\">")
				if err != nil {
					return err
				}
				var_22 := `$`
				_, err = templBuffer.WriteString(var_22)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span><input type=\"number\" name=\"amount\" id=\"amount-input-3\" class=\"flex-1 min-w-0 block w-full px-3 py-2 rounded-none rounded-r-md border border-gray-300\" placeholder=\"Set amount\"")
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
				_, err = templBuffer.WriteString("></div></div></div><div class=\"flex justify-end items-center\"><!--")
				if err != nil {
					return err
				}
				var_23 := ` Added flex layout to align button to the right and center it vertically `
				_, err = templBuffer.WriteString(var_23)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("-->")
				if err != nil {
					return err
				}
				var var_24 = []any{bigAdd}
				err = templ.RenderCSSItems(ctx, templBuffer, var_24...)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_24).String()))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var_25 := `Approve`
				_, err = templBuffer.WriteString(var_25)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</button></div></div></form>")
				if err != nil {
					return err
				}
			} else {
				_, err = templBuffer.WriteString("<div><form hx-post=\"/fines/approve\" hx-swap=\"outerHTML\" ethod=\"POST\"><input type=\"hidden\" name=\"fid\" value=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", f.Fine.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"><input type=\"hidden\" name=\"approved\" value=\"off\"><input type=\"hidden\" name=\"amount\"")
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
				_, err = templBuffer.WriteString("><button type=\"submit\" class=\"text-3xl bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200\">")
				if err != nil {
					return err
				}
				var_26 := `Decline`
				_, err = templBuffer.WriteString(var_26)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</button></form></div>")
				if err != nil {
					return err
				}
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
			var var_27 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_27...)
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
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_27).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_28 := `edit`
			_, err = templBuffer.WriteString(var_28)
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
		var_29 := templ.GetChildren(ctx)
		if var_29 == nil {
			var_29 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<td colspan=\"7\"><div class=\"border rounded-lg flex flex-col items-center p-4 space-y-4 w-full mx-auto text-3xl\"><div class=\"text-center w-full\">")
		if err != nil {
			return err
		}
		var var_30 string = f.Player.Name
		_, err = templBuffer.WriteString(templ.EscapeString(var_30))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_31 := `- `
		_, err = templBuffer.WriteString(var_31)
		if err != nil {
			return err
		}
		var var_32 string = fmt.Sprintf("$%v - ", f.Fine.Amount)
		_, err = templBuffer.WriteString(templ.EscapeString(var_32))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_33 := `-  `
		_, err = templBuffer.WriteString(var_33)
		if err != nil {
			return err
		}
		var var_34 string = f.Fine.Reason
		_, err = templBuffer.WriteString(templ.EscapeString(var_34))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><p class=\" w-full text-gray-700 text-sm\">")
		if err != nil {
			return err
		}
		var_35 := `Add (optional) context for this fine:`
		_, err = templBuffer.WriteString(var_35)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p><div class=\"w-full\"><label class=\"block w-full text-lg font-semibold\">")
		if err != nil {
			return err
		}
		var_36 := `Context:`
		_, err = templBuffer.WriteString(var_36)
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
		_, err = templBuffer.WriteString("\" class=\"px-4 py-2 border border-gray-300 rounded-lg w-full focus:ring-blue-500 focus:border-blue-500\" placeholder=\"Context\"></label></div><div class=\"mt-2\"><div hx-get=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/match-list?type=select&matchId=%d", f.Fine.MatchId)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-trigger=\"load once\"></div></div><!--")
		if err != nil {
			return err
		}
		var_37 := `<div>OR</div>
			<div class="w-full" id="dateInputDiv">
				<label class="block ">
					Date/Time:
					<input
 						type="datetime-local"
 						id="fineAt"
 						name="fineAt"
 						value={ f.Fine.FineAt.Format("2006-01-02T15:04") }
 						class="px-2 py-1 border rounded"
					/>
					<p class="italic text-md">(defaults to create time)</p>
				</label>
			</div>`
		_, err = templBuffer.WriteString(var_37)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("--><input type=\"hidden\" name=\"fid\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><button class=\"px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600\" hx-post=\"/fines/context\" hx-swap=\"outerHTML\" hx-include=\"closest tr\" type=\"submit\">")
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
		_, err = templBuffer.WriteString("<div class=\"border rounded-lg flex flex-col items-center p-4 space-y-4\" id=\"contest-form\"><p class=\"text-lg font-semibold w-full text-gray-700\">")
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

func fineEditForm(f FineWithPlayer, isFineMaster bool) templ.Component {
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
		_, err = templBuffer.WriteString("<form hx-post=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/edit/%d", f.Fine.ID)))
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
		_, err = templBuffer.WriteString("\"><h1 class=\"text-lg\">")
		if err != nil {
			return err
		}
		var_50 := `Edit Fine`
		_, err = templBuffer.WriteString(var_50)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><h3><pre>")
		if err != nil {
			return err
		}
		var var_51 string = F(`
Name: %s
Reason: %s
Context: %s
Amount: %v
Contest: %s
FineAt: %s
PlayerID: %d
		`,
			f.Player.Name,
			f.Fine.Reason,
			f.Fine.Context,
			f.Fine.Amount,
			f.Fine.Contest,
			f.Fine.FineAt,
			f.Fine.PlayerID,
		)
		_, err = templBuffer.WriteString(templ.EscapeString(var_51))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</pre></h3><div class=\"p-4\"><label for=\"reason\" class=\"block text-lg font-semibold text-gray-700\">")
		if err != nil {
			return err
		}
		var_52 := `Reason`
		_, err = templBuffer.WriteString(var_52)
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
		_, err = templBuffer.WriteString("\" class=\"mt-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 w-full\" placeholder=\"Reason\"><label for=\"context\" class=\"block mt-4 text-lg font-semibold text-gray-700\">")
		if err != nil {
			return err
		}
		var_53 := `Context for the fine`
		_, err = templBuffer.WriteString(var_53)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"context\" class=\"mt-1 px-3 py-2 border border\n 				value={ f.Fine.Context }-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 w-full\" placeholder=\"Context for the fine\"><div class=\"mt-2\"><div hx-get=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/match-list?type=select&matchId=%d", f.Fine.MatchId)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-trigger=\"load once\" hx-target=\"this\"></div>")
		if err != nil {
			return err
		}
		if f.Fine.MatchId == 0 {
			_, err = templBuffer.WriteString("<span class=\"text-sm\">")
			if err != nil {
				return err
			}
			var_54 := `⚠️`
			_, err = templBuffer.WriteString(var_54)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div><!--")
		if err != nil {
			return err
		}
		var_55 := ` Hidden Inputs `
		_, err = templBuffer.WriteString(var_55)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("--><input type=\"playerId\" name=\"playerId\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", f.Player.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"hidden\"><input type=\"fid\" name=\"fid\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", f.Fine.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"hidden\"><input type=\"amount\" name=\"amount\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", f.Fine.Amount)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"hidden\"></div>")
		if err != nil {
			return err
		}
		if isFineMaster {
			_, err = templBuffer.WriteString("<div class=\"px-6 py-4\"><input type=\"number\" name=\"amount\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", f.Fine.Amount)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"px-2 py-1 border rounded w-full\" placeholder=\"Amount\"></div> <div><select name=\"approved\" class=\"px-2 py-1 border rounded\"><option value=\"true\"")
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
			var_56 := `Approved`
			_, err = templBuffer.WriteString(var_56)
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
			var_57 := `Not Approved`
			_, err = templBuffer.WriteString(var_57)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option></select></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<div class=\"mt-10\">")
		if err != nil {
			return err
		}
		var var_58 = []any{bigPri}
		err = templ.RenderCSSItems(ctx, templBuffer, var_58...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_58).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" type=\"submit\">")
		if err != nil {
			return err
		}
		var_59 := `Save`
		_, err = templBuffer.WriteString(var_59)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></div><div class=\"mt-10 flex justify-between\">")
		if err != nil {
			return err
		}
		var var_60 = []any{fmt.Sprintf("%s w-3/5", sec)}
		err = templ.RenderCSSItems(ctx, templBuffer, var_60...)
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
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_60).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_61 := `Cancel`
		_, err = templBuffer.WriteString(var_61)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button>")
		if err != nil {
			return err
		}
		var var_62 = []any{del}
		err = templ.RenderCSSItems(ctx, templBuffer, var_62...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_62).String()))
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
		var_63 := `Delete`
		_, err = templBuffer.WriteString(var_63)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></div></form>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func fineEditRow(f FineWithPlayer, isFineMaster bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_64 := templ.GetChildren(ctx)
		if var_64 == nil {
			var_64 = templ.NopComponent
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
		_, err = templBuffer.WriteString("\" class=\"bg-white divide-y divide-gray-200\" hx-target=\"this\" hx-swap=\"innerHTML\"><td class=\"px-6 py-4\" colspan=\"10\">")
		if err != nil {
			return err
		}
		err = fineEditForm(f, isFineMaster).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</td></tr>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
