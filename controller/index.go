package controller

import (
	"ReadBook/model"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"strings"
	"time"
)

// 数据结构定义
type NewIn struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Author   string `json:"author"`
}

type NewUp struct {
	Category      string `json:"category"`
	Title         string `json:"title"`
	BookURL       string `json:"book_url"`
	LatestChapter string `json:"latest_chapter"`
	ChapterURL    string `json:"chapter_url"`
	Author        string `json:"author"`
	UpdateTime    string `json:"update_time"`
}

type NewMsg struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type MoreData struct {
	NewIn  []NewIn  `json:"new_in"`
	NewUp  []NewUp  `json:"new_up"`
	NewMsg []NewMsg `json:"new_msg"`
}

type AllData struct {
	Recommend []model.Book     `json:"推荐阅读"`
	Category  []model.BookElse `json:"分类"`
	MoreData                   // 嵌套最新入库/更新/资讯
}

// 合并后的核心解析函数
func parseAllData() (AllData, error) {
	baseURL := "http://wz09.eeyyk.cn"
	var result AllData

	// 创建单个收集器
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.AllowURLRevisit(),
	)

	// 1. 解析推荐阅读
	c.OnHTML("div.panel-body div.row div.col-xs-4.book-coverlist", func(e *colly.HTMLElement) {
		book := model.Book{}
		if coverStyle := e.ChildAttr("div.col-sm-5 a.thumbnail", "style"); coverStyle != "" {
			if start := strings.Index(coverStyle, "url(") + 4; start > 4 {
				if end := strings.Index(coverStyle[start:], ")"); end > 0 {
					book.CoverURL = baseURL + coverStyle[start:start+end]
				}
			}
		}
		e.ForEach("div.caption h4.fs-16 a", func(_ int, el *colly.HTMLElement) {
			book.Title = el.Text
			book.DetailURL = baseURL + el.Attr("href")
		})
		book.Author = e.ChildText("div.caption small.fs-14")
		book.Intro = strings.TrimSpace(e.ChildText("div.caption p.fs-12"))
		if book.Title != "" {
			result.Recommend = append(result.Recommend, book)
		}
	})

	// 2. 解析分类数据
	c.OnHTML("a.list-group-item", func(e *colly.HTMLElement) {
		book := model.BookElse{}
		book.URL = e.Attr("href")
		book.Title = e.Attr("title")
		if book.Title == "" {
			fullText := strings.TrimSpace(e.Text)
			authorText := strings.TrimSpace(e.ChildText("span.pull-right.fs-12"))
			if authorText != "" {
				book.Title = strings.TrimSpace(strings.ReplaceAll(fullText, authorText, ""))
			} else {
				book.Title = fullText
			}
		}
		book.Author = strings.TrimSpace(e.ChildText("span.pull-right.fs-12"))
		if book.Title != "" {
			result.Category = append(result.Category, book)
		}
	})

	// 3. 解析最新入库
	c.OnHTML("div:nth-child(3) > div:nth-child(1) > div:nth-child(1) table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(idx int, row *colly.HTMLElement) {
			cells := row.ChildTexts("td")
			if len(cells) < 3 {
				return
			}
			category := strings.TrimSpace(cells[0])
			bookTitle := row.ChildText("td:nth-child(2) a")
			bookHref := row.ChildAttr("td:nth-child(2) a", "href")
			author := strings.TrimSpace(cells[2])
			if bookHref != "" && bookTitle != "" {
				result.NewIn = append(result.NewIn, NewIn{
					Category: category,
					Title:    bookTitle,
					URL:      baseURL + bookHref,
					Author:   author,
				})
			}
		})
	})

	// 4. 解析最新更新
	c.OnHTML("div:nth-child(3) > div:nth-child(2) > div table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(idx int, row *colly.HTMLElement) {
			cells := row.ChildTexts("td")
			if len(cells) < 5 {
				return
			}
			category := strings.TrimSpace(cells[0])
			bookTitle := row.ChildText("td:nth-child(2) a")
			bookHref := row.ChildAttr("td:nth-child(2) a", "href")
			chapterTitle := row.ChildText("td:nth-child(3) a")
			chapterHref := row.ChildAttr("td:nth-child(3) a", "href")
			author := strings.TrimSpace(cells[3])
			updateTime := strings.TrimSpace(cells[4])
			if bookHref != "" && chapterHref != "" {
				result.NewUp = append(result.NewUp, NewUp{
					Category:      category,
					Title:         bookTitle,
					BookURL:       baseURL + bookHref,
					LatestChapter: chapterTitle,
					ChapterURL:    baseURL + chapterHref,
					Author:        author,
					UpdateTime:    updateTime,
				})
			}
		})
	})

	// 5. 解析最新资讯
	c.OnHTML("div:nth-child(3) > div:nth-child(1) > div:nth-child(2) table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(idx int, row *colly.HTMLElement) {
			title := row.ChildAttr("td a", "title")
			if title == "" {
				title = strings.TrimSpace(row.ChildText("td a"))
			}
			href := row.ChildAttr("td a", "href")
			if href != "" && title != "" {
				result.NewMsg = append(result.NewMsg, NewMsg{
					Title: title,
					URL:   baseURL + href,
				})
			}
		})
	})

	// 调试与错误处理
	c.OnRequest(func(r *colly.Request) {
		log.Printf("正在请求: %s", r.URL)
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("请求错误: %v, 状态码: %d", err, r.StatusCode)
	})

	// 执行单次请求
	if err := c.Visit(baseURL); err != nil {
		return result, fmt.Errorf("访问失败: %w", err)
	}

	return result, nil
}

// 合并后的API接口
func GetAllBooks(c *gin.Context) {
	data, err := parseAllData()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

// XPath辅助函数（保留调试能力）
func getXPath(sel *goquery.Selection) (string, error) {
	if sel == nil {
		return "", nil
	}
	var xpath string
	for {
		parent := sel.Parent()
		if parent.Length() == 0 {
			break
		}
		index := sel.Index() + 1
		tagName := goquery.NodeName(sel)
		xpath = fmt.Sprintf("/%s[%d]%s", tagName, index, xpath)
		sel = parent
	}
	return xpath, nil
}

type BookDetails struct {
	Category   string   `json:"分类"`
	BookCover  string   `json:"book_cover"`
	BookName   string   `json:"book_name"`
	BookAuthor string   `json:"book_author"`
	BookTag1   string   `json:"book_tag1"`
	BookTag2   string   `json:"book_tag2"`
	Intro1     string   `json:"intro1"`
	Intro2     string   `json:"intro2"`
	Intro3     string   `json:"intro3"`
	NewChapter string   `json:"new_chapter"`
	NewTime    string   `json:"new_time"`
	NewList    []string `json:"new_list"`
	Related    string   `json:"相关"`
	Books      []Book   `json:"books"`
}

type Book struct {
	BookLink     string `json:"book_link"`
	CoverImg     string `json:"cover_img"`
	BookName     string `json:"book_name"`
	Author       string `json:"author"`
	Introduction string `json:"introduction"`
}

// GetBookDetails 获取特定书籍的详细信息
func GetBookDetails(c *gin.Context) {
	bookID := c.Param("bookID")
	chapterID := c.Param("chapterID")
	baseURL := "http://wz09.eeyyk.cn"
	bookURL := fmt.Sprintf("%s/book/%s/%s/", baseURL, bookID, chapterID)

	var details BookDetails
	co := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	co.OnHTML("html", func(e *colly.HTMLElement) {
		// 提取分类
		details.Category = strings.TrimSpace(e.DOM.Find("ol.breadcrumb li:nth-child(2) a").Text())

		// 提取书籍基本信息
		panel := e.DOM.Find("div.panel-body .row")
		if panel.Length() > 0 {
			// 修复封面选择器：去掉多余的img标签名
			details.BookCover = baseURL + panel.Find("img.img-thumbnail").AttrOr("src", "")

			// 更准确地定位内容区域
			contentDiv := panel.Find("div.col-sm-10.pl0")

			// 提取基本信息
			details.BookName = contentDiv.Find("h1.bookTitle").Text()
			details.BookAuthor = contentDiv.Find("p.booktag a.red").Text()

			// 提取标签信息
			tags := contentDiv.Find("p.booktag span")
			if tags.Length() >= 2 {
				details.BookTag1 = tags.Eq(0).Text()
				details.BookTag2 = tags.Eq(1).Text()
			}

			// 提取简介内容
			introParagraphs := contentDiv.Find("p.text-justify").NextAll()
			if introParagraphs.Length() >= 1 {
				details.Intro1 = introParagraphs.Eq(0).Text()
			}
			if introParagraphs.Length() >= 2 {
				details.Intro2 = introParagraphs.Eq(1).Text()
			}
			if introParagraphs.Length() >= 3 {
				details.Intro3 = introParagraphs.Eq(2).Text()
			}

			// 检查简介3是否过短
			if len(details.Intro3) < 10 {
				details.Intro3 = "夫为人子者：出必告，反必面，所游必有常，所习必有业。恒言不称老。年长以倍则父事之，十年以长则兄事之，五年以长则肩随之。群居五人，则长者必异席。"
			}
		}

		// 最新章节
		details.NewChapter = e.DOM.Find("a.text-danger").Text()
		details.NewTime = time.Now().Format("2006-01-02")

		// 章节列表
		e.DOM.Find("dl.panel-chapterlist a").Each(func(_ int, s *goquery.Selection) {
			details.NewList = append(details.NewList, strings.TrimSpace(s.Text()))
		})

		// 相关阅读
		details.Related = details.Category + "相关阅读"

		// 相关书籍
		e.DOM.Find("div.col-xs-4.book-coverlist").Each(func(_ int, s *goquery.Selection) {
			book := Book{}
			coverA := s.Find("a.thumbnail")
			if href, ok := coverA.Attr("href"); ok {
				book.BookLink = baseURL + href
			}
			if style, ok := coverA.Attr("style"); ok && strings.Contains(style, "url(") {
				book.CoverImg = baseURL + strings.Trim(style[strings.Index(style, "url(")+4:], ")")
			}

			titleA := s.Find("h4.fs-16 a")
			if titleA.Length() > 0 {
				book.BookName = titleA.Text()
				if href, ok := titleA.Attr("href"); ok {
					book.BookLink = baseURL + href
				}
			}

			book.Author = s.Find("small.fs-14").Text()
			book.Introduction = s.Find("p.fs-12").Text()
			details.Books = append(details.Books, book)
		})
	})

	co.OnError(func(r *colly.Response, err error) {
		log.Printf("请求错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
	})

	if err := co.Visit(bookURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, details)
}
