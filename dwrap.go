package dwrap

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
)

const doc = "dwrap forces a public function to begin by deferring a call to a specific function"

const name = "dwrap"

func NewAnalyzer(pkgPath string, funcName string, ignorePkgs ...string) *analysis.Analyzer {
	ignored := make(map[string]struct{})
	for _, pkg := range ignorePkgs {
		ignored[pkg] = struct{}{}
	}
	r := runner{
		pkgPath:    pkgPath,
		funcName:   funcName,
		ignorePkgs: ignored,
	}
	return &analysis.Analyzer{
		Name: name,
		Doc:  doc,
		Run:  r.run,
		Requires: []*analysis.Analyzer{
			commentmap.Analyzer,
		},
	}
}

type runner struct {
	pkgPath    string
	funcName   string
	ignorePkgs map[string]struct{}
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
	if _, ok := r.ignorePkgs[pass.Pkg.Path()]; ok {
		return nil, nil
	}
	cmaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if cmaps.IgnorePos(decl.Pos(), name) {
				continue
			}
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if !fn.Name.IsExported() {
				continue
			}
			if fn.Body == nil {
				continue
			}
			if fn.Body.List == nil {
				continue
			}
			if !returnsError(pass, fn) {
				continue
			}
			first := fn.Body.List[0]
			switch first := first.(type) {
			case *ast.DeferStmt:
				switch f := first.Call.Fun.(type) {
				case *ast.Ident:
					got := pass.TypesInfo.ObjectOf(f)
					want := analysisutil.ObjectOf(pass, r.pkgPath, r.funcName)
					if got != want {
						pass.Reportf(decl.Pos(), "should call defer %s.%s on the first line", r.pkgName(), r.funcName)
						continue
					}
					continue
				case *ast.SelectorExpr:
					got := pass.TypesInfo.ObjectOf(f.Sel)
					want := analysisutil.ObjectOf(pass, r.pkgPath, r.funcName)
					if got != want {
						pass.Reportf(decl.Pos(), "should call defer %s.%s on the first line", r.pkgName(), r.funcName)
						continue
					}
					continue
				}
			default:
				pass.Reportf(decl.Pos(), "should call defer %s.%s on the first line", r.pkgName(), r.funcName)
			}
		}
	}
	return nil, nil
}

func (r runner) pkgName() string {
	pp := strings.Split(r.pkgPath, "/")
	return pp[len(pp)-1]
}

func returnsError(pass *analysis.Pass, fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil {
		return false
	}
	for _, arg := range fn.Type.Results.List {
		if types.Identical(pass.TypesInfo.TypeOf(arg.Type), types.Universe.Lookup("error").Type()) {
			return true
		}
	}
	return false
}