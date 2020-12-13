package cmd

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate model",
	RunE:  gen,
}

var projectPath = "github.com/inari111/moconv"

var domainModelPath = "github.com/inari111/moconv/domain/user"

var outputPath = "./examples/infra/persistence/postgres"

type Field struct {
	Name string
	Type ast.Expr
	Tag  map[string]string
}

func NewField(
	name string,
	t ast.Expr,
	tag map[string]string,
) *Field {
	return &Field{
		Name: name,
		Type: t,
		Tag:  tag,
	}
}

type Params struct {
	StructName string
	Fields     []*Field
}

func gen(cmd *cobra.Command, args []string) error {
	generate()

	/**
	 * ここから下のコードはデバッグ用
	 */
	filename := "./examples/domain/user/user.go"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return err
	}

	//for _, d := range f.Decls {
	//	ast.Print(fset, d)
	//	fmt.Println()
	//}

	conf := types.Config{
		Importer: importer.Default(),
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  nil,
		Selections: nil,
		Scopes:     nil,
		InitOrder:  nil,
	}

	pkg, err := conf.Check("./examples/domain/user", fset, []*ast.File{f}, info)
	if err != nil {
		return err
	}

	exprMap := make(map[ast.Expr]types.Type)
	for expr, obj := range info.Types {
		fmt.Printf("expr: %+v\n", expr)
		fmt.Printf("obj: %+v\n", obj)
		fmt.Println()

		exprMap = map[ast.Expr]types.Type{
			expr: obj.Type,
		}
	}
	fmt.Printf("exprMap: %+v\n", exprMap)
	fmt.Println("----------------------------------")

	for id, obj := range info.Defs {
		fmt.Printf("id.Name: %+v\n", id.Name)
		fmt.Printf("obj: %+v\n", obj)
		//fmt.Printf("%+v\n", obj.Type())
		//fmt.Printf("%+v\n", obj.String())
		//fmt.Printf("%+v\n", obj.Name())
		//fmt.Printf("%+v\n", obj.Id())
		fmt.Println()
	}
	fmt.Println("----------------------------------")

	for id, obj := range info.Uses {
		fmt.Printf("id.Name: %+v\n", id.Name)
		fmt.Printf("obj: %+v\n", obj.Type().String())
		fmt.Println()
	}

	scope := pkg.Scope()
	for _, name := range scope.Names() {
		fmt.Printf("pkg.Scope name: %+v\n", name)
		obj := scope.Lookup(name)
		fmt.Printf("pkg.Scope obj: %+v\n", obj)
		fmt.Printf("pkg.Scope Type: %+v\n", obj.Type())
		fmt.Printf("pkg.Scope Id: %+v\n", obj.Id())
		fmt.Printf("pkg.Scope Name: %+v\n", obj.Name())
		fmt.Printf("pkg.Scope String: %+v\n", obj.String())
		fmt.Printf("pkg.Scope Parent: %+v\n", obj.Parent())
		fmt.Println()
	}

	structName := ""
	fields := []*Field{}
	ast.Inspect(f, func(n ast.Node) bool {
		switch nt := n.(type) {
		case *ast.Ident:
			obj := info.ObjectOf(nt)
			if obj == nil {
				fmt.Println("obj == nil だった")
				fmt.Println()
				return true
			}
			if _, ok := obj.(*types.TypeName); !ok {
				fmt.Println("not *types.TypeName")
				fmt.Println()
				return true
			}
			typ := obj.Type()
			fmt.Printf("obj.Name: %+v\n", obj.Name())
			fmt.Printf("obj.Id: %+v\n", obj.Id())
			fmt.Printf("obj.Type: %+v\n", typ)
			fmt.Println()
			//t = obj.Type().String()

		default:

			x, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if v, ok := x.Type.(*ast.StructType); ok {
				if structName == "" {
					structName = x.Name.String()
				}
				for _, field := range v.Fields.List {
					fmt.Printf("name: %+v\n", field.Names)
					fmt.Printf("type: %+v\n", field.Type)
					name := field.Names[0]
					t := exprMap[field.Type]
					fmt.Printf("type: %+v\n", t)

					fields = append(fields, NewField(name.Name, field.Type, nil))
				}
			}
		}

		return true
	})
	fmt.Printf("structName: %+v\n", structName)

	return nil
}

// コード生成部分
// 本来はastから必要な情報を抽出して生成したかった
func generate() error {
	f := jen.NewFilePath("./examples/infra/persistence/postgres")
	f.Type().Id("User").Struct(
		jen.Id("ID").String().Tag(map[string]string{"db": "id"}),
		jen.Id("Email").String().Tag(map[string]string{"db": "email"}),
		jen.Id("Password").String().Tag(map[string]string{"db": "password"}),
		jen.Id("CreatedAt").Qual("time", "Time").Tag(map[string]string{"db": "created_at"}),
		jen.Id("UpdatedAt").Qual("time", "Time").Tag(map[string]string{"db": "updated_at"}),
	)

	f.Func().Params(
		jen.Id("u").Op("*").Id("User"),
	).Id("ToDomain").Params().Op("*").Qual("github.com/inari111/moconv/examples/domain/user", "User").Block(
		jen.Return(jen.Op("&").Qual("github.com/inari111/moconv/examples/domain/user", "User").Block(
			jen.Id("ID").Op(":").Qual("github.com/inari111/moconv/examples/domain", "UserID").Call(jen.Qual("u", "ID")).Op(","),
			jen.Id("Email").Op(":").Qual("u", "Email").Op(","),
			jen.Id("Password").Op(":").Qual("u", "Password").Op(","),
			jen.Id("CreatedAt").Op(":").Qual("u", "CreatedAt").Op(","),
			jen.Id("UpdatedAt").Op(":").Qual("u", "UpdatedAt").Op(","),
		)),
	)

	f.Func().Id("NewUser").Params(
		jen.Id("u").Op("*").Id("user.User"),
	).Op("*").Op("User").Block(
		jen.Return(jen.Op("&").Op("User").Block(
			jen.Id("ID").Op(":").Qual("u", "ID").Dot("String").Call().Op(","),
			jen.Id("Email").Op(":").Qual("u", "Email").Op(","),
			jen.Id("Password").Op(":").Qual("u", "Password").Op(","),
			jen.Id("CreatedAt").Op(":").Qual("u", "CreatedAt").Op(","),
			jen.Id("UpdatedAt").Op(":").Qual("u", "UpdatedAt").Op(","),
		)),
	)

	fmt.Printf("%#v", f)
	f.Save("./examples/infra/persistence/postgres/user.go")
	return nil
}

//func generate2(p Params) error {
//	f := jen.NewFilePath(outputPath)
//	fields := []jen.Code{}
//	for _, v := range p.Fields {
//		var f *jen.Statement
//		switch v.Type {
//		//case ast.Expr():
//		//	f = jen.Id(v.Name).String().Tag(v.Tag)
//		}
//
//		fields = append(fields, f)
//	}
//	f.Type().Id(p.StructName).Struct(fields...)
//
//	fmt.Printf("%#v", f)
//	f.Save("./examples/infra/persistence/postgres/user.go")
//	return nil
//}
