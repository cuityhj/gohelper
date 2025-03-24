package postgresql

import (
	"errors"
	"regexp"
)

var sqlStateRegexp = regexp.MustCompile(`\(SQLSTATE\s*(\S+)\)`)

func Error(err error) error {
	if err == nil {
		return nil
	}

	subValues := sqlStateRegexp.FindAllStringSubmatch(err.Error(), -1)
	if len(subValues) == 1 && len(subValues[0]) == 2 {
		if err_, ok := codeErrors[subValues[0][1]]; ok {
			return err_
		}
	}

	return err
}

var codeErrors = map[string]error{
	"00000": errors.New("successful completion (00000)"),

	"01000": errors.New("warning (01000)"),
	"0100C": errors.New("warning: dynamic result sets returned (0100C)"),
	"01008": errors.New("warning: implicit zero bit padding (01008)"),
	"01003": errors.New("warning: null value eliminated in set function (01003)"),
	"01007": errors.New("warning: privilege not granted (01007)"),
	"01006": errors.New("warning: privilege not revoked (01006)"),
	"01004": errors.New("warning: string data right truncation (01004)"),
	"01P01": errors.New("warning: deprecated feature (01P01)"),

	"02000": errors.New("no data (02000)"),
	"02001": errors.New("no data: no additional dynamic result sets returned (02001)"),

	"03000": errors.New("sql statement not yet complete (03000)"),

	"08000": errors.New("connection exception (08000)"),
	"08003": errors.New("connection exception: connection does not exist (08003)"),
	"08006": errors.New("connection exception: connection failure (08006)"),
	"08001": errors.New("connection exception: sqlclient unable to establish sqlconnection (08001)"),
	"08004": errors.New("connection exception: sqlserver rejected establishment of sqlconnection (08004)"),
	"08007": errors.New("connection exception: transaction resolution unknown (08007)"),
	"08P01": errors.New("connection exception: protocol violation (08P01)"),

	"09000": errors.New("triggered action exception (09000)"),

	"0A000": errors.New("feature not supported (0A000)"),

	"0B000": errors.New("invalid transaction initiation (0B000)"),

	"0F000": errors.New("locator exception (0F000)"),
	"0F001": errors.New("locator exception: invalid locator specification (0F001)"),

	"0L000": errors.New("invalid grantor (0L000)"),
	"0LP01": errors.New("invalid grantor: invalid grant operation (0LP01)"),

	"0P000": errors.New("invalid role specification (0P000)"),

	"0Z000": errors.New("diagnostics exception (0Z000)"),
	"0Z002": errors.New("diagnostics exception: stacked diagnostics accessed without active handler (0Z002)"),

	"20000": errors.New("case not found (20000)"),

	"21000": errors.New("cardinality violation (21000)"),

	"22000": errors.New("data exception (22000)"),
	"2202E": errors.New("data exception: array subscript error (2202E)"),
	"22021": errors.New("data exception: character not in repertoire (22021)"),
	"22008": errors.New("data exception: datetime field overflow (22008)"),
	"22012": errors.New("data exception: division by zero (22012)"),
	"22005": errors.New("data exception: error in assignment (22005)"),
	"2200B": errors.New("data exception: escape character conflict (2200B)"),
	"22022": errors.New("data exception: indicator overflow (22022)"),
	"22015": errors.New("data exception: interval field overflow (22015)"),
	"2201E": errors.New("data exception: invalid argument for logarithm (2201E)"),
	"22014": errors.New("data exception: invalid argument for ntile function (22014)"),
	"22016": errors.New("data exception: invalid argument for nth value function (22016)"),
	"2201F": errors.New("data exception: invalid argument for power function (2201F)"),
	"2201G": errors.New("data exception: invalid argument for width bucket function (2201G)"),
	"22018": errors.New("data exception: invalid character value for cast (22018)"),
	"22007": errors.New("data exception: invalid datetime format (22007)"),
	"22019": errors.New("data exception: invalid escape character (22019)"),
	"2200D": errors.New("data exception: invalid escape octet (2200D)"),
	"22025": errors.New("data exception: invalid escape sequence (22025)"),
	"22P06": errors.New("data exception: nonstandard use of escape character (22P06)"),
	"22010": errors.New("data exception: invalid indicator parameter value (22010)"),
	"22023": errors.New("data exception: invalid parameter value (22023)"),
	"22013": errors.New("data exception: invalid preceding or following size (22013)"),
	"2201B": errors.New("data exception: invalid regular expression (2201B)"),
	"2201W": errors.New("data exception: invalid row count in limit clause (2201W)"),
	"2201X": errors.New("data exception: invalid row count in result offset clause (2201X)"),
	"2202H": errors.New("data exception: invalid tablesample argument (2202H)"),
	"2202G": errors.New("data exception: invalid tablesample repeat (2202G)"),
	"22009": errors.New("data exception: invalid time zone displacement value (22009)"),
	"2200C": errors.New("data exception: invalid use of escape character (2200C)"),
	"2200G": errors.New("data exception: most specific type mismatch (2200G)"),
	"22004": errors.New("data exception: null value not allowed (22004)"),
	"22002": errors.New("data exception: null value no indicator parameter (22002)"),
	"22003": errors.New("data exception: numeric value out of range (22003)"),
	"2200H": errors.New("data exception: sequence generator limit exceeded (2200H)"),
	"22026": errors.New("data exception: string data length mismatch (22026)"),
	"22001": errors.New("data exception: string data right truncation (22001)"),
	"22011": errors.New("data exception: substring error (22011)"),
	"22027": errors.New("data exception: trim error (22027)"),
	"22024": errors.New("data exception: unterminated c string (22024)"),
	"2200F": errors.New("data exception: zero length character string (2200F)"),
	"22P01": errors.New("data exception: floating point exception (22P01)"),
	"22P02": errors.New("data exception: invalid text representation (22P02)"),
	"22P03": errors.New("data exception: invalid binary representation (22P03)"),
	"22P04": errors.New("data exception: bad copy file format (22P04)"),
	"22P05": errors.New("data exception: untranslatable character (22P05)"),
	"2200L": errors.New("data exception: not an xml document (2200L)"),
	"2200M": errors.New("data exception: invalid xml document (2200M)"),
	"2200N": errors.New("data exception: invalid xml content (2200N)"),
	"2200S": errors.New("data exception: invalid xml comment (2200S)"),
	"2200T": errors.New("data exception: invalid xml processing instruction (2200T)"),
	"22030": errors.New("data exception: duplicate json object key value (22030)"),
	"22031": errors.New("data exception: invalid argument for sql json datetime function (22031)"),
	"22032": errors.New("data exception: invalid json text (22032)"),
	"22033": errors.New("data exception: invalid sql json subscript (22033)"),
	"22034": errors.New("data exception: more than one sql json item (22034)"),
	"22035": errors.New("data exception: no sql json item (22035)"),
	"22036": errors.New("data exception: non numeric sql json item (22036)"),
	"22037": errors.New("data exception: non unique keys in a json object (22037)"),
	"22038": errors.New("data exception: singleton sql json item required (22038)"),
	"22039": errors.New("data exception: sql json array not found (22039)"),
	"2203A": errors.New("data exception: sql json member not found (2203A)"),
	"2203B": errors.New("data exception: sql json number not found (2203B)"),
	"2203C": errors.New("data exception: sql json object not found (2203C)"),
	"2203D": errors.New("data exception: too many json array elements (2203D)"),
	"2203E": errors.New("data exception: too many json object members (2203E)"),
	"2203F": errors.New("data exception: sql json scalar required (2203F)"),

	"23000": errors.New("integrity constraint violation (23000)"),
	"23001": errors.New("integrity constraint violation: restrict violation (23001)"),
	"23502": errors.New("integrity constraint violation: not null violation (23502)"),
	"23503": errors.New("integrity constraint violation: foreign key violation (23503)"),
	"23505": errors.New("integrity constraint violation: unique violation (23505)"),
	"23514": errors.New("integrity constraint violation: check violation (23514)"),
	"23P01": errors.New("integrity constraint violation: exclusion violation (23P01)"),

	"24000": errors.New("invalid cursor state (24000)"),

	"25000": errors.New("invalid transaction state (25000)"),
	"25001": errors.New("invalid transaction state: active sql transaction (25001)"),
	"25002": errors.New("invalid transaction state: branch transaction already active (25002)"),
	"25008": errors.New("invalid transaction state: held cursor requires same isolation level (25008)"),
	"25003": errors.New("invalid transaction state: inappropriate access mode for branch transaction (25003)"),
	"25004": errors.New("invalid transaction state: inappropriate isolation level for branch transaction (25004)"),
	"25005": errors.New("invalid transaction state: no active sql transaction for branch transaction (25005)"),
	"25006": errors.New("invalid transaction state: read only sql transaction (25006)"),
	"25007": errors.New("invalid transaction state: schema and data statement mixing not supported (25007)"),
	"25P01": errors.New("invalid transaction state: no active sql transaction (25P01)"),
	"25P02": errors.New("invalid transaction state: in failed sql transaction (25P02)"),
	"25P03": errors.New("invalid transaction state: idle in transaction session timeout (25P03)"),

	"26000": errors.New("invalid sql statement name (26000)"),

	"27000": errors.New("triggered data change violation (27000)"),

	"28000": errors.New("invalid authorization specification (28000)"),
	"28P01": errors.New("invalid authorization specification: invalid password (28P01)"),

	"2B000": errors.New("dependent privilege descriptors still exist (2B000)"),
	"2BP01": errors.New("dependent privilege descriptors still exist: dependent objects still exist (2BP01)"),

	"2D000": errors.New("invalid transaction termination (2D000)"),

	"2F000": errors.New("sql routine exception (2F000)"),
	"2F005": errors.New("sql routine exception: function executed no return statement (2F005)"),
	"2F002": errors.New("sql routine exception: modifying sql data not permitted (2F002)"),
	"2F003": errors.New("sql routine exception: prohibited sql statement attempted (2F003)"),
	"2F004": errors.New("sql routine exception: reading sql data not permitted (2F004)"),

	"34000": errors.New("invalid cursor name (34000)"),

	"38000": errors.New("external routine exception (38000)"),
	"38001": errors.New("external routine exception: containing sql not permitted (38001)"),
	"38002": errors.New("external routine exception: modifying sql data not permitted (38002)"),
	"38003": errors.New("external routine exception: prohibited sql statement attempted (38003)"),
	"38004": errors.New("external routine exception: reading sql data not permitted (38004)"),

	"39000": errors.New("external routine invocation exception (39000)"),
	"39001": errors.New("external routine invocation exception: invalid sqlstate returned (39001)"),
	"39004": errors.New("external routine invocation exception: null value not allowed (39004)"),
	"39P01": errors.New("external routine invocation exception: trigger protocol violated (39P01)"),
	"39P02": errors.New("external routine invocation exception: srf protocol violated (39P02)"),
	"39P03": errors.New("external routine invocation exception: event trigger protocol violated (39P03)"),

	"3B000": errors.New("savepoint exception (3B000)"),
	"3B001": errors.New("savepoint exception: invalid savepoint specification (3B001)"),

	"3D000": errors.New("invalid catalog name (3D000)"),

	"3F000": errors.New("invalid schema name (3F000)"),

	"40000": errors.New("transaction rollback (40000)"),
	"40002": errors.New("transaction rollback: transaction integrity constraint violation (40002)"),
	"40001": errors.New("transaction rollback: serialization failure (40001)"),
	"40003": errors.New("transaction rollback: statement completion unknown (40003)"),
	"40P01": errors.New("transaction rollback: deadlock detected (40P01)"),

	"42000": errors.New("syntax error or access rule violation (42000)"),
	"42601": errors.New("syntax error or access rule violation: syntax error (42601)"),
	"42501": errors.New("syntax error or access rule violation: insufficient privilege (42501)"),
	"42846": errors.New("syntax error or access rule violation: cannot coerce (42846)"),
	"42803": errors.New("syntax error or access rule violation: grouping error (42803)"),
	"42P20": errors.New("syntax error or access rule violation: windowing error (42P20)"),
	"42P19": errors.New("syntax error or access rule violation: invalid recursion (42P19)"),
	"42830": errors.New("syntax error or access rule violation: invalid foreign key (42830)"),
	"42602": errors.New("syntax error or access rule violation: invalid name (42602)"),
	"42622": errors.New("syntax error or access rule violation: name too long (42622)"),
	"42939": errors.New("syntax error or access rule violation: reserved name (42939)"),
	"42804": errors.New("syntax error or access rule violation: datatype mismatch (42804)"),
	"42P18": errors.New("syntax error or access rule violation: indeterminate datatype (42P18)"),
	"42P21": errors.New("syntax error or access rule violation: collation mismatch (42P21)"),
	"42P22": errors.New("syntax error or access rule violation: indeterminate collation (42P22)"),
	"42809": errors.New("syntax error or access rule violation: wrong object type (42809)"),
	"428C9": errors.New("syntax error or access rule violation: generated always (428C9)"),
	"42703": errors.New("syntax error or access rule violation: undefined column (42703)"),
	"42883": errors.New("syntax error or access rule violation: undefined function (42883)"),
	"42P01": errors.New("syntax error or access rule violation: undefined table (42P01)"),
	"42P02": errors.New("syntax error or access rule violation: undefined parameter (42P02)"),
	"42704": errors.New("syntax error or access rule violation: undefined object (42704)"),
	"42701": errors.New("syntax error or access rule violation: duplicate column (42701)"),
	"42P03": errors.New("syntax error or access rule violation: duplicate cursor (42P03)"),
	"42P04": errors.New("syntax error or access rule violation: duplicate database (42P04)"),
	"42723": errors.New("syntax error or access rule violation: duplicate function (42723)"),
	"42P05": errors.New("syntax error or access rule violation: duplicate prepared statement (42P05)"),
	"42P06": errors.New("syntax error or access rule violation: duplicate schema (42P06)"),
	"42P07": errors.New("syntax error or access rule violation: duplicate table (42P07)"),
	"42712": errors.New("syntax error or access rule violation: duplicate alias (42712)"),
	"42710": errors.New("syntax error or access rule violation: duplicate object (42710)"),
	"42702": errors.New("syntax error or access rule violation: ambiguous column (42702)"),
	"42725": errors.New("syntax error or access rule violation: ambiguous function (42725)"),
	"42P08": errors.New("syntax error or access rule violation: ambiguous parameter (42P08)"),
	"42P09": errors.New("syntax error or access rule violation: ambiguous alias (42P09)"),
	"42P10": errors.New("syntax error or access rule violation: invalid column reference (42P10)"),
	"42611": errors.New("syntax error or access rule violation: invalid column definition (42611)"),
	"42P11": errors.New("syntax error or access rule violation: invalid cursor definition (42P11)"),
	"42P12": errors.New("syntax error or access rule violation: invalid database definition (42P12)"),
	"42P13": errors.New("syntax error or access rule violation: invalid function definition (42P13)"),
	"42P14": errors.New("syntax error or access rule violation: invalid prepared statement definition (42P14)"),
	"42P15": errors.New("syntax error or access rule violation: invalid schema definition (42P15)"),
	"42P16": errors.New("syntax error or access rule violation: invalid table definition (42P16)"),
	"42P17": errors.New("syntax error or access rule violation: invalid object definition (42P17)"),

	"44000": errors.New("with check option violation (44000)"),

	"53000": errors.New("insufficient resources (53000)"),
	"53100": errors.New("insufficient resources: disk full (53100)"),
	"53200": errors.New("insufficient resources: out of memory (53200)"),
	"53300": errors.New("insufficient resources: too many connections (53300)"),
	"53400": errors.New("insufficient resources: configuration limit exceeded (53400)"),

	"54000": errors.New("program limit exceeded (54000)"),
	"54001": errors.New("program limit exceeded: statement too complex (54001)"),
	"54011": errors.New("program limit exceeded: too many columns (54011)"),
	"54023": errors.New("program limit exceeded: too many arguments (54023)"),

	"55000": errors.New("object not in prerequisite state (55000)"),
	"55006": errors.New("object not in prerequisite state: object in use (55006)"),
	"55P02": errors.New("object not in prerequisite state: can not change runtime param (55P02)"),
	"55P03": errors.New("object not in prerequisite state: lock not available (55P03)"),
	"55P04": errors.New("object not in prerequisite state: unsafe new enum value usage (55P04)"),

	"57000": errors.New("operator intervention (57000)"),
	"57014": errors.New("operator intervention: query canceled (57014)"),
	"57P01": errors.New("operator intervention: admin shutdown (57P01)"),
	"57P02": errors.New("operator intervention: crash shutdown (57P02)"),
	"57P03": errors.New("operator intervention: cannot connect now (57P03)"),
	"57P04": errors.New("operator intervention: database dropped (57P04)"),
	"57P05": errors.New("operator intervention: idle session timeout (57P05)"),

	"58000": errors.New("system error (58000)"),
	"58030": errors.New("system error: io error (58030)"),
	"58P01": errors.New("system error: undefined file (58P01)"),
	"58P02": errors.New("system error: duplicate file (58P02)"),

	"72000": errors.New("snapshot too old (72000)"),

	"F0000": errors.New("config file error (F0000)"),
	"F0001": errors.New("config file error: lock file exists (F0001)"),

	"HV000": errors.New("foreign data wrapper error (HV000)"),
	"HV005": errors.New("foreign data wrapper error: column name not found (HV005)"),
	"HV002": errors.New("foreign data wrapper error: dynamic parameter value needed (HV002)"),
	"HV010": errors.New("foreign data wrapper error: function sequence error (HV010)"),
	"HV021": errors.New("foreign data wrapper error: inconsistent descriptor information (HV021)"),
	"HV024": errors.New("foreign data wrapper error: invalid attribute value (HV024)"),
	"HV007": errors.New("foreign data wrapper error: invalid column name (HV007)"),
	"HV008": errors.New("foreign data wrapper error: invalid column number (HV008)"),
	"HV004": errors.New("foreign data wrapper error: invalid data type (HV004)"),
	"HV006": errors.New("foreign data wrapper error: invalid data type descriptors (HV006)"),
	"HV091": errors.New("foreign data wrapper error: invalid descriptor field identifier (HV091)"),
	"HV00B": errors.New("foreign data wrapper error: invalid handle (HV00B)"),
	"HV00C": errors.New("foreign data wrapper error: invalid option index (HV00C)"),
	"HV00D": errors.New("foreign data wrapper error: invalid option name (HV00D)"),
	"HV090": errors.New("foreign data wrapper error: invalid string length or buffer length (HV090)"),
	"HV00A": errors.New("foreign data wrapper error: invalid string format (HV00A)"),
	"HV009": errors.New("foreign data wrapper error: invalid use of null pointer (HV009)"),
	"HV014": errors.New("foreign data wrapper error: too many handles (HV014)"),
	"HV001": errors.New("foreign data wrapper error: out of memory (HV001)"),
	"HV00P": errors.New("foreign data wrapper error: no schemas (HV00P)"),
	"HV00J": errors.New("foreign data wrapper error: option name not found (HV00J)"),
	"HV00K": errors.New("foreign data wrapper error: reply handle (HV00K)"),
	"HV00Q": errors.New("foreign data wrapper error: schema not found (HV00Q)"),
	"HV00R": errors.New("foreign data wrapper error: table not found (HV00R)"),
	"HV00L": errors.New("foreign data wrapper error: unable to create execution (HV00L)"),
	"HV00M": errors.New("foreign data wrapper error: unable to create reply (HV00M)"),
	"HV00N": errors.New("foreign data wrapper error: unable to establish connection (HV00N)"),

	"P0000": errors.New("plpgsql error (P0000)"),
	"P0001": errors.New("plpgsql error: raise exception (P0001)"),
	"P0002": errors.New("plpgsql error: no data found (P0002)"),
	"P0003": errors.New("plpgsql error: too many rows (P0003)"),
	"P0004": errors.New("plpgsql error: assert failure (P0004)"),

	"XX000": errors.New("internal error (XX000)"),
	"XX001": errors.New("internal error: data corrupted (XX001)"),
	"XX002": errors.New("internal error: index corrupted (XX002)"),
}
