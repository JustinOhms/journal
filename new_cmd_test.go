	"bufio"
	"bytes"
	"time"
		const fileData = `
#~ git commit msg
# Subject/Title
Some data that shouldn't be modified
`

		tmpl, err := template.New(jd).Parse(fileData)
		c.Assume(err, IsNil)

			j, err := newEntry(jd, tmpl, nil, nil, &Command{})
		c.Specify("will commit the entry to the git repository", func() {
			j, err := newEntry(jd, entryTmpl, func() time.Time {
				return time.Time{}
			}, nil, &Command{})
			c.Assume(err, IsNil)

			o, err := GitCommand(jd, "show", "--no-color", "--pretty=format:\"%s%b\"").Output()
			c.Assume(err, IsNil)

			actualData := bytes.NewBuffer(o)
			expectedData := bytes.NewBuffer(make([]byte, 0, 1024))

			commitMsgTmpl, err := template.New("commitMsgTmpl").Parse(
				`"Title(will be used as commit message)"
diff --git a/{{.Filename}} b/{{.Filename}}
new file mode 100644
index 0000000..df0111e
--- /dev/null
+++ b/{{.Filename}}
@@ -0,0 +1,9 @@
+{{.OpenedAt}}
+
+#~ Title(will be used as commit message)
+TODO Make this some random quote or something stupid
+
+## [active] An Idea
+An idea carries over from entry to entry if it is active.
+
+{{.ClosedAt}}
`)
			c.Assume(err, IsNil)
			c.Assume(commitMsgTmpl.Execute(expectedData, j), IsNil)

			c.Expect(actualData.String(), Equals, expectedData.String())

			// Helps with debugging the test
			// Outputs the first line that doesn't match
			actualDataSc, expectedDataSc := bufio.NewScanner(actualData), bufio.NewScanner(expectedData)
			for actualDataSc.Scan() && expectedDataSc.Scan() {
				c.Expect(actualDataSc.Text(), Equals, expectedDataSc.Text())
				if actualDataSc.Text() != expectedDataSc.Text() {
					break
				}
			}
		})

				_, err := newEntry(jd, entryTmpl, nil, nil, &Command{})