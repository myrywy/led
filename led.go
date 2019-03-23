package main

import (
    "fmt"
    "os"
    "strconv"
    "errors"
    "time"
)


type BinaryIo interface {
    readPinValue() (int64, error) 
    setPinValue(value int) error 
}

type Pin int

var VALUE_PATH string = "/sys/class/gpio/gpio%d/value"

func (pin_number Pin) readPinValue() (int64, error) {
    pin_value_path := fmt.Sprintf(VALUE_PATH, pin_number)
    f, err := os.OpenFile(pin_value_path, os.O_RDONLY, 0755)
    if err != nil {
        return 0, errors.New("Cannot access value file.")
    }
    defer f.Close()
    file_content := make([]byte, 1)
    bytes_read, err := f.Read(file_content)
    if err != nil || bytes_read != 1 {
        return 0, errors.New("Error reading value.")
    }
    led_state, err := strconv.ParseInt(string(file_content), 10, 8)
    if err != nil {
        return 0, errors.New("Value pseudo-file content invalid.")
    }
    return led_state, nil
}

func (pin_number Pin) setPinValue(value int) error {
    pin_value_path := fmt.Sprintf(VALUE_PATH, pin_number)
    f, err := os.OpenFile(pin_value_path, os.O_WRONLY, 0755)
    if err != nil {
        return errors.New("Cannot access value file.")
    }
    defer f.Close()
    file_content := []byte(fmt.Sprintf("%d", value))
    _, err = f.Write(file_content)
    if err != nil {
        return errors.New("Error writing value.")
    }
    return nil
}

func emitBits(pin Pin, bits chan bool, period int) {
    tick := time.Tick(time.Duration(period) * time.Millisecond)
    for bit := range bits {
        <- tick
        if bit {
            go pin.setPinValue(1)
        }else{
            go pin.setPinValue(0)
        }
    }
}

func clockSignal() chan bool {
    bit_stream := make(chan bool, 32)
    go func () {
        for {
            bit_stream <- true
            bit_stream <- false
        } 
    } ()
    return bit_stream
}

func emitClockSignal(pin Pin, period int) {
    bit_stream := clockSignal()
    emitBits(pin, bit_stream, period)
}

func watchPin(pin Pin, period int, minimal_stability int, action func(bool, bool)) {
    stable_state := false
    state_unknown := true
    deviated_from_stable_counter := 0
    tick := time.Tick(time.Duration(period) * time.Millisecond)
    for {
        <- tick
        v, err := pin.readPinValue()
        if err != nil {
            panic(err)
        }
        current_state := false
        if v > 0 {
            current_state = true
        }
        if state_unknown {
            state_unknown = false
            stable_state = current_state
        } else if !(stable_state == current_state) {
            deviated_from_stable_counter++
            if deviated_from_stable_counter > minimal_stability {
                deviated_from_stable_counter = 0
                action(stable_state, current_state)
                stable_state = current_state
            }
        }
    }
}

type I2CTransmitter struct {
    current_byte_sent uint
    bit_number uint
    scl_pin Pin
    sda_pin Pin
    waitinig_for_ack bool
}

func (transmitter I2CTransmitter) init(scl_pin Pin, sda_pin Pin) {
    transmitter.bit_number = 0
    transmitter.waitinig_for_ack = false
    transmitter.scl_pin = scl_pin
    transmitter.sda_pin = sda_pin
}

func (transmitter I2CTransmitter) clockTransitionAction(previous bool, current bool) {
    if current == false && transmitter.waitinig_for_ack == false {
        var current_bit uint
        current_bit = transmitter.current_byte_sent & (1 << (7-transmitter.bit_number))
        transmitter.bit_number++
        err := transmitter.sda_pin.setPinValue(int(current_bit))
        if err != nil {
            panic(err)
        }
        if transmitter.bit_number == 8 {
            transmitter.waitinig_for_ack = true
            transmitter.bit_number = 0
        }
    }else if current == false && transmitter.waitinig_for_ack {
        var current_bit uint
        err := transmitter.sda_pin.setPinValue(int(current_bit))
        if err != nil {
            panic(err)
        }
        transmitter.waitinig_for_ack = false
    }
}

func main() {
    go emitClockSignal(21, 2000)
    watchPin(Pin(21), 11, 3, func(prev bool, current bool){fmt.Println("Transition from ", prev, " to ", current)} )
}


