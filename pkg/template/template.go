// Package template provides an Abstract Template struct which contains all the information needed for a generator (e.g. an HTML generator or a PDF generator), when combined with the user requested libraries.
package template

import (
	"github.com/liturgiko/doxa/pkg/enums/calendarTypes"
	"github.com/liturgiko/doxa/pkg/enums/directiveTypes"
	"github.com/liturgiko/doxa/pkg/enums/idTypes"
	"github.com/liturgiko/doxa/pkg/enums/statuses"
	"github.com/liturgiko/doxa/pkg/enums/templateTypes"
	"github.com/liturgiko/doxa/pkg/ldp"
	"time"
)
const IDPathDelimiter = "/"
/*
ATEM is an Abstract Template.
Use the function NewATEM() to get a pointer to an ATEM.  Otherwise the maps will be nil.
ID is the identifier for the template and should match a corresponding path in the file system or an ID in a database.
Type is block, book, or service.
Status values are na, draft, review, or final.
Calendar indicates whether the Gregorian or Julian calendar is to be used with this template.
Month, Day, Year are used when the Type = service.
HtmlCss is the file path or database ID for the css file to be used for the template.
PDF contains information for creating PDF files, e.g. title, Header and Footer information, and the CSS to use.
LDP holds the liturgical day properties for the specified Year, Month, and Day.
GenLibs is a slice that holds the primary library and fallback libraries to be bound to each topic-key. The fallback libraries are used if the combination of the primary library and a topic/key does not exist.
Versions is a map whose key is a library and value is the library's acronym. The acronym is used as the value for an @Ver directive to insert the version.
Values is a map whose keys are topic/key and the value is a slice with the index of each element corresponding to the index in the GenLibs slice. The value retrieved from a database is stored as well.
Paragraphs holds an array of Paragraph. The information in a paragraph should be used by a generator of HTML or a PDF (or anything else that has rows), the information in a Paragraph is used 1..n times depending on how many libraries have been requested by the user. Each table row in an HTML document or PDF file will have a cell for each requested library, and content as specified by the paragraph.
 */
type ATEM struct {
	ID               string
	Type             templateTypes.TemplateType
	Status           statuses.Status
	Calendar         calendarTypes.CalendarType
	Month, Day, Year int
	HtmlCss          string
	PDF              *PDF
	LDP              ldp.LDP
	GenLibs			 []GenLib
	Versions 		 map[string]string
	Values           map[string]TKVal
	Paragraphs       []*Paragraph
}
func NewATEM() *ATEM {
	atem := new(ATEM)
	atem.Versions = make(map[string]string)
	atem.Values = make(map[string]TKVal)
	return atem
}
// SetLDPYMD sets the Liturgical Day Properties to the supplied month, day, and year
func (a *ATEM) SetLDPYMD(month, day, year int, calendarType calendarTypes.CalendarType) error {
	l, err := ldp.NewLDPYMD(year, month, day, calendarType)
	if err != nil {
		return err
	}
	a.LDP = l
	a.Year = year
	a.Month = month
	a.Day = day
	return nil
}
// AddLTKValue adds a value to the map of LTKValues if the key (topic/key) does not already exist.
func (a *ATEM) AddTKValues(topicKey string, tkVal TKVal) {
	if _, ok := a.Values[topicKey]; ok {
		// ignore request if already set for this topicKey
	} else {
		if &tkVal != nil {
			a.Values[topicKey] = tkVal
		}
	}
}
// SetLDPMD sets the Liturgical Day Properties to current year, and supplied month and day
func (a *ATEM) SetLDPMD(month, day int, calendarType calendarTypes.CalendarType) error {
	l, err := ldp.NewLDPMD(month, day, calendarType)
	if err != nil {
		return err
	}
	a.LDP = l
	a.Year = l.TheDay.Year()
	a.Month = month
	a.Day = day
	return nil
}
// SetLDP sets the Liturgical Day Properties to the year, month, and day in the template
func (a *ATEM) SetLDP() error {
	var l ldp.LDP
	var err error
	if a.Year == 0 {
		l, err = ldp.NewLDPMD(a.Month, a.Day, a.Calendar)
		if err != nil {
			return err
		}
	} else {
		l, err = ldp.NewLDPYMD(a.Year, a.Month, a.Day, a.Calendar)
		if err != nil {
			return err
		}
	}
	a.LDP = l
	a.Year = l.TheDay.Year()
	a.Month = int(l.TheDay.Month())
	a.Day = l.TheDay.Day()
	return nil
}
func (a *ATEM) AddParagraph(p Paragraph) {
	a.Paragraphs = append(a.Paragraphs, &p)
}
func (a *ATEM) AddGenLib(genLib GenLib) {
	a.GenLibs = append(a.GenLibs, genLib)
}
// LTKVal holds a resolved Library/Topic/Key and its value.
// Resolved means that it was retrieved from a database.
// Although an ID is library/topic/key, the library is split out
// to make it easy to get the Version acronym if needed.
// To get the complete ID, use the ID() method.
type LTKVal struct {
	Library string
	TopicKey string
	Value string
}
// ID returns a concatenation of Library and TopicKey as a valid lml
func (ltkv *LTKVal) ID() string {
	return ltkv.Library + IDPathDelimiter + ltkv.TopicKey
}
// TKVal holds an slice of LTKVal.  The length of the slice corresponds to the number of columns to be generated, which are actually cells in a row.
type TKVal struct {
	Values []LTKVal
}
func (tkVal *TKVal) Add(ltkVal LTKVal) {
	tkVal.Values = append(tkVal.Values, ltkVal)
}
// GenLib holds the user's request to generate using a specific library and a slice of fallback libraries to use of the primary library + topic-key does not exist.
type GenLib struct {
	Primary string
	FallBacks []string
}
// AddFallback appends a library to the splice of fallback libraries.
func (g *GenLib) AddFallBack(library string) {
	g.FallBacks = append(g.FallBacks, library)
}
// All flattens the GenLib by returning a string slice of the libraries.
// The primary library is at the zero index, the fallback libraries are the rest.
func (g *GenLib) All() []string {
	var all []string
	all = append(all, g.Primary)
	for _, library := range g.FallBacks {
		all = append(all, library)
	}
	return all
}
type PDF struct {
	CSS string
	PageNbr int
	Title string
	Headers []Header
	Footers []Footer
}
// AddHeader appends a header to the PDF struct's slice of Headers
func (p *PDF) AddHeader(header Header) {
	p.Headers = append(p.Headers, header)
}
// AddFooter appends a footer to the PDF struct's slice of Footers
func (p *PDF) AddFooter(footer Footer) {
	p.Footers = append(p.Footers, footer)
}
// Parity indicates whether a Header/Footer is even or odd or both
type Parity int
const (
	Both Parity = iota
	Even
	Odd
)
// Position designates slots in a header/footer. The slots are left, center, or right.
type Position int
const (
	Left Position = iota
	Center
	Right
)
// Header indicates parity (whether it is for even or odd pages or both)
// and the content of each of three slots: left, center, right.
// A slot can be empty.  To determine whether a slot has content, call the functions HasLeftSlot, HasCenterSlot, and HasRightSlot.
type Header struct {
	Parity Parity
	Left Slot
	Center Slot
	Right Slot
}
func (h *Header) AddLeftDirective(directive PDFDecorator) {
	h.Left.Directives = append(h.Left.Directives, directive)
}
func (h *Header) AddCenterDirective(directive PDFDecorator) {
	h.Center.Directives = append(h.Center.Directives, directive)
}
func (h *Header) AddRightDirective(directive PDFDecorator) {
	h.Right.Directives = append(h.Right.Directives, directive)
}
func (h *Header) HasLeftSlot() bool {
	return len(h.Left.Directives) > 0
}
func (h *Header) HasCenterSlot() bool {
	return len(h.Center.Directives) > 0
}
func (h *Header) HasRightSlot() bool {
	return len(h.Right.Directives) > 0
}
// NewFooter returns a header with parity set to Both.
func NewHeader() *Header {
	var header = new(Header)
	header.Parity = Both
	return header
}
// NewHeaderEven returns a header with parity set to Even.
// That is, it is for even numbered pages.
func NewHeaderEven() *Header {
	var header = new(Header)
	header.Parity = Even
	return header
}
// NewHeaderOdd returns a header with parity set to Odd.
// That is, it is for odd numbered pages.
func NewHeaderOdd() *Header {
	var header = new(Header)
	header.Parity = Odd
	return header
}
// Slot contains information about a left, center, or right slot of a header or footer
type Slot struct {
	Directives []PDFDecorator
}
func (s *Slot) AddDirective(directive PDFDecorator) {
	s.Directives = append(s.Directives, directive)
}
// Footer indicates parity (whether it is for even or odd pages or both)
// and the content of each of three slots: left, center, right.
// A slot can be empty.  To determine whether a slot has content, call the functions HasLeftSlot, HasCenterSlot, and HasRightSlot.
type Footer struct {
	Parity Parity
	Left Slot
	Center Slot
	Right Slot
}
// NewFooter returns a footer with parity set to Both.
func NewFooter() *Footer {
	var footer = new(Footer)
	footer.Parity = Both
	return footer
}
// NewFooterEven returns a footer with parity set to Even.
// That is, it is for even numbered pages.
func NewFooterEven() *Footer {
	var footer = new(Footer)
	footer.Parity = Even
	return footer
}
// NewFooterOdd returns a footer with parity set to Odd.
// That is, it is for odd numbered pages.
func NewFooterOdd() *Footer {
	var footer = new(Footer)
	footer.Parity = Odd
	return footer
}
func (f *Footer) AddLeftDirective(directive PDFDecorator) {
	f.Left.Directives = append(f.Left.Directives, directive)
}
func (f *Footer) AddCenterDirective(directive PDFDecorator) {
	f.Center.Directives = append(f.Center.Directives, directive)
}
func (f *Footer) AddRightDirective(directive PDFDecorator) {
	f.Right.Directives = append(f.Right.Directives, directive)
}
func (f *Footer) HasLeftSlot() bool {
	return len(f.Left.Directives) > 0
}
func (f *Footer) HasCenterSlot() bool {
	return len(f.Center.Directives) > 0
}
func (f *Footer) HasRightSlot() bool {
	return len(f.Right.Directives) > 0
}
type PDFDecorator interface {
	DirType() directiveTypes.DirectiveType
	DirClass() string
}
type PDFDateDecorator interface {
	PDFDecorator
	Value() time.Time
}
type PDFLiteralDecorator interface {
	PDFDecorator
	Value () string
}
type PDFLookupDecorator interface {
	PDFDecorator
	Value () Lookup
}
// PDFDirective indicates what is to be inserted into a header or footer.
// Type indicates the type of the directive.
// The Class string holds the name of the span CSS class to be applied.
//
type PDFDirective struct {
	Type    directiveTypes.DirectiveType
	Class   string
}
func (p *PDFDirective) DirType() directiveTypes.DirectiveType {
	return p.Type
}
func (p *PDFDirective) DirClass() string {
	return p.Class
}
type PDFDateDirective struct {
	PDFDirective
	Date time.Time
}
func (p *PDFDateDirective) Value() time.Time {
	return p.Date
}
type PDFLookupDirective struct {
	PDFDirective
	Lookup Lookup
}
func (p *PDFLookupDirective) Value() Lookup {
	return p.Lookup
}
type PDFLiteralDirective struct {
	PDFDirective
	Literal string
}
func (p *PDFLiteralDirective) Value() string {
	return p.Literal
}
func NewDateDirective(class string, time time.Time) *PDFDateDirective {
	var dir = new(PDFDateDirective)
	dir.Type = directiveTypes.InsertDate
	dir.Class = class
	dir.Date = time
	return dir
}
func NewPageNbrDirective(class string) *PDFDirective {
	var dir = new(PDFDirective)
	dir.Type = directiveTypes.InsertPageNbr
	dir.Class = class
	return dir
}
func NewLiteralDirective(class, literal string) *PDFLiteralDirective {
	var dir = new(PDFLiteralDirective)
	dir.Type = directiveTypes.InsertLiteral
	dir.Class = class
	dir.Literal = literal
	return dir
}
func NewLookupDirective(library int) *PDFLookupDirective {
	var dir = new(PDFLookupDirective)
	dir.Type = directiveTypes.InsertLookup
	var lookup = new(Lookup)
	lookup.Library = library
	dir.Lookup = *lookup
	return dir
}
func (p *PDFLookupDirective) AddLookupTK(idType idTypes.IDType, class, topicKey string) error {
	var err error
	var lookupID = new(LookupTopicKey)
	lookupID.Type = idType
	lookupID.Class = class
	lookupID.TopicKey = topicKey
	p.Lookup.TopicKeys = append(p.Lookup.TopicKeys,*lookupID)
	return err
}
// Lookup provides information to do a database lookup to insert the result in a header/footer.
type Lookup struct {
	TopicKeys []LookupTopicKey
	Library   int
}
// LookupTopicKey indicates the type of lookup (RID or SID) and the Topic-Key to use and the CSS style class to use.
type LookupTopicKey struct {
	Type     idTypes.IDType
	Class string
	TopicKey string
	OverrideDay int
	OverrideMode int
}
type Spanner interface {
	CssClass() string
}
type TextSpanner interface {
	Spanner
	TKType() idTypes.IDType
}
type TextLiteralSpanner interface {
	TextSpanner
}
type TextRidSpanner interface {
	TextSpanner
}
type TextSidSpanner interface {
	TextSpanner
}
/*
TextSpan contains information for the formatting of a span that will contain the value from a database read using an ID, or contains a literal value.
Class provides the CSS classname
Type indicates whether it is a nid, sid, or rid.
TextSpan
 */
type TextSpan struct {
	Span
	Type idTypes.IDType
}
func (tks *TextSpan) TKType() idTypes.IDType {
	return tks.Type
}
type LiteralSpan struct {
	TextSpan
	Value string
}
// SidTkSpan is used by the generator to use a specific ID for a library/topic/key database read.
// The library value is added by the generator.
type SidTkSpan struct {
	TextSpan
	TopicKey string
}
/*
RidTkSpan is used by the generator to compute a relative topic for a library/topic/key database read.
The library value is added by the generator.
OverrideDay is used for a RID. Zero = no override. Otherwise it is a value from 1-7, representing the days of the week
OverrideMode is used for a RID. Zero = no override.  Otherwise it is a value from 1-8, representing the liturgical modes
*/
type RidTkSpan struct {
	TextSpan
	TopicKey     string
	OverrideDay  int
	OverrideMode int
}
// Create ID returns a database ID for the ltx table by prefixing the topic/key with the library so it is complete for a database lookup.
func (r *RidTkSpan) CreateID(library string) string {
	return library + r.TopicKey
}
func (r *RidTkSpan) HasDayOverride() bool {
	return r.OverrideDay > 0
}
func (r *RidTkSpan) HasModeOverride() bool {
	return r.OverrideMode > 0
}

// Span contains the information for an inline text.
// Class is the CSS Class name.
// TopicKeys during generation are prefixed with a library and used to obtain a value from the ltx table in a database.
// For example, actor/Priest -> gr_gr_cog/actor/Priest or en_us_dedes/actors/Priest, etc.
// ChildSpans are spans embedded within a span.
// TextSpans and ChildSpans are mutually exclusive.
// TextSpans can be thought of as the terminal nodes of a span tree.
type Span struct {
	Class      string
	Type       idTypes.IDType
	// if Type = nid has:
	Literal string
	// if Type = sid or rid has:
	TopicKey string
	// if Type = rid can have:
	ModeOverride, DayOverride int
	ChildSpans []Span
}
func (s *Span) HasChildSpans() bool {
	return len(s.ChildSpans) > 0
}
// NewNid returns a Span with Type = NID, class "nid", and Literal set to the value parameter.
// AddNid returns an error if the span has ChildSpans.
// TextSpans and ChildSpans are mutually exclusive.
func NewNid(value string) *Span {
	s := new(Span)
	s.Type = idTypes.NID
	s.Class = "nid"
	s.Literal = value
	return s
}
// NewRid returns a Span with Type = RID, Class = "kvp", and TopicKey set to the topic/key parameter.
// At generation time, the TopicKey will be prefixed with the library so that an inspection of the HTML (for example) shows the ID used to get the value.
// The value will also be added at generation time.
func NewRid(topicKey string, modeOverride, dayOverride int) *Span {
	s := new(Span)
	s.Type = idTypes.RID
	s.Class = "kvp"
	s.TopicKey = topicKey
	s.ModeOverride = modeOverride
	s.DayOverride = dayOverride
	return s
}
// AddSid will create a SidTkSpan and set the class to "kvp", and the TopicKey to the topic/key parameter.
// At generation time, the TopicKey will be prefixed with the library so that an inspection of the HTML (for example) shows the ID used to get the value.
// The value will also be added at generation time.
func NewSid(topicKey string) *Span {
	s := new(Span)
	s.Type = idTypes.SID
	s.Class = "kvp"
	s.TopicKey = topicKey
	return s
}
// AddChildSpan adds a span as a child, and returns an error if the span has TopicKeys.
// TopicKeys and ChildSpans are mutually exclusive.
func (s *Span) AddChildSpan(span Span) error {
	s.ChildSpans = append(s.ChildSpans,span)
	return nil
}
// Paragraph holds information that can be used to create the contents of a cell of a table row.
// Class is the CSS class for the paragraph.
// Spans contain the information for creating inline texts within the paragraph.
// Version contains information to create a span that will contain an acronym for the library used to retrieve the text values.
type Paragraph struct {
	Class string
	Spans[]Span
	Version Span
}
func (p *Paragraph) AddSpan(span Span) {
	p.Spans = append(p.Spans,span)
}
// AddVersion creates a TextSpan whose class = "versiondesignation" and an inner sid TextSpan whose topic/key is "properties/version.designation"
// The version will display as an acronym for the library used during generation.
// For example, en_us_dedes will show [SD] as the version.
func (p *Paragraph) AddVersion() {
	vs := new(Span)
	vs.Type = idTypes.SID
	vs.Class = "versiondesignation"
	vs.TopicKey = "properties/version.designation"
	p.Version = *vs
}
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}