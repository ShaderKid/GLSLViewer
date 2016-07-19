# GLSLViewer

## Install dependencies
### Mac
```bash
$ brew install glfw3
$ go get "github.com/go-gl/gl/v4.1-core/gl"
$ go get "github.com/go-gl/glfw/v3.2/glfw"
```

## Run
```bash
$ go run app.go
```

## Build
### bin
```bash
$ go build app.go
```

### .app
```bash
$ go build app.go -o glslviewer
$ mkdir -p GLSLViewer.app/Contents/MacOS GLSLViewer.app/Contents/Resources
$ mv glslviewer GLSLViewer.app/Contents/MacOS/glslviewer
```