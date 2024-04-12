// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "fmt"

func matchSelector(matches []Match, selectedMatchId uint) templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"flex flex-row\"><label><select id=\"matchId\" name=\"matchId\" class=\"transform scale-15 border border-gray-300 rounded-md text-gray-700 flex-grow mb-2\"><option>")
		if err != nil {
			return err
		}
		var_2 := `NA`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option>")
		if err != nil {
			return err
		}
		for _, m := range matches {
			_, err = templBuffer.WriteString("<option value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", m.ID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"")
			if err != nil {
				return err
			}
			if selectedMatchId == m.ID {
				_, err = templBuffer.WriteString(" selected")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			var_3 := `vs `
			_, err = templBuffer.WriteString(var_3)
			if err != nil {
				return err
			}
			var var_4 string = m.Opponent
			_, err = templBuffer.WriteString(templ.EscapeString(var_4))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</select></label></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func matchListPage(matches []Match, isOpen bool) templ.Component {
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
		_, err = templBuffer.WriteString("<html hx-boost=\"true\">")
		if err != nil {
			return err
		}
		err = header().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body class=\"text-2xl\"><div class=\"container mx-auto bg-gray-200 shadow-xl m-10\">")
		if err != nil {
			return err
		}
		if len(matches) > 0 {
			for _, match := range matches {
				_, err = templBuffer.WriteString("<div>")
				if err != nil {
					return err
				}
				var var_6 string = match.Location
				_, err = templBuffer.WriteString(templ.EscapeString(var_6))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div> <div>")
				if err != nil {
					return err
				}
				var var_7 string = match.Opponent
				_, err = templBuffer.WriteString(templ.EscapeString(var_7))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div> <div>")
				if err != nil {
					return err
				}
				var var_8 string = fmt.Sprintf("%+v", match)
				_, err = templBuffer.WriteString(templ.EscapeString(var_8))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
		} else {
			var_9 := `No Matches? Add one to get started`
			_, err = templBuffer.WriteString(var_9)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<div id=\"add-match-container\" hx-get=\"/match?isOpen=true\" hx-trigger=\"click\" hx-swap=\"OuterHTML\" class=\"w-full text-center\">")
		if err != nil {
			return err
		}
		var_10 := `loading add match`
		_, err = templBuffer.WriteString(var_10)
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
		_, err = templBuffer.WriteString("<body class=\"text-2xl\"><div class=\"container mx-auto bg-gray-200 shadow-xl m-10\"><h1 class=\" font-bold mb-4  text-center\">")
		if err != nil {
			return err
		}
		var_12 := `🔨 Baileys Hammer 🔨`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><div id=\"live-match\" class=\"max-w-4xl mx-auto my-8 p-4 border border-gray-200 rounded-lg shadow\"><h2 class=\"text-lg font-semibold mb-4\">")
		if err != nil {
			return err
		}
		var_13 := `Live Match - `
		_, err = templBuffer.WriteString(var_13)
		if err != nil {
			return err
		}
		var var_14 string = data.Match.Location
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
		var var_16 string = data.Match.Opponent
		_, err = templBuffer.WriteString(templ.EscapeString(var_16))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h2><div class=\"mb-4\"><div id=\"match-details\" class=\"mb-4\"></div><h3 class=\"text-md font-semibold mb-2\">")
		if err != nil {
			return err
		}
		var_17 := `Events`
		_, err = templBuffer.WriteString(var_17)
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
			var_18 := `loading events...`
			_, err = templBuffer.WriteString(var_18)
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
		var_19 := templ.GetChildren(ctx)
		if var_19 == nil {
			var_19 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div>")
		if err != nil {
			return err
		}
		var_20 := `[List event: `
		_, err = templBuffer.WriteString(var_20)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var var_21 string = fmt.Sprintf("%v", len(events))
		_, err = templBuffer.WriteString(templ.EscapeString(var_21))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_22 := `]`
		_, err = templBuffer.WriteString(var_22)
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

func createMatch(closeLink templ.SafeURL) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_23 := templ.GetChildren(ctx)
		if var_23 == nil {
			var_23 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"container mx-auto bg-gray-200 shadow-xl m-10\"><form action=\"/match\" method=\"POST\" class=\"max-w-xl mx-auto my-8 p-4 border border-gray-200 rounded-lg shadow\"><h2 class=\"text-lg font-semibold mb-4\">")
		if err != nil {
			return err
		}
		var_24 := `Create New Match`
		_, err = templBuffer.WriteString(var_24)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h2><div class=\"mb-4\"><label for=\"location\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_25 := `Location`
		_, err = templBuffer.WriteString(var_25)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"location\" id=\"location\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"mb-4\"><label for=\"start-time\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_26 := `Start Time`
		_, err = templBuffer.WriteString(var_26)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"datetime-local\" name=\"startTime\" id=\"start-time\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"mb-4\"><label for=\"opponent\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_27 := `Opponent`
		_, err = templBuffer.WriteString(var_27)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"opponent\" id=\"opponent\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"mb-4\"><label for=\"subtitle\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_28 := `Subtitle (optional)`
		_, err = templBuffer.WriteString(var_28)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"subtitle\" id=\"subtitle\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"flex items-center justify-between mt-4\">")
		if err != nil {
			return err
		}
		var var_29 = []any{bigPri}
		err = templ.RenderCSSItems(ctx, templBuffer, var_29...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
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
		var_30 := `Create Match`
		_, err = templBuffer.WriteString(var_30)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></div><div class=\"flex items-center justify-between mt-4\">")
		if err != nil {
			return err
		}
		var var_31 = []any{bigSec}
		err = templ.RenderCSSItems(ctx, templBuffer, var_31...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<a href=\"")
		if err != nil {
			return err
		}
		var var_32 templ.SafeURL = closeLink
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_32)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_31).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_33 := `Close`
		_, err = templBuffer.WriteString(var_33)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a></div></form></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func editMatch(closeLink templ.SafeURL, match Match) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_34 := templ.GetChildren(ctx)
		if var_34 == nil {
			var_34 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"container mx-auto bg-gray-200 shadow-xl m-10\"><form action=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/match/%d", match)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" method=\"POST\" class=\"max-w-xl mx-auto my-8 p-4 border border-gray-200 rounded-lg shadow\"><h2 class=\"text-lg font-semibold mb-4\">")
		if err != nil {
			return err
		}
		var_35 := `Edit Match - vs `
		_, err = templBuffer.WriteString(var_35)
		if err != nil {
			return err
		}
		var var_36 string = match.Opponent
		_, err = templBuffer.WriteString(templ.EscapeString(var_36))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h2><div class=\"mb-4\"><label for=\"location\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_37 := `Location`
		_, err = templBuffer.WriteString(var_37)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"location\" id=\"location\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(match.Location))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></div><div class=\"mb-4\"><label for=\"start-time\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_38 := `Start Time`
		_, err = templBuffer.WriteString(var_38)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"datetime-local\" name=\"startTime\" id=\"start-time\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(match.StartTime.Format("2006-01-02T15:04")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></div><div class=\"mb-4\"><label for=\"opponent\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_39 := `Opponent`
		_, err = templBuffer.WriteString(var_39)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"opponent\" id=\"opponent\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(match.Opponent))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></div><div class=\"mb-4\"><label for=\"subtitle\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_40 := `Subtitle (optional)`
		_, err = templBuffer.WriteString(var_40)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"subtitle\" id=\"subtitle\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(match.Subtitle))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></div><div class=\"flex items-center justify-between mt-4\">")
		if err != nil {
			return err
		}
		var var_41 = []any{bigPri}
		err = templ.RenderCSSItems(ctx, templBuffer, var_41...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_41).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_42 := `Edit Match`
		_, err = templBuffer.WriteString(var_42)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></div><div class=\"flex items-center justify-between mt-4\">")
		if err != nil {
			return err
		}
		var var_43 = []any{bigSec}
		err = templ.RenderCSSItems(ctx, templBuffer, var_43...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<a href=\"")
		if err != nil {
			return err
		}
		var var_44 templ.SafeURL = closeLink
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_44)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_43).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_45 := `Close`
		_, err = templBuffer.WriteString(var_45)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a></div></form></div>")
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
		var_46 := templ.GetChildren(ctx)
		if var_46 == nil {
			var_46 = templ.NopComponent
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
		var_47 := `[createNewEvent`
		_, err = templBuffer.WriteString(var_47)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_48 := `: `
		_, err = templBuffer.WriteString(var_48)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var var_49 string = fmt.Sprintf("%v", matchId)
		_, err = templBuffer.WriteString(templ.EscapeString(var_49))
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString(" ")
		if err != nil {
			return err
		}
		var_50 := `]`
		_, err = templBuffer.WriteString(var_50)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString(" <h2 class=\"text-lg font-semibold mb-4\">")
		if err != nil {
			return err
		}
		var_51 := `Create New Event`
		_, err = templBuffer.WriteString(var_51)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</h2><div class=\"mb-4\"><label for=\"event-name\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_52 := `Event Name`
		_, err = templBuffer.WriteString(var_52)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"eventName\" id=\"event-name\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></div><div class=\"mb-4\"><label for=\"event-type\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_53 := `Event Type`
		_, err = templBuffer.WriteString(var_53)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</label><select name=\"eventType\" id=\"event-type\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"><option value=\"subbed-off\">")
		if err != nil {
			return err
		}
		var_54 := `Subbed Off`
		_, err = templBuffer.WriteString(var_54)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option><option value=\"subbed-on\">")
		if err != nil {
			return err
		}
		var_55 := `Subbed On`
		_, err = templBuffer.WriteString(var_55)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option><option value=\"goal\">")
		if err != nil {
			return err
		}
		var_56 := `Goal`
		_, err = templBuffer.WriteString(var_56)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option><option value=\"assist\">")
		if err != nil {
			return err
		}
		var_57 := `Assist`
		_, err = templBuffer.WriteString(var_57)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option><option value=\"own-goal\">")
		if err != nil {
			return err
		}
		var_58 := `Own Goal`
		_, err = templBuffer.WriteString(var_58)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</option></select></div><div class=\"mb-4\"><label for=\"event-time-offset\" class=\"block text-sm font-medium text-gray-700\">")
		if err != nil {
			return err
		}
		var_59 := `Event Time Offset`
		_, err = templBuffer.WriteString(var_59)
		if err != nil {
			return err
		}
				_, err = templBuffer.WriteString("</label><select name=\"eventTimeOffset\" id=\"event-time-offset\" class=\"mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500\"></select></div><button type=\"submit\" class=\"inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500\">")
		if err != nil {
			return err
		}
		var_60 := `Create Event`
		_, err = templBuffer.WriteString(var_60)
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
