package sysinfo

import (
	"bytes"
	"fmt"
	"github.com/dullkingsman/go-pkg/utils"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type CloudProvider uint

const (
	CloudProviderUnknown CloudProvider = iota
	CloudProviderAWS
	CloudProviderGCP
	CloudProviderAzure
	CloudProviderVmware
	CloudProviderVirtualBox
	CloudProviderKvmQemu
)

// getCloudInstanceID retrieves the instance ID based on the detected cloud provider.
func getCloudInstanceID(provider CloudProvider) *string {
	var client = &http.Client{Timeout: 5 * time.Second} // Longer timeout for actual data retrieval

	var instanceID *string

	switch provider {
	case CloudProviderAWS:
		// For AWS, first get a token for IMDSv2
		var tokenReq, _ = http.NewRequest("PUT", "http://169.254.169.254/latest/api/token", nil)

		tokenReq.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")

		tokenResp, tokenErr := client.Do(tokenReq)

		if tokenErr != nil || tokenResp.StatusCode != http.StatusOK {
			fmt.Printf("Error getting AWS IMDSv2 token: %v\n", tokenErr)
			return nil
		}

		defer func(Body io.ReadCloser) {
			var err = Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}(tokenResp.Body)

		token, _ := io.ReadAll(tokenResp.Body)

		var req, _ = http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/instance-id", nil)

		req.Header.Set("X-aws-ec2-metadata-token", string(token))

		var resp, err = client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			defer func(Body io.ReadCloser) {
				var err = Body.Close()

				if err != nil {
					fmt.Println(err)
				}
			}(resp.Body)

			var body, _ = io.ReadAll(resp.Body)

			instanceID = utils.PtrOf(strings.TrimSpace(string(body)))
		} else {
			fmt.Printf("Error getting AWS instance ID: %v\n", err)
		}

	case CloudProviderGCP:
		var req, _ = http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/id", nil)

		req.Header.Set("Metadata-Flavor", "Google")

		var resp, err = client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			defer func(Body io.ReadCloser) {
				var err = Body.Close()

				if err != nil {
					fmt.Println(err)
				}
			}(resp.Body)

			var body, _ = io.ReadAll(resp.Body)

			instanceID = utils.PtrOf(strings.TrimSpace(string(body)))
		} else {
			fmt.Printf("Error getting GCP instance ID: %v\n", err)
		}

	case CloudProviderAzure:
		var req, _ = http.NewRequest("GET", "http://169.254.169.254/metadata/instance/compute/vmId?api-version=2021-02-01", nil)
		req.Header.Set("Metadata", "true")
		var resp, err = client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			defer func(Body io.ReadCloser) {
				var err = Body.Close()

				if err != nil {
					fmt.Println(err)
				}
			}(resp.Body)

			var body, _ = io.ReadAll(resp.Body)

			instanceID = utils.PtrOf(strings.TrimSpace(string(body)))

		} else {
			fmt.Printf("Error getting Azure instance ID: %v\n", err)
		}
	}

	return instanceID
}

// getCloudProvider attempts to detect the current cloud provider.
// It checks for well-known metadata service IP addresses and headers.
func getCloudProvider(product ProductIdentifier) CloudProvider {
	if product.Name != nil && *product.Name != "" {
		var lowerProductName = strings.ToLower(*product.Name)

		switch true {
		case strings.Contains(lowerProductName, strings.ToLower("HVM domU")),
			strings.Contains(lowerProductName, strings.ToLower("Amazon")),
			strings.Contains(lowerProductName, strings.ToLower("EC2")),
			strings.Contains(lowerProductName, strings.ToLower("Amazon EC2")):
			return CloudProviderAWS
		case strings.Contains(lowerProductName, strings.ToLower("Google Compute Engine")),
			strings.Contains(lowerProductName, strings.ToLower("Google")):
			return CloudProviderGCP
		case strings.Contains(lowerProductName, strings.ToLower("Virtual Machine")),
			strings.Contains(lowerProductName, strings.ToLower("Azure")):
			return CloudProviderAzure
		case strings.Contains(lowerProductName, strings.ToLower("VMware Virtual Platform")),
			strings.Contains(lowerProductName, strings.ToLower("VMware")):
			return CloudProviderVmware
		case strings.Contains(lowerProductName, strings.ToLower("VirtualBox")):
			return CloudProviderVirtualBox
		case strings.Contains(lowerProductName, strings.ToLower("Standard PC")),
			strings.Contains(lowerProductName, strings.ToLower("KVM")),
			strings.Contains(lowerProductName, strings.ToLower("QEMU")),
			strings.Contains(lowerProductName, strings.ToLower("QEMU Virtual Machine")):
			return CloudProviderKvmQemu
		}
	}

	var client = &http.Client{Timeout: 2 * time.Second}

	var req, _ = http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/", nil)

	resp, err := client.Do(req)

	if err == nil {
		defer func(Body io.ReadCloser) {
			var err = Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)

		if resp.StatusCode == http.StatusOK && resp.Header.Get("Server") == "EC2" {
			return CloudProviderAWS
		}
	}

	req, _ = http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/", nil)

	req.Header.Set("Metadata-Flavor", "Google")

	resp, err = client.Do(req)

	if err == nil {
		defer func(Body io.ReadCloser) {
			var err = Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)

		if resp.StatusCode == http.StatusOK {
			return CloudProviderGCP
		}
	}

	req, _ = http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)

	req.Header.Set("Metadata", "true")

	resp, err = client.Do(req)

	if err == nil {
		defer func(Body io.ReadCloser) {
			var err = Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)

		if resp.StatusCode == http.StatusOK {
			var body, _ = io.ReadAll(resp.Body)

			if bytes.Contains(body, []byte(`"azEnvironment"`)) { // A common field in Azure IMDS response
				return CloudProviderAzure
			}
		}
	}

	return CloudProviderUnknown
}

func getProductIdentifier() ProductIdentifier {
	var product = ProductIdentifier{}

	switch runtime.GOOS {
	case "linux":
		var name, err = os.ReadFile("/sys/class/dmi/id/product_name")

		if err != nil {
			fmt.Printf("error reading product name: %v\n", err)
		} else {
			product.Name = utils.PtrOf(string(name))
		}

		uuid, err := os.ReadFile("/sys/class/dmi/id/product_uuid")

		if err != nil {
			fmt.Printf("error reading product uuid: %v\n", err)
		} else {
			product.Uuid = utils.PtrOf(string(uuid))
		}

		serial, err := os.ReadFile("/sys/class/dmi/id/product_serial")

		if err != nil {
			fmt.Printf("error reading product serial: %v\n", err)
		} else {
			product.Serial = utils.PtrOf(string(serial))
		}
	default:
		fmt.Printf("Hardware UUID retrieval not implemented for %s\n", runtime.GOOS)
	}

	return product
}

// getMachineID attempts to retrieve the OS-level machine ID.
func getMachineID() *string {
	var machineID *string

	switch runtime.GOOS {
	case "linux":
		var content, err = os.ReadFile("/etc/machine-id")

		if err == nil {
			if len(content) > 0 {
				machineID = utils.PtrOf(strings.TrimSpace(string(content)))
			}
		} else {
			fmt.Printf("error reading /etc/machine-id: %v\n", err)
		}
	default:
		fmt.Printf("Machine ID retrieval not implemented for %s\n", runtime.GOOS)
	}

	return machineID
}

// getMacAddress retrieves the MAC address of the first non-loopback network interface.
func getMacAddress() *string {
	var interfaces, err = net.Interfaces()

	if err != nil {
		fmt.Printf("Error getting network interfaces: %v\n", err)
		return nil
	}

	for _, _interface := range interfaces {
		if _interface.Flags&net.FlagUp != 0 && _interface.Flags&net.FlagLoopback == 0 {
			var mac = _interface.HardwareAddr.String()

			if mac != "" {
				return &mac
			}
		}
	}

	return nil
}
