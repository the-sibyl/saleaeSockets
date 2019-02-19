package main

import(
	"fmt"
	"net"
	"time"
)

const (
	DefaultAddressString string = "localhost:10429"
	DefaultReceiveBufferSize int = 4096
	DefaultIOTimeoutMilliseconds = 1000
)

func main() {
	sd := NewSaleaeDevice()
	sd.initialize()

	for ;; {
		fmt.Println("num samples test:\n")
		sd.sendCommand("get_num_samples", false)
		fmt.Println("connected devices test:\n")
		sd.sendCommand("get_connected_devices", false)
		fmt.Println("Done with test! Repeating in a few seconds...")
		time.Sleep(time.Second * 10)
	}
}

type SaleaeDevice struct {
	// "host:port" format
	addressString string
	// Socket connection
	saleaeConn *net.TCPConn
	// Input buffer for the socket
	tcpInputBuffer []byte
	// Number of bytes received
	readLength int
	// Timer for I/O timeouts
//	ioTimer *time.Timer
	// Duration for this timer
//	timerDuration time.Duration
	// I/O Timeout Duration
	ioTimeoutDuration time.Duration
}

// Create a new Saleae "object"
func NewSaleaeDevice() (*SaleaeDevice) {
	var dev SaleaeDevice
	dev.addressString = DefaultAddressString
	return &dev
}

func (dev *SaleaeDevice) initialize() error {
	// Input buffer byte slice
	dev.tcpInputBuffer = make([]byte, DefaultReceiveBufferSize)

	// Create an address for the socket
	tcpAddr, err := net.ResolveTCPAddr("tcp", dev.addressString)
	if err != nil {
		return err
	}

	// Attempt a connection
	dev.saleaeConn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	// TODO: clean this up
//	dev.saleaeConn.SetReadDeadline(time.Millisecond * DefaultIOTimeoutMilliseconds)

//	dev.timerDuration = time.Millisecond * DefaultIOTimeoutMilliseconds
//	dev.ioTimer = time.NewTimer(dev.timerDuration)
	dev.ioTimeoutDuration = time.Millisecond * DefaultIOTimeoutMilliseconds

	return nil
}

func (dev *SaleaeDevice) sendCommand(cmd string, largeTransfer bool) (err error) {
	fmt.Println("Called")
	// Null character termination
	_, err = dev.saleaeConn.Write([]byte(cmd + "\000"))
	if err != nil {
		return err
	}

	// TODO: Look for ack

	rc := make(chan bool)

	//dev.ioTimer.Reset(dev.timerDuration)

	socketReadLoop:
	for ;; {
		go func(readChan chan bool) {
				fmt.Println("READCHAN!")
				dev.saleaeConn.SetReadDeadline(time.Now().Add(dev.ioTimeoutDuration))
				dev.readLength, _ = dev.saleaeConn.Read(dev.tcpInputBuffer)
				fmt.Println("done reading")
				fmt.Println("Count: ", dev.readLength)

				// If ACK, return true to the channel
				fmt.Println("setting readChan true")
				readChan <- true

				time.Sleep(time.Millisecond * 50)
				fmt.Println("persisting")
		} (rc)

		// 
		select {
			// Read channel
			case <- rc:
				fmt.Println("read channel!")
				break socketReadLoop

				/*
			// Timer channel
			case <- dev.ioTimer.C:
				fmt.Println("timer expired")
				break socketReadLoop
				*/
		}
	}

	fmt.Println("returning")
	//dev.readLength = 0


	return nil

}

func (dev *SaleaeDevice) GetNumSamples() {

}

// set_trigger
// set_num_samples
// get_num_samples
// set_capture_seconds
// set_sample_rate
// get_sample_rate
// get_all_sample_rates
// get_performance
// set_performance
// get_capture_pretrigger_buffer_size
// set_capture_pretriffer_buffer_size
// get_connected_devices
// select_active_device
// get_active_channels
// set_active_channels
// reset_active_channels
// get_digital_voltage_options
// set_digital_voltage_option
// capture
// stop_capture
// capture_to_file
// is_processing_complete
// save_to_file
// load_from_file
// close_all_tabs
// export_data2
// export_data
// get_analyzers
// export_analyzer
// is_analyzer_complete
// get_capture_range
// get_viewstate
// set_viewstate (set_viewstate mode)
