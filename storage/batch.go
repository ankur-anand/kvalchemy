package storage

import (
	"encoding/binary"
	"hash/crc32"

	"github.com/ankur-anand/kvalchemy/storage/wrecord"
	"github.com/rosedblabs/wal"
	"github.com/segmentio/ksuid"
)

// Batch can be used to save a batch of key and values.
// Batch must be commited for write to be visible.
// UnCommited Batch will never be visible to reader until commited.
type Batch struct {
	key      []byte
	batchID  []byte // a private Batch ID
	lastPos  *wal.ChunkPosition
	err      error
	commited bool
	engine   *Engine
	checksum uint32 // Rolling checksum
}

// NewBatch return an initialized batch, with a start marker in wal.
func (e *Engine) NewBatch(key []byte) (*Batch, error) {
	uuid, err := ksuid.New().MarshalBinary()
	if err != nil {
		return nil, err
	}

	index := e.globalCounter.Add(1)

	// start the batch marker in wal
	record := walRecord{
		hlc:          index,
		key:          key,
		value:        nil,
		op:           wrecord.LogOperationOpBatchStart,
		batchID:      uuid,
		lastBatchPos: nil,
	}

	encoded, err := record.fbEncode()
	if err != nil {
		return nil, err
	}

	chunkPos, err := e.wal.Write(encoded)
	if err != nil {
		return nil, err
	}

	return &Batch{
		key:     key,
		batchID: uuid,
		lastPos: chunkPos,
		err:     err,
		engine:  e,
	}, nil
}

func (b *Batch) Key() []byte {
	return b.key
}

func (b *Batch) BatchID() []byte {
	return b.batchID
}

func (b *Batch) LastPos() wal.ChunkPosition {
	return *b.lastPos
}

// Put the value for the key.
func (b *Batch) Put(value []byte) error {
	if b.err != nil {
		return b.err
	}

	index := b.engine.globalCounter.Add(1)

	record := &walRecord{
		hlc:          index,
		key:          b.key,
		value:        value,
		op:           wrecord.LogOperationOPBatchInsert,
		batchID:      b.batchID,
		lastBatchPos: b.lastPos,
	}

	// Encode and compress WAL record
	encoded, err := record.fbEncode()

	if err != nil {
		b.err = err
		return err
	}

	// Write to WAL
	chunkPos, err := b.engine.wal.Write(encoded)

	if err != nil {
		b.err = err
		return err
	}

	b.lastPos = chunkPos
	// Update the rolling checksum

	b.checksum = crc32.Update(b.checksum, crc32.IEEETable, value)
	return nil
}

// Commit the given Batch to wal.
func (b *Batch) Commit() error {
	if b.err != nil {
		return b.err
	}

	return b.engine.persistKeyValue(b.key, marshalChecksum(b.checksum), wrecord.LogOperationOpBatchCommit, b.batchID, b.lastPos)
}

func marshalChecksum(checksum uint32) []byte {
	buf := make([]byte, 4) // uint32 takes 4 bytes
	binary.LittleEndian.PutUint32(buf, checksum)
	return buf
}

func unmarshalChecksum(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(data)
}
