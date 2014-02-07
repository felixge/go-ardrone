package navdata

type State uint32

const (
	STATE_FLYING                       State = 1 << 0  // FLY MASK : (0) ardrone is landed, (1) ardrone is flying
	STATE_VIDEO_ENABLED                      = 1 << 1  // VIDEO MASK : (0) video disable, (1) video enable
	STATE_VISION_ENABLED                     = 1 << 2  // VISION MASK : (0) vision disable, (1) vision enable
	STATE_CONTROL_ALGORITHM                  = 1 << 3  // CONTROL ALGO : (0) euler angles control, (1) angular speed control
	STATE_ALTITUDE_CONTROL_ALGORITHM         = 1 << 4  // ALTITUDE CONTROL ALGO : (0) altitude control inactive (1) altitude control active
	STATE_START_BUTTON_STATE                 = 1 << 5  // USER feedback : Start button state
	STATE_CONTROL_COMMAND_ACK                = 1 << 6  // Control command ACK : (0) None, (1) one received
	STATE_CAMERA_READY                       = 1 << 7  // CAMERA MASK : (0) camera not ready, (1) Camera ready
	STATE_TRAVELLING_ENABLED                 = 1 << 8  // Travelling mask : (0) disable, (1) enable
	STATE_USB_READY                          = 1 << 9  // USB key : (0) usb key not ready, (1) usb key ready
	STATE_NAVDATA_DEMO                       = 1 << 10 // Navdata demo : (0) All navdata, (1) only navdata demo
	STATE_NAVDATA_BOOTSTRAP                  = 1 << 11 // Navdata bootstrap : (0) options sent in all or demo mode, (1) no navdata options sent
	STATE_MOTOR_PROBLEM                      = 1 << 12 // Motors status : (0) Ok, (1) Motors problem
	STATE_COMMUNICATION_LOST                 = 1 << 13 // Communication Lost : (1) com problem, (0) Com is ok
	STATE_SOFTWARE_FAULT                     = 1 << 14 // Software fault detected - user should land as quick as possible (1)
	STATE_LOW_BATTERY                        = 1 << 15 // VBat low : (1) too low, (0) Ok
	STATE_USER_EMERGENCY_LANDING             = 1 << 16 // User Emergency Landing : (1) User EL is ON, (0) User EL is OF
	STATE_TIMER_ELAPSED                      = 1 << 17 // Timer elapsed : (1) elapsed, (0) not elapsed
	STATE_MAGNOMETER_NEEDS_CALIBRATION       = 1 << 18 // Magnetometer calibration state : (0) Ok, no calibration needed, (1) not ok, calibration needed
	STATE_ANGLES_OUT_OF_RANGE                = 1 << 19 // Angles : (0) Ok, (1) out of range
	STATE_TOO_MUCH_WIND                      = 1 << 20 // WIND MASK: (0) ok, (1) Too much wind
	STATE_ULTRASONIC_SENSOR_DEAF             = 1 << 21 // Ultrasonic sensor : (0) Ok, (1) deaf
	STATE_CUTOUT_DETECTED                    = 1 << 22 // Cutout system detection : (0) Not detected, (1) detected
	STATE_PIC_VERSION_NUMBER_OK              = 1 << 23 // PIC Version number OK : (0) a bad version number, (1) version number is OK
	STATE_AT_CODEC_THREAD_ON                 = 1 << 24 // ATCodec thread ON : (0) thread OFF (1) thread ON
	STATE_NAVDATA_THREAD_ON                  = 1 << 25 // Navdata thread ON : (0) thread OFF (1) thread ON
	STATE_VIDEO_THREAD_ON                    = 1 << 26 // Video thread ON : (0) thread OFF (1) thread ON
	STATE_ACQUISITION_THREAD_ON              = 1 << 27 // Acquisition thread ON : (0) thread OFF (1) thread ON
	STATE_CONTROL_WATCHDOG_DELAY             = 1 << 28 // CTRL watchdog : (1) delay in control execution (> 5ms), (0) control is well scheduled
	STATE_ADC_WATCHDOG_DELAY                 = 1 << 29 // ADC Watchdog : (1) delay in uart2 dsr (> 5ms), (0) uart2 is good
	STATE_COM_WATCHDOG_PROBLEM               = 1 << 30 // Communication Watchdog : (1) com problem, (0) Com is ok
	STATE_EMERGENCY_LANDING                  = 1 << 31 // Emergency landing : (0) no emergency (1) emergency
)
