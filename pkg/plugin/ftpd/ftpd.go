package main

import (
	"bufio"
	"fmt"
	"net"
	"path"
	"strconv"
	"strings"

	"github.com/riotpot/internal/globals"
	lr "github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/services"
)

// Modified from: https://github.com/shenfeng/ftpd.go

var Plugin string

const (
	name        = "FTP"
	network     = globals.TCP
	port_number = 21
)

func init() {
	Plugin = "Ftpd"
}

func Ftpd() services.Service {
	mx := services.NewPluginService(name, port_number, network)
	// Default user is root
	return &FTP{
		Service: mx,
		pasv:    true,
		root:    "/",
	}
}

type FTP struct {
	services.Service // Anonymous fields from the mixin
	command          net.Conn
	data             net.Conn
	pasv             bool // Passive mode
	username         string
	root             string
	cwd              string
	filepath         *treeNode // Filepath tree structure
}

type treeNode struct {
	name     string
	children []*treeNode
}

func (c *FTP) Run() (err error) {

	port := fmt.Sprintf(":%d", c.GetPort())
	listener, err := net.Listen(c.GetNetwork().String(), port)
	if err != nil {
		lr.Log.Fatal().Err(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			lr.Log.Fatal().Err(err)
		}
		c.command = conn
		c.serve()
	}

}

func (c *FTP) reply(msg string) {
	fmt.Fprintf(c.command, msg+"\r\n")
}

func (c *FTP) serve() {
	localAddr := c.command.LocalAddr().(*net.TCPAddr)
	c.reply("220 Connected to " + localAddr.IP.String())
	reader := bufio.NewReader(c.command)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			lr.Log.Fatal().Err(err)
			return
		}
		// Split line into command and argument
		command := strings.TrimSpace(strings.ToLower(line[0:4]))
		msgs := strings.Split(strings.Trim(line, "\r\n "), " ")[1:]
		if command == "exit" {
			c.command.Close()
			return
		}
		c.handle(command, msgs)
	}

}

func (c *FTP) handle(command string, msgs []string) {
	switch command {
	case "user":
		c.username = msgs[0]
		c.cwd = path.Join("/home/", c.username)
		c.filepath = &treeNode{
			name: "home",
		}
		c.filepath.children = []*treeNode{
			{
				name: c.username,
			},
		}

		c.reply("331 Username ok, send password.")
	case "pass":
		c.reply("230 Login successful.")
	case "syst":
		c.reply("215 UNIX Type: L8.")
	case "type":
		c.reply("200 Type set to binary.")
	case "port":
		port(c, msgs)
	case "pwd":
		c.reply(fmt.Sprintf("257 \"%s\" is the current directory.", c.cwd))
	case "cwd":
		dir := msgs[0]
		cwd(c, dir)
	case "list":
		var dir string
		if len(msgs) == 0 {
			dir = c.cwd
		} else {
			dir = msgs[0]
		}
		list(c, dir)
	case "mkd":
		dir := msgs[0]
		result := updateTree(c, dir)
		if result == "257" {
			c.reply(fmt.Sprintf("257 \"%s\" created.", path.Join(c.cwd, dir)))
		} else {
			c.reply(fmt.Sprintf("%s Create directory operation failed.", result))
		}
	case "quit":
		c.reply("221 bye")
		c.command.Close()
	}

}

func list(c *FTP, dir string) {
	c.reply("150 Here comes the directory listing.")

	children_list, exist := listChildren(c, dir)
	if exist {
		// Init names with . and ..
		names := []string{".", ".."}
		for _, node := range children_list {
			names = append(names, node.name)
		}

		reply_msg := strings.Join(names, "\t")
		c.data.Close()
		c.reply(reply_msg)
	}

}

func cwd(c *FTP, dir string) {
	if dir == ".." {
		if c.cwd != c.root {
			cwd := strings.Split(c.cwd, "/")
			c.cwd = strings.Join(cwd[0:len(cwd)-1], "/")
		}
		c.reply("250 Directory successfully changed.")
	} else {
		dir = getAbsolutePath(c.cwd, dir)
		_, pathExist := pathExists(c, dir)
		if pathExist {
			c.cwd = dir
			c.reply("250 Directory successfully changed.")
		} else {
			c.reply("550 Failed to change directory.")
		}
	}

}

func getAbsolutePath(cwd string, dir string) string {
	if dir[0] != '/' {
		dir = path.Join(cwd, dir)
	}
	return dir
}

func pathExists(c *FTP, dir string) (*treeNode, bool) {
	node := c.filepath
	dir = getAbsolutePath("", dir)
	dirList := strings.Split(dir, "/")
	for _, d := range dirList[1:] {
		if d == "" || d == "home" {
			continue
		}
		found := false
		for _, child := range node.children {
			if child.name == d {
				found = true
				node = child
				break
			}
		}
		if !found {
			return nil, false
		}
	}
	return node, true
}

func port(c *FTP, msgs []string) {
	nums := strings.Split(msgs[0], ",")

	// Extract the high and low bytes of the port number and convert them to integers
	highByte, _ := strconv.ParseInt(nums[4], 10, 32)
	lowByte, _ := strconv.ParseInt(nums[5], 10, 32)

	// Calculate the port number=
	port := highByte*256 + lowByte

	ip := strings.Join(nums[0:4], ".") + ":" + strconv.Itoa(int(port))

	// Attempt to establish a TCP connection to the specified IP address and port
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		// If connection fails, set passive mode and send an error reply
		c.pasv = true
		c.reply(fmt.Sprintf("501 Can't connect to a foreign address: %s", err))
	} else {
		// If connection succeeds, assign the connection to the data field of FTP and send a success reply

		c.data = conn
		c.reply("200 PORT command successful. Consider using PASV.")
	}
}

func listChildren(c *FTP, dir string) ([]*treeNode, bool) {
	dir = getAbsolutePath(c.cwd, dir)

	node, exist := pathExists(c, dir)
	children := []*treeNode{}
	if exist {
		children = node.children
	}
	return children, exist

}

func updateTree(c *FTP, dir string) string {
	dir = getAbsolutePath(c.cwd, dir)
	parDir := path.Dir(dir)

	parDirNode, exist := pathExists(c, parDir)

	if !exist {
		return "550"
	}

	newNode := &treeNode{
		name:     path.Base(dir),
		children: []*treeNode{},
	}
	parDirNode.children = append(parDirNode.children, newNode)
	return "257"
}
