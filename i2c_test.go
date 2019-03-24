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

func checkByteEmitted(t *testing.T, transmitter *I2CTransmitter, sda_pin *MockPin, byte_value uint) {
	sda_pin.set_calls = make([]int,0)
	transmitter.current_byte_sent = byte_value
	for i := 0; i < 8; i++ {
		transmitter.clockTransitionAction(false, true)
		if len(sda_pin.set_calls) != i {
			t.Errorf("Set SDA state in high SCL state.")
		}
		transmitter.clockTransitionAction(true, false)
		if len(sda_pin.set_calls) != i+1 {
			t.Errorf("Didn't set SDA state for bit no %d", i)
		}
		if transmitter.waiting_for_ack && i != 7{
			t.Errorf("Waiting for ack before a byte is sent.")
		}
	}
	sent_value := uint(
			128 * sda_pin.set_calls[0] + 
			64 * sda_pin.set_calls[1] + 
			32 * sda_pin.set_calls[2] + 
			16 * sda_pin.set_calls[3] + 
			8 * sda_pin.set_calls[4] + 
			4 * sda_pin.set_calls[5] + 
			2 * sda_pin.set_calls[6] + 
			sda_pin.set_calls[7])

	if sent_value != byte_value {
		t.Errorf("Wrong byte sent. Should be %d, was %d", byte_value, sent_value)
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

func TestClockTransitionAction(t *testing.T) {
	var transmitter I2CTransmitter
	sda_pin := MockPin{0, make([]int,0)}
	scl_pin := MockPin{0, make([]int,0)}
	transmitter.init(&scl_pin, &sda_pin)
	checkByteEmitted(t, &transmitter, &sda_pin, 255)
	checkByteEmitted(t, &transmitter, &sda_pin, 128)
	checkByteEmitted(t, &transmitter, &sda_pin, 64)
	checkByteEmitted(t, &transmitter, &sda_pin, 32)
	checkByteEmitted(t, &transmitter, &sda_pin, 16)
	checkByteEmitted(t, &transmitter, &sda_pin, 8)
	checkByteEmitted(t, &transmitter, &sda_pin, 4)
	checkByteEmitted(t, &transmitter, &sda_pin, 2)
	checkByteEmitted(t, &transmitter, &sda_pin, 0)
	checkByteEmitted(t, &transmitter, &sda_pin, 43)
	checkByteEmitted(t, &transmitter, &sda_pin, 23)
	checkByteEmitted(t, &transmitter, &sda_pin, 87)
	checkByteEmitted(t, &transmitter, &sda_pin, 111)
}