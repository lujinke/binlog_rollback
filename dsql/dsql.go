package dsql

import (
	"regexp"
	"strings"
)

type TableInfo struct {
	Database string
	Table    string
}

type SqlInfo struct {
	SqlStr      string
	UseDatabase string
	Tables      []TableInfo
	sqlType     string
}

var (
	useRe         = regexp.MustCompile(`(?is)^\s*use\s+` + identPattern() + `\s*;?\s*$`)
	insertRe      = regexp.MustCompile(`(?is)^\s*insert\s+(?:ignore\s+)?(?:into\s+)?` + tablePattern())
	updateRe      = regexp.MustCompile(`(?is)^\s*update\s+` + tablePattern())
	deleteRe      = regexp.MustCompile(`(?is)^\s*delete\s+from\s+` + tablePattern())
	ddlTableRe    = regexp.MustCompile(`(?is)^\s*(?:create|alter|drop|truncate|rename)\s+(?:table\s+(?:if\s+(?:not\s+)?exists\s+)?)?` + tablePattern())
	createDBRe    = regexp.MustCompile(`(?is)^\s*create\s+database\b`)
	dropDBRe      = regexp.MustCompile(`(?is)^\s*drop\s+database\b`)
	ddlLeadingRe  = regexp.MustCompile(`(?is)^\s*(?:create|alter|drop|truncate|rename)\b`)
	spaceCollapse = regexp.MustCompile(`\s+`)
)

func identPattern() string {
	return "`?([a-zA-Z0-9_$]+)`?"
}

func tablePattern() string {
	ident := "`?([a-zA-Z0-9_$]+)`?"
	return `(?:` + ident + `\s*\.\s*)?` + ident
}

func ParseSqlsForSqlInfo(_ interface{}, sqlStr string, currentDB string) ([]*SqlInfo, string, error) {
	trimmed := strings.TrimSpace(sqlStr)
	if trimmed == "" {
		return nil, "", nil
	}

	if m := useRe.FindStringSubmatch(trimmed); len(m) > 1 {
		return nil, cleanIdent(m[1]), nil
	}

	info := &SqlInfo{SqlStr: trimmed, UseDatabase: currentDB}
	lower := strings.ToLower(spaceCollapse.ReplaceAllString(trimmed, " "))

	switch {
	case strings.HasPrefix(lower, "insert"):
		info.sqlType = "insert"
		info.Tables = extractTable(insertRe, trimmed, currentDB)
	case strings.HasPrefix(lower, "update"):
		info.sqlType = "update"
		info.Tables = extractTable(updateRe, trimmed, currentDB)
	case strings.HasPrefix(lower, "delete"):
		info.sqlType = "delete"
		info.Tables = extractTable(deleteRe, trimmed, currentDB)
	case ddlLeadingRe.MatchString(trimmed):
		info.sqlType = "ddl"
		info.Tables = extractTable(ddlTableRe, trimmed, currentDB)
	default:
		return nil, "", nil
	}

	return []*SqlInfo{info}, "", nil
}

func (s *SqlInfo) Copy() *SqlInfo {
	if s == nil {
		return nil
	}
	cp := *s
	cp.Tables = append([]TableInfo(nil), s.Tables...)
	return &cp
}

func (s *SqlInfo) IsDml() bool {
	return s != nil && (s.sqlType == "insert" || s.sqlType == "update" || s.sqlType == "delete")
}

func (s *SqlInfo) IsDdl() bool {
	return s != nil && s.sqlType == "ddl"
}

func (s *SqlInfo) IsDatabaseDDL() bool {
	if s == nil || !s.IsDdl() {
		return false
	}
	return createDBRe.MatchString(s.SqlStr) || dropDBRe.MatchString(s.SqlStr)
}

func (s *SqlInfo) GetDmlName() string {
	if s == nil || !s.IsDml() {
		return ""
	}
	return s.sqlType
}

func (s *SqlInfo) GetDatabasesAll(sep string) string {
	if s == nil {
		return ""
	}
	seen := map[string]bool{}
	var dbs []string
	for _, tb := range s.Tables {
		if tb.Database != "" && !seen[tb.Database] {
			seen[tb.Database] = true
			dbs = append(dbs, tb.Database)
		}
	}
	return strings.Join(dbs, sep)
}

func (s *SqlInfo) GetFullTablesAll(sep string) string {
	if s == nil {
		return ""
	}
	var tables []string
	for _, tb := range s.Tables {
		if tb.Database != "" {
			tables = append(tables, tb.Database+"."+tb.Table)
		} else {
			tables = append(tables, tb.Table)
		}
	}
	return strings.Join(tables, sep)
}

func extractTable(re *regexp.Regexp, sqlStr string, currentDB string) []TableInfo {
	m := re.FindStringSubmatch(sqlStr)
	if len(m) < 2 {
		return nil
	}
	db := currentDB
	table := ""
	if len(m) >= 3 && m[2] != "" {
		db = cleanIdent(m[1])
		table = cleanIdent(m[2])
	} else {
		table = cleanIdent(m[1])
	}
	if table == "" {
		return nil
	}
	return []TableInfo{{Database: db, Table: table}}
}

func cleanIdent(s string) string {
	return strings.Trim(strings.TrimSpace(s), "`")
}
