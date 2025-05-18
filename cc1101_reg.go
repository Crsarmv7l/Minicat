package main

const (
	// CC1101 physical layer properties
CC1101_FREQUENCY_STEP_SIZE         =            396.7285156
CC1101_MAX_PACKET_LENGTH           =            255
CC1101_FIFO_SIZE                   =            64
CC1101_CRYSTAL_FREQ                =            26.0
CC1101_DIV_EXPONENT                =            16

// CC1101 SPI commands
CC1101_CMD_READ_SINGLE             =            0x80
CC1101_CMD_READ_BURST              =            0xC0
CC1101_CMD_WRITE                   =            0x00
CC1101_CMD_WRITE_BURST             =            0x40
CC1101_CMD_RESET                   =            0x30
CC1101_CMD_FSTXON                  =            0x31
CC1101_CMD_TX                      =            0x35
CC1101_CMD_IDLE                    =            0x36
CC1101_CMD_FLUSH_RX                =            0x3A
CC1101_CMD_FLUSH_TX                =            0x3B

// CC1101 register map
CC1101_REG_IOCFG2                  =            0x00
CC1101_REG_IOCFG0                  =            0x02
CC1101_REG_SYNC1                   =            0x04
CC1101_REG_SYNC0                   =            0x05
CC1101_REG_PKTLEN                  =            0x06
CC1101_REG_PKTCTRL1                =            0x07
CC1101_REG_PKTCTRL0                =            0x08


CC1101_REG_ADDR                    =            0x09
CC1101_REG_CHANNR                  =            0x0A
CC1101_REG_FSCTRL1                 =            0x0B
CC1101_REG_FSCTRL0                 =            0x0C
CC1101_REG_FREQ2                   =            0x0D
CC1101_REG_FREQ1                   =            0x0E
CC1101_REG_FREQ0                   =            0x0F
CC1101_REG_MDMCFG4                 =            0x10
CC1101_REG_MDMCFG3                 =            0x11
CC1101_REG_MDMCFG2                 =            0x12
CC1101_REG_MDMCFG1                 =            0x13
CC1101_REG_MDMCFG0                 =            0x14
CC1101_REG_DEVIATN                 =            0x15
CC1101_REG_MCSM2                   =            0x16
CC1101_REG_MCSM1                   =            0x17
CC1101_REG_MCSM0                   =            0x18
CC1101_REG_FOCCFG                  =            0x19
CC1101_REG_BSCFG                   =            0x1A
CC1101_REG_AGCCTRL2                =            0x1B
CC1101_REG_AGCCTRL1                =            0x1C
CC1101_REG_AGCCTRL0                =            0x1D
CC1101_REG_WOREVT1                 =            0x1E
CC1101_REG_WOREVT0                 =            0x1F
CC1101_REG_WORCTRL                 =            0x20
CC1101_REG_FREND1                  =            0x21
CC1101_REG_FREND0                  =            0x22
CC1101_REG_FSCAL3                  =            0x23
CC1101_REG_FSCAL2                  =            0x24
CC1101_REG_FSCAL1                  =            0x25
CC1101_REG_FSCAL0                  =            0x26
CC1101_REG_RCCTRL1                 =            0x27
CC1101_REG_RCCTRL0                 =            0x28
CC1101_REG_FSTEST                  =            0x29
CC1101_REG_PTEST                   =            0x2A
CC1101_REG_AGCTEST                 =            0x2B
CC1101_REG_TEST2                   =            0x2C
CC1101_REG_TEST1                   =            0x2D
CC1101_REG_TEST0                   =            0x2E
CC1101_REG_PARTNUM                 =            0x30
CC1101_REG_VERSION                 =            0x31
CC1101_REG_FREQEST                 =            0x32
CC1101_REG_LQI                     =            0x33
CC1101_REG_RSSI                    =            0x34
CC1101_REG_MARCSTATE               =            0x35
CC1101_REG_WORTIME1                =            0x36
CC1101_REG_WORTIME0                =            0x37
CC1101_REG_PKTSTATUS               =            0x38
CC1101_REG_VCO_VC_DAC              =            0x39
CC1101_REG_TXBYTES                 =            0x3A
CC1101_REG_RXBYTES                 =            0x3B
CC1101_REG_RCCTRL1_STATUS          =            0x3C
CC1101_REG_RCCTRL0_STATUS          =            0x3D
CC1101_REG_PATABLE                 =            0x3E
CC1101_REG_FIFO                    =            0x3F
CC1101_MARC_STATE_IDLE             =            0x01
CC1101_FS_AUTOCAL_IDLE_TO_RXTX     =            0b00010000
CC1101_PIN_CTRL_OFF                =            0b00000000
CC1101_LENGTH_CONFIG_VARIABLE      =            0b00000001
CC1101_CRC_AUTOFLUSH_OFF           =            0b00000000
CC1101_APPEND_STATUS_ON            =            0b00000100
CC1101_APPEND_STATUS_OFF           =            0b00000000
CC1101_ADR_CHK_NONE                =            0b00000000
CC1101_WHITE_DATA_OFF              =            0b00000000 
CC1101_PKT_FORMAT_NORMAL           =            0b00000000
CC1101_SYNC_MODE_NONE              =            0b00000000
CC1101_LENGTH_CONFIG_FIXED         =            0b00000000 
CC1101_CRC_OFF                     =            0b00000000
CC1101_CRC_ON                      =            0b00000100
CC1101_NUM_PREAMBLE_2              =            0b00000000
CC1101_DEVICE_ADDR                 =            0x00 
CC1101_GDOX_HIGH_Z                 =            0x2E        

// defaults
CC1101_DEFAULT_FREQ                =            434.0
CC1101_DEFAULT_BR                  =            4.8
CC1101_DEFAULT_FREQDEV             =            5.0
CC1101_DEFAULT_RXBW                =            58.0
CC1101_DEFAULT_POWER               =            10
CC1101_DEFAULT_PREAMBLELEN         =            16
CC1101_DEFAULT_SW1                 =            0x12
CC1101_DEFAULT_SW2 		   =		0xAD
CC1101_DEFAULT_SW_LEN              =            2
CC1101_MOD_FORMAT_2_FSK            =           0b00000000
CC1101_MOD_FORMAT_ASK_OOK          =           0b00110000
)
