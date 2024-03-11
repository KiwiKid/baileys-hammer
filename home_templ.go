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

var pri = "bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"

var sec = "bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"

var add = "bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"

var bigBtnTxt = "mx-auto items-center justify-center w-4/5 text-center py-2 px-4 text-lg rounded-md border hover:bg-opacity-75 focus:outline-none"
var bigPri = fmt.Sprintf("%s %s", bigBtnTxt, pri)
var bigSec = fmt.Sprintf("%s %s", bigBtnTxt, sec)

var bigAdd = fmt.Sprintf("%s %s", bigBtnTxt, add)

func header() templ.Component {
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
		_, err = templBuffer.WriteString("<head><title>")
		if err != nil {
			return err
		}
		var_2 := `🔨 Baileys Hammer 🔨`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</title><link href=\"https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css\" rel=\"stylesheet\"><script src=\"https://unpkg.com/hyperscript.org@0.9.12\">")
		if err != nil {
			return err
		}
		var_3 := ``
		_, err = templBuffer.WriteString(var_3)
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

func pageHeader() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_4 := templ.GetChildren(ctx)
		if var_4 == nil {
			var_4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"bg-yellow-400 border-l-4 border-yellow-800 text-yellow-800 p-2 \" role=\"alert\"><p class=\"font-bold\"><details><summary class=\"text-center\">")
		if err != nil {
			return err
		}
		var_5 := `🚧 Under Construction (Version `
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		var var_6 string = `0.1`
		_, err = templBuffer.WriteString(templ.EscapeString(var_6))
		if err != nil {
			return err
		}
		var_7 := `) 🚧`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</summary><div><h1>")
		if err != nil {
			return err
		}
		var_8 := `todo:`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><ul><li>")
		if err != nil {
			return err
		}
		var_9 := `- Add ability to contest fine`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li><li>")
		if err != nil {
			return err
		}
		var_10 := `- Add Fine "Paid" status`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li><li>")
		if err != nil {
			return err
		}
		var_11 := `- Move setting of amount to during finemaster approval`
		_, err = templBuffer.WriteString(var_11)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li><li>")
		if err != nil {
			return err
		}
		var_12 := `- Add delete player`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li></ul></div></details></p></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func home(players []PlayerWithFines, pFines []PresetFine, qp HomeQueryParams) templ.Component {
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
		_, err = templBuffer.WriteString("<html hx-boost=\"true\">")
		if err != nil {
			return err
		}
		err = header().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body class=\"text-2xl\">")
		if err != nil {
			return err
		}
		err = pageHeader().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div class=\"container mx-auto bg-gray-200 shadow-xl m-10\"><h1 class=\" font-bold mb-4  text-center\">")
		if err != nil {
			return err
		}
		var_14 := `🔨 Baileys Hammer 🔨`
		_, err = templBuffer.WriteString(var_14)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><div class=\"bg-sepia-200 shadow-xl m-10 rounded-lg\"><h3 class=\"text-2xl font-bold text-center\">")
		if err != nil {
			return err
		}
		var_15 := `Fines`
		_, err = templBuffer.WriteString(var_15)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h3><ul class=\"list-inside space-y-3 text-lg font-handwriting text-brown-900\">")
		if err != nil {
			return err
		}
		for _, pf := range pFines {
			_, err = templBuffer.WriteString("<li class=\"pl-6 border-l-4 border-gold-700 hover:bg-sepia-300 transition duration-300 ease-in-out\">")
			if err != nil {
				return err
			}
			var var_16 string = fmt.Sprintf("$%v - %s", pf.Amount, pf.Reason)
			_, err = templBuffer.WriteString(templ.EscapeString(var_16))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			if !pf.Approved {
				_, err = templBuffer.WriteString("<span class=\"ml-2 inline-block bg-red-100 text-red-800 font-bold px-2 py-1 rounded-full text-sm shadow-sm\">")
				if err != nil {
					return err
				}
				var_17 := `(pending approval)`
				_, err = templBuffer.WriteString(var_17)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</li>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</ul></div><div class=\"flex  bg-gray-100 mx-auto shadow-xl m-10\"><ul>")
		if err != nil {
			return err
		}
		for _, p := range players {
			_, err = templBuffer.WriteString("<li class=\"m-4\"><div _=\"on click toggle .hidden on next &lt;section/&gt;\" class=\"cursor-pointer p-2 bg-gray-200 rounded hover:bg-gray-300 transform hover:-translate-y-1 hover:scale-10\n								5\">")
			if err != nil {
				return err
			}
			var var_18 string = p.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_18))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_19 := `- `
			_, err = templBuffer.WriteString(var_19)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var var_20 string = fmt.Sprintf("$%d (%d)", p.TotalFines, p.TotalFineCount)
			_, err = templBuffer.WriteString(templ.EscapeString(var_20))
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
				var var_21 string = f.Reason
				_, err = templBuffer.WriteString(templ.EscapeString(var_21))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_22 string = fmt.Sprintf("$%.0f - %s", f.Amount, humanize.Time(f.CreatedAt))
				_, err = templBuffer.WriteString(templ.EscapeString(var_22))
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
		err = fineAdd(baseUrl, qp.FinesOpen, players, pFines, false).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><div hx-get=\"/fines\" hx-trigger=\"load once\" hx-swap=\"OuterHTML\" class=\"w-full text-center\">")
		if err != nil {
			return err
		}
		var_23 := `loading latest..`
		_, err = templBuffer.WriteString(var_23)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><script src=\"https://unpkg.com/htmx.org\">")
		if err != nil {
			return err
		}
		var_24 := ``
		_, err = templBuffer.WriteString(var_24)
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

func fineAdd(baseUrl string, isOpen bool, players []PlayerWithFines, presetFines []PresetFine, isFineMaster bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_25 := templ.GetChildren(ctx)
		if var_25 == nil {
			var_25 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if isOpen {
			_, err = templBuffer.WriteString("<div class=\"px-8 py-6 text-left bg-gray-200 shadow-xl m-10\" id=\"fine-add\"><h3 class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			if isFineMaster {
				var_26 := `Fine a Player:`
				_, err = templBuffer.WriteString(var_26)
				if err != nil {
					return err
				}
			} else {
				var_27 := `Submit a Fine`
				_, err = templBuffer.WriteString(var_27)
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</h3><form hx-post=\"/fines\" class=\"mt-4\"><div class=\"mt-4\"><label class=\"text-lg font-semibold\">")
			if err != nil {
				return err
			}
			var_28 := `Select Player:`
			_, err = templBuffer.WriteString(var_28)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><select name=\"playerId\" class=\"mt-1 w-full border-gray-300  bg-white rounded-md shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50\">")
			if err != nil {
				return err
			}
			for _, p := range players {
				_, err = templBuffer.WriteString("<option value=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v", p.PlayerID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var var_29 string = fmt.Sprintf("%s", p.Name)
				_, err = templBuffer.WriteString(templ.EscapeString(var_29))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</select></div><div class=\"grid grid-cols-1 gap-4\"><div class=\"border-t pt-4\"><label class=\"text-lg font-semibold\">")
			if err != nil {
				return err
			}
			var_30 := `Select Fine`
			_, err = templBuffer.WriteString(var_30)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div class=\"mt-2\"><select id=\"presetFineId\" name=\"presetFineId\" class=\" bg-white  w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50\"><option selected value=\"\">")
			if err != nil {
				return err
			}
			var_31 := `-- Select Fine --`
			_, err = templBuffer.WriteString(var_31)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option>")
			if err != nil {
				return err
			}
			for _, fp := range presetFines {
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
				var var_32 string = fmt.Sprintf("%s ($%v)", fp.Reason, fp.Amount)
				_, err = templBuffer.WriteString(templ.EscapeString(var_32))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("<option value=\"-1\">")
			if err != nil {
				return err
			}
			var_33 := `-- Create New --`
			_, err = templBuffer.WriteString(var_33)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option></select></div><section class=\"hidden\" id=\"newFine\"><div class=\"border-t pt-4\"><label class=\"text-2xl font-bold text-center\">")
			if err != nil {
				return err
			}
			var_34 := `New Fine`
			_, err = templBuffer.WriteString(var_34)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div class=\"mt-4\"><label for=\"reason\" class=\"block\">")
			if err != nil {
				return err
			}
			var_35 := `Reason`
			_, err = templBuffer.WriteString(var_35)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"reason\" id=\"reason\" placeholder=\"Reason for the fine\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500\"></div><div class=\"mt-4\"><label for=\"amount\" class=\"block\">")
			if err != nil {
				return err
			}
			var_36 := `Amount ($)`
			_, err = templBuffer.WriteString(var_36)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"amount\" id=\"amount\" placeholder=\"Amount\" class=\"w-full px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500\"></div>")
			if err != nil {
				return err
			}
			if isFineMaster {
				_, err = templBuffer.WriteString("<label><input type=\"hidden\" name=\"approved\" value=\"on\"></label>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("<div class=\"mt-4\"><label class=\"block text-lg\"><input type=\"checkbox\" name=\"should_be_a_rule\" class=\"form-checkbox h-6 w-6 mt-2 text-blue-600 border-gray-300 rounded focus:ring-blue-500\">")
			if err != nil {
				return err
			}
			if isFineMaster {
				var_37 := `Make this a preset fine (This fine could apply again )`
				_, err = templBuffer.WriteString(var_37)
				if err != nil {
					return err
				}
			} else {
				var_38 := `This fine could apply again `
				_, err = templBuffer.WriteString(var_38)
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</label></div></div></section></div></div><div class=\"grid grid-cols-1 p-4 gap-4 mt-10\">")
			if err != nil {
				return err
			}
			var var_39 = []any{bigPri}
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
			var_40 := `Fine a Player`
			_, err = templBuffer.WriteString(var_40)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button>")
			if err != nil {
				return err
			}
			var var_41 = []any{bigSec}
			err = templ.RenderCSSItems(ctx, templBuffer, var_41...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a href=\"")
			if err != nil {
				return err
			}
			var var_42 templ.SafeURL = makeSafeUrlWithAnchor(baseUrl, false, false, false, "fine-add")
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_42)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" class=\"")
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
			var_43 := `Close`
			_, err = templBuffer.WriteString(var_43)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></form></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"flex justify-center w-full\" id=\"fine-add\">")
			if err != nil {
				return err
			}
			var var_44 = []any{bigPri}
			err = templ.RenderCSSItems(ctx, templBuffer, var_44...)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<a class=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_44).String()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" href=\"")
			if err != nil {
				return err
			}
			var var_45 templ.SafeURL = makeSafeUrlWithAnchor(baseUrl, true, false, false, "fine-add")
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_45)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			if isFineMaster {
				var_46 := `Fine a Player`
				_, err = templBuffer.WriteString(var_46)
				if err != nil {
					return err
				}
			} else {
				var_47 := `Suggest a Fine`
				_, err = templBuffer.WriteString(var_47)
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
		var_48 := `
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
		_, err = templBuffer.WriteString(var_48)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
