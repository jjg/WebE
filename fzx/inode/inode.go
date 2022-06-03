package inode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jjg/WebE/fzx/utils"
)

type Inode struct {
	FzxPath          string
	Fingerprint      string
	StorageLocation  *string
	Created          time.Time // TODO: This should probably get split into created/updatd
	Version          int
	Private          bool
	Encrypted        bool
	AccessKey        string
	ContentType      string
	FileSize         int
	BlockSize        int64
	BlocksReplicated int
	InodeReplicated  int
	Blocks           []string
	// New fzx properties
	Status int
}

func (i *Inode) Save() error {

	var err error
	var inodeJson []byte

	i.Created = time.Now()
	i.Fingerprint = utils.StringToSha1(i.FzxPath)

	// Write the contents of this inode to storage.
	inodeJson, err = json.Marshal(i)
	inodeFilename := fmt.Sprintf("%v/%v.json", *i.StorageLocation, i.Fingerprint)

	// TODO: See if this is the best way to write a file.
	// TODO: Do some error handling around this write.
	ioutil.WriteFile(inodeFilename, inodeJson, 0644)

	return err
}

func (i *Inode) Load(storageLocation *string, fzxPath string) error {
	var err error
	var inodeJson []byte

	i.StorageLocation = storageLocation
	i.FzxPath = fzxPath
	i.Fingerprint = utils.StringToSha1(i.FzxPath)

	inodeFilename := fmt.Sprintf("%v/%v.json", *i.StorageLocation, i.Fingerprint)
	if inodeJson, err = ioutil.ReadFile(inodeFilename); err != nil {
		// TODO: Inspect err to set Status more accurately.
		i.Status = http.StatusNotFound
		return err
	}

	if err = json.Unmarshal([]byte(inodeJson), i); err != nil {
		return err
	}

	i.Status = http.StatusOK

	return err
}

func (i *Inode) JsonString() (string, error) {
	if inodeJson, err := json.Marshal(i); err != nil {
		return "", err
	} else {
		return string(inodeJson), err
	}
}
