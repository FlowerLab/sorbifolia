package gomod

import (
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
)

const inputHTML = `
<!doctype html><html lang="en">
<head>
<meta charset="UTF-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="go-import" content="{{.Main}} git {{.Repo}}">
<meta name="go-source" content="{{.Main}} {{.Repo}} {{.Repo}}/tree/{{.Branch}}{/dir} {{.Repo}}/blob/{{.Branch}}{/dir}/{file}#L{line}">
<title>{{.PkgName}} - Golang Package</title>
<script src="https://cdn.jsdelivr.net/npm/marked@4.0.14/marked.min.js"></script>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/github-markdown-css@5.1.0/github-markdown.min.css">
<style>
body {
	box-sizing: border-box;
	min-width: 200px;
	max-width: 980px;
	margin: 0 auto;
	padding: 45px;
	background-color: #F6F6F6;
}

#content {
  padding: 30px 60px 55px;
}

.markdown-body {
	background-color: #FFFFFF;
}

.github-corner:hover .octo-arm {
animation:octocat-wave 560ms ease-in-out
}

@keyframes octocat-wave{
0%,100%{transform:rotate(0)}20%,60%{transform:rotate(-25deg)}40%,80%{transform:rotate(10deg)}
}

@media (max-width:500px){
.github-corner:hover 
.octo-arm{animation:none}
.github-corner .octo-arm{
animation:octocat-wave 560ms ease-in-out
}
}

</style>
</head>
<body>


<a href="https://pkg.go.dev/{{.PkgName}}" class="github-corner" aria-label="View doc on GoDev">
    <svg width="80" height="80" viewBox="0 0 250 250"
         style="fill:#fff; color:#151513; position: absolute; top: 0; border: 0; left: 0;"
         aria-hidden="true">
        <path d="M40.2,101.1c-0.4,0-0.5-0.2-0.3-0.5l2.1-2.7c0.2-0.3,0.7-0.5,1.1-0.5l35.7,0c0.4,0,0.5,0.3,0.3,0.6l-1.7,2.6
    c-0.2,0.3-0.7,0.6-1,0.6L40.2,101.1z" fill="currentColor" style="transform-origin: 100px 1px;" class="octo-arm"/>
        <path d="M25.1,110.3c-0.4,0-0.5-0.2-0.3-0.5l2.1-2.7c0.2-0.3,0.7-0.5,1.1-0.5l45.6,0c0.4,0,0.6,0.3,0.5,0.6l-0.8,2.4
					c-0.1,0.4-0.5,0.6-0.9,0.6L25.1,110.3z" fill="currentColor" style="transform-origin: 100px 1px;"
              class="octo-arm"/>
        <path d="M49.3,119.5c-0.4,0-0.5-0.3-0.3-0.6l1.4-2.5c0.2-0.3,0.6-0.6,1-0.6l20,0c0.4,0,0.6,0.3,0.6,0.7l-0.2,2.4
					c0,0.4-0.4,0.7-0.7,0.7L49.3,119.5z" fill="currentColor" style="transform-origin: 100px 1px;"
              class="octo-arm"/>
        <path d="M153.1,99.3c-6.3,1.6-10.6,2.8-16.8,4.4c-1.5,0.4-1.6,0.5-2.9-1c-1.5-1.7-2.6-2.8-4.7-3.8c-6.3-3.1-12.4-2.2-18.1,1.5
						c-6.8,4.4-10.3,10.9-10.2,19c0.1,8,5.6,14.6,13.5,15.7c6.8,0.9,12.5-1.5,17-6.6c0.9-1.1,1.7-2.3,2.7-3.7c-3.6,0-8.1,0-19.3,0
						c-2.1,0-2.6-1.3-1.9-3c1.3-3.1,3.7-8.3,5.1-10.9c0.3-0.6,1-1.6,2.5-1.6c5.1,0,23.9,0,36.4,0c-0.2,2.7-0.2,5.4-0.6,8.1
						c-1.1,7.2-3.8,13.8-8.2,19.6c-7.2,9.5-16.6,15.4-28.5,17c-9.8,1.3-18.9-0.6-26.9-6.6c-7.4-5.6-11.6-13-12.7-22.2
						c-1.3-10.9,1.9-20.7,8.5-29.3c7.1-9.3,16.5-15.2,28-17.3c9.4-1.7,18.4-0.6,26.5,4.9c5.3,3.5,9.1,8.3,11.6,14.1
						C154.7,98.5,154.3,99,153.1,99.3z" fill="currentColor" style="transform-origin: 20px 120px;"
              class="octo-arm"/>
        <path d="M186.2,154.6c-9.1-0.2-17.4-2.8-24.4-8.8c-5.9-5.1-9.6-11.6-10.8-19.3c-1.8-11.3,1.3-21.3,8.1-30.2
						c7.3-9.6,16.1-14.6,28-16.7c10.2-1.8,19.8-0.8,28.5,5.1c7.9,5.4,12.8,12.7,14.1,22.3c1.7,13.5-2.2,24.5-11.5,33.9
						c-6.6,6.7-14.7,10.9-24,12.8C191.5,154.2,188.8,154.3,186.2,154.6z M210,114.2c-0.1-1.3-0.1-2.3-0.3-3.3
						c-1.8-9.9-10.9-15.5-20.4-13.3c-9.3,2.1-15.3,8-17.5,17.4c-1.8,7.8,2,15.7,9.2,18.9c5.5,2.4,11,2.1,16.3-0.6
						C205.2,129.2,209.5,122.8,210,114.2z" fill="currentColor" style="transform-origin: 120px 140px;"
              class="octo-arm"/>
    </svg>
</a>

<a href="{{.Repo}}" class="github-corner" aria-label="View source on GitHub">
    <svg width="80" height="80" viewBox="0 0 250 250"
         style="fill:#151513; color:#fff; position: absolute; top: 0; border: 0; right: 0;" aria-hidden="true">
        <path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z"></path>
        <path d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2"
              fill="currentColor" style="transform-origin: 130px 106px;" class="octo-arm"></path>
        <path d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1,77.6 C213.8,80.2 216.3,84.9 216.3,84.9 C212.7,93.1 206.9,96.0 205.4,96.6 C205.1,102.4 203.0,107.8 198.3,112.5 C181.9,128.9 168.3,122.5 157.7,114.1 C157.9,116.9 156.7,120.9 152.7,124.9 L141.0,136.5 C139.8,137.7 141.6,141.9 141.8,141.8 Z"
              fill="currentColor" class="octo-body"></path>
    </svg>
</a>


<article class="markdown-body"><div id="content"></div></article>
<script>
var xhr = new XMLHttpRequest();
xhr.open("get", "{{.ReadMe}}", true);
xhr.onreadystatechange = function(){
	if(xhr.status === 200){
		document.getElementById('content').innerHTML = marked.parse(xhr.responseText);
	}
}
xhr.send();
</script>
</body>
</html>
`

type packageData struct {
	Main    string
	PkgName string
	Repo    string
	Branch  string
	ReadMe  string
}

var tpl = template.Must(template.New("pkg").Parse(inputHTML))

func (pd packageData) Write(w io.Writer) error {
	return tpl.Execute(w, pd)
}

func (pd packageData) WriteToFile() error {
	filename := path.Join(pd.PkgName + ".html")
	dir, _ := filepath.Split(filename)
	_ = os.MkdirAll(dir, os.ModePerm)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = pd.Write(file)
	_ = file.Close()

	return err
}
