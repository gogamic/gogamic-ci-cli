package functions

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path"

	"golang.org/x/crypto/ssh"
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

func RunServerCommands(server_ip string, commands []string, ssh_key_path *string, b64 bool, b64_key string) {
	key, err := getKeyFile(ssh_key_path, b64, b64_key)

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
		HandleErr(err, "Failed to connect")
	}

	session, err := client.NewSession()
	if err != nil {
		HandleErr(err, "Failed to create session")
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
}

func getKeyFile(ssh_key_path *string, b64 bool, b64_key string) (key ssh.Signer, err error) {
	var buf []byte
	if b64_key != "null" {
		buf = []byte(b64_key)
		b64 = true
	} else {
		var file string
		if ssh_key_path != nil {
			cwd, _ := os.Getwd()
			file = path.Join(cwd, *ssh_key_path)
		} else {
			usr, _ := user.Current()
			file = path.Join(usr.HomeDir, "/.ssh/id_rsa")
		}

		buf, err = ioutil.ReadFile(file)
		if err != nil {
			HandleErr(err, "ssh key doesn't exists")
		}
	}

	if b64 {
		buf, err = Base64Decode(buf)
		if err != nil {
			HandleErr(err, "unable to parse ssh key in base 64 decoding")
		}
		/* fmt.Println("successfully decoded")
		fmt.Println(string(buf)) */
	}

	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

func GetCommands(backend_name string, backend_image string) ([]string, error) {
	if backend_name == "dokku" {
		commands := []string{
			fmt.Sprintf("sudo docker pull %s", backend_image),
			fmt.Sprintf("sudo docker tag %s dokku/%s", backend_image, backend_name),
			fmt.Sprintf("sudo dokku git:from-image %s dokku/%s:latest", backend_name, backend_name),
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
	log.Fatalf("%s: %s", message, err.Error())
}

/* func ValidateYAMLFile(file_path string) error {
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
} */

func Base64Decode(message []byte) (b []byte, err error) {
	var l int
	b = make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, err = base64.StdEncoding.Decode(b, message)
	if err != nil {
		return
	}
	return b[:l], nil
}
