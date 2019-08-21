package enumcase

import (
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "enumcase",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const Doc = "enumcase ..."

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.SwitchStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		sw, ok := n.(*ast.SwitchStmt)
		if !ok {
			return
		}

		// When the `switch` statement has `default` clause, finishes processing immediately.
		// Because `default` clause is exhaustive.
		for _, stmt := range sw.Body.List {
			cc, ok := stmt.(*ast.CaseClause)
			if !ok {
				continue
			}

			// When `cc.List` is `nil`, it means the `cc` represents `default` clause.
			//
			// https://golang.org/pkg/go/ast/#CaseClause
			if cc.List == nil {
				return
			}
		}

		t := pass.TypesInfo.TypeOf(sw.Tag)
		consts := GetRelatedConsts(pass, t)
		if len(consts) == 0 {
			return
		}

		covers := make(map[*types.Const]bool)
		for _, c := range consts {
			covers[c] = false
		}

		for _, stmt := range sw.Body.List {
			cc, ok := stmt.(*ast.CaseClause)
			if !ok {
				continue
			}

			for _, e := range cc.List {
				var id *ast.Ident
				switch e1 := e.(type) {
				case *ast.SelectorExpr:
					id = e1.Sel
				case *ast.Ident:
					id = e1
				}
				if id == nil {
					continue
				}
				o := pass.TypesInfo.ObjectOf(id)
				c, ok := o.(*types.Const)
				if !ok {
					continue
				}
				if _, ok := covers[c]; ok {
					covers[c] = true
				}
			}
		}

		var uncovers []string
		for c, ok := range covers {
			if !ok {
				uncovers = append(uncovers, c.Name())
			}
		}
		if len(uncovers) == 0 {
			return
		}

		sort.Strings(uncovers)
		typeName := types.TypeString(t, types.RelativeTo(pass.Pkg))
		pass.Reportf(sw.Pos(), "missing case(s) to %s value(s): %s", typeName, strings.Join(uncovers, ", "))
	})

	return nil, nil
}

func GetRelatedConsts(pass *analysis.Pass, t types.Type) []*types.Const {
	named, ok := t.(*types.Named)
	if !ok {
		return nil
	}

	u := named.Underlying()
	if u == nil {
		return nil
	}
	if _, ok = u.(*types.Basic); !ok {
		return nil
	}

	pkg := named.Obj().Pkg()
	scope := pkg.Scope()

	var consts []*types.Const
	for _, name := range scope.Names() {
		o := scope.Lookup(name)
		c, ok := o.(*types.Const)
		if !ok {
			continue
		}
		if c.Type() != t || !c.Exported() {
			continue
		}

		consts = append(consts, c)
	}

	return consts
}
