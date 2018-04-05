package meizitu

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

const meizitu_url string = "http://www.meizitu.com/"
const meizitu_cate_base_url string = "http://www.meizitu.com/a/"

type Meizitu struct {
	mainPageContent  string           //主页内容
	mainPageDocument goquery.Document //主页的goquery.Document
	categories       []Category       //分类切片
	selectedCategory Category
	currentAlbums    Queue
	currentPage      int
}

type Category struct {
	name     string //分类名称
	url      string //分类链接
	contents Queue  //分类分页链接
}

type CategoryNavs struct {
	name string
	url  string
}

type Album struct {
	name   string
	url    string
	imgUrl string
	imgs   Queue
}

type Img struct {
	name string
	url  string
}

func fetchContentByUrl(url string) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	content := string(bodyBytes)
	reg, err := regexp.Compile(`charset=(.*?)"`)
	if err != nil {
		panic(err.Error())
	}
	charset := reg.FindStringSubmatch(content)[1]
	if charset == "gb2312" {
		charset = "GBK"
	}
	decoder := mahonia.NewDecoder(charset)
	return decoder.ConvertString(content)
}

func fetchContentByReader(reader io.ReadCloser) string {
	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err.Error())
	}
	content := string(bodyBytes)
	reg, err := regexp.Compile(`charset=(.*?)"`)
	if err != nil {
		panic(err.Error())
	}
	charset := reg.FindStringSubmatch(content)[1]
	if charset == "gb2312" {
		charset = "GBK"
	}
	decoder := mahonia.NewDecoder(charset)
	return decoder.ConvertString(content)
}

func (meizi *Meizitu) fetchMainContent() {
	meizi.mainPageContent = fetchContentByUrl(meizitu_url)
	document, err := goquery.NewDocumentFromReader(strings.NewReader(meizi.mainPageContent))
	if err != nil {
		panic(err.Error())
	}
	meizi.mainPageDocument = *document
}

func (meizi *Meizitu) getCategories() {
	dom := meizi.mainPageDocument
	selection := dom.Find("div.tags span")
	//fmt.Println(selection)
	selection.Each(func(i int, selection *goquery.Selection) {
		selection.Find("a").Each(func(i int, selection *goquery.Selection) {
			//fmt.Println(selection.Text())
			url, exist := selection.Attr("href")
			if !exist {
				url = ""
			}
			meizi.categories = append(meizi.categories, Category{selection.Text(), url, Queue{}})
		})
	})
}

//region old fetchCategoryPages
//func (meizi *Meizitu) fetchCategoryPages(index int) {
//	meizi.selectedCategory = meizi.categories[index]
//	url := meizi.selectedCategory.url
//	ret := fetchContentByUrl(url)
//	document, err := goquery.NewDocumentFromReader(strings.NewReader(ret))
//	if err != nil {
//		panic(err.Error())
//	}
//	document.Find("ul.wp-list").Each(func(i int, selection *goquery.Selection) {
//		sel := selection.Find("h3.tit a")
//		url, _ := sel.Attr("href")
//		name := sel.First().Text()
//		fmt.Println(url, " ", name)
//	})
//	document.Find("div#wp_page_numbers ul").Each(func(i int, selection *goquery.Selection) {
//		name := selection.Text()
//
//		url, exists := selection.Attr("href")
//		if exists {
//			meizi.c
//		}
//
//	})
//
//}
//endregion

func (meizi *Meizitu) fetchCategoryPages(cate *Category, ch chan<- int) {
	urlString := cate.url

	//document.Find("ul.wp-list").Each(func(i int, selection *goquery.Selection) {
	//	sel := selection.Find("h3.tit a")
	//	url, _ := sel.Attr("href")
	//	name := sel.First().Text()
	//	cate.contents.Push(CategoryNavs{name, url})
	//})

	//document.Find("div#wp_page_numbers ul").Each(func(i int, selection *goquery.Selection) {
	//	//selection.Find("li").Each(func(i int, selection *goquery.Selection) {
	//	//	name := selection.Text()
	//	//
	//	//	url, exists := selection.Attr("href")
	//	//	if exists {
	//	//		cate.contents = append(cate.contents, CategoryContent{name, meizitu_cate_base_url + url})
	//	//	} else {
	//	//		cate.contents = append(cate.contents, CategoryContent{name, ""})
	//	//	}
	//	//})
	//
	//	fmt.Print(selection)
	//})

	//brow := surf.NewBrowser()
	//err := brow.SendGET(url)
	//if err!= nil{
	//	panic(err.Error())
	//}
	//response := brow.Response
	//closer := response.Body
	//defer closer.Close()
	//document, e := goquery.NewDocumentFromReader(closer)
	//if e!= nil{
	//	panic(e)
	//}
	//document.Find("div#wp_page_numbers ul").Each(func(i int, selection *goquery.Selection) {
	//	//	//selection.Find("li").Each(func(i int, selection *goquery.Selection) {
	//	//	//	name := selection.Text()
	//	//	//
	//	//	//	url, exists := selection.Attr("href")
	//	//	//	if exists {
	//	//	//		cate.contents = append(cate.contents, CategoryContent{name, meizitu_cate_base_url + url})
	//	//	//	} else {
	//	//	//		cate.contents = append(cate.contents, CategoryContent{name, ""})
	//	//	//	}
	//	//	//})
	//	v := selection.Text()
	//	if v!=""{
	//		fmt.Println(v)
	//	}
	//})

	resp, err := http.PostForm("http://127.0.0.1:12458", url.Values{"url": {urlString}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//document, err := goquery.NewDocumentFromReader(resp.Body)  won't work  don't know why

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}
	documentString := string(bytes)

	if err != nil {
		panic(err)
	}

	fmt.Println(cate.url)
	fmt.Println(documentString)
	document, err := goquery.NewDocumentFromReader(strings.NewReader(documentString))
	if err != nil {
		panic(err)
	}

	fmt.Println(document)

	//document.Find("div#wp_page_numbers ul").Each(func(i int, selection *goquery.Selection) {
	//
	//	name := selection.Text()
	//	fmt.Println(name, "  ")
	//
	//	//url, exists := selection.Attr("href")
	//	//if exists {
	//	//
	//	//} else {
	//	//}
	//
	//})
	ch <- 1
}

func (meizi *Meizitu) seek(page int) {
	//baseUrl := meizi.selectedCategory.url
	//strs := strings.Split(baseUrl, ".")
	//strs = strs[:len(strs)-1]
	//baseUrl = strings.Join(strs, ".") // instead of replace .html with empty string  because suffix can be .htm etc.
	//
	//if page >= 0 && page < meizi.selectedCategory.contents.Size() {
	//	meizi.currentPage = page
	//} else {
	//	meizi.currentPage = 0
	//}
}

func (meizi *Meizitu) Run() {
	var cateNumber int
	chFetchCategoryContent := make(chan int)

	meizi.fetchMainContent()
	meizi.getCategories()
	for i:=0; i< len(meizi.categories); i++ {
		go meizi.fetchCategoryPages(&meizi.categories[i], chFetchCategoryContent)
	}
	meizi.seek(0)

	sum := 0
	for i := range chFetchCategoryContent {
		sum += i
		if sum >= len(meizi.categories) {
			break
		}
	}

	fmt.Println("请输入你喜欢的分类编号： ")
	for index, cate := range meizi.categories {
		fmt.Println(index, ": ", cate.name)
	}

	//n, err := fmt.Scanf("%d", &cateNumber)
	_, err := fmt.Scanf("%d", &cateNumber)

	if err != nil {
		panic(err.Error())
	}
	if cateNumber >= 0 && cateNumber < len(meizi.categories) {
		// 输入值有效
		meizi.selectedCategory = meizi.categories[cateNumber]
	}

	//meizi.seek(1)

	defer func() {
		close(chFetchCategoryContent)
	}()
}
