package led 

type I2CTransmitter struct {
    current_byte_sent uint
    bit_number uint
    scl_pin BinaryIo
    sda_pin BinaryIo
    waiting_for_ack bool
}

func (transmitter *I2CTransmitter) init(scl_pin BinaryIo, sda_pin BinaryIo) {
    transmitter.bit_number = 0
    transmitter.waiting_for_ack = false
    transmitter.scl_pin = scl_pin
    transmitter.sda_pin = sda_pin
}

func (transmitter *I2CTransmitter) clockTransitionAction(previous bool, current bool) {
    if current == false && transmitter.waiting_for_ack == false {
        var current_bit uint
        current_bit = (transmitter.current_byte_sent >> (7-transmitter.bit_number)) & 1
        transmitter.bit_number++
        err := transmitter.sda_pin.setPinValue(int(current_bit))
        if err != nil {
            panic(err)
        }
        if transmitter.bit_number == 8 {
            transmitter.waiting_for_ack = true
            transmitter.bit_number = 0
        }
    }else if current == false && transmitter.waiting_for_ack {
        err := transmitter.sda_pin.setPinValue(1)
        if err != nil {
            panic(err)
        }
        transmitter.waiting_for_ack = false
    }
}