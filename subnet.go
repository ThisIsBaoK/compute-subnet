// This program perform subnet calculation (step-by-step) for IPv4 address.

package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

func main() {
	// PerformSubnet("192.168.63.0/24")
	// fmt.Println()
	// fmt.Println()
	PerformSubnet("160.142.222.158/19")
	// fmt.Println()
	// fmt.Println()
	// PerformSubnet("60.210.14.230/19")
	// fmt.Println()
	// fmt.Println()
	// PerformSubnet("201.22.45.89/26")
	// fmt.Println()
	// fmt.Println()
	// PerformSubnet("201.22.45.89/26")
	// fmt.Println()
	// fmt.Println()
	// PerformSubnet("240.0.0.0/19") // Class E address, expect to fail.
}

func PerformSubnet(cidr string) {
	fmt.Printf("========== Perform subnetting on %s ==========\n", cidr)

	// Parse the CIDR notation.
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	mask := ipnet.Mask

	if ip == nil || mask == nil {
		fmt.Println("Invalid input. Please provide a valid network address and subnet mask.")
		return
	}

	// Calculate the network address and broadcast address.
	network := ip.Mask(mask)
	broadcast := make(net.IP, len(network))

	for i := 0; i < len(network); i++ {
		broadcast[i] = network[i] | ^mask[i]
	}

	// Determine the class of the network address based on the first octet.
	networkClass := ""
	firstOctet := network[0]

	switch {
	case firstOctet >= 1 && firstOctet <= 126:
		networkClass = "Class A"
	case firstOctet >= 128 && firstOctet <= 191:
		networkClass = "Class B"
	case firstOctet >= 192 && firstOctet <= 223:
		networkClass = "Class C"
	default:
		fmt.Println("Could not subnet since the IP address does not belong to class A, B, or C.")
		return
	}

	// Calculate the first and last host addresses.
	firstHost := IncrementIP(network)
	lastHost := DecrementIP(broadcast)

	// Print the results.
	fmt.Printf("Or a mask of %s\n", net.IP(mask).String())

	fmt.Println()
	fmt.Println()

	fmt.Printf("%15s %39s\n", "IP Address:", IPv4ToBinFormat(ip))
	fmt.Printf("%15s %39s\n", "Mask:", IPv4ToBinFormat(net.IP(mask)))
	fmt.Printf("%15s-%39s\n", " ", strings.Repeat("-", 39))
	fmt.Printf("%15s %39s\n", "Id:", IPv4ToBinFormat(network))
	fmt.Printf("%15s+%39s\n", " ", "1")
	fmt.Printf("%15s-%39s\n", " ", strings.Repeat("-", 39))
	fmt.Printf("%15s %39s\n", "1st:", IPv4ToBinFormat(firstHost))

	fmt.Println()
	fmt.Println()

	fmt.Printf("%15s %39s\n", "IP Address:", IPv4ToBinFormat(ip))
	fmt.Printf("%15s %39s\n", "Inverted Mask:", IPv4ToBinFormat(net.IP(InvertIPMask(mask))))
	fmt.Printf("%15s-%39s\n", " ", strings.Repeat("-", 39))
	fmt.Printf("%15s %39s\n", "b/c:", IPv4ToBinFormat(broadcast))
	fmt.Printf("%15s-%39s\n", " ", "1")
	fmt.Printf("%15s-%39s\n", " ", strings.Repeat("-", 39))
	fmt.Printf("%15s %s\n", "Last:", IPv4ToBinFormat(lastHost))

	fmt.Println()
	fmt.Println()

	fmt.Printf("%-25s%-25s%-25s%-25s\n", "", "Last octet (binary)", "Last octet (decimal)", "IP Address")
	fmt.Printf("%-25s%-25s%-25d%-25s\n", "Id", OctetToBin(network[3]), network[3], network.String())
	fmt.Printf("%-25s%-25s%-25d%-25s\n", "1st usable address", OctetToBin(firstHost[3]), firstHost[3], firstHost.String())
	fmt.Printf("%-25s%-25s%-25d%-25s\n", "Last usable address", OctetToBin(lastHost[3]), lastHost[3], lastHost.String())
	fmt.Printf("%-25s%-25s%-25d%-25s\n", "b/c address", OctetToBin(broadcast[3]), broadcast[3], broadcast.String())

	fmt.Println()
	fmt.Println()

	fmt.Printf("Network class: %s\n", networkClass)
}

func IPv4ToBinFormat(ipv4 net.IP) string {
	ipv4 = ipv4.To4()

	var buf bytes.Buffer

	for i, octet := range ipv4 {
		buf.WriteString(OctetToBin(octet))

		if i < 3 {
			buf.WriteString(".")
		}
	}

	return buf.String()
}

func OctetToBin(octet byte) string {
	return BinToNimbles(fmt.Sprintf("%08b", octet))
}

func IncrementIP(ip net.IP) net.IP {
	result := make(net.IP, len(ip))
	copy(result, ip)
	for i := len(ip) - 1; i >= 0; i-- {
		result[i]++
		if result[i] != 0 {
			break
		}
	}
	return result
}

func DecrementIP(ip net.IP) net.IP {
	result := make(net.IP, len(ip))
	copy(result, ip)
	for i := len(ip) - 1; i >= 0; i-- {
		result[i]--
		if result[i] != 255 {
			break
		}
	}
	return result
}

func BinToNimbles(binNumber string) string {
	var result bytes.Buffer

	noSpaceBin := strings.ReplaceAll(binNumber, " ", "")
	leftMostNimbleEndIndex := len(noSpaceBin) % 4

	if leftMostNimbleEndIndex != 0 {
		result.WriteString(noSpaceBin[:leftMostNimbleEndIndex])

		if leftMostNimbleEndIndex+1 < len(noSpaceBin) {
			result.WriteByte(' ')
		}
	}

	for i, c := range noSpaceBin[leftMostNimbleEndIndex:] {
		if i != 0 && i%4 == 0 {
			result.WriteByte(' ')
		}

		result.WriteRune(c)
	}

	return result.String()
}

func InvertIPMask(mask net.IPMask) net.IPMask {
	var invertedMask net.IPMask

	for i := 0; i < len(mask); i++ {
		invertedMask = append(invertedMask, ^mask[i])
	}

	return invertedMask
}
