package led

import (
	"testing"
)

type MockPin struct {
	read_calls int
	set_calls []int
}

func (p *MockPin) readPinValue() (int64, error) {
	p.read_calls++
	return 0, nil
}

func (p *MockPin) setPinValue(value int) error {
	p.set_calls = append(p.set_calls, value)
	return nil
}

func TestClockTransitionAction(t *testing.T) {
	var transmitter I2CTransmitter
	sda_pin := MockPin{0, make([]int,0)}
	scl_pin := MockPin{0, make([]int,0)}
	transmitter.init(&scl_pin, &sda_pin)
	transmitter.clockTransitionAction(false, true)
	if len(sda_pin.set_calls) != 0 {
		t.Errorf("Shouldn't set after scl transition to 1, sets history: ", sda_pin.set_calls)
	}
	transmitter.current_byte_sent=128
	transmitter.clockTransitionAction(true, false)
	if len(sda_pin.set_calls) != 1 || sda_pin.set_calls[0] != 1 {
		t.Errorf("Bit no %d doesn't match", 0, sda_pin.set_calls)
	}
	if transmitter.waiting_for_ack {
		t.Errorf("Waiting for ack before a byte is sent.")
	}
	for i := 1; i < 8; i++ {
		transmitter.clockTransitionAction(true, false)
		if len(sda_pin.set_calls) != i+1 || sda_pin.set_calls[i] != 0 {
			t.Errorf("Bit no %d doesn't match: ", i, sda_pin.set_calls)
		}
		if transmitter.waiting_for_ack && i != 7{
			t.Errorf("Waiting for ack before a byte is sent.")
		}
	}
	if !transmitter.waiting_for_ack {
		t.Errorf("Should be waititng for acknowledgement after the last bit but it's not.")
	}
	transmitter.clockTransitionAction(true, false)
	if len(sda_pin.set_calls) != 9 || sda_pin.set_calls[8] != 1 {
		t.Errorf("Didn't set 1 state in SDA when waiting for ack.")
	}
	if transmitter.bit_number != 0 || transmitter.waiting_for_ack {
		t.Errorf("Transmitter didn't return to the initial stete after sending a byte")
	}
}