package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadpage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadpage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	errr := templates.ExecuteTemplate(w, "view.html", p)
	if errr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadpage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, err := template.ParseFiles("edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "edit.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page")}
	// p1.save()
	// p2, _ := loadpage("TestPage")
	// fmt.Println(string(p2.Body))
	http.HandleFunc("/view/", viewHandler)

	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
