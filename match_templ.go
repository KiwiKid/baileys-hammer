// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "fmt"

func matchListPage(matches []Match, isOpen bool) templ.Component {
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
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_2 := `[matchListPage:`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var var_3 string = fmt.Sprintf("%+v", matches)
		_, err = templBuffer.WriteString(templ.EscapeString(var_3))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_4 := `]`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <body class=\"text-2xl\"><div class=\"container mx-auto bg-gray-200 shadow-xl m-10\">")
		if err != nil {
			return err
		}
		if len(matches) > 0 {
			for _, match := range matches {
				_, err = templBuffer.WriteString("<div>")
				if err != nil {
					return err
				}
				var var_5 string = match.Location
				_, err = templBuffer.WriteString(templ.EscapeString(var_5))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div> <div>")
				if err != nil {
					return err
				}
				var var_6 string = match.Opponent
				_, err = templBuffer.WriteString(templ.EscapeString(var_6))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div> <div>")
				if err != nil {
					return err
				}
				var var_7 string = fmt.Sprintf("%+v", match)
				_, err = templBuffer.WriteString(templ.EscapeString(var_7))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
		} else {
			var_8 := `No Matches? Add one to get started`
			_, err = templBuffer.WriteString(var_8)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<div id=\"add-match-container\"")
		if err != nil {
			return err
		}
		if isOpen {
			_, err = templBuffer.WriteString(" hx-get=\"/match?isOpen=true\"")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString(" hx-get=\"/match?isOpen=false\"")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(" hx-trigger=\"load once\" hx-swap=\"OuterHTML\" class=\"w-full text-center\">")
		if err != nil {
			return err
		}
		var_9 := `loading add match`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div></div>")
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

func matchPage(data MatchPageData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_10 := templ.GetChildren(ctx)
		if var_10 == nil {
			var_10 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<html hx-boost=\"true\">")
		if err != nil {
			return err
		}
		var_11 := `[matchPage`
		_, err = templBuffer.WriteString(var_11)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_12 := `: `
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var var_13 string = fmt.Sprintf("%v", data)
		_, err = templBuffer.WriteString(templ.EscapeString(var_13))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_14 := `]`
		_, err = templBuffer.WriteString(var_14)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		err = header().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body class=\"text-2xl\"><div class=\"container mx-auto bg-gray-200 shadow-xl m-10\"><h1 class=\" font-bold mb-4  text-center\">")
		if err != nil {
			return err
		}
		var_15 := `🔨 Baileys Hammer 🔨`
		_, err = templBuffer.WriteString(var_15)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><div id=\"live-match\" class=\"max-w-4xl mx-auto my-8 p-4 border border-gray-200 rounded-lg shadow\"><h2 class=\"text-lg font-semibold mb-4\">")
		if err != nil {
			return err
		}
		var_16 := `Live Match - `
		_, err = templBuffer.WriteString(var_16)
		if err != nil {
			return err
		}
		var var_17 string = data.Match.Location
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
		var var_19 string = data.Match.Opponent
		_, err = templBuffer.WriteString(templ.EscapeString(var_19))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h2><div class=\"mb-4\"><div id=\"match-details\" class=\"mb-4\"></div><h3 class=\"text-md font-semibold mb-2\">")
		if err != nil {
			return err
		}
		var_20 := `Events`
		_, err = templBuffer.WriteString(var_20)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h3>")
		if err != nil {
			return err
		}
		if data.Match.ID > 0 {
			_, err = templBuffer.WriteString("<div class=\"list-disc list-inside\"")
			if err != nil {
				return err
			}
			if data.isOpen {
				_, err = templBuffer.WriteString(" hx-get=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/match/%d/event?isOpen=true", data.Match.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
			} else {
				_, err = templBuffer.WriteString(" hx-get=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/match/%d/event?isOpen=false", data.Match.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString(" hx-trigger=\"load once\"></div> <div hx-get=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/match/%d/events", data.Match.ID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-trigger=\"load once\">")
			if err != nil {
				return err
			}
			var_21 := `loading events...`
			_, err = templBuffer.WriteString(var_21)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div></div></div>")
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

func listEvents(events []MatchEvent) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_22 := templ.GetChildren(ctx)
		if var_22 == nil {
			var_22 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div>")
		if err != nil {
			return err
		}
		var_23 := `[List event: `
		_, err = templBuffer.WriteString(var_23)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var var_24 string = fmt.Sprintf("%v", len(events))
		_, err = templBuffer.WriteString(templ.EscapeString(var_24))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_25 := `]`
		_, err = templBuffer.WriteString(var_25)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func createMatch(isOpen bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_26 := templ.GetChildren(ctx)
		if var_26 == nil {
			var_26 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div>")
		if err != nil {
			return err
		}
		var_27 := `[createMatch`
		_, err = templBuffer.WriteString(var_27)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_28 := `: `
		_, err = templBuffer.WriteString(var_28)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var var_29 string = fmt.Sprintf("%v", isOpen)
		_, err = templBuffer.WriteString(templ.EscapeString(var_29))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_30 := `]`
		_, err = templBuffer.WriteString(var_30)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		if isOpen {
			_, err = templBuffer.WriteString("<div class=\"container mx-auto bg-gray-200 shadow-xl m-10\"><form action=\"/match\" method=\"POST\" class=\"max-w-xl mx-auto my-8 p-4 border border-gray-200 rounded-lg shadow\"><h2 class=\"text-lg font-semibold mb-4\">")
			if err != nil {
				return err
			}
			var_31 := `Create New Match`
			_, err = templBuffer.WriteString(var_31)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2><div class=\"mb-4\"><label for=\"location\" class=\"block text-sm font-medium text-gray-700\">")
			if err != nil {
				return err
			}
			var_32 := `Location`
			_, err = templBuffer.WriteString(var_32)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"location\" id=\"location\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"mb-4\"><label for=\"start-time\" class=\"block text-sm font-medium text-gray-700\">")
			if err != nil {
				return err
			}
			var_33 := `Start Time`
			_, err = templBuffer.WriteString(var_33)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"datetime-local\" name=\"startTime\" id=\"start-time\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"mb-4\"><label for=\"opponent\" class=\"block text-sm font-medium text-gray-700\">")
			if err != nil {
				return err
			}
			var_34 := `Opponent`
			_, err = templBuffer.WriteString(var_34)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"opponent\" id=\"opponent\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"mb-4\"><label for=\"subtitle\" class=\"block text-sm font-medium text-gray-700\">")
			if err != nil {
				return err
			}
			var_35 := `Subtitle (optional)`
			_, err = templBuffer.WriteString(var_35)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"subtitle\" id=\"subtitle\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"flex items-center justify-between mt-4\">")
			if err != nil {
				return err
			}
			var var_36 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_36...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
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
			var_37 := `Create Match`
			_, err = templBuffer.WriteString(var_37)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div><div class=\"flex items-center justify-between mt-4\">")
			if err != nil {
				return err
			}
			var var_38 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_38...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button hx-get=\"/match-list?isOpen=false\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_38).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_39 := `Close`
			_, err = templBuffer.WriteString(var_39)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div></form></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"flex justify-center w-full\">")
			if err != nil {
				return err
			}
			var var_40 = []any{bigAdd}
			err = templ.RenderCSSItems(ctx, templBuffer, var_40...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_41 templ.SafeURL = templ.SafeURL("/match-list?isOpen=true")
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_41)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_40).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_42 := `Add Match`
			_, err = templBuffer.WriteString(var_42)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

/*

func createNewEvent(matchId uint64) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_43 := templ.GetChildren(ctx)
		if var_43 == nil {
			var_43 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
				_, err = templBuffer.WriteString("<form action=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/match/%d/event", matchId)))
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("\" method=\"POST\" class=\"max-w-xl mx-auto my-8 p-4 border border-gray-200 rounded-lg shadow\">")
		if err != nil {
			return err
		}
		var_44 := `[createNewEvent`
		_, err = templBuffer.WriteString(var_44)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_45 := `: `
		_, err = templBuffer.WriteString(var_45)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var var_46 string = fmt.Sprintf("%v", matchId)
		_, err = templBuffer.WriteString(templ.EscapeString(var_46))
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_47 := `]`
		_, err = templBuffer.WriteString(var_47)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString(" <h2 class=\"text-lg font-semibold mb-4\">")
		if err != nil {
			return err
		}
		var_48 := `Create New Event`
		_, err = templBuffer.WriteString(var_48)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</h2><div class=\"mb-4\"><label for=\"event-name\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_49 := `Event Name`
		_, err = templBuffer.WriteString(var_49)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"eventName\" id=\"event-name\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"mb-4\"><label for=\"event-type\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_50 := `Event Type`
		_, err = templBuffer.WriteString(var_50)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</label><select name=\"eventType\" id=\"event-type\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"><option value=\"subbed-off\">")
		if err != nil {
			return err
		}
		var_51 := `Subbed Off`
		_, err = templBuffer.WriteString(var_51)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option><option value=\"subbed-on\">")
		if err != nil {
			return err
		}
		var_52 := `Subbed On`
		_, err = templBuffer.WriteString(var_52)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option><option value=\"goal\">")
		if err != nil {
			return err
		}
		var_53 := `Goal`
		_, err = templBuffer.WriteString(var_53)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option><option value=\"assist\">")
		if err != nil {
			return err
		}
		var_54 := `Assist`
		_, err = templBuffer.WriteString(var_54)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option><option value=\"own-goal\">")
		if err != nil {
			return err
		}
		var_55 := `Own Goal`
		_, err = templBuffer.WriteString(var_55)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option></select></div><div class=\"mb-4\"><label for=\"event-time-offset\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_56 := `Event Time Offset`
		_, err = templBuffer.WriteString(var_56)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</label><select name=\"eventTimeOffset\" id=\"event-time-offset\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></select></div><button type=\"submit\" class=\"inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500\">")
		if err != nil {
			return err
		}
		var_57 := `Create Event`
		_, err = templBuffer.WriteString(var_57)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</button></form>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

*/
