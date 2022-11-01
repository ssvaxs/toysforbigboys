package toysforbigboys

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ssvaxs/rosreestrcn"
)

type (
	Command struct {
		URL   string
		Metod string
		Descr string
	}

	CommandList []Command

	inFile []byte
)

func (f *inFile) GetCN() ([]string, error) {
	if len(*f) == 0 {
		return nil, fmt.Errorf("No data in a file")
	}

	return strings.Split(fmt.Sprint(*f), "\n"), nil
}

func main() {
	s := CommandList{
		{URL: "/info_rosreestr_by_cn", Metod: "POST", Descr: "Gets a txt-file 'file' with a list of cadastral numbers (separated by enter) " +
			"and returns a xslx-file with a data"},
	}

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.LoadHTMLGlob("templates/*")
	// start page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"MyList": s,
		})
	})

	// a request to rosreestr by cadastral numbers
	r.POST(s[0].URL, func(c *gin.Context) {
		f, err := c.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}

		file, err := f.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		bFile := new([]byte)
		_, err = file.Read(*bFile)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		rosreestrcn.RosreestrData(&inFile(*bFile))
		c.JSON(http.StatusOK, gin.H{
			"message": "Your file has been successfully uploaded.",
		})
	})

	r.Run()
}
