package scraper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"github.com/gocolly/colly"
)

type Chapter struct {
	chapterNum string
	cdnUrls    []string
}

func ScrapeChapter(chapterUrl string) []string {
	cdns := []string{}

	imageScraper := colly.NewCollector()

	//define callback functions
	/*
		imageScraper.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})
	*/

	imageScraper.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	/*

		imageScraper.OnResponse(func(r *colly.Response) {
			fmt.Println("Visited", r.Request.URL)
		})
	*/

	imageScraper.OnHTML("picture", func(e *colly.HTMLElement) {
		cdns = append(cdns, e.ChildAttr("img", "src"))
	})
	imageScraper.Visit(chapterUrl)
	return cdns
}

func DownloadImage(url string, chapterNum string, pageNum string) {
	os.Chdir("../manga")
	img, _ := os.Create(chapterNum + "_" + pageNum + ".jpg")
	defer img.Close()

	test, err := http.Get(url)
	if err != nil {
		fmt.Println("Error")
		return
	}
	defer test.Body.Close()

	io.Copy(img, test.Body)
}

func FindChapterUrl(mangaName string, mangaUrl string, chapter string) string {
	var chapterUrl string
	sc := colly.NewCollector()

	sc.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
	/*
		sc.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		sc.OnResponse(func(r *colly.Response) {
			fmt.Println("Visited", r.Request.URL)
		})
	*/
	sc.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.ForEach("div", func(_ int, el *colly.HTMLElement) {

			if el.Text == mangaName+"  Chapter "+chapter {
				chapterUrl = e.Attr("href")
			}
		})
	})

	sc.Visit(mangaUrl)
	return ("https://properlinker.com" + chapterUrl)
}

func DownloadChapter(chapterNum string, chapterUrl string) {
	chapter := Chapter{chapterNum: chapterNum}
	chapter.cdnUrls = ScrapeChapter(chapterUrl)
	for i, v := range chapter.cdnUrls {
		DownloadImage(v, chapter.chapterNum, strconv.Itoa(i))
	}
}

func Start(chapterNum string, mangaName string) {
	var mangaUrls map[string]string = map[string]string{
		"One Piece": "https://properlinker.com/mangas/5/one-piece",
	}

	//io

	fmt.Print("Enter the chapter number: ")
	fmt.Scanln(&chapterNum)
	DownloadChapter(chapterNum, FindChapterUrl(mangaName, mangaUrls[mangaName], chapterNum))
	//loading bar
	fmt.Println("Done! Downloaded chapter " + chapterNum)
}
