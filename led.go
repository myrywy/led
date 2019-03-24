package led

import (
    "fmt"
    "time"
)


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


func main() {
    go emitClockSignal(21, 2000)
    watchPin(Pin(21), 11, 3, func(prev bool, current bool){fmt.Println("Transition from ", prev, " to ", current)} )
}


