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
	"strconv"
	"errors"
)

const meizitu_url string = "http://www.meizitu.com/"
const meizitu_cate_base_url string = "http://www.meizitu.com/a/"

type Meizitu struct {
	mainPageContent  string           //主页内容
	mainPageDocument goquery.Document //主页的goquery.Document
	categories       []Category       //分类切片
	selectedCategory *Category
	currentAlbums    Queue
	currentPage      int
}

type Category struct {
	name     string //分类名称
	url      string //分类链接
	contents Queue  //分类分页链接
}

type CategoryNav struct {
	name string
	url  string
}

type Album struct {
	name   string
	url    string
	icon string
	imgs   Queue
}

type Img struct {
	name string
	url  string
}

func (meizi *Meizitu) fetchContentByUrl(url string) string {
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

func (meizi *Meizitu) fetchContentByReader(reader io.ReadCloser) string {
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
	meizi.mainPageContent = meizi.fetchContentByUrl(meizitu_url)
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

//var same string

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

	var (
		resp *http.Response
		err error
	)
	for {
		resp, err = http.PostForm("http://127.0.0.1:12458", url.Values{"url": {urlString}})
		if err != nil {
			continue
		} else {
			break
		}
	}

	defer resp.Body.Close()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	//document, err := goquery.NewDocumentFromReader(resp.Body)  won't work  don't know why

	//fmt.Printf("\nres.Body %p\n", &resp.Body)
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}
	documentString := string(bytes)

	//if same == ""{
	//	same = documentString
	//} else {
	//	if same == documentString{
	//		fmt.Println("true")
	//	}else {
	//		fmt.Println("false")
	//	}
	//	same = documentString
	//}

	if err != nil {
		panic(err)
	}

	fmt.Println(urlString)
	//fmt.Printf("\n%p", &documentString)
	document, err := goquery.NewDocumentFromReader(strings.NewReader(documentString))
	if err != nil {
		panic(err)
	}
	//fmt.Printf("\ndocument %p\n", &document)
	//text := document.Find("title").Text()
	//fmt.Println(text)
	document.Find("div#wp_page_numbers ul").Find("li").Each(func(i int, selection *goquery.Selection) {
		if selection.HasClass("thisclass") {
			cate.contents.Push(CategoryNav{"本页", urlString})
		}
		aLable := selection.Find("a")
		//fmt.Println("aLable.Size() ", aLable.Size())
		aLableText := aLable.Text()
		if aLable.Size() == 1{
			_, err := strconv.Atoi(aLableText)
			if err == nil{
				//说明不是特殊链接
				_url, exists := aLable.Attr("href")
				if exists{
					cate.contents.Push(CategoryNav{aLableText, meizitu_cate_base_url + _url})
				}
			}

		}
	})

	ch <- 1
}

func (meizi *Meizitu) fetchAlbumImgs(url string) []dataType{
	var imgs Queue

	docString := meizi.fetchContentByUrl(url)
	document, e := goquery.NewDocumentFromReader(strings.NewReader(docString))
	if e!=nil{
		panic(e)
	}
	document.Find("div#picture p").Find("img").Each(func(i int, selection *goquery.Selection) {
		picUrl, exists := selection.Attr("src")
		if exists{
			name, exists := selection.Attr("alt")
			if exists{
				imgs.Push(Img{name, picUrl})
			}
		}
	})
	return imgs.data
}

func (meizi *Meizitu) fetchCurrentAlbums(url string){
	docString := meizi.fetchContentByUrl(url)
	document, e := goquery.NewDocumentFromReader(strings.NewReader(docString))
	if e!=nil{
		panic(e)
	}
	meizi.currentAlbums.Empty()
	document.Find("ul.wp-list").Find("li").Each(func(i int, selection *goquery.Selection) {
		album := Album{}
		iconUrl, exists := selection.Find("div.pic a img").Attr("src")
		if exists {
			album.icon = iconUrl
		}
		aLable := selection.Find("h3.tit a")
		albumUrl, exists := aLable.Attr("href")
		if exists{
			album.url = albumUrl
		}

		album.name = aLable.Text()

		album.imgs.Push(meizi.fetchAlbumImgs(album.url))
		meizi.currentAlbums.Push(album)
	})
}

func (meizi *Meizitu) seek(page int) {
	//baseUrl := meizi.selectedCategory.url
	//strs := strings.Split(baseUrl, ".")
	//strs = strs[:len(strs)-1]
	//baseUrl = strings.Join(strs, ".") // instead of replace .html with empty string  because suffix can be .htm etc.
	//
	length := meizi.selectedCategory.contents.Size()

	if page<0 || page > length + 1  {
		//无效则置一
		meizi.currentPage = page

	} else if page == 0{
		//上一页
		if meizi.currentPage > 0{
			meizi.currentPage -= 1
		} else {
			meizi.currentPage = length - 1
		}
	} else if page == length + 1 {
		//下一页
		if meizi.currentPage < length - 1{
			meizi.currentPage += 1
		} else {
			meizi.currentPage = 0
		}
	} else {
		meizi.currentPage = page - 1
	}

	cateNav, e := meizi.selectedCategory.getCurrentPage(meizi.currentPage)
	if e != nil{
		fmt.Println(e)
	}
	meizi.fetchCurrentAlbums(cateNav.url)
}

func (cate *Category) getCurrentPage(pageNum int) (cateNav CategoryNav, err error){
	if pageNum >= cate.contents.Size() || pageNum < 0{
		err = errors.New("pageNum " + strconv.Itoa(pageNum) + " out of range")
	} else {
		_cateNav := cate.contents.data[pageNum]
		if __cateNav, ok := _cateNav.(CategoryNav); ok{
			return __cateNav, nil
		} else {
			err = errors.New("data in cate.contents.data of this pageNum is invalid")
		}
	}
	return
}

func (meizi *Meizitu) PageControl(){
	iter := meizi.selectedCategory.contents.Iterator()
	var info string

	info += "输入编号进行分页跳转：\n"
	info += "0：上一页 "

	iter.Each(func(index int, data dataType) {
		info += strconv.Itoa(index+1) + "：第" + strconv.Itoa(index+1) + "页 "
	})

	i :=meizi.selectedCategory.contents.Size()+1
	info += strconv.Itoa(i) + "：下一页"
	for {
		fmt.Println(info)
		var pageNum int

		n, err := fmt.Scanf("%d", &pageNum)
		if n==1 && err ==nil{
			if pageNum > i || pageNum < 0{
				continue
			} else {
				meizi.seek(pageNum)
				break
			}

		} else {
			fmt.Println(err, n)
			return
		}
	}

}

func (meizi *Meizitu) Run() {
	var cateNumber int
	chFetchCategoryContent := make(chan int)

	meizi.fetchMainContent()
	meizi.getCategories()
	for i:=0; i< len(meizi.categories); i++ {
		go meizi.fetchCategoryPages(&meizi.categories[i], chFetchCategoryContent)
	}

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
	meizi.changeCategory(cateNumber)

	//iter := meizi.selectedCategory.contents.Iterator()
	//iter.Each(func(index int, data dataType) {
	//	if c, ok := data.(CategoryNav);ok{
	//		meizi.fetchCurrentAlbums(c.url)
	//	}
	//})
	for{
		meizi.PageControl()

		iterator := meizi.currentAlbums.Iterator()
		iterator.Each(func(index int, data dataType) {
			if c, ok := data.(Album); ok {
				fmt.Println(c)
			}
		})
	}




	defer func() {
		close(chFetchCategoryContent)
	}()
}


func (meizi *Meizitu) changeCategory(cateNumber int) {
	if cateNumber >= 0 && cateNumber < len(meizi.categories) {
		// 输入值有效
		meizi.selectedCategory = &meizi.categories[cateNumber]
	}
}
