package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>{{ .Title }}</title>
</head>
<body>
{{ .Body }}
</body>
<p>Filename: {{ .Filename }}</p>
</html>
`
)

// content type represents the HTML content to add into the template
type content struct {
	Title    string
	Body     template.HTML
	Filename string
}

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter file name:")
		filename, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		filename = strings.TrimSpace(filename)
		if err := run(filename, os.Stderr, *skipPreview); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if err := run(*filename, os.Stderr, *skipPreview); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func run(filename string, out io.Writer, skipPreview bool) error {
	// Read all the data from the input file and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	htmlData, err := parseContent(input, filename)
	if err != nil {
		return err
	}

	temp, err := os.CreateTemp("./temp", "mdp*.html")
	if err != nil {
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err = saveHTML(outName, htmlData); err != nil {
		return err
	}
	if skipPreview {
		return nil
	}
	defer os.Remove(outName)
	return preview(outName)
}

func parseContent(input []byte, filename string) ([]byte, error) {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Allow user to choose default template by setting environment variable
	tmpl := defaultTemplate
	value, ok := os.LookupEnv("DEFAULT_TEMPLATE")
	if ok {
		data, err := os.ReadFile(value)
		if err != nil {
			return nil, err
		}
		tmpl = string(data)
	}

	// Parse the contents of the defaultTemplate const into a new Template
	t, err := template.New("mdp").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	// Instantiate the content type, adding the title and body
	c := content{
		Title:    "Markdown Preview Tool",
		Body:     template.HTML(body),
		Filename: filename,
	}

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// Execute the template with the content type
	if err = t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func saveHTML(outFname string, data []byte) error {
	// Write the bytes to the file
	return os.WriteFile(outFname, data, 0644)
}

func preview(fname string) error {
	cName := ""
	var cParams []string
	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "wslview"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}
	// Append filename to parameters slice
	cParams = append(cParams, fname)
	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	// Open the file using default program
	err = exec.Command(cPath, cParams...).Run()
	// Give the browser some time to open the file before deleting it
	time.Sleep(1 * time.Second)
	return err
}
