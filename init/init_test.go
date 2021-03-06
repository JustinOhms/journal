package init_test

import (
	"errors"
	"fmt"
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"github.com/ghthor/journal/entry"
	"github.com/ghthor/journal/git"
	"github.com/ghthor/journal/git/gittest"
	"github.com/ghthor/journal/idea"
	jinit "github.com/ghthor/journal/init"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func CanBeInitialized(actual, ignoredParam interface{}) (canBeInitialized bool, pos gospec.Message, neg gospec.Message, err error) {
	jd, isString := actual.(string)
	if !isString {
		err = errors.New(fmt.Sprintf("%v isn't a string", actual))
		return
	}

	canBeInitialized, _ = jinit.CanBeInitialized(jd)

	pos = gospec.Messagef(canBeInitialized, "%v can be initialized", actual)
	neg = gospec.Messagef(canBeInitialized, "%v cannot be initialized", actual)

	return
}

func HasBeenInitialized(actual, ignoredParam interface{}) (hasBeenInitialized bool, pos gospec.Message, neg gospec.Message, err error) {
	jd, isString := actual.(string)
	if !isString {
		err = errors.New(fmt.Sprintf("%v isn't a string", actual))
		return
	}

	hasBeenInitialized = jinit.HasBeenInitialized(jd)

	pos = gospec.Messagef(hasBeenInitialized, "%v has been initialized", actual)
	neg = gospec.Messagef(hasBeenInitialized, "%v has not been initialized", actual)

	return
}

func DescribeInit(c gospec.Context) {
	c.Specify("a journal", func() {
		tmpJournal := func() (directory string, commitable git.Commitable, cleanUp func()) {
			directory, err := ioutil.TempDir("", "journal_init_")
			c.Assume(err, IsNil)

			c.Assume(directory, Not(HasBeenInitialized))

			commitable, err = jinit.Journal(directory)
			c.Assume(err, IsNil)

			cleanUp = func() {
				c.Assume(os.RemoveAll(directory), IsNil)
			}

			return
		}

		c.Specify("that is initialized", func() {
			jd, commitable, cleanUp := tmpJournal()
			defer cleanUp()

			c.Assume(jd, HasBeenInitialized)

			c.Specify("is a git repository", func() {
				c.Expect(jd, gittest.IsAGitRepository)

				c.Specify("that contains", func() {
					c.Specify("an entry directory", func() {
						info, err := os.Stat(filepath.Join(jd, "entry"))
						c.Assume(err, IsNil)
						c.Expect(info.IsDir(), IsTrue)

						c.Specify("that can have entries", func() {
							c.Expect(jd, HasBeenInitialized)

							ne := entry.New(filepath.Join(jd, "entry/"))
							oe, err := ne.Open(time.Now(), nil)
							c.Assume(err, IsNil)

							_, err = oe.Close(time.Now())
							c.Assume(err, IsNil)

							c.Expect(jd, HasBeenInitialized)
						})
					})

					c.Specify("an idea directory store", func() {
						ids, err := idea.NewDirectoryStore(filepath.Join(jd, "idea/"))
						c.Assume(err, IsNil)

						c.Specify("that can have ideas", func() {
							c.Expect(jd, HasBeenInitialized)

							ids.SaveIdea(&idea.Idea{
								Name: "An Idea",
								Body: "A Body\n",
							})

							c.Expect(jd, HasBeenInitialized)
						})
					})
				})
			})

			c.Specify("has commitable changes", func() {
				c.Assume(git.IsClean(jd), Not(IsNil))

				c.Expect(git.Commit(commitable), IsNil)
				c.Expect(git.IsClean(jd), IsNil)
			})
		})

		tmpDir := func() (directory string, cleanUp func()) {
			directory, err := ioutil.TempDir("", "journal_init_")
			c.Assume(err, IsNil)

			cleanUp = func() {
				c.Assume(os.RemoveAll(directory), IsNil)
			}
			return
		}

		c.Specify("that can be initialized", func() {
			c.Specify("is an empty directory", func() {
				directory, cleanUp := tmpDir()
				defer cleanUp()

				c.Assume(directory, CanBeInitialized)

				c.Specify("inside a git repository", func() {
					jd := filepath.Join(directory, "journal")

					c.Assume(os.MkdirAll(jd, 0755), IsNil)
					c.Assume(git.Init(directory), IsNil)

					c.Expect(jd, CanBeInitialized)

					_, err := jinit.Journal(jd)
					c.Assume(err, IsNil)
					c.Expect(jd, HasBeenInitialized)
					c.Expect(jd, Not(gittest.IsAGitRepository))
				})

				c.Specify("NOT inside a git repository", func() {
					c.Expect(directory, CanBeInitialized)

					_, err := jinit.Journal(directory)
					c.Assume(err, IsNil)
					c.Expect(directory, HasBeenInitialized)
					c.Expect(directory, gittest.IsAGitRepository)
				})
			})

			c.Specify("is an empty git repository", func() {
				directory, cleanUp := tmpDir()
				defer cleanUp()

				c.Assume(git.Init(directory), IsNil)

				c.Expect(directory, CanBeInitialized)

				_, err := jinit.Journal(directory)
				c.Assume(err, IsNil)
				c.Expect(directory, HasBeenInitialized)
			})

			c.Specify("is a directory", func() {
				c.Specify("that doesn't exist", func() {
					d := filepath.Join(os.TempDir(), "doesnotexist")

					_, err := os.Stat(d)
					c.Assume(os.IsNotExist(err), IsTrue)

					c.Expect(d, CanBeInitialized)
				})

				c.Specify("that doesn't exist but is within an existing git repository", func() {
					base, cleanUp := tmpDir()
					defer cleanUp()

					c.Assume(git.Init(base), IsNil)

					jd := filepath.Join(base, "doesntexistyet")
					c.Expect(jd, CanBeInitialized)

					_, err := jinit.Journal(jd)
					c.Assume(err, IsNil)
					c.Expect(jd, HasBeenInitialized)
					c.Expect(jd, Not(gittest.IsAGitRepository))
				})
			})
		})

		c.Specify("that can NOT be initialized", func() {
			directory, cleanUp := tmpDir()
			defer cleanUp()

			c.Specify("is a file", func() {
				f, err := ioutil.TempFile(directory, "notadirectory")
				c.Assume(err, IsNil)

				c.Expect(f.Name(), Not(CanBeInitialized))

				_, err = jinit.Journal(f.Name())
				c.Assume(err, Not(IsNil))
				c.Expect(err.Error(), Equals, fmt.Sprintf("\"%s\" isn't a directory", f.Name()))
			})

			c.Specify("is a directory", func() {
				c.Specify("that already contains", func() {
					c.Specify("an entry directory", func() {
						c.Expect(directory, CanBeInitialized)

						entryDir := filepath.Join(directory, "entry")
						c.Assume(os.MkdirAll(entryDir, 0755), IsNil)

						c.Expect(directory, Not(CanBeInitialized))

						_, err := jinit.Journal(directory)
						c.Assume(err, Not(IsNil))
						c.Expect(err.Error(), Equals, fmt.Sprintf("\"%s\" isn't an empty directory", directory))
					})

					c.Specify("an idea directory store", func() {
						ideaDir := filepath.Join(directory, "idea")
						c.Assume(os.Mkdir(ideaDir, 0755), IsNil)

						_, _, err := idea.InitDirectoryStore(ideaDir)
						c.Assume(err, IsNil)

						c.Expect(directory, Not(CanBeInitialized))

						_, err = jinit.Journal(directory)
						c.Assume(err, Not(IsNil))
						c.Expect(err.Error(), Equals, fmt.Sprintf("\"%s\" isn't an empty directory", directory))
					})
				})
			})
		})
	})
}
