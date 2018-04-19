package main

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"net/http"
	"github.com/djimenez/iconv-go"
	"io/ioutil"
	"path"
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"sync"
	"io"
	"strings"
	"log"
	"time"
)

func sohuHandle(url string, char string, categoryId int64, childId int64) {

	doc, err := createDoc(url)
	if err != nil {
		trace(fmt.Sprintf("%s createDoc failed : %v", url, err))
	}

	doc.Find(".list-box li").Each(func(i int, s *goquery.Selection) {
		art := new(AuArticle)
		//title := s.Find(".content-title a").Text()
		href, _ := s.Find(".content-title a").Attr("href")
		img, _ := s.Find(".content-pic img").Attr("src")
		art.CategoryId = categoryId
		art.ChildId = childId
		art.ArticleFocus = img
		art.ArticleThumb = img
		art.FromUrl = href
		sohuDetailHandle(href, char, art)
	})
}

func sohuBatchHandle(src string, char string, categoryId int64, childId int64, start int, end int) {

	for i := start; i <= end; i++ {
		url := fmt.Sprintf(src, i)
		doc, err := createDoc(url)
		if err != nil {
			trace(fmt.Sprintf("%s createDoc failed : %v", url, err))
		}

		doc.Find(".list-box li").Each(func(i int, s *goquery.Selection) {
			art := new(AuArticle)
			//title := s.Find(".content-title a").Text()
			href, _ := s.Find(".content-title a").Attr("href")
			img, _ := s.Find(".content-pic img").Attr("src")
			art.CategoryId = categoryId
			art.ChildId = childId
			art.ArticleFocus = img
			art.ArticleThumb = img
			art.FromUrl = href
			sohuDetailHandle(href, char, art)
		})

	}
}

func sohuDetailHandle(url string, char string, art *AuArticle) {
	doc, err := createDoc(url)
	if err != nil {
		trace(fmt.Sprintf("%s createDoc failed : %v", url, err))
	}
	art.ArticleTitle = doc.Find("div.news-title h1").Text()
	art.Author = doc.Find("span.writer").Text()
	pushTime := doc.Find("span.time").Text()
	art.CreateTime, _ = formatTime(pushTime)
	art.ArticleContent, _ = doc.Find("#contentText").Html()
	art.ArticleDesc = doc.Find("#description").Text()
	art2utf8(art, char)
	if len(art.ArticleTitle) > 0 {
		b, _ := engine.Table("au_article_crawler").Get(&AuArticle{ArticleTitle: art.ArticleTitle})
		if !b {
			art.Status = 2
			art.Type = 2
			art.FromName = "搜狐汽车"
			affected, err := engine.Table("au_article_crawler").Insert(art)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(affected, art.ArticleTitle)
		} else {
			fmt.Println(0, art.ArticleTitle)
		}
	}

}

func art2utf8(art *AuArticle, char string) {
	v := reflect.ValueOf(art).Elem()
	k := v.Type()
	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)
		if value, ok := val.Interface().(string); ok {
			content, _ := iconv.ConvertString(value, char, "utf-8")
			field := v.FieldByName(k.Field(i).Name)
			if field.IsValid() {
				field.SetString(content)
			}
		}
	}
}

func createDoc(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get images failed : %s", err.Error())
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doc create from reader error %s ", err.Error())
	}
	return doc, nil
}

func saveFile(src string) (string, error) {
	resp, err := http.Get(src)
	if err != nil {
		return "", fmt.Errorf("http get images failed : %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fileMd5 := hex.EncodeToString(byte2string(md5.Sum(body)))
	fileSuffix := path.Ext(src)
	url, err := antUploadImage(fileSuffix, body, fileMd5)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return url, nil
}

func articleHandle() {
	var (
		ch chan int = make(chan int, MaxOpenGoRoutine)
		wg sync.WaitGroup
		pageSize int = 10
	)

	count, err := engine.Where("status = ? AND `type` = ?", 1, 2).Count(&AuArticle{})
	if err != nil {
		log.Println("query article count failed :", err)
		return
	}
	page := int(count) / pageSize
	for i := 0; i <= page; i++ {
		offset := pageSize * i
		wg.Add(1)
		go func(offset, pageSize int) {
			defer wg.Done()
			ch <- 1
			getArticleList(offset, pageSize)
			<-ch
		}(offset, pageSize)
	}
	wg.Wait()
}

func getArticleList(offset int, pageSize int) {
	rows, err := engine.Where("status = ? AND `type` = ?", 1, 2).Limit(pageSize, offset).Rows(&AuArticle{})
	if err != nil {
		log.Println("query article failed :", err)
		return
	}

	defer rows.Close()
	art := new(AuArticle)
	for rows.Next() {
		err = rows.Scan(art)
		if err != nil {
			log.Println("get row failed :", err)
		} else {

			if strings.Contains(art.ArticleThumb, "isou365") {
				continue
			}

			//正则匹配的方式
			//reg, _ := regexp.Compile(`http://([^"]+(?:jpg|gif|png|bmp|jpeg))`)
			//result := reg.FindAllString(art.ArticleContent, -1)
			//fmt.Println(b,result)

			//DOM的方式
			var r io.Reader = strings.NewReader(art.ArticleContent)
			doc, err := goquery.NewDocumentFromReader(r)
			if err != nil {
				log.Printf("doc create from reader error %s \n", err.Error())
			} else {
				doc.Find("img").Each(func(i int, s *goquery.Selection) {
					src, _ := s.Attr("src")
					result, err := saveFile(src)
					antUrl := getImgUrl() + "/" + result
					if err != nil {
						log.Printf("%s save failed : %v \n", src, err)
					} else {
						art.ArticleContent = strings.Replace(art.ArticleContent, src, antUrl, -1)
					}
				})

				if len(art.ArticleFocus) > 0 {
					focusImg, err := saveFile(art.ArticleFocus)
					antFocusImg := getImgUrl() + "/" + focusImg
					if err != nil {
						log.Printf("%s save failed : %v \n", antFocusImg, err)
					} else {
						art.ArticleFocus = antFocusImg
					}
				}

				if len(art.ArticleThumb) > 0 {
					thumbImg, err := saveFile(art.ArticleThumb)
					antThumbImg := getImgUrl() + "/" + thumbImg
					if err != nil {
						log.Printf("%s save failed : %v \n", antThumbImg, err)
					} else {
						art.ArticleThumb = antThumbImg
					}
				}

				affected, err := engine.Id(art.ArticleId).Update(art)
				fmt.Println(affected, err)
			}
		}
	}
}


func article(articleId int){
	art := new(AuArticle)
	result, err := engine.Table("au_article_crawler").Id(articleId).Get(art)
	if !result {
		trace(fmt.Sprintf("get article failed : %v", err))
		return
	}
	var r io.Reader = strings.NewReader(art.ArticleContent)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		trace(fmt.Sprintf("doc create from reader error %s ", err.Error()))
	} else {
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			result, err := saveFile(src)
			antUrl := getImgUrl() + "/" + result
			fmt.Println(src, antUrl)
			if err != nil {
				trace(fmt.Sprintf("%s save failed : %v \n", src, err))
			} else {
				art.ArticleContent = strings.Replace(art.ArticleContent, src, antUrl, -1)
			}
		})

		if len(art.ArticleFocus) > 0 {
			focusImg, err := saveFile(art.ArticleFocus)
			if err != nil {
				trace(fmt.Sprintf("%s save failed : %v \n", focusImg, err))
			} else {
				art.ArticleFocus = focusImg
			}
		}

		if len(art.ArticleThumb) > 0 {
			thumbImg, err := saveFile(art.ArticleThumb)
			if err != nil {
				trace(fmt.Sprintf("%s save failed : %v \n", thumbImg, err))
			} else {
				art.ArticleThumb = thumbImg
			}
		}

		if art.CreateTime == 0 {
			art.CreateTime = time.Now().Unix()
		}

		art.ArticleId = 0
		affected, err := engine.Table("au_article").Insert(art)
		trace(fmt.Sprintf("%d, %s, %v", affected, art.ArticleTitle, err))
	}
}