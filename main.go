package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/fsnotify/fsnotify"
)

var fileName string

func initWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	fmt.Println("🔭  Watching for changes in the HTML template")
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					buildPDF()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("e rror:", err)
			}
		}
	}()
	err = watcher.Add("html/" + fileName + ".html")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func buildPDF() {
	fmt.Println("🤖  Building PDF...")
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	pdfg.PageSize.Set(wkhtml.PageSizeA4)
	pdfg.Dpi.Set(300)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginTop.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.Orientation.Set(wkhtml.OrientationLandscape)

	htmlFile, err := ioutil.ReadFile("html/" + fileName + ".html")
	if err != nil {
		log.Fatal(err)
	}

	pdfg.AddPage(wkhtml.NewPageReader(bytes.NewBuffer(htmlFile)))
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile("pdf/" + fileName + ".pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅  Done!")
}

func prepareAssets() {
	fmt.Println("📁  Preparing assets path")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("📁  Curent working path: %s\n", dir)

	input, err := ioutil.ReadFile("html/" + fileName + ".html")
	if err != nil {
		log.Fatalln(err)
	}

	var re = regexp.MustCompile(`-PWD-`)
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "-PWD-") {
			fmt.Printf("%s contains PWD\n", line)
			lines[i] = re.ReplaceAllString(line, dir)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("html/"+fileName+".html", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.StringVar(&fileName, "filename", "", "HTML template filename, generated PDF output will be of the same name")
	flag.Parse()

	if len(fileName) == 0 {
		log.Fatal("Provide -filename flag")
	}

	fmt.Println("🚀  Started!")
	fmt.Printf("📁  Filename is set to '%s.html'\n", fileName)

	prepareAssets()
	buildPDF()
	initWatcher()
}
