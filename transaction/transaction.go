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

/**
 * Transaction File:
 *
 * This log records every transaction and
 * associated hash keys.
 *
 * It is loaded as a blob for 2 reasons:
 *
 * 1. Data should only be appended at the
 *    end of the file
 *
 * 2. The old data should only be used to
 *    produce the verification sha256 hash keys,
 *    for which a `[]byte` is enough
 */
type File struct {
  file_addr string
  content []byte
}

func New(file_addr string) *File {
  content, _ := ioutil.ReadFile(file_addr)
  return &File{file_addr, content}
}

/**
 * Operation Interface:
 *
 * Any new transaction model
 * must fulfill 2 requirementes:
 *
 * 1. Be convertable to a JSON string
 * 2. Implement this interface
 */
type Op_t interface {
    name() string
}

/**
 * Transaction Instance:
 *
 * Represents a single transaction
 */
type Instance struct {
  // The operation unique name:
  Op string

  // Contains all specific information that
  // will be recorded by this operation:
  OpInfo Op_t

  // The URLs of those responsible to attest
  // the validity of the transaction:
  SponsorList []string

  // Date and time of the transaction:
  Date time.Time
}

/*
 * Transaction Signature:
 *
 * A hash signature to document the Author's
 * current version of the transaction file
 * and send to someone else.
 */
type Signature struct {
  Author string
  Hash hash.Hash
}

/**
 * Writes current file blob to disk.
 */
func (log *File) save() {
  // Save it to the file:
  err := ioutil.WriteFile(log.file_addr, log.content, 0644)
  if err != nil {
    panic(err)
  }
}

/**
 * Discards any changes made on current
 * file's content, and reload it from disk.
 */
func (log *File) Reload() {
  data, err := ioutil.ReadFile(log.file_addr)
  if err != nil {
    fmt.Printf("Could not read log file!")
    panic(err)
  }

  log.content = data
}


/**
 * Push:
 *
 * Appends data to the end of the content blob
 */
func (log *File) push(line []byte) {
  log.content = append(append(log.content, line...), '\n')
}

/**
 * Commit:
 *
 * Saves the current content blob state followed
 * by a set of signatures to attest its validity.
 */
func (log *File) Commit(sig ...Signature) {
  // Append the key received from the counter-part:
  for _, item := range sig {
    hash := fmt.Sprintf("%s:%x", item.Author, item.Hash.Sum(nil))
    log.push([]byte(hash))
  }

  // Save it to the file:
  log.save()
}

/**
 * Make Transaction:
 *
 * Receives an operation add a timestamp and build
 * a transaction.
 */
func (log *File) Make(Op Op_t) (Instance, hash.Hash) {
  t := Instance {
    Op : Op.name(),
    OpInfo : Op,
    Date : time.Now(),
  }

  h := log.Accept(t)

  return t, h
}

/**
 * Accept Transaction:
 *
 * Accepts a Transaction Instance produced by some other user,
 * and returns a sha1 that can be used as a signature, for
 * validating the new transaction.
 */
func (log *File) Accept(t Instance) hash.Hash {
  // Convert the transaction to JSON:
  bytes, err := json.Marshal(t)
  if err != nil {
    panic(err)
  }

  // Save it to the buffer:
  log.push(bytes)

  // Generate the trust_key to be sent to the counter-part:
  h := sha256.New()
  h.Write(log.content)

  // And return it:
  return h
}

