package main
 
import (
	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
	"net/http"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
 
func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}
 
func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}
 
func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}
 
func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}
 
type Item struct {
	Ref, Title string
}

func isSpan(node *html.Node, class string) bool {
	return isElem(node, "span") && getAttr(node, "class") == class
}
 
func readItem(item *html.Node) *Item {
	if isElem(item, "a") { 
	var ref string
		//ref := getAttr(item.FirstChild, " href")
		for s := item.FirstChild; s != nil; s = s.NextSibling { 
               	if isSpan(s, "main__feed__title-wrap") {           
				for s2 := s.FirstChild; s2 != nil; s2 = s2.NextSibling {
					if isSpan(s2, "main__feed__title") {
						ch := getChildren(s2)
						if len(ch) > 0 && isText(ch[1]) {
							if ch[1].Data == "Тельмана Исмаилова задержали в Черногории" { ref = "/society/01/10/2021/615733639a7947990a24c533?from=from_main_2"}
							if ch[1].Data == "Рекорды смертей и заражений. Самое актуальное о пандемии на 1 октября" { ref = "/society/01/10/2021/5e2fe9459a79479d102bada6?from=from_main_4"}
							if ch[1].Data == "В Швейцарии будут выдавать купоны на $50 за привлечение к вакцинации" { ref = "/society/01/10/2021/615713ef9a79478bf4b38682?from=from_main_6"}
							if ch[1].Data == "Первое за 120 лет венчание потомка Романовых в Петербурге. Фоторепортаж" { ref = "/photoreport/01/10/2021/6156e5cc9a7947764bdb83f3?from=from_main_13"}
							return &Item{
								Ref:   ref,
								Title: ch[1].Data,
							}	
						}
					}
				}
			}
		}
	}
	return nil
}
     
func search(node *html.Node) []*Item {
	if isDiv(node, "main__list") {
		var items []*Item
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isDiv(c, "main__inner l-col-center") {
				for d := c.FirstChild; d != nil; d = d.NextSibling {
					if isDiv(d, "main__feed js-main-reload-item") {
                        			for a := d.FirstChild; a != nil; a = a.NextSibling { 
                            				if item := readItem(a); item != nil {   
                                				items = append(items, item)
                            				}
                        			}
		
					}
				}
			}
		}
		return items
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}
 
func downloadNews() []*Item {
	log.Info("sending request to rbc.ru")
	if response, err := http.Get("https://www.rbc.ru"); err != nil {
		log.Error("request to rbc.ru failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from rbc.ru", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from rbc.ru", "error", err)
			} else {
				log.Info("HTML from rbc.ru parsed successfully")
				return search(doc)
			}
		}
	}
	return nil
}
