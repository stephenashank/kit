package generator

import (
	"fmt"

	"strings"

	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/dave/jennifer/jen"
	"github.com/kujtimiihoxha/gk-cli/fs"
	"github.com/kujtimiihoxha/gk-cli/parser"
	"github.com/kujtimiihoxha/gk-cli/utils"
)

type Gen interface {
	Generate() error
}

type BaseGenerator struct {
	srcFile *jen.File
	code    *PartialGenerator
	fs      *fs.DefaultFs
}

func (b *BaseGenerator) InitPg() {
	b.code = NewPartialGenerator(b.srcFile.Empty())
}
func (b *BaseGenerator) CreateFolderStructure(path string) error {
	e, err := b.fs.Exists(path)
	if err != nil {
		return err
	}
	if !e {
		logrus.Debug(fmt.Sprintf("Creating missing folder structure : %s", path))
		return b.fs.MkdirAll(path)
	}
	return nil
}

// GenerateNameBySample is used to generate a variable name using a sample.
//
// The exclude parameter represents the names that it can not use.
//
// E.x  sample = "hello" this will return the name "h" if it is not in any NamedTypeValue name.
func (b *BaseGenerator) GenerateNameBySample(sample string, exclude []parser.NamedTypeValue) string {
	sn := 1
	name := utils.ToLowerFirstCamelCase(sample)[:sn]
	for _, v := range exclude {
		if v.Name == name {
			sn++
			if sn > len(sample) {
				sample = string(len(sample) - sn)
			}
			name = utils.ToLowerFirstCamelCase(sample)[:sn]
		}
	}
	return name
}
func (b *BaseGenerator) EnsureThatWeUseQualifierIfNeeded(tp string, imp []parser.NamedTypeValue) string {
	if t := strings.Split(tp, "."); len(t) > 0 {
		s := t[0]
		for _, v := range imp {
			i, _ := strconv.Unquote(v.Type)
			if strings.HasSuffix(i, s) || v.Name == s {
				return i
			}
		}
		return ""
	}
	return ""
}

type PartialGenerator struct {
	raw *jen.Statement
}

func NewPartialGenerator(st *jen.Statement) *PartialGenerator {
	if st != nil {
		return &PartialGenerator{
			raw: st,
		}
	}
	return &PartialGenerator{
		raw: &jen.Statement{},
	}
}
func (p *PartialGenerator) appendMultilineComment(c []string) {
	for i, v := range c {
		if i != len(c)-1 {
			p.raw.Comment(v).Line()
			continue
		}
		p.raw.Comment(v)
	}
}
func (p *PartialGenerator) Raw() *jen.Statement {
	return p.raw
}
func (p *PartialGenerator) String() string {
	return p.raw.GoString()
}
func (p *PartialGenerator) appendInterface(name string, methods []jen.Code) {
	p.raw.Type().Id(name).Interface(methods...).Line()
}

func (p *PartialGenerator) appendStruct(name string, fields ...jen.Code) {
	p.raw.Type().Id(name).Struct(fields...).Line()
}
func (p *PartialGenerator) NewLine() {
	p.raw.Line()
}

func (p *PartialGenerator) appendFunction(name string, stp *jen.Statement,
	parameters []jen.Code, results []jen.Code, oneResponse string, body ...jen.Code) {
	p.raw.Func()
	if stp != nil {
		p.raw.Params(stp)
	}
	if name != "" {
		p.raw.Id(name)
	}
	p.raw.Params(parameters...)
	if oneResponse != "" {
		p.raw.Id(oneResponse)
	} else if len(results) > 0 {
		p.raw.Params(results...)
	}
	p.raw.Block(body...)
}
