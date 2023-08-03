// Code generated by "enumer -type=TokenType"; DO NOT EDIT.

package tokens

import (
	"fmt"
	"strings"
)

const _TokenTypeName = "ILLEGALEOFIDENTINTEGERSTRINGASSIGNPLUSMINUSBANGASTERISKSLASHLTLTEQGTGTEQEQUALNOTEQUALORANDCOMMASEMICOLONLPRARENTRPARENTLBRACERBRACEFUNCTIONLETIFELSERETURNTRUEFALSE"

var _TokenTypeIndex = [...]uint8{0, 7, 10, 15, 22, 28, 34, 38, 43, 47, 55, 60, 62, 66, 68, 72, 77, 85, 87, 90, 95, 104, 112, 119, 125, 131, 139, 142, 144, 148, 154, 158, 163}

const _TokenTypeLowerName = "illegaleofidentintegerstringassignplusminusbangasteriskslashltlteqgtgteqequalnotequalorandcommasemicolonlprarentrparentlbracerbracefunctionletifelsereturntruefalse"

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenTypeIndex)-1) {
		return fmt.Sprintf("TokenType(%d)", i)
	}
	return _TokenTypeName[_TokenTypeIndex[i]:_TokenTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _TokenTypeNoOp() {
	var x [1]struct{}
	_ = x[ILLEGAL-(0)]
	_ = x[EOF-(1)]
	_ = x[IDENT-(2)]
	_ = x[INTEGER-(3)]
	_ = x[STRING-(4)]
	_ = x[ASSIGN-(5)]
	_ = x[PLUS-(6)]
	_ = x[MINUS-(7)]
	_ = x[BANG-(8)]
	_ = x[ASTERISK-(9)]
	_ = x[SLASH-(10)]
	_ = x[LT-(11)]
	_ = x[LTEQ-(12)]
	_ = x[GT-(13)]
	_ = x[GTEQ-(14)]
	_ = x[EQUAL-(15)]
	_ = x[NOTEQUAL-(16)]
	_ = x[OR-(17)]
	_ = x[AND-(18)]
	_ = x[COMMA-(19)]
	_ = x[SEMICOLON-(20)]
	_ = x[LPRARENT-(21)]
	_ = x[RPARENT-(22)]
	_ = x[LBRACE-(23)]
	_ = x[RBRACE-(24)]
	_ = x[FUNCTION-(25)]
	_ = x[LET-(26)]
	_ = x[IF-(27)]
	_ = x[ELSE-(28)]
	_ = x[RETURN-(29)]
	_ = x[TRUE-(30)]
	_ = x[FALSE-(31)]
}

var _TokenTypeValues = []TokenType{ILLEGAL, EOF, IDENT, INTEGER, STRING, ASSIGN, PLUS, MINUS, BANG, ASTERISK, SLASH, LT, LTEQ, GT, GTEQ, EQUAL, NOTEQUAL, OR, AND, COMMA, SEMICOLON, LPRARENT, RPARENT, LBRACE, RBRACE, FUNCTION, LET, IF, ELSE, RETURN, TRUE, FALSE}

var _TokenTypeNameToValueMap = map[string]TokenType{
	_TokenTypeName[0:7]:          ILLEGAL,
	_TokenTypeLowerName[0:7]:     ILLEGAL,
	_TokenTypeName[7:10]:         EOF,
	_TokenTypeLowerName[7:10]:    EOF,
	_TokenTypeName[10:15]:        IDENT,
	_TokenTypeLowerName[10:15]:   IDENT,
	_TokenTypeName[15:22]:        INTEGER,
	_TokenTypeLowerName[15:22]:   INTEGER,
	_TokenTypeName[22:28]:        STRING,
	_TokenTypeLowerName[22:28]:   STRING,
	_TokenTypeName[28:34]:        ASSIGN,
	_TokenTypeLowerName[28:34]:   ASSIGN,
	_TokenTypeName[34:38]:        PLUS,
	_TokenTypeLowerName[34:38]:   PLUS,
	_TokenTypeName[38:43]:        MINUS,
	_TokenTypeLowerName[38:43]:   MINUS,
	_TokenTypeName[43:47]:        BANG,
	_TokenTypeLowerName[43:47]:   BANG,
	_TokenTypeName[47:55]:        ASTERISK,
	_TokenTypeLowerName[47:55]:   ASTERISK,
	_TokenTypeName[55:60]:        SLASH,
	_TokenTypeLowerName[55:60]:   SLASH,
	_TokenTypeName[60:62]:        LT,
	_TokenTypeLowerName[60:62]:   LT,
	_TokenTypeName[62:66]:        LTEQ,
	_TokenTypeLowerName[62:66]:   LTEQ,
	_TokenTypeName[66:68]:        GT,
	_TokenTypeLowerName[66:68]:   GT,
	_TokenTypeName[68:72]:        GTEQ,
	_TokenTypeLowerName[68:72]:   GTEQ,
	_TokenTypeName[72:77]:        EQUAL,
	_TokenTypeLowerName[72:77]:   EQUAL,
	_TokenTypeName[77:85]:        NOTEQUAL,
	_TokenTypeLowerName[77:85]:   NOTEQUAL,
	_TokenTypeName[85:87]:        OR,
	_TokenTypeLowerName[85:87]:   OR,
	_TokenTypeName[87:90]:        AND,
	_TokenTypeLowerName[87:90]:   AND,
	_TokenTypeName[90:95]:        COMMA,
	_TokenTypeLowerName[90:95]:   COMMA,
	_TokenTypeName[95:104]:       SEMICOLON,
	_TokenTypeLowerName[95:104]:  SEMICOLON,
	_TokenTypeName[104:112]:      LPRARENT,
	_TokenTypeLowerName[104:112]: LPRARENT,
	_TokenTypeName[112:119]:      RPARENT,
	_TokenTypeLowerName[112:119]: RPARENT,
	_TokenTypeName[119:125]:      LBRACE,
	_TokenTypeLowerName[119:125]: LBRACE,
	_TokenTypeName[125:131]:      RBRACE,
	_TokenTypeLowerName[125:131]: RBRACE,
	_TokenTypeName[131:139]:      FUNCTION,
	_TokenTypeLowerName[131:139]: FUNCTION,
	_TokenTypeName[139:142]:      LET,
	_TokenTypeLowerName[139:142]: LET,
	_TokenTypeName[142:144]:      IF,
	_TokenTypeLowerName[142:144]: IF,
	_TokenTypeName[144:148]:      ELSE,
	_TokenTypeLowerName[144:148]: ELSE,
	_TokenTypeName[148:154]:      RETURN,
	_TokenTypeLowerName[148:154]: RETURN,
	_TokenTypeName[154:158]:      TRUE,
	_TokenTypeLowerName[154:158]: TRUE,
	_TokenTypeName[158:163]:      FALSE,
	_TokenTypeLowerName[158:163]: FALSE,
}

var _TokenTypeNames = []string{
	_TokenTypeName[0:7],
	_TokenTypeName[7:10],
	_TokenTypeName[10:15],
	_TokenTypeName[15:22],
	_TokenTypeName[22:28],
	_TokenTypeName[28:34],
	_TokenTypeName[34:38],
	_TokenTypeName[38:43],
	_TokenTypeName[43:47],
	_TokenTypeName[47:55],
	_TokenTypeName[55:60],
	_TokenTypeName[60:62],
	_TokenTypeName[62:66],
	_TokenTypeName[66:68],
	_TokenTypeName[68:72],
	_TokenTypeName[72:77],
	_TokenTypeName[77:85],
	_TokenTypeName[85:87],
	_TokenTypeName[87:90],
	_TokenTypeName[90:95],
	_TokenTypeName[95:104],
	_TokenTypeName[104:112],
	_TokenTypeName[112:119],
	_TokenTypeName[119:125],
	_TokenTypeName[125:131],
	_TokenTypeName[131:139],
	_TokenTypeName[139:142],
	_TokenTypeName[142:144],
	_TokenTypeName[144:148],
	_TokenTypeName[148:154],
	_TokenTypeName[154:158],
	_TokenTypeName[158:163],
}

// TokenTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func TokenTypeString(s string) (TokenType, error) {
	if val, ok := _TokenTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _TokenTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to TokenType values", s)
}

// TokenTypeValues returns all values of the enum
func TokenTypeValues() []TokenType {
	return _TokenTypeValues
}

// TokenTypeStrings returns a slice of all String values of the enum
func TokenTypeStrings() []string {
	strs := make([]string, len(_TokenTypeNames))
	copy(strs, _TokenTypeNames)
	return strs
}

// IsATokenType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i TokenType) IsATokenType() bool {
	for _, v := range _TokenTypeValues {
		if i == v {
			return true
		}
	}
	return false
}
