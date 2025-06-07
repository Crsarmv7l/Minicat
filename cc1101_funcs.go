package main

import (
	"machine"
	"strconv"
	"time"
    	"math"
	"tinygo.org/x/drivers/delay"
)


type radio struct {
	spi        *machine.SPI
	cs          machine.Pin
	GDO0        machine.Pin
	GDO2        machine.Pin
	frequency   float32
	modulation  int8
	bitRate     float32
	freqdev     float32
	rxbw        float32
	pwr	    int8
}

func newradio(spi *machine.SPI, cs machine.Pin, gdo0 machine.Pin, gdo2 machine.Pin) *radio {
	spi.Configure(machine.SPIConfig{
		Frequency: 4000000,
		LSBFirst:  false,
		Mode:      1,
	})
	
	cs.Configure(machine.PinConfig{Mode: machine.PinOutput})
	cs.High()
	
	return &radio {
	spi:         spi,
	cs:          cs,
	GDO0:        gdo0,
	GDO2:        gdo2,
	frequency:   CC1101_DEFAULT_FREQ,
	modulation:  CC1101_MOD_FORMAT_2_FSK,
	bitRate:     CC1101_DEFAULT_BR,
	freqdev:     CC1101_DEFAULT_FREQDEV,
	rxbw:        CC1101_DEFAULT_RXBW,
	pwr:	     CC1101_DEFAULT_POWER,
	}
}


func (r *radio) ReadReg(reg uint8) uint8 {
	if reg > CC1101_REG_TEST0 && reg < CC1101_REG_PATABLE {
    		reg |= CC1101_CMD_READ_BURST
   	}else {
   		reg |= CC1101_CMD_READ_SINGLE
   	}
	tx := []byte{reg, 0x00}
	rx := make([]byte, 2)
	r.cs.Low()
	r.spi.Tx(tx, rx)
	r.cs.High()
	return rx[1]
}

func (r *radio) WriteReg(reg uint8, value uint8) {
	tx := []byte{reg, value}
	r.cs.Low()
	r.spi.Tx(tx, nil)
	r.cs.High()
}

func (r *radio) WriteRegBurst(reg uint8, data []byte) {
	reg |= CC1101_CMD_WRITE_BURST
	r.cs.Low()
	r.spi.Tx([]byte{reg}, nil)
	r.spi.Tx(data, nil)
	r.cs.High()
}

func (r *radio) Strobe(reg uint8) {
	r.cs.Low()
	r.spi.Tx([]byte{reg}, nil)
	r.cs.High()
}

func (r *radio) setFreq(freq float32) {
	if !(freq >= 300.0 && freq <= 348.0 ||
	    freq >= 387.0 && freq <= 464.0 ||
	    freq >= 779.0 && freq <= 928.0) {
		println("Frequency out of Params")
		return
	}
	r.frequency = freq
	r.standby()
	
	var FRF uint32= uint32((freq * 65536) / 26.0)
	r.WriteReg(CC1101_REG_FREQ2, uint8((FRF & 0xFF0000) >> 16))
	r.WriteReg(CC1101_REG_FREQ1, uint8((FRF & 0x00FF00) >> 8))
	r.WriteReg(CC1101_REG_FREQ0, uint8(FRF & 0x0000FF))
	
	r.setOutputPower(r.pwr)
}

func (r *radio) setOOK(enableOOK bool) {
	r.standby()
	if enableOOK {
		r.setRegValue(CC1101_REG_MDMCFG2, CC1101_MOD_FORMAT_ASK_OOK, 6, 4)
		r.setRegValue(CC1101_REG_FREND0, 1, 2, 0)
		//Per datasheet
		r.setRegValue(CC1101_REG_FOCCFG, 0, 1, 0)
       		r.modulation = CC1101_MOD_FORMAT_ASK_OOK
	} else {
		r.setRegValue(CC1101_REG_MDMCFG2, CC1101_MOD_FORMAT_2_FSK, 6, 4)
		r.setRegValue(CC1101_REG_FREND0, 0, 2, 0)
		r.setRegValue(CC1101_REG_FOCCFG, 2, 1, 0)
		r.modulation = CC1101_MOD_FORMAT_2_FSK
	}
	
	r.setOutputPower(r.pwr)
}

func (r *radio) setBitrate(br float32) {
	if br < 0.025 || br > 600.00 {
		println("Bitrate out of Params")
		return
	}
	
	r.standby()
	
	e, m := getExpMant((br * 1000.0), 256, 28, 14)
	if e == 0 && m == 0 {
		println("Error Calculating Bitrate Exponent and Mantissa")
		return
	}
	
	r.setRegValue(CC1101_REG_MDMCFG4, e, 3, 0)
	r.WriteReg(CC1101_REG_MDMCFG3, m)
	r.bitRate = br
}

func getExpMant(target float32, mantOffset uint16, divExp uint8, expMax uint8) (uint8, uint8) {
	var e uint8 = 0
	var origin float32 = float32((float32(mantOffset) * CC1101_CRYSTAL_FREQ * 1000000.0)/float32((uint32(1) << divExp)))
	for e =expMax; e >= 0; e -- {
  		var intervalStart float32 = float32((uint32(1) << e)) * origin
  		if target >= intervalStart {
  			stepSize  := float32(intervalStart/float32(mantOffset))
  			m  := uint8((target - intervalStart)/stepSize)
  			return e, m
  		}
  	}
  	return 0, 0
}

func (r *radio) setRegValue(reg uint8, value uint8, msb uint8, lsb uint8) {
	if((msb > 7) || (lsb > 7) || (lsb > msb)) {
    		println("Invalid MSB/LSB mask")
    		return
  	}
  	currentValue := r.ReadReg(reg)
  	mask := uint8 (^((0xFF << (msb + 1)) | (0xFF >> (8 - lsb))))
  	if !((currentValue & mask) == (value & mask)) {
    		newValue := (currentValue & ^mask) | (value & mask)
  		r.WriteReg(reg, newValue)
  	}
}

func (r *radio) getRegValue(reg uint8, msb uint8, lsb uint8) uint8 {
  	if((msb > 7) || (lsb > 7) || (lsb > msb)) {
    		println("Invalid MSB/LSB mask")
    		return(0)
  	}
  	rawVal := r.ReadReg(reg)
  	maskedValue := rawVal & ((0b11111111 << lsb) & (0b11111111 >> (7 - msb)))
  	return(maskedValue)
}

func (r *radio) begin() {
	r.config()
	r.setFreq(r.frequency)
	r.setBitrate(r.bitRate)
	r.setRxBandwidth(r.rxbw)
	r.setFrequencyDeviation(r.freqdev)
	r.setOutputPower(r.pwr)
	
}

func (r *radio) setRxBandwidth(bw float32) {
    if bw < 58.0 || bw > 812.0 {
        println("Invalid RXBW value")
        return
    } else {
        r.standby()
        var e int8 = 0
        var m int8 = 0
        for e = 3; e >= 0; e-- {
            for m = 3; m >= 0; m-- {
                point := (CC1101_CRYSTAL_FREQ * 1000000.0)/(8.0 * float32(m + 4) * float32(uint32(1)<<e))
                if math.Abs (float64(bw * 1000.0 - point)) <= 1000.0 {
                    r.setRegValue(CC1101_REG_MDMCFG4, uint8((e << 6) | (m << 4)), 7, 4)
                    r.rxbw = bw
                    break
                }
            }
        }
    }
}

func (r *radio) setSyncWord(data []byte) {
	if len(data) != 2 {
		println("Syncword must be a string consisting of 2 bytes")
		return
	}
	for _, element := range data {
		if element == 0x00 {
			println("Invalid SyncWord")
			return
		}
	}
	//errbits, CarrierSense need to add
	r.WriteReg(CC1101_REG_SYNC1, data[0])
  	r.WriteReg(CC1101_REG_SYNC0, data[1])

}

func (r *radio) setFrequencyDeviation(dev float32) {
    if dev < 0.0 {
        dev = 1.587
    }
    if dev != 0 {
        if dev >= 1.587 && dev <= 380.8 {
            r.standby()
            e, m := getExpMant((dev * 1000.0), 8, 17, 7)
	        if e == 0 && m == 0 {
		        println("Error Calculating Freq Dev Exponent and Mantissa")
		        return
	        }
            r.setRegValue(CC1101_REG_DEVIATN, (e << 4), 6, 4)
            r.setRegValue(CC1101_REG_DEVIATN, m, 2, 0)
            r.freqdev = dev
        }
    }
}

func (r *radio) config() {
	r.Strobe(CC1101_CMD_RESET)
	time.Sleep(150 * time.Millisecond)
	r.standby()
	r.setRegValue(CC1101_REG_MCSM0, CC1101_FS_AUTOCAL_IDLE_TO_RXTX, 5, 4)
	r.setRegValue(CC1101_REG_MCSM0, CC1101_PIN_CTRL_OFF, 1, 1)
  	r.setRegValue(CC1101_REG_IOCFG0, CC1101_GDOX_HIGH_Z, 5, 0)
  	r.setRegValue(CC1101_REG_IOCFG2, CC1101_GDOX_HIGH_Z, 5, 0)

  	r.packetMode(true)
}

func (r *radio) packetMode(tx bool) {
	r.setRegValue(CC1101_REG_PKTCTRL1, CC1101_CRC_AUTOFLUSH_OFF | CC1101_APPEND_STATUS_ON | CC1101_ADR_CHK_NONE, 3, 0)
	r.setRegValue(CC1101_REG_PKTCTRL0, CC1101_WHITE_DATA_OFF | CC1101_PKT_FORMAT_NORMAL, 6, 4)
  	r.setRegValue(CC1101_REG_PKTCTRL0, CC1101_CRC_ON | CC1101_LENGTH_CONFIG_VARIABLE, 2, 0)
	if tx {
  		r.Strobe(CC1101_CMD_FLUSH_TX)
  	} else {
  		r.Strobe(CC1101_CMD_FLUSH_RX)
  	}
}

func (r *radio) setOutputPower(power int8) {
	allowedPwrs := []int8{ -30, -20, -15, -10, 0, 5, 7, 10 }
	
	if power <= -30 {
		power = -30
	} else if power >= 10 {
		power = 10
	} else {
		for i := 0; i < 8; i++ {
			if allowedPwrs[i] > power {
				r.pwr = allowedPwrs[i]
				power = allowedPwrs[i]
				break
			}
		}
	}
	var f uint8 = 0
	if r.frequency < 374.0 {
		f = 0
	} else if r.frequency < 650.0 {
		f = 1
	} else if r.frequency < 891.5 {
		f = 2
	} else {
		f = 3
	}
	var paTable [8][4]uint8 = [8][4]uint8{
    				{0x12, 0x12, 0x03, 0x03},
    				{0x0D, 0x0E, 0x0F, 0x0E},
    				{0x1C, 0x1D, 0x1E, 0x1E},
    				{0x34, 0x34, 0x27, 0x27},
    				{0x51, 0x60, 0x50, 0x8E},
    				{0x85, 0x84, 0x81, 0xCD},
    				{0xCB, 0xC8, 0xCB, 0xC7},
    				{0xC2, 0xC0, 0xC2, 0xC0},
				}
	 for i := 0; i < 8; i++ {
    		if power == allowedPwrs[i] {
    			if r.modulation == CC1101_MOD_FORMAT_ASK_OOK {
      			 	paValues :=[]byte{0x00, paTable[i][f]}
      			 	r.WriteRegBurst(CC1101_REG_PATABLE, paValues)
      			} else {
      				r.WriteReg(CC1101_REG_PATABLE, paTable[i][f])
      			  }
    		}
    	}
}
	//Need carriersense option
	//Need funcs for each
func (r *radio) promiscuous () {
	//set pqt 0
	r.setRegValue(CC1101_REG_PKTCTRL1, uint8(0 << 5), 7, 5)
	//disable syncword filter
	r.setRegValue(CC1101_REG_MDMCFG2, CC1101_SYNC_MODE_NONE, 2, 0)
	//disable crc
	r.setRegValue(CC1101_REG_PKTCTRL0, CC1101_CRC_OFF, 2, 2)
	//disable addr
	r.setRegValue(CC1101_REG_PKTCTRL1, CC1101_DEVICE_ADDR, 1, 0)
}

func (r *radio) standby() {
	r.Strobe(CC1101_CMD_IDLE)
	for {
		if r.getRegValue(CC1101_REG_MARCSTATE, 4, 0) == CC1101_MARC_STATE_IDLE {
			break
		}
	}
}

func (r *radio) transmitAsync(data *[]byte) {

	r.Strobe(CC1101_CMD_IDLE)
	r.WriteReg(CC1101_REG_IOCFG1, CC1101_GDOX_HIGH_Z)
	r.WriteReg(CC1101_REG_IOCFG2, CC1101_GDOX_HIGH_Z)
	r.GDO0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	r.setRegValue(CC1101_REG_PKTCTRL0, 0x30, 5, 0)
	r.setRegValue(CC1101_REG_PKTCTRL0, 0x02, 1, 0)
	
	//High baud in order to oversample data
	//duration is calc on OG bitRate so its fine
	e, m := getExpMant((36400 * 1000.0), 256, 28, 14)
	r.setRegValue(CC1101_REG_MDMCFG4, e, 3, 0)
	r.WriteReg(CC1101_REG_MDMCFG3, m)
	
	//delay per bit 
	duration := time.Duration((1/r.bitRate) * 1000)
	
  	r.Strobe(CC1101_CMD_TX)
  	for {
		if r.getRegValue(CC1101_REG_MARCSTATE, 4, 0) == 0x13 {
			break
		}
	}
  	
  	for _, b := range (*data) {
  		for i := 7; i >= 0; i-- {
			r.GDO0.Set(((b>>uint(i))&0x01) != 0)
 			delay.Sleep(duration * time.Microsecond)
  		}
  	}
  	
  	r.GDO0.Set(false)
  	r.packetMode(true)
  	//reset baud to user value
  	r.setBitrate(r.bitRate)
  	r.setRegValue(CC1101_REG_IOCFG0, CC1101_GDOX_HIGH_Z, 5, 0)
  	println("Sent")
}



func (r *radio) transmit (data *[]byte, length uint8) {
	r.standby()
	r.promiscuous()
	r.setRegValue(CC1101_REG_PKTCTRL0, CC1101_LENGTH_CONFIG_FIXED, 1, 0)
	r.Strobe(CC1101_CMD_FLUSH_TX)
	r.Strobe(CC1101_CMD_FSTXON)
	for {
		if r.getRegValue(CC1101_REG_MARCSTATE, 4, 0) == 0x12 {
			break
		}
	}
	if length > uint8(CC1101_MAX_PACKET_LENGTH) {
		println ("Tx error, max packet size is 255 bytes")
		for{}
	}
	r.WriteReg(CC1101_REG_PKTLEN, length)
	//Calc delay for a byte to leave the fifo
    	var duration uint32 = uint32(8000 / r.bitRate)

	initialWrite := min(length, uint8(CC1101_FIFO_SIZE))
  	r.WriteRegBurst(CC1101_REG_FIFO, (*data)[:initialWrite])
  	
  	datasent := initialWrite
  	
  	r.Strobe(CC1101_CMD_TX)
  	
  	for datasent < length {
		//Give time for at least a byte to be out of the fifo
		delay.Sleep(time.Duration(duration) * time.Microsecond)
		
		fifobytes := r.getRegValue(CC1101_REG_TXBYTES, 6, 0)
		if fifobytes < CC1101_FIFO_SIZE {
			bytesToWrite := min(uint8(CC1101_FIFO_SIZE - fifobytes), (length - datasent))
        		r.WriteRegBurst(CC1101_REG_FIFO, (*data)[datasent:datasent+bytesToWrite])
        		datasent += bytesToWrite
		}
  	}
  	for {
  		if r.getRegValue(CC1101_REG_TXBYTES, 6, 0) == 0 {
  			r.standby()
  			r.Strobe(CC1101_CMD_FLUSH_TX)
  			break
  		}
  		delay.Sleep(time.Duration(duration) * time.Microsecond)
  	}
  	println("Sent")
}

func parseFloat(s string) float32 {
	val, _ := strconv.ParseFloat(s, 32)
	return float32(val)
}

func parseSignedInt(s string) int32 {
	val, _ := strconv.Atoi(s)
	return int32(val)
}

