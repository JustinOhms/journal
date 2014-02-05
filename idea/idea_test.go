	"github.com/ghthor/journal/git"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
func init() {
	var err error
	if _, err = os.Stat("_test/"); os.IsNotExist(err) {
		err = os.Mkdir("_test/", 0755)
	}

	if err != nil {
		log.Fatal(err)
	}
}


	c.Specify("an idea directory", func() {
		makeEmptyDirectory := func(prefix string) string {
			d, err := ioutil.TempDir("_test", prefix+"_")
			c.Assume(err, IsNil)
			return d
		}

		makeIdeaDirectory := func(prefix string) (*IdeaDirectory, string) {
			d := makeEmptyDirectory(prefix)

			// Verify the directory isn't an IdeaDirectory
			_, err := NewIdeaDirectory(d)
			c.Assume(IsInvalidIdeaDirectoryError(err), IsTrue)

			// Initialize the directory
			id, _, err := InitIdeaDirectory(d)
			c.Assume(err, IsNil)
			c.Assume(id, Not(IsNil))

			// Verify the directory has been initialized
			id, err = NewIdeaDirectory(d)
			c.Assume(err, IsNil)
			c.Assume(id, Not(IsNil))

			return id, d
		}

		c.Specify("can be initialized", func() {
			d := makeEmptyDirectory("idea_directory_init")

			id, commitable, err := InitIdeaDirectory(d)
			c.Assume(err, IsNil)
			c.Expect(id, Not(IsNil))

			c.Expect(id.directory, Equals, d)

			c.Specify("only once", func() {
				_, _, err = InitIdeaDirectory(d)
				c.Expect(err, Equals, ErrInitOnExistingIdeaDirectory)
			})

			c.Specify("and is commitable", func() {
				c.Expect(commitable, Not(IsNil))
				c.Expect(commitable.WorkingDirectory(), Equals, d)
				c.Expect(commitable.Changes(), ContainsAll, []git.ChangedFile{
					git.ChangedFile("nextid"),
					git.ChangedFile("active"),
				})
				c.Expect(commitable.CommitMsg(), Equals, "idea directory initialized")

				c.Assume(git.Init(d), IsNil)
				c.Expect(git.Commit(commitable), IsNil)

				o, err := git.Command(d, "show", "--no-color", "--pretty=format:\"%s%b\"").Output()
				c.Assume(err, IsNil)

				c.Expect(string(o), Equals,
					`"idea directory initialized"
diff --git a/active b/active
new file mode 100644
index 0000000..e69de29
diff --git a/nextid b/nextid
new file mode 100644
index 0000000..d00491f
--- /dev/null
+++ b/nextid
@@ -0,0 +1 @@
+1
`)
			})
		})

		c.Specify("contains an index of the next available id", func() {
			_, d := makeIdeaDirectory("idea_directory_spec")

			data, err := ioutil.ReadFile(filepath.Join(d, "nextid"))
			c.Expect(err, IsNil)
			c.Expect(string(data), Equals, "1\n")
		})

		c.Specify("contains an index of active ideas", func() {
			_, d := makeIdeaDirectory("idea_directory_spec")

			_, err := os.Stat(filepath.Join(d, "active"))
			c.Expect(err, IsNil)
		})

		c.Specify("contains ideas stored in a files", func() {
			c.Specify("with the id as the filename", func() {
			})
		})

		c.Specify("can create a new idea", func() {
			c.Specify("by assigning the next available id to the idea", func() {
			})

			c.Specify("by incrementing the next available id", func() {
			})

			c.Specify("by writing the idea to a file", func() {
				c.Specify("with the id as the filename", func() {
				})

				c.Specify("and return a commitable change for the new idea file", func() {
				})
			})

			c.Specify("and if the idea's status is active", func() {
				c.Specify("will add the idea's id to the active index", func() {
					c.Specify("and will return a commitable change for modifying the index", func() {
					})
				})
			})

			c.Specify("and if the idea's status isn't active", func() {
				c.Specify("will not add the idea's id to the active index", func() {
				})
			})
		})

		c.Specify("can update an existing idea", func() {
			c.Specify("by writing the idea to the file", func() {
				c.Specify("with the id as the filename", func() {
				})

				c.Specify("and will return a commitable change for the modified idea file", func() {
				})
			})

			c.Specify("and if the idea's status is active", func() {
				c.Specify("will add the idea's id to the active index", func() {
					c.Specify("and will return a commitable change for modifying the index", func() {
					})
				})
			})

			c.Specify("and if the idea's status isn't active", func() {
				c.Specify("will not add the idea's id to the active index", func() {
				})
			})
		})
	})