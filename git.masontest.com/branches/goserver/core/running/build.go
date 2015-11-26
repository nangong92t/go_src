package running

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"git.masontest.com/branches/goserver/workers/revel"
    "git.masontest.com/branches/goserver/plugins"
)

var importErrorPattern = regexp.MustCompile("cannot find package \"([^\"]+)\"")
var runtimePath = os.Getenv("GOPATH") + "/src/git.masontest.com/branches/goserver/runtime"

// Build the app:
// 1. Generate the the main.go file.
// 2. Run the appropriate "go build" command.
// Requires that revel.Init has been called previously.
// Returns the path to the built binary, and an error if there was a problem building it.
func Build(appEnv, appName, appPath string, logger *plugins.ServerLog) {
	revel.Init(appEnv, appPath, "")

	// First, clear the generated files (to avoid them messing with ProcessSource).
	cleanSource(appName)

	sourceInfo, compileError := ProcessSource(revel.CodePaths)
	if compileError != nil {
        panic(compileError)
	}

    checkIsHaveArgs := func(ctrl []*TypeInfo) bool {
        result      := false

	    L: for _, c := range ctrl {
		    for _, a:= range c.MethodSpecs {
                if len(a.Args) > 0 {
                    result  = true
                    break L
                }
			}
	    }

        return result
    }

    ctrlSpecs       := sourceInfo.ControllerSpecs()

	// Generate two source files.
	templateArgs := map[string]interface{}{
        "AppName":        appName,
        "AppPath":        appPath,
        "AppEnv":         appEnv,
		"Controllers":    ctrlSpecs,
		"ValidationKeys": sourceInfo.ValidationKeys,
		"ImportPaths":    calcImportAliases(sourceInfo),
		"TestSuites":     sourceInfo.TestSuites(),
        "IsArgInCtrl":    checkIsHaveArgs(ctrlSpecs),
	}

	genSource(runtimePath, appName + ".go", MAIN, templateArgs)
	// genSource("routes", "routes.go", ROUTES, templateArgs)

}

// Build the app link
func BuildAppLink(appNames []string) {
	genSource(runtimePath, "Processer.go", APPLINK, map[string]interface{}{"Apps": appNames})
}

// Try to define a version string for the compiled app
// The following is tried (first match returns):
// - Read a version explicitly specified in the APP_VERSION environment
//   variable
// - Read the output of "git describe" if the source is in a git repository
// If no version can be determined, an empty string is returned.
func getAppVersion() string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}

	// Check for the git binary
	if gitPath, err := exec.LookPath("git"); err == nil {
		// Check for the .git directory
		gitDir := path.Join(revel.BasePath, ".git")
		info, err := os.Stat(gitDir)
		if (err != nil && os.IsNotExist(err)) || !info.IsDir() {
			return ""
		}
		gitCmd := exec.Command(gitPath, "--git-dir="+gitDir, "describe", "--always", "--dirty")
		revel.TRACE.Println("Exec:", gitCmd.Args)
		output, err := gitCmd.Output()

		if err != nil {
			revel.WARN.Println("Cannot determine git repository version:", err)
			return ""
		}

		return "git-" + strings.TrimSpace(string(output))
	}

	return ""
}

func cleanSource(appName string) {
    os.Remove(runtimePath + "/" + appName + ".go")
}

// genSource renders the given template to produce source code, which it writes
// to the given directory and file.
func genSource(dir, filename, templateSource string, args map[string]interface{}) {
	sourceCode := revel.ExecuteTemplate(
		template.Must(template.New("").Parse(templateSource)),
		args)

	// Create a fresh dir.
	// tmpPath := path.Join(revel.AppPath, dir)

	// Create the file
	file, err := os.Create(path.Join(dir, filename))
	defer file.Close()
	if err != nil {
		revel.ERROR.Fatalf("Failed to create file: %v", err)
	}
	_, err = file.WriteString(sourceCode)
	if err != nil {
		revel.ERROR.Fatalf("Failed to write to file: %v", err)
	}
}

// Looks through all the method args and returns a set of unique import paths
// that cover all the method arg types.
// Additionally, assign package aliases when necessary to resolve ambiguity.
func calcImportAliases(src *SourceInfo) map[string]string {
	aliases := make(map[string]string)
	typeArrays := [][]*TypeInfo{src.ControllerSpecs(), src.TestSuites()}
	for _, specs := range typeArrays {
		for _, spec := range specs {
			addAlias(aliases, spec.ImportPath, spec.PackageName)

			for _, methSpec := range spec.MethodSpecs {
				for _, methArg := range methSpec.Args {
					if methArg.ImportPath == "" {
						continue
					}

					addAlias(aliases, methArg.ImportPath, methArg.TypeExpr.PkgName)
				}
			}
		}
	}

	// Add the "InitImportPaths", with alias "_"
	for _, importPath := range src.InitImportPaths {
		if _, ok := aliases[importPath]; !ok {
			aliases[importPath] = "_"
		}
	}

	return aliases
}

func addAlias(aliases map[string]string, importPath, pkgName string) {
	alias, ok := aliases[importPath]
	if ok {
		return
	}
	alias = makePackageAlias(aliases, pkgName)
	aliases[importPath] = alias
}

func makePackageAlias(aliases map[string]string, pkgName string) string {
	i := 0
	alias := pkgName
	for containsValue(aliases, alias) {
		alias = fmt.Sprintf("%s%d", pkgName, i)
		i++
	}
	return alias
}

func containsValue(m map[string]string, val string) bool {
	for _, v := range m {
		if v == val {
			return true
		}
	}
	return false
}

// Parse the output of the "go build" command.
// Return a detailed Error.
func newCompileError(output []byte) *revel.Error {
	errorMatch := regexp.MustCompile(`(?m)^([^:#]+):(\d+):(\d+:)? (.*)$`).
		FindSubmatch(output)
	if errorMatch == nil {
		revel.ERROR.Println("Failed to parse build errors:\n", string(output))
		return &revel.Error{
			SourceType:  "Go code",
			Title:       "Go Compilation Error",
			Description: "See console for build error.",
		}
	}

	// Read the source for the offending file.
	var (
		relFilename    = string(errorMatch[1]) // e.g. "src/revel/sample/app/controllers/app.go"
		absFilename, _ = filepath.Abs(relFilename)
		line, _        = strconv.Atoi(string(errorMatch[2]))
		description    = string(errorMatch[4])
		compileError   = &revel.Error{
			SourceType:  "Go code",
			Title:       "Go Compilation Error",
			Path:        relFilename,
			Description: description,
			Line:        line,
		}
	)

	fileStr, err := revel.ReadLines(absFilename)
	if err != nil {
		compileError.MetaError = absFilename + ": " + err.Error()
		revel.ERROR.Println(compileError.MetaError)
		return compileError
	}

	compileError.SourceLines = fileStr
	return compileError
}

const MAIN = `// GENERATED CODE - DO NOT EDIT
package runtime

import (
    {{if eq .IsArgInCtrl true}}"reflect"{{end}}
    "errors"
    "runtime/debug"
    "fmt"
    "net/http"
    "net/url" 

    "git.masontest.com/branches/goserver/plugins"
	"git.masontest.com/branches/goserver/workers/revel"{{range $k, $v := $.ImportPaths}}
	{{$v}} "{{$k}}"{{end}}
)

func init() {
    logger = plugins.NewServerLog("statics")
    Register{{.AppName}}Controllers()

    showRequest, _  := http.NewRequest("GET", "", nil)
    revelRequest    = revel.NewRequest(showRequest)
    revelParams     = &revel.Params{Values: make(url.Values)}
}

func Register{{.AppName}}Controllers() {
    revel.Init("{{.AppEnv}}", "{{.AppPath}}", "")

    revel.RunStartupHooks()

	{{range $i, $c := .Controllers}}
	revel.RegisterController((*{{index $.ImportPaths .ImportPath}}.{{.StructName}})(nil),
		[]*revel.MethodType{
			{{range .MethodSpecs}}&revel.MethodType{
				Name: "{{.Name}}",
				Args: []*revel.MethodArg{ {{range .Args}}
					&revel.MethodArg{Name: "{{.Name}}", Type: reflect.TypeOf((*{{index $.ImportPaths .ImportPath | .TypeExpr.TypeName}})(nil)) },{{end}}
				},
				RenderArgNames: map[int][]string{ {{range .RenderCalls}}
					{{.Line}}: []string{ {{range .Names}}
						"{{.}}",{{end}}
					},{{end}}
				},
			},
			{{end}}
		})
	{{end}}
	revel.DefaultValidationKeys = map[string]map[int]string{ {{range $path, $lines := .ValidationKeys}}
		"{{$path}}": { {{range $line, $key := $lines}}
			{{$line}}: "{{$key}}",{{end}}
		},{{end}}
	}
	revel.TestSuites = []interface{}{ {{range .TestSuites}}
		(*{{index $.ImportPaths .ImportPath}}.{{.StructName}})(nil),{{end}}
	}
}

func Process{{.AppName}}(className, method string, params []interface{}) (result interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            logger.Add("Error in worker, %s: %s", r, debug.Stack())
            switch x := r.(type) {
                case string:
                    err = errors.New(x)
                case error:
                    err = x
                default:
                    err = errors.New("Unknown error")
            }

            result  = nil
        }
    }()

    c   := revel.NewController(revelRequest, nil)
    c.Params = revelParams
    if c.Result != nil { c.Result.Clean() }

    if err = c.SetAction(className, method); err != nil {
        return nil, errors.New(fmt.Sprintf("Failed to set action: %s", err))
    }

    // mapping the params
    for i, value := range params {
        if i < len(c.MethodType.Args) {
            arg := c.MethodType.Args[i]
            if values, ok := value.([]interface{}); ok {
                for _, one := range values {
                    c.Params.Add(arg.Name+"[]", fmt.Sprintf("%v", one))
                }
            } else {
                c.Params.Set(arg.Name, fmt.Sprintf("%v", value))
            }
        } else {
            return nil, errors.New("Too many parameters to " + className + " trying to add " + method) 
        }
    }

    revel.Filters[0](c, revel.Filters[1:])

    res     := c.Result.(revel.RenderRpcResult)

    // remove old params
    for i, value := range params {
        arg := c.MethodType.Args[i]
        if _, ok := value.([]interface{}); ok {
            c.Params.Del(arg.Name+"[]")
        } else {
            c.Params.Del(arg.Name)
        }
    }
    return res.Data, res.Error
}
`
const ROUTES = `// GENERATED CODE - DO NOT EDIT
package routes

import "git.masontest.com/branches/goserver/workers/revel"

{{range $i, $c := .Controllers}}
type t{{.StructName}} struct {}
var {{.StructName}} t{{.StructName}}

{{range .MethodSpecs}}
func (_ t{{$c.StructName}}) {{.Name}}({{range .Args}}
		{{.Name}} {{if .ImportPath}}interface{}{{else}}{{.TypeExpr.TypeName ""}}{{end}},{{end}}
		) string {
	args := make(map[string]string)
	{{range .Args}}
	revel.Unbind(args, "{{.Name}}", {{.Name}}){{end}}
	return revel.MainRouter.Reverse("{{$c.StructName}}.{{.Name}}", args).Url
}
{{end}}
{{end}}

`

const APPLINK   = `// GENERATED CODE - DO NOT EDIT
package runtime

import (
    "errors"

	"git.masontest.com/branches/goserver/protocols"
    "git.masontest.com/branches/goserver/plugins"
	"git.masontest.com/branches/goserver/workers/revel"
)

var (
    logger  *plugins.ServerLog
    revelRequest *revel.Request
    revelParams *revel.Params
)

func Processer(appName string, req *protocols.RequestData) (interface{}, error) {
    switch appName {
    {{range $i, $n := .Apps}}
        case "{{$n}}":
	        return Process{{$n}}(req.Class, req.Method, req.Params){{end}}
    }

    return nil, errors.New("Now Found out this worker app")
}
`
