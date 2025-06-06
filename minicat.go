package main

import (
	"time"
	"encoding/hex"
	"machine"
)

var buf [270]byte
func input() (string, string) {
	var pos int = 0

	for {
		if machine.Serial.Buffered() > 0 {
			c, _ := machine.Serial.ReadByte()
			if c != 10 {
				buf[pos] = c
				pos++
			} else {
				break
			}
		}
	}
	start := -1
	end := -1
	for i := 0; i < pos; i++ {
		if buf[i] == '(' {
			start = i
		} else if buf[i] == ')' {
			end = i
		}
	}
	if start != -1 && end != -1 && start < end {
		for i := 0; i < start; i++ {
            		if buf[i] >= 'a' && buf[i] <= 'z' {
                		buf[i] = buf[i] - 32 
            		}
        	}
		cmd := string(buf[:start]) 
		param := string(buf[start+1:end]) // Param as-is
		return cmd, param
	} else {
		
		for {
			if machine.Serial.Buffered() > 0 {
				machine.Serial.ReadByte()
			} else {
				break 
			}
		}
		print("Bad Parse")
		return "", ""
	}
}

func main() {
	
	cs := machine.PA05
	gdo0 := machine.GPIO25
	gdo2 := machine.GPIO26
	spi := machine.SPI0
	
	radio := newradio(spi, cs, gdo0, gdo2)
		
	for {
		//Wait for serial to be ready
		time.Sleep(10 * time.Millisecond)
		if machine.Serial.DTR() == true && machine.Serial.RTS() == true {
			break
		}
	}
	//Verify CC1101 found
	if radio.ReadReg(CC1101_REG_PARTNUM) == 0x00 && radio.ReadReg(CC1101_REG_VERSION) == 0x14 {
		radio.begin()
		println("CC1101 Init Success")
		println("Enter Commands:")
	} else {
		println("CC1101 Init Error")
		for{}
	}
	
	for {
		
		cmd, param := input()
		
		if cmd != "" {
			switch {
			case cmd == "RFRECV":
			
			case cmd == "RFXMIT":
				if len(param) %2 != 0 {
					param = "0" + param
				}
				decoded, _ := hex.DecodeString(param)
				radio.transmit(&decoded, uint8(len(decoded)))
			case cmd == "SETFREQ":
				radio.setFreq(parseFloat(param))
				println(cmd)
			case cmd == "SETMDMDRATE":
				radio.setBitrate(parseFloat(param))
				println(cmd)
			case cmd == "SETMDMMODULATION":
				if param == "MOD_ASK_OOK" {
					radio.setOOK(true)
				} else {
					radio.setOOK(false)
				}
				println(cmd)
			case cmd == "SETMDMDEVIATN":
				radio.setFrequencyDeviation(parseFloat(param))
			case cmd == "SETMDMCHANBW":
				radio.setRxBandwidth(parseFloat(param))
			case cmd == "SETPOWER":
				radio.setOutputPower(int8(parseSignedInt(param)))
			case cmd == "SETMDMSYNCWORD":
				decoded, _ := hex.DecodeString(param)
				radio.setSyncWord(decoded)
				println(cmd)
			case cmd == "SETENABLEPKTDATAWHITENING":
				println(cmd)
			case cmd == "SETENABLEPKTCRC":
				println(cmd)
			case cmd == "SETENABLEMDMMANCHESTER":
				println(cmd)
			case cmd == "SETPKTPQT":
				println(cmd)
			case cmd == "LOWBALL":
				println(cmd)
			}
		}
	}
}


