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

func DownloadImage(mangaName string, url string, chapterNum string, pageNum string) {
	os.Chdir("./manga")
	img, _ := os.Create(mangaName + "_" + chapterNum + "_" + pageNum + ".jpg")
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

			if el.Text == mangaName+" Chapter "+chapter {
				chapterUrl = e.Attr("href")
			}
		})
	})

	sc.Visit(mangaUrl)
	return ("https://properlinker.com" + chapterUrl)
}

func DownloadChapter(mangaName string, chapterNum string, chapterUrl string) {
	chapter := Chapter{chapterNum: chapterNum}
	chapter.cdnUrls = ScrapeChapter(chapterUrl)
	for i, v := range chapter.cdnUrls {
		DownloadImage(mangaName, v, chapter.chapterNum, strconv.Itoa(i))
	}
}

func Start(chapterNum string, mangaName string) {
	var mangaUrls map[string]string = map[string]string{
		"One Piece ":                      "https://properlinker.com/mangas/5/one-piece",
		"Chainsaw Man":                    "https://properlinker.com/mangas/13/chainsaw-man",
		"Attack on Titan":                 "https://properlinker.com/mangas/8/attack-on-titan",
		"Black Clover":                    "https://properlinker.com/mangas/3/black-clover",
		"Bleach":                          "https://properlinker.com/mangas/2/bleach",
		"Demon Slayer: Kimetsu no Yaiba ": "https://properlinker.com/mangas/19/demon-slayer-kimetsu-no-yaiba",
		"Hunter X Hunter":                 "https://properlinker.com/mangas/15/hunter-x-hunter",
		"Jujutsu Kaisen":                  "https://properlinker.com/mangas/4/jujutsu-kaisen",
		"My Hero Academia":                "https://properlinker.com/mangas/6/my-hero-academia",
		"Spy X Family":                    "https://properlinker.com/mangas/23/spy-x-family",
	}

	//io

	DownloadChapter(mangaName, chapterNum, FindChapterUrl(mangaName, mangaUrls[mangaName], chapterNum))
	//loading bar
}
