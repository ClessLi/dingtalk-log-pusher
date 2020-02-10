package main

import (
	"flag"
	"fmt"
	"github.com/hpcloud/tail"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	filepath = flag.String("f", `app.log`, "log `filepath`")
	env      = flag.String("e", "", "Name of task execution `env`ironment")
	token    = flag.String("t", "", "`token` for dingTalk openAPI")
	hostname = flag.String("h", "", "`hostname` of dingTalk openAPI")
	port     = flag.String("p", "443", "`port` of dingTalk openAPI")
)

//var env = "SIT"

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

func main() {
	flag.Parse()

	isExist, pathErr := PathExists(*filepath)
	if !isExist {
		if pathErr != nil {
			fmt.Println("The logfile", *filepath, "is not found.")
		} else {
			fmt.Println("Unkown error of the logfile.")
		}
		flag.Usage()
		os.Exit(1)
	}

	if *env == "" || *token == "" || *hostname == "" {
		flag.Usage()
		os.Exit(1)
	}
	//filepath := `test.log`

	cert := ConvertByte2String([]byte(`证书即将到期的应用有`), UTF8)

	dingTalkAPIURL := `https://` + *hostname + `:` + *port + `/robot/send?access_token=` + *token
	logRegStr := `^\s*(\d{4}-\d{2}-\d{2}\s*\d{2}:\d{2}:\d*).*` + cert + `\[(.*)\]`

	logReg, regerr := regexp.Compile(logRegStr)

	if regerr != nil {
		fmt.Println(regerr)
		os.Exit(1)
	}

	tails, err := tail.TailFile(*filepath, tail.Config{
		//tails, err := tail.TailFile(filepath, tail.Config{
		ReOpen: true,
		Follow: true,
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	var msg *tail.Line
	var ok bool

	for true {
		msg, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(10 * time.Millisecond)
			continue
		}

		match := logReg.FindStringSubmatch(msg.Text)
		if len(match) != 3 {
			fmt.Println("match err, msg:", msg.Text)
			continue
		}
		t := match[1]
		m := match[2]

		dret, err := sendDingTalkNotification(dingTalkAPIURL, t, m)

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(time.Now(), "-", dret)

	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		return false, nil
	}
}

func sendDingTalkNotification(url string, logTime string, data string) (string, error) {
	context := `{"msgtype": "text",
		"text": {"content": "` + *env + `环境portal后管日志：\n\t\t日志时间：` + logTime + `\n\t\t证书即将到期的应用有：[ ` + data + ` ]"}
	}`
	//context := `{"msgtype": "text",
	//	"text": {"content": "` + env + `环境portal后管日志：\n\t\t日志时间：` + logTime + `证书即将到期的应用有：[` + data + `]"}
	//}`

	req, err := http.NewRequest("POST", url, strings.NewReader(context))
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return "resonse error", err
	}

	defer resp.Body.Close()

	body, rerr := ioutil.ReadAll(resp.Body)
	if rerr != nil {
		return "HttpPOST error", rerr
	}

	return string(body), nil
}
