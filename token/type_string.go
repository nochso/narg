// Code generated by "stringer -type Type"; DO NOT EDIT

package token

import "fmt"

const _Type_name = "EOFLinefeedWhitespaceCommentBraceOpenBraceCloseQuotedValueUnquotedValueInvalidValueMissingClosingQuoteInvalidValueMissingSeparator"

var _Type_index = [...]uint8{0, 3, 11, 21, 28, 37, 47, 58, 71, 102, 130}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return fmt.Sprintf("Type(%d)", i)
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
