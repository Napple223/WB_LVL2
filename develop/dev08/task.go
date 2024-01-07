package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Структура процесса
type process struct {
	name string
	*os.Process
}

// Струкура команды.
type cmd struct {
	name string
	args []string
}

// Метод для запуска исполнения записанных команд.
func (c cmd) run(pss map[int]process) (string, error) {
	switch c.name {
	case "leave":
		os.Exit(0)
	case "cd":
		return "", c.cd()
	case "pwd":
		return os.Getwd()
	case "echo":
		return c.echo(), nil
	case "kill":
		return "", c.kill(pss)
	case "ps":
		return getPS(pss), nil
	case "exec":
		return "", c.execute(pss)
	default:
		return "", errors.New(c.name + ": cmd not found")
	}
	return "", nil
}

// Метод для смены директории.
func (c cmd) cd() error {
	path := ""
	switch {
	case len(c.args) == 0:
		p, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = p
	case strings.HasPrefix(c.args[0], "~"):
		p, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = p + c.args[0][len("~"):]
	default:
		path = c.args[0]
	}
	return os.Chdir(path)
}

// Метод - эхо. Что написал в консоль, то и получил обратно.
func (c cmd) echo() string {
	return strings.Join(c.args, " ")
}

// Метод для убийства созданных процессов.
func (c cmd) kill(pss map[int]process) error {
	if len(c.args) == 0 {
		return errors.New("no pid specified")
	}
	ps := make([]process, len(c.args))
	var id int
	for idx, p := range c.args {
		pid, err := strconv.Atoi(p)
		if err != nil {
			return errors.New("wrong pid")
		}
		id = pid
		pr, ok := pss[pid]
		if !ok {
			return errors.New("pid №" + c.args[idx] + "not found")
		}
		ps[idx] = pr
	}

	for _, p := range ps {
		err := p.Kill()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		delete(pss, id)
	}
	return nil
}

// Функция для получения списка всех запущенных процессов.
func getPS(pss map[int]process) string {
	var b strings.Builder
	b.WriteString("PID\tCMD")
	for pid, pr := range pss {
		s := strconv.Itoa(pid)
		b.WriteString("\n" + s + "\t" + pr.name)
	}
	return b.String()
}

// Метод для создания процесса.
// Как я понял в го не нужно вызывать отдельно fork.
func (c cmd) execute(pss map[int]process) error {
	if len(c.args) == 0 {
		return errors.New("no command specified")
	}
	cm := exec.Command(c.args[0], c.args[1:]...)
	cm.Stdin = os.Stdin
	cm.Stdout = os.Stdout
	cm.Stderr = os.Stderr
	err := cm.Run()
	if err != nil {
		return err
	}
	p, err := os.FindProcess(cm.Process.Pid)
	if err != nil {
		return err
	}
	pss[cm.Process.Pid] = process{c.args[0], p}
	return nil
}

// Функция для распечатывания текущей директории в консоли.
func shellPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(path, homeDir) {
		path = "~" + path[len(homeDir):]
	}
	return path, nil
}

// Функция для парсинга переданных в консоли аргументов.
// Поддердавает пайпы.
func parseInput(input string) ([]cmd, error) {
	pipe := strings.FieldsFunc(input, func(r rune) bool {
		return r == '|'
	})
	if len(pipe) == 0 {
		return nil, errors.New("unexpected token")
	}
	c := []cmd{}
	for _, p := range pipe {
		cmdS := strings.Fields(p)
		switch len(cmdS) {
		case 0:
			return nil, errors.New("unexpected token")
		case 1:
			com := cmd{name: cmdS[0]}
			c = append(c, com)
		default:
			com := cmd{
				name: cmdS[0],
				args: cmdS[1:],
			}
			c = append(c, com)
		}
	}
	return c, nil
}

// Функция, ставящая в очередь на выполнение команды
// из пайпа.
func execPipe(pipe []cmd, pss map[int]process) error {
	for _, p := range pipe {
		out, err := p.run(pss)
		if err != nil {
			return err
		}
		if out != "" {
			fmt.Println(out)
		}
	}
	return nil
}

func main() {
	//Мап для хранения процессов.
	pss := make(map[int]process)
	//Получаем pid процесса.
	pid := os.Getpid()
	//Получаем процесс.
	p, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//Пишем процесс в мап процессов.
	pss[pid] = process{"shell", p}
	//Бесконечно читаем stdin, пока не напишем leave.
	reader := bufio.NewScanner(os.Stdin)
	for {
		path, err := shellPath()
		switch {
		case err != nil:
			fmt.Print("$ ")
		case err == nil:
			fmt.Print(path + "$ ")
		}
		fmt.Println("for exit print leave")
		if reader.Scan() {
			input := reader.Text()
			pipe, err := parseInput(input)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			err = execPipe(pipe, pss)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
	}
}
