package lib

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// 模拟namespace 网络命名空间隔离
// 只支持linux操作系统 不支持window、mac
// 还未完成、待完成

const self = "/proc/self/exe" // 在Linux中 代表当前执行的程序
const alpine = "/root/alpine" // 路径

var RootCmd = &cobra.Command{
	Use:   "usage",
	Short: "short description",
	Long:  `long description`,
	Run: func(cmd *cobra.Command, args []string) {
		// execute function
	},
}

var runCommand = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd := exec.Command(self, "exec", "/bin/sh")
		runCmd.SysProcAttr = &syscall.SysProcAttr{
			// fork 出一个进程
			// syscall.CLONE_NEWUTS 主机名、域名隔离
			// syscall.CLONE_NEWNS 挂在点隔离
			// syscall.CLONE_NEWUSER 用户隔离
			// syscall.CLONE_NEWPID 进程隔离
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID,
			UidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getuid(),
					Size:        1,
				},
			},
			GidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getgid(),
					Size:        1,
				},
			},
		}
		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		if err := runCmd.Start(); err != nil {
			log.Fatalln("err:", err)
		}
		runCmd.Wait()
	},
}

var execCommand = &cobra.Command{
	Use: "exec",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalln("error args")
		}
		runArgs := make([]string, 0)
		if len(args) > 1 {
			runArgs = args[1:]
		}

		CheckError("chroot:", syscall.Chroot(alpine))
		CheckError("chdir:", os.Chdir("/"))
		CheckError("mount proc error:", syscall.Mount("proc", "/proc", "proc", 0, ""))

		// log.Println("runArgs", runArgs)
		// log.Println("child uid:", os.Geteuid())
		// log.Println("child gid:", os.Getegid())
		runCmd := exec.Command(args[0], runArgs...)
		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Env = []string{ENV}
		if err := runCmd.Start(); err != nil {
			log.Fatal(err)
		}
		runCmd.Wait()
	},
}

func CheckError(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func Execute()  {
	RootCmd.AddCommand(runCommand, execCommand)

	if err := RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
