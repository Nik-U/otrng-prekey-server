package prekeyserver

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

// Design:
// - a production level storage mechanism that stores data into separate files
// - the storage is very simple:
// - each from will be hashed into a hex-string using sha256
//   if the path is /var/data/otrng-prekey-storage
//   a specific data point will be 1234567ABCDEFADBDCDCDCDCDCDAAAAAA123538892322332ABDCDFFFFACDADDD
// - all information about it will be stored in:
// - /var/data/otrng-prekey-storage/1234/1234567ABCDEFADBDCDCDCDCDCDAAAAAA123538892322332ABDCDFFFFACDADDD
// - inside of this directory will be a .lock file if it's currently used
// - it will have one file for each instance tag, for example:
//   - for instance tag 0x1234AABB, will be directory
//     /var/data/otrng-prekey-storage/1234/1234567ABCDEFADBDCDCDCDCDCDAAAAAA123538892322332ABDCDFFFFACDADDD/1234AABB
//   - inside of this directory will be two files, called cp.bin and pp.bin, containing the binary representation of
//     the client profile and prekey profile, if they exist
//   - there will be a directory called pm, where each pm entry will be of the form /1245AABB.bin, containing the raw
//     binary representation of the identifier.
//   in each prefix directory, like /var/data/otrng-prekey-storage/1234, there can be a .lock file as well. This is only there
//   when the directory is about to be added to, or when the directory will be removed
//   same as the top level directory, if /var/data/otrng-prekey-storage/1234 will be removed, we will create a .lock file, and remove the directory and then remove the .lock file
//   - same thing when we create the directory.
//   - The cleanup method is in charge of deleting empty directories - we don't care about that in the main path

type fileStorageFactory struct {
	desc string
}

func (fsf *fileStorageFactory) createStorage() storage {
	return createFileStorageFrom(fsf.desc)
}

type fileStorage struct {
	path string
}

func isFileStorageDescriptor(desc string) bool {
	return strings.HasPrefix(desc, "dir:")
}

func createFileStorageFrom(descriptor string) *fileStorage {
	// TODO: this should return an error if the directory doesn't exist
	path := strings.TrimPrefix(descriptor, "dir:")
	return &fileStorage{
		path: path,
	}
}

func (fs *fileStorage) writeData(user, file string, itag uint32, data []byte) error {
	userDir := fs.getOrCreateDirFor(user)
	fs.lock(userDir)
	defer fs.unlock(userDir)

	itagDir := fs.getOrCreateInstanceTagDir(userDir, itag)
	return ioutil.WriteFile(path.Join(itagDir, file), data, 0600)
}

func (fs *fileStorage) storeClientProfile(user string, cp *clientProfile) error {
	return fs.writeData(user, "cp.bin", cp.instanceTag, cp.serialize())
}

func (fs *fileStorage) storePrekeyProfile(user string, pp *prekeyProfile) error {
	return fs.writeData(user, "pp.bin", pp.instanceTag, pp.serialize())
}

func formatUint32(v uint32) string {
	return fmt.Sprintf("%08X", v)
}

func entryExists(entry string) bool {
	_, err := os.Stat(entry)
	return err == nil
}

// TODO: maybe add pid to the lock file name?

// lock will place a lock file in the named directory
//  if there is already a lock file there, it will wait
//  until it's been removed
func (fs *fileStorage) lock(dirName string) {
	lockFile := path.Join(dirName, ".lock")
	for entryExists(lockFile) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
	ioutil.WriteFile(lockFile, []byte{}, 0600)
}

func (fs *fileStorage) unlock(dirName string) {
	os.Remove(path.Join(dirName, ".lock"))
}

func (fs *fileStorage) composeDirNameFor(user string) (string, string) {
	hex := fmt.Sprintf("%X", sha256.Sum256([]byte(user)))
	first := path.Join(fs.path, hex[0:4])
	return first, path.Join(first, hex)
}

func (fs *fileStorage) getDirFor(user string) (string, bool) {
	_, nm := fs.composeDirNameFor(user)
	return nm, entryExists(nm)
}

func (fs *fileStorage) getInstanceTagDir(userDir string, itag uint32) string {
	return path.Join(userDir, formatUint32(itag))
}

func (fs *fileStorage) getPmDir(itagDir string) string {
	return path.Join(itagDir, "pm")
}

func (fs *fileStorage) getOrCreateDirFor(user string) string {
	dir, ok := fs.getDirFor(user)
	if ok {
		return dir
	}

	pref, us := fs.composeDirNameFor(user)
	if !entryExists(pref) {
		fs.lock(fs.path)
		os.Mkdir(pref, 0600)
		fs.unlock(fs.path)
	}
	fs.lock(pref)
	os.Mkdir(us, 0600)
	fs.unlock(pref)
	return us
}

// getOrCreateInstanceTagDir assumes that the user dir is already locked
func (fs *fileStorage) getOrCreateInstanceTagDir(userDir string, itag uint32) string {
	name := path.Join(userDir, formatUint32(itag))
	if !entryExists(name) {
		os.Mkdir(name, 0600)
	}
	return name
}

// getOrCreatePmDir assumes that the user dir is already locked
func (fs *fileStorage) getOrCreatePmDir(userDir string) string {
	name := path.Join(userDir, "pm")
	if !entryExists(name) {
		os.Mkdir(name, 0600)
	}
	return name
}

func (fs *fileStorage) storePrekeyMessages(user string, pms []*prekeyMessage) error {
	userDir := fs.getOrCreateDirFor(user)
	fs.lock(userDir)
	defer fs.unlock(userDir)

	itagDir := fs.getOrCreateInstanceTagDir(userDir, pms[0].instanceTag)
	pmDir := fs.getOrCreatePmDir(itagDir)
	for _, pm := range pms {
		if e := ioutil.WriteFile(path.Join(pmDir, formatUint32(pm.identifier)+".bin"), pm.serialize(), 0600); e != nil {
			return e
		}
	}
	return nil
}

func (fs *fileStorage) numberStored(user string, itag uint32) uint32 {
	userDir, ok := fs.getDirFor(user)
	if !ok {
		return 0
	}
	fs.lock(userDir)
	defer fs.unlock(userDir)

	pmDir := fs.getPmDir(fs.getInstanceTagDir(userDir, itag))
	files, err := ioutil.ReadDir(pmDir)
	if err != nil {
		return 0
	}
	count := uint32(0)
	for _, f := range files {
		if !f.IsDir() && path.Ext(f.Name()) == ".bin" && len(f.Name()) == 12 {
			count++
		}
	}

	return count
}

func isUint32Hex(name string) bool {
	res, _ := regexp.MatchString("^[0-9A-F]{8}$", name)
	return res
}

func (fs *fileStorage) retrieveFor(user string) []*prekeyEnsemble {
	userDir, ok := fs.getDirFor(user)
	if !ok {
		return nil
	}
	fs.lock(userDir)
	defer fs.unlock(userDir)

	files, err := ioutil.ReadDir(userDir)
	if err != nil {
		return nil
	}

	entries := []*prekeyEnsemble{}

	for _, f := range files {
		if f.IsDir() && isUint32Hex(f.Name()) {
			itagDir := path.Join(userDir, f.Name())
			cpFile := path.Join(itagDir, "cp.bin")
			ppFile := path.Join(itagDir, "pp.bin")
			pmDir := path.Join(itagDir, "pm")
			if entryExists(cpFile) &&
				entryExists(ppFile) &&
				entryExists(pmDir) {
				pmFiles, _ := ioutil.ReadDir(pmDir)
				if len(pmFiles) > 0 {
					pmFile := path.Join(pmDir, pmFiles[0].Name())
					pm, e1 := ioutil.ReadFile(pmFile)
					cp, e2 := ioutil.ReadFile(cpFile)
					pp, e3 := ioutil.ReadFile(ppFile)
					if e1 == nil && e2 == nil && e3 == nil {
						pmR := &prekeyMessage{}
						cpR := &clientProfile{}
						ppR := &prekeyProfile{}
						_, ok1 := pmR.deserialize(pm)
						_, ok2 := cpR.deserialize(cp)
						_, ok3 := ppR.deserialize(pp)
						if ok1 && ok2 && ok3 {
							defer os.Remove(pmFile)
							entries = append(entries, &prekeyEnsemble{
								cp: cpR,
								pp: ppR,
								pm: pmR,
							})
						}
					}
				}
			}
		}
	}

	return entries
}

func (fs *fileStorage) cleanup() {
	// For each user:
	//    for each instance tag:
	//       - check if the client profile has expired - if so, remove it
	//       - check if the prekey profile has expired - if so, remove it
	//       - check if the prekey message directory is empty - if so, remove it
	//       - check if the instance tag directory is empty - if so, remove it
	//    - check if the user directory has any instance tags - if not, remove it
	// Check the prefix directory for the user - if it's empty, remove it
	// TODO: implement
}
