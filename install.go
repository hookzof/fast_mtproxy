package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func cmd(cmd string) string {
	if out, err := exec.Command("sh", "-c", cmd).Output(); err != nil && err.Error() != "exit status 1" &&
		err.Error() != "exit status 2" {
		log.Println("[error]", err, "("+cmd+")")
		return ""
	} else {
		return string(out)
	}
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getIP() string {
	if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		defer func() {
			err := conn.Close()
			if err != nil {
				log.Println("[error]", err)
			}
		}()

		return conn.LocalAddr().(*net.UDPAddr).IP.String()
	} else {
		log.Println("[error]", err, "(Couldn't identify IP | Не удалось определить IP)")
		return ""
	}
}

func getTrueIP(ver string) string {
	ip := strings.Trim(cmd("curl ifconfig.co -"+ver), "\n")
	if net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("The platform is not supported | Платформа не поддерживается")
		_, _ = fmt.Scanln()
		return
	}

	log.Println("          Starting | Начало работы")

	portStats := flag.String("p", "", "Is the local port for stats")
	port := flag.String("H", "443", "Is the port, used by clients to connect to the proxy")
	secret := flag.String("S", randomHex(16), "Secret")
	tag := flag.String("P", "", "Ad tag get here @MTProxybot")
	domain := flag.String("D", "www.google.com", "Domain with TLS 1.3 support")

	start := flag.String("start", "", "Start server")
	stop := flag.String("stop", "", "Stop server")
	restart := flag.String("restart", "", "Restart server")
	enable := flag.String("enable", "", "Enable server")
	disable := flag.String("disable", "", "Disable server")
	remove := flag.String("delete", "", "Delete server")

	ipv6 := flag.Bool("6", false, "Activation of ipv6")

	flag.Parse()

	path := "/etc/systemd/system/MTProxy-" + *port + ".service"

	if *start != "" {
		cmd("systemctl start MTProxy-" + *start + ".service")

		log.Println(" Server is started | Сервер запущен")
		log.Println(" Program completed | Программа завершена")
		return
	}

	if *stop != "" {
		cmd("systemctl stop MTProxy-" + *stop + ".service")

		log.Println(" Server is stopped | Сервер остановлен")
		log.Println(" Program completed | Программа завершена")
		return
	}

	if *restart != "" {
		cmd("systemctl restart MTProxy-" + *restart + ".service")

		log.Println("  Server is restarted | Сервер перезапущен")
		log.Println("    Program completed | Программа завершена")
		return
	}

	if *enable != "" {
		cmd("systemctl daemon-reload")
		cmd("systemctl restart MTProxy-" + *enable + ".service")
		cmd("systemctl enable MTProxy-" + *enable + ".service")

		log.Println("  Server is enabled | Сервер включен")
		log.Println("  Program completed | Программа завершена")
		return
	}

	if *disable != "" {
		cmd("systemctl stop MTProxy-" + *disable + ".service")
		cmd("systemctl disable MTProxy-" + *disable + ".service")

		log.Println(" Server is disabled | Сервер отключен")
		log.Println("  Program completed | Программа завершена")
		return
	}

	if *remove != "" {
		cmd("systemctl stop MTProxy-" + *remove + ".service")
		cmd("systemctl disable MTProxy-" + *remove + ".service")
		cmd("rm /etc/systemd/system/MTProxy-" + *remove + ".service")

		log.Println("Uninstall complete | Удаление завершено")
		log.Println(" Program completed | Программа завершена")
		return
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Println("A server with such a port is already created, to rewrite it?")
		log.Println("Сервер с таким портом уже создан, перезаписать его?")

		answer := ""
		fmt.Print("\n\nY/N: ")
		_, _ = fmt.Scan(&answer)

		switch answer {
		case "y", "Y":
			cmd("systemctl stop MTProxy-" + *port + ".service")
			cmd("systemctl disable MTProxy-" + *port + ".service")
			cmd("rm " + path)
		case "n", "N":
			log.Println("Program completed | Программа завершена")
			return
		default:
			log.Println("Invalid input | Некорректный ввод")
			log.Println("Program completed | Программа завершена")
			return
		}
	}

	log.Println("  Dependency check | Проверка зависимостей")

	if _, err := os.Stat("/etc/centos-release"); !os.IsNotExist(err) {
		cmd("yum update")
		cmd("yum -y install openssl-devel zlib-devel qrencode")
		cmd("yum -y groupinstall \"Development Tools\"")
	} else if _, err := os.Stat("/etc/fedora-release"); !os.IsNotExist(err) {
		cmd("yum update")
		cmd("yum -y install openssl-devel zlib-devel qrencode")
		cmd("yum -y groupinstall \"Development Tools\"")
	} else {
		cmd("apt update")
		cmd("apt -y install git make build-essential libssl-dev zlib1g-dev qrencode")
	}

	log.Println("        Installing | Установка")
	cmd("git clone https://github.com/hookzof/MTProxy && cd MTProxy && make && cd objs/bin && " +
		"curl -s https://core.telegram.org/getProxySecret -o proxy-secret && " +
		"curl -s https://core.telegram.org/getProxyConfig -o proxy-multi.conf")

	cmd("cd /opt && mkdir mtproxy")
	cmd("cp MTProxy/objs/bin/mtproto-proxy /opt/mtproxy/mtproto-proxy")
	cmd("cp MTProxy/objs/bin/proxy-multi.conf /opt/mtproxy/proxy-multi.conf")
	cmd("cp MTProxy/objs/bin/proxy-secret /opt/mtproxy/proxy-secret")

	cmd("rm -r MTProxy")

	log.Println("Creating a service | Создание службы")
	cmd("touch " + path)

	options := ""
	if *portStats != "" {
		options += " -p " + *portStats
	}

	if *tag != "" {
		options += " -P " + *tag
	}

	if *domain != "www.google.com" {
		options += " -D " + *domain
	} else {
		options += " -D www.google.com"
	}

	v4 := getTrueIP("4")

	if getIP()[:3] == "10." {
		if v4 != "" {
			options += " --nat-info " + getIP() + ":" + v4
		} else {
			log.Println("Couldn't get a real ipv4")
			return
		}
	}

	if *ipv6 {
		options += " -6"
	}

	config := `[Unit]
Description=MTProxy
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/mtproxy
ExecStart=/opt/mtproxy/mtproto-proxy -u nobody -H ` + *port + " -S " + *secret + options + ` --aes-pwd proxy-secret proxy-multi.conf
Restart=on-failure
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target`

	cmd("echo \"" + config + "\" >> " + path)

	cmd("systemctl daemon-reload")
	cmd("systemctl restart MTProxy-" + *port + ".service")
	cmd("systemctl enable MTProxy-" + *port + ".service")

	src := []byte(*domain)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	log.Println(" Program completed | Программа завершена")

	fmt.Println("\n\nServer file path | Путь файлов сервера — /opt/mtproxy/")
	fmt.Println("Config file path | Путь конфиг файла   — " + path)

	if v4 != "" {
		fmt.Println("\n\nIPv4:\ntg://proxy?server=" + v4 + "&port=" + *port + "&secret=ee" + *secret + string(dst) + "\n")
		out, _ := exec.Command("sh", "-c", "qrencode -t ansiutf8 -l L \"tg://proxy?server="+v4+"&port="+*port+"&secret=ee"+*secret+string(dst)+"\"").Output()
		fmt.Println(string(out))
	} else {
		fmt.Println("\n\nCouldn't get a real ipv4")
	}

	if *ipv6 {
		v6 := getTrueIP("6")

		if v6 != "" {
			fmt.Println("IPv6:\ntg://proxy?server=" + v6 + "&port=" + *port + "&secret=ee" + *secret + string(dst) + "\n")
			out, _ := exec.Command("sh", "-c", "qrencode -t ansiutf8 -l L \"tg://proxy?server="+v6+"&port="+*port+"&secret=ee"+*secret+string(dst)+"\"").Output()
			fmt.Println(string(out))
		} else {
			fmt.Println("Couldn't get a real ipv6")
		}
	}
}
