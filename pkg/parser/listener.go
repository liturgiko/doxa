package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/liturgiko/doxa/pkg/enums/calendarTypes"
	"github.com/liturgiko/doxa/pkg/enums/idTypes"
	"github.com/liturgiko/doxa/pkg/enums/positions"
	"github.com/liturgiko/doxa/pkg/enums/statuses"
	"github.com/liturgiko/doxa/pkg/enums/templateTypes"
	"github.com/liturgiko/doxa/pkg/template"
	lml "gitlab.com/ocmc/liturgiko/lml-go/parser"
	"strconv"
	"strings"
)

/*
LMListener implements the BaseLMLListener for the Liturgical Markup Language.
 */
type LMLListener struct {
	lml.BaseLMLListener
	ATEM      *template.ATEM
	Resolver  Resolver
}
// A Resolver provides an methods for obtaining database values.
// The Close method tells the resolver to close its database connection.
// The ExistsTK method returns true if there is at least one record in the database that has this topic/key.
// The Values method gets the value of each combination of GenLib and topic/key.
// The Versions method gets the value of each combination of GenLib and the topic/key 'properties/version.designation'
type Resolver interface {
	Close()
	ExistsTK(topicKey string) bool
	Values(topicKey string, genLibs []template.GenLib) template.TKVal
	Versions(genLibs []template.GenLib) map[string]string
}

func NewLMLListener(genLibs []template.GenLib, resolver Resolver) (*LMLListener, error) {
	l := new(LMLListener)
	l.Resolver = resolver
	l.ATEM = template.NewATEM()
	l.ATEM.Calendar = calendarTypes.Gregorian // can be overridden if set explicitly in template
	l.ATEM.PDF = new(template.PDF)
	l.ATEM.GenLibs = genLibs
	if resolver != nil && genLibs != nil {
		l.ATEM.Versions = resolver.Versions(genLibs)
	}
	return l, nil
}
func (l *LMLListener) VisitErrorNode(node antlr.ErrorNode) {
	fmt.Print("")
}
// VisitTerminal is called when a terminal node is visited.
func (l *LMLListener) VisitTerminal(node antlr.TerminalNode) {
	fmt.Print("")
}

// EnterEveryRule is called when any rule is entered.
func (l *LMLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Print("")
}

// ExitEveryRule is called when any rule is exited.
func (l *LMLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Print("")
}

// EnterTemplate is called when production template is entered.
func (l *LMLListener) EnterTemplate(ctx *lml.TemplateContext) {
	fmt.Print("")
}

// ExitTemplate is called when production template is exited.
func (l *LMLListener) ExitTemplate(ctx *lml.TemplateContext) {
	if l.Resolver != nil {l.Resolver.Close()}
}

// EnterProperty is called when production property is entered.
func (l *LMLListener) EnterProperty(ctx *lml.PropertyContext) {
	fmt.Print("")
}

// ExitProperty is called when production property is exited.
func (l *LMLListener) ExitProperty(ctx *lml.PropertyContext) {
	fmt.Print("")
}

// EnterPropertyBlock is called when production propertyBlock is entered.
func (l *LMLListener) EnterPropertyBlock(ctx *lml.PropertyBlockContext) {
	fmt.Print("")
}

// ExitPropertyBlock is called when production propertyBlock is exited.
func (l *LMLListener) ExitPropertyBlock(ctx *lml.PropertyBlockContext) {
	fmt.Print("")
}

// EnterStatement is called when production statement is entered.
func (l *LMLListener) EnterStatement(ctx *lml.StatementContext) {
	fmt.Print("")
}

// ExitStatement is called when production statement is exited.
func (l *LMLListener) ExitStatement(ctx *lml.StatementContext) {
	fmt.Print("")
}

// EnterInsert is called when production insert is entered.
func (l *LMLListener) EnterInsert(ctx *lml.InsertContext) {
	if ctx.STRING() != nil {
		id,err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.STRING().GetSymbol(),nil)
		}
		// TODO: add check to ensure the id is a valid path to a file (if using file based templates) or is a valid database ID if storing templates in a database.
		l.ATEM.ID = id
	}
}

// ExitInsert is called when production insert is exited.
func (l *LMLListener) ExitInsert(ctx *lml.InsertContext) {
	fmt.Print("")
}

// EnterPara is called when production para is entered.
/*
para: PARA_STYLE ( nid | rid | sid | span )+ INSERT_VER?;
span: SPAN_STYLE ( nid | rid | sid | span )+;
 */
func (l *LMLListener) EnterPara(ctx *lml.ParaContext) {
	paragraph = new(template.Paragraph)
	paragraph.Class = ctx.PARA_STYLE().GetText()
	if ctx.INSERT_VER() != nil {
		paragraph.AddVersion()
	}
	span = nil
}

// ExitPara is called when production para is exited.
func (l *LMLListener) ExitPara(ctx *lml.ParaContext) {
	if span != nil {
		paragraph.AddSpan(*span)
	}
	l.ATEM.AddParagraph(*paragraph)
	span = nil
	fmt.Print("")
}

// EnterSpan is called when production span is entered.
/*
para: PARA_STYLE ( nid | rid | sid | span )+ INSERT_VER?;
span: SPAN_STYLE ( nid | rid | sid | span )+;
Ideas:
- Maybe we need a span stack.
	If the global span is nil, create one.
	Else push the current span onto the stack.
	For ExitPara, pop the stack (which should have len == 1) and add
	the span to the para.
 */
func (l *LMLListener) EnterSpan(ctx *lml.SpanContext) {
	span = new(template.Span)
	span.Class = ctx.SPAN_STYLE().GetText()
	fmt.Print("")
}
// ExitSpan is called when production span is exited.
func (l *LMLListener) ExitSpan(ctx *lml.SpanContext) {
	if pspan != nil {
		pspan.AddChildSpan(*span)
	} else {
		paragraph.AddSpan(*span)
	}
	span = nil
	fmt.Print("")
}

// EnterMedia is called when production media is entered.
func (l *LMLListener) EnterMedia(ctx *lml.MediaContext) {
	fmt.Print("")
}

// ExitMedia is called when production media is exited.
func (l *LMLListener) ExitMedia(ctx *lml.MediaContext) {
	fmt.Print("")
}

// EnterNid is called when production nid is entered.
func (l *LMLListener) EnterNid(ctx *lml.NidContext) {
	if ctx.STRING() == nil {
		ctx.GetParser().NotifyErrorListeners("nid value cannot be empty",ctx.GetStart(),nil)
	} else {
		if value, err := strconv.Unquote(ctx.STRING().GetText()); err == nil  {
			nid := template.NewNid(value)
			if span == nil {
				if pspan == nil {
					paragraph.AddSpan(*nid)
				} else {
					pspan.AddChildSpan(*nid)
				}
			} else {
				span.AddChildSpan(*nid)
			}
			//if pspan == nil {
			//	if span == nil {
			//		paragraph.AddSpan(*nid)
			//	} else {
			//		span.AddChildSpan(*nid)
			//	}
			//} else {
			//	if span == nil  {
			//		pspan.AddChildSpan(*nid)
			//	} else {
			//		span.AddChildSpan(*nid)
			//	}
			//}
		} else {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%v",err),ctx.GetStart(),nil)
		}
	}
	fmt.Print("")
}

// ExitNid is called when production nid is exited.
func (l *LMLListener) ExitNid(ctx *lml.NidContext) {
	fmt.Print("")
}

/**
Enter 3
 */
// EnterPspan is called when production pspan is entered.
func (l *LMLListener) EnterPspan(ctx *lml.PspanContext) {
	if span != nil {
		spans.Push(span)
	}
	if pspan != nil {
		pspans.Push(pspan)
	}
	pspan = new(template.Span)
}
// ExitPspan is called when production pspan is exited.
func (l *LMLListener) ExitPspan(ctx *lml.PspanContext) {
	if item, ok := spans.Pop(); ok {
		popped := item.(*template.Span)
		for _, s := range pspan.ChildSpans {
			popped.AddChildSpan(s)
		}
		span = popped
	}
	if pspans.Empty() {
		pspan = nil
	} else {
		if item, ok := pspans.Pop(); ok {
			popped := item.(*template.Span)
			pspan = popped
		}
	}
}

// EnterPosition is called when production position is entered.
func (l *LMLListener) EnterPosition(ctx *lml.PositionContext) {
	buildingLeft = false
	buildingCenter = false
	buildingRight =false

	if ctx.PositionType() == nil {
		ctx.GetParser().NotifyErrorListeners("nil position error",ctx.GetStart(),nil)
	} else {
		value := ctx.PositionType().GetText()
		value = strings.Title(strings.ToLower(value))
		slotPosition, err := positions.PositionString(value)
		if err != nil {
			msg := fmt.Sprintf("invalid header/footer slot position \"%s\", expected one of %v", value, positions.PositionValues())
			ctx.GetParser().NotifyErrorListeners(msg,ctx.GetStart(),nil)
		} else {
			switch slotPosition {
			case positions.Left:
				buildingLeft = true
			case positions.Center:
				buildingCenter = true
			case positions.Right:
				buildingRight = true
			}
		}
	}
}

// ExitPosition is called when production position is exited.
func (l *LMLListener) ExitPosition(ctx *lml.PositionContext) {
	buildingLeft = false
	buildingCenter = false
	buildingRight = false
}
func AddDirective(d template.PDFDecorator) {
	if buildingLeft {
		if buildingHeader {
			pageHeader.AddLeftDirective(d)
		} else {
			pageFooter.AddLeftDirective(d)
		}
	} else if buildingCenter {
		if buildingHeader {
			pageHeader.AddCenterDirective(d)
		} else {
			pageFooter.AddCenterDirective(d)
		}
	} else { // buildingRight
		if buildingHeader {
			pageHeader.AddRightDirective(d)
		} else {
			pageFooter.AddRightDirective(d)
		}
	}

}

// EnterDirective is called when production directive is entered.
func (l *LMLListener) EnterDirective(ctx *lml.DirectiveContext) {
	if ctx.INSERT_DATE() != nil {
		dir := template.NewDateDirective("span.date", l.ATEM.LDP.TheDay)
		AddDirective(dir)
	}
	if ctx.INSERT_PAGE_NUMBER() != nil {
		AddDirective(template.NewPageNbrDirective("span.pageNbr"))
	}
}

// ExitDirective is called when production directive is exited.
func (l *LMLListener) ExitDirective(ctx *lml.DirectiveContext) {
	fmt.Print("")
}

// EnterLookup is called when production lookup is entered.
func (l *LMLListener) EnterLookup(ctx *lml.LookupContext) {
	buildingLookup = true
	lib, err := strconv.Atoi(ctx.INTEGER().GetText())
	if err != nil || (lib == 0 || lib > 3) {
		msg := fmt.Sprintf("invalid language number %s, expected 1, 2, or 3", ctx.INTEGER().GetText())
		ctx.GetParser().NotifyErrorListeners(msg,ctx.GetStart(),nil)
		lookupDirective = template.NewLookupDirective(-1)
	} else {
		lookupDirective = template.NewLookupDirective(lib)
	}
}

// ExitLookup is called when production lookup is exited.
func (l *LMLListener) ExitLookup(ctx *lml.LookupContext) {
	if buildingFooter {
		if buildingLeft {
			pageFooter.AddLeftDirective(lookupDirective)
		} else if buildingCenter {
			pageFooter.AddCenterDirective(lookupDirective)
		} else {
			pageFooter.AddRightDirective(lookupDirective)
		}
	}
	if buildingHeader {
		if buildingLeft {
			pageHeader.AddLeftDirective(lookupDirective)
		} else if buildingCenter {
			pageHeader.AddCenterDirective(lookupDirective)
		} else {
			pageHeader.AddRightDirective(lookupDirective)
		}
	}
	buildingLookup = false
}

// EnterRid is called when production rid is entered.
func (l *LMLListener) EnterRid(ctx *lml.RidContext) {
	var dayOverride, modeOverride int
	var err error
	for _, override := range ctx.AllOverride() {
		o := override.GetStart().GetText()
		switch o {
		case "@Day":
			dayOverride, err = strconv.Atoi(ctx.GetStop().GetText())
			if err != nil {
				ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("day override error %v",err),override.GetStop(),nil)
			}
			if dayOverride < 1 || dayOverride > 7 {
				ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("expected a value between 1 and 7 for day override, but got %d",dayOverride),ctx.GetStop(),nil)
			}
		case "@Mode":
			modeOverride, err = strconv.Atoi(ctx.GetStop().GetText())
			if err != nil {
				ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("mode override error %v",err),ctx.GetStop(),nil)
			}
			if modeOverride < 1 || modeOverride > 8 {
				ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("expected a value between 1 and 8 for mode override, but got %d",modeOverride),ctx.GetStop(),nil)
			}
		default:
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("expected @Day or @Mode, but got %s",o),ctx.GetStop(),nil)
		}
	}
	if ctx.STRING() != nil {
		id,err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.STRING().GetSymbol(),nil)
		}
		parts := strings.Split(id, "/")
		switch len(parts) {
		case 1:
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("mismatched input '%s' expecting at least one forward slash in topic/key path",id),ctx.STRING().GetSymbol(),nil)
		case 2:
			if modeOverride > 0 || dayOverride > 0 {
				if ! strings.HasPrefix(parts[0],"oc") {
					ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("rid directives (@Mode or @Day) may only be used for topics starting with 'oc' (i.e. Octoechos)"),ctx.STRING().GetSymbol(),nil)
				}
			}
			relativeTopic := l.ATEM.LDP.RelativeTopic(parts[0], modeOverride, dayOverride)
			if l.Resolver != nil {
				if ! l.Resolver.ExistsTK(relativeTopic + template.IDPathDelimiter + parts[1]) {
					ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("not found topic/key '%s' does not exist in topic-key rings",id),ctx.STRING().GetSymbol(),nil)
				}
			}
		default:
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("mismatched input '%s' expecting only one forward slash in topic/key path",id),ctx.STRING().GetSymbol(),nil)
		}
		if l.Resolver != nil {
			fmt.Println(fmt.Sprintf("%s: %v",id, l.ATEM.GenLibs))
			l.ATEM.AddTKValues(id, l.Resolver.Values(id, l.ATEM.GenLibs))
		}
		if buildingLookup {
			lookupDirective.AddLookupTK(idTypes.RID, "", id)
		} else {
			rid := template.NewRid(id, modeOverride, dayOverride)
			if span == nil {
				if pspan == nil {
					paragraph.AddSpan(*rid)
				} else {
					pspan.AddChildSpan(*rid)
				}
			} else {
				span.AddChildSpan(*rid)
			}
		}
	}
}

// ExitRid is called when production rid is exited.
func (l *LMLListener) ExitRid(ctx *lml.RidContext) {
	fmt.Print("")
}

// EnterSid is called when production sid is entered.
func (l *LMLListener) EnterSid(ctx *lml.SidContext) {
	if ctx.STRING() != nil {
		id,err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.STRING().GetSymbol(),nil)
		}
		parts := strings.Split(id, "/")
		switch len(parts) {
		case 1:
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("mismatched input '%s' expecting at least one forward slash in topic/key path",id),ctx.STRING().GetSymbol(),nil)
		case 2:
			if l.Resolver != nil {
				if !l.Resolver.ExistsTK(id) {
					ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("not found topic/key '%s' does not exist in topic-key rings", id), ctx.STRING().GetSymbol(), nil)
				}
			}
		default:
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("mismatched input '%s' expecting only one forward slash in topic/key path",id),ctx.STRING().GetSymbol(),nil)
		}
		if l.Resolver != nil {
			l.ATEM.AddTKValues(id, l.Resolver.Values(id, l.ATEM.GenLibs))
		}
		if buildingLookup {
			lookupDirective.AddLookupTK(idTypes.SID, "", id)
		} else {
			sid := template.NewSid(id)
			if span == nil {
				if pspan == nil {
					paragraph.AddSpan(*sid)
				} else {
					pspan.AddChildSpan(*sid)
				}
			} else {
				span.AddChildSpan(*sid)
			}
			//if pspan == nil {
			//	if span == nil {
			//		paragraph.AddSpan(*sid)
			//	} else {
			//		span.AddChildSpan(*sid)
			//	}
			//} else {
			//	if span == nil  {
			//		pspan.AddChildSpan(*sid)
			//	} else {
			//		span.AddChildSpan(*sid)
			//	}
			//}
			fmt.Print("")
		}
	}
}

// ExitSid is called when production sid is exited.
func (l *LMLListener) ExitSid(ctx *lml.SidContext) {
	fmt.Print("")
}
func (l *LMLListener) EnterTmplCalendar(ctx *lml.TmplCalendarContext) {
	if ctx != nil {
		value, err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%v ",err),ctx.GetStart(),nil)
		} else {
			switch 	strings.ToLower(value) {
			case "gregorian":
				l.ATEM.Calendar = calendarTypes.Gregorian
			case "julian":
				l.ATEM.Calendar = calendarTypes.Julian
			default:
				ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("invalid calendar type '%s'. Expected one of %v",value,calendarTypes.CalendarTypeValues()),ctx.STRING().GetSymbol(),nil)
			}
		}
	}
}

// ExitTmplCalendar is called when production tmplCalendar is exited.
func (l *LMLListener) ExitTmplCalendar(ctx *lml.TmplCalendarContext) {}

func (l *LMLListener) EnterTmplHtmlCss(ctx *lml.TmplHtmlCssContext) {
	if ctx.STRING() != nil {
		value, err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.STRING().GetSymbol(),nil)
		}
		// TODO: check for existence of css
		l.ATEM.HtmlCss = value
	}
}

func (l *LMLListener) ExitTmplHtmlCss(ctx *lml.TmplHtmlCssContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplPdfCss(ctx *lml.TmplPdfCssContext) {
	if ctx.STRING() != nil {
		value, err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.STRING().GetSymbol(),nil)
		}
		// TODO: check for existence of css
		l.ATEM.PDF.CSS = value
	}
}

func (l *LMLListener) ExitTmplPdfCss(ctx *lml.TmplPdfCssContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplDay(ctx *lml.TmplDayContext) {
	if ctx.INTEGER() != nil {
		value, err := strconv.Atoi(ctx.INTEGER().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s", err), ctx.INTEGER().GetSymbol(), nil)
		}
		if value > 0 && value < 32 {
			l.ATEM.Day = value
			l.ATEM.SetLDP()
		} else {
			msg := fmt.Sprintf("invalid Day %d, expected value between 1 and 31", value)
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s", msg), ctx.INTEGER().GetSymbol(), nil)
		}
	}
}

func (l *LMLListener) ExitTmplDay(ctx *lml.TmplDayContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplID(ctx *lml.TmplIDContext) {
	// if ctx.STRING is nil, it will be reported by the ErrorListener
	if ctx.STRING() != nil {
		id,err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.STRING().GetSymbol(),nil)
		}
		if len(id) == 0 {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("mismatched input '%s' expecting value between quotes",id),ctx.STRING().GetSymbol(),nil)
		}
		if ! strings.Contains(id, "/") {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("mismatched input '%s' expecting at least one forward slash in ID path",id),ctx.STRING().GetSymbol(),nil)
		}
		l.ATEM.ID = id
	}
}

func (l *LMLListener) ExitTmplID(ctx *lml.TmplIDContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplMonth(ctx *lml.TmplMonthContext) {
	if ctx.INTEGER() != nil {
		value , err := strconv.Atoi(ctx.INTEGER().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.INTEGER().GetSymbol(),nil)
		}
		if value > 0 && value < 13 {
			l.ATEM.Month = value
		} else {
			msg := fmt.Sprintf("invalid Month %d, expected value between 1 and 12",value)
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",msg),ctx.INTEGER().GetSymbol(),nil)
		}
	}
}

func (l *LMLListener) ExitTmplMonth(ctx *lml.TmplMonthContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplPageHeader(ctx *lml.TmplPageHeaderContext) {
	buildingHeader = true
	buildingFooter = false
	pageHeader = new(template.Header)
	pageHeader.Parity = template.Both
	buildingCenter = false
	buildingLeft = false
	buildingRight = false
}

func (l *LMLListener) ExitTmplPageHeader(ctx *lml.TmplPageHeaderContext) {
	l.ATEM.PDF.AddHeader(*pageHeader)
}

func (l *LMLListener) EnterTmplPageFooter(ctx *lml.TmplPageFooterContext) {
	buildingHeader = false
	buildingFooter = true
	pageFooter = new(template.Footer)
	pageFooter.Parity = template.Both
	buildingCenter = false
	buildingLeft = false
	buildingRight = false
}

func (l *LMLListener) ExitTmplPageFooter(ctx *lml.TmplPageFooterContext) {
	l.ATEM.PDF.AddFooter(*pageFooter)
}

func (l *LMLListener) EnterTmplPageHeaderEven(ctx *lml.TmplPageHeaderEvenContext) {
	fmt.Print("")
}

func (l *LMLListener) ExitTmplPageHeaderEven(ctx *lml.TmplPageHeaderEvenContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplPageFooterEven(ctx *lml.TmplPageFooterEvenContext) {
	fmt.Print("")
}

func (l *LMLListener) ExitTmplPageFooterEven(ctx *lml.TmplPageFooterEvenContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplPageHeaderOdd(ctx *lml.TmplPageHeaderOddContext) {
	fmt.Print("")
}

func (l *LMLListener) ExitTmplPageHeaderOdd(ctx *lml.TmplPageHeaderOddContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplPageFooterOdd(ctx *lml.TmplPageFooterOddContext) {
	fmt.Print("")
}

func (l *LMLListener) ExitTmplPageFooterOdd(ctx *lml.TmplPageFooterOddContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplPageNumber(ctx *lml.TmplPageNumberContext) {
	if ctx.INTEGER() != nil {
		value , err := strconv.Atoi(ctx.INTEGER().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.INTEGER().GetSymbol(),nil)
		}
		l.ATEM.PDF.PageNbr = value
	}
}

// ExitT_set_page_number is called when production t_set_page_number is exited.
func (l *LMLListener) ExitTmplPageNumber(ctx *lml.TmplPageNumberContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplStatus(ctx *lml.TmplStatusContext) {
	if ctx.STRING() != nil {
		raw, err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.STRING().GetSymbol(),nil)
		}
		value := strings.Title(strings.ToLower(raw))
		status, err := statuses.StatusString(value)
		if err != nil {
			msg := fmt.Sprintf("invalid template Status \"%s\", expected one of %s", raw, statuses.StatusValues())
			ctx.GetParser().NotifyErrorListeners(msg,ctx.STRING().GetSymbol(),nil)
		} else {
			l.ATEM.Status = status
		}
	}
}

func (l *LMLListener) ExitTmplStatus(ctx *lml.TmplStatusContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplTitle(ctx *lml.TmplTitleContext) {
	if ctx.STRING() != nil {
		value, err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.STRING().GetSymbol(),nil)
		}
		l.ATEM.PDF.Title = value
	}
}

func (l *LMLListener) ExitTmplTitle(ctx *lml.TmplTitleContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplType(ctx *lml.TmplTypeContext) {
	if ctx.STRING() != nil {
		raw, err := strconv.Unquote(ctx.STRING().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.STRING().GetSymbol(),nil)
		}
		value := strings.Title(strings.ToLower(raw))
		tmplType, err := templateTypes.TemplateTypeString(value)
		if err != nil {
			msg := fmt.Sprintf("invalid template Type \"%s\", expected one of %s", raw, templateTypes.TemplateTypeValues())
			ctx.GetParser().NotifyErrorListeners(msg,ctx.STRING().GetSymbol(),nil)
		} else {
			l.ATEM.Type = tmplType
		}
	}
}

func (l *LMLListener) ExitTmplType(ctx *lml.TmplTypeContext) {
	fmt.Print("")
}

func (l *LMLListener) EnterTmplYear(ctx *lml.TmplYearContext) {
	if ctx.INTEGER() != nil {
		value , err := strconv.Atoi(ctx.INTEGER().GetText())
		if err != nil {
			ctx.GetParser().NotifyErrorListeners(fmt.Sprintf("%s",err),ctx.INTEGER().GetSymbol(),nil)
		}
		l.ATEM.Year = value
		l.ATEM.SetLDP()
	}
}

func (l *LMLListener) ExitTmplYear(ctx *lml.TmplYearContext) {
	fmt.Print("")
}

// EnterBlock is called when production block is entered.
func (l *LMLListener) EnterBlock(ctx *lml.BlockContext) {
	fmt.Print("")
}

// ExitBlock is called when production block is exited.
func (l *LMLListener) ExitBlock(ctx *lml.BlockContext) {
	fmt.Print("")
}

// EnterSwitchBlock is called when production switchBlock is entered.
func (l *LMLListener) EnterSwitchBlock(ctx *lml.SwitchBlockContext) {
	fmt.Print("")
}

// ExitSwitchBlock is called when production switchBlock is exited.
func (l *LMLListener) ExitSwitchBlock(ctx *lml.SwitchBlockContext) {
	fmt.Print("")
}

// EnterSwitchBlockStatementGroup is called when production switchBlockStatementGroup is entered.
func (l *LMLListener) EnterSwitchBlockStatementGroup(ctx *lml.SwitchBlockStatementGroupContext) {
	fmt.Print("")
}

// ExitSwitchBlockStatementGroup is called when production switchBlockStatementGroup is exited.
func (l *LMLListener) ExitSwitchBlockStatementGroup(ctx *lml.SwitchBlockStatementGroupContext) {
	fmt.Print("")
}

// EnterSwitchLabel is called when production switchLabel is entered.
func (l *LMLListener) EnterSwitchLabel(ctx *lml.SwitchLabelContext) {
	fmt.Print("")
}

// ExitSwitchLabel is called when production switchLabel is exited.
func (l *LMLListener) ExitSwitchLabel(ctx *lml.SwitchLabelContext) {
	fmt.Print("")
}

// EnterIntegerExpression is called when production integerExpression is entered.
func (l *LMLListener) EnterIntegerExpression(ctx *lml.IntegerExpressionContext) {
	fmt.Print("")
}

// ExitIntegerExpression is called when production integerExpression is exited.
func (l *LMLListener) ExitIntegerExpression(ctx *lml.IntegerExpressionContext) {
	fmt.Print("")
}

// EnterIntegerList is called when production integerList is entered.
func (l *LMLListener) EnterIntegerList(ctx *lml.IntegerListContext) {
	fmt.Print("")
}

// ExitIntegerList is called when production integerList is exited.
func (l *LMLListener) ExitIntegerList(ctx *lml.IntegerListContext) {
	fmt.Print("")
}

// EnterDowExpression is called when production dowExpression is entered.
func (l *LMLListener) EnterDowExpression(ctx *lml.DowExpressionContext) {
	fmt.Print("")
}

// ExitDowExpression is called when production dowExpression is exited.
func (l *LMLListener) ExitDowExpression(ctx *lml.DowExpressionContext) {
	fmt.Print("")
}

// EnterDowList is called when production dowList is entered.
func (l *LMLListener) EnterDowList(ctx *lml.DowListContext) {
	fmt.Print("")
}

// ExitDowList is called when production dowList is exited.
func (l *LMLListener) ExitDowList(ctx *lml.DowListContext) {
	fmt.Print("")
}

// EnterLdpInt is called when production ldpInt is entered.
func (l *LMLListener) EnterLdpInt(ctx *lml.LdpIntContext) {
	fmt.Print("")
}

// ExitLdpInt is called when production ldpInt is exited.
func (l *LMLListener) ExitLdpInt(ctx *lml.LdpIntContext) {
	fmt.Print("")
}

// EnterMonthDay is called when production monthDay is entered.
func (l *LMLListener) EnterMonthDay(ctx *lml.MonthDayContext) {
	fmt.Print("")
}

// ExitMonthDay is called when production monthDay is exited.
func (l *LMLListener) ExitMonthDay(ctx *lml.MonthDayContext) {
	fmt.Print("")
}

// EnterExpression is called when production expression is entered.
func (l *LMLListener) EnterExpression(ctx *lml.ExpressionContext) {
	fmt.Println(ctx)
}

// ExitExpression is called when production expression is exited.
func (l *LMLListener) ExitExpression(ctx *lml.ExpressionContext) {
	fmt.Print("")
}

var pageHeader *template.Header
var pageFooter *template.Footer
var lookupDirective *template.PDFLookupDirective
var paragraph *template.Paragraph
var span *template.Span
var pspan *template.Span
var spans = arraystack.New()
var pspans = arraystack.New()
var buildingHeader, buildingFooter, buildingLeft, buildingCenter, buildingRight, buildingLookup bool