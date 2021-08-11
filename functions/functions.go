package functions

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path"

	"github.com/gogamic/gogamic-ci-cli/utils"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

func CheckIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		// fmt.Printf("IP Address: %s - Invalid\n", ip)
		return errors.New("invalid ip address")
	} else {
		// fmt.Printf("IP Address: %s - Valid\n", ip)
		return nil
	}
}

func RunServerCommands(server_ip string, commands []string, ssh_key_path *string) error {
	key, err := getKeyFile(ssh_key_path)

	if err != nil {
		panic(err)
	}

	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", server_ip), config)
	if err != nil {
		return fmt.Errorf("Failed to connect: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: " + err.Error())
	}

	in, err := session.StdinPipe()
	if err != nil {
		HandleErr(err, "unable to read stdin")
	}
	out, err := session.StdoutPipe()
	if err != nil {
		HandleErr(err, "unable to read stdout")
	}
	err_log, err := session.StderrPipe()
	if err != nil {
		HandleErr(err, "unable to read stderr")
	}

	// session.Stderr = out // this will send stderr to the same pipe
	err = session.Shell()

	if err != nil {
		HandleErr(err, "Unable to open a shell")
	}

	for _, cmd := range commands {
		// log.Printf("running command => %s", cmd)
		fmt.Fprintf(in, "%s;", cmd)
	}
	// fmt.Printf("%T", out)
	// fmt.Printf("out: %v\n", out)

	defer session.Close()
	defer readoutput(out, err_log)
	defer in.Close()

	return nil
}

func getKeyFile(ssh_key_path *string) (key ssh.Signer, err error) {
	var file string
	if ssh_key_path != nil {
		cwd, _ := os.Getwd()
		file = path.Join(cwd, *ssh_key_path)
	} else {
		usr, _ := user.Current()
		file = path.Join(usr.HomeDir, "/.ssh/id_rsa")
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

func GetCommands(backend *utils.Config) ([]string, error) {
	var backend_name = backend.Config.Backend
	if backend_name == "dokku" {
		commands := []string{
			fmt.Sprintf("sudo docker pull %s", backend.Config.Image),
			fmt.Sprintf("sudo docker tag %s dokku/%s", backend.Config.Image, backend.Name),
			fmt.Sprintf("sudo dokku git:from-image %s dokku/%s:latest", backend.Name, backend.Name),
		}
		return commands, nil
	} else if backend_name == "docker" {
		commands := []string{"", "", "", ""}
		return commands, nil
	} else if backend_name == "test" {
		commands := []string{"whoami", "pwd"}
		return commands, nil
	} else {
		return nil, fmt.Errorf("%s is not a valid backend", backend_name)
	}

}

func readoutput(out io.Reader, err_log io.Reader) {
	fmt.Print("Started to read! \n")
	out_buf := new(bytes.Buffer)
	out_buf.ReadFrom(out)
	fmt.Printf("out: %s", out_buf.String())

	err_buf := new(bytes.Buffer)
	err_buf.ReadFrom(err_log)
	err_buf_string := err_buf.String()
	if len(err_buf_string) > 0 {
		log.Print(err_buf_string)
	}
}

func HandleErr(err error, message string) {
	log.Panicf("%s: %s", message, err.Error())
}

func ValidateYAMLFile(file_path string) error {
	data, err := ParseYAMLFile(file_path)
	if err != nil {
		return fmt.Errorf("unable to parse file: %s", err.Error())
	}

	err = CheckIPAddress(data.Config.IP)
	if err != nil {
		return fmt.Errorf("invalid ip address: %s", err.Error())
	}
	return nil
}

func ParseYAMLFile(file_path string) (*utils.Config, error) {
	var config utils.Config
	current_dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	var config_path = path.Join(current_dir + fmt.Sprintf("/%s", file_path))

	yamlFile, err := ioutil.ReadFile(config_path)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return &config, nil
}
