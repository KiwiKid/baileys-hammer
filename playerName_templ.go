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

func playerNames(players []Player) templ.Component {
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
		_, err = templBuffer.WriteString("<div>")
		if err != nil {
			return err
		}
		var_2 := `✅`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		for _, p := range players {
			_, err = templBuffer.WriteString("<div id=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(F("pname-%d", p.ID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-swap-oob=\"true\">")
			if err != nil {
				return err
			}
			var var_3 string = p.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_3))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
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

func playerName(player Player) templ.Component {
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
		_, err = templBuffer.WriteString("<div>")
		if err != nil {
			return err
		}
		var var_5 string = player.Name
		_, err = templBuffer.WriteString(templ.EscapeString(var_5))
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

func playerInputSelector(players []Player, playerId uint64, inputType string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_6 := templ.GetChildren(ctx)
		if var_6 == nil {
			var_6 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		switch inputType {
		case "potd":
			_, err = templBuffer.WriteString("<div id=\"player-select-potd\" class=\"mt-4\">")
			if err != nil {
				return err
			}
			if len(UsePlayerOfTheDayName(ctx)) > 0 {
				_, err = templBuffer.WriteString("<label class=\"text-lg font-semibold\">")
				if err != nil {
					return err
				}
				var var_7 string = UsePlayerOfTheDayName(ctx)
				_, err = templBuffer.WriteString(templ.EscapeString(var_7))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" <select hx-ext=\"tomselect\" name=\"playerOfTheDay\" ts-max-items=\"1\" ts-item-class=\"text-2xl py-3\" ts-option-class=\"text-2xl py-3\" class=\"mt-1 w-full border-gray-300  bg-white rounded-md shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50\"><option value=\"\">")
				if err != nil {
					return err
				}
				var_8 := `N/A`
				_, err = templBuffer.WriteString(var_8)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
				for _, p := range players {
					_, err = templBuffer.WriteString("<option")
					if err != nil {
						return err
					}
					if p.ID == uint(playerId) {
						_, err = templBuffer.WriteString(" selected")
						if err != nil {
							return err
						}
					}
					_, err = templBuffer.WriteString(" value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(F("%d", p.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var var_9 string = p.Name
					_, err = templBuffer.WriteString(templ.EscapeString(var_9))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</option>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("</select></label>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		case "dotd":
			_, err = templBuffer.WriteString("<div id=\"player-select-dud-of-day\" class=\"mt-4\">")
			if err != nil {
				return err
			}
			if len(UseDudOfTheDayName(ctx)) > 0 {
				_, err = templBuffer.WriteString("<label class=\"text-lg font-semibold\">")
				if err != nil {
					return err
				}
				var var_10 string = UseDudOfTheDayName(ctx)
				_, err = templBuffer.WriteString(templ.EscapeString(var_10))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" <select id=\"player-select-dud-of-day\" name=\"dudOfTheDay\" hx-ext=\"tomselect\" ts-max-items=\"1\" ts-item-class=\"text-3xl py-3\" ts-option-class=\"text-3xl py-3\"><option value=\"\">")
				if err != nil {
					return err
				}
				var_11 := `N/A`
				_, err = templBuffer.WriteString(var_11)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
				for _, p := range players {
					_, err = templBuffer.WriteString("<option")
					if err != nil {
						return err
					}
					if p.ID == uint(playerId) {
						_, err = templBuffer.WriteString(" selected")
						if err != nil {
							return err
						}
					}
					_, err = templBuffer.WriteString(" value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(F("%d", p.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var var_12 string = p.Name
					_, err = templBuffer.WriteString(templ.EscapeString(var_12))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</option>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("</select></label>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		default:
			err = errMsg(fmt.Sprintf("No type for %s", inputType)).Render(ctx, templBuffer)
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

func playerHasEventCount(playerId uint, Events []MatchEvent, eventType string) uint {
	var count uint = 0
	for _, event := range Events {
		if event.PlayerId == playerId && event.EventType == eventType {
			count = count + 1
		}
	}
	return count
}

func getEventCount(Events []MatchEvent, eventType string) uint {
	var count uint = 0
	for _, event := range Events {
		if event.EventType == eventType {
			count = count + 1
		}
	}
	return count

}

func presentPlayerEvent(eventCount uint, p Player) string {
	switch eventCount {
	case 0:
		return ""
	case 1:
		return F("%s, ", p.Name)
	default:
		return F("%sx%d, ", p.Name, eventCount)
	}
}

func playerEventInputSelector(players []Player, events []MatchEvent, eventType string) templ.Component {
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
		switch eventType {
		case "injury":
			_, err = templBuffer.WriteString("<div id=\"player-select-injuries\" class=\"mt-4\">")
			if err != nil {
				return err
			}
			if len(UseInjuryCounterTrackerName(ctx)) > 0 {
				_, err = templBuffer.WriteString("<label class=\"text-lg font-semibold\">")
				if err != nil {
					return err
				}
				var var_14 string = UseInjuryCounterTrackerName(ctx)
				_, err = templBuffer.WriteString(templ.EscapeString(var_14))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var_15 := `(`
				_, err = templBuffer.WriteString(var_15)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				for _, p := range players {
					var var_16 string = presentPlayerEvent(playerHasEventCount(p.ID, events, eventType), p)
					_, err = templBuffer.WriteString(templ.EscapeString(var_16))
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var_17 := `)`
				_, err = templBuffer.WriteString(var_17)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" <select hx-ext=\"tomselect\" ts-no-delete=\"true\" ts-item-class=\"text-3xl py-3\" ts-option-class=\"text-3xl py-3\" name=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("eventType%s", Title(eventType))))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" multiple class=\"mt-1 w-full border-gray-300  bg-white rounded-md shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50\"><option value=\"\">")
				if err != nil {
					return err
				}
				var_18 := `N/A`
				_, err = templBuffer.WriteString(var_18)
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
					_, err = templBuffer.WriteString(templ.EscapeString(F("%d", p.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var var_19 string = p.Name
					_, err = templBuffer.WriteString(templ.EscapeString(var_19))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</option>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("</select></label>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		case "goal":
			_, err = templBuffer.WriteString("<div id=\"player-select-goals-for\" class=\"mt-4\">")
			if err != nil {
				return err
			}
			if UseShowGoalScorerMatchList(ctx) {
				_, err = templBuffer.WriteString("<label class=\"text-lg font-semibold\">")
				if err != nil {
					return err
				}
				var_20 := `New Goal Scorer: (`
				_, err = templBuffer.WriteString(var_20)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				for _, p := range players {
					var var_21 string = presentPlayerEvent(playerHasEventCount(p.ID, events, eventType), p)
					_, err = templBuffer.WriteString(templ.EscapeString(var_21))
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var_22 := `)`
				_, err = templBuffer.WriteString(var_22)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" <select hx-ext=\"tomselect\" ts-no-delete=\"true\" ts-item-class=\"text-3xl py-3\" ts-option-class=\"text-3xl py-3\" name=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("eventType%s", Title(eventType))))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" multiple class=\"mt-1 w-full border-gray-300  bg-white rounded-md shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50\">")
				if err != nil {
					return err
				}
				for _, p := range players {
					_, err = templBuffer.WriteString("<option value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(F("%d", p.ID)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var var_23 string = p.Name
					_, err = templBuffer.WriteString(templ.EscapeString(var_23))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</option>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("<option value=\"\">")
				if err != nil {
					return err
				}
				var_24 := `N/A`
				_, err = templBuffer.WriteString(var_24)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option></select></label>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		case "assist":
			_, err = templBuffer.WriteString("<div id=\"player-select-assist-for\" class=\"mt-4\">")
			if err != nil {
				return err
			}
			if UseShowGoalAssister(ctx) {
				_, err = templBuffer.WriteString("<label class=\"text-lg font-semibold\">")
				if err != nil {
					return err
				}
				var_25 := `New Assists:`
				_, err = templBuffer.WriteString(var_25)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" <select hx-ext=\"tomselect\" ts-no-delete=\"true\" ts-duplicates=\"true\" ts-item-class=\"text-3xl py-3\" ts-option-class=\"text-3xl py-3\" name=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("eventType%s", Title(eventType))))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" multiple class=\"mt-1 w-full border-gray-300  bg-white rounded-md shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50\"><option value=\"\">")
				if err != nil {
					return err
				}
				var_26 := `N/A`
				_, err = templBuffer.WriteString(var_26)
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
					_, err = templBuffer.WriteString(templ.EscapeString(F("%d", p.ID)))
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
					_, err = templBuffer.WriteString("</option>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("</select></label> ")
				if err != nil {
					return err
				}
				for _, p := range players {
					var var_28 string = presentPlayerEvent(playerHasEventCount(p.ID, events, eventType), p)
					_, err = templBuffer.WriteString(templ.EscapeString(var_28))
					if err != nil {
						return err
					}
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		case "conceded-goal":
			_, err = templBuffer.WriteString("<div id=\"player-select-conceded-goal\" class=\"mt-4\">")
			if err != nil {
				return err
			}
			if UseShowGoalScorerMatchList(ctx) {
				_, err = templBuffer.WriteString("<label class=\"text-lg font-semibold\">")
				if err != nil {
					return err
				}
				var_29 := `Opponent Goals (`
				_, err = templBuffer.WriteString(var_29)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_30 string = F("%d", getEventCount(events, eventType))
				_, err = templBuffer.WriteString(templ.EscapeString(var_30))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var_31 := `):`
				_, err = templBuffer.WriteString(var_31)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" <select hx-ext=\"tomselect\" ts-no-delete=\"true\" ts-max-options=\"99\" ts-duplicates=\"true\" ts-item-class=\"text-3xl py-3\" ts-option-class=\"text-3xl py-3\" name=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("eventType%s", Title(eventType))))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" multiple class=\"mt-1 w-full border-gray-300  bg-white rounded-md shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50\"><option value=\"yes\">")
				if err != nil {
					return err
				}
				var_32 := `Opponent Goal`
				_, err = templBuffer.WriteString(var_32)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option><option value=\"yes1\">")
				if err != nil {
					return err
				}
				var_33 := `Opponent Goal`
				_, err = templBuffer.WriteString(var_33)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option><option value=\"yes2\">")
				if err != nil {
					return err
				}
				var_34 := `Opponent Goal`
				_, err = templBuffer.WriteString(var_34)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option><option value=\"yes3\">")
				if err != nil {
					return err
				}
				var_35 := `Opponent Goal`
				_, err = templBuffer.WriteString(var_35)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option><option value=\"yes4\">")
				if err != nil {
					return err
				}
				var_36 := `Opponent Goal`
				_, err = templBuffer.WriteString(var_36)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option><option value=\"yes5\">")
				if err != nil {
					return err
				}
				var_37 := `Opponent Goal`
				_, err = templBuffer.WriteString(var_37)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option><option value=\"yes6\">")
				if err != nil {
					return err
				}
				var_38 := `Opponent Goal`
				_, err = templBuffer.WriteString(var_38)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option><option value=\"yes7\">")
				if err != nil {
					return err
				}
				var_39 := `Opponent Goal`
				_, err = templBuffer.WriteString(var_39)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option><option value=\"yes8\">")
				if err != nil {
					return err
				}
				var_40 := `Opponent Goal`
				_, err = templBuffer.WriteString(var_40)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option></select></label>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		default:
			err = errMsg(fmt.Sprintf("No eventType for %s", eventType)).Render(ctx, templBuffer)
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

func playerRoleSelector(player PlayerWithFines, config *Config, msg string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_41 := templ.GetChildren(ctx)
		if var_41 == nil {
			var_41 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"w-full\" id=\"players-ss\"><form class=\"todo\" hx-post=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/players?playerId=%d", player.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\"><div class=\"p-2\"><input type=\"hidden\" name=\"ID\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", player.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><label for=\"role\" class=\"block mt-2\">")
		if err != nil {
			return err
		}
		var_42 := `Name      `
		_, err = templBuffer.WriteString(var_42)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <input type=\"text\" name=\"Name\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(player.Name))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" id=\"name\" placeholder=\"Player name\" class=\"w-full  text-3xl px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></label>")
		if err != nil {
			return err
		}
		if UseRoles(ctx) {
			_, err = templBuffer.WriteString("<label for=\"role\" class=\"block mt-2\">")
			if err != nil {
				return err
			}
			var_43 := `Role      `
			_, err = templBuffer.WriteString(var_43)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" <input type=\"text\" name=\"role\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(player.Role))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" id=\"role\" placeholder=\"Role\" class=\"w-full  text-3xl px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\"></label> <label for=\"role\" class=\"block mt-2\">")
			if err != nil {
				return err
			}
			var_44 := `Role Description      `
			_, err = templBuffer.WriteString(var_44)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" <input type=\"text\" name=\"roleDescription\" id=\"roleDescription\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(player.RoleDescription))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" placeholder=\"Role Description\" class=\"w-full text-3xl px-4 py-2 mt-2 border rounded-md focus:outline-none focus:ring-1 focus:ring-blue-600\">")
			if err != nil {
				return err
			}
			if len(UseInjuryCounterTrackerName(ctx)) > 0 {
				_, err = templBuffer.WriteString("<div class=\"flex flex-row mt-2\">")
				if err != nil {
					return err
				}
				var var_45 = []any{bigPri}
				err = templ.RenderCSSItems(ctx, templBuffer, var_45...)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("<button type=\"submit\" hx-get=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/match/%d/event")))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" class=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_45).String()))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var_46 := `Add Injury`
				_, err = templBuffer.WriteString(var_46)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</button></div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</label>")
			if err != nil {
				return err
			}
		}
		if len(msg) > 0 {
			err = success(msg).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<div class=\"flex flex-row mt-2\">")
		if err != nil {
			return err
		}
		var var_47 = []any{bigPri}
		err = templ.RenderCSSItems(ctx, templBuffer, var_47...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_47).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_48 := `Update Player`
		_, err = templBuffer.WriteString(var_48)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></div><div class=\"flex flex-row mt-2\">")
		if err != nil {
			return err
		}
		var var_49 = []any{bigDel}
		err = templ.RenderCSSItems(ctx, templBuffer, var_49...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button type=\"submit\" hx-delete=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/players?playerId=%d", player.ID)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-confirm=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("Are you sure you want to delete %s?", player.Name)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"")
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
		var_50 := `Delete `
		_, err = templBuffer.WriteString(var_50)
		if err != nil {
			return err
		}
		var var_51 string = player.Name
		_, err = templBuffer.WriteString(templ.EscapeString(var_51))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></div></div></form></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
