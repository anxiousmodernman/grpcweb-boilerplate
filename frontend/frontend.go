package main

import (
	"strings"

	"honnef.co/go/js/dom"

	"github.com/anxiousmodernman/grpcweb-boilerplate/proto/client"
	"github.com/gopherjs/gopherjs/js"
	vecty "github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// Build this snippet with GopherJS, minimize the output and
// write it to html/frontend.js
//go:generate gopherjs build frontend.go -m -o html/frontend.js

// Integrate generated JS into a Go file for static loading.
//go:generate bash -c "go run assets_generate.go"

// This constant is very useful for interacting with
// the DOM dynamically
var document = dom.GetWindow().Document().(dom.HTMLDocument)

// Define no-op main since it doesn't run when we want it to
func main() {}

// Ensure our setup() gets called as soon as the DOM has loaded
func init() {
	document.AddEventListener("DOMContentLoaded", false, func(_ dom.Event) {
		go setup()
	})
}

// Setup is where we do the real work.
func setup() {
	p := &Page{}

	w := js.Global.Get("window")
	w = w.Call("addEventListener", "resize", func(e vecty.Event) { // avoid duplicate body
		// TODO: use debounce func here
		dims.Width = js.Global.Get("window").Get("innerWidth").Int64()
		dims.Height = js.Global.Get("window").Get("innerHeight").Int64()
		vecty.Rerender(p)
	})
	// This is the address to the server, and should be used
	// when creating clients.
	serverAddr := strings.TrimSuffix(document.BaseURI(), "/")

	// TODO: Use functions exposed by generated interface
	_ = client.NewProxyClient(serverAddr)

	//document.Body().SetInnerHTML(`<div><h2>GopherJS gRPC-Web is great!</h2></div>`)
	vecty.RenderBody(p)
}

type Dims struct {
	Width, Height int64
}

var dims Dims

// Page ...
type Page struct {
	vecty.Core
}

// Render implements vecty.Component for Page.
func (p *Page) Render() vecty.ComponentOrHTML {

	return elem.Body(
		vecty.Markup(
			vecty.Style("margin", "0"),
			vecty.Style("padding", "0"),
			vecty.Style("background", "#ccc"),
		),
		elem.Header(
			&NavComponent{},
		),
	)
}

// NavComponent ...
type NavComponent struct {
	vecty.Core
	Items []*NavItem
}

// Render ...
func (n *NavComponent) Render() vecty.ComponentOrHTML {

	var ulstyle vecty.MarkupList

	ulstyle = vecty.Markup(

		vecty.Style("list-style", "none"),
		vecty.Style("background-color", "#444"),
		vecty.Style("text-align", "center"),
		vecty.Style("padding", "0"),
		vecty.Style("margin", "0"),
	)
	if dims.Width > 600 {
		ulstyle = vecty.Markup(
			vecty.Style("list-style", "none"),
			vecty.Style("background-color", "#444"),
			vecty.Style("margin", "auto"),
			vecty.Style("width", "100%"),
			vecty.Style("overflow", "auto"),
		)
	}

	return elem.Div(
		elem.UnorderedList(
			&NavItem{Name: "first"},
			&NavItem{Name: "second"},
			&NavItem{Name: "third"},
			&NavItem{Name: "fourth"},
			&NavItem{Name: "fifth"},
			&NavItem{Name: "shoe"},
			ulstyle,
		),
	)
}

// NavItem ...
type NavItem struct {
	vecty.Core
	hovered bool
	active  bool
	Name    string
}

// Render ...
func (ni *NavItem) Render() vecty.ComponentOrHTML {
	var listyle vecty.MarkupList
	listyle = vecty.Markup(
		vecty.Style("font-family", "'Oswald', sans-serif"),
		vecty.Style("font-size", "1.2em"),
		vecty.Style("line-height", "40px"),
		vecty.Style("height", "40px"),
		vecty.Style("border-bottom", "1px solid #888"),
	)

	if dims.Width > 600 {
		listyle = vecty.Markup(
			vecty.Style("font-family", "'Oswald', sans-serif"),
			vecty.Style("font-size", "1.4em"),
			vecty.Style("line-height", "50px"),
			vecty.Style("height", "50px"),
			vecty.Style("width", "120px"),
			vecty.Style("float", "left"),
		)
	}

	var colr = "#fff"
	var bckgrnd = "#444"
	if ni.hovered {
		colr = "#005f5f"
		bckgrnd = "#005f5f"
	}

	var astyle vecty.MarkupList
	astyle = vecty.Markup(
		vecty.Style("text-decoration", "none"),
		vecty.Style("color", colr),
		vecty.Style("background-color", bckgrnd),
		vecty.Style("display", "block"),
		vecty.Style("transition", ".2s background-color"),
	)
	mo := event.MouseEnter(func(e *vecty.Event) {
		ni.hovered = true
		vecty.Rerender(ni)
	})
	ml := event.MouseLeave(func(e *vecty.Event) {
		ni.hovered = false
		vecty.Rerender(ni)
	})

	return elem.ListItem(
		listyle,
		elem.Anchor(
			vecty.Markup(mo),
			vecty.Markup(ml),
			astyle,
			vecty.Markup(vecty.Attribute("href", "#")),
			vecty.Text(ni.Name),
		),
	)
}

// MediaQuery ...
type MediaQuery struct {
	// Between
	Common []vecty.ComponentOrHTML
	Ranged []*MediaQueryStyle
}

// AddCommon ...
func (mq *MediaQuery) AddCommon(styles ...vecty.ComponentOrHTML) *MediaQuery {
	// for _, s := range styles {
	// 	mq.Common = append(mq.Common, s)
	// }
	mq.Common = append(mq.Common, styles...)
	return mq
}

func (mq *MediaQuery) AddRanged(min, max int, styles ...vecty.ComponentOrHTML) *MediaQuery {

	mq.Ranged = append(mq.Ranged, &MediaQueryStyle{Min: min, Max: max, Styles: styles})
	return mq
}

func (mq *MediaQuery) Apply() []vecty.List {
	var ret []vecty.ComponentOrHTML
	for _, c := range mq.Common {
		ret = append(ret, c)
	}

	// for _, s := range

	return nil
}

type MediaQueryStyle struct {
	Min, Max int
	Styles   []vecty.ComponentOrHTML
}

// Range ...
type Range struct {
	Min, Max int
}

// TODO fix this
func (r *Range) Within(val int) bool {
	if val >= r.Min && val <= r.Max {
		return true
	}
	return false
}

// if between 0-600
// if gt 600

/*
  .nav li {
    width: 120px;
    border-bottom: none;
    height: 50px;
    line-height: 50px;
    font-size: 1.4em;
  }


  .nav li {
    display: inline-block;
    margin-right: -4px;
  }

  .nav li {
    float: left;
  }
  .nav ul {
    overflow: auto;
    width: 600px;
    margin: 0 auto;
  }
  .nav {
    background-color: #444;
  }
*/
