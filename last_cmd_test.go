package main

import (
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"os"
	"time"
)

func DescribeLastCmd(c gospec.Context) {
	c.Specify("can find the last entry in the journal", func() {
		jd, err := tmpGitRepository("journal_test")
		c.Assume(err, IsNil)

		defer func() { c.Assume(os.RemoveAll(jd), IsNil) }()

		// Create Some Entries
		for i := 0; i < 4; i++ {
			_, err := newEntry(jd, entryTmpl, func() time.Time {
				return time.Date(2000, time.January, i, 0, 0, 0, 0, time.UTC)
			}, nil, &Command{})
			c.Assume(err, IsNil)
		}

		// Create the entry weren't supposed to find
		latest, err := newEntry(jd, entryTmpl, func() time.Time {
			return time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)
		}, nil, &Command{})
		c.Assume(err, IsNil)

		filename, err := lastEntryFilename(jd)
		c.Assume(err, IsNil)

		c.Expect(filename, Equals, latest.Filename)
	})
}
