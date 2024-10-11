package operator

const (
	Unknown Operator = ""

	// Equal NotEqual LessThan GreaterThan LessThanOrEq GreaterThanOrEq
	// IsDistinct IsNotDistinct IsNull IsNotNull IsTrue IsNotTrue IsFalse IsNotFalse IsUnknown IsNotUnknown
	//
	// https://www.postgresql.org/docs/16/functions-comparison.html
	Equal           Operator = "="               // any
	NotEqual        Operator = "!="              // any
	LessThan        Operator = "<"               // any
	GreaterThan     Operator = ">"               // any
	LessThanOrEq    Operator = "<="              // any
	GreaterThanOrEq Operator = ">="              // any
	IsDistinct      Operator = "IS DISTINCT"     // bool
	IsNotDistinct   Operator = "IS NOT DISTINCT" // bool
	IsNull          Operator = "IS NULL"         // bool
	IsNotNull       Operator = "IS NOT NULL"     // bool
	IsTrue          Operator = "IS TRUE"         // bool
	IsNotTrue       Operator = "IS NOT TRUE"     // bool
	IsFalse         Operator = "IS FALSE"        // bool
	IsNotFalse      Operator = "IS NOT FALSE"    // bool
	IsUnknown       Operator = "IS UNKNOWN"      // bool
	IsNotUnknown    Operator = "IS NOT UNKNOWN"  // bool

	// StartsWith
	//
	// https://www.postgresql.org/docs/16/functions-string.html
	StartsWith Operator = "^@" // text

	// Like NotLike Regex RegexI NotRegex NotRegexI
	//
	// https://www.postgresql.org/docs/16/functions-matching.html
	Like      Operator = "LIKE"     // text
	NotLike   Operator = "NOT LIKE" // text
	Regex     Operator = "~"        // text
	RegexI    Operator = "~*"       // case-insensitive: text
	NotRegex  Operator = "!~"       // text
	NotRegexI Operator = "!~*"      // case-insensitive: text

	// Contain ContainBy Overlap
	//
	// https://www.postgresql.org/docs/16/functions-array.html
	Contain   Operator = "@>" // array, jsonb
	ContainBy Operator = "<@" // array, jsonb
	Overlap   Operator = "&&" // array, subnet

	// In NotIn
	//
	// https://www.postgresql.org/docs/16/functions-comparisons.html
	In    Operator = "IN"     // array
	NotIn Operator = "NOT IN" // array
	// expression operator ANY|ALL (array expression)

	// Exist
	// TODO JSON https://www.postgresql.org/docs/16/functions-json.html
	_Contain               = Contain
	_ContainBy             = ContainBy
	Exist         Operator = "?"  // text
	ContainKey    Operator = "?|" // array
	ContainAllKey Operator = "?&" // array
	//  @? @@ json path 不理解

	// SubnetContain SubnetContainOrEq SubnetContainBy SubnetContainByOrEq SubnetOverlap
	//
	// https://www.postgresql.org/docs/16/functions-net.html
	SubnetContain       Operator = ">>"  // subnet
	SubnetContainOrEq   Operator = ">>=" // subnet
	SubnetContainBy     Operator = "<<"  // subnet
	SubnetContainByOrEq Operator = "<<=" // subnet
	SubnetOverlap                = Overlap
)

var availableOp = map[string]Operator{
	"equal":                 Equal,
	"not_equal":             NotEqual,
	"less_than":             LessThan,
	"greater_than":          GreaterThan,
	"less_than_or_equal":    LessThanOrEq,
	"greater_than_or_equal": GreaterThanOrEq,

	"is_distinct":     IsDistinct,
	"is_not_distinct": IsNotDistinct,
	"is_null":         IsNull,
	"is_not_null":     IsNotNull,
	"is_true":         IsTrue,
	"is_not_true":     IsNotTrue,
	"is_false":        IsFalse,
	"is_not_false":    IsNotFalse,
	"is_unknown":      IsUnknown,
	"is_not_unknown":  IsNotUnknown,

	"starts_with": StartsWith,

	"like":     Like,
	"not_like": NotLike,

	"regex":       Regex,
	"not_regex":   NotRegex,
	"regex_i":     RegexI,    // case-insensitive
	"not_regex_i": NotRegexI, // case-insensitive

	"contain":    Contain,
	"contain_by": ContainBy,
	"overlap":    Overlap,

	"in":     In,
	"not_in": NotIn,

	"exist":           Exist,
	"contain_key":     ContainKey,
	"contain_all_key": ContainAllKey,

	"subnet_contain":          SubnetContain,
	"subnet_contain_or_eq":    SubnetContainOrEq,
	"subnet_contain_by":       SubnetContainBy,
	"subnet_contain_by_or_eq": SubnetContainByOrEq,
	"subnet_overlap":          SubnetOverlap,

	// abbreviation:
	"eq":      Equal,
	"ne":      NotEqual,
	"less":    LessThan,
	"greater": GreaterThan,

	"less_than_or_eq": LessThanOrEq,
	"less_or_eq":      LessThanOrEq,
	"loe":             LessThanOrEq,

	"greater_than_or_eq": GreaterThanOrEq,
	"greater_or_eq":      GreaterThanOrEq,
	"goe":                GreaterThanOrEq,
}
