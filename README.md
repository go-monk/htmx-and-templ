The full code mentioned below can be found at https://github.com/go-monk/html-and-templ.

# 1) The Simplest Web Page

The simplest web page is just some data in HTML format transferred over the HTTP protocol and rendered in your browser. To implement this, we only need to create the data (stored in the `page` variable) and start an HTTP server that returns that data:

```go
var page = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Title</title>
</head>
<body>
	<h1>Heading 1</h1>
	<p>Paragraph</p>
</body>
</html>
`

func main() {
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(page))
}
```

To see for yourself start the server

```sh
cd 1
go run .
```

and visit http://localhost:8080.

# 2) Templating with `html/template`

That's nice and simple. But feels kind of hardcoded, static... What if we want some part of the page to change? We can use [html/template](https://pkg.go.dev/html/template) to add variables to the page:

```html
<!-- from 2/page.html -->

<body>
	<p>{{.Paragraph}}</p>
</body>
```

In this case, the variable part is `{{.Paragraph}}`. Here's how we fill it in:

```go
// from 2/main.go

type Page struct {
	Paragraph string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("page.html"))
	p := Page{
		Paragraph: time.Now().String(),
	}
	err := templ.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

Before sending the page to the browser (or other HTTP client), we replace the `{{.Paragraph}}` part in data from the `page.html` with the current time. Cool!

However, to refresh the time, we need to reload the whole page. What if we don’t want to do that? What if we want to refresh just part of the page — the variable part? This is called AJAX in web parlance: Asynchronous JavaScript and XML (it’s mostly JSON not XML these days). AJAX allows your web page to communicate with the server without a full reload. It's commonly used in modern web apps for smooth, dynamic user experiences.

# 3) htmx

Enter [htmx](https://htmx.org). Finally. It simplifies AJAX-like behavior by using just HTML attributes — no messing with JavaScript required!

First, add these into `page.html`:

```html
<script src="https://unpkg.com/htmx.org@2.0.4"></script>

<body>
	<p>Last full page reload</p>
	{{.FullReloadTime}}

	<p>Current time is</p>
	<div id="result">{{.FullReloadTime}}</div>

	<button 
		hx-get="/time" 
		hx-target="#result" 
		hx-swap="innerHTML">
		Refresh
	</button>
</body>
```

* The `<script>` tag includes the HTMX library from a CDN.
* `hx-get="/time"` makes a GET request to `/time`.
* `hx-target="#result"` injects the response into the `<div>` with `id="result"`.
* `hx-swap="innerHTML"` replaces the content inside that div.

We also need to implement the handler for the `/time` path:

```go
http.HandleFunc("/time", TimeHandler)

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, time.Now().Format(time.DateTime))
}
```

Now, when you reload the page in your browser, both timestamps will update. But if you click the "Refresh" button only the bottom timestamp updates:

![htmx demo](htmx.gif)

# 4) templ

If you don't like working with `html/template` standard library package for some reason, you can have a look at [templ](https://templ.guide).

To start using it for our project, replace `page.html` with `page.templ` that looks like this:

```go
package main

templ PageLayout(data Page) {
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<script src="https://unpkg.com/htmx.org@2.0.4"></script>
		<title>Title</title>
	</head>
	<body>
		@TimeDisplay(data)
	</body>
	</html>
}

templ TimeDisplay(data Page) {
	<p>Last full page reload</p>
	{data.FullReloadTime}
	<p>Current time is</p>
	<div id="result">{data.FullReloadTime}</div>
	<button 
		hx-get="/time" 
		hx-target="#result" 
		hx-swap="innerHTML">
		Refresh
	</button>
}
```

As you can see with templ you can define templates like Go functions that take arguments. You can also compose these components by using `@`.  

Next, change the `Handler` in `main.go` like this:

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	data := Page{
		FullReloadTime: time.Now().Format(time.DateTime),
	}
	PageLayout(data).Render(r.Context(), w)
}
```

Finally you need to get the templ package and binary and generate Go code from the `page.templ` file:

```sh
go get github.com/a-h/templ # add templ package
go install github.com/a-h/templ/cmd/templ@latest # install templ binary
templ generate # generate Go code from templ
```

Now you can run your Go server as before:

```sh
go run .
```