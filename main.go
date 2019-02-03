package main

import(
	"fmt"
	"net"
	"time"
)

const (
	DefaultAddressString string = "localhost:10429"
	DefaultReceiveBufferSize int = 4096
)

func main() {
	sd := NewSaleaeDevice()
	sd.initialize()

	for ;; {
		sd.sendCommand()
		time.Sleep(time.Second)
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

	return nil
}

func (dev *SaleaeDevice) sendCommand() (err error) {
	dev.readLength, err = dev.saleaeConn.Write([]byte("get_num_samples\000"))
	if err != nil {
		return err
	}

	fmt.Println("Count: ", dev.readLength)

	// TODO: Look for ack
	for ;dev.readLength != 0; {
		dev.readLength, _ = dev.saleaeConn.Read(dev.tcpInputBuffer)
		fmt.Println(string(dev.tcpInputBuffer[:dev.readLength]))

		dev.readLength = 0
	}

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
