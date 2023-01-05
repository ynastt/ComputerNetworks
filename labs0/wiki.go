package main

import (
  "html/template"//для обработки HTML-шаблонов только валидный html 
                 //генерируется действиями шаблона!
  "io/ioutil"
  "log"
  "net/http" //для создания веб-приложений
  //"regexp"
)

type Page struct {
  Title string
  Body  []byte
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)//записываем срез байтов в файл
                                                 //0600 =>файл должен быть создан с разрешением на чтение и 
                                                 //запись только для текущего пользователя
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

//обработчик для просмотра страницы.
//Если нет сохраненной страницы (т.е. с сохраненными данными из формы), то 
//перенаправляет на страницу для редактирования формы
func viewHandler(w http.ResponseWriter, 
                 r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  workWithTemplate(w, "view", p)
}

//обработчик для отображения формы редактирования страницы
func editHandler(w http.ResponseWriter, 
                 r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  workWithTemplate(w, "edit", p)
}

//обработчик для сохранения введенных через форму данных
//перенаправляет на страницу для просмотра только что сохранившихся данных формы
func saveHandler(w http.ResponseWriter, 
                 r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), 
               http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

//Функция template.Must - оболочка, которая паникует когда передано ненулевое значение error, а в противном случае возвращает *Template без изменений
//чтобы вызывать парсер один раз, а не каждый раз при открытии страницы
var templates = template.Must(
                   template.ParseFiles("edit.html", 
                                       "view.html"))
                                       
//функция для работы с шаблонами: читает содержимое edit.html и view.html
func workWithTemplate(w http.ResponseWriter, 
                    tmp string, p *Page) {
  err := templates.ExecuteTemplate(w, tmp+".html", p)
  if err != nil {
    http.Error(w, err.Error(), 
               http.StatusInternalServerError)
  }
}

/*
var validPath = regexp.MustCompile(
                   "^/(edit|save|view)/([a-zA-Z0-9]+)$")
                   
//
func makeHandler(funcHandler func(http.ResponseWriter, 
                 *http.Request, string)) http.HandlerFunc {
  	return func(w http.ResponseWriter, r *http.Request) {
    		m := validPath.FindStringSubmatch(r.URL.Path)
    		if m == nil {
      			http.NotFound(w, r)
      			return
    		}
    		funcHandler(w, r, m[2])
  		}
}*/
func makeHandler(funcHandler func(http.ResponseWriter, 
                 *http.Request, string)) http.HandlerFunc {
  	return func(w http.ResponseWriter, r *http.Request) {
    		//m := r.URL.Path
    		//if m == nil {
      		//	http.NotFound(w, r)
      		//	return
    		//}
    		funcHandler(w, r, r.URL.Path[4:])
  		}
}

func main() {
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  
  log.Fatal(http.ListenAndServe(":8000", nil))
}

