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

func fineList(fines []FineWithPlayer, page int, isFineMaster bool) templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"container mx-auto text-center\"><div class=\"text-3xl p-10\">")
		if err != nil {
			return err
		}
		var_2 := `Latest Fines`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div><table class=\"min-w-full divide-y divide-gray-200\"><thead class=\"bg-gray-50\"><tr><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">")
		if err != nil {
			return err
		}
		var_3 := `Reason`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">")
		if err != nil {
			return err
		}
		var_4 := `Amount`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">")
		if err != nil {
			return err
		}
		var_5 := `Player`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">")
		if err != nil {
			return err
		}
		var_6 := `Approved`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">")
		if err != nil {
			return err
		}
		var_7 := `When`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th></tr></thead><tbody class=\"bg-white divide-y divide-gray-200\">")
		if err != nil {
			return err
		}
		for _, f := range fines {
			_, err = templBuffer.WriteString("<tr")
			if err != nil {
				return err
			}
			if f.Fine.Approved {
				_, err = templBuffer.WriteString(" class=\"bg-white divide-y divide-gray-200\"")
				if err != nil {
					return err
				}
			} else {
				_, err = templBuffer.WriteString(" class=\"bg-yellow-200 divide-y divide-gray-200\"")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("><td class=\"px-6 py-4 whitespace-nowrap text-sm text-gray-900\">")
			if err != nil {
				return err
			}
			var var_8 string = f.Fine.Reason
			_, err = templBuffer.WriteString(templ.EscapeString(var_8))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td class=\"px-6 py-4 whitespace-nowrap text-sm text-gray-900\">")
			if err != nil {
				return err
			}
			var var_9 string = fmt.Sprintf("%v", f.Fine.Amount)
			_, err = templBuffer.WriteString(templ.EscapeString(var_9))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td>")
			if err != nil {
				return err
			}
			var var_10 string = f.Player.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_10))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td>")
			if err != nil {
				return err
			}
			if f.Fine.Approved {
				_, err = templBuffer.WriteString("<div>")
				if err != nil {
					return err
				}
				var_11 := `✅`
				_, err = templBuffer.WriteString(var_11)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			} else if isFineMaster {
				_, err = templBuffer.WriteString("<button hx-post=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/approve?fid=%d", f.Fine.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var_12 := `☐`
				_, err = templBuffer.WriteString(var_12)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</button>")
				if err != nil {
					return err
				}
			} else {
				var_13 := `(Pending approval)`
				_, err = templBuffer.WriteString(var_13)
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</td><td>")
			if err != nil {
				return err
			}
			var var_14 string = humanize.Time(f.Fine.CreatedAt)
			_, err = templBuffer.WriteString(templ.EscapeString(var_14))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td></tr>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</tbody></table><!--")
		if err != nil {
			return err
		}
		var_15 := `<div class="py-3">
			<button hx-get={ fmt.Sprintf("/load-more?page=%d", page +1) } hx-target="this" hx-swap="outerHTML" class="px-4 py-2 bg-blue-500 text-white font-semibold rounded hover:bg-blue-700">
				Load More
			</button>
		</div>`
		_, err = templBuffer.WriteString(var_15)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("--></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
