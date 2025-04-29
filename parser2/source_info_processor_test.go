package parser2

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/wiselike/revel-cmd/logger"
	"github.com/wiselike/revel-cmd/model"
	"golang.org/x/tools/go/packages"
)

// helper builds a minimal SourceInfoProcessor together with
// an *ast.FuncDecl and its *packages.Package wrapper that points
// to the supplied Go source snippet.
func buildTestObjects(src string, funcName string) (*SourceInfoProcessor, *ast.FuncDecl, *packages.Package, error) {
	// 1. Parse source into AST.
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		return nil, nil, nil, err
	}
	// 2. Locate the desired function declaration.
	var funcDecl *ast.FuncDecl
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == funcName {
			funcDecl = fd
			break
		}
	}
	if funcDecl == nil {
		return nil, nil, nil, errFunctionNotFound(funcName)
	}

	// 3. Create a minimal *packages.Package; only the fields used by
	//    getValidation are required (Fset, Syntax, PkgPath, Name).
	pkg := &packages.Package{
		Fset:    fset,
		Syntax:  []*ast.File{file},
		PkgPath: "github.com/example/app/controllers",
		Name:    "controllers",
	}

	// 4. Stub out SourceProcessor with just what getValidation touches
	//    (log + importMap).
	sp := &SourceProcessor{
		log:       logger.New("unit", "test"),
		importMap: map[string]string{"revel": model.RevelImportPath},
	}

	// 5. Real object under test.
	sip := NewSourceInfoProcessor(sp)

	return sip, funcDecl, pkg, nil
}

// --- Tests ------------------------------------------------------------

func TestGetValidation_ReceiverSelector(t *testing.T) {
	src := `package controllers

import "github.com/wiselike/revel"

type MyController struct {
    *revel.Controller
}

func (c *MyController) Index(name string) {
    c.Validation.Required(name)
}`

	sip, fd, pkg, err := buildTestObjects(src, "Index")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	keys := sip.getValidation(fd, pkg)
	if len(keys) != 1 {
		t.Fatalf("expected 1 validation key, got %d", len(keys))
	}
	for _, k := range keys {
		if k != "name" {
			t.Fatalf("expected key 'name', got %q", k)
		}
	}
}

func TestGetValidation_ReceiverSelector2(t *testing.T) {
	src := `package controllers

import "github.com/wiselike/revel"

type MyController struct {
    *revel.Controller
}

func (c *MyController) HandleUpload(avatar []byte) revel.Result {
	// Validation rules.
	c.Validation.Required(avatar)
	c.Validation.MinSize(avatar, 2*KB).
		Message("Minimum a file size of 2KB expected")
	c.Validation.MaxSize(avatar, 2*MB).
		Message("File cannot be larger than 2MB")

	// Check format of the file.
	conf, format, err := image.DecodeConfig(bytes.NewReader(avatar))
	c.Validation.Required(err == nil).Key("avatar").
		Message("Incorrect file format")
	c.Validation.Required(format == "jpeg" || format == "png").Key("avatar").
		Message("JPEG or PNG file format is expected")

	// Check resolution.
	c.Validation.Required(conf.Height >= 150 && conf.Width >= 150).Key("avatar").
		Message("Minimum allowed resolution is 150x150px")

	// Handle errors.
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect((*Single).Upload)
	}

	return c.RenderJSON(FileInfo{
		ContentType: c.Params.Files["avatar"][0].Header.Get("Content-Type"),
		Filename:    c.Params.Files["avatar"][0].Filename,
		RealFormat:  format,
		Resolution:  fmt.Sprintf("%dx%d", conf.Width, conf.Height),
		Size:        len(avatar),
		Status:      "Successfully uploaded",
	})
}`

	sip, fd, pkg, err := buildTestObjects(src, "HandleUpload")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	keys := sip.getValidation(fd, pkg)
	if len(keys) != 6 {
		t.Fatalf("expected 6 validation key, got %d", len(keys))
	}
}

func TestGetValidation_ParamPointer(t *testing.T) {
	src := `package controllers

import "github.com/wiselike/revel"

func ValidateUser(v *revel.Validation) {
    var email string
    v.Required(email)
}`

	sip, fd, pkg, err := buildTestObjects(src, "ValidateUser")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	keys := sip.getValidation(fd, pkg)
	if len(keys) != 1 {
		t.Fatalf("expected 1 validation key, got %d", len(keys))
	}
	for _, k := range keys {
		if k != "email" {
			t.Fatalf("expected key 'email', got %q", k)
		}
	}
}

// errFunctionNotFound is a small helper error type so callers can get a
// sensible message if the sample snippet doesn't compile as expected.

type errFunctionNotFound string

func (e errFunctionNotFound) Error() string { return "function not found: " + string(e) }
