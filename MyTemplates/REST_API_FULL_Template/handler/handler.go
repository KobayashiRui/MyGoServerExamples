package handler
import (
	"html/template"
	"net/http"
	"log"
	//"time"
	//"errors"
	//"strconv"
)

type Home struct {
	ErrorLog *log.Logger
	InfoLog *log.Logger 
}

func (h *Home) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFound(w)
		return
	}

	ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	if err != nil {
		h.ErrorLog.Println(err.Error())
		//http.Error(w, "Internal Server Error", 500)
		serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		h.ErrorLog.Println(err.Error())
		serverError(w, err)
	}
}


