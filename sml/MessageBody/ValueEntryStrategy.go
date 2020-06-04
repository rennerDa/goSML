package MessageBody

type ValueEntryStrategy interface {
	Responsible(currentByteInArray int, message []byte) bool
	ExtractStringValue(currentByteInArray int, message []byte) string
	ExtractIntValue(currentByteInArray int, message []byte) uint64
}
