package main

type IndexType int

const (
	IndexType_BTREE IndexType = iota
	IndexType_HASH
	IndexType_GIST
	IndexType_SPGIST
	IndexType_GIN
	IndexType_BRIN
)

// create table
type Table struct {
	TableName string
	//tableConstraints []TableConstraint
	Columns []Column
}

type Column struct {
	Name       string
	ColumnType ColumnType
	//constraint ColumnConstraint
}

type ColumnType struct {
	TypePath []string
	TypeMods []int // FIXME I believe these don't have to be just integers
}

type TableConstraint struct {
	Name       Optional[string]
	Check      Optional[CheckConstraint]
	Unique     Optional[UniqueConstraint]
	PrimaryKey Optional[PrimaryKeyConstraint]
	Exclusion  Optional[ExclusionConstraint]
	ForeignKey Optional[ForeignKeyConstraint]
}

type ColumnConstraint struct {
	Name              Optional[string]
	NotNull           bool
	Check             Optional[CheckConstraint]
	DefaultValue      Optional[DefaultConstraint]
	GeneratedValue    Optional[GeneratedValueConstraint]
	GeneratedIdentity Optional[GeneratedIdentityConstraint]
	ForeignKey        Optional[ForeignKeyConstraint]
	// collation
}

type CheckConstraint struct {
	Expression string
	NoInherit  bool
}

type DefaultConstraint struct {
	Expression string
}

type GeneratedValueConstraint struct {
	Expression string
}

type GeneratedIdentityConstraint struct {
	Always bool // false is By Default
	// TODO sequence options
}

type UniqueConstraint struct {
	ColumnNames []string // only needed for table constraints
	Nulls       bool
	NotDistinct bool
	// TODO index parameters
}

type PrimaryKeyConstraint struct {
	ColumnNames     []string // only needed for table constraints
	IndexParameters string   // TODO
}

type ForeignKeyConstraint struct {
	ColumnNames           []string // only needed for table constraints
	ReferencingTableName  string
	ReferencingColumnName string
	MatchFull             bool
	MatchPartial          bool
	MatchSimple           bool
	OnDeleteAction        string // TODO referential action
	OnUpdateAction        string // TODO referential action
}

type ExclusionConstraint struct {
	IndexType  IndexType
	Exclusions []string // seems to be a list of columnName and operator pairs
}

// create index
type Index struct {
	Name                string
	Unique              bool // only btree indexes can be set as unique
	IndexType           IndexType
	TableName           string
	ColumnNames         []string
	IncludedColumnNames []string
	// TODO ordering
	// TODO indexing expressions
	// TODO partial indexes
	// TODO operator classes
	// TODO collation
}

type Type interface {
}

// create extension
type Extension struct {
	ExtensionName string
}

// create enum type
type Enum struct {
	Name   string
	Values []string
}
