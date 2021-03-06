//go:generate enumer -type=DayOfSeason -json -text -yaml -sql
// 1) go get github.com/alvaroloes/enumer
// 2) in the enum subfolder for this enum: go generate

// Package dsos provides an enum of Liturgical Days of the Season (i.e. the Triodion)
package dsos

type DayOfSeason int
const (
	D01 DayOfSeason = iota + 1
	D02
	D03
	D04
	D05
	D06
	D07
	D08
	D09
	D10
	D11
	D12
	D13
	D14
	D15
	D16
	D17
	D18
	D19
	D20
	D21
	D22
	D23
	D24
	D25
	D26
	D27
	D28
	D29
	D30
	D31
	D32
	D33
	D34
	D35
	D36
	D37
	D38
	D39
	D40
	D41
	D42
	D43
	D44
	D45
	D46
	D47
	D48
	D49
	D50
	D51
	D52
	D53
	D54
	D55
	D56
	D57
	D58
	D59
	D60
	D61
	D62
	D63
	D64
	D65
	D66
	D67
	D68
	D69
	D70
)
