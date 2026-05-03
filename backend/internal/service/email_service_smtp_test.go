//go:build unit

package service

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailService_TestSMTPConnectionWithConfig_StartTLS(t *testing.T) {
	var sawStartTLS atomic.Bool
	addr := startTestSMTPServer(t, true, func(command string) {
		if strings.EqualFold(command, "STARTTLS") {
			sawStartTLS.Store(true)
		}
	})
	host, port := splitTCPAddr(t, addr)

	svc := NewEmailService(nil, nil)
	err := svc.TestSMTPConnectionWithConfig(&SMTPConfig{
		Host:     host,
		Port:     port,
		Username: "user",
		Password: "pass",
		UseTLS:   false,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "starttls failed")
	require.True(t, sawStartTLS.Load(), "expected connection test to attempt STARTTLS")
}

func TestEmailService_TestSMTPConnectionWithConfig_PlainAuth(t *testing.T) {
	var sawAuth atomic.Bool
	addr := startTestSMTPServer(t, false, func(command string) {
		if strings.HasPrefix(strings.ToUpper(command), "AUTH ") {
			sawAuth.Store(true)
		}
	})
	host, port := splitTCPAddr(t, addr)

	svc := NewEmailService(nil, nil)
	err := svc.TestSMTPConnectionWithConfig(&SMTPConfig{
		Host:     host,
		Port:     port,
		Username: "user",
		Password: "pass",
		UseTLS:   false,
	})

	require.NoError(t, err)
	require.True(t, sawAuth.Load(), "expected connection test to authenticate")
}

func startTestSMTPServer(t *testing.T, advertiseStartTLS bool, onCommand func(string)) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()

		reader := bufio.NewReader(conn)
		_, _ = fmt.Fprint(conn, "220 localhost ESMTP\r\n")
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			command := strings.TrimSpace(line)
			onCommand(command)

			upper := strings.ToUpper(command)
			switch {
			case strings.HasPrefix(upper, "EHLO ") || strings.HasPrefix(upper, "HELO "):
				if advertiseStartTLS {
					_, _ = fmt.Fprint(conn, "250-localhost\r\n250-STARTTLS\r\n250 AUTH PLAIN\r\n")
				} else {
					_, _ = fmt.Fprint(conn, "250-localhost\r\n250 AUTH PLAIN\r\n")
				}
			case upper == "STARTTLS":
				_, _ = fmt.Fprint(conn, "220 Ready to start TLS\r\n")
				return
			case strings.HasPrefix(upper, "AUTH "):
				_, _ = fmt.Fprint(conn, "235 Authentication successful\r\n")
			case upper == "QUIT":
				_, _ = fmt.Fprint(conn, "221 Bye\r\n")
				return
			default:
				_, _ = fmt.Fprint(conn, "250 OK\r\n")
			}
		}
	}()

	return listener.Addr().String()
}

func splitTCPAddr(t *testing.T, addr string) (string, int) {
	t.Helper()

	host, portString, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	var port int
	_, err = fmt.Sscanf(portString, "%d", &port)
	require.NoError(t, err)
	return host, port
}
