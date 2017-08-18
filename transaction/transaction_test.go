package transaction

import (
  "testing"
  "fmt"
  "os"
  "io/ioutil"
  // "crypto/sha256"
)

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

func TestTransaction(t *testing.T) {
  // Initialize 2 logs with different data:
  ioutil.WriteFile("test1.log", []byte("seed1\n"), 0644)
  ioutil.WriteFile("test2.log", []byte("seed2\n"), 0644)

  // Instantiate "transaction.File" for each:
  log1 := New("test1.log")
  log2 := New("test2.log")
  defer os.Remove("test1.log")
  defer os.Remove("test2.log")

  me := "vingarcia00@gmail.com"
  him := "c.veloso.mg@gmail.com"

  // Make a transaction:
  t1, h1 := log1.Make(Transfer{me, him, 100,"US$"})
  fmt.Printf("\n%v\n", t1)
  fmt.Printf("%x\n", h1.Sum(nil))

  h2 := log2.Accept(t1)
  fmt.Printf("%x\n\n", h2.Sum(nil))

  // Commit the transactions by exchanging trust_keys:
  log1.Commit(Signature{him, h2})
  log2.Commit(Signature{me, h1})

  // Now read both files to check for correctess:
  f1, err := ioutil.ReadFile("test1.log")
  if err != nil {
    t.Error(err)
  }
  f2, err := ioutil.ReadFile("test2.log")
  if err != nil {
    t.Error(err)
  }

  i:=0
  // Consume the seed# line:
  for ; f1[i] != '\n'; i++ {}
  i++
  // Make sure the second line is exactly the same:
  for ; f1[i] != '\n'; i++ {
    if f1[i] != f2[i] {
      t.Error("The transactions are not equal!")
    }
  }
  i++
  // The third line must correspond to the hash code:
  h1_s := fmt.Sprintf("%s:%x", me, h1.Sum(nil))
  h2_s := fmt.Sprintf("%s:%x", him, h2.Sum(nil))
  offset := i
  for i=0; f1[i] != '\n'; i++ {
    if f1[i+offset] != h2_s[i] {
      t.Error("Hash incorrect for file 1!")
    }
    if f2[i+offset] != h1_s[i] {
      t.Error("Hash incorrect for file 2!")
    }
  }
}
