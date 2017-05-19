%{
// THIS IS A GENERATED FILE. DO NOT EDIT.
package parse

var TheSpecification *Node

%}

%token T_AMPERSAND
%token T_ANY
%token T_ASTERISK
%token T_ATTRIBUTE
%token T_BOOLEAN
%token T_CASE
%token T_CHAR
%token T_CHARACTER_LITERAL
%token T_CIRCUMFLEX
%token T_COLON
%token T_COMMA
%token T_CONST
%token T_CONTEXT
%token T_DEFAULT
%token T_DOUBLE
%token T_ENUM
%token T_EQUAL
%token T_EXCEPTION
%token T_FALSE
%token T_FIXED
%token T_FIXED_PT_LITERAL
%token T_FLOAT
%token T_FLOATING_PT_LITERAL
%token T_GREATER_THAN_SIGN
%token T_IDENTIFIER
%token T_IN
%token T_INOUT
%token T_INTEGER_LITERAL
%token T_INTERFACE
%token T_LEFT_CURLY_BRACKET
%token T_LEFT_PARANTHESIS
%token T_LEFT_SQUARE_BRACKET
%token T_LESS_THAN_SIGN
%token T_LONG
%token T_MINUS_SIGN
%token T_MODULE
%token T_OCTET
%token T_ONEWAY
%token T_OUT
%token T_PERCENT_SIGN
%token T_PLUS_SIGN
%token T_PRINCIPAL
%token T_RAISES
%token T_READONLY
%token T_RIGHT_CURLY_BRACKET
%token T_RIGHT_PARANTHESIS
%token T_RIGHT_SQUARE_BRACKET
%token T_SCOPE
%token T_SEMICOLON
%token T_SEQUENCE
%token T_SHIFTLEFT
%token T_SHIFTRIGHT
%token T_SHORT
%token T_SOLIDUS
%token T_STRING
%token T_STRING_LITERAL
%token T_PRAGMA
%token T_STRUCT
%token T_SWITCH
%token T_TILDE
%token T_TRUE
%token T_OBJECT
%token T_TYPEDEF
%token T_UNION
%token T_UNSIGNED
%token T_VERTICAL_LINE
%token T_VOID
%token T_WCHAR
%token T_WSTRING
%token T_UNKNOWN
%token T_ABSTRACT
%token T_VALUETYPE
%token T_TRUNCATABLE
%token T_SUPPORTS
%token T_CUSTOM
%token T_PUBLIC
%token T_PRIVATE
%token T_FACTORY
%token T_NATIVE
%token T_VALUEBASE

%union {
  value string
  nodePtr *Node
}

%type <nodePtr> definition definitions declarator declarators module type_dcl struct_type member_list member simple_declarator type_spec simple_type_spec constr_type_spec base_type_spec template_type_spec object_type
%type <nodePtr> char_type wide_char_type boolean_type scoped_names
%type <nodePtr> integer_type signed_int unsigned_int signed_long_int signed_short_int signed_longlong_int unsigned_long_int unsigned_short_int unsigned_longlong_int
%type <nodePtr> floating_pt_type
%type <nodePtr> sequence_type string_type wide_string_type
%type <nodePtr> const_dcl const_type enum_type enumerators
%type <nodePtr> union_type switch_type_spec switch_body case element_spec
%type <value> T_IDENTIFIER T_INTEGER_LITERAL T_STRING_LITERAL T_FLOATING_PT_LITERAL T_CHARACTER_LITERAL enumerator
%type <value> scoped_name literal T_string_literal primary_expr const_exp positive_int_const case_label

%%

/*1*/
specification
	: /*empty*/
	{

	}
	| definitions
	{
	  TheSpecification = $1
	}
	;


definitions
	: definition
	{
	defs := &Node{}
	defs.Add($1)
	$$ = defs
	}
	| definitions definition
	{
	$1.Add($2)
	$$ = $1
	}
	;

/*2*/
definition
	: type_dcl T_SEMICOLON
	| const_dcl T_SEMICOLON
	| except_dcl T_SEMICOLON { panic("Exception is not supported")}
	| interface T_SEMICOLON { panic("Interface is not supported")}
	| module T_SEMICOLON
	| value T_SEMICOLON  { panic("Value is not supported")}
	;

/*3*/
module
	: T_MODULE T_IDENTIFIER T_LEFT_CURLY_BRACKET definitions T_RIGHT_CURLY_BRACKET
    {
        mod := $4
        mod.Name = $2
        mod.Kind = KindModule
        $$ = mod
    }
	;

/*4*/
interface
	: interface_dcl
	| forward_dcl
	;

/*5*/
interface_dcl
	: interface_header T_LEFT_CURLY_BRACKET interface_body
                                         T_RIGHT_CURLY_BRACKET
	;

/*6*/
forward_dcl
	: T_INTERFACE T_IDENTIFIER
	| T_ABSTRACT T_INTERFACE T_IDENTIFIER
	;

/*7*/
interface_header
	: T_INTERFACE T_IDENTIFIER
	| T_INTERFACE T_IDENTIFIER interface_inheritance_spec
	| T_ABSTRACT T_INTERFACE T_IDENTIFIER
	| T_ABSTRACT T_INTERFACE T_IDENTIFIER interface_inheritance_spec
	;

/*8*/
interface_body
	: /*empty*/
	| exports
	;

exports
	: export
	| export exports
	;

/*9*/
export
	: type_dcl T_SEMICOLON
	| const_dcl T_SEMICOLON
	| except_dcl T_SEMICOLON
	| attr_dcl T_SEMICOLON
	| op_dcl T_SEMICOLON
	;

/*10*/
interface_inheritance_spec
	: T_COLON interface_names
	;

interface_names
	: scoped_names
	;

scoped_names
	: scoped_name
	{
	    $$ = &Node{Name: "scoped_names"}
	    $$.Add(&Node{Name: $1})
	}
	| scoped_names T_COMMA scoped_name
	{
	    $1.Add(&Node{Name: $3})
	    $$ = $1
	}
	;

/*11*/
interface_name
	: scoped_name
	;

/*12*/
scoped_name
	: T_IDENTIFIER
	{
	    $$ = $1
	}
    | T_SCOPE T_IDENTIFIER
	{
        $$ = "_" + $2
	}
	| scoped_name T_SCOPE T_IDENTIFIER
	{
        $$ = $$ + "_" + $3
	}
	;

/*13*/
value
	: value_dcl
	| value_abs_dcl
	| value_box_dcl
	| value_forward_dcl
	;

/*14*/
value_forward_dcl
	: T_VALUETYPE T_IDENTIFIER
	| T_ABSTRACT T_VALUETYPE T_IDENTIFIER
	;

/*15*/
value_box_dcl
	: T_VALUETYPE T_IDENTIFIER type_spec
	;

/*16*/
value_abs_dcl
	: T_ABSTRACT T_VALUETYPE T_IDENTIFIER
		T_LEFT_CURLY_BRACKET value_body T_RIGHT_CURLY_BRACKET
	| T_ABSTRACT T_VALUETYPE T_IDENTIFIER value_inheritance_spec
		T_LEFT_CURLY_BRACKET value_body T_RIGHT_CURLY_BRACKET
	;

value_body
	: /*empty*/
	| exports
	;

/*17*/
value_dcl
	: value_header T_LEFT_CURLY_BRACKET value_elements T_RIGHT_CURLY_BRACKET
	| value_header T_LEFT_CURLY_BRACKET T_RIGHT_CURLY_BRACKET
	;

value_elements
	: value_element
	| value_element value_elements
	;

/*18*/
value_header
	: T_VALUETYPE T_IDENTIFIER value_inheritance_spec
	| T_CUSTOM T_VALUETYPE T_IDENTIFIER value_inheritance_spec
	| T_VALUETYPE T_IDENTIFIER
	| T_CUSTOM T_VALUETYPE T_IDENTIFIER
	;

/*19*/
value_inheritance_spec
	: T_COLON value_inheritance_bases
	| T_COLON value_inheritance_bases T_SUPPORTS interface_names
	| T_SUPPORTS interface_names
	;

value_inheritance_bases
	: value_name
	| value_name T_COMMA value_names
	| T_TRUNCATABLE value_name
	| T_TRUNCATABLE value_name T_COMMA value_names
	;

value_names
	: scoped_names
	;

/*20*/
value_name
	: scoped_name
	;

/*21*/
value_element
	: export
	| state_member
	| init_dcl
	;

/*22*/
state_member
	: T_PUBLIC type_spec declarators T_SEMICOLON
	| T_PRIVATE type_spec declarators T_SEMICOLON
	;

/*23*/
init_dcl
	: T_FACTORY T_IDENTIFIER
		T_LEFT_PARANTHESIS init_param_decls T_RIGHT_PARANTHESIS
		T_SEMICOLON
	;

/*24*/
init_param_decls
	: init_param_decl
	| init_param_decl T_COMMA init_param_decls
	;

/*25*/
init_param_decl
	: init_param_attribute param_type_spec simple_declarator
	;

/*26*/
init_param_attribute
	: T_IN
	;

/*27*/
const_dcl
	: T_CONST const_type T_IDENTIFIER T_EQUAL const_exp
	{
        $$ = &Node{Name: $3, Kind: KindConst, TypeName: $2.Name, Value: $5}
	}
	;

/*28*/
const_type
	: integer_type
	| char_type
//	| wide_char_type
	| boolean_type
	| floating_pt_type
	| string_type
//	| wide_string_type
//	| fixed_pt_const_type
	| scoped_name
	{
	  $$ = &Node{Name: $1}
	}
//	| octet_type
	;

/*29*/
const_exp
//	: or_expr {}
	: primary_expr
	;

/*30*/
or_expr
	: xor_expr
	| or_expr T_VERTICAL_LINE xor_expr
	;

/*31*/
xor_expr
	: and_expr
	| xor_expr T_CIRCUMFLEX and_expr
	;

/*32*/
and_expr
	: shift_expr
	| and_expr T_AMPERSAND shift_expr
	;

/*33*/
shift_expr
	: add_expr
	| shift_expr T_SHIFTRIGHT add_expr
	| shift_expr T_SHIFTLEFT add_expr
	;

/*34*/
add_expr
	: mult_expr
	| add_expr T_PLUS_SIGN mult_expr
	| add_expr T_MINUS_SIGN mult_expr
	;

/*35*/
mult_expr
	: unary_expr
	| mult_expr T_ASTERISK unary_expr
	| mult_expr T_SOLIDUS unary_expr
	| mult_expr T_PERCENT_SIGN unary_expr
	;

/*36*/
/*37*/
unary_expr
	: T_MINUS_SIGN primary_expr
	| T_PLUS_SIGN primary_expr
	| T_TILDE primary_expr
	| primary_expr
	;

/*38*/
primary_expr
	: scoped_name
	| literal
	| T_LEFT_PARANTHESIS const_exp T_RIGHT_PARANTHESIS
	{
	    $$ = $2
	}
	;

/*39*/
/*40*/
literal
	: T_INTEGER_LITERAL
	| T_string_literal
	{
	$$ = "\"" + $1 + "\""
	}
	| T_CHARACTER_LITERAL
	| T_FIXED_PT_LITERAL {}
	| T_FLOATING_PT_LITERAL
	| T_TRUE
	{
	$$ = "true"
	}
	| T_FALSE /*boolean_literal*/
	{
	$$ = "false"
	}
	;

/*41*/
positive_int_const
	: const_exp
	;

/*42*/
/*43*/
type_dcl
	: T_TYPEDEF type_spec declarators
	{
	    panic("Typedefs are not supported.")
	}
	| struct_type
	| union_type
	| enum_type
	| T_NATIVE simple_declarator
	{
	    panic("Native is not supported.")
	}
	;

/*44*/
type_spec
	: simple_type_spec
	| constr_type_spec
	;

/*45*/
simple_type_spec
	: base_type_spec
	{
        $1.Kind = KindBaseMember
        $1.TypeName = $1.Name
        $$ = $1
	}
	| template_type_spec
	| scoped_name
	{
        $$ = &Node{Name: $1, TypeName: $1, Kind: KindMember}
	}
	;

/*46*/
base_type_spec
	: floating_pt_type
	| integer_type
	| char_type
	| wide_char_type {}
	| boolean_type
//	| octet_type {}
//	| any_type {}
	| object_type {}
//	| value_base_type {}
//	| principal_type {}
	;

/*47*/
template_type_spec
	: sequence_type
	| string_type
	{
	    $1.TypeName = $1.Name
	    $$ = $1
	}
//	| wide_string_type
//	| fixed_pt_type
	;

/*48*/
constr_type_spec
	: struct_type
	| union_type {}
	| enum_type {}
	;

/*49*/
declarators
	: declarator
	{
        decls := &Node{}
        decls.Add($1)
        $$ = decls
	}
	| declarators T_COMMA declarator
	{
        $1.Add($3)
        $$ = $1
	}
	;

/*50*/
declarator
	: simple_declarator
	| complex_declarator { panic("Complex declarators are not supported") }
	;

/*51*/
simple_declarator
	: T_IDENTIFIER
	{
	  $$ = &Node{Name: $1}
	}
	;

/*52*/
complex_declarator
	: array_declarator { panic("Array declarators are not supported") }
	;

/*53*/
floating_pt_type
	: T_FLOAT { $$ = &Node{Name: "float32"} }
	| T_DOUBLE { $$ = &Node{Name: "float64"} }
	| T_LONG T_DOUBLE  { $$ = &Node{Name: "float128"} }
	;

/*54*/
integer_type
	: signed_int
	| unsigned_int
	;

/*55*/
signed_int
	: signed_long_int
	| signed_short_int
	| signed_longlong_int
	;

/*56*/
signed_short_int
	: T_SHORT { $$ = &Node{Name: "int16"} }
	;

/*57*/
signed_long_int
	: T_LONG { $$ = &Node{Name: "int32"} }
	;

/*58*/
signed_longlong_int
	: T_LONG T_LONG { $$ = &Node{Name: "int64"} }
	;

/*59*/
unsigned_int
	: unsigned_long_int
	| unsigned_short_int
	| unsigned_longlong_int
	;

/*60*/
unsigned_short_int
	: T_UNSIGNED T_SHORT { $$ = &Node{Name: "uint16"} }
	;

/*61*/
unsigned_long_int
	: T_UNSIGNED T_LONG { $$ = &Node{Name: "uint32"} }
	;

/*62*/
unsigned_longlong_int
	: T_UNSIGNED T_LONG T_LONG { $$ = &Node{Name: "uint64"} }
	;

/*63*/
char_type
	: T_CHAR { $$ = &Node{Name: "char"} }
	;

/*64*/
wide_char_type
	: T_WCHAR { $$ = &Node{Name: "char"} }
	;

/*65*/
boolean_type
	: T_BOOLEAN { $$ = &Node{Name: "bool"} }
	;

/*66*/
octet_type
	: T_OCTET
	;

/*67*/
any_type
	: T_ANY
	;

/*68*/
object_type
	: T_OBJECT {}
	;

/*69*/
struct_type
	: T_STRUCT T_IDENTIFIER T_LEFT_CURLY_BRACKET member_list T_RIGHT_CURLY_BRACKET
    {
        $$ = $4
        $$.Name = $2
        $$.Kind = KindType
    }
	;

/*70*/
member_list
	: member
	{
	    $$ = &Node{}
	    $$.Add($1.Children()...)
	}
	| member_list member
	{
	    $$.Add($2.Children()...)
	}
	;

/*71*/
member
	: type_spec declarators T_SEMICOLON
	{
	    members := &Node{Name: "<members>"}
	    for _, v := range $2.Children() {
            members.Add(&Node{Name: v.Name, TypeName: $1.TypeName, Kind: $1.Kind, Length: $1.Length})
	    }
        $$ = members
	}
	;

/*72*/
union_type
	: T_UNION T_IDENTIFIER T_SWITCH T_LEFT_PARANTHESIS switch_type_spec T_RIGHT_PARANTHESIS T_LEFT_CURLY_BRACKET switch_body T_RIGHT_CURLY_BRACKET
	{
	    $$ = &Node{Name: $2, Kind: KindUnionType, TypeName: $5.Name}
	    for _, v := range $8.Children() {
    	    $$.Add(v)
    	}
	}
	;

/*73*/
switch_type_spec
	: integer_type
	| char_type
	| boolean_type
	| enum_type
	| scoped_name
	{
	    $$ = &Node{Name: $1}
	}
	;

/*74*/
switch_body
	: case
	{
	    $$ = &Node{Name: "<switch_body>"}
	    $$.Add($1)
	}
	| switch_body case
	{
	    $$ = $1
	    $$.Add($2)
	}
	;

/*75*/
case
	: case_label case
	{
	    panic("Multiple case labels is not supported in a union.")
	}
	| case_label element_spec T_SEMICOLON
	{
	    $$ = $2
	    $$.Value = $1
	}
	| case_label T_PRAGMA element_spec T_SEMICOLON   /* New */
	{
	    panic("case pragma is not supported in a union.")
	}
	;

/*76*/
case_label
	: T_CASE const_exp T_COLON
	{
        $$ = $2
	}
	| T_DEFAULT T_COLON
	{
	    panic("Default is not supported in a union.")
	}
	;

/*77*/
element_spec
	: type_spec declarator
	{
	    $$ = &Node{Name: $2.Name, Kind: $1.Kind, TypeName: $1.Name}
	}
	;

/*78*/
enum_type
	: T_ENUM T_IDENTIFIER T_LEFT_CURLY_BRACKET enumerators T_RIGHT_CURLY_BRACKET
	{
    	$$ = $4
    	$$.Name = $2
	}
	;

enumerators
	: enumerator
	{
        $$ = &Node{Kind: KindEnum}
        $$.Add(&Node{Name: $1})
    }
	| enumerators T_COMMA enumerator
	{
        $1.Add(&Node{Name: $3})
        $$ = $1
	}
	;

/*79*/
enumerator
	: T_IDENTIFIER
	;

/*80*/
sequence_type
	: T_SEQUENCE T_LESS_THAN_SIGN simple_type_spec T_COMMA positive_int_const T_GREATER_THAN_SIGN
    {
	    $$ = &Node{Name: $3.Name, TypeName: $3.Name, Kind: $3.Kind, Length: $5}
    }
	| T_SEQUENCE T_LESS_THAN_SIGN simple_type_spec T_GREATER_THAN_SIGN
	{
        panic("Unbounded sequence types are not supported.")
	}
	;

/*81*/
string_type
//	: T_STRING T_LESS_THAN_SIGN positive_int_const T_GREATER_THAN_SIGN
	: T_STRING T_LESS_THAN_SIGN T_INTEGER_LITERAL T_GREATER_THAN_SIGN
	{
	    $$ = &Node{Name: "string", Kind: KindBaseMember}
	    $$.Add(&Node{Name: $3})
	}
	| T_STRING
	{
	    $$ = &Node{Name: "string", Kind: KindBaseMember}
	}
	;

/*82*/
wide_string_type
	: T_WSTRING T_LESS_THAN_SIGN positive_int_const T_GREATER_THAN_SIGN
	{}
	| T_WSTRING
	{}
	;

/*83*/
array_declarator
	: T_IDENTIFIER fixed_array_sizes
	;

fixed_array_sizes
	: fixed_array_size
	| fixed_array_size fixed_array_sizes
	;

/*84*/
fixed_array_size
	: T_LEFT_SQUARE_BRACKET positive_int_const T_RIGHT_SQUARE_BRACKET
	;

/*85*/
attr_dcl
	: T_ATTRIBUTE param_type_spec simple_declarators
	| T_READONLY T_ATTRIBUTE param_type_spec simple_declarators
	;

simple_declarators
	: simple_declarator
	| simple_declarator T_COMMA simple_declarators
	;

/*86*/
except_dcl
	: T_EXCEPTION T_IDENTIFIER T_LEFT_CURLY_BRACKET members
                                          T_RIGHT_CURLY_BRACKET
                                          {}
	;

members
	: /*empty*/
	| member members
	;

/*87*/
op_dcl
	: op_attribute op_type_spec T_IDENTIFIER parameter_dcls
                                       raises_expr context_expr
	;

/*88*/
op_attribute
	: /*empty*/
	| T_ONEWAY
	;

/*89*/
op_type_spec
	: param_type_spec
	| T_VOID
	;

/*90*/
parameter_dcls
	: T_LEFT_PARANTHESIS param_dcls T_RIGHT_PARANTHESIS
	| T_LEFT_PARANTHESIS T_RIGHT_PARANTHESIS
	;

param_dcls
	: param_dcl
	| param_dcl T_COMMA param_dcls
	;

/*91*/
param_dcl
	: param_attribute param_type_spec simple_declarator
	;

/*92*/
param_attribute
	: T_IN
	| T_OUT
	| T_INOUT
	;

/*93*/
raises_expr
	: /*empty*/
	| T_RAISES T_LEFT_PARANTHESIS scoped_names T_RIGHT_PARANTHESIS
	;

/*94*/
context_expr
	: /*empty*/
	| T_CONTEXT T_LEFT_PARANTHESIS string_literals T_RIGHT_PARANTHESIS
	;

string_literals
	: T_string_literal
	| T_string_literal T_COMMA string_literals
	;

T_string_literal
	: T_STRING_LITERAL
	| T_STRING_LITERAL T_string_literal
	;

/*95*/
param_type_spec
	: base_type_spec
	| string_type
	| wide_string_type
	| scoped_name
	;

/*96*/
fixed_pt_type
	: T_FIXED T_LESS_THAN_SIGN positive_int_const T_COMMA
              T_INTEGER_LITERAL T_GREATER_THAN_SIGN
	;

/*97*/
fixed_pt_const_type
	: T_FIXED
	;

/*98*/
value_base_type
	: T_VALUEBASE {}
	;

/* New production for Principal */
principal_type
	: T_PRINCIPAL {}
	;

%%