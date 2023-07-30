package middleware

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc { //membutuhkan next echo.HandlerFunc) echo.HandlerFunc
	return func(c echo.Context) error {
		file, err := c.FormFile("input-project-image") // tangkap file dari input

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		
         // buka file
		sumber, err := file.Open() //untuk buka file

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		fmt.Println("sumber:", sumber)

		defer sumber.Close() // lifo fungsi nya untuk jalankan di akhir tapi bebas di taro mau di atas => manfaat agar tidak kebocoran memory

		tpFile, err := ioutil.TempFile("uploads", "image-*.png") //mnetapkan nama file

		defer tpFile.Close()

		fmt.Println("tpFile", tpFile)

		kopi, err := io.Copy(tpFile, sumber)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		fmt.Println("kopi:", kopi)


		data := tpFile.Name()
		filename := data[8:]
		fmt.Println("nama utuh", filename)

		c.Set("dataFile", filename)

		return next(c) // lnjutin func echo context handle func

	
	}
}
