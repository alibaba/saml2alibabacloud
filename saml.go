package saml2aws

import (
	"fmt"

	"github.com/beevik/etree"
)

const (
	assertionTag          = "Assertion"
	attributeStatementTag = "AttributeStatement"
	attributeTag          = "Attribute"
	attributeValueTag     = "AttributeValue"
)

//ErrMissingElement is the error type that indicates an element and/or attribute is
//missing. It provides a structured error that can be more appropriately acted
//upon.
type ErrMissingElement struct {
	Tag, Attribute string
}

//ErrMissingAssertion indicates that an appropriate assertion element could not
//be found in the SAML Response
var (
	ErrMissingAssertion = ErrMissingElement{Tag: assertionTag}
)

func (e ErrMissingElement) Error() string {
	if e.Attribute != "" {
		return fmt.Sprintf("missing %s attribute on %s element", e.Attribute, e.Tag)
	}
	return fmt.Sprintf("missing %s element", e.Tag)
}

// ExtractAwsRoles given an assertion document extract the aws roles
func ExtractAwsRoles(data []byte) ([]string, error) {

	awsroles := []string{}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(data); err != nil {
		return awsroles, err
	}

	// log.Printf("root tag: %s", doc.Root().Tag)

	assertionElement := doc.FindElement(".//Assertion")
	if assertionElement == nil {
		return nil, ErrMissingAssertion
	}

	// log.Printf("tag: %s", assertionElement.Tag)

	//Get the actual assertion attributes
	attributeStatement := assertionElement.FindElement(childPath(assertionElement.Space, attributeStatementTag))
	if attributeStatement == nil {
		return nil, ErrMissingElement{Tag: attributeStatementTag}
	}

	// log.Printf("tag: %s", attributeStatement.Tag)

	attributes := attributeStatement.FindElements(childPath(assertionElement.Space, attributeTag))
	for _, attribute := range attributes {
		if attribute.SelectAttrValue("Name", "") != "https://aws.amazon.com/SAML/Attributes/Role" {
			continue
		}
		atributeValues := attribute.FindElements(childPath(assertionElement.Space, attributeValueTag))
		for _, attrValue := range atributeValues {
			awsroles = append(awsroles, attrValue.Text())
		}
	}

	return awsroles, nil
}

func childPath(space, tag string) string {
	if space == "" {
		return "./" + tag
	}
	//log.Printf("query = %s", "./"+space+":"+tag)
	return "./" + space + ":" + tag
}
