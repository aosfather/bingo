package utils

import (
	"os"
	//"syscall"
	"reflect"
)

func IsFileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

var os_Chown = os.Chown

func chown(name string, info os.FileInfo) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	f.Close()
	stat := info.Sys() //.(*syscall.Stat_t)
	//	if
	if reflect.TypeOf(stat).Elem().Name() == "Stat_t" {
		uid := reflect.ValueOf(stat).Elem().FieldByName("Uid").Uint()
		gid := reflect.ValueOf(stat).Elem().FieldByName("Gid").Uint()
		return os_Chown(name, int(uid), int(gid))
	}
	return nil
}