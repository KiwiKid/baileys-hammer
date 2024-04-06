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
	"log"
)

var baseUrl = "/"

func makeSafeUrlWithAnchor(baseUrl string, finesOpen bool, playersOpen bool, presetFinesOpen bool, anchorTag string) templ.SafeURL {
	url := makeUrl(baseUrl, finesOpen, playersOpen, presetFinesOpen) + "#" + anchorTag
	return templ.SafeURL(url)
}

func makeSafeUrl(baseUrl string, finesOpen bool, playersOpen bool, presetFinesOpen bool) templ.SafeURL {
	url := makeUrl(baseUrl, finesOpen, playersOpen, presetFinesOpen)
	return templ.SafeURL(url)
}

func makeUrl(fbaseUrl string, finesOpen bool, playersOpen bool, presetFinesOpen bool) string {

	hp := HomeQueryParams{
		FinesOpen:       finesOpen,
		PlayerOpen:      playersOpen,
		PresetFinesOpen: presetFinesOpen,
	}

	url, err := GenerateUrl(fbaseUrl, &hp)
	if err != nil {
		log.Fatalf("Generate url error: %+v", err)
	}
	return *url
}

func downArrow() templ.Component {
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
		_, err = templBuffer.WriteString("<svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-5 w-5\" viewBox=\"0 0 20 20\" fill=\"currentColor\"><path fill-rule=\"evenodd\" d=\"M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z\" clip-rule=\"evenodd\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

var pri = "bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"

var sec = "bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"

var add = "bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"
var del = "bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"

var bigBtnTxt = "mx-auto items-center justify-center w-4/5 text-center py-2 px-4 text-lg rounded-md border hover:bg-opacity-75 focus:outline-none"
var bigPri = fmt.Sprintf("%s %s", bigBtnTxt, pri)
var bigSec = fmt.Sprintf("%s %s", bigBtnTxt, sec)
var bigAdd = fmt.Sprintf("%s %s", bigBtnTxt, add)
var bigDel = fmt.Sprintf("%s %s", bigBtnTxt, del)

var S = fmt.Sprint
var F = fmt.Sprintf

func header() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_2 := templ.GetChildren(ctx)
		if var_2 == nil {
			var_2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<head><title>")
		if err != nil {
			return err
		}
		var_3 := `🔨 Baileys Hammer 🔨`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</title><link href=\"https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css\" rel=\"stylesheet\"><script src=\"https://unpkg.com/hyperscript.org@0.9.12\">")
		if err != nil {
			return err
		}
		var_4 := ``
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></head>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func pageFooter() templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"bg-yellow-400 border-l-4 border-yellow-800 text-yellow-800 p-2 mt-10 hidden\" role=\"alert\"><p class=\"font-bold\"><details><summary class=\"text-center\">")
		if err != nil {
			return err
		}
		var_6 := `🚧 Under Construction (Version `
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		var var_7 string = `0.2`
		_, err = templBuffer.WriteString(templ.EscapeString(var_7))
		if err != nil {
			return err
		}
		var_8 := `) 🚧`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</summary><div><h1>")
		if err != nil {
			return err
		}
		var_9 := `todo:`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><ul><li>")
		if err != nil {
			return err
		}
		var_10 := `- Add Fine Status: `
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li><li>")
		if err != nil {
			return err
		}
		var_11 := `--- "Paid"`
		_, err = templBuffer.WriteString(var_11)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li><li>")
		if err != nil {
			return err
		}
		var_12 := `--- "Declined"`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li><li>")
		if err != nil {
			return err
		}
		var_13 := `- Add option to delete player`
		_, err = templBuffer.WriteString(var_13)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li></ul></div></details></p><div id=\"any-errors\"></div><script src=\"https://unpkg.com/htmx.org\">")
		if err != nil {
			return err
		}
		var_14 := ``
		_, err = templBuffer.WriteString(var_14)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/htmx.org@1.9.11/dist/ext/response-targets.js\">")
		if err != nil {
			return err
		}
		var_15 := ``
		_, err = templBuffer.WriteString(var_15)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func home(players []PlayerWithFines, approvedPFines []PresetFine, pendingPFines []PresetFine, qp HomeQueryParams) templ.Component {
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
		_, err = templBuffer.WriteString("<html hx-boost=\"true\">")
		if err != nil {
			return err
		}
		err = header().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body class=\"text-2xl\"><div class=\"container mx-auto bg-gray-200 shadow-xl\"><h1 class=\" font-bold mb-4  text-center\">")
		if err != nil {
			return err
		}
		var_17 := `🔨 Baileys Hammer 🔨`
		_, err = templBuffer.WriteString(var_17)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1>")
		if err != nil {
			return err
		}
		if len(approvedPFines) > 0 {
			_, err = templBuffer.WriteString("<div class=\"bg-sepia-200 shadow-xl m-10 rounded-lg\"><h3 class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			var_18 := `Fines`
			_, err = templBuffer.WriteString(var_18)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3><ul class=\"list-inside space-y-3 text-lg font-handwriting text-brown-900\">")
			if err != nil {
				return err
			}
			for _, pf := range approvedPFines {
				_, err = templBuffer.WriteString("<li class=\"pl-6 border-l-4 border-gold-700 hover:bg-sepia-300 transition duration-300 ease-in-out\">")
				if err != nil {
					return err
				}
				var var_19 string = fmt.Sprintf("$%v - %s", pf.Amount, pf.Reason)
				_, err = templBuffer.WriteString(templ.EscapeString(var_19))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</li>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</ul></div>")
			if err != nil {
				return err
			}
		}
		if len(pendingPFines) > 0 {
			_, err = templBuffer.WriteString("<div class=\"w-full flex justify-center items-center\"><div _=\"on click toggle .hidden on next &lt;section/&gt;\" class=\"flex justify-center items-center cursor-pointer\">")
			if err != nil {
				return err
			}
			var var_20 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_20...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<h3 class=\"")
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
			var_21 := `Pending Fines `
			_, err = templBuffer.WriteString(var_21)
			if err != nil {
				return err
			}
			var var_22 string = fmt.Sprintf("(%d)", len(pendingPFines))
			_, err = templBuffer.WriteString(templ.EscapeString(var_22))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3></div></div> <section class=\"bg-sepia-200 shadow-xl m-10 rounded-lg hidden\"><ul class=\"list-inside space-y-3 text-lg font-handwriting text-brown-900\">")
			if err != nil {
				return err
			}
			for _, pf := range pendingPFines {
				_, err = templBuffer.WriteString("<li class=\"pl-6 border-l-4 border-gold-700 hover:bg-sepia-300 transition duration-300 ease-in-out\">")
				if err != nil {
					return err
				}
				var var_23 string = fmt.Sprintf("$%v - %s", pf.Amount, pf.Reason)
				_, err = templBuffer.WriteString(templ.EscapeString(var_23))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" <span class=\"ml-2 inline-block bg-red-100 text-red-800 font-bold px-2 py-1 rounded-full text-sm shadow-sm\">")
				if err != nil {
					return err
				}
				var_24 := `(pending approval)`
				_, err = templBuffer.WriteString(var_24)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span></li>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</ul></section>")
			if err != nil {
				return err
			}
		}
		err = fineAdd(baseUrl, qp.FinesOpen, players, approvedPFines, false).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div class=\"flex  bg-gray-100 mx-auto shadow-xl m-10\"><div class=\"w-full mt-10\"><h3 class=\"text-2xl font-bold text-center\">")
		if err != nil {
			return err
		}
		var_25 := `Leaderboard`
		_, err = templBuffer.WriteString(var_25)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h3><ul>")
		if err != nil {
			return err
		}
		for _, p := range players {
			_, err = templBuffer.WriteString("<li class=\"m-4\">")
			if err != nil {
				return err
			}
			var var_26 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_26...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<div _=\"on click toggle .hidden on next &lt;section/&gt;\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_26).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var var_27 string = p.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_27))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_28 := `- `
			_, err = templBuffer.WriteString(var_28)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_29 string = fmt.Sprintf("$%d (%d)", p.TotalFines, p.TotalFineCount)
			_, err = templBuffer.WriteString(templ.EscapeString(var_29))
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
				var var_30 string = f.Reason
				_, err = templBuffer.WriteString(templ.EscapeString(var_30))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_31 string = fmt.Sprintf("$%.0f - %s", f.Amount, humanize.Time(f.FineAt))
				_, err = templBuffer.WriteString(templ.EscapeString(var_31))
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
		_, err = templBuffer.WriteString("</ul></div></div></div><div id=\"fine-list-container\" hx-get=\"/fines\" hx-target=\"#fine-list-container\" hx-swap=\"outerHTML\" hx-trigger=\"load once\" class=\"w-full text-center\">")
		if err != nil {
			return err
		}
		var_32 := `loading latest..`
		_, err = templBuffer.WriteString(var_32)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><script src=\"https://unpkg.com/htmx.org\">")
		if err != nil {
			return err
		}
		var_33 := ``
		_, err = templBuffer.WriteString(var_33)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script>")
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

func fineAdd(baseUrl string, isOpen bool, players []PlayerWithFines, presetFines []PresetFine, isFineMaster bool) templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"container mx-auto bg-gray-200 shadow-xl m-10\">")
		if err != nil {
			return err
		}
		if isOpen {
			_, err = templBuffer.WriteString("<div class=\"px-8 py-6 text-left bg-gray-200 m-10\" id=\"fine-add\"><h3 class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			if isFineMaster {
				var_35 := `Fine a Player:`
				_, err = templBuffer.WriteString(var_35)
				if err != nil {
					return err
				}
			} else {
				var_36 := `Submit a Fine`
				_, err = templBuffer.WriteString(var_36)
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</h3><form hx-post=\"/fines\" class=\"mt-4\"><div class=\"mt-4\"><div class=\"border-t pt-4\"><label class=\"text-lg font-semibold\"><div class=\"mt-2\"><select id=\"presetFineId\" name=\"presetFineId\" class=\" bg-white  w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50\"><option selected value=\"\">")
			if err != nil {
				return err
			}
			var_37 := `-- Select Fine --`
			_, err = templBuffer.WriteString(var_37)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option><option value=\"-1\">")
			if err != nil {
				return err
			}
			var_38 := `-- Fine is not listed here --`
			_, err = templBuffer.WriteString(var_38)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option>")
			if err != nil {
				return err
			}
			for _, fp := range presetFines {
				if fp.Approved {
					_, err = templBuffer.WriteString("<option value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", fp.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var var_39 string = fmt.Sprintf("%s ($%v)", fp.Reason, fp.Amount)
					_, err = templBuffer.WriteString(templ.EscapeString(var_39))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</option>")
					if err != nil {
						return err
					}
				}
			}
			_, err = templBuffer.WriteString("</select></div></label><section class=\"hidden\" id=\"newFine\"><div class=\"pt-4 border border-2 rounded\"><label class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			var_40 := `New Fine`
			_, err = templBuffer.WriteString(var_40)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div class=\"mt-4\"><label for=\"reason\" class=\"block\">")
			if err != nil {
				return err
			}
			var_41 := `Title`
			_, err = templBuffer.WriteString(var_41)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"reason\" id=\"reason\" placeholder=\"Reason for the fine\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500\"></div><div class=\"mt-4\"><label for=\"context\" class=\"block\">")
			if err != nil {
				return err
			}
			var_42 := `Details:`
			_, err = templBuffer.WriteString(var_42)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"context\" id=\"context\" placeholder=\"More Context for the fine\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500\"></div>")
			if err != nil {
				return err
			}
			if isFineMaster {
				_, err = templBuffer.WriteString("<div class=\"mt-4\"><label for=\"amount\" class=\"block\">")
				if err != nil {
					return err
				}
				var_43 := `Amount ($)`
				_, err = templBuffer.WriteString(var_43)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"amount\" id=\"amount\" placeholder=\"Amount\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500\"></div> <label><input type=\"hidden\" name=\"approved\" value=\"on\"></label>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("<div class=\"grid grid-cols-2 gap-4\"><div class=\"flex items-center justify-center p-4 border border-gray-200 rounded-lg\"><input type=\"radio\" id=\"oneOffFine\" name=\"fineOption\" value=\"oneOffFine\" class=\"form-radio text-blue-600 transform scale-15\"><label for=\"oneOffFine\" class=\"ml-2 text-gray-800\">")
			if err != nil {
				return err
			}
			var_44 := `One Off Fine`
			_, err = templBuffer.WriteString(var_44)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label></div><div class=\"flex items-center justify-center p-4 border border-gray-200 rounded-lg\"><input type=\"radio\" id=\"applyAgain\" name=\"fineOption\" value=\"applyAgain\" class=\"form-radio text-blue-600 transform scale-15\"><label for=\"applyAgain\" class=\"ml-2 text-gray-800\">")
			if err != nil {
				return err
			}
			var_45 := `Could Apply Again`
			_, err = templBuffer.WriteString(var_45)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label></div></div></div></section><div class=\"mt-4\"><label>")
			if err != nil {
				return err
			}
			var_46 := `Who:`
			_, err = templBuffer.WriteString(var_46)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><select name=\"playerId\" class=\"mt-1 w-full border-gray-300  bg-white rounded-md shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50\"><option selected value=\"\">")
			if err != nil {
				return err
			}
			var_47 := `-- no player - just suggest fineable offence --`
			_, err = templBuffer.WriteString(var_47)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option>")
			if err != nil {
				return err
			}
			for _, p := range players {
				_, err = templBuffer.WriteString("<option value=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", p.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var var_48 string = fmt.Sprintf("%s", p.Name)
				_, err = templBuffer.WriteString(templ.EscapeString(var_48))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</select></div></div></div><div class=\"grid grid-cols-1 p-4 gap-4 mt-10\">")
			if err != nil {
				return err
			}
			var var_49 = []any{bigAdd}
			err = templ.RenderCSSItems(ctx, templBuffer, var_49...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_49).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_50 := `Add Fine`
			_, err = templBuffer.WriteString(var_50)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button>")
			if err != nil {
				return err
			}
			var var_51 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_51...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_52 templ.SafeURL = makeSafeUrlWithAnchor(baseUrl, false, false, false, "fine-add")
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_52)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"")
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
			var_53 := `Close`
			_, err = templBuffer.WriteString(var_53)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></form></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"flex justify-center w-full p-4\" id=\"fine-add\">")
			if err != nil {
				return err
			}
			var var_54 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_54...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_54).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" href=\"")
			if err != nil {
				return err
			}
			var var_55 templ.SafeURL = makeSafeUrlWithAnchor(baseUrl, true, false, false, "fine-add")
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_55)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			if isFineMaster {
				var_56 := `Fine a Player`
				_, err = templBuffer.WriteString(var_56)
				if err != nil {
					return err
				}
			} else {
				var_57 := `Suggest a Fine`
				_, err = templBuffer.WriteString(var_57)
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</a></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<script>")
		if err != nil {
			return err
		}
		var_58 := `
		window.fpSelect = document.getElementById('presetFineId')
		if(window.fpSelect != null){
			fpSelect.addEventListener('change', function() {
				const section = document.getElementById('newFine');
				if (this.value == '-1') { // Change '2' to the value of the option that should show the section
					section.classList.remove('hidden');
				} else {
					section.classList.add('hidden');
				}
			});
		}else{
			console.warn('no fpSelect')
		}
	`
		_, err = templBuffer.WriteString(var_58)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
