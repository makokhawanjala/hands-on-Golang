package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// Device represents a connected network device
type Device struct {
	MAC         string
	IP          string
	Hostname    string
	FirstSeen   time.Time
	LastSeen    time.Time
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
	Connections int
	IsActive    bool
}

// TrafficStats holds network statistics
type TrafficStats struct {
	TotalBytes    uint64
	TotalPackets  uint64
	HTTPRequests  uint64
	HTTPSRequests uint64
	DNSQueries    uint64
	TCPConns      uint64
	UDPConns      uint64
}

// NetworkMonitor manages network monitoring
type NetworkMonitor struct {
	devices       map[string]*Device
	stats         *TrafficStats
	mutex         sync.RWMutex
	interfaceName string
	handle        *pcap.Handle
	localMAC      string
	stopChan      chan struct{}
}

// NewNetworkMonitor creates a new network monitor
func NewNetworkMonitor(iface string) *NetworkMonitor {
	return &NetworkMonitor{
		devices:       make(map[string]*Device),
		stats:         &TrafficStats{},
		interfaceName: iface,
		stopChan:      make(chan struct{}),
	}
}

// Start begins monitoring network traffic
func (nm *NetworkMonitor) Start() error {
	// Open device for packet capture
	handle, err := pcap.OpenLive(nm.interfaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open device: %v", err)
	}
	nm.handle = handle

	// Get local MAC address
	iface, err := net.InterfaceByName(nm.interfaceName)
	if err == nil {
		nm.localMAC = iface.HardwareAddr.String()
	}

	// Start packet processing
	go nm.processPackets()

	// Start statistics reporter
	go nm.reportStats()

	// Start device cleanup (mark inactive devices)
	go nm.cleanupDevices()

	log.Printf("Started monitoring on interface: %s", nm.interfaceName)
	return nil
}

// Stop stops the network monitor
func (nm *NetworkMonitor) Stop() {
	close(nm.stopChan)
	if nm.handle != nil {
		nm.handle.Close()
	}
}

// processPackets processes captured packets
func (nm *NetworkMonitor) processPackets() {
	packetSource := gopacket.NewPacketSource(nm.handle, nm.handle.LinkType())

	for {
		select {
		case <-nm.stopChan:
			return
		case packet := <-packetSource.Packets():
			nm.analyzePacket(packet)
		}
	}
}

// analyzePacket analyzes a single packet
func (nm *NetworkMonitor) analyzePacket(packet gopacket.Packet) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	// Update total stats
	nm.stats.TotalPackets++
	nm.stats.TotalBytes += uint64(len(packet.Data()))

	// Extract Ethernet layer
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return
	}
	eth := ethLayer.(*layers.Ethernet)

	// Extract IP layer
	var srcIP, dstIP string
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip := ipLayer.(*layers.IPv4)
		srcIP = ip.SrcIP.String()
		dstIP = ip.DstIP.String()
	}

	// Track source device
	srcMAC := eth.SrcMAC.String()
	if srcMAC != nm.localMAC {
		nm.updateDevice(srcMAC, srcIP, uint64(len(packet.Data())), 1, true)
	}

	// Track destination device
	dstMAC := eth.DstMAC.String()
	if dstMAC != nm.localMAC && (eth.DstMAC[0]&0x01) == 0 { // Not multicast
		nm.updateDevice(dstMAC, dstIP, uint64(len(packet.Data())), 1, false)
	}

	// Analyze transport layer
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		nm.stats.TCPConns++

		// Detect HTTP/HTTPS
		if tcp.DstPort == 80 || tcp.SrcPort == 80 {
			nm.stats.HTTPRequests++
		}
		if tcp.DstPort == 443 || tcp.SrcPort == 443 {
			nm.stats.HTTPSRequests++
		}
	}

	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		nm.stats.UDPConns++

		// Detect DNS
		if udp.DstPort == 53 || udp.SrcPort == 53 {
			nm.stats.DNSQueries++
		}
	}
}

// updateDevice updates device information
func (nm *NetworkMonitor) updateDevice(mac, ip string, bytes, packets uint64, isSending bool) {
	device, exists := nm.devices[mac]
	if !exists {
		device = &Device{
			MAC:       mac,
			IP:        ip,
			FirstSeen: time.Now(),
			IsActive:  true,
		}
		nm.devices[mac] = device
		log.Printf("New device detected: MAC=%s IP=%s", mac, ip)
	}

	device.LastSeen = time.Now()
	device.IsActive = true

	if ip != "" && device.IP == "" {
		device.IP = ip
	}

	if isSending {
		device.BytesSent += bytes
		device.PacketsSent += packets
	} else {
		device.BytesRecv += bytes
		device.PacketsRecv += packets
	}
}

// reportStats periodically reports statistics
func (nm *NetworkMonitor) reportStats() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-nm.stopChan:
			return
		case <-ticker.C:
			nm.printStats()
		}
	}
}

// cleanupDevices marks inactive devices
func (nm *NetworkMonitor) cleanupDevices() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-nm.stopChan:
			return
		case <-ticker.C:
			nm.mutex.Lock()
			for _, device := range nm.devices {
				if time.Since(device.LastSeen) > 10*time.Minute {
					device.IsActive = false
				}
			}
			nm.mutex.Unlock()
		}
	}
}

// printStats prints current statistics
func (nm *NetworkMonitor) printStats() {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	fmt.Println("\n" + string([]byte{0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x4A})) // Clear screen
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("           MAWINGU NETWORK MONITOR - STATISTICS")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("Total Traffic: %.2f MB | Packets: %d\n",
		float64(nm.stats.TotalBytes)/(1024*1024), nm.stats.TotalPackets)
	fmt.Printf("HTTP: %d | HTTPS: %d | DNS: %d\n",
		nm.stats.HTTPRequests, nm.stats.HTTPSRequests, nm.stats.DNSQueries)
	fmt.Printf("TCP Connections: %d | UDP: %d\n", nm.stats.TCPConns, nm.stats.UDPConns)

	fmt.Println("\n─────────────────────────────────────────────────────────────")
	fmt.Println("                     CONNECTED DEVICES")
	fmt.Println("─────────────────────────────────────────────────────────────")

	// Sort devices by last seen
	type devEntry struct {
		mac    string
		device *Device
	}
	var devList []devEntry
	for mac, dev := range nm.devices {
		devList = append(devList, devEntry{mac, dev})
	}
	sort.Slice(devList, func(i, j int) bool {
		return devList[i].device.LastSeen.After(devList[j].device.LastSeen)
	})

	activeCount := 0
	for _, entry := range devList {
		dev := entry.device
		if dev.IsActive {
			activeCount++
		}

		status := "🟢"
		if !dev.IsActive {
			status = "🔴"
		}

		totalMB := float64(dev.BytesSent+dev.BytesRecv) / (1024 * 1024)
		fmt.Printf("%s MAC: %s\n", status, dev.MAC)
		fmt.Printf("   IP: %-15s | Active: %v\n", dev.IP, dev.IsActive)
		fmt.Printf("   Sent: %.2f MB (%d packets) | Recv: %.2f MB (%d packets)\n",
			float64(dev.BytesSent)/(1024*1024), dev.PacketsSent,
			float64(dev.BytesRecv)/(1024*1024), dev.PacketsRecv)
		fmt.Printf("   Total: %.2f MB | Last Seen: %s ago\n",
			totalMB, time.Since(dev.LastSeen).Round(time.Second))
		fmt.Println()
	}

	fmt.Printf("Active Devices: %d | Total Devices: %d\n", activeCount, len(nm.devices))
	fmt.Println("═══════════════════════════════════════════════════════════════")
}

// GetDevices returns all tracked devices
func (nm *NetworkMonitor) GetDevices() []*Device {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	devices := make([]*Device, 0, len(nm.devices))
	for _, dev := range nm.devices {
		devices = append(devices, dev)
	}
	return devices
}

// GetStats returns current statistics
func (nm *NetworkMonitor) GetStats() *TrafficStats {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	statsCopy := *nm.stats
	return &statsCopy
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: sudo ./mawingu-monitor <interface>")
		fmt.Println("\nAvailable interfaces:")
		interfaces, err := pcap.FindAllDevs()
		if err != nil {
			log.Fatal(err)
		}
		for _, iface := range interfaces {
			fmt.Printf("  - %s", iface.Name)
			if len(iface.Addresses) > 0 {
				fmt.Printf(" (%s)", iface.Addresses[0].IP)
			}
			fmt.Println()
		}
		os.Exit(1)
	}

	iface := os.Args[1]
	monitor := NewNetworkMonitor(iface)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	if err := monitor.Start(); err != nil {
		log.Fatalf("Failed to start monitor: %v", err)
	}

	fmt.Println("\nPress Ctrl+C to stop monitoring...")

	<-sigChan
	fmt.Println("\n\nStopping monitor...")
	monitor.Stop()

	// Print final stats
	monitor.printStats()
}
