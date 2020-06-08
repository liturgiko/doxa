package template

import (
	"encoding/json"
	"fmt"
	"github.com/liturgiko/doxa/pkg/enums/idTypes"
	"github.com/liturgiko/doxa/pkg/enums/statuses"
	"github.com/liturgiko/doxa/pkg/enums/templateTypes"
	"testing"
	"time"
)

func TestParagraph(t *testing.T) {
	var p Paragraph
	var s1 Span
	s1.AddChildSpan(*NewNid("Literal Text"))
	s1.AddChildSpan(*NewSid("actors/Priest"))
	s1.AddChildSpan(*NewRid("oc.*/ocVE.ApolTheotokionVM.text",0,0))
	p.AddSpan(s1)
	p.AddVersion()
	j, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(j))

}
func TestALT(t *testing.T) {
	alt := new(ATEM)
	alt.ID = "x/y"
	alt.Type = templateTypes.Service
	alt.Status = statuses.Draft
	alt.HtmlCss = "ages.html.css"
	pdf := new(PDF)
	pdf.Title = "I am the title"
	pdf.CSS = "ages.pdf.css"
	pdf.PageNbr = 1
	alt.PDF = pdf
	headerEven := NewHeaderEven()
	headerEven.AddLeftDirective(NewLiteralDirective(".it","This is a test"))
	headerEven.AddCenterDirective(NewDateDirective(".it", time.Now()))
	pdf.AddHeader(*headerEven)
	headerOdd := NewHeaderOdd()
	lookupDirective := NewLookupDirective(1)
	err := lookupDirective.AddLookupTK(idTypes.SID, "actor", "actors/priest")
	headerOdd.AddRightDirective(lookupDirective)
	headerOdd.AddCenterDirective(NewPageNbrDirective(".it"))
	pdf.AddHeader(*headerOdd)
	footerEven := NewFooterEven()
	footerEven.AddCenterDirective(NewDateDirective(".it", time.Now()))
	pdf.AddFooter(*footerEven)
	footerOdd := NewFooterOdd()
	footerOdd.AddCenterDirective(NewPageNbrDirective(".it"))
	pdf.AddFooter(*footerOdd)
	j, err := json.MarshalIndent(alt, "", " ")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(j))
}