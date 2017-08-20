package transaction

import (
  "testing"
  "fmt"
  "crypto/sha256"
  "time"
  . "github.com/smartystreets/goconvey/convey"
)

var (
  warnBkp = warn
  formatBkp = format
  jsonDumpBkp = jsonDump
  readFileBkp = readFile
  writeFileBkp = writeFile
  nowBkp = now
)
func mocksReset() {
  warn = warnBkp
  format = formatBkp
  jsonDump = jsonDumpBkp
  readFile = readFileBkp
  writeFile = writeFileBkp
  now = nowBkp
}

// Build a custom transaction type:
type Transfer struct {
  From, To string
  Value int
  Coin string
}
// Implement the Op_t interface:
func (s Transfer) name() string {
  return "transfer"
}

func TestNew(t *testing.T) {
  Convey("", t, func() {
    mocksReset()

    Convey("Testing Log Manipulation", func() {
      var mockFile = []byte("content\n")
      // Stub readFile:
      readFile = func(_ string) ([]byte, error) {
        return mockFile, nil
      }

      // Stub writeFile:
      writeFile = func(_ string, content []byte) error {
        mockFile = content
        return nil
      }

      Convey("It should build correctly", func() {
        file := New("myfile")

        So(file.addr, ShouldEqual, "myfile")
        So(string(file.content), ShouldResemble, "content\n")
      })

      Convey("It should save correctly", func() {
        file := New("myfile")
        file.content = []byte("new stuff")
        file.save()

        So(string(mockFile), ShouldResemble, "new stuff")
      })

      Convey("It should reload correctly", func() {
        file := New("myfile")
        file.content = []byte("new stuff")
        file.Reload()

        So(file.addr, ShouldEqual, "myfile")
        So(string(file.content), ShouldResemble, "content\n")
      })

      Convey("It should push lines correctly", func() {
        file := New("myfile")
        file.push([]byte("new stuff"))

        So(string(file.content), ShouldResemble,
          "content\nnew stuff\n",
        )

        file.push([]byte("more stuff"))
        So(string(file.content), ShouldResemble,
          "content\nnew stuff\nmore stuff\n",
        )
      })

      Convey("It should sign commits before saving", func() {
        file := New("myfile")

        // With no arguments:
        file.Commit()
        So(string(mockFile), ShouldResemble, "content\n")

        // With arguments:
        h1 := sha256.New(); h1.Write([]byte("test1"))
        h2 := sha256.New(); h2.Write([]byte("test2"))
        file.Commit(
          Signature{
            Author: "me",
            Hash: h1,
          },
          Signature{
            Author: "him",
            Hash: h2,
          },
        )

        So(string(mockFile), ShouldResemble,
          fmt.Sprintf(
            "content\n%s:%x\n%s:%x\n",
            "me", h1.Sum(nil),
            "him", h2.Sum(nil),
          ),
        )
      })

      Convey("It should make transactions", func() {
        // Mock time function so its value won't change:
        timeNow := now()
        now = func() time.Time {
          return timeNow
        }

        file := New("myfile")

        me := "vingarcia00@gmail.com"
        him := "c.veloso.mg@gmail.com"

        // Make a transaction:
        t := Transfer{me, him, 100,"US$"}
        instance, h1 := file.Make(t)

        So(instance.Op, ShouldResemble, t.name())
        So(instance.OpInfo, ShouldResemble, t)
        So(instance.SponsorList, ShouldEqual, nil)
        So(instance.Date, ShouldEqual, timeNow)

        // The log file should now contain the transaction:
        bytes, err := jsonDump(instance)
        So(err, ShouldEqual, nil)
        So(string(file.content), ShouldResemble,
          fmt.Sprintf("content\n%s\n", bytes),
        )

        // The hash of the instance should match
        // the contents of the current log file:
        h := sha256.New()
        h.Write(file.content)
        So(h1, ShouldResemble, h)
      })
    })

    Convey("Testing transaction IO operations", func() {
    })
  })
}
