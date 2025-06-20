// modules/flake/sonyflake.go
package flake

import (
	"os"
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

func NewSonyflake() *sonyflake.Sonyflake {
    startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

    machineIDFunc := func() (uint16, error) {
        raw := os.Getenv("SONYFLAKE_MACHINE_ID")
        if raw == "" {
            return 1, nil
        }
        id, err := strconv.ParseUint(raw, 10, 16)
        return uint16(id), err
    }

    settings := sonyflake.Settings{
        StartTime: startTime,
        MachineID: machineIDFunc,
    }

    fl := sonyflake.NewSonyflake(settings)
    if fl == nil {
        panic("failed to initialize Sonyflake")
    }
    return fl
}
