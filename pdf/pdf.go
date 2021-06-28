package pdf

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"

	_ "image/gif"
	_ "image/png"

	"github.com/lirix360/ReadmangaGrabber/data"

	"github.com/phpdave11/gofpdf"
)

// CreatePDF - ...
func CreatePDF(chapterPath string, savedFiles []string, wg *sync.WaitGroup) error {
	defer wg.Done()

	var opt gofpdf.ImageOptions

	pdf := gofpdf.New("P", "mm", "A4", "")

	for _, file := range savedFiles {
		width, height := resizeToFit(file)

		imageFile, err := convertImg(file)
		if err != nil {
			data.WSChan <- data.WSData{
				Cmd: "updateLog",
				Payload: map[string]interface{}{
					"type": "err",
					"text": "-- Ошибка при создании PDF файла (" + chapterPath + ".pdf):" + err.Error(),
				},
			}
			return err
		}

		if width < height {
			pdf.AddPage()
			pdf.ImageOptions(imageFile, (data.PDFOpts.A4Width-width)/2, (data.PDFOpts.A4Height-height)/2, width, height, false, opt, 0, "")
		} else {
			pdf.AddPageFormat("L", pdf.GetPageSizeStr("A4"))
			pdf.ImageOptions(imageFile, (data.PDFOpts.A4Height-width)/2, (data.PDFOpts.A4Width-height)/2, width, height, false, opt, 0, "")
		}
	}

	err := pdf.OutputFileAndClose(chapterPath + ".pdf")
	if err != nil {
		log.Println("Ошибка создания PDF файла ("+chapterPath+".pdf):", err.Error())
		data.WSChan <- data.WSData{
			Cmd: "updateLog",
			Payload: map[string]interface{}{
				"type": "err",
				"text": "-- Ошибка при создании PDF файла (" + chapterPath + ".pdf):" + err.Error(),
			},
		}
		return err
	}

	err = os.RemoveAll(chapterPath + "/pdf")
	if err != nil {
		log.Println(err)
	}

	return nil
}

func resizeToFit(imgFilename string) (float64, float64) {
	var widthScale, heightScale float64

	width, height := getImageDimension(imgFilename)

	if width < height {
		widthScale = data.PDFOpts.MaxWidth / width
		heightScale = data.PDFOpts.MaxHeight / height
	} else {
		widthScale = data.PDFOpts.MaxHeight / width
		heightScale = data.PDFOpts.MaxWidth / height
	}

	scale := math.Min(widthScale, heightScale)

	return math.Round(pixelsToMM(scale * width)), math.Round(pixelsToMM(scale * height))
}

func convertImg(srcImg string) (string, error) {
	srcPath := filepath.Dir(srcImg)
	dstPath := filepath.Join(srcPath, "pdf")
	srcFile := filepath.Base(srcImg)
	dstFile := filepath.Join(dstPath, srcFile+".jpg")

	imgFile, _ := os.Open(srcImg)

	imgSrc, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("-- Skipping file ("+srcImg+"):", err)
		imgFile.Close()
		return "", err
	}

	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		os.MkdirAll(dstPath, 0755)
	}

	newImg := image.NewRGBA(imgSrc.Bounds())

	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)

	jpgFile, _ := os.Create(dstFile)

	var opt jpeg.Options
	opt.Quality = 90

	err = jpeg.Encode(jpgFile, newImg, &opt)
	if err != nil {
		log.Println(err)
		imgFile.Close()
		jpgFile.Close()
		return "", err
	}

	imgFile.Close()
	jpgFile.Close()

	return dstFile, nil
}

func getImageDimension(imagePath string) (float64, float64) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Fatalln("TEST", err)
	}

	return float64(image.Width), float64(image.Height)
}

func pixelsToMM(val float64) float64 {
	return float64(val * data.PDFOpts.MmInInch / data.PDFOpts.DPI)
}