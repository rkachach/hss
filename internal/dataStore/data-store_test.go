package dataStore

import (
	"math/rand"
	"os"
	"testing"
	"github.com/rkachach/hss/cmd/config"
)

var store OsFileSystem = OsFileSystem{}


func createFile(filePath string, t *testing.T) FileInfo {
  info, err := store.StartFileUpload(filePath, map[string]string{})
  if err != nil {
    t.FailNow()
  }
  return info
}

func TestReadFile(t *testing.T) {
  filePath := "foofile"
  info, err := store.StartFileUpload(filePath, map[string]string{})
  if err != nil {
    t.FailNow()
  }
  defer store.DeleteFile(filePath)
  if info.Size != 0 {
    t.Error("Size wasn't initialized to 0")
  }

  data := [10][100]byte{}
  // Fill data with parts to write
  for i := 0; i < 10; i++ {
    for j := 0; j < 100; j++ {
      data[i][j] = byte(rand.Int31())
    }
    info, err = store.WriteFilePart(filePath, data[i][:], 0)
    if err != nil {
      t.Error("Error writing file ", err)
    }
  }

  readData, err := store.ReadFile(filePath)
  if err != nil {
    t.Error("Error writing file ", err)
  }
  for i := 0; i < 10; i++ {
    for j := 0; j < 100; j++ {
      if data[i][j] != readData[i*100 + j] {
        t.Errorf("Failed reading byte number=%d, block=%d, col=%d", i*100+j, i, j)
      }
    }
  }
}

func TestDeleteFile(t *testing.T){
  filePath := "foofile"
  createFile(filePath, t)
  defer store.DeleteFile(filePath)

  err := store.DeleteFile(filePath)
  if err == nil {
    t.FailNow()
  }
}

func TestStartFileUpload(t *testing.T) {
  filePath := "foofile"
  info, err := store.StartFileUpload(filePath, map[string]string{})
  if err != nil {
    t.FailNow()
  }
  defer store.DeleteFile(filePath)
  if info.Size != 0 {
    t.Error("Size wasn't initialized to 0")
  }
}

func TestWriteFilePart(t *testing.T) {
  filePath := "foofile"
  info, err := store.StartFileUpload(filePath, map[string]string{})
  if err != nil {
    t.FailNow()
  }
  defer store.DeleteFile(filePath)
  if info.Size != 0 {
    t.Error("Size wasn't initialized to 0")
  }

  data := []byte{1, 2, 3}
  info, err = store.WriteFilePart(filePath, data, 0)
  if err != nil {
    t.Error("Error writing file ", err)
  }
  expectedSize := 3
  if info.Size != int64(expectedSize) {
    t.Errorf("File size is not ok expected=%d got=%d", expectedSize, info.Size)
  }


}

func TestReadFileInfo(t *testing.T) {}
func TestUpdateFileInfo(t *testing.T) {}

func TestDeleteDirectory(t *testing.T) {
  err := store.CreateDirectory("testdir", map[string]string{})
  if err != nil {
    t.Error("Error creating directory")
  }
  err = store.DeleteDirectory("testdir")
  if err != nil {
    t.Error("Error removing dir")
  }
}

func TestCreateDirectory(t *testing.T) {
  err := store.CreateDirectory("testdir", map[string]string{})
  if err != nil {
    t.Error("Error creating directory")
  }
  defer store.DeleteDirectory("testdir")
}
func TestGetDirectoryInfo(t *testing.T) {
  metadata := map[string]string{}
  metadata["foo"] = "var"
  err := store.CreateDirectory("testdir", metadata)
  if err != nil {
    t.Error("Error creating directory")
  }
  defer store.DeleteDirectory("testdir")

  dirInfo, err := store.GetDirectoryInfo("testdir")
  if dirInfo.Name != "testdir" {
    t.Error("Wrong name")
  }
  if dirInfo.Metadata["foo"] != "var" {
    t.Error("Wrong metadata")
  }
}

func TestListDirectory(t *testing.T) {
  t.Skip()
  err := store.CreateDirectory("testdir", map[string]string{})
  if err != nil {
    t.Error("Error creating directory")
  }
  defer store.DeleteDirectory("testdir")
  defer store.DeleteDirectory("testdir")

  createFile("testdir/a", t)
  store.WriteFilePart("testdir/a", []byte{}, 0)
  createFile("testdir/b", t)
  store.WriteFilePart("testdir/b", []byte{}, 0)
  createFile("testdir/c", t)
  store.WriteFilePart("testdir/c", []byte{}, 0)

  files, err := store.ListDirectory("testdir")
  if err != nil {
    t.Error("Error listing files")
  }

  if len(files) != 3 {
    println(files)
    t.Errorf("Got wrong number of files %d!=3", len(files))
  }

  if files[0] != "testdir/a" || files[1] != "testdir/b" || files[2] != "testdir/c"{
    t.Error("Got different file names")
  } 
}

func TestIsMetadataFile(t *testing.T) {}


func TestMain(m *testing.M) {
  println(os.Getwd())
  // hacky I know, I don't want to deal with go right now
  config.ReadConfig("../../config/config.json")
  config.InitLogger()

  exitCode := m.Run()

  os.Exit(exitCode)
}
