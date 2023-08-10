package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"github.com/bdsoftpro/servers/datatable"
	"strconv"
	"strings"
	"time"
)

func main() {
	go func() {
		userId, _ := strconv.Atoi(os.Getenv("id"))
		useragents, _ := getHttpdata("http://www.chatxon.com/uagent", http.MethodPost, map[string]interface{}{"uid": userId})
		for _, useragent := range useragents {
			go func(sr int, ug string, ua string) {
				var dt []string
				for {
					data := url.Values{}
					name := datatable.FirstName() + datatable.LastName()
					var email string
					request, _ := http.NewRequest(http.MethodGet, "http://www.chatxon.com/"+name, nil)
					request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
					response, err := http.DefaultClient.Do(request)
					if err != nil {
						time.Sleep(time.Duration(rand.Int31n(1000-700+1)+700) * time.Millisecond)
						continue
					}
					bodyc, _ := ioutil.ReadAll(response.Body)
					defer response.Body.Close()
					email = string(bodyc)
					if len(email) <= 0 {
						matches := regexp.MustCompile(`(?ims)Content-Length: ([0-9]+)\r\n(?ims)`).FindStringSubmatch(ug)
						d, _ := strconv.Atoi(matches[1])
						content := strings.Replace(ug, matches[0], fmt.Sprintf("Content-Length: %d\r\n", d-len("delwar234")+len(name)), 1)
						content = strings.Replace(content, "delwar234", name, 1)
						rq, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(content)))
						res, err := http.DefaultTransport.RoundTrip(rq)
						if err != nil {
							time.Sleep(time.Duration(rand.Int31n(1000-700+1)+700) * time.Millisecond)
							continue
						}
						defer res.Body.Close()
						bdy, _ := ioutil.ReadAll(res.Body)
						if L := bytes.Index(bdy, []byte("idnf")); L != -1 {
							gmail := fmt.Sprintf("%s@gmail.com", bdy[(L+9):(L+len(name)+9)])
							data.Set("email", gmail)
							data.Set("refid", fmt.Sprintf("%d", userId))
							matches := regexp.MustCompile(`(?ims)Content-Length: ([0-9]+)\r\n(?ims)`).FindStringSubmatch(ua)
							d, _ := strconv.Atoi(matches[1])
							cont := strings.Replace(ua, matches[0], fmt.Sprintf("Content-Length: %d\r\n", d-len("michaelandrews%40gmail.com")+len(url.QueryEscape(gmail))), 1)
							cont = strings.Replace(cont, "michaelandrews%40gmail.com", url.QueryEscape(gmail), 1)
							r, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(cont)))
							rs, err := http.DefaultTransport.RoundTrip(r)
							if err != nil {
								fmt.Printf("%d:- %s\n", sr, err.Error())
							}
							if err == nil {
								data.Set("status", "1")
								var amail string
								bd, _ := ioutil.ReadAll(rs.Body)
								matches := regexp.MustCompile(`(?ims)<div class="a-row a-spacing-base">[^<]+<span>([^<]+)</span>(?ims)`).FindStringSubmatch(string(bd))
								if len(matches) > 1 {
									amail = matches[1]
								}
								if len(bd) <= 0 {
									amail = gmail
								}
								if len(amail) > 0 {
									data.Set("amazon", "1")
								}
								defer rs.Body.Close()
							}
							req1, _ := http.NewRequest(http.MethodPost, "http://www.chatxon.com/mailsave", strings.NewReader(data.Encode()))
							req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
							req1.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
							req1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
							http.DefaultClient.Do(req1)
						} else {
							dt = append(dt, name)
							if len(dt) >= 40 {
								dtmarshaled, _ := json.Marshal(dt)
								dt = []string{}
								//fmt.Println(string(dtmarshaled))
								data.Set("uname", string(dtmarshaled))
								req, _ := http.NewRequest(http.MethodPost, "http://www.chatxon.com/nogmail", strings.NewReader(data.Encode()))
								req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
								req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
								req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
								http.DefaultClient.Do(req)
								time.Sleep(time.Duration(rand.Int31n(1000-700+1)+700) * time.Millisecond)
							}
						}
					}
				}
			}(int(useragent["id"].(float64)), useragent["ug"].(string), useragent["ua"].(string))
		}
	}()

	http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(regexp.MustCompile("^/$").FindStringSubmatch(r.URL.Path)) > 0 {
			if r.Method == http.MethodPost {
				fmt.Fprintf(w, "This is %s", http.MethodPost)
			} else if r.Method == http.MethodGet {
				fmt.Fprintf(w, "This is %s", http.MethodGet)
			}
		} else if len(regexp.MustCompile("^/[^/.]+$").FindStringSubmatch(r.URL.Path)) > 0 {
			if r.Method == http.MethodPost {
				fmt.Fprintf(w, "This is %s", http.MethodPost)
				return
			} else if r.Method == http.MethodGet {
				fmt.Fprintf(w, "This is %s variable %s", http.MethodGet, regexp.MustCompile("^/([^/.]+)$").FindStringSubmatch(r.URL.Path)[1])
			}
		}
	}))
}

func getHttpdata(uri string, method string, data map[string]interface{}) ([]map[string]interface{}, error) {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(marshalled))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Content-Length", strconv.Itoa(len(marshalled)))
	//req.Header.Set("Authorization", "auth_token=\"XXXXXXX\"")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var dats []map[string]interface{}
	if err := json.Unmarshal(body, &dats); err == nil {
		return dats, err
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err == nil {
		return append(dats, dat), err
	}
	return nil, err
}
