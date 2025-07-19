package popx

import (
	"io"
	"io/fs"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/huanggze/x/logrusx"
	"github.com/ory/pop/v6"
)

type (
	MigrationBox struct {
		*Migrator

		Dir              fs.FS
		l                *logrusx.Logger
		migrationContent MigrationContent
		goMigrations     Migrations
	}
	MigrationContent   func(mf Migration, c *pop.Connection, r []byte, usingTemplate bool) (string, error)
	MigrationBoxOption func(*MigrationBox) *MigrationBox
)

var emptySQLReplace = regexp.MustCompile(`(?m)^(\s*--.*|\s*)$`)

func isMigrationEmpty(content string) bool {
	return len(strings.ReplaceAll(emptySQLReplace.ReplaceAllString(content, ""), "\n", "")) == 0
}

// NewMigrationBox creates a new migration box.
func NewMigrationBox(dir fs.FS, m *Migrator, opts ...MigrationBoxOption) (*MigrationBox, error) {
	mb := &MigrationBox{
		Migrator:         m,
		Dir:              dir,
		l:                m.l,
		migrationContent: ParameterizedMigrationContent(nil),
	}

	for _, o := range opts {
		mb = o(mb)
	}

	txRunner := func(b []byte) func(Migration, *pop.Connection, *pop.Tx) error {
		return func(mf Migration, c *pop.Connection, tx *pop.Tx) error {
			content, err := mb.migrationContent(mf, c, b, true)
			if err != nil {
				return errors.Wrapf(err, "error processing %s", mf.Path)
			}
			if isMigrationEmpty(content) {
				m.l.WithField("migration", mf.Path).Trace("This is usually ok - ignoring migration because content is empty. This is ok!")
				return nil
			}
			if _, err = tx.Exec(content); err != nil {
				return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
			}
			return nil
		}
	}

	autoCommitRunner := func(b []byte) func(Migration, *pop.Connection) error {
		return func(mf Migration, c *pop.Connection) error {
			content, err := mb.migrationContent(mf, c, b, true)
			if err != nil {
				return errors.Wrapf(err, "error processing %s", mf.Path)
			}
			if isMigrationEmpty(content) {
				m.l.WithField("migration", mf.Path).Trace("This is usually ok - ignoring migration because content is empty. This is ok!")
				return nil
			}
			if _, err = c.RawQuery(content).ExecWithCount(); err != nil {
				return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
			}
			return nil
		}
	}

	err := mb.findMigrations(txRunner, autoCommitRunner)
	if err != nil {
		return mb, err
	}

	for _, migration := range mb.goMigrations {
		mb.Migrations[migration.Direction] = append(mb.Migrations[migration.Direction], migration)
	}

	if err := mb.check(); err != nil {
		return nil, err
	}
	return mb, nil
}

func (fm *MigrationBox) findMigrations(
	runner func([]byte) func(mf Migration, c *pop.Connection, tx *pop.Tx) error,
	runnerNoTx func([]byte) func(mf Migration, c *pop.Connection) error,
) error {
	return fs.WalkDir(fm.Dir, ".", func(p string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		if info.IsDir() {
			return nil
		}

		match, err := ParseMigrationFilename(info.Name())
		if err != nil {
			if strings.HasPrefix(err.Error(), "unsupported dialect") {
				fm.l.Tracef("This is usually ok - ignoring migration file %s because dialect is not supported: %s", info.Name(), err.Error())
				return nil
			}
			return errors.WithStack(err)
		}

		if match == nil {
			fm.l.Tracef("This is usually ok - ignoring migration file %s because it does not match the file pattern.", info.Name())
			return nil
		}

		f, err := fm.Dir.Open(p)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return errors.WithStack(err)
		}

		mf := Migration{
			Path:       p,
			Version:    match.Version,
			Name:       match.Name,
			DBType:     match.DBType,
			Direction:  match.Direction,
			Type:       match.Type,
			Content:    string(content),
			Autocommit: match.Autocommit,
		}

		if match.Autocommit {
			mf.RunnerNoTx = runnerNoTx(content)
		} else {
			mf.Runner = runner(content)
		}

		fm.Migrations[mf.Direction] = append(fm.Migrations[mf.Direction], mf)
		mod := sort.Interface(fm.Migrations[mf.Direction])
		if mf.Direction == "down" {
			mod = sort.Reverse(mod)
		}
		sort.Sort(mod)
		return nil
	})
}

// hasDownMigrationWithVersion checks if there is a migration with the given
// version.
func (fm *MigrationBox) hasDownMigrationWithVersion(version string) bool {
	for _, down := range fm.Migrations["down"] {
		if version == down.Version {
			return true
		}
	}
	return false
}

// check checks that every "up" migration has a corresponding "down" migration.
func (fm *MigrationBox) check() error {
	for _, up := range fm.Migrations["up"] {
		if !fm.hasDownMigrationWithVersion(up.Version) {
			return errors.Errorf("migration %s has no corresponding down migration", up.Version)
		}
	}

	for _, m := range fm.Migrations {
		for _, n := range m {
			if err := n.Valid(); err != nil {
				return err
			}
		}
	}
	return nil
}
