package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/cheggaaa/pb"
)

var filepath *string = flag.String("filepath", "maps.txt", "default is maps.txt")

func DownloadFile(srcURL string, dest string) error {
	file := path.Base(srcURL)
	log.Printf("Downloading file %s from %s\n", file, srcURL)

	var path bytes.Buffer
	path.WriteString(dest)
	path.WriteString("/")
	path.WriteString(file)

	out, err := os.Create(path.String())
	if err != nil {
		return err
	}
	defer out.Close()

	res, err := http.Get(srcURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bar := pb.New64(res.ContentLength).Start()
	defer bar.Finish()

	r := bar.NewProxyReader(res.Body)

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	log.Printf("Download completed!")

	return nil
}

func main() {
	flag.Parse()

	fp, err := os.Open(*filepath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Failed read \"%s\" file!\n", *filepath)
		os.Exit(1)
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 4096)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
			os.Exit(1)
		}

		er := DownloadFile(string(line), "./")
		if er != nil {
			fmt.Fprintln(os.Stdout, er)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
