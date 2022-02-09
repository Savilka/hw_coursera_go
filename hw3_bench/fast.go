package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"io"
	"os"
	"regexp"
	"sync"
)

// go test -bench . -benchmem -cpuprofile cpu.out -memprofile mem.out -memprofilerate=1
// go tool pprof  hw3_bench.test.exe cpu.out
// go tool pprof  hw3_bench.test.exe mem.out
// вам надо написать более быструю оптимальную этой функции

type user struct {
	Browsers []string `json:"browsers,intern"`
	Company  string   `json:"company,intern"`
	Country  string   `json:"country,intern"`
	Email    string   `json:"email,intern"`
	Job      string   `json:"job,intern"`
	Name     string   `json:"name,intern"`
	Phone    string   `json:"phone,intern"`
}

var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson42239ddeDecodeHw3BenchJson(in *jlexer.Lexer, out *user) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.StringIntern())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "company":
			out.Company = string(in.StringIntern())
		case "country":
			out.Country = string(in.StringIntern())
		case "email":
			out.Email = string(in.StringIntern())
		case "job":
			out.Job = string(in.StringIntern())
		case "name":
			out.Name = string(in.StringIntern())
		case "phone":
			out.Phone = string(in.StringIntern())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson42239ddeEncodeHw3BenchJson(out *jwriter.Writer, in user) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		out.RawString(prefix[1:])
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"company\":"
		out.RawString(prefix)
		out.String(string(in.Company))
	}
	{
		const prefix string = ",\"country\":"
		out.RawString(prefix)
		out.String(string(in.Country))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"job\":"
		out.RawString(prefix)
		out.String(string(in.Job))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"phone\":"
		out.RawString(prefix)
		out.String(string(in.Phone))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v user) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson42239ddeEncodeHw3BenchJson(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v user) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson42239ddeEncodeHw3BenchJson(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *user) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson42239ddeDecodeHw3BenchJson(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *user) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson42239ddeDecodeHw3BenchJson(l, v)
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func FastSearch(out io.Writer) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	//reader := bufio.NewReader(file)

	r := regexp.MustCompile("@")
	android := regexp.MustCompile("Android")
	msie := regexp.MustCompile("MSIE")
	var seenBrowsers []string
	uniqueBrowsers := 0

	scanner := bufio.NewScanner(file)

	users := make([]user, 1000)
	var userStruct user
	idx := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		err := easyjson.Unmarshal(line, &userStruct)
		if err != nil {
			panic(err)
		}

		users[idx] = userStruct
		users[idx].Browsers = make([]string, len(userStruct.Browsers))
		copy(users[idx].Browsers, userStruct.Browsers)

		idx++
	}
	//for {
	//	//line, err = reader.ReadBytes('\n')
	//	line := bufPool.Get().([]byte)
	//	line = append(line, reader.ReadBytes('\n'))
	//	if err == io.EOF {
	//		break
	//	}
	//
	//	// fmt.Printf("%v %v\n", err, line)
	//	err := easyjson.Unmarshal(line, &userStruct)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	users[idx] = userStruct
	//	users[idx].Browsers = make([]string, len(userStruct.Browsers))
	//	copy(users[idx].Browsers, userStruct.Browsers)
	//
	//	idx++
	//}
	fmt.Fprintln(out, "found users:")
	for i, user := range users {

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

		for _, browserRaw := range browsers {
			browser := browserRaw
			if ok := android.MatchString(browser); ok {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browserRaw := range browsers {
			browser := browserRaw

			if ok := msie.MatchString(browser); ok {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := r.ReplaceAllString(user.Email, " [at] ")
		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
