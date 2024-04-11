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

func tomSelectLinks() templ.Component {
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
		_, err = templBuffer.WriteString("<link href=\"https://cdn.jsdelivr.net/npm/tom-select@2.3.1/dist/css/tom-select.css\" rel=\"stylesheet\"><script src=\"https://cdn.jsdelivr.net/npm/tom-select@2.3.1/dist/js/tom-select.complete.min.js\">")
		if err != nil {
			return err
		}
		var_2 := ``
		_, err = templBuffer.WriteString(var_2)
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

func getIdStr(id string) string {
	return fmt.Sprintf("#%s", id)
}

func success(msg string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_3 := templ.GetChildren(ctx)
		if var_3 == nil {
			var_3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"bg-green-100 border-l-4 border-green-500 text-green-700 p-4 rounded-lg\"><p class=\"text-lg font-semibold\">")
		if err != nil {
			return err
		}
		var var_4 string = msg
		_, err = templBuffer.WriteString(templ.EscapeString(var_4))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func errMsg(msg string) templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-lg\"><p class=\"text-lg font-semibold\">")
		if err != nil {
			return err
		}
		var var_6 string = msg
		_, err = templBuffer.WriteString(templ.EscapeString(var_6))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func fineSuperSelect(players []PlayerWithFines, approvedPFines []PresetFine) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_7 := templ.GetChildren(ctx)
		if var_7 == nil {
			var_7 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"container mx-auto bg-gray-200 shadow-xl m-10\" id=\"fine-ss\" hx-get=\"/fines/add\" hx-trigger=\"pageLoaded\" hx-target=\"#fine-ss\"><form id=\"ss-form\" hx-post=\"/fines-multi\" method=\"POST\" class=\"flex flex-col space-y-4 bg-white shadow-md p-6 rounded-lg\"><p>")
		if err != nil {
			return err
		}
		var_8 := `Select players and fines below:`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p><div class=\"flex flex-row\"><label for=\"select-fine\" class=\"mt-2 pr-2 font-semibold text-gray-700\">")
		if err != nil {
			return err
		}
		var_9 := `Fines`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><select id=\"select-fine\" name=\"pfines[]\" multiple placeholder=\"Select fine(s)...\" class=\"border border-gray-300 rounded-md text-gray-700 flex-grow mb-2\"><option value=\"\">")
		if err != nil {
			return err
		}
		var_10 := `Select a fine...`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option>")
		if err != nil {
			return err
		}
		for _, apf := range approvedPFines {
			_, err = templBuffer.WriteString("<option value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", apf.ID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var var_11 string = apf.Reason
			_, err = templBuffer.WriteString(templ.EscapeString(var_11))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</select></div><div class=\"flex flex-row\"><label for=\"select-player\" class=\"mt-2 pr-2 font-semibold text-gray-700\">")
		if err != nil {
			return err
		}
		var_12 := `Players`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><select id=\"select-player\" name=\"players[]\" multiple placeholder=\"Select player(s)...\" class=\"border border-gray-300 rounded-md text-gray-700 flex-grow mb-2\"><option value=\"\">")
		if err != nil {
			return err
		}
		var_13 := `Select a player...`
		_, err = templBuffer.WriteString(var_13)
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
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", p.ID)))
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
			_, err = templBuffer.WriteString("</option>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</select></div>")
		if err != nil {
			return err
		}
		var var_15 = []any{bigAdd}
		err = templ.RenderCSSItems(ctx, templBuffer, var_15...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_15).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_16 := `Multi-Fine`
		_, err = templBuffer.WriteString(var_16)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button><div id=\"results-container\"></div><script>")
		if err != nil {
			return err
		}
		var_17 := `
			var settings = {};
			new TomSelect("#select-fine",{
				maxOptions: 20,
				create: true,
				persist: false,
				plugins: {
					no_active_items: 'true',
					remove_button: {
						title:'Remove this fine',
					}
				},
				createFilter: function(input) {
					var match = input.match(/^[^,]*$/); // Example filter: disallow commas in input
					if(match) return !this.options.hasOwnProperty(input);
					return false;
				},
				onOptionAdd: function(value, item) {
					this.lock();
					fetch('/fines/add', { // Replace with your actual endpoint URL
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
						body: JSON.stringify({ reason: value }),
					})
					.then(response => {
						if (response.ok) {
							htmx.trigger("#ss-form", "pageLoaded")
							return response.json();
							
						} else {
							throw new Error('Server responded with an error');
						}
					})
					.then(data => {
						console.log(data.message); // Log the success message
						// The item is already added to the select; you might want to do something else here
					})
					.catch(error => {
						console.error('Error adding fine:', error);
						this.removeItem(value); // Remove the item if the server request failed
					})
					.finally(() => {
						this.unlock(); // Re-enable the select
					});
				},
			});
			new TomSelect("#select-player", {
				maxOptions: 20,
				plugins: {
					remove_button:{
						title:'Remove this player'
					}
				},
			});
			`
		_, err = templBuffer.WriteString(var_17)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></form></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func fineSuperSelectResults(players []PlayerWithFines, approvedPFines []PresetFine, newFines []Fine) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_18 := templ.GetChildren(ctx)
		if var_18 == nil {
			var_18 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		err = fineSuperSelect(players, approvedPFines).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		if len(newFines) > 0 {
			err = success(fmt.Sprintf("Added %d Fines", len(newFines))).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			for _, nf := range newFines {
				_, err = templBuffer.WriteString("<div>")
				if err != nil {
					return err
				}
				var var_19 string = fmt.Sprintf("%+v", nf)
				_, err = templBuffer.WriteString(templ.EscapeString(var_19))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"bg-orange-100 border-l-4 border-orange-500 text-orange-700 p-4\" role=\"alert\"><p>")
			if err != nil {
				return err
			}
			var_20 := `o fines added? Make sure to select fines/players above`
			_, err = templBuffer.WriteString(var_20)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p></div>")
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
