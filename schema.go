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
	tableName        string
	tableConstraints []TableConstraint
	columns          []Column
}

type Column struct {
	name       string
	columnType ColumnType
	constraint ColumnConstraint
}

type ColumnType struct {
	baseType string
}

type TableConstraint struct {
	name       Optional[string]
	check      Optional[CheckConstraint]
	unique     Optional[UniqueConstraint]
	primaryKey Optional[PrimaryKeyConstraint]
	exclusion  Optional[ExclusionConstraint]
	foreignKey Optional[ForeignKeyConstraint]
}

type ColumnConstraint struct {
	name              Optional[string]
	notNull           bool
	check             Optional[CheckConstraint]
	defaultValue      Optional[DefaultConstraint]
	generatedValue    Optional[GeneratedValueConstraint]
	generatedIdentity Optional[GeneratedIdentityConstraint]
	foreignKey        Optional[ForeignKeyConstraint]
	// collation
}

type CheckConstraint struct {
	expression string
	noInherit  bool
}

type DefaultConstraint struct {
	expression string
}

type GeneratedValueConstraint struct {
	expression string
}

type GeneratedIdentityConstraint struct {
	always bool // false is By Default
	// TODO sequence options
}

type UniqueConstraint struct {
	columnNames []string // only needed for table constraints
	nulls       bool
	notDistinct bool
	// TODO index parameters
}

type PrimaryKeyConstraint struct {
	columnNames     []string // only needed for table constraints
	indexParameters string   // TODO
}

type ForeignKeyConstraint struct {
	columnNames           []string // only needed for table constraints
	referencingTableName  string
	referencingColumnName string
	matchFull             bool
	matchPartial          bool
	matchSimple           bool
	onDeleteAction        string // TODO referential action
	onUpdateAction        string // TODO referential action
}

type ExclusionConstraint struct {
	indexType  IndexType
	exclusions []string // seems to be a list of columnName and operator pairs
}

// create index
type Index struct {
	name                string
	unique              bool // only btree indexes can be set as unique
	indexType           IndexType
	tableName           string
	columnNames         []string
	includedColumnNames []string
	// TODO ordering
	// TODO indexing expressions
	// TODO partial indexes
	// TODO operator classes
	// TODO collation
}

// create extension
type Extension struct {
	extensionName string
}
