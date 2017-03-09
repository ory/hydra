package sqlparse

import (
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type SqlParseSuite struct {
}

var _ = Suite(&SqlParseSuite{})

func (s *SqlParseSuite) TestSemicolons(c *C) {
	type testData struct {
		line   string
		result bool
	}

	tests := []testData{
		{
			line:   "END;",
			result: true,
		},
		{
			line:   "END; -- comment",
			result: true,
		},
		{
			line:   "END   ; -- comment",
			result: true,
		},
		{
			line:   "END -- comment",
			result: false,
		},
		{
			line:   "END -- comment ;",
			result: false,
		},
		{
			line:   "END \" ; \" -- comment",
			result: false,
		},
	}

	for _, test := range tests {
		r := endsWithSemicolon(test.line)
		c.Assert(r, Equals, test.result)
	}
}

func (s *SqlParseSuite) TestSplitStatements(c *C) {
	type testData struct {
		sql       string
		upCount   int
		downCount int
	}

	tests := []testData{
		{
			sql:       functxt,
			upCount:   2,
			downCount: 2,
		},
		{
			sql:       multitxt,
			upCount:   2,
			downCount: 2,
		},
	}

	for _, test := range tests {
		migration, err := ParseMigration(strings.NewReader(test.sql))
		c.Assert(err, IsNil)
		c.Assert(migration.UpStatements, HasLen, test.upCount)
		c.Assert(migration.DownStatements, HasLen, test.downCount)
	}
}

func (s *SqlParseSuite) TestIntentionallyBadStatements(c *C) {
	for _, test := range intenionallyBad {
		_, err := ParseMigration(strings.NewReader(test))
		c.Assert(err, NotNil)
	}
}

func (s *SqlParseSuite) TestCustomTerminator(c *C) {
	LineSeparator = "GO"
	defer func() { LineSeparator = "" }()

	type testData struct {
		sql       string
		upCount   int
		downCount int
	}

	tests := []testData{
		{
			sql:       functxtSplitByGO,
			upCount:   2,
			downCount: 2,
		},
		{
			sql:       multitxtSplitByGO,
			upCount:   2,
			downCount: 2,
		},
	}

	for _, test := range tests {
		migration, err := ParseMigration(strings.NewReader(test.sql))
		c.Assert(err, IsNil)
		c.Assert(migration.UpStatements, HasLen, test.upCount)
		c.Assert(migration.DownStatements, HasLen, test.downCount)
	}
}

var functxt = `-- +migrate Up
CREATE TABLE IF NOT EXISTS histories (
  id                BIGSERIAL  PRIMARY KEY,
  current_value     varchar(2000) NOT NULL,
  created_at      timestamp with time zone  NOT NULL
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION histories_partition_creation( DATE, DATE )
returns void AS $$
DECLARE
  create_query text;
BEGIN
  FOR create_query IN SELECT
      'CREATE TABLE IF NOT EXISTS histories_'
      || TO_CHAR( d, 'YYYY_MM' )
      || ' ( CHECK( created_at >= timestamp '''
      || TO_CHAR( d, 'YYYY-MM-DD 00:00:00' )
      || ''' AND created_at < timestamp '''
      || TO_CHAR( d + INTERVAL '1 month', 'YYYY-MM-DD 00:00:00' )
      || ''' ) ) inherits ( histories );'
    FROM generate_series( $1, $2, '1 month' ) AS d
  LOOP
    EXECUTE create_query;
  END LOOP;  -- LOOP END
END;         -- FUNCTION END
$$
language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
drop function histories_partition_creation(DATE, DATE);
drop TABLE histories;
`

// test multiple up/down transitions in a single script
var multitxt = `-- +migrate Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE post;

-- +migrate Up
CREATE TABLE fancier_post (
    id int NOT NULL,
    title text,
    body text,
    created_on timestamp without time zone,
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE fancier_post;
`

// raise error when statements are not explicitly ended
var intenionallyBad = []string{
	// first statement missing terminator
	`-- +migrate Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
)

-- +migrate Down
DROP TABLE post;

-- +migrate Up
CREATE TABLE fancier_post (
    id int NOT NULL,
    title text,
    body text,
    created_on timestamp without time zone,
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE fancier_post;
`,

	// second half of first statement missing terminator
	`-- +migrate Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

SELECT 'No ending semicolon'

-- +migrate Down
DROP TABLE post;

-- +migrate Up
CREATE TABLE fancier_post (
    id int NOT NULL,
    title text,
    body text,
    created_on timestamp without time zone,
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE fancier_post;
`,

	// second statement missing terminator
	`-- +migrate Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE post

-- +migrate Up
CREATE TABLE fancier_post (
    id int NOT NULL,
    title text,
    body text,
    created_on timestamp without time zone,
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE fancier_post;
`,

	// trailing text after explicit StatementEnd
	`-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);
-- +migrate StatementBegin
SELECT 'no semicolon'

-- +migrate Down
DROP TABLE post;

-- +migrate Up
CREATE TABLE fancier_post (
    id int NOT NULL,
    title text,
    body text,
    created_on timestamp without time zone,
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE fancier_post;
`,
}

// Same as functxt above but split by GO lines
var functxtSplitByGO = `-- +migrate Up
CREATE TABLE IF NOT EXISTS histories (
  id                BIGSERIAL  PRIMARY KEY,
  current_value     varchar(2000) NOT NULL,
  created_at      timestamp with time zone  NOT NULL
)
GO

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION histories_partition_creation( DATE, DATE )
returns void AS $$
DECLARE
  create_query text;
BEGIN
  FOR create_query IN SELECT
      'CREATE TABLE IF NOT EXISTS histories_'
      || TO_CHAR( d, 'YYYY_MM' )
      || ' ( CHECK( created_at >= timestamp '''
      || TO_CHAR( d, 'YYYY-MM-DD 00:00:00' )
      || ''' AND created_at < timestamp '''
      || TO_CHAR( d + INTERVAL '1 month', 'YYYY-MM-DD 00:00:00' )
      || ''' ) ) inherits ( histories );'
    FROM generate_series( $1, $2, '1 month' ) AS d
  LOOP
    EXECUTE create_query;
  END LOOP;  -- LOOP END
END;         -- FUNCTION END
$$
GO
/* while GO wouldn't be used in a statement like this, I'm including it for the test */
language plpgsql
-- +migrate StatementEnd

-- +migrate Down
drop function histories_partition_creation(DATE, DATE)
GO
drop TABLE histories
GO
`

// test multiple up/down transitions in a single script, split by GO lines
var multitxtSplitByGO = `-- +migrate Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
)
GO

-- +migrate Down
DROP TABLE post
GO

-- +migrate Up
CREATE TABLE fancier_post (
    id int NOT NULL,
    title text,
    body text,
    created_on timestamp without time zone,
    PRIMARY KEY(id)
)
GO

-- +migrate Down
DROP TABLE fancier_post
GO
`
