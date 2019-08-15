package generator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var rootapi []APIDoc
var astPkgs []*ast.Package
var basicTypes = map[string]string{
	"bool":    "Boolean",
	"uint":    "Number",
	"uint8":   "Number",
	"uint16":  "Number",
	"uint32":  "Number",
	"uint64":  "Number",
	"int":     "Number",
	"int8":    "Number",
	"int16":   "Number",
	"int32":   "Number",
	"int64":   "Number",
	"uintptr": "Number",
	"float32": "Number",
	"float64": "Number",
	"string":  "String",
}

const (
	astTypeArray  = "array"
	astTypeObject = "object"
	astTypeMap    = "map"
)

func init() {
	astPkgs = make([]*ast.Package, 0)
}

// Scan scan file
func Scan(currpath string, filename string) (result *[]APIDoc, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join(currpath, filename), nil, parser.ParseComments)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if f.Comments != nil {
		for _, d := range f.Decls {
			var tempAPI = APIDoc{
				Params:  map[string]Param{},
				Success: map[string]Param{},
			}
			switch specDecl := d.(type) {
			case *ast.FuncDecl:
				for _, s := range strings.Split(specDecl.Doc.Text(), "\n") {
					if strings.HasPrefix(s, "@route") {
						elements := strings.TrimSpace(s[len("@router"):])
						r := strings.SplitN(elements, " ", 2)
						if len(r) != 2 {
							return nil, errors.New("wrong route")
						}
						tempAPI.Path = r[0]
						tempAPI.Method = r[1]
					} else if strings.HasPrefix(s, "@name") {
						tempAPI.APIName = strings.TrimSpace(s[len("@name"):])
						tempAPI.Title = strings.TrimSpace(s[len("@name"):])
					} else if strings.HasPrefix(s, "@version") {
						tempAPI.Version = strings.TrimSpace(s[len("@version"):])
					} else if strings.HasPrefix(s, "@group") {
						tempAPI.Group = strings.TrimSpace(s[len("@group"):])
					} else if strings.HasPrefix(s, "@group") {
						tempAPI.Group = strings.TrimSpace(s[len("@group"):])
					} else if strings.HasPrefix(s, "@title") {
						tempAPI.Title = strings.TrimSpace(s[len("@title"):])
					} else if strings.HasPrefix(s, "@in") {
						err := inputHandle(&tempAPI, s)
						if err != nil {
							return nil, err
						}
					} else if strings.HasPrefix(s, "@out") {
						err := outputHandle(&tempAPI, s)
						if err != nil {
							return nil, err
						}
					}
				}
				rootapi = append(rootapi, tempAPI)
			case *ast.GenDecl:
				continue
			}
		}
	}
	return &rootapi, nil
}

func inputHandle(tempAPI *APIDoc, s string) error {
	result, err := checker("@in", s)
	if err != nil {
		return err
	}
	if strings.HasPrefix(result[0], "[]") {
		if isBasicType(result[0][2:]) {
			tempAPI.Params[result[1]] = Param{
				Type:        basicTypes[result[0][2:]] + "[]",
				Description: result[2],
			}
			return nil
		}
	}
	if isBasicType(result[0]) {
		tempAPI.Params[result[1]] = Param{
			Type:        basicTypes[result[0]],
			Description: result[2], // 检查是否越界
		}
		return nil
	}
	strs := strings.Split(result[1], ".")
	packageName := strs[0]
	objectname := strs[len(strs)-1]
	if strings.HasPrefix(packageName, "[]") {
		log.Fatalln("doesn'h support input array")
	}
	params := handleObject(objectname, packageName)
	if params == nil {
		log.Fatalln("you should run this under your project root path")
	}
	for _, v := range *params {
		tempAPI.Params[v.Name] = v
	}
	return nil
}

func outputHandle(tempAPI *APIDoc, s string) error {
	result, err := checker("@out", s)
	if err != nil {
		return err
	}
	if isBasicType(result[0]) {
		tempAPI.Success[result[1]] = Param{
			Type:        basicTypes[result[0]],
			Description: result[2], // 检查是否越界
		}
		return nil
	}
	strs := strings.Split(result[1], ".")
	packageName := strs[0]
	objectname := strs[len(strs)-1]
	if strings.HasPrefix(packageName, "[]") {
		tempAPI.IsSuccessArray = true
		packageName = strings.TrimPrefix(packageName, "[]")
	}
	params := handleObject(objectname, packageName)
	for _, v := range *params {
		tempAPI.Success[v.Name] = v
	}
	return nil
}

func handleObject(objectname string, packageName string) *[]Param {
	var params *[]Param
L:
	for _, pkg := range astPkgs {
		if packageName == pkg.Name {
			for _, fl := range pkg.Files {
				for k, d := range fl.Scope.Objects {
					if d.Kind == ast.Typ {
						if k != objectname {
							// Still searching for the right object
							continue
						}
						params = parseObject(d, packageName)
						// When we've found the correct object, we can stop searching
						break L
					}
				}
			}
		}
	}
	return params
}

func checker(tag string, s string) (result []string, err error) {
	elements := strings.TrimSpace(s[len(tag):])
	result = strings.SplitN(elements, " ", 3)
	if len(result) < 2 {
		return nil, errors.New("wrong " + tag)
	}
	return result, nil
}

func isBasicType(Type string) bool {
	if _, ok := basicTypes[Type]; ok {
		return true
	}
	return false
}

// ParsePackagesFromDir parsePackagesFromDir
func ParsePackagesFromDir(dirpath string) {
	c := make(chan error)

	go func() {
		filepath.Walk(dirpath, func(fpath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !fileInfo.IsDir() {
				return nil
			}

			// skip folder if it's a 'vendor' folder within dirpath or its child,
			// all 'tests' folders and dot folders wihin dirpath
			d, _ := filepath.Rel(dirpath, fpath)
			if !(d == "vendor" || strings.HasPrefix(d, "vendor"+string(os.PathSeparator))) &&
				!strings.Contains(d, "tests") &&
				!(d[0] == '.') {
				err = parsePackageFromDir(fpath)
				if err != nil {
					// Send the error to through the channel and continue walking
					c <- fmt.Errorf("error while parsing directory: %s", err.Error())
					return nil
				}
			}
			return nil
		})
		close(c)
	}()

	for err := range c {
		log.Fatalln(err.Error())
	}
}

func parsePackageFromDir(path string) error {
	fileSet := token.NewFileSet()
	folderPkgs, err := parser.ParseDir(fileSet, path, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
	}, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, v := range folderPkgs {
		astPkgs = append(astPkgs, v)
	}

	return nil
}

func parseObject(d *ast.Object, packageName string) *[]Param {
	ts, ok := d.Decl.(*ast.TypeSpec)
	if !ok {
		log.Fatalf("Unknown type without TypeSec: %v", d)
	}
	switch t := ts.Type.(type) {
	case *ast.ArrayType:
		log.Fatalln("Current does't support arrayType")
	case *ast.Ident:
		log.Fatalln("Current does't support Ident")
	case *ast.StructType:
		return parseStruct(t, packageName)
	}
	return nil
}

func parseStruct(st *ast.StructType, packageName string) *[]Param {
	var result []Param
	if st.Fields.List != nil {
		for _, field := range st.Fields.List { // 默认所有的field都是入参，所以要专门定义dto,vo,po
			var temp Param
			if field.Tag == nil {
				continue
			}
			stag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			switch t := field.Type.(type) {
			case *ast.ArrayType:
				temp.Name = stag.Get("bson")
				child := handleObject(fmt.Sprint(t.Elt), packageName)
				if !isBasicType(fmt.Sprint(t.Elt)) {
					temp.Type = "Object[]"
					temp.Child = *child
				} else {
					temp.Type = basicTypes[fmt.Sprint(t.Elt)] + "[]"
				}
				temp.Description = stag.Get("description")
			case *ast.Ident:
				if t.Obj != nil {
					child := parseObject(t.Obj, packageName)
					temp.Name = stag.Get("bson")
					temp.Type = "Object"
					temp.Child = *child
					temp.Description = stag.Get("description")
				} else {
					temp.Name = stag.Get("bson")
					temp.Type = basicTypes[fmt.Sprint(field.Type)]
					temp.Description = stag.Get("description")
				}
			}
			result = append(result, temp)
		}
	} else {
		log.Fatalln("Struct doesn't have any fileds")
		return nil
	}
	return &result
}
