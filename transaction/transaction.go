// package transaction
package transaction

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "hash"
  "crypto/sha256"
  "time"
)

const log_file = "transactions.json"

type File struct {
  file_addr string
  content []byte
}

func New(file_addr string) *File {
  content, _ := ioutil.ReadFile(file_addr)
  return &File{file_addr, content}
}

type Op_t interface { name() string }

type Instance struct {
  Op string
  // Defined by the specific type of Op:
  OpInfo Op_t

  // Stakeholders URLs:
  From string
  To string

  // Date and time of the operation:
  Date time.Time
}

func (log *File) load() {
  data, err := ioutil.ReadFile(log_file)
  if err != nil {
    fmt.Printf("Could not read log file!")
    panic(err)
  }

  log.content = data
}

func (log *File) save() {
  // Save it to the file:
  err := ioutil.WriteFile(log.file_addr, log.content, 0644)
  if err != nil {
    panic(err)
  }
}

func (log *File) addLine(line []byte) {
  log.content = append(append(log.content, line...), '\n')
}

func (log *File) Commit(trust_key hash.Hash) {
  // Append the key received from the counter-part:
  hash := fmt.Sprintf("%x", trust_key.Sum(nil))
  log.addLine([]byte(hash))
  // Save it to the file:
  log.save()
}

func (log *File) Reset() {
  log.load()
}

func (log *File) Make(F, T string, Op Op_t) (Instance, hash.Hash) {
  t := Instance {
    Op : Op.name(),
    OpInfo : Op,
    From : F,
    To : T,
    Date : time.Now(),
  }

  h := log.Accept(t)

  return t, h
}

func (log *File) Accept(t Instance) hash.Hash {
  // Convert the transaction to JSON:
  bytes, err := json.Marshal(t)
  if err != nil {
    panic(err)
  }

  // Save it to the buffer:
  log.addLine(bytes)

  // Generate the trust_key to be sent to the counter-part:
  h := sha256.New()
  h.Write(log.content)

  // And return it:
  return h
}

