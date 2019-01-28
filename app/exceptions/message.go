package exceptions

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func Message() gin.HandlerFunc {

	return MessageToBrowser(gin.DefaultErrorWriter)
}

func MessageToBrowser(out io.Writer) gin.HandlerFunc {

	return func(c *gin.Context) {

		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				stack := stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					exceptionMsg := fmt.Sprintf("%s\n%s", err, string(httprequest))
					c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", render(exceptionMsg))
				} else if gin.IsDebugging() {
					exceptionMsg := fmt.Sprintf("%s\n%s", err, stack)
					c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", render(exceptionMsg))
				} else {
					exceptionMsg := fmt.Sprintf("%s\n%s", err, stack)
					c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", render(exceptionMsg))
				}
			}
		}()

		c.Next()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func render(content string) []byte {

	style := `
		<style>
			body {
				padding: 32px;
				background-color: black;
				color: white;
				font-size: 18px;
			}

			p {
				margin: 6px;
			}

			.path {
				color: #BCD42A;
			}

			.line {
				color: red;
			}

			.padding-width {
				padding: 0px 20px;
			}
		</style>`

	data := strings.Split(content, "\n")

	for idx, value := range data {
		if strings.HasPrefix(value, "/") {
			if strings.Contains(value, ".go:") || regexp.MustCompile(`\:\d+\s\(`).MatchString(value) {
				tmp := value

				tmp = regexp.MustCompile(`^\/`).ReplaceAllString(tmp, "<span class=\"path\">/")
				tmp = regexp.MustCompile(`\:`).ReplaceAllString(tmp, "</span>&nbsp;&nbsp;&nbsp;<span class=\"line\"><b>line:")
				tmp = regexp.MustCompile(`\s\(`).ReplaceAllString(tmp, "</b></span>&nbsp;&nbsp;&nbsp;(")

				data[idx] = tmp
			}
		}
	}

	content = strings.Join(data, "\n")
	content = strings.Replace("<p>"+content+"</p>", "\n", "</p><p>", -1)
	content = strings.Replace(content, "\t", "<span class=\"padding-width\"></span>", -1)
	content = strings.Replace(content, "\t", "<span class=\"padding-width\"></span>", -1)

	return []byte(style + content)
}
