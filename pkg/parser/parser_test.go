package parser

import (
	"encoding/json"
	"fmt"
	"github.com/liturgiko/doxa/pkg/enums/statuses"
	"github.com/liturgiko/doxa/pkg/enums/templateTypes"
	"github.com/liturgiko/doxa/pkg/template"
	"testing"
)

type MockResolver struct {
	tkMap map[string]bool // used by the mock ExistsTK method
	values map[string]string // used by the mock Values and Versions methods
}
func (r *MockResolver) AddTK(tk string)  {
	r.tkMap[tk] = true
}
func (r *MockResolver) AddValue(id, value string)  {
	r.values[id] = value
}
func (r *MockResolver) Close()  {}
func (r *MockResolver) ExistsTK(topicKey string) bool {
	return r.tkMap[topicKey]
}
func (r *MockResolver) Values(topicKey string, genLibs []template.GenLib) template.TKVal  {
	var tkVal = new(template.TKVal)
	for _, genLib := range genLibs {
		for _, library := range genLib.All() {
			var ltkValue template.LTKVal
			ltkValue.Library = library
			ltkValue.TopicKey = topicKey
			ltkValue.Value = r.values[ltkValue.ID()]
		}
	}
	return *tkVal
}
func (r *MockResolver) Versions(genLibs []template.GenLib) map[string]string  {
	var mock = map[string]string{}
	for _, genLib := range genLibs {
		for _, library := range genLib.All() {
			mock[library] = r.values[library + template.IDPathDelimiter + "properties/version.designation"]
		}
	}
	return mock
}
var genLibs []template.GenLib
var resolver MockResolver
func init() {
	var genLib = new(template.GenLib)
	genLib.Primary = "gr_gr_cog"
	genLibs = append(genLibs, *genLib)
	genLib = new(template.GenLib)
	genLib.Primary = "en_us_goa"
	genLib.AddFallBack("en_us_dedes")

	resolver.tkMap = make(map[string]bool)
	resolver.values = make(map[string]string)
	resolver.tkMap["template.titles/ve.pdf.header"] = true
	resolver.tkMap["da.d2/daVE.OnTheEveningBefore"] = true
	resolver.tkMap["actors/Deacon"] = true
	resolver.tkMap["rubrical/InALowVoice"] = true
	resolver.tkMap["rubrical/Thrice"] = true
}
func TestLML_Errors(t *testing.T) {
	input := `ID == "xy"
Type = "services"
Status = "drafty"
Month = 66
Day = 33
if BadID == 1 { insert "actors/Priest"" }`
	lml, err := NewLMLParser("a/b",input, genLibs, &resolver)
	if err != nil {
		t.Errorf("%s", err)
	}
	for _, error := range lml.WalkTemplate() {
		fmt.Println(error.StringVerbose())
	}
}
func TestTemplate(t *testing.T) {
	input := `
ID = "ages/Dated-Services/m01/d06/se.m01.d06.li" // used as id in database
Type = "service" // types are block, book, or service. 
Status = "draft" // status values are na, draft, review, or final. 
Calendar = "Gregorian" // if omitted is set to Gregorian.  Can be set to Julian.
HtmlCss = "ages/css/html.css"
PdfCss = "ages/css/pdf.css"
Title = "Divine Liturgy" // optional. Used for searching db for template
Month = 1 // required if template type is a service
Day = 6  // required if template type is a service
Year = 2020 // optional for a service. If year is omitted or set to zero, the current year will be used.
PageHeader = center @Lookup sid "template.titles/ve.pdf.header" rid "da/daVE.OnTheEveningBefore" lang 1
PageFooter = left @PageNbr right @Date // You can use one or more of left, center, and right.
SetPageNumber = 1 
{
// 1a
//p.actor sid "actors/Priest"
// 1b
//p.actor sid "actors/Priest" sid "rubrical/InALowVoice"
// 2a
//p.dialog span.it sid "prayers/res04p"
// 2b
//p.dialog span.it sid "prayers/res04p" span.rubric sid "rubrical/Thrice"
// 3
//p.actor span.rditbd sid "actors/Deacon" (span.rubric sid "rubrical/InALowVoice") span.bk nid "But, loud enough to be heard."
// 4
//p.actor sid "actors/Deacon" (span.rubric sid "rubrical/InALowVoice") span.bk nid "But, loud enough to be heard."
// 5
//p.actor sid "actors/Deacon" (span.rubric sid "rubrical/InALowVoice") nid "But, loud enough to be heard."
// 6
//p.actor sid "actors/Deacon" (span.rubric sid "rubrical/InALowVoice" span.rubric sid "rubrical/Thrice") nid "But, loud enough to be heard."
// 7
//p.actor span.it sid "actors/Deacon" (span.rubric sid "rubrical/InALowVoice" span.rubric sid "rubrical/Thrice") nid "But, loud enough to be heard."
// 8
//p.actor span.it sid "actors/Deacon" (span.rubric sid "rubrical/InALowVoice"  sid "rubrical/Thrice") nid "But, loud enough to be heard."
// 9
p.actor span.it sid "actors/Deacon" (span.rubric sid "rubrical/InALowVoice" ( span.bl sid "rubrical/Thrice" ) ) nid "But, loud enough to be heard."


 // p.actor sid "actors/Priest"
  //p.dialog sid "eu.lichrysbasil/euLI.Key0109.text"
 // p.dialog span.it sid "prayers/res04p"
  //p.hymn rid "oc.*/ocVE.ApolTheotokionVM.text" @Mode 1 @Day 1 @Ver // @Mode and @Day override the mode and day for the service date.
  //p.actor span.rditbd sid "actors/Deacon" span.rubric sid "rubrical/InALowVoice" span.bk nid "But, loud enough to be heard."
  // p.actor span.rditbd sid "actors/Deacon" (span.rubric sid "rubrical/InALowVoice") span.bk nid "But, loud enough to be heard."
}`
	lml, err := NewLMLParser("a/b",input, genLibs, &resolver)
	if err != nil {
		t.Errorf("%s", err)
	}
	for _, error := range lml.WalkTemplate() {
		t.Error(error.StringVerbose())
	}
	if lml.Listener.ATEM.ID != "ages/Dated-Services/m01/d06/se.m01.d06.li" {
		t.Errorf(fmt.Sprintf("got ID = %s, expected %s", lml.Listener.ATEM.ID, "ages/Dated-Services/m01/d06/se.m01.d06.li"))
	}
	if lml.Listener.ATEM.Type != templateTypes.Service {
		t.Errorf(fmt.Sprintf("got Type = %s, expected %s", lml.Listener.ATEM.Type.String(), templateTypes.Service))
	}
	if lml.Listener.ATEM.Status != statuses.Draft {
		t.Errorf(fmt.Sprintf("got Status = %s, expected %s", lml.Listener.ATEM.Status.String(), statuses.Draft))
	}
	if lml.Listener.ATEM.HtmlCss != "ages/css/html.css" {
		t.Errorf(fmt.Sprintf("got HtmlCss = %s, expected %s", lml.Listener.ATEM.HtmlCss, "ages/css/html.css"))
	}
	if lml.Listener.ATEM.PDF.CSS != "ages/css/pdf.css" {
		t.Errorf(fmt.Sprintf("got PdfCss = %s, expected %s", lml.Listener.ATEM.PDF.CSS, "ages/css/pdf.css"))
	}
	if lml.Listener.ATEM.PDF.Title != "Divine Liturgy" {
		t.Errorf(fmt.Sprintf("got Title = %s, expected %s", lml.Listener.ATEM.PDF.Title, "Divine Liturgy"))
	}
	if lml.Listener.ATEM.Month != 1 {
		t.Errorf(fmt.Sprintf("got Month = %d, expected %d", lml.Listener.ATEM.Month, 1))
	}
	if lml.Listener.ATEM.Day != 6 {
		t.Errorf(fmt.Sprintf("got Day = %d, expected %d", lml.Listener.ATEM.Day, 6))
	}
	if lml.Listener.ATEM.Year != 2020 {
		t.Errorf(fmt.Sprintf("got Year = %d, expected %d", lml.Listener.ATEM.Year, 2020))
	}
	if lml.Listener.ATEM.PDF.PageNbr != 1 {
		t.Errorf(fmt.Sprintf("got SetPageNumber = %d, expected %d", lml.Listener.ATEM.PDF.PageNbr, 1))
	}
	if lml.Listener.ATEM.LDP.TheDay.Year() != 2020 {
		t.Errorf(fmt.Sprintf("got LDP year = %d, expected %d", lml.Listener.ATEM.LDP.TheDay.Year(), 2020))
	}
	if lml.Listener.ATEM.LDP.TheDay.Month() != 1 {
		t.Errorf(fmt.Sprintf("got LDP month = %d, expected %d", lml.Listener.ATEM.LDP.TheDay.Month(), 1))
	}
	if lml.Listener.ATEM.LDP.TheDay.Day() != 6 {
		t.Errorf(fmt.Sprintf("got LDP day = %d, expected %d", lml.Listener.ATEM.LDP.TheDay.Day(), 6))
	}
	json, err := json.MarshalIndent(lml.Listener.ATEM, "", " ")
	fmt.Println(string(json))
}
func TestLML_Tokens(t *testing.T) {
	input := `ID = "x/y"`
	lml, err := NewLMLParser("x/y", input, genLibs, &resolver)
	if err != nil  {
		t.Errorf("%s", err)
	}
	for _, t := range lml.Tokens() {
		fmt.Printf("%d:%d %s (%q)\n", t.GetLine(), t.GetColumn(),
			lml.Lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	}
	for _, e := range lml.ErrorListener.Errors {
		fmt.Println(e)
	}
}
