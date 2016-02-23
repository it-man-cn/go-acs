package main

// Start Commond eg: ./client 1 5000 localhost:8080
// first parameter：beginning userId
// second parameter: amount of clients
// third parameter: comet server ip

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	begin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	num, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	url := os.Args[3]

	for i := begin; i < begin+num; i++ {
		key := fmt.Sprintf("16000%d", i)
		quit := make(chan bool, 1)
		go startClient(url, key, quit)
	}

	var exit chan bool
	<-exit
}

var data = `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0"><SOAP-ENV:Header><cwmp:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.Inform958064891.1658176487</cwmp:ID><cwmp:NoMoreRequests>0</cwmp:NoMoreRequests></SOAP-ENV:Header><SOAP-ENV:Body><cwmp:Inform><DeviceId xsi:type="cwmp:DeviceIdStruct"><Manufacturer xsi:type="xsd:string">ACS</Manufacturer><OUI xsi:type="xsd:string">OO11AB</OUI><ProductClass xsi:type="xsd:string">it-man</ProductClass><SerialNumber xsi:type="xsd:string">%s</SerialNumber></DeviceId><Event SOAP-ENC:arrayType="cwmp:EventStruct[1]"><EventStruct><EventCode xsi:type="xsd:string">2 PERIODIC</EventCode><CommandKey></CommandKey></EventStruct></Event><MaxEnvelopes xsi:type="xsd:unsignedInt">0</MaxEnvelopes><CurrentTime xsi:type="xsd:string">2015-01-21T14:46:07</CurrentTime><RetryCount xsi:type="xsd:unsignedInt">0</RetryCount><ParameterList SOAP-ENC:arrayType="cwmp:ParameterValueStruct[11]"><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.DeviceInfo.ProvisioningCode</Name><Value xsi:type="string">it-man</Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.ManagementServer.URL</Name><Value xsi:type="string"></Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.ManagementServer.Password</Name><Value xsi:type="string"></Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.ManagementServer.ConnectionRequestURL</Name><Value xsi:type="string">http://192.168.16.68:5400</Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.DeviceSummary</Name><Value xsi:type="string"></Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.DeviceInfo.SpecVersion</Name><Value xsi:type="string">V1.0</Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.ManagementServer.UDPConnectionRequestAddress</Name><Value xsi:type="string">200.200.202.68:1036</Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.ManagementServer.ConnectionRequestPassword</Name><Value xsi:type="string"></Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.DeviceInfo.SoftwareVersion</Name><Value xsi:type="string">it-man.bin-150120</Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress</Name><Value xsi:type="string">192.168.16.68</Value></ParameterValueStruct><ParameterValueStruct><Name xsi:type="string">InternetGatewayDevice.DeviceInfo.HardwareVersion</Name><Value xsi:type="string">V1.0</Value></ParameterValueStruct></ParameterList></cwmp:Inform></SOAP-ENV:Body></SOAP-ENV:Envelope>`

func startClient(url, key string, quit chan bool) {
	for {
		http.Post(url, "text/xml", strings.NewReader(fmt.Sprintf(data, key)))
		time.Sleep(30 * time.Second)
	}

}
