package v1

import (
	"adams549659584/go-proxy-bingai/common"
	"net/http"
	"os"
	"strings"
	"time"

	binglib "github.com/Harry-zklcdc/bing-lib"
	"github.com/Harry-zklcdc/bing-lib/lib/aes"
	"github.com/Harry-zklcdc/bing-lib/lib/hex"
	"github.com/Harry-zklcdc/bing-lib/lib/request"
)

func init() {
	VERCEL_URL := os.Getenv("VERCEL_URL")
	if VERCEL_URL == "" {
		go func() {
			time.Sleep(200 * time.Millisecond)
			t, _ := getCookie("", "", "")
			common.Logger.Info("BingAPI Ready!")
			globalChat = binglib.NewChat(t).SetBingBaseUrl("http://localhost:" + common.PORT).SetSydneyBaseUrl("ws://localhost:" + common.PORT).SetBypassServer(common.BypassServer)
			globalImage = binglib.NewImage(t).SetBingBaseUrl("http://localhost:" + common.PORT).SetBypassServer(common.BypassServer)
		}()
	} else {
		t, _ := getCookie("", "", "")
		globalChat = binglib.NewChat(t).SetBypassServer(common.BypassServer)
		globalImage = binglib.NewImage(t).SetBypassServer(common.BypassServer)
	}
}

func getCookie(reqCookie, convId, rid string) (cookie string, err error) {
	cookie = reqCookie
	if common.AUTH_KEY != "" {
		cookie += "; " + common.AUTH_KEY_COOKIE_NAME + "=" + common.AUTH_KEY
	}
	c := request.NewRequest()
	var res *request.Client
	if os.Getenv("VERCEL_URL") == "" {
		res = c.SetUrl("http://localhost:"+common.PORT+"/chat?q=Bing+AI&showconv=1&FORM=hpcodx&ajaxhist=0&ajaxserp=0&cc=us").
			SetHeader("User-Agent", common.User_Agent).
			SetHeader("Cookie", cookie).Do()
	} else {
		res = c.SetUrl("https://www.bing.com/chat?q=Bing+AI&showconv=1&FORM=hpcodx&ajaxhist=0&ajaxserp=0&cc=us").
			SetHeader("User-Agent", common.User_Agent).
			SetHeader("Cookie", cookie).Do()

	}
	headers := res.GetHeaders()
	for k, v := range headers {
		if strings.ToLower(k) == "set-cookie" {
			for _, i := range v {
				cookie += strings.Split(i, "; ")[0] + "; "
			}
		}
	}
	cookie = strings.TrimLeft(strings.Trim(cookie, "; "), "; ")
	IG := strings.ToUpper(hex.NewHex(32))
	T, err := aes.Encrypt(common.AUTHOR, IG)
	if err != nil {
		return
	}
	resp, status, err := binglib.Bypass(common.BypassServer, reqCookie, "local-gen-"+hex.NewUUID(), IG, convId, rid, T)
	if err != nil || status != http.StatusOK {
		common.Logger.Error("Bypass Error: %v", err)
		return
	}
	cookie = resp.Result.Cookies
	if len(common.USER_TOKEN_LIST) == 0 {
		cookie += "; _U=" + hex.NewHex(128)
	}
	if common.AUTH_KEY != "" {
		cookie += "; " + common.AUTH_KEY_COOKIE_NAME + "=" + common.AUTH_KEY
	}
	return cookie, nil
}
