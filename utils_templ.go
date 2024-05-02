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
	"time"
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
		_, err = templBuffer.WriteString("</script><script src=\"https://kiwikid.github.io/hx-tomselect/hx-tom-select.js\">")
		if err != nil {
			return err
		}
		var_3 := ``
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><!--")
		if err != nil {
			return err
		}
		var_4 := `<script>
	(function() {   
    /** stable build*/
    const version = '06'

    /**
     * @typedef {Object} SupportedAttribute
     * Defines an attribute supported by a configuration modification system.
     * @property {string} key - The key of the configuration attribute to modify.
     * @property {ConfigChange} configChange - The modifications to apply to the TomSelect configuration.
     */

    /**
     * @typedef {'simple' | 'callback'} AttributeType
     */

      /**
     * @typedef {function(HTMLElement, Object):void} CallbackFunction
     * Description of what the callback does and its parameters.
     * @param {string} a - The first number parameter.
     
     */
      /**
     * @typedef {Object} AttributeConfig
     * Defines an attribute supported by a configuration modification system.
     * @property {string} key - The key of the configuration attribute to modify.
     * @property {string} _description
     * @property {ConfidenceLevel} _isBeta
     * @property {CallbackFunction|string|null} configChange - The modifications to apply to the TomSelect configuration.
     * 
     */



    /**
     * @type {SupportedAttribute[]}
     */

    /**
     * @typedef {'ts-max-items' | 'ts-max-options' | 'ts-create' | 'ts-sort' | 'ts-sort-direction' | 'ts-allow-empty-option', 'ts-clear-after-add', 'ts-raw-config', 'ts-create-on-blur', 'ts-no-delete'} TomSelectConfigKey
     * Defines the valid keys for configuration options in TomSelect.
     * Each key is a string literal corresponding to a specific property that can be configured in TomSelect.
     */

    /**
     * @type {Array<AttributeConfig>}
     */
    const attributeConfigs = [
        {
            key: 'ts-create',
            configChange: 'create',
            _description: 'Allow creating new items'
        },{
            key: 'ts-create-on-blur',
            configChange: 'createOnBlur'
        },{
            key: 'ts-create-filter',
            configChange:  (elm, config) => ({
                createFilter: function(input) {
                    try {
                        const filter = elm.getAttribute('ts-create-filter')
                        const matchEx = filter == "true" ? /^[^,]*$/ : elm.getAttribute('ts-create-filter')
                        var match = input.match(matchEx); // Example filter: disallow commas in input
                        if(match) return !this.options.hasOwnProperty(input);
                        elm.setAttribute('tom-select-warning', JSON.stringify(err));
                        return false;
                    } catch (err) {
                        return false
                    }
                }
            })
        },{
            key: 'ts-delimiter',
            configChange: 'delimiter'
        },{
            key: 'ts-highlight',
            configChange: 'highlight'
        },{
            key: 'ts-multiple',
            configChange: 'multiple'
        },{
            key: 'ts-persist',
            configChange: 'persist'
        },{
            key: 'ts-open-on-focus',
            configChange: 'openOnFocus'
        },{
            key: 'ts-max-items',
            configChange: 'maxItems'
        },{
            key: 'ts-hide-selected',
            configChange: 'hideSelected'
        },{
            key: 'tx-close-after-select',
            configChange: 'closeAfterSelect'
        },{
            key: 'tx-duplicates',
            configChange: 'duplicates'
        },
        {
            key: 'ts-max-options',
            configChange: 'maxOptions'
        },{
            key: 'ts-sort',
            configChange: (elm, config) => ({
                sortField: {
                    field: elm.getAttribute('ts-sort'),
                },
            })
        },{
            key: 'ts-sort-direction',
            configChange: (elm, config) => ({
                sortField: {
                    direction: elm.getAttribute('ts-sort-direction') ?? 'asc'
                },
            })
        },{
            key: 'ts-allow-empty-option',
            type: 'simple',
            configChange: 'allowEmptyOption'
        },{
            key: 'ts-clear-after-add',
            configChange: {
                create: true,
                onItemAdd: function() {
					debugger
                    this.setTextboxValue('');
               //     this.refreshOptions();
                }
            }
        },{
            key: 'ts-remove-button-title',
            configChange: (elm, config) => deepAssign(config,{
                plugins: {
                    remove_button: {
                        title: elm.getAttribute('ts-remove-button-title') == 'true' ? 'Remove this item' : elm.getAttribute('ts-remove-button-title')
                    }
                },
            })
        },{
            key: 'ts-delete-confirm',
            configChange: (elm, config) => ({
                onDelete: function(values) {
                    if(elm.getAttribute('ts-delete-confirm') == "true"){
                        return confirm(values.length > 1 ? 'Are you sure you want to remove these ' + values.length + ' items?' : 'Are you sure you want to remove "' + values[0] + '"?');
                    }else {
                        return confirm(elm.getAttribute('ts-delete-confirm'));
                    }
                    
                }
            })
        },{
            key: 'ts-add-post-url',
            configChange: (elm, config) => ({
                    onOptionAdd: function(value, item) {
                        this.lock();
                        const valueKeyName = elm.getAttribute('ts-add-post-url-body-value') ?? 'value'
                        const body = {}
                        body[valueKeyName] = value
                        fetch(elm.getAttribute('ts-add-post-url'), {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                            },
                            body: JSON.stringify(body),
                        })
						.then((res) => {
							if (!res.ok) {
								throw new Error(` + "`" + `HTTP status ${res.status}` + "`" + `); // Throw an error if response is not ok.
							}
							return res.text(); // Return the response as text for HTMX processing.
						})
                        .then((responseHtml) => htmx.process(elm, responseHtml))
                        .catch(error => {
                            console.error('Error adding item', error)
                            elm.setAttribute('tom-select-warning', ` + "`" + `ts-add-post-url - Error processing item: ${JSON.stringify(error)}` + "`" + `);
                            this.removeItem(value);
                        })
                        .finally(() => {
                            this.unlock();
                        });
                }
            }),
            _isBeta: true,
        },{
            key: 'ts-add-post-url-body-value',
            configChange: '',
            _isBeta: true,
        },
        {
            key: 'ts-no-active',
            configChange: {
                plugins: ['no_active_items'],
                persist: false,
                create: true
            }
        },{
            key: 'ts-remove-selector-on-select',
            type: 'simple',
            configChange: null
        },{
            key: 'ts-no-delete',
            configChange: {
                onDelete: () => { return false},
            }
        },{
            key: 'ts-option-class',
            configChange: 'optionClass'
        },{
            key: 'ts-option-class-ext',
            configChange: (elm, config) => ({
                'optionClass': ` + "`" + `${elm.getAttribute('ts-option-class-ext')} option` + "`" + `
            })
        },{
            key: 'ts-item-class',
            configChange: 'itemClass'
        },{
            key: 'ts-item-class-ext',
            configChange:(elm, config) => ({
                key: 'ts-option-class-ext',
                configChange: {
                    'itemClass': ` + "`" + `${elm.getAttribute('ts-option-class-ext')} item` + "`" + `
                }
            })
        },
        {
            key: 'ts-raw-config',
            configChange: (elm, config) => elm.getAttribute('ts-raw-config')
        }
    ]

    /**
     * Deeply assigns properties to an object, merging any existing nested properties.
     * 
     * @param {Object} target The target object to which properties will be assigned.
     * @param {Object} updates The updates to apply. This object can contain deeply nested properties.
     * @returns {Object} The updated target object.
     */
    function deepAssign(target, updates) {
        Object.keys(updates).forEach(key => {
            if (typeof updates[key] === 'object' && updates[key] !== null && !Array.isArray(updates[key])) {
                if (!target[key]) target[key] = {};
                deepAssign(target[key], updates[key]);
            } else {
                target[key] = updates[key];
            }
        });
        return target;
    }

    function attachTomSelect(s){
        try {
            if(s.attributes?.length == 0){
                throw new Error("no attributes on select?")
            }
            
            let config = {
                maxItems: 999,
                plugins: {}
            };
            const debug = s.getAttribute('hx-ext')?.split(',').map(item => item.trim()).includes('debug');
            if (debug) { console.log(s.attributes) }

            Array.from(s.attributes).forEach((a) => {
                const attributeConfig = attributeConfigs.find((ac) => ac.key == a.name)
                if (attributeConfig != null){
                    let configChange = {}
                    if(typeof attributeConfig.configChange == 'string'){
                        configChange[attributeConfig.configChange] = a.value
                    }else if(typeof attributeConfig.configChange == 'function'){
                        configChange = attributeConfig.configChange(s, config)
                    }else if(typeof attributeConfig.configChange == 'object'){
                        configChange = attributeConfig.configChange
                    }else if(a.name.startsWith('ts-')) {
                        s.setAttribute('tom-select-warning', ` + "`" + `Invalid config key found: ${attr.name}` + "`" + `);
                        console.warn(` + "`" + `Could not find config match:${JSON.stringify(attributeConfig)}` + "`" + `)
                    }
                
                    deepAssign(config, configChange)
                }else if(a.name.startsWith('ts-')){
                    console.warn(` + "`" + `Invalid config key found: ${a.name}` + "`" + `);
                    s.setAttribute(` + "`" + `tom-select-warning_${a.name}` + "`" + `, ` + "`" + `Invalid config key found` + "`" + `);
                }
            })

        if (debug) {  console.info('hx-tomselect - tom-select-success - config', config) }
        const ts = new TomSelect(s, config);
        s.setAttribute('tom-select-success', ` + "`" + `success` + "`" + `);
        s.setAttribute('hx-tom-select-version', ` + "`" + `hx-ts-${version}_ts-${ts.version}` + "`" + `);

    } catch (err) {
        s.setAttribute('tom-select-error', JSON.stringify(err));
        console.error(` + "`" + `htmx-tomselect - Failed to load hx-tomsselect ${err}` + "`" + `);
    }
    }

    htmx.defineExtension('tomselect', {
        // This is doing all the tom-select attachment at this stage, but relies on this full document scan (would prefer onLoad of speicfic content):
        onEvent: function (name, evt) {
            if (name === "htmx:afterProcessNode") {
                const newSelects = document.querySelectorAll('select[hx-ext*="tomselect"]:not([tom-select-success]):not([tom-select-error])')
                newSelects.forEach((s) => {
                    attachTomSelect(s)
                })
            }
        },
        onLoad: function (content) {
            console.log('onLoad')
                    const newSelects = content.querySelectorAll('select[hx-ext*="tomselect"]:not([tom-select-success]):not([tom-select-error])')
                    newSelects.forEach((s) => {
                        attachTomSelect(s)
                    })

            // When the DOM changes, this block ensures TomSelect will reflect the current html state (i.e. new <option selected></option> will be respected)
            // Still evaulating the need of this
             /*   const selectors = document.querySelectorAll('select[hx-ext*="tomselect"]')
            selectors.forEach((s) => {
                console.log('SYNC RAN')
                s.tomselect.clear();
                s.tomselect.clearOptions();
                s.tomselect.sync(); 
            })*/
        },
		beforeHistorySave: function(){
			document.querySelectorAll('select[hx-ext*="tomselect"]')
            	.forEach(elt => elt.tomselect.destroy())
		}
    });

})();
	</script>`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("-->")
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

func contextSuccess(matchId uint64, contextStr string, fineAt time.Time) templ.Component {
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
		var_6 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			_, err = templBuffer.WriteString("<h1>")
			if err != nil {
				return err
			}
			var_7 := `Added Context`
			_, err = templBuffer.WriteString(var_7)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1> <div>")
			if err != nil {
				return err
			}
			var var_8 string = fmt.Sprintf("%d", matchId)
			_, err = templBuffer.WriteString(templ.EscapeString(var_8))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div> <div>")
			if err != nil {
				return err
			}
			var var_9 string = contextStr
			_, err = templBuffer.WriteString(templ.EscapeString(var_9))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div> <div>")
			if err != nil {
				return err
			}
			if fineAt.After(twoWeeksAgo) {
				var var_10 string = humanize.Time(fineAt)
				_, err = templBuffer.WriteString(templ.EscapeString(var_10))
				if err != nil {
					return err
				}
			} else {
				var var_11 string = fineAt.Format("2006-01-02T15:04")
				_, err = templBuffer.WriteString(templ.EscapeString(var_11))
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = successComp().Render(templ.WithChildren(ctx, var_6), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func successComp() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_12 := templ.GetChildren(ctx)
		if var_12 == nil {
			var_12 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"bg-green-100 border-l-4 border-green-500 text-green-700 p-4 rounded-lg\"><p class=\"text-lg font-semibold\">")
		if err != nil {
			return err
		}
		err = var_12.Render(ctx, templBuffer)
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

func success(msg string) templ.Component {
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
		if len(msg) > 0 {
			_, err = templBuffer.WriteString("<div class=\"bg-green-100 border-l-4 border-green-500 text-green-700 p-4 rounded-lg\"><p class=\"text-lg font-semibold\">")
			if err != nil {
				return err
			}
			var var_14 string = msg
			_, err = templBuffer.WriteString(templ.EscapeString(var_14))
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

func warning(msg string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_15 := templ.GetChildren(ctx)
		if var_15 == nil {
			var_15 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if len(msg) > 0 {
			_, err = templBuffer.WriteString("<div class=\"p-4 mb-4 text-blue-800 border border-blue-300 rounded-lg bg-blue-50 dark:bg-gray-800 dark:text-blue-400 dark:border-blue-800\"><p class=\"text-lg font-semibold\">")
			if err != nil {
				return err
			}
			var var_16 string = msg
			_, err = templBuffer.WriteString(templ.EscapeString(var_16))
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

func errMsg(msg string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_17 := templ.GetChildren(ctx)
		if var_17 == nil {
			var_17 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<div class=\"bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-lg\"><p class=\"text-lg font-semibold\">")
		if err != nil {
			return err
		}
		var var_18 string = msg
		_, err = templBuffer.WriteString(templ.EscapeString(var_18))
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

func niceDate(date *time.Time) string {
	return date.Format("January 2, 2006")
}

func isSelected(selectedIds []uint, id uint) bool {
	for _, sid := range selectedIds {
		if sid == id {
			return true
		}
	}
	return false
}

func fineSuperSelect(players []PlayerWithFines, approvedPFines []PresetFine, selectedFineIds []uint) templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"w-full mx-auto bg-gray-200 shadow-xl m-10\" id=\"fine-ss\" hx-get=\"/fines/add\" hx-trigger=\"pageLoaded\" hx-target=\"#fine-ss\"><form id=\"ss-form\" hx-post=\"/fines-multi\" method=\"POST\" class=\"flex flex-col space-y-4 bg-white shadow-md p-6 rounded-lg\"><p class=\"font-bold text-3xl\">")
		if err != nil {
			return err
		}
		var_20 := `Select players and fines:`
		_, err = templBuffer.WriteString(var_20)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p><div class=\"flex flex-row\"><select id=\"select-fine\" hx-ext=\"tomselect\" ts-persist=\"false\" ts-create=\"true\" ts-create-filter=\"true\" ts-create-on-blur=\"true\" ts-open-on-focus=\"true\" ts-add-post-url=\"/fines/add\" ts-add-post-url-body-value=\"reason\" ts-item-class=\"text-3xl py-3\" ts-option-class=\"text-3xl w-full py-3\" tx-max-items=\"99\" name=\"pfines[]\" multiple required placeholder=\"Select fine(s)...\" class=\"text-3xl border border-gray-300 rounded-md text-gray-700 flex-grow mb-2\"><option value=\"\">")
		if err != nil {
			return err
		}
		var_21 := `Select a fine...`
		_, err = templBuffer.WriteString(var_21)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option>")
		if err != nil {
			return err
		}
		for _, apf := range approvedPFines {
			if isSelected(selectedFineIds, apf.ID) {
				_, err = templBuffer.WriteString("<option selected value=\"")
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
				var var_22 string = apf.Reason
				_, err = templBuffer.WriteString(templ.EscapeString(var_22))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
			} else {
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
				var var_23 string = apf.Reason
				_, err = templBuffer.WriteString(templ.EscapeString(var_23))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
			}
		}
		_, err = templBuffer.WriteString("</select></div><div class=\"flex flex-row\"><select id=\"select-player\" required hx-ext=\"tomselect\" tx-max-items=\"99\" name=\"players[]\" ts-item-class=\"text-3xl py-3\" tx-close-after-select=\"true\" ts-option-class=\"text-3xl w-full py-3\" multiple placeholder=\"Select player(s)...\" class=\"text-3xl  border border-gray-300 rounded-md text-gray-700 flex-grow mb-2\"><option value=\"\">")
		if err != nil {
			return err
		}
		var_24 := `Select a player...`
		_, err = templBuffer.WriteString(var_24)
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
			var var_25 string = p.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_25))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</select></div><div class=\"flex flex-row\"><input type=\"text\" name=\"context\" value=\"\" class=\"px-4 py-2 border border-gray-300 rounded-lg w-full focus:ring-blue-500 focus:border-blue-500\" placeholder=\"Context for the fine\"></div>")
		if err != nil {
			return err
		}
		var var_26 = []any{bigAdd}
		err = templ.RenderCSSItems(ctx, templBuffer, var_26...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<button type=\"submit\" class=\"")
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
		var_27 := `Create Fine(s)`
		_, err = templBuffer.WriteString(var_27)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button><div id=\"results-container\"></div></form></div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func getFineIds(fines []Fine) []uint {
	var ids = []uint{}
	for _, f := range fines {
		ids = append(ids, f.ID)
	}
	return ids
}

func fineSuperSelectResults(players []PlayerWithFines, approvedPFines []PresetFine, newFines []Fine) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_28 := templ.GetChildren(ctx)
		if var_28 == nil {
			var_28 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		err = fineSuperSelect(players, approvedPFines, getFineIds(newFines)).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		if len(newFines) > 0 {
			_, err = templBuffer.WriteString("<div class=\"text-2xl list-none bg-green-100 border-l-4 border-green-500 text-green-700 p-4 rounded-lg\"><h1 class=\"text-3xl\">")
			if err != nil {
				return err
			}
			var var_29 string = fmt.Sprintf("Added %d Fines:", len(newFines))
			_, err = templBuffer.WriteString(templ.EscapeString(var_29))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1>")
			if err != nil {
				return err
			}
			if len(newFines) > 0 {
				for _, p := range players {
					for _, nf := range newFines {
						if p.ID == nf.PlayerID {
							_, err = templBuffer.WriteString("<details><summary>")
							if err != nil {
								return err
							}
							var var_30 string = fmt.Sprintf("%s - %s - %s", nf.Reason, p.Name, nf.Context)
							_, err = templBuffer.WriteString(templ.EscapeString(var_30))
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString("</summary> ")
							if err != nil {
								return err
							}
							var var_31 string = fmt.Sprintf("%+v", nf)
							_, err = templBuffer.WriteString(templ.EscapeString(var_31))
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString("</details> ")
							if err != nil {
								return err
							}
							var var_32 = []any{smPri}
							err = templ.RenderCSSItems(ctx, templBuffer, var_32...)
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString("<button hx-get=\"")
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines/edit/%d?isEdit=form", nf.ID)))
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-target=\"this\" class=\"")
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
							if len(nf.Context) == 0 {
								var_33 := `Add Context	`
								_, err = templBuffer.WriteString(var_33)
								if err != nil {
									return err
								}
							} else {
								var_34 := `Edit Context`
								_, err = templBuffer.WriteString(var_34)
								if err != nil {
									return err
								}
							}
							_, err = templBuffer.WriteString("</button> ")
							if err != nil {
								return err
							}
							var var_35 = []any{del}
							err = templ.RenderCSSItems(ctx, templBuffer, var_35...)
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString("<button class=\"")
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_35).String()))
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString("\" hx-confirm=\"Are you sure you want to delete the fine by this player?\" hx-delete=\"")
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/fines?fid=%d", nf.ID)))
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString("\" hx-swap=\"outerHTML\" hx-target=\"this\">")
							if err != nil {
								return err
							}
							var_36 := `Undo`
							_, err = templBuffer.WriteString(var_36)
							if err != nil {
								return err
							}
							_, err = templBuffer.WriteString("</button>")
							if err != nil {
								return err
							}
						}
					}
				}
			} else {
				_, err = templBuffer.WriteString("<div>")
				if err != nil {
					return err
				}
				var_37 := `No fines created? `
				_, err = templBuffer.WriteString(var_37)
				if err != nil {
					return err
				}
				var var_38 string = fmt.Sprintf("%d %d", len(approvedPFines))
				_, err = templBuffer.WriteString(templ.EscapeString(var_38))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<div class=\"bg-orange-100 border-l-4 border-orange-500 text-orange-700 p-4\" role=\"alert\"><p>")
			if err != nil {
				return err
			}
			var_39 := `no fines added? Make sure to select fines/players above`
			_, err = templBuffer.WriteString(var_39)
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
