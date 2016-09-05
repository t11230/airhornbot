package perms

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/ramendb"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	permsCollName = "permslist"
)

type permCollection struct {
	*mgo.Collection
}

type userPermsCollection struct {
	*mgo.Collection
}

type PermsHandle struct {
	GuildID       string
	Namespace     string
	UserPermsColl userPermsCollection
}

type Perm struct {
	Name string
}

type UserPerm struct {
	UserID string
	Perms  []Perm `bson:",omitempty"`
}

func CreatePerm(perm string) error {
	permsColl := permCollection{ramendb.GetCollection("permsdb", permsCollName)}
	c, err := permsColl.Find(&Perm{Name: perm}).Count()
	if err != nil {
		log.Errorf("Error checking for perm: %v", err)
		return err
	}
	// Only create the perm if it doesn't exist
	if c > 0 {
		return nil
	}

	err = permsColl.Insert(&Perm{Name: perm})
	if err != nil {
		log.Error("Error creating perm: %v", err)
		return err
	}
	return nil
}

func PermExists(perm string) (bool, error) {
	permsColl := permCollection{ramendb.GetCollection("permsdb", permsCollName)}
	c, err := permsColl.Find(&Perm{Name: perm}).Count()
	if err != nil {
		log.Errorf("Error checking for perm: %v", err)
		return false, err
	}

	return (c > 0), nil
}

func PermsList() ([]Perm, error) {
	var result []Perm

	permsColl := permCollection{ramendb.GetCollection("permsdb", permsCollName)}
	err := permsColl.Find(nil).All(&result)
	if err != nil {
		log.Errorf("Error getting perm list: %v", err)
		return nil, err
	}

	return result, nil
}

func GetPermsHandle(guildId string, namespace string) *PermsHandle {
	userPermsCollName := "userperms"

	return &PermsHandle{
		GuildID:       guildId,
		Namespace:     namespace,
		UserPermsColl: userPermsCollection{ramendb.GetCollection(guildId, userPermsCollName)},
	}
}

func (h *PermsHandle) AddPerm(userId string, perm string) error {
	log.Debugf("Adding Perm %v to %v", perm, userId)

	user := &UserPerm{UserID: userId}
	result := &UserPerm{}

	err := h.UserPermsColl.Find(user).One(result)
	if err == mgo.ErrNotFound {
		user.Perms = []Perm{Perm{Name: perm}}
		h.UserPermsColl.Insert(user)
		return nil
	} else if err != nil {
		log.Errorf("Error finding UserPerm: %v", err)
		return err
	}

	result.Perms = append(result.Perms, Perm{Name: perm})

	err = h.UserPermsColl.Update(user, bson.M{"$set": result})
	if err != nil {
		return err
	}
	return nil
}

func (h *PermsHandle) RemovePerm(userId string, perm string) error {
	log.Debugf("Removing Perm %v from %v", perm, userId)

	hasPerm, err := h.CheckPerm(userId, perm)
	if err != nil {
		return err
	}
	if !hasPerm {
		return errors.New("User does not have that perm")
	}

	user := &UserPerm{UserID: userId}
	result := &UserPerm{}

	err = h.UserPermsColl.Find(user).One(result)
	if err != nil {
		log.Errorf("Error finding UserPerm: %v", err)
		return err
	}

	hasPerm, permIndex := permsContains(result.Perms, Perm{Name: perm})

	if !hasPerm {
		err = errors.New("User does not have perm")
		log.Error(err)
		return err
	}

	result.Perms = append(result.Perms[:permIndex], result.Perms[permIndex+1:]...)

	err = h.UserPermsColl.Remove(user)
	if err != nil {
		log.Errorf("Error removing user perms: %v", err)
		return err
	}

	if len(result.Perms) == 0 {
		return nil
	}

	err = h.UserPermsColl.Insert(user, result)
	if err != nil {
		log.Errorf("Error updating user perms: %v", err)
		return err
	}
	return nil
}

func (h *PermsHandle) CheckPerm(userId string, perm string) (bool, error) {
	log.Debugf("Checking perm %v for %v", perm, userId)

	result := &UserPerm{}

	err := h.UserPermsColl.Find(&UserPerm{UserID: userId}).One(result)
	if err == mgo.ErrNotFound {
		return false, nil
	} else if err != nil {
		log.Errorf("Error finding UserPerm: %v", err)
		return false, err
	}

	hasPerm, _ := permsContains(result.Perms, Perm{Name: perm})
	return hasPerm, nil
}

func permsContains(perms []Perm, find Perm) (bool, int) {
	for i, item := range perms {
		if item.Name == find.Name {
			return true, i
		}
	}
	return false, -1
}
