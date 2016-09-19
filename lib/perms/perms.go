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
	permQuery := permsColl.Find(&Perm{Name: perm})

	c, err := permQuery.Count()
	if err != nil {
		log.Errorf("Error checking for perm: %v", err)
		return err
	}
	// Only create the perm if it doesn't exist
	if c > 0 {
		return nil
	}

	result := &Perm{Name: perm}

	err = permsColl.Insert(result)
	if err != nil {
		log.Error("Error creating perm: %v", err)
		return err
	}
	return nil
}

func PermExists(perm string) (*Perm, error) {
	permsColl := permCollection{ramendb.GetCollection("permsdb", permsCollName)}
	permQuery := permsColl.Find(&Perm{Name: perm})
	c, err := permQuery.Count()
	if err != nil {
		log.Errorf("Error checking for perm: %v", err)
		return nil, err
	}

	result := &Perm{}

	if c > 0 {
		err = permQuery.One(result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, nil
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

func GetPermsHandle(guildId string) *PermsHandle {
	userPermsCollName := "userperms"

	return &PermsHandle{
		GuildID:       guildId,
		UserPermsColl: userPermsCollection{ramendb.GetCollection(guildId, userPermsCollName)},
	}
}

func (h *PermsHandle) AddPerm(userId string, perm *Perm) error {
	log.Debugf("Adding Perm %v to %v", perm, userId)

	if perm == nil {
		err := errors.New("Perm was nil")
		log.Error(err)
		return err
	}

	hasPerm := h.CheckPerm(userId, perm)
	if hasPerm {
		return errors.New("User already has that perm")
	}

	user := &UserPerm{UserID: userId}
	result := &UserPerm{}

	err := h.UserPermsColl.Find(user).One(result)
	if err == mgo.ErrNotFound {
		user.Perms = []Perm{*perm}
		h.UserPermsColl.Insert(user)
		return nil
	} else if err != nil {
		log.Errorf("Error finding UserPerm: %v", err)
		return err
	}

	result.Perms = append(result.Perms, *perm)

	err = h.UserPermsColl.Update(user, bson.M{"$set": result})
	if err != nil {
		return err
	}
	return nil
}

func (h *PermsHandle) RemovePerm(userId string, perm *Perm) error {
	log.Debugf("Removing Perm %v from %v", perm, userId)

	if perm == nil {
		err := errors.New("Perm was nil")
		log.Error(err)
		return err
	}

	hasPerm := h.CheckPerm(userId, perm)
	if !hasPerm {
		return errors.New("User does not have that perm")
	}

	user := &UserPerm{UserID: userId}
	result := &UserPerm{}

	err := h.UserPermsColl.Find(user).One(result)
	if err != nil {
		log.Errorf("Error finding UserPerm: %v", err)
		return err
	}

	hasPerm, permIndex := permsContains(result.Perms, perm)

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

func (h *PermsHandle) CheckPerm(userId string, perm *Perm) bool {
	log.Debugf("Checking perm %v for %v", perm, userId)

	result := &UserPerm{}

	err := h.UserPermsColl.Find(&UserPerm{UserID: userId}).One(result)
	if err == mgo.ErrNotFound {
		return false
	} else if err != nil {
		log.Errorf("Error finding UserPerm: %v", err)
		return false
	}

	hasPerm, _ := permsContains(result.Perms, perm)
	return hasPerm
}

func permsContains(perms []Perm, find *Perm) (bool, int) {
	if find == nil {
		return false, -1
	}
	for i, item := range perms {
		if item.Name == find.Name {
			return true, i
		}
	}
	return false, -1
}
