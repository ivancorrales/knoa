package mutator

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	DefAttributeNameFormat = `("[A-Za-z_]+[A-Za-z0-9_./-]*"|[A-Za-z_]+[A-Za-z0-9_/-]*)`
	arrayIndexExprStr      = `([0-9]+|\*)`
)

type Parser struct {
	RegExp          *regexp.Regexp
	AttributeRegExp *regexp.Regexp
	Strict          bool
}

func RegExpFromAttributeFormat(attributeFormat string) *regexp.Regexp {
	regExpStr := fmt.Sprintf(`^(?P<parent>(((\.)?%s|\[%s\]))*)((\.)(?P<attribute>%s)|(\[(?P<index>%s)\]))$`,
		attributeFormat, arrayIndexExprStr, attributeFormat, arrayIndexExprStr)
	return regexp.MustCompile(regExpStr)
}

func RegExpsFromAttributeFormat(attributeFormat string) (*regexp.Regexp, *regexp.Regexp) {
	regExpStr := fmt.Sprintf(`^(?P<parent>(((\.)?%s|\[%s\]))*)((\.)(?P<attribute>%s)|(\[(?P<index>%s)\]))$`,
		attributeFormat, arrayIndexExprStr, attributeFormat, arrayIndexExprStr)

	return regexp.MustCompile(regExpStr), regexp.MustCompile(fmt.Sprintf(`^(?P<attribute>%s)$`, attributeFormat))
}

func (p *Parser) Parse(pathExpr string) (*Mutator, error) {
	match := p.RegExp.FindStringSubmatch(pathExpr)
	if match == nil {
		attrMatch := p.AttributeRegExp.FindStringSubmatch(pathExpr)
		if attrMatch != nil {
			return &Mutator{
				child: &Mutator{
					name: pathExpr,
				},
			}, nil
		}
		if p.Strict {
			log.Panicf("invalid Path  '%v'. Path doesn't match defined format", pathExpr)
		}
		return nil, fmt.Errorf("invalid path '%s'", pathExpr)
	}
	subMatchMap := map[string]string{}
	for i, name := range p.RegExp.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	attr := subMatchMap["attribute"]
	if strings.HasPrefix(attr, ".") {
		attr = attr[1:]
	}
	parentExpr := subMatchMap["parent"]
	if strings.HasSuffix(parentExpr, ".") {
		parentExpr = parentExpr[:len(parentExpr)-1]
	}
	arrayIndex := subMatchMap["index"]
	m := &Mutator{
		name: attr,
	}
	if arrayIndex != "" {
		m.index = arrayIndex
		parent := &Mutator{}
		var err error
		if parentExpr != "" {
			parent, err = p.Parse(parentExpr)
			if parent == nil {
				parent = &Mutator{
					name: parentExpr,
				}
			}
		}
		addToBottom(parent, m)
		return parent, err
	}
	if parentExpr != "" {
		if attr != "" {
			if attr[0] == '"' && attr[len(attr)-1] == '"' {
				m.name = attr[1 : len(attr)-1]
			}
		}
		var err error
		parent, err := p.Parse(parentExpr)
		if parent == nil {
			parent = &Mutator{
				name: parentExpr,
			}
		}
		addToBottom(parent, m)
		return parent, err
	}
	return m, nil
}

func addToBottom(parent *Mutator, child *Mutator) {
	if parent.child == nil {
		parent.child = child
	} else {
		addToBottom(parent.Child(), child)
	}
}
