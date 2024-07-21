package terminal

type TerminalType string

const (
	Air     TerminalType = "ait"
	Rail    TerminalType = "rail"
	Road    TerminalType = "road"
	Sailing TerminalType = "sailing" // port
)

type Terminal struct {
	ID      uint
	Name    string
	Type    TerminalType
	City    string
	Country string
}
